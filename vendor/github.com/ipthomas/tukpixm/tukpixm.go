// tukpixm provides a golang implementtion of an IHE PIXm PDQ Client
//
// There is currently no authentication implemented. The func (i *PIXmQuery) newRequest() error is used to handle the http request/response and should be amended according to your authentication requirements
//
// Struct PIXmQuery implements the tukpixm.PDQ interface
//
//	type PIXmQuery struct {
//		Count        int          `json:"count"`
//		PID          string       `json:"pid"`
//		PIDOID       string       `json:"pidoid"`
//		PIX_URL      string       `json:"pixurl"`
//		NHS_OID      string       `json:"nhsoid"`
//		Region_OID   string       `json:"regionoid"`
//		Timeout      int64        `json:"timeout"`
//		StatusCode   int          `json:"statuscode"`
//		Response     []byte       `json:"response"`
//		PIXmResponse PIXmResponse `json:"pixmresponse"`
//		Patients     []PIXPatient `json:"patients"`
//	}
//
//	 PID is the MRN or NHS ID or Regional/xds ID and is required
//	 Region_OID is the Regional/XDS OID and is required
//	 PIX_URL is the PIXm WS end point and is required.
//	 PID_OID is required if the PID is not an NHS ID. If pid length = 10 and no PID_OID is provided, the pid is assumed to be a NHS ID and the PID_OID is set to the NHS offical NHS ID OID (2.16.840.1.113883.2.1.4.1)
//	 Timeout sets the http context timeout in seconds and is optional. Default is 5 secs
//	 NHS_OID is optional. Default is 2.16.840.1.113883.2.1.4.1
//	 Count will be set from the pixm response to the number of patients found matching the query
//	 Response will contain the PIXm response in []byte format
//	 PIXmResponse will contain the initialised PIXmResponse struc from the Response []byte
//	 StatusCode will be set from the PIXm Server http response header statuscode
//	 []Patients is any array of PIXPatient structs containing all matched patients. Hopefully just 1 !!
//
//	Example usage:
//		pdq := tukpixm.PIXmQuery{
//			PID:        "9999999468",
//			Region_OID: "2.16.840.1.113883.2.1.3.31.2.1.1",
//			PIX_URL:    "http://spirit-test-01.tianispirit.co.uk:8081/SpiritPIXFhir/r4/Patient",
//		}
//		if err = tukpixm.PDQ(&pdq); err == nil {
//			log.Printf("Patient %s %s is registered", pdq.Patients[0].GivenName, pdq.Patients[0].FamilyName)
//		} else {
//			log.Println(err.Error())
//		}
//
//	Running the above example produces the following Log output:
//
//	2022/09/12 14:02:55.510679 tukpixm.go:188: HTTP GET Request Headers
//
//	2022/09/12 14:02:55.510834 tukpixm.go:190: {
//	  "Accept": [
//	    "*/*"
//	  ],
//	  "Connection": [
//	    "keep-alive"
//	  ],
//	  "Content-Type": [
//	    "application/json"
//	  ]
//	}
//
// 2022/09/12 14:02:55.510860 tukpixm.go:191: HTTP Request
// URL = http://spirit-test-01.tianispirit.co.uk:8081/SpiritPIXFhir/r4/Patient?identifier=2.16.840.1.113883.2.1.4.1%7C9999999468&_format=json&_pretty=true
// 2022/09/12 14:02:55.851605 tukpixm.go:194: HTML Response - Status Code = 200
//
//	{
//	  "resourceType": "Bundle",
//	  "id": "53c44d32-fb2c-4dfb-b819-db2150e6fa87",
//	  "type": "searchset",
//	  "total": 1,
//	  "link": [ {
//	    "relation": "self",
//	    "url": "http://spirit-test-01.tianispirit.co.uk:8081/SpiritPIXFhir/r4/Patient?_format=json&_pretty=true&identifier=2.16.840.1.113883.2.1.4.1%7C9999999468"
//	  } ],
//	  "entry": [ {
//	    "fullUrl": "http://spirit-test-01.tianispirit.co.uk:8081/SpiritPIXFhir/r4/Patient/VFNVSy4xNjYxOTc2MjMwMjYxMSYyLjE2Ljg0MC4xLjExMzg4My4yLjEuMy4zMS4yLjEuMS4xLjMuMS4x",
//	    "resource": {
//	      "resourceType": "Patient",
//	      "id": "VFNVSy4xNjYxOTc2MjMwMjYxMSYyLjE2Ljg0MC4xLjExMzg4My4yLjEuMy4zMS4yLjEuMS4xLjMuMS4x",
//	      "extension": [ {
//	        "url": "http://hl7.org/fhir/StructureDefinition/patient-citizenship",
//	        "valueCodeableConcept": {
//	          "coding": [ {
//	            "code": "GBR"
//	          } ]
//	        }
//	      }, {
//	        "url": "http://hl7.org/fhir/StructureDefinition/patient-nationality",
//	        "valueCodeableConcept": {
//	          "coding": [ {
//	            "code": "GBR"
//	          } ]
//	        }
//	      } ],
//	      "identifier": [ {
//	        "system": "urn:oid:2.16.840.1.113883.2.1.4.1",
//	        "value": "9999999468"
//	      }, {
//	        "use": "usual",
//	        "system": "urn:oid:2.16.840.1.113883.2.1.3.31.2.1.1.1.3.1.1",
//	        "value": "TSUK.16619762302611"
//	      }, {
//	        "system": "urn:oid:2.16.840.1.113883.2.1.3.31.2.1.1",
//	        "value": "REG.1MWU5C92M2"
//	      } ],
//	      "active": true,
//	      "name": [ {
//	        "use": "official",
//	        "family": "Testpatient",
//	        "given": [ "Nhs" ]
//	      } ],
//	      "telecom": [ {
//	        "system": "phone",
//	        "value": "07777661324",
//	        "use": "work"
//	      }, {
//	        "system": "email",
//	        "value": "nhs.testpatient@nhs.net",
//	        "use": "work"
//	      } ],
//	      "gender": "male",
//	      "birthDate": "1962-04-04",
//	      "address": [ {
//	        "line": [ "Preston Road" ],
//	        "city": "Preston",
//	        "state": "Lancashire",
//	        "postalCode": "PR1 1PR",
//	        "country": "GBR"
//	      } ],
//	      "maritalStatus": {
//	        "coding": [ {
//	          "code": "D"
//	        } ]
//	      },
//	      "multipleBirthBoolean": false
//	    }
//	  } ]
//	}
//
// 2022/09/12 14:02:55.852334 tukpixm.go:102: 1 Patient Entries in Response
// 2022/09/12 14:02:55.852392 tukpixm.go:122: Set NHS ID 9999999468 2.16.840.1.113883.2.1.4.1
// 2022/09/12 14:02:55.852427 tukpixm.go:117: Set PID TSUK.16619762302611 2.16.840.1.113883.2.1.3.31.2.1.1.1.3.1.1
// 2022/09/12 14:02:55.852455 tukpixm.go:112: Set Reg ID REG.1MWU5C92M2 2.16.840.1.113883.2.1.3.31.2.1.1
// 2022/09/12 14:02:55.852546 tukpixm.go:149: Added Patient 9999999468 to response
// 2022/09/12 14:02:55.852569 main.go:84: Patient Nhs Testpatient is registered
package tukpixm

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
	"strings"
	"time"

	cnst "github.com/ipthomas/tukcnst"
	// github.com/aws/aws-lambda-go
)

