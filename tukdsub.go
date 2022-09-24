package tukdsub

import (
	"bytes"
	"encoding/xml"
	"errors"
	"log"
	"strings"
	"text/template"

	"github.com/ipthomas/tukcnst"
	"github.com/ipthomas/tukdbint"
	"github.com/ipthomas/tukhttp"
	"github.com/ipthomas/tukpdq"
	"github.com/ipthomas/tukutil"
)

type DSUBEvent struct {
	Action                  string
	TUK_DB_Connection       tukdbint.TukDBConnection
	DSUB_Broker_URL         string
	DSUB_Consumer_URL       string
	DSUB_Topic              string
	DSUB_Expression         string
	DSUB_Ack_Template       string
	DSUB_Subscribe_Template string
	DSUB_Cancel_Template    string
	PDQ_SERVER_URL          string
	PDQ_SERVER_TYPE         string
	REG_OID                 string
	NHS_OID                 string
	Message                 string
	Notify                  DSUBNotifyMessage
	Event                   tukdbint.Event
	Request                 []byte
	Response                string
	BrokerRef               string
	UUID                    string
}

// DSUBSubscribeResponse is an IHE DSUB Subscribe Message compliant struct
type DSUBSubscribeResponse struct {
	XMLName        xml.Name `xml:"Envelope"`
	Text           string   `xml:",chardata"`
	S              string   `xml:"s,attr"`
	A              string   `xml:"a,attr"`
	Xsi            string   `xml:"xsi,attr"`
	Wsnt           string   `xml:"wsnt,attr"`
	SchemaLocation string   `xml:"schemaLocation,attr"`
	Header         struct {
		Text   string `xml:",chardata"`
		Action string `xml:"Action"`
	} `xml:"Header"`
	Body struct {
		Text              string `xml:",chardata"`
		SubscribeResponse struct {
			Text                  string `xml:",chardata"`
			SubscriptionReference struct {
				Text    string `xml:",chardata"`
				Address string `xml:"Address"`
			} `xml:"SubscriptionReference"`
		} `xml:"SubscribeResponse"`
	} `xml:"Body"`
}

// DSUBSubscribe struct provides method dsubSubscribe.NewEvent() to create and send a IHE DSUB Subscribe SOAP message to a IHE DSUB Broker
type DSUBSubscribe struct {
	BrokerURL          string
	ConsumerURL        string
	Topic              string
	Expression         string
	Request            []byte
	BrokerRef          string
	Subscribe_Template string
}

// DSUBCancel struct provides method dsubCancel.NewEvent() to create and send a IHE DSUB Cancel SOAP message to a IHE DSUB Broker
type DSUBCancel struct {
	BrokerURL       string
	BrokerRef       string
	UUID            string
	Request         []byte
	Cancel_Template string
}

// DSUBAck struct provides method dsubAck.NewEvent() to create and send a DSUB Ack SOAP messge to IHE DSUB Broker
type DSUBAck struct {
	BrokerURL    string
	Request      []byte
	Ack_Template string
}

