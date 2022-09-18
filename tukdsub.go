package tukdsub

import (
	"C"
	"bytes"
	"encoding/xml"
	"errors"
	"log"
	"strings"
	"text/template"

	"github.com/ipthomas/tukcnst"

	dbint "github.com/ipthomas/tukdbint"
	"github.com/ipthomas/tukhttp"
	pixm "github.com/ipthomas/tukpdq"
	util "github.com/ipthomas/tukutil"
)

var (
	REG_OID                 = ""
	CS                      = make(map[string]string)
	DSUB_ACK_TEMPLATE       = "<SOAP-ENV:Envelope xmlns:SOAP-ENV='http://www.w3.org/2003/05/soap-envelope' xmlns:s='http://www.w3.org/2001/XMLSchema' xmlns:xsi='http://www.w3.org/2001/XMLSchema-instance'><SOAP-ENV:Body/></SOAP-ENV:Envelope>"
	DSUB_SUBSCRIBE_TEMPLATE = "{{define \"subscribe\"}}<SOAP-ENV:Envelope xmlns:SOAP-ENV='http://www.w3.org/2003/05/soap-envelope' xmlns:xsi='http://www.w3.org/2001/XMLSchema-instance' xmlns:s='http://www.w3.org/2001/XMLSchema' xmlns:wsa='http://www.w3.org/2005/08/addressing'><SOAP-ENV:Header><wsa:Action SOAP-ENV:mustUnderstand='true'>http://docs.oasis-open.org/wsn/bw-2/NotificationProducer/SubscribeRequest</wsa:Action><wsa:MessageID>urn:uuid:{{newuuid}}</wsa:MessageID><wsa:ReplyTo SOAP-ENV:mustUnderstand='true'><wsa:Address>http://www.w3.org/2005/08/addressing/anonymous</wsa:Address></wsa:ReplyTo><wsa:To>{{.BrokerUrl}}</wsa:To></SOAP-ENV:Header><SOAP-ENV:Body><wsnt:Subscribe xmlns:wsnt='http://docs.oasis-open.org/wsn/b-2' xmlns:a='http://www.w3.org/2005/08/addressing' xmlns:rim='urn:oasis:names:tc:ebxml-regrep:xsd:rim:3.0' xmlns:wsa='http://www.w3.org/2005/08/addressing'><wsnt:ConsumerReference><wsa:Address>{{.ConsumerUrl}}</wsa:Address></wsnt:ConsumerReference><wsnt:Filter><wsnt:TopicExpression Dialect='http://docs.oasis-open.org/wsn/t-1/TopicExpression/Simple'>ihe:FullDocumentEntry</wsnt:TopicExpression><rim:AdhocQuery id='urn:uuid:742790e0-aba6-43d6-9f1f-e43ed9790b79'><rim:Slot name='{{.Topic}}'><rim:ValueList><rim:Value>('{{.Expression}}')</rim:Value></rim:ValueList></rim:Slot></rim:AdhocQuery></wsnt:Filter></wsnt:Subscribe></SOAP-ENV:Body></SOAP-ENV:Envelope>{{end}}"
	DSUB_CANCEL_TEMPLATE    = "{{define \"cancel\"}}<soap:Envelope xmlns:soap='http://www.w3.org/2003/05/soap-envelope'><soap:Header><Action xmlns='http://www.w3.org/2005/08/addressing' soap:mustUnderstand='true'>http://docs.oasis-open.org/wsn/bw-2/SubscriptionManager/UnsubscribeRequest</Action><MessageID xmlns='http://www.w3.org/2005/08/addressing' soap:mustUnderstand='true'>urn:uuid:{{.UUID}}</MessageID><To xmlns='http://www.w3.org/2005/08/addressing' soap:mustUnderstand='true'>{{.BrokerRef}}</To><ReplyTo xmlns='http://www.w3.org/2005/08/addressing' soap:mustUnderstand='true'><Address>http://www.w3.org/2005/08/addressing/anonymous</Address></ReplyTo></soap:Header><soap:Body><Unsubscribe xmlns='http://docs.oasis-open.org/wsn/b-2' xmlns:ns2='http://www.w3.org/2005/08/addressing' xmlns:ns3='http://docs.oasis-open.org/wsrf/bf-2' xmlns:ns4='urn:oasis:names:tc:ebxml-regrep:xsd:rim:3.0' xmlns:ns5='urn:oasis:names:tc:ebxml-regrep:xsd:rs:3.0' xmlns:ns6='urn:oasis:names:tc:ebxml-regrep:xsd:lcm:3.0' xmlns:ns7='http://docs.oasis-open.org/wsn/t-1' xmlns:ns8='http://docs.oasis-open.org/wsrf/r-2'/></soap:Body></soap:Envelope>{{end}}"
)

