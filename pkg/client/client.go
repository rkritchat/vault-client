package vault_client

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
)

type VaultClient interface {
	LodeConfig(interface{}) (map[string]interface{}, error)
}

type vaultClient struct {
	conf   Config
	client HttpI
}

func NewClient(conf Config, client HttpI) VaultClient {
	return &vaultClient{conf: conf, client: client}
}

type VaultResponse struct {
	Data Data `json:"data"`
}
type Data struct {
	Data interface{} `json:"data"`
}

func (c *vaultClient) LodeConfig(resultStructure interface{}) (map[string]interface{}, error) {
	config, err := c.conf.GetConfig()
	if err != nil {
		return nil, err
	}
	url := fmt.Sprintf("%v/v1/secret/data/%v", config[vaultURL], config[vaultPath]) //support v2 only
	request, _ := http.NewRequest(http.MethodGet, url, nil)
	request.Header.Set(xVaultToken, config[vaultToken])

	result, err := c.client.Do(request)
	if result != nil {
		defer result.Body.Close()
	}

	if err != nil {
		return nil, err
	}

	if result.StatusCode != 200 {
		return nil, errors.New(fmt.Sprintf("vault response status code not ok, status code[%s]", strconv.Itoa(result.StatusCode)))
	}

	var resp VaultResponse
	resp.Data.Data = resultStructure

	if err := json.NewDecoder(result.Body).Decode(&resp); err != nil {
		return nil, err
	}
	return resp.Data.Data.(map[string]interface{}), nil
}