// DSUBNotifyMessage is an IHE DSUB Notify Message compliant struct
type DSUBNotifyMessage struct {
	XMLName             xml.Name `xml:"Notify"`
	Text                string   `xml:",chardata"`
	Xmlns               string   `xml:"xmlns,attr"`
	Xsd                 string   `xml:"xsd,attr"`
	Xsi                 string   `xml:"xsi,attr"`
	NotificationMessage struct {
		Text                  string `xml:",chardata"`
		SubscriptionReference struct {
			Text    string `xml:",chardata"`
			Address struct {
				Text  string `xml:",chardata"`
				Xmlns string `xml:"xmlns,attr"`
			} `xml:"Address"`
		} `xml:"SubscriptionReference"`
		Topic struct {
			Text    string `xml:",chardata"`
			Dialect string `xml:"Dialect,attr"`
		} `xml:"Topic"`
		ProducerReference struct {
			Text    string `xml:",chardata"`
			Address struct {
				Text  string `xml:",chardata"`
				Xmlns string `xml:"xmlns,attr"`
			} `xml:"Address"`
		} `xml:"ProducerReference"`
		Message struct {
			Text                 string `xml:",chardata"`
			SubmitObjectsRequest struct {
				Text               string `xml:",chardata"`
				Lcm                string `xml:"lcm,attr"`
				RegistryObjectList struct {
					Text            string `xml:",chardata"`
					Rim             string `xml:"rim,attr"`
					ExtrinsicObject struct {
						Text       string `xml:",chardata"`
						A          string `xml:"a,attr"`
						ID         string `xml:"id,attr"`
						MimeType   string `xml:"mimeType,attr"`
						ObjectType string `xml:"objectType,attr"`
						Slot       []struct {
							Text      string `xml:",chardata"`
							Name      string `xml:"name,attr"`
							ValueList struct {
								Text  string   `xml:",chardata"`
								Value []string `xml:"Value"`
							} `xml:"ValueList"`
						} `xml:"Slot"`
						Name struct {
							Text            string `xml:",chardata"`
							LocalizedString struct {
								Text  string `xml:",chardata"`
								Value string `xml:"value,attr"`
							} `xml:"LocalizedString"`
						} `xml:"Name"`
						Description    string `xml:"Description"`
						Classification []struct {
							Text                 string `xml:",chardata"`
							ClassificationScheme string `xml:"classificationScheme,attr"`
							ClassifiedObject     string `xml:"classifiedObject,attr"`
							ID                   string `xml:"id,attr"`
							NodeRepresentation   string `xml:"nodeRepresentation,attr"`
							ObjectType           string `xml:"objectType,attr"`
							Slot                 []struct {
								Text      string `xml:",chardata"`
								Name      string `xml:"name,attr"`
								ValueList struct {
									Text  string   `xml:",chardata"`
									Value []string `xml:"Value"`
								} `xml:"ValueList"`
							} `xml:"Slot"`
							Name struct {
								Text            string `xml:",chardata"`
								LocalizedString struct {
									Text  string `xml:",chardata"`
									Value string `xml:"value,attr"`
								} `xml:"LocalizedString"`
							} `xml:"Name"`
						} `xml:"Classification"`
						ExternalIdentifier []struct {
							Text                 string `xml:",chardata"`
							ID                   string `xml:"id,attr"`
							IdentificationScheme string `xml:"identificationScheme,attr"`
							ObjectType           string `xml:"objectType,attr"`
							RegistryObject       string `xml:"registryObject,attr"`
							Value                string `xml:"value,attr"`
							Name                 struct {
								Text            string `xml:",chardata"`
								LocalizedString struct {
									Text  string `xml:",chardata"`
									Value string `xml:"value,attr"`
								} `xml:"LocalizedString"`
							} `xml:"Name"`
						} `xml:"ExternalIdentifier"`
					} `xml:"ExtrinsicObject"`
				} `xml:"RegistryObjectList"`
			} `xml:"SubmitObjectsRequest"`
		} `xml:"Message"`
	} `xml:"NotificationMessage"`
}

type DSUB_Interface interface {
	newEvent() error
}

func NewEvent(i DSUB_Interface) error {
	return i.newEvent()
}

func (i *DSUBEvent) newEvent() error {
	switch i.Action {
	case tukcnst.CREATE:
		return i.newDSUBNotify()
	}
	return errors.New("action not supported")
}

