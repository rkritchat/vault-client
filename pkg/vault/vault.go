package vault

import (
	"github.com/rkritchat/vault-client/pkg/client"
	"github.com/rkritchat/vault-client/pkg/conf"
	"github.com/rkritchat/vault-client/pkg/http_i"
	"log"
	"net/http"
	"os"
	"strconv"
)

const (
	reloadConfException  = "Reload config Exception"
	vaultClientIsDisable = "vault-client is disable"
	vaultDisable         = "VAULT.DISABLE"
	empty                = ""
)

type Vault interface {
	Load()(*vault,error)
	Reload(w http.ResponseWriter, r *http.Request)
}

type vault struct {
	i      interface{}
	conf   conf.Config
	client client.Vault
}

func NewVault(i interface{}) Vault {
	config := conf.NewConfig()
	return &vault{i: i, conf: config, client: client.NewVault(config, http_i.NewHttpInterface())}
}

func (v *vault) Load() (*vault, error) {
	isDisable, _ := strconv.ParseBool(os.Getenv(vaultDisable))
	if isDisable {
		return nil, nil
	}
	response, err := v.client.LodeConfig()
	if err != nil {
		return nil, err
	}
	_, err = v.conf.SetConfig(v.i, response)
	if err != nil {
		return nil, err
	}
	return v, nil
}

func (v *vault) Reload(w http.ResponseWriter, r *http.Request)  {
	isDisable, _ := strconv.ParseBool(os.Getenv(vaultDisable))
	if isDisable {
		_, _ = w.Write([]byte(vaultClientIsDisable))
		return
	}

	config, err := v.client.LodeConfig()
	if err != nil {
		_, _ = w.Write([]byte(reloadConfException))
	} else {
		change, err := v.conf.SetConfig(v.i, config)

		if err != nil {
			log.Printf("Err %v", err)
			_, _ = w.Write([]byte(reloadConfException))
			return
		}

		if change != empty {
			_, _ = w.Write([]byte("[" + change[:len(change)-1] + "]"))
			return
		}
		_, _ = w.Write([]byte("[]"))
	}
}
