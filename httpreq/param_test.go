package httpreq

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func StartHttpServer() {
	http.HandleFunc("/normal", func(resp http.ResponseWriter, req *http.Request) {
		resp.Write([]byte("normal"))
	})
	http.HandleFunc("/method/", func(resp http.ResponseWriter, req *http.Request) {
		resp.Write([]byte(req.Method))
	})
	http.HandleFunc("/uri", func(resp http.ResponseWriter, req *http.Request) {
		m := map[string]string{}
		for k, v := range req.URL.Query() {
			m[k] = strings.Join(v, ";")
		}
		data, _ := json.Marshal(m)
		resp.Write(data)
	})
	http.HandleFunc("/header/", func(resp http.ResponseWriter, req *http.Request) {
		resp.Write([]byte(req.Header.Get("Header")))
	})
	http.HandleFunc("/body/", func(resp http.ResponseWriter, req *http.Request) {
		resp.Header().Set("ct", req.Header.Get("Content-Type"))
		io.Copy(resp, req.Body)
		// req.Body.Close()
	})
	http.HandleFunc("/file", func(resp http.ResponseWriter, req *http.Request) {
		resp.Write([]byte("file"))
	})
	go func() {
		http.ListenAndServe(":8080", nil)
	}()
}

var baseUrl = "http://localhost:8080"

func TestMain(t *testing.M) {
	StartHttpServer()
	t.Run()
}

func Test_Normal(t *testing.T) {
	p := Param{
		URL: baseUrl + "/normal",
	}
	resp, err := p.HttpReq()
	assert.NoError(t, err)
	buf, err := ioutil.ReadAll(resp.Body)
	assert.NoError(t, err)
	assert.Equal(t, "normal", string(buf))
}

func Test_Method(t *testing.T) {
	p := Param{
		URL:    baseUrl + "/method/",
		Method: http.MethodPut,
	}
	resp, err := p.HttpReq()
	assert.NoError(t, err)
	buf, err := ioutil.ReadAll(resp.Body)
	assert.NoError(t, err)
	assert.Equal(t, http.MethodPut, string(buf))
}

func Test_URIParam(t *testing.T) {
	p := Param{
		URL: baseUrl + "/uri",
		URIParam: map[string]string{
			"a": "1",
			"b": "2",
		},
	}
	resp, err := p.HttpReq()
	assert.NoError(t, err)
	buf, err := ioutil.ReadAll(resp.Body)
	assert.NoError(t, err)
	m := map[string]string{}
	err = json.Unmarshal(buf, &m)
	assert.NoError(t, err)
	assert.Equal(t, "1", m["a"])
	assert.Equal(t, "2", m["b"])
}

func Test_Header(t *testing.T) {
	p := Param{
		URL: baseUrl + "/header/",
		Header: map[string]string{
			"Header": "header",
		},
	}
	resp, err := p.HttpReq()
	assert.NoError(t, err)
	buf, err := ioutil.ReadAll(resp.Body)
	assert.NoError(t, err)
	assert.Equal(t, "header", string(buf))
}

func Test_EmptyBody(t *testing.T) {
	p := Param{
		URL: baseUrl + "/body/",
	}
	resp, err := p.HttpReq()
	assert.NoError(t, err)
	buf, err := ioutil.ReadAll(resp.Body)
	assert.NoError(t, err)
	assert.Equal(t, "", string(buf))
}

func Test_JsonBody(t *testing.T) {
	type TestData struct {
		A int
		B string
		C []int32
	}
	testMap := TestData{
		A: 1,
		B: "2",
		C: []int32{1, 2, 3},
	}
	buf, _ := json.Marshal(testMap)
	p := Param{
		URL: baseUrl + "/body/",
		BodyContent: BodyContent{
			ContentType: "application/json",
			Content:     buf,
		},
	}
	resp, err := p.HttpReq()
	assert.NoError(t, err)
	buf, err = ioutil.ReadAll(resp.Body)
	assert.NoError(t, err)
	result := TestData{}
	err = json.Unmarshal(buf, &result)
	assert.NoError(t, err)
	assert.Equal(t, testMap, result)
}

func Test_WWWFormData(t *testing.T) {
	p := Param{
		URL: baseUrl + "/body/",
		BodyContent: BodyContent{
			ContentType: "application/x-www-form-urlencoded",
			Content:     map[string]string{"a": "1", "b": "2"},
		},
	}
	resp, err := p.HttpReq()
	assert.NoError(t, err)
	buf, err := ioutil.ReadAll(resp.Body)
	assert.NoError(t, err)
	assert.Equal(t, "a=1&b=2", string(buf))
}

func Test_FormData(t *testing.T) {
	p := Param{
		URL: baseUrl + "/body/",
		BodyContent: BodyContent{
			ContentType: "multipart/form-data",
			Content: map[string]FormDataValue{
				"a": {
					Value: "1",
					Type:  FormDataTypeStr,
				},
				"file": {
					Value:    "http://localhost:8080/file",
					Type:     FormDataTypeFile,
					FileName: "filename.dat",
				},
			},
		},
	}
	resp, err := p.HttpReq()
	assert.NoError(t, err)
	r := resp.Header.Get("Ct")
	assert.Equal(t, "multipart/form-data; boundary=", r[:len("multipart/form-data; boundary=")])
	r = r[len("multipart/form-data; boundary="):]
	reader := multipart.NewReader(resp.Body, r)
	dataMap, err := reader.ReadForm(1 << 20)
	assert.NoError(t, err)
	assert.Equal(t, "1", dataMap.Value["a"][0])
	fileData := dataMap.File["file"][0]
	assert.Equal(t, "filename.dat", fileData.Filename)
	file, err := fileData.Open()
	assert.NoError(t, err)
	buf, err := ioutil.ReadAll(file)
	assert.NoError(t, err)
	assert.Equal(t, "file", string(buf))
}
