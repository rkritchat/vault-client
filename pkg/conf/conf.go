package conf

import (
	"errors"
	"github.com/rkritchat/vault-client/pkg/constant"
	"os"
	"reflect"
)

const (
	tagConf = "conf"
	tagJson = "json"
)

var (
	//error
	tagConfIsRequired = errors.New("tag conf is required")
	notFoundConfig    = func(param string) error {
		return errors.New("not found config [" + param + "] in Vault")
	}
	missingMandatoryConf = func(param string) error {
		return errors.New("[" + param + "] is required")
	}
)

type Values interface {
	GetConfig() map[string]string
	SetConfig(interface{}, map[string]interface{}) (string, error)
}

type values struct {
	Needed []string
	Url    string
	Path   string
	Token  string
}

func Default() (Values, error) {
	c := new(values)
	c.Needed = []string{
		constant.VaultURL,
		constant.VaultPath,
		constant.VaultToken,
	}

	for _, v := range c.Needed {
		if os.Getenv(v) == "" {
			return nil, missingMandatoryConf(v)
		}
	}

	c.Url = os.Getenv(constant.VaultURL)
	c.Path = os.Getenv(constant.VaultPath)
	c.Token = os.Getenv(constant.VaultToken)
	return c, nil
}

func (c values) GetConfig() map[string]string {
	storage := make(map[string]string)
	storage[c.Needed[0]] = c.Url
	storage[c.Needed[1]] = c.Path
	storage[c.Needed[2]] = c.Token
	return storage
}

func (c values) SetConfig(confStruct interface{}, vaultResponse map[string]interface{}) (string, error) {
	t := reflect.TypeOf(confStruct)
	change := ""
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		fieldConf := field.Tag.Get(tagConf)
		filedJson := field.Tag.Get(tagJson)

		if fieldConf == constant.Empty {
			return "", tagConfIsRequired
		}

		if vaultResponse[filedJson] == nil {
			return "", notFoundConfig(filedJson)
		}

		if os.Getenv(fieldConf) != vaultResponse[filedJson].(string) {
			_ = os.Setenv(fieldConf, vaultResponse[filedJson].(string))
			change += filedJson + constant.Space
		}
	}
	return change, nil
}
