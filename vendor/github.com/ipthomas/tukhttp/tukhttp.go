package tukhttp

import (
	"context"
	"io"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/ipthomas/tukcnst"
	"github.com/ipthomas/tukutil"
)

type HTTPRequest struct {
	Act          string
	Version      int
	Server       string
	Method       string
	URL          string
	X_Api_Key    string
	X_Api_Secret string
	PID_OID      string
	PID          string
	SOAPAction   string
	ContentType  string
	Body         []byte
	Timeout      int
	StatusCode   int
	Response     []byte
	DebugMode    bool
}
type TukHTTPInterface interface {
	newRequest() error
}

func NewRequest(i TukHTTPInterface) error {
	return i.newRequest()
}
func (i *HTTPRequest) newRequest() error {
	var err error
	var header http.Header
	header.Add(tukcnst.ACCEPT, tukcnst.ALL)
	header.Add(tukcnst.CONTENT_TYPE, i.ContentType)
	header.Add(tukcnst.CONNECTION, tukcnst.KEEP_ALIVE)
	if i.X_Api_Key != "" && i.X_Api_Secret != "" {
		header.Add("X-API-KEY", i.X_Api_Key)
		header.Add("X-API-SECRET", i.X_Api_Secret)
	}
	if i.SOAPAction != "" {
		header.Set(tukcnst.SOAP_ACTION, i.SOAPAction)
	}
	if i.Timeout == 0 {
		i.Timeout = 15
	}
	i.logReq(header, i.URL, string(i.Body))
	if err = i.sendHttpRequest(header); err != nil {
		log.Println(err.Error())
	}
	i.logRsp(i.StatusCode, string(i.Response))
	return err

}
func (i *HTTPRequest) sendHttpRequest(header http.Header) error {
	var err error
	var req *http.Request
	var rsp *http.Response
	var bytes []byte
	req.Header = header
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(i.Timeout)*time.Second)
	defer cancel()
	if i.Method == http.MethodPost && string(i.Body) != "" {
		req, err = http.NewRequest(i.Method, i.URL, strings.NewReader(string(i.Body)))
	} else {
		req, err = http.NewRequest(i.Method, i.URL, nil)
	}
	if err == nil {
		if rsp, err = http.DefaultClient.Do(req.WithContext(ctx)); err == nil {
			defer rsp.Body.Close()
			if bytes, err = io.ReadAll(rsp.Body); err == nil {
				i.StatusCode = rsp.StatusCode
				i.Response = bytes
			}
		}
	}
	return err
}
func (i *HTTPRequest) logReq(headers interface{}, url string, body string) {
	if i.DebugMode {
		log.Printf("HTTP Request Headers")
		tukutil.Log(headers)
		log.Printf("\nHTTP Request\n-- URL = %s", url)
		if body != "" {
			log.Printf("\n-- Body:\n%s", body)
		}
	}
}
func (i *HTTPRequest) logRsp(statusCode int, response string) {
	if i.DebugMode {
		log.Printf("HTML Response - Status Code = %v\n-- Response--\n%s", statusCode, response)
	}
}
