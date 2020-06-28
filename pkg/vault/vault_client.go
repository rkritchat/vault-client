package vault

import (
	"github.com/rkritchat/vault-client/pkg/client"
	"github.com/rkritchat/vault-client/pkg/conf"
	"github.com/rkritchat/vault-client/pkg/constant"
	"log"
	"net/http"
	"os"
)

type Vault interface {
	Reload() func(w http.ResponseWriter, r *http.Request)
}

type vault struct {
	i           interface{}
	value       conf.Values
	vaultClient client.VaultClient
}

func NewVault(value conf.Values, i interface{}) Vault {
	vaultClient := client.NewClient(value)
	if os.Getenv(constant.VaultDisable) == constant.Empty || os.Getenv(constant.VaultDisable) == constant.False {
		if response, err := vaultClient.LodeConfig(i); err != nil {
			log.Fatal(err)
		} else {
			_ = value.SetConfig(i, response)
		}
	}
	return &vault{
		i:           i,
		value:       value,
		vaultClient: vaultClient,
	}
}

func (v vault) Reload() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if os.Getenv(constant.VaultDisable) == constant.True {
			_, _ = w.Write([]byte("vault-client is not enable."))
			return
		}

		config, err := v.vaultClient.LodeConfig(v.i)
		if err != nil {
			_, _ = w.Write([]byte("Reload config Exception"))
		} else {
			change := v.value.SetConfig(v.i, config)
			if change != "" {
				_, _ = w.Write([]byte("[" + change[:len(change)-1] + "]"))
				return
			}
			_, _ = w.Write([]byte("[]"))
		}
	}
}
