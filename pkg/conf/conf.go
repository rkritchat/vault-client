package conf

import (
	"github.com/rkritchat/vault-client/pkg/constant"
	"log"
	"os"
	"reflect"
)

const (
	tagConf = "conf"
	tagJson = "json"
)

type Values interface {
	GetConfig() map[string]string
	SetConfig(interface{}, map[string]interface{}) string
}

type values struct {
	Needed []string
	Url    string
	Path   string
	Token  string
}

func Default() Values {
	c := new(values)
	c.Needed = []string{
		constant.VaultURL,
		constant.VaultPath,
		constant.VaultToken,
	}

	for _, v := range c.Needed {
		if os.Getenv(v) == "" {
			log.Fatalf("%v is required", v)
			return nil
		}
	}

	c.Url = os.Getenv(constant.VaultURL)
	c.Path = os.Getenv(constant.VaultPath)
	c.Token = os.Getenv(constant.VaultToken)
	return c
}

func (c values) GetConfig() map[string]string {
	storage := make(map[string]string)
	storage[c.Needed[0]] = c.Url
	storage[c.Needed[1]] = c.Path
	storage[c.Needed[2]] = c.Token
	return storage
}

func (c values) SetConfig(confStruct interface{}, vaultResponse map[string]interface{}) string {
	t := reflect.TypeOf(confStruct)
	change := ""
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		fieldConf := field.Tag.Get(tagConf)
		filedJson := field.Tag.Get(tagJson)

		if fieldConf == constant.Empty {
			log.Fatalf("Tag conf is required.")
		}

		if vaultResponse[filedJson] == nil {
			log.Fatalf("Not found config [%v] in Vault", filedJson)
		}

		if os.Getenv(fieldConf) != vaultResponse[filedJson].(string) {
			_ = os.Setenv(fieldConf, vaultResponse[filedJson].(string))
			change += filedJson + constant.Space
		}
	}
	return change
}
