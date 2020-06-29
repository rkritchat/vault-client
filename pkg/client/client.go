package client

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/rkritchat/vault-client/pkg/conf"
	"github.com/rkritchat/vault-client/pkg/constant"
	"log"
	"net/http"
	"strconv"
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

var (
	vaultStatusCodNotOk = func(param int) error {
		return errors.New("vault response status code not ok, status code[" + strconv.Itoa(param) + "]")
	}
)

func (c *vaultClient) LodeConfig(resultStructure interface{}) (map[string]interface{}, error) {
	config := c.conf.GetConfig()
	url := fmt.Sprintf("%v/v1/secret/data/%v", config[constant.VaultURL], config[constant.VaultPath]) //support v2 only
	request, _ := http.NewRequest(http.MethodGet, url, nil)
	request.Header.Set(constant.XVaultToken, config[constant.VaultToken])
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
		return nil, vaultStatusCodNotOk(result.StatusCode)
	}

	var resp VaultResponse
	resp.Data.Data = resultStructure

	if err := json.NewDecoder(result.Body).Decode(&resp); err != nil {
		return nil, err
	}
	return resp.Data.Data.(map[string]interface{}), nil
}
