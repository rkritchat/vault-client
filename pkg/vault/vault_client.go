package vault

import (
	"github.com/rkritchat/vault-client/pkg/client"
	"github.com/rkritchat/vault-client/pkg/conf"
	"log"
)

func NewVault(i interface{}){
	c := conf.NewConf()
	cli := client.NewClient(c, i)
	vaultResponse, err := cli.LodConfig()
	if err!=nil{
		log.Fatal(err)
	}
	c.SetConfig(i, vaultResponse)
}