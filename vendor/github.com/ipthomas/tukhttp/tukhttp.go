package tukhttp

import (
	"bytes"
	"context"
	"io"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/ipthomas/tukcnst"
	"github.com/ipthomas/tukutil"
)

type CGLRequest struct {
	Request    string
	X_Api_Key  string
	StatusCode int
	Response   []byte
}
type PIXmRequest struct {
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
	URL        string
	Act        string
	Resource   string
	Timeout    int64
	StatusCode int
	Body       []byte
	Response   []byte
}
type ClientRequest struct {
	HttpRequest  *http.Request
	ServerURL    string `json:"serverurl"`
	Act          string `json:"act"`
	User         string `json:"user"`
	Org          string `json:"org"`
	Orgoid       string `json:"orgoid"`
	Role         string `json:"role"`
	NHS          string `json:"nhs"`
	PID          string `json:"pid"`
	PIDOrg       string `json:"pidorg"`
	PIDOID       string `json:"pidoid"`
	FamilyName   string `json:"familyname"`
	GivenName    string `json:"givenname"`
	DOB          string `json:"dob"`
	Gender       string `json:"gender"`
	ZIP          string `json:"zip"`
	Status       string `json:"status"`
	XDWKey       string `json:"xdwkey"`
	ID           int    `json:"id"`
	Task         string `json:"task"`
	Pathway      string `json:"pathway"`
	Version      int    `json:"version"`
	ReturnFormat string `json:"returnformat"`
}
type TukHTTPInterface interface {
	newRequest() error
}

