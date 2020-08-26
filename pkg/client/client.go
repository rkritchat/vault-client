package client

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/rkritchat/vault-client/pkg/conf"
	"github.com/rkritchat/vault-client/pkg/http_i"
	"net/http"
	"strconv"
)

const (
	vaultURL    = "VAULT.URL"
	vaultPath   = "VAULT.PATH"
	vaultToken  = "VAULT.TOKEN"
	xVaultToken = "X-Vault-Token"
)

type Vault interface {
	LodeConfig() (map[string]interface{}, error)
}

type vault struct {
	conf   conf.Config
	client http_i.HttpI
}

func NewVault(conf conf.Config, c http_i.HttpI) Vault {
	return &vault{conf: conf, client: c}
}

type VaultResponse struct {
	Data Data `json:"data"`
}
type Data struct {
	Data interface{} `json:"data"`
}

func (v *vault) LodeConfig() (map[string]interface{}, error) {
	config, err := v.conf.GetConfig()
	if err != nil {
		return nil, err
	}
	url := fmt.Sprintf("%v/v1/secret/data/%v", config[vaultURL], config[vaultPath]) //support v2 only
	request, _ := http.NewRequest(http.MethodGet, url, nil)
	request.Header.Set(xVaultToken, config[vaultToken])

	result, err := v.client.Do(request)
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
	if err := json.NewDecoder(result.Body).Decode(&resp); err != nil {
		return nil, err
	}
	return resp.Data.Data.(map[string]interface{}), nil
}
