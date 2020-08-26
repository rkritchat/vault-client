package vault

import (
	"errors"
	"github.com/rkritchat/vault-client/pkg/client"
	"github.com/rkritchat/vault-client/pkg/conf"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"net/http/httptest"
	"os"
	"testing"
)

type test struct {
	Test1 string `json:"test1" conf:"test1"`
}

type clientMock struct {
	mock.Mock
}

func (m *clientMock) LodeConfig() (map[string]interface{}, error) {
	args := m.Called()
	if args.Get(0) != nil {
		return args.Get(0).(map[string]interface{}), nil
	} else {
		return nil, args.Error(1)
	}
}

type confMock struct {
	mock.Mock
}

func (m *confMock) GetConfig() (map[string]string, error) {
	args := m.Called()
	if args.Get(0) != nil {
		return args.Get(0).(map[string]string), nil
	} else {
		return nil, args.Error(1)
	}
}
func (m *confMock) SetConfig(i interface{}, mockMap map[string]interface{}) (string, error) {
	args := m.Called(i, mockMap)
	if args.Get(0) != nil {
		return args.Get(0).(string), nil
	} else {
		return "", args.Error(1)
	}
}

func Test_Load(t *testing.T) {
	_ = os.Setenv("VAULT.URL", "http://test.me")
	_ = os.Setenv("VAULT.PATH", "path")
	_ = os.Setenv("VAULT.TOKEN", "token")
	tt := []struct {
		name   string
		i      interface{}
		conf   conf.Config
		client client.Vault
		expect error
	}{
		{
			name:   "Case success",
			i:      test{},
			conf:   mockConfMock("ok"),
			client: mockClientMock("ok"),
			expect: nil,
		},
		{
			name:   "Case error while load config",
			i:      test{},
			conf:   mockConfMock("ok"),
			client: mockClientMock("!ok"),
			expect: errors.New("test_err"),
		},
		{
			name:   "Case error while set config",
			i:      test{},
			conf:   mockConfMock("set_conf_fail"),
			client: mockClientMock("ok"),
			expect: errors.New("test_err"),
		},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			v := vault{i: tc.i, client: tc.client, conf: tc.conf}
			load, err := v.Load()
			if tc.expect != nil {
				assert.NotNil(t, err)
			} else {
				assert.Nil(t, err)
				assert.NotNil(t, load)
			}
		})
	}
}

func Test_Reload(t *testing.T) {
	tt := []struct {
		name   string
		i      interface{}
		conf   conf.Config
		client client.Vault
		expect string
	}{
		{
			name:   "Case reload success",
			i:      test{},
			conf:   mockConfMock("ok"),
			client: mockClientMock("ok"),
			expect: "[test]",
		},
		{
			name:   "Case error while load config",
			i:      test{},
			conf:   mockConfMock("ok"),
			client: mockClientMock("!ok"),
			expect: "Reload config Exception",
		},
		{
			name:   "Case error while set config",
			i:      test{},
			conf:   mockConfMock("set_conf_fail"),
			client: mockClientMock("ok"),
			expect: "Reload config Exception",
		},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			v := vault{i: tc.i, client: tc.client, conf: tc.conf}
			v.Reload(w, nil)
			assert.Equal(t, tc.expect,w.Body.String())
		})
	}
}

func mockClientMock(tc string) *clientMock {
	c := new(clientMock)
	switch tc {
	case "ok":
		c.On("LodeConfig").Return(mockMapResult(), nil)
	case "!ok":
		c.On("LodeConfig").Return(nil, errors.New("test_err"))
	}
	return c
}

func mockConfMock(tc string) *confMock {
	c := new(confMock)
	switch tc {
	case "ok":
		c.On("GetConfig").Return(mockVaultResp("ok"), nil)
		c.On("SetConfig", test{}, mockMapResult()).Return("test2", nil)
	case "set_conf_fail":
		c.On("SetConfig", test{}, mockMapResult()).Return(nil, errors.New("test_err"))
	}
	return c
}

func mockMapResult() map[string]interface{} {
	m := make(map[string]interface{})
	m["test1"] = "for testing"
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