type PIXmResponse struct {
	ResourceType string `json:"resourceType"`
	ID           string `json:"id"`
	Type         string `json:"type"`
	Total        int    `json:"total"`
	Link         []struct {
		Relation string `json:"relation"`
		URL      string `json:"url"`
	} `json:"link"`
	Entry []struct {
		FullURL  string `json:"fullUrl"`
		Resource struct {
			ResourceType string `json:"resourceType"`
			ID           string `json:"id"`
			Identifier   []struct {
				Use    string `json:"use,omitempty"`
				System string `json:"system"`
				Value  string `json:"value"`
			} `json:"identifier"`
			Active bool `json:"active"`
			Name   []struct {
				Use    string   `json:"use"`
				Family string   `json:"family"`
				Given  []string `json:"given"`
			} `json:"name"`
			Gender    string `json:"gender"`
			BirthDate string `json:"birthDate"`
			Address   []struct {
				Use        string   `json:"use"`
				Line       []string `json:"line"`
				City       string   `json:"city"`
				PostalCode string   `json:"postalCode"`
				Country    string   `json:"country"`
			} `json:"address"`
		} `json:"resource"`
	} `json:"entry"`
}
type PIXmQuery struct {
	Count        int          `json:"count"`
	PID          string       `json:"pid"`
	PIDOID       string       `json:"pidoid"`
	PIX_URL      string       `json:"pixurl"`
	NHS_OID      string       `json:"nhsoid"`
	Region_OID   string       `json:"regionoid"`
	Timeout      int64        `json:"timeout"`
	StatusCode   int          `json:"statuscode"`
	Response     []byte       `json:"response"`
	PIXmResponse PIXmResponse `json:"pixmresponse"`
	Patients     []PIXPatient `json:"patients"`
}
type PIXPatient struct {
	PIDOID     string `json:"pidoid"`
	PID        string `json:"pid"`
	REGOID     string `json:"regoid"`
	REGID      string `json:"regid"`
	NHSOID     string `json:"nhsoid"`
	NHSID      string `json:"nhsid"`
	GivenName  string `json:"givenname"`
	FamilyName string `json:"familyname"`
	Gender     string `json:"gender"`
	BirthDate  string `json:"birthdate"`
	Street     string `json:"street"`
	Town       string `json:"town"`
	City       string `json:"city"`
	State      string `json:"state"`
	Country    string `json:"country"`
	Zip        string `json:"zip"`
}
type PIXmInterface interface {
	pdq() error
}