// (i *DSUBEvent) newDSUBNotify creates a DSUBNotifyMessage from the DSUBEvent.Message and creates a dbint.Event from the DSUBNotifyMessage values
// It then checks for Tuk DB Subscriptions matching the brokerref in the populated DSUBEvent.BrokerRef and creates a Tuk DB Event for each subscription
// A DSUB ack response is always sent back to the DSUB broker regardless of success
// If no subscriptions are found a DSUB cancel message is sent to the DSUB Broker
//
//	Example
//		dsubEvent := tukdsub.DSUBEvent{Message: string(`notiftMessage.Body')}
//		dsubEvent.NewEvent()
func (i *DSUBEvent) newDSUBNotify() error {
	log.Printf("Processing DSUB Broker Notfy Message\n%s", i.Message)
	i.Response = i.DSUB_Ack_Template
	if err := i.NewDSUBNotifyMessage(); err == nil {
		if i.Event.BrokerRef == "" {
			log.Println("no subscription ref found in notification message")
			return nil
		}
		log.Printf("Found Subscription Reference %s. Setting Event state from Notify Message", i.Event.BrokerRef)
		i.initEvent()
		if i.Event.XdsPid == "" {
			log.Println("no pid found in notification message")
			return nil
		}
		tukdbint.NewDBEvent(&i.TUK_DB_Connection)
		log.Printf("Checking for TUK Event subscriptions with Broker Ref = %s", i.Event.BrokerRef)
		tukdbSub := tukdbint.Subscription{BrokerRef: i.Event.BrokerRef}
		tukdbSubs := tukdbint.Subscriptions{Action: tukcnst.SELECT}
		tukdbSubs.Subscriptions = append(tukdbSubs.Subscriptions, tukdbSub)
		if err = tukdbint.NewDBEvent(&tukdbSubs); err == nil {
			log.Printf("TUK Event Subscriptions Count : %v", tukdbSubs.Count)
			if tukdbSubs.Count > 0 {
				log.Printf("Obtaining NHS ID. Using %s", i.Event.XdsPid+":"+i.REG_OID)
				pdq := tukpdq.PDQQuery{
					Server:     i.PDQ_SERVER_TYPE,
					REG_ID:     i.Event.XdsPid,
					Server_URL: i.PDQ_SERVER_URL,
					REG_OID:    i.REG_OID,
				}
				if err = tukpdq.PDQ(&pdq); err != nil {
					log.Println(err.Error())
					return nil
				}
				if pdq.Count == 0 {
					log.Println("no patient returned for pid " + i.Event.XdsPid)
					return nil
				}
				if len(pdq.Patients[0].NHSID) != 10 {
					log.Println("no valid nhs id returned in pix query for pid " + i.Event.XdsPid)
					return nil
				}
				for _, dbsub := range tukdbSubs.Subscriptions {
					if dbsub.Id > 0 {
						log.Printf("Creating event for %s %s Subsription for Broker Ref %s", dbsub.Pathway, dbsub.Expression, dbsub.BrokerRef)
						i.Event.Pathway = dbsub.Pathway
						i.Event.Topic = dbsub.Topic
						i.Event.NhsId = pdq.Patients[0].NHSID
						tukevs := tukdbint.Events{Action: "insert"}
						tukevs.Events = append(tukevs.Events, i.Event)
						if err = tukdbint.NewDBEvent(&tukevs); err == nil {
							log.Printf("Created TUK DB Event for Pathway %s Expression %s Broker Ref %s", i.Event.Pathway, i.Event.Expression, i.Event.BrokerRef)
						}
					}
				}
			} else {
				log.Printf("No Subscriptions found with brokerref = %s. Sending Cancel request to Broker", i.Event.BrokerRef)
				dsubCancel := DSUBCancel{BrokerRef: i.Event.BrokerRef, UUID: tukutil.NewUuid()}
				dsubCancel.newEvent()
			}
		}
	}
	return nil
}

