package conf

import (
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
	SetConfig(interface{}, map[string]interface{})
}

type values struct{
	Needed []string
	Url string
	Path string
	Token string
}

func NewConf() Values {
	c := new(values)
	c.Needed = []string{
		"VAULT.URL",
		"VAULT.PATH",
		"VAULT.TOKEN",
	}

	for _,v := range c.Needed{
		if os.Getenv(v) == ""{
			log.Fatalf("%v is required", v)
			return nil
		}
	}

	c.Url = os.Getenv("VAULT.URL")
	c.Path = os.Getenv("VAULT.PATH")
	c.Token =  os.Getenv("VAULT.TOKEN")
	return c
}

func (c values) GetConfig() map[string]string{
	storage := make(map[string]string)
	storage[c.Needed[0]] = c.Url
	storage[c.Needed[1]] = c.Path
	storage[c.Needed[2]] = c.Token
	return storage
}

func (c values) SetConfig(confStruct interface{}, vaultResponse map[string]interface{}){
	t := reflect.TypeOf(confStruct)
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		fieldConf := field.Tag.Get(tagConf)
		filedJson := field.Tag.Get(tagJson)

		if fieldConf == ""{
			log.Fatalf("Invalid fieldName...")
		}
		log.Printf("--int-- [%v]:%v \n", fieldConf, vaultResponse[filedJson])
		_ = os.Setenv(fieldConf, vaultResponse[filedJson].(string))
	}
}


