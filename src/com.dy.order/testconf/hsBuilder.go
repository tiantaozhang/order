package testconf

import (
	"bytes"
	"encoding/json"
	"github.com/Centny/gwf/log"
	"github.com/Centny/gwf/routing"
	"io/ioutil"
	"net/http"
	"net/url"
	"org.cny.uas/usr"
	"strings"
)

func NewUrlValues() *UrlValues {
	v := &UrlValues{}
	v.V = make(map[string][]string)
	return v
}

type UrlValues struct {
	V url.Values
}

func (vmap *UrlValues) Add(k, v string) *UrlValues {
	vmap.V[k] = []string{v}
	return vmap
}

func DoPost(fn func(hs *routing.HTTPSession) routing.HResult, v url.Values, uid int64, uname string) []byte {
	hs, rs := HsBuilder("POST", "", v, uid, uname)
	fn(hs)
	log.I(string(rs.Bytes()))
	return rs.Bytes()
}

func DoPostMap(fn func(hs *routing.HTTPSession) routing.HResult, v url.Values, uid int64, uname string) map[string]interface{} {
	rs := DoPost(fn, v, uid, uname)
	var parse = make(map[string]interface{})
	json.Unmarshal(rs, &parse)
	return parse
}

func DoPostCode(fn func(hs *routing.HTTPSession) routing.HResult, v url.Values, uid int64, uname string) int {
	rs := DoPostMap(fn, v, uid, uname)
	return int(rs["code"].(float64))
}

func DoPostBody(fn func(hs *routing.HTTPSession) routing.HResult, v url.Values, uid int64, uname, body string) []byte {
	hs, rs := HsBuilderBody("POST", "", v, uid, uname, body)
	fn(hs)
	log.I(string(rs.Bytes()))
	return rs.Bytes()
}

func DoPostBodyMap(fn func(hs *routing.HTTPSession) routing.HResult, v url.Values, uid int64, uname, body string) map[string]interface{} {
	rs := DoPostBody(fn, v, uid, uname, body)
	var parse = make(map[string]interface{})
	json.Unmarshal(rs, &parse)
	return parse
}

func DoPostBodyCode(fn func(hs *routing.HTTPSession) routing.HResult, v url.Values, uid int64, uname, body string) int {
	rs := DoPostBodyMap(fn, v, uid, uname, body)
	return int(rs["code"].(float64))
}

func HsBuilder(method, url_ string, v url.Values, uid int64, uname string) (*routing.HTTPSession, *bytes.Buffer) {
	var hs *routing.HTTPSession = &routing.HTTPSession{}

	hs.S = Session{uid, map[string]interface{}{}}
	rs := bytes.NewBuffer(make([]byte, 0))
	hs.W = writer{rs}
	v_ := v.Encode()
	if "GET" == method || "get" == method {
		url_ += "?" + v_
	}
	req, _ := http.NewRequest(method, url_, strings.NewReader(v_))
	if method == "POST" {
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded; param=value")
	}
	hs.R = req
	hs.Kvs = map[string]interface{}{
		"USR": &usr.Usr{
			Tid: uid,
			Usr: uname,
		},
	}
	return hs, rs
}

func HsBuilderBody(method, url_ string, v url.Values, uid int64, uname string, body string) (*routing.HTTPSession, *bytes.Buffer) {
	hs, rs := HsBuilder(method, url_, v, uid, uname)
	hs.R.Body = ioutil.NopCloser(strings.NewReader(body))
	return hs, rs
}

type writer struct {
	B *bytes.Buffer
}

func (w writer) Write(b []byte) (int, error) {
	br := bytes.NewReader(b)
	in, err := br.WriteTo(w.B)
	return int(in), err
}

func (w writer) Header() http.Header {
	return http.Header{}
}

func (w writer) WriteHeader(n int) {
}

type Session struct {
	uid int64
	val map[string]interface{}
}

func (s Session) Val(key string) interface{} {
	return s.val[key]
}
func (s Session) Set(key string, val interface{}) {
	s.val[key] = val
}
func (s Session) SetVal(key string, val interface{}) {
	s.val[key] = val
}
func (s Session) Flush() error {
	return nil
}
