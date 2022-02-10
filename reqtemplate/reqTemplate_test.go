package reqtemplate

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_LoadVariable(t *testing.T) {
	p := `{    \"url\":\"${{url}}\",    \"method\":\"${{method}}\",    \"uriParam\":{\"key\":\"${{value1}}\"    } ,    \"header\":{\"key\":\"${{value2}}\"    } ,    \"httpProxy\":\"${{proxy1}}\",    \"socketProxy\":\"${{proxy2}}\",    \"timeout\":\"${{timeout}}\",    \"variable\":{\"key\":\"${{value2}}\"    } ,    \"callbackAddr\":\"${{address}}\",    \"contentType\":\"${{contentType}}\",    \"content\":\"\\\"${{test}}\\\"\",    \"respType\":0}`
	r := ReqParamTemplate{
		ParamsTemp: p,
	}
	r.LoadVariable()
	m := map[string]struct{}{
		"url":         {},
		"method":      {},
		"value1":      {},
		"value2":      {},
		"proxy1":      {},
		"proxy2":      {},
		"timeout":     {},
		"address":     {},
		"contentType": {},
		"test":        {},
	}
	for _, k := range r.ReqParam {
		_, ok := m[k]
		assert.True(t, ok)
		delete(m, k)
	}
	assert.Equal(t, 0, len(m))
}

func Test_LoadSet(t *testing.T) {
	p := `"{    \"url\":\"${{url}}\",    \"method\":\"${{method}}\",    \"uriParam\":{\"key\":\"${{value1}}\"    } ,    \"header\":{\"key\":\"${{value2}}\"    } ,    \"httpProxy\":\"${{proxy1}}\",    \"socketProxy\":\"${{proxy2}}\",    \"timeout\":\"${{timeout}}\",    \"variable\":{\"key\":\"${{value2}}\"    } ,    \"callbackAddr\":\"${{address}}\",    \"contentType\":\"${{contentType}}\",    \"content\":\"\\\"${{test}}\\\"\",    \"respType\":0}"`
	r := ReqParamTemplate{
		ParamsTemp: p,
	}
	r.LoadVariable()
	m := map[string]interface{}{
		"url":         "http://www.baidu.com",
		"method":      "GET",
		"value1":      "value1",
		"value2":      "value2",
		"proxy1":      "proxy1",
		"proxy2":      "proxy2",
		"timeout":     "timeout",
		"address":     "address",
		"contentType": "contentType",
		"test":        1,
	}
	p1, err := r.Set(m)
	assert.NoError(t, err)
	assert.Equal(t, "http://www.baidu.com", p1.URL)
	assert.Equal(t, "GET", p1.Method)
	assert.Equal(t, "value1", p1.URIParam["key"])
	assert.Equal(t, "value2", p1.Header["key"])
	assert.Equal(t, "proxy1", p1.HttpProxy)
	assert.Equal(t, "proxy2", p1.SocketProxy)
	assert.Equal(t, "timeout", p1.Timeout)
	assert.Equal(t, "address", p1.CallbackAddr)
	assert.Equal(t, "contentType", string(p1.ContentType))
	assert.Equal(t, "\"1\"", p1.Content)
}
