package conf

import (
	"errors"
	"fmt"
	"os"
	"reflect"
)

const (
	tagConf       = "conf"
	tagJson       = "json"
	empty         = ""
	space         = " "
	vaultURL      = "VAULT.URL"
	vaultPath     = "VAULT.PATH"
	vaultToken    = "VAULT.TOKEN"
	invalidStruct = "invalid struct, tag conf is required"
)

type Config interface {
	GetConfig() (map[string]string, error)
	SetConfig(interface{}, map[string]interface{}) (string, error)
}

type config struct {
	Url   string
	Path  string
	Token string
}

func NewConfig() Config {
	c := new(config)
	c.Url = os.Getenv(vaultURL)
	c.Path = os.Getenv(vaultPath)
	c.Token = os.Getenv(vaultToken)
	return c
}

func (c config) GetConfig() (map[string]string, error) {
	err := c.validateConfig()
	if err != nil {
		return nil, err
	}
	storage := make(map[string]string)
	storage[vaultURL] = c.Url
	storage[vaultPath] = c.Path
	storage[vaultToken] = c.Token
	return storage, nil
}

func (c config) SetConfig(confStruct interface{}, vaultResponse map[string]interface{}) (string, error) {
	t := reflect.TypeOf(confStruct)
	change := empty
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		fieldConf := field.Tag.Get(tagConf)
		filedJson := field.Tag.Get(tagJson)

		if fieldConf == empty {
			return empty, errors.New(invalidStruct)
		}

		if vaultResponse[filedJson] == nil {
			return empty, errors.New(fmt.Sprintf("not found config [%s] in Vault", filedJson))
		}

		if os.Getenv(fieldConf) != vaultResponse[filedJson].(string) {
			_ = os.Setenv(fieldConf, vaultResponse[filedJson].(string))
			change += filedJson + space
		}
	}
	return change, nil
}

func (c config) validateConfig() error {
	if len(c.Url) == 0 {
		return errors.New(fmt.Sprintf("[%s] is required", vaultURL))
	}
	if len(c.Path) == 0 {
		return errors.New(fmt.Sprintf("[%s] is required", vaultPath))
	}
	if len(c.Token) == 0 {
		return errors.New(fmt.Sprintf("[%s] is required", vaultToken))
	}
	return nil
}
