package runner

import (
	"github.com/0fv/req/reqtemplate"
)

type DataSource interface {
	Load(key string) (reqtemplate.ReqParamTemplate, error)
}

type OutPut interface {
	Set(key string, value interface{})
	Get(key string) interface{}
	Delete(key string)
}