type DSUBEvent struct {
	Message string            `json:"message"`
	Notify  DSUBNotifyMessage `json:"notify"`
	Event   dbint.Event       `json:"event"`
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
	BrokerUrl   string
	ConsumerUrl string
	Topic       string
	Expression  string
	Request     []byte
	BrokerRef   string
}

// DSUBCancel struct provides method dsubCancel.NewEvent() to create and send a IHE DSUB Cancel SOAP message to a IHE DSUB Broker
type DSUBCancel struct {
	BrokerRef string
	UUID      string
	Request   []byte
}

// DSUBAck struct provides method dsubAck.NewEvent() to create and send a DSUB Ack SOAP messge to IHE DSUB Broker
type DSUBAck struct {
	Request []byte
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
type CodeSystem struct {
	CS_Map map[string]string
}
type DSUB_Template struct {
	Act       string
	Templates map[string]string
}
type DSUB_Interface interface {
	newEvent() error
}

func SetRegOID(regoid string) {
	REG_OID = regoid
}
func NewDsubEvent(i DSUB_Interface) error {
	return i.newEvent()
}

// (i *CodeSystem) newEvent() sets the DSUB package codesystem which is merged with the tukutil.Codesystem. If a CS entry exists, the provided CS entry replaces the existing cs entry
func (i *CodeSystem) newEvent() error {
	for k, v := range i.CS_Map {
		CS[k] = v
	}
	if regoid, ok := CS["REG_OID"]; ok {
		REG_OID = regoid

	}
	return nil
}

// (i *DSUB_Template) newEvent()
//
//	i.Act must equal 'select' or 'update'
//		'select' sets the Ack,Subscribe and Cancel i.Templates[string]string map values
//		'update' overides the default Ack,Subscribe and Cancel templates with the values in the i.Templates[string]string map.
func (i *DSUB_Template) newEvent() error {
	switch i.Act {
	case tukcnst.SELECT:
		i.Templates = make(map[string]string)
		i.Templates[tukcnst.DSUB_ACK_TEMPLATE] = DSUB_ACK_TEMPLATE
		i.Templates[tukcnst.DSUB_SUBSCRIBE_TEMPLATE] = DSUB_SUBSCRIBE_TEMPLATE
		i.Templates[tukcnst.DSUB_CANCEL_TEMPLATE] = DSUB_CANCEL_TEMPLATE
	case tukcnst.UPDATE:
		DSUB_ACK_TEMPLATE = i.Templates[tukcnst.DSUB_ACK_TEMPLATE]
		DSUB_SUBSCRIBE_TEMPLATE = i.Templates[tukcnst.DSUB_SUBSCRIBE_TEMPLATE]
		DSUB_CANCEL_TEMPLATE = i.Templates[tukcnst.DSUB_CANCEL_TEMPLATE]
	}
	return nil
}

// (i *DSUBEvent) NewEvent creates a DSUBNotifyMessage from the DSUBEvent.Message and creates a dbint.Event from the DSUBNotifyMessage values
// It then checks for Tuk DB Subscriptions matching the brokerref in the populated DSUBEvent.BrokerRef and creates a Tuk DB Event for each subscription
// A DSUB ack response is always sent back to the DSUB broker regardless of success
// If no subscriptions are found a DSUB cancel message is sent to the DSUB Broker
//
//	Example
//		dsubEvent := tukdsub.DSUBEvent{Message: string(`notiftMessage.Body')}
//		dsubEvent.NewEvent()
func (i *DSUBEvent) newEvent() error {
	log.Printf("Processing DSUB Broker Notfy Message\n%s", i.Message)
	if err := i.newDSUBNotifyMessage(); err == nil {
		if i.Event.BrokerRef == "" {
			return errors.New("no subscription ref found in notification message")
		}
		log.Printf("Found Subscription Reference %s. Setting Event state from Notify Message", i.Event.BrokerRef)
		i.initEvent()
		if i.Event.XdsPid == "" {
			return errors.New("no pid found in notification message")
		}
		log.Printf("Checking for TUK Event subscriptions with Broker Ref = %s", i.Event.BrokerRef)
		tukdbSub := dbint.Subscription{BrokerRef: i.Event.BrokerRef}
		tukdbSubs := dbint.Subscriptions{Action: tukcnst.SELECT}
		tukdbSubs.Subscriptions = append(tukdbSubs.Subscriptions, tukdbSub)
		if err = dbint.NewDBEvent(&tukdbSubs); err == nil {
			log.Printf("TUK Event Subscriptions Count : %v", tukdbSubs.Count)
			if tukdbSubs.Count > 0 {
				log.Printf("Obtaining NHS ID. Using %s", i.Event.XdsPid+":"+REG_OID)
				pdq := pixm.PDQQuery{
					Server:     tukcnst.PIXv3,
					REG_ID:     i.Event.XdsPid,
					Server_URL: getCodeSystemVal(tukcnst.PIX_URL),
					REG_OID:    REG_OID,
				}
				if err = pixm.PDQ(&pdq); err != nil {
					return err
				}
				if pdq.Count == 0 {
					return errors.New("no patient returned for pid " + i.Event.XdsPid)
				}
				if len(pdq.Patients[0].NHSID) != 10 {
					return errors.New("no valid nhs id returned in pix query for pid " + i.Event.XdsPid)
				}
				for _, dbsub := range tukdbSubs.Subscriptions {
					if dbsub.Id > 0 {
						log.Printf("Creating event for %s %s Subsription for Broker Ref %s", dbsub.Pathway, dbsub.Expression, dbsub.BrokerRef)
						i.Event.Pathway = dbsub.Pathway
						i.Event.Topic = dbsub.Topic
						i.Event.NhsId = pdq.Patients[0].NHSID
						tukevs := dbint.Events{Action: "insert"}
						tukevs.Events = append(tukevs.Events, i.Event)
						if err = dbint.NewDBEvent(&tukevs); err == nil {
							log.Printf("Created TUK DB Event for Pathway %s Expression %s Broker Ref %s", i.Event.Pathway, i.Event.Expression, i.Event.BrokerRef)
						}
					}
				}
			} else {
				log.Printf("No Subscriptions found with brokerref = %s. Sending Cancel request to Broker", i.Event.BrokerRef)
				dsubCancel := DSUBCancel{BrokerRef: i.Event.BrokerRef, UUID: util.NewUuid()}
				dsubCancel.newEvent()
			}
		}
	}
	return nil
}

// InitDSUBEvent initialise the DSUBEvent struc with values parsed from the DSUBNotifyMessage
func (i *DSUBEvent) initEvent() {
	i.Event.Creationtime = util.Time_Now()
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
						i.Event.User = i.Event.User + util.PrettyAuthorPerson(ap) + ","
					}
					i.Event.User = strings.TrimSuffix(i.Event.User, ",")
				case tukcnst.AUTHOR_INSTITUTION:
					for _, ai := range s.ValueList.Value {
						i.Event.Org = i.Event.Org + util.PrettyAuthorInstitution(ai) + ","
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
	util.Log(i)
}

// NewDSUBNotifyMessage creates an IHE DSUB Notify message struc from the notfy element in the DSUB message
func (i *DSUBEvent) newDSUBNotifyMessage() error {
	dsubNotify := DSUBNotifyMessage{}
	if i.Message == "" {
		return errors.New("message is empty")
	}
	notifyElement := util.GetXMLNodeList(i.Message, tukcnst.DSUB_NOTIFY_ELEMENT)
	if notifyElement == "" {
		return errors.New("unable to locate notify element in received message")
	}
	if err := xml.Unmarshal([]byte(notifyElement), &dsubNotify); err != nil {
		return err
	}
	i.Event = dbint.Event{BrokerRef: dsubNotify.NotificationMessage.SubscriptionReference.Address.Text}
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
	tmplt, err := template.New(tukcnst.CANCEL).Funcs(util.TemplateFuncMap()).Parse(DSUB_CANCEL_TEMPLATE)
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
		URL:        util.GetCodeSystemVal(tukcnst.DSUB_BROKER_URL),
		SOAPAction: tukcnst.SOAP_ACTION_UNSUBSCRIBE_REQUEST,
		Timeout:    1,
		Body:       i.Request,
	}
	return tukhttp.NewRequest(&soapReq)
}

