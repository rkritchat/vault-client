package client

import (
	"encoding/json"
	"fmt"
	"github.com/rkritchat/vault-client/pkg/conf"
	"net/http"
	"time"
)

type Client interface {
	LodConfig() (interface{},error)
}

type client struct{
	conf conf.Values
	resultStructure interface{}
}

func NewClient(conf conf.Values, resultStructure interface{}) Client {
	return &client{
		conf: conf,
		resultStructure: resultStructure,
	}
}

func (c *client) LodConfig() (interface{},error){
	config := c.conf.GetConfig()
	url := fmt.Sprintf("%v/data/%v", config["vault.url"], config["vault.path"]) //v2 only
	fmt.Print("url : " + url)
	request, _ := http.NewRequest(http.MethodGet, url, nil)
	request.Header.Set("X-Vault-Token", config["vault.token"])
	cli := &http.Client{
		Timeout: 30 * time.Second,
	}

	result, err := cli.Do(request)
	if err!=nil{
		return nil, err
	}

	if err := json.NewDecoder(result.Body).Decode(&c.resultStructure); err != nil {
		return nil, err
	}
	return c.resultStructure, nil
}
