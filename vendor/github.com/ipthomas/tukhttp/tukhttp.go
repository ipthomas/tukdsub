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

var DebugMode = true

type CGLRequest struct {
	URL          string
	X_Api_Key    string
	X_Api_Secret string
	StatusCode   int
	Response     []byte
}
type PIXmRequest struct {
	URL        string
	PID_OID    string
	PID        string
	Timeout    int64
	StatusCode int
	Response   []byte
}
type MHDRequest struct {
	URL        string
	PID_OID    string
	PID        string
	Timeout    int64
	StatusCode int
	Response   []byte
}
type SOAPRequest struct {
	URL        string
	SOAPAction string
	Timeout    int64
	StatusCode int
	Body       []byte
	Response   []byte
}
type AWS_APIRequest struct {
	URL         string
	Version     int
	Act         string
	Resource    string
	ContentType string
	Timeout     int64
	StatusCode  int
	Body        []byte
	Response    []byte
}

type TukHTTPInterface interface {
	newRequest() error
}

var debugMode bool = false

func SetDebugMode(debug bool) {
	debugMode = debug
}
func NewRequest(i TukHTTPInterface) error {
	return i.newRequest()
}
func (i *SOAPRequest) newRequest() error {
	var err error
	var header *http.Header
	if i.Timeout == 0 {
		i.Timeout = 15
	}
	if i.SOAPAction != "" {
		header.Set(tukcnst.SOAP_ACTION, i.SOAPAction)
	}
	header.Set(tukcnst.CONTENT_TYPE, tukcnst.SOAP_XML)
	header.Set(tukcnst.ACCEPT, tukcnst.ALL)
	header.Set(tukcnst.CONNECTION, tukcnst.KEEP_ALIVE)
	i.logRequest(*header)
	if i.StatusCode, i.Response, err = sendHttpRequest(http.MethodPost, header, i.URL, string(i.Body), int(i.Timeout)); err != nil {
		log.Println(err.Error())
		debugMode = true
	}
	i.logResponse()
	return err
}
func (i *PIXmRequest) newRequest() error {
	var err error
	var header *http.Header
	if i.Timeout == 0 {
		i.Timeout = 15
	}
	header.Set(tukcnst.ACCEPT, tukcnst.ALL)
	header.Set(tukcnst.CONNECTION, tukcnst.KEEP_ALIVE)
	i.URL = i.URL + "?identifier=" + i.PID_OID + "|" + i.PID + tukcnst.FORMAT_JSON_PRETTY
	i.logRequest(*header)
	if i.StatusCode, i.Response, err = sendHttpRequest(http.MethodGet, header, i.URL, "", int(i.Timeout)); err != nil {
		log.Println(err.Error())
		debugMode = true
	}
	i.logResponse()
	return err
}
func (i *MHDRequest) newRequest() error {
	var err error
	var header *http.Header
	if i.Timeout == 0 {
		i.Timeout = 15
	}
	i.URL = i.URL + "?patient.identifier=urn:oid:=" + i.PID_OID + "|" + i.PID + tukcnst.FORMAT_JSON_PRETTY
	i.logRequest(*header)
	if i.StatusCode, i.Response, err = sendHttpRequest(http.MethodGet, header, i.URL, "", int(i.Timeout)); err != nil {
		log.Println(err.Error())
		debugMode = true
	}
	i.logResponse()
	return err
}
func (i *CGLRequest) newRequest() error {
	var err error
	var header *http.Header
	header.Set(tukcnst.ACCEPT, tukcnst.APPLICATION_JSON)
	header.Set("X-API-KEY", i.X_Api_Key)
	header.Set("X-API-SECRET", i.X_Api_Secret)
	i.logRequest(*header)
	if i.StatusCode, i.Response, err = sendHttpRequest(http.MethodGet, header, i.URL, "", 15); err != nil {
		log.Println(err.Error())
		debugMode = true
	}
	i.logResponse()
	return err
}
func (i *AWS_APIRequest) newRequest() error {
	var err error
	var header *http.Header
	header.Add(tukcnst.CONTENT_TYPE, i.ContentType)
	if i.Timeout == 0 {
		i.Timeout = 5
	}
	i.logRequest(*header)
	if i.StatusCode, i.Response, err = sendHttpRequest(http.MethodPost, header, i.URL+i.Resource, string(i.Body), int(i.Timeout)); err != nil {
		log.Println(err.Error())
		debugMode = true
	}
	i.logResponse()
	return err
}
func sendHttpRequest(method string, header *http.Header, url string, body string, timeout int) (int, []byte, error) {
	var err error
	var req *http.Request
	var rsp *http.Response
	var bytes []byte
	var bodycontent *strings.Reader
	req.Header = *header
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(timeout)*time.Second)
	defer cancel()
	if method == http.MethodPost && body != "" {
		bodycontent = strings.NewReader(body)
	}
	if req, err = http.NewRequest(method, url, bodycontent); err == nil {
		if rsp, err = http.DefaultClient.Do(req.WithContext(ctx)); err == nil {
			defer rsp.Body.Close()
			if bytes, err = io.ReadAll(rsp.Body); err == nil {
				return rsp.StatusCode, bytes, err
			}
		}
	}
	if err != nil {
		log.Println(err.Error())
		debugMode = true
	}
	return rsp.StatusCode, bytes, err
}
func (i *AWS_APIRequest) logRequest(headers http.Header) {
	logReq(i.StatusCode, i.URL, string(i.Body))
}
func (i *AWS_APIRequest) logResponse() {
	logRsp(i.StatusCode, string(i.Response))
}
func (i *SOAPRequest) logRequest(headers http.Header) {
	logReq(headers, i.URL, string(i.Body))
}
func (i *SOAPRequest) logResponse() {
	logRsp(i.StatusCode, string(i.Response))

}
func (i *PIXmRequest) logRequest(headers http.Header) {
	logReq(headers, i.URL, "")

}
func (i *MHDRequest) logRequest(headers http.Header) {
	logReq(headers, i.URL, "")

}
func (i *PIXmRequest) logResponse() {
	logRsp(i.StatusCode, string(i.Response))

}
func (i *MHDRequest) logResponse() {
	logRsp(i.StatusCode, string(i.Response))

}
func (i *CGLRequest) logRequest(headers http.Header) {
	logReq(headers, i.URL, "")

}
func (i *CGLRequest) logResponse() {
	logRsp(i.StatusCode, string(i.Response))
}
func logReq(headers interface{}, url string, body string) {
	if debugMode {
		log.Printf("HTTP GET Request Headers")
		tukutil.Log(headers)
		log.Printf("\nHTTP Request\n-- URL = %s", url)
		if body != "" {
			log.Printf("\n-- Body:\n%s", body)
		}
	}
}
func logRsp(statusCode int, response string) {
	if debugMode {
		log.Printf("HTML Response - Status Code = %v\n-- Response--\n%s", statusCode, response)
	}
}
