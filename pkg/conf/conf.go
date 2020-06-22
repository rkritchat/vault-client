package conf

import (
	"fmt"
	"log"
	"os"
	"reflect"
	"strconv"
)

const tagName = "conf"

type Values interface {
	GetConfig() map[string]string
	SetConfig(interface{})
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
		"vault.url",
		"vault.path",
		"vault.token",
	}

	for _,v := range c.Needed{
		if os.Getenv(v) == ""{
			log.Fatalf("%v is required", v)
			return nil
		}
	}

	c.Url = os.Getenv("vault.url")
	c.Path = os.Getenv("vault.path")
	c.Token =  os.Getenv("vault.token")
	return c
}

func (c values) GetConfig() map[string]string{
	storage := make(map[string]string)
	storage[c.Needed[0]] = c.Url
	storage[c.Needed[1]] = c.Path
	storage[c.Needed[2]] = c.Url
	return storage
}

func (c values) SetConfig(confStruct interface{}){
	t := reflect.TypeOf(confStruct)
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		fieldName := field.Tag.Get(tagName)
		if fieldName == ""{
			log.Fatalf("Invalid fieldName...")
		}

		switch reflect.ValueOf(confStruct).Field(i).Kind() {
		case reflect.Float32, reflect.Float64:
			value := reflect.ValueOf(confStruct).Field(i).Float()
			log.Printf("--float-- [%v]:%v \n", fieldName, value)

			valueStr := fmt.Sprintf("%.2f", value)
			_ = os.Setenv(fieldName, valueStr)
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			value := reflect.ValueOf(confStruct).Field(i).Int()
			log.Printf("--int-- [%v]:%v \n", fieldName, value)

			valueStr := strconv.FormatInt(value, 10)
			_ = os.Setenv(fieldName, valueStr)
		case reflect.String:
			value := reflect.ValueOf(confStruct).Field(i).String()
			log.Printf("--string-- [%v]:%v \n", fieldName, value)

			_ = os.Setenv(fieldName, value)
		default:
			log.Fatal("type not match")
		}
	}
}


