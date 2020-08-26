package vault_client

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func Test_NewConfig(t *testing.T) {
	tt := []struct {
		name   string
		expect error
	}{
		{
			name:   "Case error",
			expect: errors.New("test_err"),
		},
		{
			name:   "Case success",
			expect: nil,
		},
	}
	for _, tc := range tt {
		if tc.name == "Case success" {
			_ = os.Setenv("VAULT.URL", "url")
			_ = os.Setenv("VAULT.PATH", "path")
			_ = os.Setenv("VAULT.TOKEN", "token")
		}
		t.Run(tc.name, func(t *testing.T) {
			_, err := NewConfig()
			if tc.expect != nil {
				assert.NotNil(t, err)
			} else {
				assert.Nil(t, err)
			}
		})
	}
}

func Test_GetConfig(t *testing.T) {
	type expect struct {
		r   map[string]string
		err error
	}
	tt := []struct {
		name   string
		keys   []string
		url    string
		path   string
		token  string
		expect expect
	}{
		{
			name:   "Case success",
			keys:   []string{"n1", "n2", "n3"},
			url:    "url",
			path:   "path",
			token:  "token",
			expect: expect{r: mockMap(), err: nil},
		},
		{
			name:   "Case success",
			keys:   []string{"n1"},
			url:    "url",
			path:   "path",
			token:  "token",
			expect: expect{r: nil, err: errors.New("test_err")},
		},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			c := config{Keys: tc.keys, Url: "url", Path: "path", Token: "token"}
			r, err := c.GetConfig()
			if tc.expect.err != nil {
				assert.NotNil(t, err)
			} else {
				assert.Equal(t, tc.expect.r, r)
				assert.Nil(t, err)
			}
		})
	}
}

func Test_SetConfig(t *testing.T) {
	type expect struct {
		change string
		err    error
	}
	//case valid struct
	type mockStruct struct {
		testConf1 string `json:"testConf1" conf:"testConf1"`
		testConf2 string `json:"testConf2" conf:"testConf2"`
	}
	//case invalid struct
	type mockInvalidStruct struct{
		testInvalid string
	}
	tt := []struct {
		name          string
		confStruct    interface{}
		vaultResponse map[string]interface{}
		expect        expect
	}{
		{
			name:          "Case success",
			confStruct:    mockStruct{},
			vaultResponse: mockVaultResp("ok"),
			expect:        expect{change: "testConf1 testConf2 ", err: nil},
		},
		{
			name:          "Case not found config from vault response",
			confStruct:    mockStruct{},
			vaultResponse: mockVaultResp("!ok"),
			expect:        expect{change: "", err: errors.New("not found config [testConf1] in Vault")},
		},
		{
			name:          "Case struct not contain conf tag",
			confStruct:    mockInvalidStruct{},
			vaultResponse: mockVaultResp("ok"),
			expect:        expect{change: "", err: errors.New("invalid struct, tag conf is required")},
		},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			c := config{}
			r, err := c.SetConfig(tc.confStruct, tc.vaultResponse)
			if tc.expect.err != nil {
				assert.NotNil(t, err)
				assert.Equal(t, tc.expect.err.Error(), err.Error())
			} else {
				assert.Nil(t, err)
				assert.Equal(t, tc.expect.change, r)
			}
		})
	}
}

func mockMap() map[string]string {
	m := make(map[string]string)
	m["n1"] = "url"
	m["n2"] = "path"
	m["n3"] = "token"
	return m
}

func mockVaultResp(tc string) map[string]interface{} {
	m := make(map[string]interface{}, 2)
	switch tc {
	case "ok":
		m["testConf1"] = "testConf1"
		m["testConf2"] = "testConf2"
	case "!ok":
		m["invalid"] = "invalid"
	}
	return m
}
