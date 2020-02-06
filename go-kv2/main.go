package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
	"github.com/hashicorp/vault/api"
)

func readSecretWithVersion(client *api.Client, pathName string, keyName string, version string) *api.Secret {
  
  //Setup the data map: map[string][]string
  var v map[string][]string
  v = make(map[string][]string) //Initialize map
  var verArray []string
  verArray = make([]string,1) //Initialize array
  verArray[0] = version
  v["version"] = verArray

  //Read secret
  secretValues, err := client.Logical().ReadWithData(pathName,v)
  if err != nil {
	fmt.Println(err)
  }
  log.Printf("secret %s -> %v", pathName, secretValues)

  fmt.Println("~~~~~ Printing individual Struct fields ~~~~~")
  fmt.Println("RequestID: %s", secretValues.RequestID)
  fmt.Println("LeaseID: %s", secretValues.LeaseID)
  fmt.Println("Data: %v", secretValues.Data)
  fmt.Println("Warnings: %v", secretValues.Warnings)  
  fmt.Println("~~~~~~~~~~~~~~~~~~~~~~~~~~")

  fmt.Println("~~~~~ Printing  Data ~~~~~")
  var data = secretValues.Data["data"]
  fmt.Println(data)
  fmt.Println("~~~~~~~~~~~~~~~~~~~~~~~~~~")

  fmt.Printf("~~~~~ Printing value for Key: %s ~~~~~\n",keyName)
  data = secretValues.Data["data"]
  var dataassert = data.(map[string]interface{})
  fmt.Println(dataassert[keyName])
  fmt.Println("~~~~~~~~~~~~~~~~~~~~~~~~~~")

  return secretValues
}

func main() {
	config := api.DefaultConfig()
	vaultClient, err := api.NewClient(config)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("==> WARNING: Don't ever write secrets to logs.")
	log.Println("==>          This is for demonstration only.")
	log.Println(vaultClient.Token())

	//Lookup SECRET_PATH Environment variable
	pathName, set := os.LookupEnv("SECRET_PATH")
	if !set {
		pathName = "secret/creds"
	}
	fmt.Printf("Using path=%s", pathName)

	//Lookup SECRET_KEY Environment variable
	keyName, set := os.LookupEnv("SECRET_KEY")
	if !set {
		keyName = "foo"
	}
	fmt.Printf("Using key=%s", keyName)

	//Lookup SECRET_VERSION Environment variable
	versionName, set := os.LookupEnv("SECRET_VERSION")
	if !set {
		versionName = "1"
	}
	fmt.Printf("Using path=%s, key=%s, version=%s", pathName, keyName, versionName)

	//Read secret
	var secret = readSecretWithVersion(vaultClient, pathName, keyName, versionName)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	// Keep token renewed
	renewer, err := vaultClient.NewRenewer(&api.RenewerInput{
		Secret: secret,
		Grace:  1 * time.Second,
	})
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Starting renewal loop")
	go renewer.Renew()
	defer renewer.Stop()

	for {
		select {
		case err := <-renewer.DoneCh():
			if err != nil {
				log.Fatal(err)
			}
		case renewal := <-renewer.RenewCh():
			log.Printf("Successfully renewed: %#v", renewal)
			//Read secret
			readSecretWithVersion(vaultClient, pathName, keyName, versionName)
		case <-quit:
			log.Fatal("Shutdown signal received, exiting...")
		}
	}

}
