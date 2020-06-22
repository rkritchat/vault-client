package vault

import (
	"github.com/rkritchat/vault-client/pkg/client"
	"github.com/rkritchat/vault-client/pkg/conf"
)

func NewVault(){
	c := conf.NewConf()
	cli := client.NewClient(c, nil)
	config, err := cli.LodConfig()
	if err!=nil{

	}
	c.SetConfig(config)
}