package httpreq

import (
	"bytes"
	"errors"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/url"
	"strings"
	"time"

	"golang.org/x/net/proxy"
)

type Param struct {
	URL          string            `json:"url"`
	Method       string            `json:"method"`
	URIParam     map[string]string `json:"uriParam"`
	Header       map[string]string `json:"header"`
	HttpProxy    string            `json:"httpProxy"`
	SocketProxy  string            `json:"socketProxy"`
	Timeout      string            `json:"timeout"`
	Variable     map[string]string `json:"variable"`
	CallbackAddr string            `json:"callbackAddr"`
	BodyContent
	RespType RespType `json:"respType"`
}

type ContentType string

const (
	ContentTypeWWWForm   ContentType = "application/x-www-form-urlencoded"
	ContentTypeForm      ContentType = "multipart/form-data"
	ContentTypeJSON      ContentType = "application/json"
	ContentTypeNoContent ContentType = ""
)

type BodyContent struct {
	ContentType ContentType `json:"contentType"`
	Content     interface{} `json:"content"`
}

var ErrContentNotCorrect = errors.New("content not correct")

func (b BodyContent) BuildBody(req *http.Request) (err error) {
	switch b.ContentType {
	case ContentTypeWWWForm:
		req.Body, err = b.buildWWWForm()
		req.Header.Set("Content-Type", string(b.ContentType))
	case ContentTypeForm:
		var contentType string
		req.Body, contentType, err = b.buildFormBody()
		req.Header.Set("Content-Type", contentType)
	case ContentTypeJSON:
		req.Body, err = b.BuildString()
	default:
		if b.Content != nil {
			req.Body, err = b.BuildString()
		}
	}
	return
}

func (b BodyContent) buildWWWForm() (io.ReadCloser, error) {
	data, ok := b.Content.(map[string]string)
	if !ok {
		return nil, ErrContentNotCorrect
	}
	val := url.Values{}
	for k, v := range data {
		val.Set(k, v)

	}
	return ioutil.NopCloser(strings.NewReader(val.Encode())), nil
}

func (b BodyContent) buildFormBody() (body io.ReadCloser, contentType string, err error) {
	data, ok := b.Content.(map[string]FormDataValue)
	if !ok {
		return
	}
	var bf bytes.Buffer
	w := multipart.NewWriter(&bf)
	for k, v := range data {
		if v.Type == FormDataTypeFile {
			var f io.Writer
			f, err = w.CreateFormFile(k, v.FileName)
			if err != nil {
				return
			}
			var resp *http.Response
			resp, err = http.Get(v.Value)
			if err != nil {
				return
			}
			io.Copy(f, resp.Body)
		} else {
			w.WriteField(k, v.Value)
		}
	}
	w.Close()
	return ioutil.NopCloser(&bf), w.FormDataContentType(), nil
}

func (b BodyContent) BuildString() (body io.ReadCloser, err error) {
	data, ok := b.Content.(string)
	if !ok {
		buf, ok := b.Content.([]byte)
		if !ok {
			return nil, ErrContentNotCorrect
		}
		return ioutil.NopCloser(bytes.NewReader(buf)), nil
	}
	return ioutil.NopCloser(strings.NewReader(data)), nil

}

type RespType uint8

const (
	RespAsync RespType = iota + 1
	RespSync
	Callback
)

func (p Param) BuildReq() (req *http.Request, err error) {
	u, err := url.Parse(p.URL)
	if err != nil {
		return
	}
	query := u.Query()
	for k, v := range p.URIParam {
		query.Set(k, v)
	}
	url := p.URL + "?" + query.Encode()
	if p.Method == "" {
		p.Method = "GET"
	}
	req, err = http.NewRequest(p.Method, url, nil)
	if err != nil {
		return
	}
	for k, v := range p.Header {
		req.Header.Add(k, v)
	}
	err = p.BuildBody(req)
	return
}

func (p Param) BuildClient() (clt *http.Client, err error) {
	var proxyURL *url.URL
	if p.HttpProxy != "" {
		proxyURL, err = url.Parse(p.HttpProxy)
		if err != nil {
			return
		}
		clt = &http.Client{
			Transport: &http.Transport{
				Proxy: http.ProxyURL(proxyURL),
			},
		}
	} else if p.SocketProxy != "" {
		var dialSocksProxy proxy.Dialer
		dialSocksProxy, err = proxy.SOCKS5("tcp", p.SocketProxy, nil, proxy.Direct)
		if err != nil {
			return
		}
		clt = &http.Client{
			Transport: &http.Transport{Dial: dialSocksProxy.Dial},
		}
	} else {
		clt = &http.Client{}
	}
	if p.Timeout != "" {
		timeout, err := time.ParseDuration(p.Timeout)
		if err == nil {
			clt.Timeout = timeout
		}
	}
	return
}

func (p Param) HttpReq() (resp *http.Response, err error) {
	var req *http.Request
	req, err = p.BuildReq()
	if err != nil {
		return
	}
	client, err := p.BuildClient()
	if err != nil {
		return
	}
	resp, err = client.Do(req)
	return
}