// InitDSUBEvent initialise the DSUBEvent struc with values parsed from the DSUBNotifyMessage
func (i *DSUBEvent) initEvent() {
	i.Event.Creationtime = tukutil.Time_Now()
	i.Event.DocName = i.Notify.NotificationMessage.Message.SubmitObjectsRequest.RegistryObjectList.ExtrinsicObject.Name.LocalizedString.Value
	i.Event.ClassCode = tukcnst.NO_VALUE
	i.Event.ConfCode = tukcnst.NO_VALUE
	i.Event.FormatCode = tukcnst.NO_VALUE
	i.Event.FacilityCode = tukcnst.NO_VALUE
	i.Event.PracticeCode = tukcnst.NO_VALUE
	i.Event.Expression = tukcnst.NO_VALUE
	i.Event.XdsPid = tukcnst.NO_VALUE
	i.Event.XdsDocEntryUid = tukcnst.NO_VALUE
	i.Event.RepositoryUniqueId = tukcnst.NO_VALUE
	i.Event.NhsId = tukcnst.NO_VALUE
	i.Event.User = tukcnst.NO_VALUE
	i.Event.Org = tukcnst.NO_VALUE
	i.Event.Role = tukcnst.NO_VALUE
	i.Event.Topic = tukcnst.NO_VALUE
	i.Event.Pathway = tukcnst.NO_VALUE
	i.Event.BrokerRef = i.Notify.NotificationMessage.SubscriptionReference.Address.Text
	i.setRepositoryUniqueId()
	for _, c := range i.Notify.NotificationMessage.Message.SubmitObjectsRequest.RegistryObjectList.ExtrinsicObject.Classification {
		log.Printf("Found Classification Scheme %s", c.ClassificationScheme)
		val := c.Name.LocalizedString.Value
		switch c.ClassificationScheme {
		case tukcnst.URN_CLASS_CODE:
			i.Event.ClassCode = val
		case tukcnst.URN_CONF_CODE:
			i.Event.ConfCode = val
		case tukcnst.URN_FORMAT_CODE:
			i.Event.FormatCode = val
		case tukcnst.URN_FACILITY_CODE:
			i.Event.FacilityCode = val
		case tukcnst.URN_PRACTICE_CODE:
			i.Event.PracticeCode = val
		case tukcnst.URN_TYPE_CODE:
			i.Event.Expression = val
		case tukcnst.URN_AUTHOR:
			for _, s := range c.Slot {
				switch s.Name {
				case tukcnst.AUTHOR_PERSON:
					for _, ap := range s.ValueList.Value {
						i.Event.User = i.Event.User + tukutil.PrettyAuthorPerson(ap) + ","
					}
					i.Event.User = strings.TrimSuffix(i.Event.User, ",")
				case tukcnst.AUTHOR_INSTITUTION:
					for _, ai := range s.ValueList.Value {
						i.Event.Org = i.Event.Org + tukutil.PrettyAuthorInstitution(ai) + ","
					}
					i.Event.Org = strings.TrimSuffix(i.Event.Org, ",")
				}
			}
		default:
			log.Printf("Unknown classication scheme %s. Skipping", c.ClassificationScheme)
		}
	}
	i.Event.Role = i.Event.PracticeCode
	i.setExternalIdentifiers()
	log.Println("Parsed DSUB Notify Message")
}

