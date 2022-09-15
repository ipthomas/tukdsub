# tukpixm
tukpixm provides a golang implementtion of IHE PIXm, IHE PIXv3 and IHE PDQv3 Consumer clients

There is currently no authentication implemented. The github.com/ipthomas/tukhttp package is used to handle the http request/response and should be amended according to your authentication requirements

Struct PDQQuery implements the tukpixm.PDQ interface

	type PDQQuery struct {
	Server     string
	MRN_ID     string
	MRN_OID    string
	NHS_ID     string
	NHS_OID    string
	REG_ID     string
	REG_OID    string
	Server_URL string
	Timeout    int64
	Used_PID   string
	Request    []byte
	Response   []byte
	StatusCode int
	Count      int
	PDQ_ID     string
	PDQ_OID    string
	Patients   []PIXPatient `json:"patients"`
}
	Server must be set to either "pixm" to perform a IHE PIXm query or "pixv3" to perform an IHE PIXv3 query or "pdqv3" to perform an IHE PDQv3 query. The github.com/ipthomas/tukcnst provides constants for each of the valid Server values i.e. tukcnst.PIXm, tukcnst.PIXv3, tukcnst.PDQv3, or you can just use strings!
	
	A patient identifier is required for use in the PDQ. This can be either the MRN id along with the associated OID or the NHS ID (if no NHS OID is provided the NHS assigned OID is used) or the XDS regional ID

	 The REG_OID is the Regional/XDS OID and is required
	 
	 Server_URL is the PIXm WS end point and is required.

	 Timeout is the http context timeout in seconds and is optional. Default is 5 secs

	 PDQ_ID will be set to the ID used for the query
	 PDQ_OID will be set to the OID used for the query
	 
	 Count will equal the number of patients found matching the query
	 Response will contain the PIXm response in []byte format
	 StatusCode will contain the http response header statuscode
	 []Patients will contain an array of PIXPatient structs containing all matched patients. Hopefully just 1 !!

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

	Example usage:
		pdq := tukpixm.PDQQuery{
		Server:     tukcnst.PIXm
		NHS_ID:     "1111111111",
		REG_OID:    "2.16.840.1.113883.2.1.3.31.2.1.1",
		Server_URL: "http://spirit-test-01.tianispirit.co.uk:8081/SpiritPIXFhir/r4/Patient",
	}
	if err = tukpixm.PDQ(&pdq); err == nil {
		if pdq.Count == 1 {
			log.Printf("Patient %s %s is registered", pdq.Patients[0].GivenName, pdq.Patients[0].FamilyName)
		} else {
			log.Printf("Found %v Patients", pdq.Count)
		}
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