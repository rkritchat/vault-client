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

func NewVault(initValues func() (conf.Values, error), i interface{}) (Vault, error) {
	values, err := initValues()
	if err != nil {
		return nil, err
	}
	vaultClient := client.NewClient(values)
	if os.Getenv(constant.VaultDisable) == constant.Empty || os.Getenv(constant.VaultDisable) == constant.False {
		if response, err := vaultClient.LodeConfig(i); err != nil {
			log.Fatal(err)
		} else {
			_, err = values.SetConfig(i, response)
			if err != nil {
				return nil, err
			}
		}
	}
	return &vault{
		i:           i,
		value:       values,
		vaultClient: vaultClient,
	}, nil
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
			change, err := v.value.SetConfig(v.i, config)

			if err != nil {
				log.Printf("Err %v", err)
				_, _ = w.Write([]byte("Exception occurred while reload config"))
				return
			}

			if change != "" {
				_, _ = w.Write([]byte("[" + change[:len(change)-1] + "]"))
				return
			}
			_, _ = w.Write([]byte("[]"))
		}
	}
}

func (v vault) DefaultConfig() (conf.Values, error) {
	defaultValue, err := conf.Default()
	if err != nil {
		return nil, err
	}
	return defaultValue, nil
}
