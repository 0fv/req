package reqtemplate

import (
	"encoding/json"
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/0fv/req/httpreq"
)

type ReqParamTemplate struct {
	ParamsTemp   string
	ReqParam     []string
	DefaultParam map[string]string
}

var variablepattern = regexp.MustCompile(`\$\{{([^\}]+)\}}`)

func (r *ReqParamTemplate) LoadVariable() {
	params := variablepattern.FindAllString(r.ParamsTemp, -1)
	set := make(map[string]struct{})
	for _, v := range params {
		if _, ok := set[v]; ok {
			continue
		}
		set[v] = struct{}{}
		param := strings.TrimLeft(v, "${{")
		param = strings.TrimRight(param, "}}")
		r.ReqParam = append(r.ReqParam, param)
	}
}

func (r ReqParamTemplate) Set(m map[string]interface{}) (httpreq.Param, error) {
	var ret httpreq.Param
	data := r.ParamsTemp
	//check variable
	for _, v := range r.ReqParam {
		if _, ok := m[v]; !ok {
			if _, ok := r.DefaultParam[v]; ok {
				m[v] = r.DefaultParam[v]
			} else {
				return ret, errors.New("variable not found")
			}
		}
	}
	//replace variable
	for k, v := range m {
		data = strings.Replace(data, "${{"+k+"}}", fmt.Sprint(v), -1)
	}
	s, err := strconv.Unquote(data)
	if err != nil {
		return ret, err
	}
	err = json.Unmarshal([]byte(s), &ret)
	return ret, err
}
