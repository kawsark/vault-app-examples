# vault-app-examples
An example vault client that reads KV2
**Note: this is an example client only and does not indicate best practices in reading secrets from Vault. Please do not use this app in production.


## Vault setup
```
# Setup environment
export VAULT_ADDR="http://vault-gke:8200"
export VAULT_TOKEN="root-or-admin-token"

# Setup KV and policy
cd go-kv2/
vault secrets enable -path=kv2 -version=2 kv
vault kv put kv2/api_token API_TOKEN=v1.abcd
vault kv put kv2/api_token API_TOKEN=v2.efgh
vault kv put kv2/api_token API_TOKEN=v3.ijkl
vault policy write app1 app1-policy.hcl
```

## Building and Running the App
```
export VAULT_TOKEN=$(vault token create -format=json -policy=app1 | jq -r .auth.client_token)
export SECRET_PATH=kv2/data/api_token
export SECRET_VERSION=2
export SECRET_KEY=API_TOKEN

go build -o vault-app .
./vault-app
```

## Example output
```
3289488Z deletion_time: destroyed:false version:2]]
Warnings: %v []
~~~~~~~~~~~~~~~~~~~~~~~~~~
~~~~~ Printing  Data ~~~~~
map[API_TOKEN:v2.efgh]
~~~~~~~~~~~~~~~~~~~~~~~~~~
~~~~~ Printing value for Key: API_TOKEN ~~~~~
v2.efgh
~~~~~~~~~~~~~~~~~~~~~~~~~~
2020/02/06 01:19:09 Starting renewal loop
2020/02/06 01:19:09 secret is not renewable
```