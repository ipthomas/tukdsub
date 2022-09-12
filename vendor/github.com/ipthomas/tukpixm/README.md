# tukpixm
tukpixm provides a golang implementtion of an IHE PIXm PDQ Client

There is currently no authentication implemented. The func (i *PIXmQuery) newRequest() error is used to handle the http request/response and should be amended according to your authentication requirements

Struct PIXmQuery implements the tukpixm.PDQ interface

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

	 PID is the MRN or NHS ID or Regional/xds ID and is required
	 Region_OID is the Regional/XDS OID and is required
	 PIX_URL is the PIXm WS end point and is required.
	 PID_OID is required if the PID is not an NHS ID. If pid length = 10 and no PID_OID is provided, the pid is assumed to be a NHS ID and the PID_OID is set to the NHS offical NHS ID OID (2.16.840.1.113883.2.1.4.1)
	 Timeout sets the http context timeout in seconds and is optional. Default is 5 secs
	 NHS_OID is optional. Default is 2.16.840.1.113883.2.1.4.1
	 Count will be set from the pixm response to the number of patients found matching the query
	 Response will contain the PIXm response in []byte format
	 PIXmResponse will contain the initialised PIXmResponse struc from the Response []byte
	 StatusCode will be set from the PIXm Server http response header statuscode
	 []Patients is any array of PIXPatient structs containing all matched patients. Hopefully just 1 !!

	Example usage:
		pdq := tukpixm.PIXmQuery{
			PID:        "9999999468",
			Region_OID: "2.16.840.1.113883.2.1.3.31.2.1.1",
			PIX_URL:    "http://spirit-test-01.tianispirit.co.uk:8081/SpiritPIXFhir/r4/Patient",
		}
		if err = tukpixm.PDQ(&pdq); err == nil {
			log.Printf("Patient %s %s is registered", pdq.Patients[0].GivenName, pdq.Patients[0].FamilyName)
		} else {
			log.Println(err.Error())
		}

	Running the above example produces the following Log output:

	2022/09/12 14:02:55.510679 tukpixm.go:188: HTTP GET Request Headers

	2022/09/12 14:02:55.510834 tukpixm.go:190: {
	  "Accept": [
	    "*/*"
	  ],
	  "Connection": [
	    "keep-alive"
	  ],
	  "Content-Type": [
	    "application/json"
	  ]
	}

2022/09/12 14:02:55.510860 tukpixm.go:191: HTTP Request
URL = http://spirit-test-01.tianispirit.co.uk:8081/SpiritPIXFhir/r4/Patient?identifier=2.16.840.1.113883.2.1.4.1%7C9999999468&_format=json&_pretty=true
2022/09/12 14:02:55.851605 tukpixm.go:194: HTML Response - Status Code = 200

	{
	  "resourceType": "Bundle",
	  "id": "53c44d32-fb2c-4dfb-b819-db2150e6fa87",
	  "type": "searchset",
	  "total": 1,
	  "link": [ {
	    "relation": "self",
	    "url": "http://spirit-test-01.tianispirit.co.uk:8081/SpiritPIXFhir/r4/Patient?_format=json&_pretty=true&identifier=2.16.840.1.113883.2.1.4.1%7C9999999468"
	  } ],
	  "entry": [ {
	    "fullUrl": "http://spirit-test-01.tianispirit.co.uk:8081/SpiritPIXFhir/r4/Patient/VFNVSy4xNjYxOTc2MjMwMjYxMSYyLjE2Ljg0MC4xLjExMzg4My4yLjEuMy4zMS4yLjEuMS4xLjMuMS4x",
	    "resource": {
	      "resourceType": "Patient",
	      "id": "VFNVSy4xNjYxOTc2MjMwMjYxMSYyLjE2Ljg0MC4xLjExMzg4My4yLjEuMy4zMS4yLjEuMS4xLjMuMS4x",
	      "extension": [ {
	        "url": "http://hl7.org/fhir/StructureDefinition/patient-citizenship",
	        "valueCodeableConcept": {
	          "coding": [ {
	            "code": "GBR"
	          } ]
	        }
	      }, {
	        "url": "http://hl7.org/fhir/StructureDefinition/patient-nationality",
	        "valueCodeableConcept": {
	          "coding": [ {
	            "code": "GBR"
	          } ]
	        }
	      } ],
	      "identifier": [ {
	        "system": "urn:oid:2.16.840.1.113883.2.1.4.1",
	        "value": "9999999468"
	      }, {
	        "use": "usual",
	        "system": "urn:oid:2.16.840.1.113883.2.1.3.31.2.1.1.1.3.1.1",
	        "value": "TSUK.16619762302611"
	      }, {
	        "system": "urn:oid:2.16.840.1.113883.2.1.3.31.2.1.1",
	        "value": "REG.1MWU5C92M2"
	      } ],
	      "active": true,
	      "name": [ {
	        "use": "official",
	        "family": "Testpatient",
	        "given": [ "Nhs" ]
	      } ],
	      "telecom": [ {
	        "system": "phone",
	        "value": "07777661324",
	        "use": "work"
	      }, {
	        "system": "email",
	        "value": "nhs.testpatient@nhs.net",
	        "use": "work"
	      } ],
	      "gender": "male",
	      "birthDate": "1962-04-04",
	      "address": [ {
	        "line": [ "Preston Road" ],
	        "city": "Preston",
	        "state": "Lancashire",
	        "postalCode": "PR1 1PR",
	        "country": "GBR"
	      } ],
	      "maritalStatus": {
	        "coding": [ {
	          "code": "D"
	        } ]
	      },
	      "multipleBirthBoolean": false
	    }
	  } ]
	}

	2022/09/12 14:02:55.852334 tukpixm.go:102: 1 Patient Entries in Response
	2022/09/12 14:02:55.852392 tukpixm.go:122: Set NHS ID 9999999468 2.16.840.1.113883.2.1.4.1
	2022/09/12 14:02:55.852427 tukpixm.go:117: Set PID TSUK.16619762302611 2.16.840.1.113883.2.1.3.31.2.1.1.1.3.1.1
	2022/09/12 14:02:55.852455 tukpixm.go:112: Set Reg ID REG.1MWU5C92M2 2.16.840.1.113883.2.1.3.31.2.1.1
	2022/09/12 14:02:55.852546 tukpixm.go:149: Added Patient 9999999468 to response
	2022/09/12 14:02:55.852569 main.go:84: Patient Nhs Testpatient is registered