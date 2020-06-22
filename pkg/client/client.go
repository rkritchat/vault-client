package client

import (
	"encoding/json"
	"fmt"
	"github.com/rkritchat/vault-client/pkg/conf"
	"net/http"
	"time"
)

type Client interface {
	LodConfig() (map[string]interface{},error)
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

type VaultResponse struct {
	Data Data `json:"data"`
}

type Data struct {
	Data interface{} `json:"data"`
}

func (c *client) LodConfig() (map[string]interface{},error){
	config := c.conf.GetConfig()
	url := fmt.Sprintf("%v/data/%v", config["VAULT.URL"], config["VAULT.PATH"]) //support v2 only
	fmt.Print("url : " + url)
	request, _ := http.NewRequest(http.MethodGet, url, nil)
	request.Header.Set("X-Vault-Token", config["VAULT.TOKEN"])
	cli := &http.Client{
		Timeout: 30 * time.Second,
	}

	result, err := cli.Do(request)
	if err!=nil{
		return nil, err
	}

	var resp VaultResponse
	resp.Data.Data = c.resultStructure

	if err := json.NewDecoder(result.Body).Decode(&resp); err != nil {
		return nil, err
	}
	return resp.Data.Data.(map[string]interface{}), nil
}
