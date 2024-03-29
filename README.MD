Vault - Client for GO language, this dependency support only HashiCorp - Vault k/v 2.0

## Installation
```shell
go get -u github.com/rkritchat/vault-client

```

## Usage
Add Vault information to your configuration. (`.env`) file:

```sh
VAULT.URL=http://127.0.0.1:8200   #Vault URL
VAULT.PATH=test-vault-client      #Secret path
VAULT.TOKEN=s.xxxxxxXxxxx         #Vault token for access kv
VAULT.DISABLE=false               #Optional make it "true" for disable vault-client and read configuration from env file instead
```

List Vault Keys/Values by using vault kv get ${secret_path}

```sh
/ # vault kv get secret/test-vault-client
.
.
.
======== Data ========
Key              Value
---              -----
database.pwd     root
database.url     vault-client
database.user    change
```


Code example

```go

package main

import (
	"fmt"
	"github.com/joho/godotenv"
	"github.com/rkritchat/vault-client/pkg/conf"
	"github.com/rkritchat/vault-client/pkg/vault"
	"os"
)

//Create Configuration struct, Tag conf is ENV key that will set to envrionment variable.
//make sure that you have "database.url", "database.user" and "database.pwd" in Vault
type Configuration struct {
	DatabaseUrl string `json:"database.url" conf:"DATABASE.URL"`
	DatabaseUser string `json:"database.user" conf:"DATABASE.USER"`
	DatabasePwd string `json:"database.pwd" conf:"DATABASE.PWD"`
}

func main(){
	_ = godotenv.Load() //read VAULT.URL | VAULT.PATH | VAULT.TOKEN from the .evn file
	err := vault.NewVault(Configuration{}).Load()
	if err!=nil{
	    log.Fatal(err)
	}
	fmt.Printf("%v | %v | %v", os.Getenv("DATABASE.URL"), os.Getenv("DATABASE.USER"),os.Getenv("DATABASE.PWD"))
}

```

Vault - Client support reload configuration from Vault.

```go
package main

import (
	"github.com/go-chi/chi"
	"github.com/joho/godotenv"
	"github.com/rkritchat/vault-client/pkg/conf"
	"github.com/rkritchat/vault-client/pkg/vault"
	"net/http"
)

//Create Configuration struct, Tag conf is ENV key that will set to envrionment variable.
//make sure that you have "database.url", "database.user" and "database.pwd" in Vault
type Configuration struct {
	DatabaseUrl string `json:"database.url" conf:"DATABASE.URL"`
	DatabaseUser string `json:"database.user" conf:"DATABASE.USER"`
	DatabasePwd string `json:"database.pwd" conf:"DATABASE.PWD"`
}

func main(){
	_ = godotenv.Load() //read VAULT.URL | VAULT.PATH | VAULT.TOKEN from the .evn file
	vaultClient := vault.NewVault(Configuration{})
        err := vaultClient.Load()
	if err!=nil{
	    log.Fatal(err)
	}
	router := chi.NewRouter()
	router.Get("/show", show)
	router.Get("/reload", vaultClient.Reload) //change me
	
	http.ListenAndServe(":9991", router)
}

func show(w http.ResponseWriter, r *http.Request){
        //show env variable here.
	w.Write([]byte(os.Getenv("DATABASE.URL")  + " | " + os.Getenv("DATABASE.USER") + " | " + os.Getenv("DATABASE.PWD")))
}

```

For reload configuration, just change value in Vault and hit reload URL.
```shell
  curl http://localhost:9991/reload
```

