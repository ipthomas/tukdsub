// tukpdq provides a golang implementtion of, IHE PIXm,IHE PIXv3 and IHE PDQv3 Client Consumers
//
// There is currently no authentication implemented. The func (i *PDQQuery) newRequest() error is used to handle the http request/response and should be amended according to your authentication requirements
//
// Struct PDQQuery implements the tukpdq.PDQ() interface
//
//	type PDQQuery struct {
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
//		pdq := tukpdq.PIXmQuery{
//			PID:        "9999999468",
//			Region_OID: "2.16.840.1.113883.2.1.3.31.2.1.1",
//			PIX_URL:    "http://spirit-test-01.tianispirit.co.uk:8081/SpiritPIXFhir/r4/Patient",
//		}
//		if err = tukpdq.PDQ(&pdq); err == nil {
//			log.Printf("Patient %s %s is registered", pdq.Patients[0].GivenName, pdq.Patients[0].FamilyName)
//		} else {
//			log.Println(err.Error())
//		}
//
//	Running the above example produces the following Log output:
//
//	2022/09/12 14:02:55.510679 tukpdq.go:188: HTTP GET Request Headers
//
//	2022/09/12 14:02:55.510834 tukpdq.go:190: {
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
// 2022/09/12 14:02:55.510860 tukpdq.go:191: HTTP Request
// URL = http://spirit-test-01.tianispirit.co.uk:8081/SpiritPIXFhir/r4/Patient?identifier=2.16.840.1.113883.2.1.4.1%7C9999999468&_format=json&_pretty=true
// 2022/09/12 14:02:55.851605 tukpdq.go:194: HTML Response - Status Code = 200
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
// 2022/09/12 14:02:55.852334 tukpdq.go:102: 1 Patient Entries in Response
// 2022/09/12 14:02:55.852392 tukpdq.go:122: Set NHS ID 9999999468 2.16.840.1.113883.2.1.4.1
// 2022/09/12 14:02:55.852427 tukpdq.go:117: Set PID TSUK.16619762302611 2.16.840.1.113883.2.1.3.31.2.1.1.1.3.1.1
// 2022/09/12 14:02:55.852455 tukpdq.go:112: Set Reg ID REG.1MWU5C92M2 2.16.840.1.113883.2.1.3.31.2.1.1
// 2022/09/12 14:02:55.852546 tukpdq.go:149: Added Patient 9999999468 to response
// 2022/09/12 14:02:55.852569 main.go:84: Patient Nhs Testpatient is registered
package tukpdq

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"errors"
	"log"
	"strconv"
	"strings"
	"text/template"

	cnst "github.com/ipthomas/tukcnst"
	"github.com/ipthomas/tukhttp"
	util "github.com/ipthomas/tukutil"
)

type PDQQuery struct {
	Server       string       `json:"server"`
	MRN_ID       string       `json:"mrnid"`
	MRN_OID      string       `json:"mrnoid"`
	NHS_ID       string       `json:"nhsid"`
	NHS_OID      string       `json:"nhsoid"`
	REG_ID       string       `json:"regid"`
	REG_OID      string       `json:"regoid"`
	Server_URL   string       `json:"serverurl"`
	Timeout      int64        `json:"timeout"`
	Used_PID     string       `json:"usedpid"`
	Used_PID_OID string       `json:"usedpidoid"`
	Request      []byte       `json:"request"`
	Response     []byte       `json:"response"`
	StatusCode   int          `json:"statuscode"`
	Count        int          `json:"count"`
	Patients     []PIXPatient `json:"patients"`
}

