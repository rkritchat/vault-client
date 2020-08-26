package vault_client

import (
	"crypto/tls"
	"net/http"
	"time"
)

type HttpI interface {
	Do(req *http.Request) (*http.Response, error)
}

type httpI struct{
	client http.Client
}

func NewHttpInterface() HttpI {
	return &httpI{
		client: initClient(),
	}
}

func (i *httpI)Do(req *http.Request) (*http.Response, error){
	return i.client.Do(req)
}

func initClient() http.Client {
	customTransport := http.DefaultTransport.(*http.Transport).Clone()
	customTransport.TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	client := http.Client{
		Timeout:   time.Duration(30) * time.Second,
		Transport: customTransport,
	}
	return client
}