func PDQ(i PIXmInterface) error {
	return i.pdq()
}
func (i *PIXmQuery) pdq() error {
	if err := i.newRequest(); err != nil {
		return err
	}
	if strings.Contains(string(i.Response), "Error") {
		return errors.New(string(i.Response))
	}
	i.PIXmResponse = PIXmResponse{}
	if err := json.Unmarshal(i.Response, &i.PIXmResponse); err != nil {
		return err
	}
	log.Printf("%v Patient Entries in Response", i.PIXmResponse.Total)
	i.Count = i.PIXmResponse.Total
	if i.Count > 0 {
		for cnt := 0; cnt < len(i.PIXmResponse.Entry); cnt++ {
			rsppat := i.PIXmResponse.Entry[cnt]
			tukpat := PIXPatient{}
			for _, id := range rsppat.Resource.Identifier {
				if id.System == cnst.URN_OID_PREFIX+i.Region_OID {
					tukpat.REGID = id.Value
					tukpat.REGOID = i.Region_OID
					log.Printf("Set Reg ID %s %s", tukpat.REGID, tukpat.REGOID)
				}
				if id.Use == "usual" {
					tukpat.PID = id.Value
					tukpat.PIDOID = strings.Split(id.System, ":")[2]
					log.Printf("Set PID %s %s", tukpat.PID, tukpat.PIDOID)
				}
				if id.System == cnst.URN_OID_PREFIX+i.NHS_OID {
					tukpat.NHSID = id.Value
					tukpat.NHSOID = i.NHS_OID
					log.Printf("Set NHS ID %s %s", tukpat.NHSID, tukpat.NHSOID)
				}
			}
			gn := ""
			for _, name := range rsppat.Resource.Name {
				for _, n := range name.Given {
					gn = gn + n + " "
				}
			}

			tukpat.GivenName = strings.TrimSuffix(gn, " ")
			tukpat.FamilyName = rsppat.Resource.Name[0].Family
			tukpat.BirthDate = strings.ReplaceAll(rsppat.Resource.BirthDate, "-", "")
			tukpat.Gender = rsppat.Resource.Gender

			if len(rsppat.Resource.Address) > 0 {
				tukpat.Zip = rsppat.Resource.Address[0].PostalCode
				if len(rsppat.Resource.Address[0].Line) > 0 {
					tukpat.Street = rsppat.Resource.Address[0].Line[0]
					if len(rsppat.Resource.Address[0].Line) > 1 {
						tukpat.Town = rsppat.Resource.Address[0].Line[1]
					}
				}
				tukpat.City = rsppat.Resource.Address[0].City
				tukpat.Country = rsppat.Resource.Address[0].Country
			}
			i.Patients = append(i.Patients, tukpat)
			log.Printf("Added Patient %s to response", tukpat.NHSID)
		}
	} else {
		log.Println("patient is not registered")
	}
	return nil
}
func (i *PIXmQuery) newRequest() error {
	if i.PID == "" || i.Region_OID == "" || i.PIX_URL == "" {
		return errors.New("invalid request, not all mandated values provided in pixmquery (pid, region_oid and pix_url)")
	}
	if i.Timeout == 0 {
		i.Timeout = 5
	}
	if i.NHS_OID == "" {
		i.NHS_OID = "2.16.840.1.113883.2.1.4.1"
	}
	if i.PIDOID == "" && len(i.PID) == 10 {
		i.PIDOID = i.NHS_OID
	}
	i.PIX_URL = i.PIX_URL + "?identifier=" + i.PIDOID + "%7C" + i.PID + cnst.FORMAT_JSON_PRETTY
	req, _ := http.NewRequest(cnst.HTTP_GET, i.PIX_URL, nil)
	req.Header.Set(cnst.CONTENT_TYPE, cnst.APPLICATION_JSON)
	req.Header.Set(cnst.ACCEPT, cnst.ALL)
	req.Header.Set(cnst.CONNECTION, cnst.KEEP_ALIVE)
	i.logRequest(req.Header)
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(i.Timeout)*time.Second)
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
func (i *PIXmQuery) logRequest(headers http.Header) {
	log.Println("HTTP GET Request Headers")
	b, _ := json.MarshalIndent(headers, "", "  ")
	log.Println(string(b))
	log.Printf("HTTP Request\nURL = %s", i.PIX_URL)
}
func (i *PIXmQuery) logResponse() {
	log.Printf("HTML Response - Status Code = %v\n%s", i.StatusCode, string(i.Response))
}