type PDQv3Response struct {
	XMLName xml.Name `xml:"Envelope"`
	S       string   `xml:"S,attr"`
	Env     string   `xml:"env,attr"`
	Header  struct {
		Action struct {
			Xmlns string `xml:"xmlns,attr"`
			S     string `xml:"S,attr"`
			Env   string `xml:"env,attr"`
		} `xml:"Action"`
		MessageID struct {
			Xmlns string `xml:"xmlns,attr"`
			S     string `xml:"S,attr"`
			Env   string `xml:"env,attr"`
		} `xml:"MessageID"`
		RelatesTo struct {
			Xmlns string `xml:"xmlns,attr"`
			S     string `xml:"S,attr"`
			Env   string `xml:"env,attr"`
		} `xml:"RelatesTo"`
		To struct {
			Xmlns string `xml:"xmlns,attr"`
			S     string `xml:"S,attr"`
			Env   string `xml:"env,attr"`
		} `xml:"To"`
	} `xml:"Header"`
	Body struct {
		PRPAIN201306UV02 struct {
			Xmlns      string `xml:"xmlns,attr"`
			ITSVersion string `xml:"ITSVersion,attr"`
			ID         struct {
				Extension string `xml:"extension,attr"`
				Root      string `xml:"root,attr"`
			} `xml:"id"`
			CreationTime struct {
				Value string `xml:"value,attr"`
			} `xml:"creationTime"`
			VersionCode struct {
				Code string `xml:"code,attr"`
			} `xml:"versionCode"`
			InteractionId struct {
				Extension string `xml:"extension,attr"`
				Root      string `xml:"root,attr"`
			} `xml:"interactionId"`
			ProcessingCode struct {
				Code string `xml:"code,attr"`
			} `xml:"processingCode"`
			ProcessingModeCode struct {
				Code string `xml:"code,attr"`
			} `xml:"processingModeCode"`
			AcceptAckCode struct {
				Code string `xml:"code,attr"`
			} `xml:"acceptAckCode"`
			Receiver struct {
				TypeCode string `xml:"typeCode,attr"`
				Device   struct {
					ClassCode      string `xml:"classCode,attr"`
					DeterminerCode string `xml:"determinerCode,attr"`
					ID             struct {
						AssigningAuthorityName string `xml:"assigningAuthorityName,attr"`
						Root                   string `xml:"root,attr"`
					} `xml:"id"`
					AsAgent struct {
						ClassCode               string `xml:"classCode,attr"`
						RepresentedOrganization struct {
							ClassCode      string `xml:"classCode,attr"`
							DeterminerCode string `xml:"determinerCode,attr"`
							ID             struct {
								AssigningAuthorityName string `xml:"assigningAuthorityName,attr"`
								Root                   string `xml:"root,attr"`
							} `xml:"id"`
						} `xml:"representedOrganization"`
					} `xml:"asAgent"`
				} `xml:"device"`
			} `xml:"receiver"`
			Sender struct {
				TypeCode string `xml:"typeCode,attr"`
				Device   struct {
					ClassCode      string `xml:"classCode,attr"`
					DeterminerCode string `xml:"determinerCode,attr"`
					ID             struct {
						Root string `xml:"root,attr"`
					} `xml:"id"`
					AsAgent struct {
						ClassCode               string `xml:"classCode,attr"`
						RepresentedOrganization struct {
							ClassCode      string `xml:"classCode,attr"`
							DeterminerCode string `xml:"determinerCode,attr"`
							ID             struct {
								Root string `xml:"root,attr"`
							} `xml:"id"`
						} `xml:"representedOrganization"`
					} `xml:"asAgent"`
				} `xml:"device"`
			} `xml:"sender"`
			Acknowledgement struct {
				TypeCode struct {
					Code string `xml:"code,attr"`
				} `xml:"typeCode"`
				TargetMessage struct {
					ID struct {
						Extension string `xml:"extension,attr"`
						Root      string `xml:"root,attr"`
					} `xml:"id"`
				} `xml:"targetMessage"`
			} `xml:"acknowledgement"`
			ControlActProcess struct {
				ClassCode string `xml:"classCode,attr"`
				MoodCode  string `xml:"moodCode,attr"`
				Code      struct {
					Code       string `xml:"code,attr"`
					CodeSystem string `xml:"codeSystem,attr"`
				} `xml:"code"`
				Subject struct {
					ContextConductionInd string `xml:"contextConductionInd,attr"`
					TypeCode             string `xml:"typeCode,attr"`
					RegistrationEvent    struct {
						ClassCode string `xml:"classCode,attr"`
						MoodCode  string `xml:"moodCode,attr"`
						ID        struct {
							NullFlavor string `xml:"nullFlavor,attr"`
						} `xml:"id"`
						StatusCode struct {
							Code string `xml:"code,attr"`
						} `xml:"statusCode"`
						Subject1 struct {
							TypeCode string `xml:"typeCode,attr"`
							Patient  struct {
								ClassCode string `xml:"classCode,attr"`
								ID        []struct {
									AssigningAuthorityName string `xml:"assigningAuthorityName,attr"`
									Extension              string `xml:"extension,attr"`
									Root                   string `xml:"root,attr"`
								} `xml:"id"`
								StatusCode struct {
									Code string `xml:"code,attr"`
								} `xml:"statusCode"`
								EffectiveTime struct {
									Value string `xml:"value,attr"`
								} `xml:"effectiveTime"`
								PatientPerson struct {
									ClassCode      string `xml:"classCode,attr"`
									DeterminerCode string `xml:"determinerCode,attr"`
									Name           struct {
										Use    string `xml:"use,attr"`
										Given  string `xml:"given"`
										Family string `xml:"family"`
									} `xml:"name"`
									AdministrativeGenderCode struct {
										Code           string `xml:"code,attr"`
										CodeSystem     string `xml:"codeSystem,attr"`
										CodeSystemName string `xml:"codeSystemName,attr"`
									} `xml:"administrativeGenderCode"`
									BirthTime struct {
										Value string `xml:"value,attr"`
									} `xml:"birthTime"`
									DeceasedInd struct {
										Value string `xml:"value,attr"`
									} `xml:"deceasedInd"`
									MultipleBirthInd struct {
										Value string `xml:"value,attr"`
									} `xml:"multipleBirthInd"`
									Addr struct {
										StreetAddressLine string `xml:"streetAddressLine"`
										City              string `xml:"city"`
										State             string `xml:"state"`
										PostalCode        string `xml:"postalCode"`
										Country           string `xml:"country"`
									} `xml:"addr"`
									BirthPlace struct {
										Addr struct {
											City string `xml:"city"`
										} `xml:"addr"`
									} `xml:"birthPlace"`
								} `xml:"patientPerson"`
								ProviderOrganization struct {
									ClassCode      string `xml:"classCode,attr"`
									DeterminerCode string `xml:"determinerCode,attr"`
									ID             struct {
										NullFlavor string `xml:"nullFlavor,attr"`
									} `xml:"id"`
									ContactParty struct {
										ClassCode string `xml:"classCode,attr"`
									} `xml:"contactParty"`
								} `xml:"providerOrganization"`
								SubjectOf1 struct {
									TypeCode              string `xml:"typeCode,attr"`
									QueryMatchObservation struct {
										ClassCode string `xml:"classCode,attr"`
										MoodCode  string `xml:"moodCode,attr"`
										Code      struct {
											Code       string `xml:"code,attr"`
											CodeSystem string `xml:"codeSystem,attr"`
										} `xml:"code"`
										Value struct {
											Value string `xml:"value,attr"`
											Xsi   string `xml:"xsi,attr"`
											Type  string `xml:"type,attr"`
										} `xml:"value"`
									} `xml:"queryMatchObservation"`
								} `xml:"subjectOf1"`
							} `xml:"patient"`
						} `xml:"subject1"`
						Custodian struct {
							TypeCode       string `xml:"typeCode,attr"`
							AssignedEntity struct {
								ClassCode string `xml:"classCode,attr"`
								ID        struct {
									NullFlavor string `xml:"nullFlavor,attr"`
								} `xml:"id"`
							} `xml:"assignedEntity"`
						} `xml:"custodian"`
					} `xml:"registrationEvent"`
				} `xml:"subject"`
				QueryAck struct {
					QueryId struct {
						Extension string `xml:"extension,attr"`
						Root      string `xml:"root,attr"`
					} `xml:"queryId"`
					StatusCode struct {
						Code string `xml:"code,attr"`
					} `xml:"statusCode"`
					QueryResponseCode struct {
						Code string `xml:"code,attr"`
					} `xml:"queryResponseCode"`
					ResultTotalQuantity struct {
						Value string `xml:"value,attr"`
					} `xml:"resultTotalQuantity"`
					ResultCurrentQuantity struct {
						Value string `xml:"value,attr"`
					} `xml:"resultCurrentQuantity"`
					ResultRemainingQuantity struct {
						Value string `xml:"value,attr"`
					} `xml:"resultRemainingQuantity"`
				} `xml:"queryAck"`
				QueryByParameter struct {
					QueryId struct {
						Extension string `xml:"extension,attr"`
						Root      string `xml:"root,attr"`
					} `xml:"queryId"`
					StatusCode struct {
						Code string `xml:"code,attr"`
					} `xml:"statusCode"`
					ResponseModalityCode struct {
						Code string `xml:"code,attr"`
					} `xml:"responseModalityCode"`
					ResponsePriorityCode struct {
						Code string `xml:"code,attr"`
					} `xml:"responsePriorityCode"`
					MatchCriterionList string `xml:"matchCriterionList"`
					ParameterList      struct {
						LivingSubjectId struct {
							Value struct {
								Extension string `xml:"extension,attr"`
							} `xml:"value"`
							SemanticsText string `xml:"semanticsText"`
						} `xml:"livingSubjectId"`
					} `xml:"parameterList"`
				} `xml:"queryByParameter"`
			} `xml:"controlActProcess"`
		} `xml:"PRPA_IN201306UV02"`
	} `xml:"Body"`
}
type PIXv3Response struct {
	XMLName xml.Name `xml:"Envelope"`
	S       string   `xml:"S,attr"`
	Env     string   `xml:"env,attr"`
	Header  struct {
		Action struct {
			Xmlns string `xml:"xmlns,attr"`
			S     string `xml:"S,attr"`
			Env   string `xml:"env,attr"`
		} `xml:"Action"`
		MessageID struct {
			Xmlns string `xml:"xmlns,attr"`
			S     string `xml:"S,attr"`
			Env   string `xml:"env,attr"`
		} `xml:"MessageID"`
		RelatesTo struct {
			Xmlns string `xml:"xmlns,attr"`
			S     string `xml:"S,attr"`
			Env   string `xml:"env,attr"`
		} `xml:"RelatesTo"`
		To struct {
			Xmlns string `xml:"xmlns,attr"`
			S     string `xml:"S,attr"`
			Env   string `xml:"env,attr"`
		} `xml:"To"`
	} `xml:"Header"`
	Body struct {
		PRPAIN201310UV02 struct {
			Xmlns      string `xml:"xmlns,attr"`
			ITSVersion string `xml:"ITSVersion,attr"`
			ID         struct {
				Extension string `xml:"extension,attr"`
				Root      string `xml:"root,attr"`
			} `xml:"id"`
			CreationTime struct {
				Value string `xml:"value,attr"`
			} `xml:"creationTime"`
			VersionCode struct {
				Code string `xml:"code,attr"`
			} `xml:"versionCode"`
			InteractionId struct {
				Extension string `xml:"extension,attr"`
				Root      string `xml:"root,attr"`
			} `xml:"interactionId"`
			ProcessingCode struct {
				Code string `xml:"code,attr"`
			} `xml:"processingCode"`
			ProcessingModeCode struct {
				Code string `xml:"code,attr"`
			} `xml:"processingModeCode"`
			AcceptAckCode struct {
				Code string `xml:"code,attr"`
			} `xml:"acceptAckCode"`
			Receiver struct {
				TypeCode string `xml:"typeCode,attr"`
				Device   struct {
					ClassCode      string `xml:"classCode,attr"`
					DeterminerCode string `xml:"determinerCode,attr"`
					ID             struct {
						AssigningAuthorityName string `xml:"assigningAuthorityName,attr"`
						Root                   string `xml:"root,attr"`
					} `xml:"id"`
					AsAgent struct {
						ClassCode               string `xml:"classCode,attr"`
						RepresentedOrganization struct {
							ClassCode      string `xml:"classCode,attr"`
							DeterminerCode string `xml:"determinerCode,attr"`
							ID             struct {
								AssigningAuthorityName string `xml:"assigningAuthorityName,attr"`
								Root                   string `xml:"root,attr"`
							} `xml:"id"`
						} `xml:"representedOrganization"`
					} `xml:"asAgent"`
				} `xml:"device"`
			} `xml:"receiver"`
			Sender struct {
				TypeCode string `xml:"typeCode,attr"`
				Device   struct {
					ClassCode      string `xml:"classCode,attr"`
					DeterminerCode string `xml:"determinerCode,attr"`
					ID             struct {
						Root string `xml:"root,attr"`
					} `xml:"id"`
					AsAgent struct {
						ClassCode               string `xml:"classCode,attr"`
						RepresentedOrganization struct {
							ClassCode      string `xml:"classCode,attr"`
							DeterminerCode string `xml:"determinerCode,attr"`
							ID             struct {
								Root string `xml:"root,attr"`
							} `xml:"id"`
						} `xml:"representedOrganization"`
					} `xml:"asAgent"`
				} `xml:"device"`
			} `xml:"sender"`
			Acknowledgement struct {
				TypeCode struct {
					Code string `xml:"code,attr"`
				} `xml:"typeCode"`
				TargetMessage struct {
					ID struct {
						Extension string `xml:"extension,attr"`
						Root      string `xml:"root,attr"`
					} `xml:"id"`
				} `xml:"targetMessage"`
			} `xml:"acknowledgement"`
			ControlActProcess struct {
				ClassCode string `xml:"classCode,attr"`
				MoodCode  string `xml:"moodCode,attr"`
				Code      struct {
					Code       string `xml:"code,attr"`
					CodeSystem string `xml:"codeSystem,attr"`
				} `xml:"code"`
				Subject struct {
					TypeCode          string `xml:"typeCode,attr"`
					RegistrationEvent struct {
						ClassCode string `xml:"classCode,attr"`
						MoodCode  string `xml:"moodCode,attr"`
						ID        struct {
							NullFlavor string `xml:"nullFlavor,attr"`
						} `xml:"id"`
						StatusCode struct {
							Code string `xml:"code,attr"`
						} `xml:"statusCode"`
						Subject1 struct {
							TypeCode string `xml:"typeCode,attr"`
							Patient  struct {
								ClassCode string `xml:"classCode,attr"`
								ID        []struct {
									Extension              string `xml:"extension,attr"`
									Root                   string `xml:"root,attr"`
									AssigningAuthorityName string `xml:"assigningAuthorityName,attr"`
								} `xml:"id"`
								StatusCode struct {
									Code string `xml:"code,attr"`
								} `xml:"statusCode"`
								PatientPerson struct {
									ClassCode      string `xml:"classCode,attr"`
									DeterminerCode string `xml:"determinerCode,attr"`
									Name           struct {
										Given  string `xml:"given"`
										Family string `xml:"family"`
									} `xml:"name"`
								} `xml:"patientPerson"`
							} `xml:"patient"`
						} `xml:"subject1"`
						Custodian struct {
							TypeCode       string `xml:"typeCode,attr"`
							AssignedEntity struct {
								ClassCode string `xml:"classCode,attr"`
								ID        struct {
									Root string `xml:"root,attr"`
								} `xml:"id"`
							} `xml:"assignedEntity"`
						} `xml:"custodian"`
					} `xml:"registrationEvent"`
				} `xml:"subject"`
				QueryAck struct {
					QueryId struct {
						Extension string `xml:"extension,attr"`
						Root      string `xml:"root,attr"`
					} `xml:"queryId"`
					StatusCode struct {
						Code string `xml:"code,attr"`
					} `xml:"statusCode"`
					QueryResponseCode struct {
						Code string `xml:"code,attr"`
					} `xml:"queryResponseCode"`
					ResultTotalQuantity struct {
						Value string `xml:"value,attr"`
					} `xml:"resultTotalQuantity"`
					ResultCurrentQuantity struct {
						Value string `xml:"value,attr"`
					} `xml:"resultCurrentQuantity"`
					ResultRemainingQuantity struct {
						Value string `xml:"value,attr"`
					} `xml:"resultRemainingQuantity"`
				} `xml:"queryAck"`
				QueryByParameter struct {
					QueryId struct {
						Extension string `xml:"extension,attr"`
						Root      string `xml:"root,attr"`
					} `xml:"queryId"`
					StatusCode struct {
						Code string `xml:"code,attr"`
					} `xml:"statusCode"`
					ResponsePriorityCode struct {
						Code string `xml:"code,attr"`
					} `xml:"responsePriorityCode"`
					ParameterList struct {
						PatientIdentifier struct {
							Value struct {
								AssigningAuthorityName string `xml:"assigningAuthorityName,attr"`
								Extension              string `xml:"extension,attr"`
								Root                   string `xml:"root,attr"`
							} `xml:"value"`
							SemanticsText string `xml:"semanticsText"`
						} `xml:"patientIdentifier"`
					} `xml:"parameterList"`
				} `xml:"queryByParameter"`
			} `xml:"controlActProcess"`
		} `xml:"PRPA_IN201310UV02"`
	} `xml:"Body"`
}
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
type PIXInterface interface {
	pdq() error
}