// (i *DSUBAck) NewEvent() creates an IHE DSUB Ack message and sends it to the DSUB broker
func (i *DSUBAck) newEvent() error {
	i.Request = []byte(DSUB_ACK_TEMPLATE)
	soapReq := tukhttp.SOAPRequest{
		URL:     util.GetCodeSystemVal(tukcnst.DSUB_BROKER_URL),
		Timeout: 2,
		Body:    i.Request,
	}
	return tukhttp.NewRequest(&soapReq)
}

// (i *DSUBSubscribe) NewEvent() creates an IHE DSUB Subscribe message and sends it to the DSUB broker
func (i *DSUBSubscribe) newEvent() error {
	if tmplt, err := template.New(tukcnst.SUBSCRIBE).Funcs(util.TemplateFuncMap()).Parse(DSUB_SUBSCRIBE_TEMPLATE); err == nil {
		var b bytes.Buffer
		if err := tmplt.Execute(&b, i); err == nil {
			i.Request = b.Bytes()
			soapReq := tukhttp.SOAPRequest{
				URL:        util.GetCodeSystemVal(tukcnst.DSUB_BROKER_URL),
				SOAPAction: tukcnst.SOAP_ACTION_SUBSCRIBE_REQUEST,
				Timeout:    2,
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

// GetTemplate returns a string of the template for the specified string input `name` valid input names are
//
//	tukcnst..DSUB_ACK_TEMPLATE
//	tukcnst..DSUB_SUBSCRIBE_TEMPLATE
//	tukcnst..DSUB_CANCEL_TEMPLATE
func getTemplate(name string) string {
	switch name {
	case tukcnst.DSUB_ACK_TEMPLATE:
		return DSUB_ACK_TEMPLATE
	case tukcnst.DSUB_SUBSCRIBE_TEMPLATE:
		return DSUB_SUBSCRIBE_TEMPLATE
	case tukcnst.DSUB_CANCEL_TEMPLATE:
		return DSUB_CANCEL_TEMPLATE
	}
	return "Invalid name"
}

// setTemplate overirdes the default template for DSUB messages. The string input `name` valid values are
//
//	tukcnst..DSUB_ACK_TEMPLATE
//	tukcnst..DSUB_SUBSCRIBE_TEMPLATE
//	tukcnst..DSUB_CANCEL_TEMPLATE
//
// The input string dsubTemplate is a 'url safe' string of the template
func setTemplate(name string, dsubTemplate string) {
	switch name {
	case tukcnst.DSUB_ACK_TEMPLATE:
		DSUB_ACK_TEMPLATE = dsubTemplate
	case tukcnst.DSUB_SUBSCRIBE_TEMPLATE:
		DSUB_SUBSCRIBE_TEMPLATE = dsubTemplate
	case tukcnst.DSUB_CANCEL_TEMPLATE:
		DSUB_CANCEL_TEMPLATE = dsubTemplate
	}
}

// GetCodeSystemVal returns the value associated with the input string value. The codesystem is initialised from the tukutil.Codesystem and the dsub codesystem (if set by SetCodesystem(cs) )
func getCodeSystemVal(key string) string {
	if val, ok := CS[key]; ok {
		return val
	}
	return key
}

// functions to support the build of C header / output files

// CGOGetCodeSystemVal supports the build of a C Header / Output for go method GetCodeSystemVal(key)
//
//export CGOGetCodeSystemVal
func CGOGetCodeSystemVal(key string) string {
	return getCodeSystemVal(key)
}

// CGOSetCodeSystemVal supports the build of a C sharp Header / Output for go method CodeSystem(key) = val
//
//export CGOSetCodeSystemVal
func CGOSetCodeSystemVal(key string, val string) {
	CS[key] = val
}

// CGOSetTemplate supports the build of a C sharp Header / Output for go method SetTemplate(name,tmplt)
//
//export CGOSetTemplate
func CGOSetTemplate(name string, tmplt string) {
	setTemplate(name, tmplt)
}

// CGOGetTemplate supports the build of a C sharp Header / Output for go method GetTemplate(name)
//
//export CGOGetTemplate
func CGOGetTemplate(name string) {
	getTemplate(name)
}

// CGONewDSUBCancel supports the build of a C sharp Header / Output for go method DSUBCancel.NewEvent()
//
//export CGONewDSUBCancel
func CGONewDSUBCancel(brokerref string, uuid string) {
	dsubcancel := DSUBCancel{
		BrokerRef: brokerref,
		UUID:      uuid,
	}
	dsubcancel.newEvent()
}

//export CGONewDSUBAck
func CGONewDSUBAck() {
	dsubAck := DSUBAck{}
	dsubAck.newEvent()
}

//export CGONewSubscription
func CGONewSubscription(brokerurl string, consumerurl string, topic string, expression string) string {
	sub := DSUBSubscribe{
		BrokerUrl:   brokerurl,
		ConsumerUrl: consumerurl,
		Topic:       topic,
		Expression:  expression,
	}
	sub.newEvent()
	return sub.BrokerRef
}