func NewRequest(i TukHTTPInterface) error {
	return i.newRequest()
}
func (i *ClientRequest) newRequest() error {
	req := i.HttpRequest
	req.ParseForm()
	i.Act = req.FormValue(tukcnst.ACT)
	i.User = req.FormValue("user")
	i.Org = req.FormValue("org")
	i.Orgoid = tukutil.GetCodeSystemVal(req.FormValue("org"))
	i.Role = req.FormValue("role")
	i.NHS = req.FormValue("nhs")
	i.PID = req.FormValue("pid")
	i.PIDOrg = req.FormValue("pidorg")
	i.PIDOID = tukutil.GetCodeSystemVal(req.FormValue("pidorg"))
	i.FamilyName = req.FormValue("familyname")
	i.GivenName = req.FormValue("givenname")
	i.DOB = req.FormValue("dob")
	i.Gender = req.FormValue("gender")
	i.ZIP = req.FormValue("zip")
	i.Status = req.FormValue("status")
	i.ID = tukutil.GetIntFromString(req.FormValue("id"))
	i.Task = req.FormValue(tukcnst.TASK)
	i.Pathway = req.FormValue(tukcnst.PATHWAY)
	i.Version = tukutil.GetIntFromString(req.FormValue("version"))
	i.XDWKey = req.FormValue("xdwkey")
	i.ReturnFormat = req.Header.Get(tukcnst.ACCEPT)
	if len(i.XDWKey) > 12 {
		i.Pathway, i.NHS = tukutil.SplitXDWKey(i.XDWKey)
	}
	return nil
}
func (i *SOAPRequest) newRequest() error {
	if i.Timeout == 0 {
		i.Timeout = 15
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(i.Timeout)*time.Second)
	defer cancel()
	req, err := http.NewRequest(http.MethodPost, i.URL, strings.NewReader(string(i.Body)))
	if err != nil {
		return err
	}
	if i.SOAPAction != "" {
		req.Header.Set(tukcnst.SOAP_ACTION, i.SOAPAction)
	}
	req.Header.Set(tukcnst.CONTENT_TYPE, tukcnst.SOAP_XML)
	req.Header.Set(tukcnst.ACCEPT, tukcnst.ALL)
	req.Header.Set(tukcnst.CONNECTION, tukcnst.KEEP_ALIVE)
	i.logRequest(req.Header)

	resp, err := http.DefaultClient.Do(req.WithContext(ctx))
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	i.StatusCode = resp.StatusCode
	i.Response, err = io.ReadAll(resp.Body)
	i.logResponse()
	return err
}
func (i *PIXmRequest) newRequest() error {
	var err error
	var req *http.Request
	if i.Timeout == 0 {
		i.Timeout = 15
	}
	i.URL = i.URL + "?identifier=" + i.PID_OID + "%7C" + i.PID + tukcnst.FORMAT_JSON_PRETTY
	if req, err = http.NewRequest(tukcnst.HTTP_GET, i.URL, nil); err == nil {
		req.Header.Set(tukcnst.CONTENT_TYPE, tukcnst.APPLICATION_JSON)
		req.Header.Set(tukcnst.ACCEPT, tukcnst.ALL)
		req.Header.Set(tukcnst.CONNECTION, tukcnst.KEEP_ALIVE)
		i.logRequest(req.Header)
		ctx, cancel := context.WithTimeout(context.Background(), time.Duration(i.Timeout)*time.Second)
		defer cancel()
		resp, err := http.DefaultClient.Do(req.WithContext(ctx))
		if err != nil {
			return err
		}
		i.StatusCode = resp.StatusCode
		if i.Response, err = io.ReadAll(resp.Body); err != nil {
			log.Println(err.Error())
		}
		defer resp.Body.Close()
		i.logResponse()
		return nil
	}
	return err
}
func (i *CGLRequest) newRequest() error {
	req, _ := http.NewRequest(tukcnst.HTTP_GET, i.Request, nil)
	req.Header.Set(tukcnst.ACCEPT, tukcnst.APPLICATION_JSON)
	req.Header.Set("X-API-KEY", i.X_Api_Key)
	i.logRequest(req.Header)
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(5)*time.Second)
	defer cancel()
	resp, err := http.DefaultClient.Do(req.WithContext(ctx))
	if err != nil {
		return err
	}
	i.StatusCode = resp.StatusCode
	i.Response, err = io.ReadAll(resp.Body)
	defer resp.Body.Close()
	i.logResponse()
	return err
}
func (i *AWS_APIRequest) newRequest() error {
	if i.Timeout == 0 {
		i.Timeout = 5
	}
	var err error
	var req *http.Request
	var resp *http.Response
	client := &http.Client{}
	if req, err = http.NewRequest(http.MethodPost, i.URL+i.Resource, bytes.NewBuffer(i.Body)); err == nil {
		req.Header.Add(tukcnst.CONTENT_TYPE, tukcnst.APPLICATION_JSON_CHARSET_UTF_8)
		ctx, cancel := context.WithTimeout(context.Background(), time.Duration(i.Timeout)*time.Second)
		defer cancel()
		i.logRequest(req.Header)
		if resp, err = client.Do(req.WithContext(ctx)); err == nil {
			log.Printf("Response Status Code %v\n", resp.StatusCode)
			if resp.StatusCode == http.StatusOK {
				i.Response, err = io.ReadAll(resp.Body)
			}
		}
	}
	defer resp.Body.Close()
	i.StatusCode = resp.StatusCode
	i.logResponse()
	return err
}
func (i *AWS_APIRequest) logRequest(headers http.Header) {
	log.Println("HTTP POST Request Headers")
	tukutil.Log(headers)
	log.Printf("HTTP Request\nURL = %s\nTimeout = %v\nMessage body\n%s", i.URL, i.Timeout, string(i.Body))
}
func (i *AWS_APIRequest) logResponse() {
	log.Printf("HTML Response - Status Code = %v\n%s", i.StatusCode, string(i.Response))
}
func (i *SOAPRequest) logRequest(headers http.Header) {
	log.Println("SOAP Request Headers")
	tukutil.Log(headers)
	log.Printf("SOAP Request\nURL = %s\nAction = %s\nTimeout = %v\n\n%s", i.URL, i.SOAPAction, i.Timeout, string(i.Body))
}
func (i *SOAPRequest) logResponse() {
	log.Printf("SOAP Response - Status Code = %v\n%s", i.StatusCode, string(i.Response))
}
func (i *PIXmRequest) logRequest(headers http.Header) {
	log.Println("HTTP GET Request Headers")
	tukutil.Log(headers)
	log.Printf("HTTP Request\nURL = %s\nTimeout = %v", i.URL, i.Timeout)
}
func (i *CGLRequest) logRequest(headers http.Header) {
	log.Println("HTTP GET Request Headers")
	tukutil.Log(headers)
	log.Printf("HTTP Request\nURL = %s - Timeout = %v", i.Request, 5)
}
func (i *PIXmRequest) logResponse() {
	log.Printf("HTML Response - Status Code = %v\n%s", i.StatusCode, string(i.Response))
}
func (i *CGLRequest) logResponse() {
	log.Printf("HTML Response - Status Code = %v\n%s", i.StatusCode, string(i.Response))
}
