package client

import (
	"encoding/json"
	"fmt"
	"github.com/rkritchat/vault-client/pkg/conf"
	"log"
	"net/http"
	"time"
)

type VaultClient interface {
	LodeConfig(interface{}) (map[string]interface{}, error)
}

type vaultClient struct {
	conf conf.Values
}

func NewClient(conf conf.Values) VaultClient {
	return &vaultClient{
		conf: conf,
	}
}

type VaultResponse struct {
	Data Data `json:"data"`
}

type Data struct {
	Data interface{} `json:"data"`
}

func (c *vaultClient) LodeConfig(resultStructure interface{}) (map[string]interface{}, error) {
	config := c.conf.GetConfig()
	url := fmt.Sprintf("%v/v1/secret/data/%v", config["VAULT.URL"], config["VAULT.PATH"]) //support v2 only
	request, _ := http.NewRequest(http.MethodGet, url, nil)
	request.Header.Set("X-Vault-Token", config["VAULT.TOKEN"])
	cli := &http.Client{
		Timeout: 30 * time.Second,
	}

	result, err := cli.Do(request)
	if result != nil {
		defer func() {
			if err := result.Body.Close(); err != nil {
				log.Fatal("Exception while close body")
			}
		}()
	}

	if err != nil {
		return nil, err
	}

	if result.StatusCode != 200 {
		log.Fatalf("Response not ok, %v\n", result.StatusCode)
	}

	var resp VaultResponse
	resp.Data.Data = resultStructure

	if err := json.NewDecoder(result.Body).Decode(&resp); err != nil {
		return nil, err
	}
	return resp.Data.Data.(map[string]interface{}), nil
}