var (
	PDQ_V3_Request_Template = "{{define \"pdqv3\"}}<S:Envelope xmlns:S='http://www.w3.org/2003/05/soap-envelope' xmlns:env='http://www.w3.org/2003/05/soap-envelope'><S:Header><To xmlns='http://www.w3.org/2005/08/addressing'>{{.Server_URL}}</To><Action xmlns='http://www.w3.org/2005/08/addressing' S:mustUnderstand='true' xmlns:S='http://www.w3.org/2003/05/soap-envelope'>urn:hl7-org:v3:PRPA_IN201305UV02</Action><ReplyTo xmlns='http://www.w3.org/2005/08/addressing'><Address>http://www.w3.org/2005/08/addressing/anonymous</Address></ReplyTo><FaultTo xmlns='http://www.w3.org/2005/08/addressing'><Address>http://www.w3.org/2005/08/addressing/anonymous</Address></FaultTo><MessageID xmlns='http://www.w3.org/2005/08/addressing'>uuid:{{newuuid}}</MessageID></S:Header><S:Body><PRPA_IN201305UV02 xmlns='urn:hl7-org:v3' ITSVersion='XML_1.0'><id extension='1663079209882' root='1.3.6.1.4.1.21998.2.1.10.15'/><creationTime value='{{simpledatetime}}'/><versionCode code='V3PR1'/><interactionId extension='PRPA_IN201305UV02' root='2.16.840.1.113883.1.6'/><processingCode code='P'/><processingModeCode code='T'/><acceptAckCode code='AL'/><receiver typeCode='RCV'><device classCode='DEV' determinerCode='INSTANCE'><id root='1.3.6.1.4.1.21367.2009.2.2.795'/><asAgent classCode='AGNT'><representedOrganization classCode='ORG' determinerCode='INSTANCE'><id root='1.3.6.1.4.1.21367.2009.2.2.1'/></representedOrganization></asAgent></device></receiver><sender typeCode='SND'><device classCode='DEV' determinerCode='INSTANCE'><id assigningAuthorityName='EHR_TIANI-SPIRIT' root='1.3.6.1.4.1.21367.2011.2.2.7919'/><asAgent classCode='AGNT'><representedOrganization classCode='ORG' determinerCode='INSTANCE'><id assigningAuthorityName='Tiani-Cisco' root='1.3.6.1.4.1.21367.2011.2.7.5572'/></representedOrganization></asAgent></device></sender><controlActProcess classCode='CACT' moodCode='EVN'><code code='PRPA_TE201305UV02' codeSystem='2.16.840.1.113883.1.6'/><queryByParameter><queryId extension='1663079209880' root='1.3.6.1.4.1.21998.2.1.10.15'/><statusCode code='new'/><responseModalityCode code='R'/><responsePriorityCode code='I'/><matchCriterionList/><parameterList><livingSubjectId><value extension='{{.Used_PID}}'/><semanticsText>LivingSubject.id</semanticsText></livingSubjectId></parameterList></queryByParameter></controlActProcess></PRPA_IN201305UV02></S:Body></S:Envelope>{{end}}"
	PIX_V3_Request_Template = "{{define \"pixv3\"}}<S:Envelope xmlns:S='http://www.w3.org/2003/05/soap-envelope' xmlns:env='http://www.w3.org/2003/05/soap-envelope'><S:Header><To xmlns='http://www.w3.org/2005/08/addressing'>{{.Server_URL}}</To><Action xmlns='http://www.w3.org/2005/08/addressing' S:mustUnderstand='true' xmlns:S='http://www.w3.org/2003/05/soap-envelope'>urn:hl7-org:v3:PRPA_IN201309UV02</Action><ReplyTo xmlns='http://www.w3.org/2005/08/addressing'><Address>http://www.w3.org/2005/08/addressing/anonymous</Address></ReplyTo><FaultTo xmlns='http://www.w3.org/2005/08/addressing'><Address>http://www.w3.org/2005/08/addressing/anonymous</Address></FaultTo><MessageID xmlns='http://www.w3.org/2005/08/addressing'>uuid:{{newuuid}}</MessageID></S:Header><S:Body><PRPA_IN201309UV02 xmlns='urn:hl7-org:v3' ITSVersion='XML_1.0'><id extension='1663059665645' root='1.3.6.1.4.1.21998.2.1.10.12'/><creationTime value='{{simpledatetime}}'/><versionCode code='V3PR1'/><interactionId extension='PRPA_IN201309UV02' root='2.16.840.1.113883.1.6'/><processingCode code='P'/><processingModeCode code='T'/><acceptAckCode code='AL'/><receiver typeCode='RCV'><device classCode='DEV' determinerCode='INSTANCE'><id root='1.3.6.1.4.1.21367.2009.2.2.795'/><asAgent classCode='AGNT'><representedOrganization classCode='ORG' determinerCode='INSTANCE'><id root='1.3.6.1.4.1.21367.2009.2.2.1'/></representedOrganization></asAgent></device></receiver><sender typeCode='SND'><device classCode='DEV' determinerCode='INSTANCE'><id assigningAuthorityName='NHS' root='1.3.6.1.4.1.21367.2011.2.2.7919'/><asAgent classCode='AGNT'><representedOrganization classCode='ORG' determinerCode='INSTANCE'><id assigningAuthorityName='ICB' root='1.3.6.1.4.1.21367.2011.2.7.5572'/></representedOrganization></asAgent></device></sender><controlActProcess classCode='CACT' moodCode='EVN'><code code='PRPA_TE201309UV02' codeSystem='2.16.840.1.113883.1.6'/><queryByParameter><queryId extension='1663059665645' root='1.3.6.1.4.1.21998.2.1.10.12'/><statusCode code='new'/><responsePriorityCode code='I'/><parameterList><patientIdentifier><value assigningAuthorityName='{{.Used_PID_OID}}' extension='{{.Used_PID}}' root='{{.Used_PID_OID}}'/><semanticsText>Patient.id</semanticsText></patientIdentifier></parameterList></queryByParameter></controlActProcess></PRPA_IN201309UV02></S:Body></S:Envelope>{{end}}"
)