// functions to support AWS Lambda deployment. Uncomment and change package name from tukpixm to main
// func main() {
// 	lambda.Start(Handle_Request)
// }
// func Handle_Request(req events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
// 	log.SetFlags(log.Ldate | log.Ltime | log.Lmicroseconds | log.Lshortfile)
// 	log.Printf("Processing API Gateway Proxy %s %s request", req.HTTPMethod, req.Path)

// 	pdq := PIXmQuery{
// 		PID:        req.QueryStringParameters["pid"],
// 		Region_OID: req.QueryStringParameters["regionoid"],
// 		PIX_URL:    req.QueryStringParameters["pixurl"],
// 	}
// 	if err := PDQ(&pdq); err != nil {
// 		pdq.StatusCode = http.StatusInternalServerError
// 		pdq.Response = []byte(err.Error())
// 	}
// 	return UnhandledRequest(pdq.StatusCode, pdq.Response)
// }
// func UnhandledRequest(status int, body interface{}) (*events.APIGatewayProxyResponse, error) {
// 	resp := events.APIGatewayProxyResponse{Headers: map[string]string{cnst.CONTENT_TYPE: cnst.APPLICATION_JSON}}
// 	resp.StatusCode = status
// 	stringBody, _ := json.Marshal(body)
// 	resp.Body = string(stringBody)
// 	return &resp, nil
// }