// NewDSUBNotifyMessage creates an IHE DSUB Notify message struc from the notfy element in the i.Message
func (i *DSUBEvent) NewDSUBNotifyMessage() error {
	dsubNotify := DSUBNotifyMessage{}
	if i.Message == "" {
		return errors.New("message is empty")
	}
	notifyElement := tukutil.GetXMLNodeList(i.Message, tukcnst.DSUB_NOTIFY_ELEMENT)
	if notifyElement == "" {
		return errors.New("unable to locate notify element in received message")
	}
	if err := xml.Unmarshal([]byte(notifyElement), &dsubNotify); err != nil {
		return err
	}
	i.Event = tukdbint.Event{BrokerRef: dsubNotify.NotificationMessage.SubscriptionReference.Address.Text}
	i.Notify = dsubNotify
	return nil
}
func (i *DSUBEvent) setRepositoryUniqueId() {
	log.Println("Searching for Repository Unique ID")
	for _, slot := range i.Notify.NotificationMessage.Message.SubmitObjectsRequest.RegistryObjectList.ExtrinsicObject.Slot {
		if slot.Name == tukcnst.REPOSITORY_UID {
			i.Event.RepositoryUniqueId = slot.ValueList.Value[0]
			return
		}
	}
}
func (i *DSUBEvent) setExternalIdentifiers() {
	for exid := range i.Notify.NotificationMessage.Message.SubmitObjectsRequest.RegistryObjectList.ExtrinsicObject.ExternalIdentifier {
		val := i.Notify.NotificationMessage.Message.SubmitObjectsRequest.RegistryObjectList.ExtrinsicObject.ExternalIdentifier[exid].Value
		ids := i.Notify.NotificationMessage.Message.SubmitObjectsRequest.RegistryObjectList.ExtrinsicObject.ExternalIdentifier[exid].IdentificationScheme
		switch ids {
		case tukcnst.URN_XDS_PID:
			i.Event.XdsPid = strings.Split(val, "^^^")[0]
		case tukcnst.URN_XDS_DOCUID:
			i.Event.XdsDocEntryUid = val
		}
	}
}

// (i *DSUBCancel) NewEvent() creates an IHE DSUB cancel message and sends it to the DSUB broker
func (i *DSUBCancel) newEvent() error {
	tmplt, err := template.New(tukcnst.CANCEL).Funcs(tukutil.TemplateFuncMap()).Parse(i.Cancel_Template)
	if err != nil {
		log.Println(err.Error())
		return err
	}
	var b bytes.Buffer
	err = tmplt.Execute(&b, i)
	if err != nil {
		log.Println(err.Error())
		return err
	}
	i.Request = b.Bytes()
	soapReq := tukhttp.SOAPRequest{
		URL:        i.BrokerURL,
		SOAPAction: tukcnst.SOAP_ACTION_UNSUBSCRIBE_REQUEST,
		Timeout:    1,
		Body:       i.Request,
	}
	return tukhttp.NewRequest(&soapReq)
}

// (i *DSUBAck) NewEvent() creates an IHE DSUB Ack message and sends it to the DSUB broker
func (i *DSUBAck) newEvent() error {
	i.Request = []byte(i.Ack_Template)
	soapReq := tukhttp.SOAPRequest{
		URL:     i.BrokerURL,
		Timeout: 2,
		Body:    i.Request,
	}
	return tukhttp.NewRequest(&soapReq)
}

// (i *DSUBSubscribe) NewEvent() creates an IHE DSUB Subscribe message and sends it to the DSUB broker
func (i *DSUBSubscribe) newEvent() error {
	if tmplt, err := template.New(tukcnst.SUBSCRIBE).Funcs(tukutil.TemplateFuncMap()).Parse(i.Subscribe_Template); err == nil {
		var b bytes.Buffer
		if err := tmplt.Execute(&b, i); err == nil {
			i.Request = b.Bytes()
			soapReq := tukhttp.SOAPRequest{
				URL:        i.BrokerURL,
				SOAPAction: tukcnst.SOAP_ACTION_SUBSCRIBE_REQUEST,
				Body:       i.Request,
			}
			if err := tukhttp.NewRequest(&soapReq); err != nil {
				return err
			}
			subrsp := DSUBSubscribeResponse{}
			if err := xml.Unmarshal(soapReq.Response, &subrsp); err == nil {
				i.BrokerRef = subrsp.Body.SubscribeResponse.SubscriptionReference.Address
				log.Printf("Broker Response. Broker Ref :  %s", subrsp.Body.SubscribeResponse.SubscriptionReference.Address)
			} else {
				return err
			}
		} else {
			return err
		}
	} else {
		return err
	}
	return nil
}
