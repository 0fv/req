package vm

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/0fv/req/httpreq"
	"github.com/0fv/req/log"
	"github.com/dop251/goja"
	"go.uber.org/zap"
)

type JSVM struct {
	vm     *goja.Runtime
	logger *zap.Logger
}

func NewVm(logger *zap.Logger) *JSVM {
	vm := goja.New()
	if logger == nil {
		logger = log.NewSimplelog()
	}
	return &JSVM{vm: vm, logger: logger}
}

func (j *JSVM) LoadReqParam(param *httpreq.Param) {
	buf, err := json.Marshal(param)
	if err != nil {
		j.logger.Error("marshal data failed", zap.Error(err))
		return
	}
	j.vm.Set("param", string(buf))
	j.vm.RunString("param = JSON.parse(param)")
}

func (j *JSVM) LoadRespVariable(resp *http.Response) {
	//header
	m := headerToMap(resp.Header)
	j.vm.Set("header", m)
	j.vm.Set("statusCode", resp.StatusCode)
	j.vm.Set("uri", resp.Request.RequestURI)
	//body
	body, err := ioutil.ReadAll(resp.Body)
	if err == nil {
		j.logger.Error("read data failed", zap.Error(err))
	}
	j.vm.Set("body", string(body))
}

func (j *JSVM) LoadParam() *httpreq.Param {
	param := &httpreq.Param{}
	d := j.vm.Get("param").Export()
	data, err := json.Marshal(d)
	if err != nil {
		j.logger.Error("marshal data failed", zap.Error(err))
		return param
	}
	err = json.Unmarshal(data, param)
	if err != nil {
		j.logger.Error("unmarshal data failed", zap.Error(err))
	}
	return param
}

func (j *JSVM) Exec(script string) ([]byte, error) {
	result, err := j.vm.RunString(script)
	if err != nil {
		j.logger.Error("run script failed", zap.Error(err))
		return []byte{}, err
	}
	m := result.Export()
	if m == nil {
		return []byte{}, nil
	}
	b, err := json.Marshal(m)
	if err != nil {
		j.logger.Error("marshal data failed", zap.Error(err))
		return []byte{}, err
	}
	return b, nil
}

func headerToMap(header http.Header) map[string]string {
	m := make(map[string]string)
	for k, v := range header {
		m[k] = strings.Join(v, ";")
	}
	return m
}
