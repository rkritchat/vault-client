Vault - Client for GO language, this dependency support only Vault k/v 2.0

## Installation
```shell
go get -u github.com/rkritchat/vault-client

```

## Usage
Add Vault information to your configuration. (`.env`) file:

```shell
VAULT.URL=http://127.0.0.1:8200/v1/secret
VAULT.PATH=test  //your secret path
VAULT.TOKEN=s.xxxxxxXxxxx
```

Code example

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
type Configuration struct {
	DatabaseUrl string `json:"database.url" conf:"DATABASE.DB_NAME"`
	DatabaseUser string `json:"database.user" conf:"DATABASE.USER"`
	DatabasePwd string `json:"database.pwd" conf:"DATABASE.PWD"`
}

func main(){
	_ = godotenv.Load() //read VAULT.URL | VAULT.PATH | VAULT.TOKEN
	vaultClient := vault.NewVault(conf.Default(), Configuration{})

	router := chi.NewRouter()
	router.Get("/test", test)
	router.Get("/reload", vaultClient.Reload()) //path reload

	http.ListenAndServe(":9991", router)
}

func test(w http.ResponseWriter, r *http.Request){
  //show env variable here.
	w.Write([]byte(os.Getenv("DATABASE.DB_NAME")  + " | " + os.Getenv("DATABASE.USER") + " | " + os.Getenv("DATABASE.PWD")))
}

```

For reload configuration, change value in Vault and hit reload URL.
```shell
  curl http://localhost:9991/reload
```

