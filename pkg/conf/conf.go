package vault_client

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
	Keys  []string
	Url   string
	Path  string
	Token string
}

func NewConfig() (Config, error) {
	c := new(config)
	c.Keys = []string{vaultURL, vaultPath, vaultToken}
	err := validateConfig(c)
	if err != nil {
		return nil, err
	}

	c.Url = os.Getenv(vaultURL)
	c.Path = os.Getenv(vaultPath)
	c.Token = os.Getenv(vaultToken)
	return c, nil
}

func validateConfig(c *config) error {
	for _, v := range c.Keys {
		if os.Getenv(v) == empty {
			return errors.New(fmt.Sprintf("[%s] is required", v))
		}
	}
	return nil
}

func (c config) GetConfig() (map[string]string, error) {
	if len(c.Keys) == 3 {
		storage := make(map[string]string)
		storage[c.Keys[0]] = c.Url
		storage[c.Keys[1]] = c.Path
		storage[c.Keys[2]] = c.Token
		return storage, nil
	}
	return nil, errors.New("config storage's length must be equals three")
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