func PDQ(i PIXInterface) error {
	return i.pdq()
}
func (i *PDQQuery) pdq() error {
	if err := i.setPDQ_ID(); err != nil {
		return err
	}
	return i.getPatient()
}
func (i *PDQQuery) setPDQ_ID() error {
	if i.Server_URL == "" {
		return errors.New("invalid request - pix server url is not set")
	}
	if i.REG_OID == "" {
		return errors.New("invalid request - reg oid is not set")
	}
	if i.Timeout == 0 {
		i.Timeout = 5
	}
	if i.NHS_OID == "" {
		i.NHS_OID = "2.16.840.1.113883.2.1.4.1"
	}
	if i.MRN_ID != "" && i.MRN_OID != "" {
		i.Used_PID = i.MRN_ID
		i.Used_PID_OID = i.MRN_OID
	} else {
		if i.NHS_ID != "" {
			i.Used_PID = i.NHS_ID
			i.Used_PID_OID = i.NHS_OID
		} else {
			if i.REG_ID != "" && i.REG_OID != "" {
				i.Used_PID = i.REG_ID
				i.Used_PID_OID = i.REG_OID
			}
		}
	}
	if i.Used_PID == "" || i.Used_PID_OID == "" {
		return errors.New("invalid request - no suitable id and oid input values found which can be used for pdq query")
	}
	return nil
}
func (i *PDQQuery) getPatient() error {
	var tmplt *template.Template
	var err error
	switch i.Server {
	case cnst.PIXv3:
		if tmplt, err = template.New(cnst.PIXv3).Funcs(util.TemplateFuncMap()).Parse(PIX_V3_Request_Template); err == nil {
			var b bytes.Buffer
			if err = tmplt.Execute(&b, i); err == nil {
				i.Request = b.Bytes()
				if err = i.newTukSOAPRequest(cnst.SOAP_ACTION_PIXV3_Request); err == nil {
					pdqrsp := PIXv3Response{}
					if err = xml.Unmarshal(i.Response, &pdqrsp); err == nil {
						if pdqrsp.Body.PRPAIN201310UV02.Acknowledgement.TypeCode.Code != "AA" {
							return errors.New("acknowledgement code not equal aa, received " + pdqrsp.Body.PRPAIN201310UV02.Acknowledgement.TypeCode.Code)
						}
						i.Count, _ = strconv.Atoi(pdqrsp.Body.PRPAIN201310UV02.ControlActProcess.QueryAck.ResultTotalQuantity.Value)
						pat := PIXPatient{}
						pat.GivenName = pdqrsp.Body.PRPAIN201310UV02.ControlActProcess.Subject.RegistrationEvent.Subject1.Patient.PatientPerson.Name.Given
						pat.FamilyName = pdqrsp.Body.PRPAIN201310UV02.ControlActProcess.Subject.RegistrationEvent.Subject1.Patient.PatientPerson.Name.Family
						i.Patients = append(i.Patients, pat)
					} else {
						log.Println(err.Error())
						return err
					}
				} else {
					log.Println(err.Error())
					return err
				}
			} else {
				log.Println(err.Error())
				return err
			}
		} else {
			log.Println(err.Error())
			return err
		}
	case cnst.PDQv3:
		if tmplt, err = template.New(cnst.PDQv3).Funcs(util.TemplateFuncMap()).Parse(PDQ_V3_Request_Template); err == nil {
			var b bytes.Buffer
			if err = tmplt.Execute(&b, i); err == nil {
				i.Request = b.Bytes()
				if err = i.newTukSOAPRequest(cnst.SOAP_ACTION_PDQV3_Request); err == nil {
					pdqrsp := PDQv3Response{}
					if err = xml.Unmarshal(i.Response, &pdqrsp); err == nil {
						if pdqrsp.Body.PRPAIN201306UV02.Acknowledgement.TypeCode.Code != "AA" {
							return errors.New("acknowledgement code not equal aa, received " + pdqrsp.Body.PRPAIN201306UV02.Acknowledgement.TypeCode.Code)
						}
						i.Count, _ = strconv.Atoi(pdqrsp.Body.PRPAIN201306UV02.ControlActProcess.QueryAck.ResultTotalQuantity.Value)
						pat := PIXPatient{}
						pat.GivenName = pdqrsp.Body.PRPAIN201306UV02.ControlActProcess.Subject.RegistrationEvent.Subject1.Patient.PatientPerson.Name.Given
						pat.FamilyName = pdqrsp.Body.PRPAIN201306UV02.ControlActProcess.Subject.RegistrationEvent.Subject1.Patient.PatientPerson.Name.Family
						pat.BirthDate = pdqrsp.Body.PRPAIN201306UV02.ControlActProcess.Subject.RegistrationEvent.Subject1.Patient.PatientPerson.BirthTime.Value
						pat.Zip = pdqrsp.Body.PRPAIN201306UV02.ControlActProcess.Subject.RegistrationEvent.Subject1.Patient.PatientPerson.Addr.PostalCode
						pat.City = pdqrsp.Body.PRPAIN201306UV02.ControlActProcess.Subject.RegistrationEvent.Subject1.Patient.PatientPerson.Addr.City
						pat.State = pdqrsp.Body.PRPAIN201306UV02.ControlActProcess.Subject.RegistrationEvent.Subject1.Patient.PatientPerson.Addr.State
						pat.Street = pdqrsp.Body.PRPAIN201306UV02.ControlActProcess.Subject.RegistrationEvent.Subject1.Patient.PatientPerson.Addr.StreetAddressLine
						i.Patients = append(i.Patients, pat)
					} else {
						log.Println(err.Error())
						return err
					}
				} else {
					log.Println(err.Error())
					return err
				}
			} else {
				log.Println(err.Error())
				return err
			}
		} else {
			log.Println(err.Error())
			return err
		}
	case cnst.PIXm:
		if err := i.newTukHttpRequest(); err != nil {
			return err
		}
		if strings.Contains(string(i.Response), "Error") {
			return errors.New(string(i.Response))
		}
		pdqrsp := PIXmResponse{}
		if err := json.Unmarshal(i.Response, &pdqrsp); err != nil {
			log.Println("Error unmarshalling i.Response")
			return err
		}
		log.Printf("%v Patient Entries in Response", pdqrsp.Total)
		i.Count = pdqrsp.Total
		if i.Count > 0 {
			for cnt := 0; cnt < len(pdqrsp.Entry); cnt++ {
				rsppat := pdqrsp.Entry[cnt]
				tukpat := PIXPatient{}
				for _, id := range rsppat.Resource.Identifier {
					if id.System == cnst.URN_OID_PREFIX+i.REG_OID {
						tukpat.REGID = id.Value
						tukpat.REGOID = i.REG_OID
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
	}
	return nil
}
func (i *PDQQuery) newTukHttpRequest() error {
	httpReq := tukhttp.PIXmRequest{
		URL:     i.Server_URL,
		PID_OID: i.Used_PID_OID,
		PID:     i.Used_PID,
		Timeout: i.Timeout,
	}
	err := tukhttp.NewRequest(&httpReq)
	i.Request = []byte(httpReq.URL)
	i.Response = httpReq.Response
	i.StatusCode = httpReq.StatusCode
	return err
}
func (i *PDQQuery) newTukSOAPRequest(soapaction string) error {
	httpReq := tukhttp.SOAPRequest{
		URL:        i.Server_URL,
		SOAPAction: soapaction,
		Body:       i.Request,
		Timeout:    i.Timeout,
	}
	err := tukhttp.NewRequest(&httpReq)
	i.Response = httpReq.Response
	i.StatusCode = httpReq.StatusCode
	return err
}
