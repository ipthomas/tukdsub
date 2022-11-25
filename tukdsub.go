package tukdsub

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"log"
	"strings"
	"text/template"

	"github.com/ipthomas/tukcnst"
	"github.com/ipthomas/tukdbint"
	"github.com/ipthomas/tukhttp"
	"github.com/ipthomas/tukpdq"
	"github.com/ipthomas/tukutil"
)

var DebugMode = false

// DSUBEvent implements NewEvent(i DSUB_Interface) error
type DSUBEvent struct {
	Action          string
	BrokerURL       string
	ConsumerURL     string
	PDQ_SERVER_URL  string
	PDQ_SERVER_TYPE string
	REG_OID         string
	NHS_OID         string
	EventMessage    string
	Expressions     []string
	Pathway         string
	RowID           int64
	DBConnection    tukdbint.TukDBConnection
	Subs            tukdbint.Subscriptions
	Request         []byte
	Response        []byte
	BrokerRef       string
	UUID            string
	Notify          DSUBNotifyMessage
	Event           tukdbint.Event
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

func New_Transaction(i DSUB_Interface) error {
	return i.newEvent()
}

func (i *DSUBEvent) newEvent() error {
	var err error
	switch i.Action {
	case tukcnst.SELECT:
		l(fmt.Sprintf("Processing Select Subscriptions for Pathway %s", i.Pathway), true)
		err = i.selectSubscriptions()
	case tukcnst.CREATE:
		l(fmt.Sprintf("Processing Create Subscriptions"), true)
		err = i.createSubscriptions()
	case tukcnst.CANCEL:
		l(fmt.Sprintf("Processing Cancel Subscriptions"), true)
		err = i.cancelSubscriptions()
	default:
		l(fmt.Sprintf("Processing Broker Notify Message\n%s", i.EventMessage), true)
		i.processBrokerEventMessage()
	}
	return err
}

// creates a DSUBNotifyMessage from the EventMessage and populates a new TUKEvent with the DSUBNotifyMessage values
// It then checks for TukDB Subscriptions matching the brokerref and creates a TUKEvent for each subscription
// A DSUB ack response is always returned regardless of success
func (i *DSUBEvent) processBrokerEventMessage() {
	i.Response = []byte(tukcnst.GO_TEMPLATE_DSUB_ACK)
	if err := i.newDSUBNotifyMessage(); err == nil {
		if i.Event.BrokerRef == "" {
			l("no subscription ref found in notification message", true)
			return
		}
		l(fmt.Sprintf("Found Subscription Reference %s Setting Event state from Notify Message", i.Event.BrokerRef), true)
		i.initTUKEvent()
		if i.Event.XdsPid == "" {
			l("no pid found in notification message", true)
			return
		}
		l(fmt.Sprintf("Checking for TUK Event subscriptions with Broker Ref = %s", i.Event.BrokerRef), true)
		tukdbSubs := tukdbint.Subscriptions{Action: tukcnst.SELECT}
		tukdbSub := tukdbint.Subscription{BrokerRef: i.Event.BrokerRef}
		tukdbSubs.Subscriptions = append(tukdbSubs.Subscriptions, tukdbSub)
		if err = tukdbint.NewDBEvent(&tukdbSubs); err == nil {
			l(fmt.Sprintf("TUK Event Subscriptions Count : %v", tukdbSubs.Count), true)
			if tukdbSubs.Count > 0 {
				l(fmt.Sprintf("Obtaining NHS ID. Using %s", i.Event.XdsPid+":"+i.REG_OID), true)
				pdq := tukpdq.PDQQuery{
					Server_Mode: i.PDQ_SERVER_TYPE,
					REG_ID:      i.Event.XdsPid,
					Server_URL:  i.PDQ_SERVER_URL,
					REG_OID:     i.REG_OID,
				}
				if err = tukpdq.New_Transaction(&pdq); err != nil {
					l(err.Error(), false)
					return
				}
				if pdq.Count == 0 {
					l("no patient returned for pid "+i.Event.XdsPid, true)
					return
				}

				for _, dbsub := range tukdbSubs.Subscriptions {
					if dbsub.Id > 0 {
						l(fmt.Sprintf("Creating event for %s %s Subsription for Broker Ref %s", dbsub.Pathway, dbsub.Expression, dbsub.BrokerRef), true)
						i.Event.Pathway = dbsub.Pathway
						i.Event.Topic = dbsub.Topic
						i.Event.NhsId = pdq.NHS_ID
						tukevs := tukdbint.Events{Action: tukcnst.INSERT}
						tukevs.Events = append(tukevs.Events, i.Event)
						if err = tukdbint.NewDBEvent(&tukevs); err == nil {
							l(fmt.Sprintf("Created TUK DB Event for Pathway %s Expression %s Broker Ref %s", i.Event.Pathway, i.Event.Expression, i.Event.BrokerRef), true)
						}
					}
				}
			} else {
				l(fmt.Sprintf("No Subscriptions found with brokerref = %s. Sending Cancel request to Broker", i.Event.BrokerRef), true)
				i.cancelSubscriptions()
			}
		}
	} else {
		l(err.Error(), false)
	}

}

// InitDSUBEvent initialise the DSUBEvent struc with values parsed from the DSUBNotifyMessage
func (i *DSUBEvent) initTUKEvent() {
	i.Event.Creationtime = tukutil.Time_Now()
	i.Event.DocName = i.Notify.NotificationMessage.Message.SubmitObjectsRequest.RegistryObjectList.ExtrinsicObject.Name.LocalizedString.Value
	i.Event.BrokerRef = i.Notify.NotificationMessage.SubscriptionReference.Address.Text
	i.setRepositoryUniqueId()
	for _, c := range i.Notify.NotificationMessage.Message.SubmitObjectsRequest.RegistryObjectList.ExtrinsicObject.Classification {
		val := c.Name.LocalizedString.Value
		switch c.ClassificationScheme {
		case tukcnst.URN_CLASS_CODE:
			i.Event.ClassCode = val
			l(fmt.Sprintf("Set Class Code %s", val), true)
		case tukcnst.URN_CONF_CODE:
			i.Event.ConfCode = val
			l(fmt.Sprintf("Set Conf Code %s", val), true)
		case tukcnst.URN_FORMAT_CODE:
			i.Event.FormatCode = val
			l(fmt.Sprintf("Set Format Code %s", val), true)
		case tukcnst.URN_FACILITY_CODE:
			i.Event.FacilityCode = val
			l(fmt.Sprintf("Set Facility Code %s", val), true)
		case tukcnst.URN_PRACTICE_CODE:
			i.Event.PracticeCode = val
			l(fmt.Sprintf("Set Practice Code %s", val), true)
		case tukcnst.URN_TYPE_CODE:
			i.Event.Expression = val
			l(fmt.Sprintf("Set Type Code %s", val), true)
		case tukcnst.URN_AUTHOR:
			for _, slot := range c.Slot {
				switch slot.Name {
				case tukcnst.AUTHOR_PERSON:
					if i.Event.User == "" {
						user := strings.ReplaceAll(slot.ValueList.Value[0], "^", " ")
						user = strings.TrimSpace(user)
						i.Event.User = i.Event.User + user + " "
						l(fmt.Sprintf("Set User %s", i.Event.User), true)
					}
				case tukcnst.AUTHOR_INSTITUTION:
					if i.Event.Org == "" {
						i.Event.Org = strings.TrimSuffix(slot.ValueList.Value[0], "^^^")
						l(fmt.Sprintf("Set Org %s", i.Event.Org), true)
					}
				case tukcnst.AUTHOR_ROLE:
					if i.Event.Role == "" {
						i.Event.Role = strings.TrimSuffix(slot.ValueList.Value[0], "^^^")
						l(fmt.Sprintf("Set Role %s", i.Event.Role), true)
					}
				case tukcnst.AUTHOR_SPECIALITY:
					for _, value := range slot.ValueList.Value {
						if !strings.HasPrefix(i.Event.Speciality, value+". ") {
							i.Event.Speciality = i.Event.Speciality + strings.TrimSuffix(value, "^^^") + ". "
						}
					}
					log.Printf("Set Speciality %s", i.Event.Speciality)
					i.Event.Speciality = strings.TrimSpace(i.Event.Speciality)
				}
			}
		case tukcnst.URN_EVENT_LIST:
			i.Event.Comments = i.Event.Comments + val + ". "
		default:
			l(fmt.Sprintf("Unknown classication scheme %s. Skipping", c.ClassificationScheme), true)
		}
	}
	i.Event.Comments = strings.TrimSpace(i.Event.Comments)
	l(fmt.Sprintf("Added Event Codes to Notes %s", i.Event.Comments), true)
	i.setExternalIdentifiers()
	l("Parsed DSUB Notify Message", true)
}

// NewDSUBNotifyMessage creates an IHE DSUB Notify message struc from the notfy element in the i.Message
func (i *DSUBEvent) newDSUBNotifyMessage() error {
	dsubNotify := DSUBNotifyMessage{}
	if i.EventMessage == "" {
		return errors.New("message is empty")
	}
	if err := xml.Unmarshal([]byte(i.EventMessage), &dsubNotify); err != nil {
		return err
	}
	i.Event = tukdbint.Event{BrokerRef: dsubNotify.NotificationMessage.SubscriptionReference.Address.Text}
	i.Notify = dsubNotify
	return nil
}
func (i *DSUBEvent) newDSUBCancelMessage() {
	tmplt, err := template.New(tukcnst.CANCEL).Funcs(tukutil.TemplateFuncMap()).Parse(tukcnst.GO_TEMPLATE_DSUB_CANCEL)
	if err != nil {
		l(err.Error(), false)
		return
	}
	var b bytes.Buffer
	err = tmplt.Execute(&b, i.BrokerRef)
	if err != nil {
		l(err.Error(), false)
		return
	}
	soapReq := tukhttp.SOAPRequest{
		URL:        i.BrokerURL,
		SOAPAction: tukcnst.SOAP_ACTION_UNSUBSCRIBE_REQUEST,
		Body:       b.Bytes(),
		Timeout:    2,
	}
	log.Printf("Sending Cancel Request to DSUB Broker %s", i.BrokerURL)
	err = tukhttp.NewRequest(&soapReq)
	if err != nil {
		l(err.Error(), false)
	}
}
func (i *DSUBEvent) setRepositoryUniqueId() {
	for _, slot := range i.Notify.NotificationMessage.Message.SubmitObjectsRequest.RegistryObjectList.ExtrinsicObject.Slot {
		if slot.Name == tukcnst.REPOSITORY_UID {
			i.Event.RepositoryUniqueId = slot.ValueList.Value[0]
			l(fmt.Sprintf("Set Repository UID %s", i.Event.RepositoryUniqueId), true)
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
			l(fmt.Sprintf("Set XDS PID %s", i.Event.XdsPid), true)
		case tukcnst.URN_XDS_DOCUID:
			i.Event.XdsDocEntryUid = val
			l(fmt.Sprintf("Set XDS Doc UID %s", i.Event.XdsDocEntryUid), true)
		}
	}
}

func (i *DSUBEvent) cancelSubscriptions() error {
	if i.Pathway == "" && i.RowID == 0 {
		log.Println("pathway or rowid not set. Sending Cancel subscription message to Broker")
		i.newDSUBCancelMessage()
		return nil
	}
	i.Subs = tukdbint.Subscriptions{Action: tukcnst.DELETE}
	delsub := tukdbint.Subscription{}
	if i.Pathway != "" {
		delsub.Pathway = i.Pathway
	} else {
		delsub.Id = i.RowID
	}
	i.Subs.Subscriptions = append(i.Subs.Subscriptions, delsub)
	return tukdbint.NewDBEvent(&i.Subs)
}

// (i *DSUBSubscribe) NewEvent() creates an IHE DSUB Subscribe request to an IHE DSUB broker. The Subscription reference is persisted in the subscriptions table.
func (i *DSUBEvent) createSubscriptions() error {
	type subreq struct {
		BrokerURL   string
		ConsumerURL string
		Topic       string
		Expression  string
	}
	i.Subs = tukdbint.Subscriptions{Action: i.Action}
	if len(i.Expressions) > 0 {
		for _, expression := range i.Expressions {
			l(fmt.Sprintf("Checking for existing Subscriptions for Expression %s", expression), true)
			expressionSubs := tukdbint.Subscriptions{Action: tukcnst.SELECT}
			expressionSub := tukdbint.Subscription{Pathway: i.Pathway, Topic: tukcnst.DSUB_TOPIC_TYPE_CODE, Expression: expression}
			expressionSubs.Subscriptions = append(expressionSubs.Subscriptions, expressionSub)
			if err := tukdbint.NewDBEvent(&expressionSubs); err != nil {
				l(err.Error(), false)
				return err
			}
			if expressionSubs.Count == 1 {
				l(fmt.Sprintf("Found %v Subscription for Pathway %s Topic %s Expression %s", expressionSubs.Count, i.Pathway, tukcnst.DSUB_TOPIC_TYPE_CODE, expression), true)
				i.Subs.Subscriptions = append(i.Subs.Subscriptions, expressionSubs.Subscriptions[1])
				i.Subs.Count = i.Subs.Count + 1
				if i.Subs.LastInsertId < int64(expressionSubs.Subscriptions[1].Id) {
					i.Subs.LastInsertId = int64(expressionSubs.Subscriptions[1].Id)
				}
			} else {
				l(fmt.Sprintf("No Subscription for Pathway %s Topic %s Expression %s found", i.Pathway, tukcnst.DSUB_TOPIC_TYPE_CODE, expression), true)
				brokerSub := subreq{
					BrokerURL:   i.BrokerURL,
					ConsumerURL: i.ConsumerURL,
					Topic:       tukcnst.DSUB_TOPIC_TYPE_CODE,
					Expression:  expression,
				}
				tmplt, err := template.New(tukcnst.SUBSCRIBE).Funcs(tukutil.TemplateFuncMap()).Parse(tukcnst.GO_TEMPLATE_DSUB_SUBSCRIBE)
				if err != nil {
					l(err.Error(), false)
					return err
				}
				var b bytes.Buffer
				err = tmplt.Execute(&b, brokerSub)
				if err != nil {
					l(err.Error(), false)
					return err
				}
				soapReq := tukhttp.SOAPRequest{
					URL:        i.BrokerURL,
					SOAPAction: tukcnst.SOAP_ACTION_SUBSCRIBE_REQUEST,
					Body:       b.Bytes(),
				}
				l(fmt.Sprintf("Sending Subscribe Request to DSUB Broker %s", i.BrokerURL), true)
				err = tukhttp.NewRequest(&soapReq)
				if err != nil {
					l(err.Error(), false)
					return err
				}
				subrsp := DSUBSubscribeResponse{}
				err = xml.Unmarshal(soapReq.Response, &subrsp)
				if err != nil {
					l(err.Error(), false)
					return err
				}
				newSub := tukdbint.Subscription{}
				newSub.BrokerRef = subrsp.Body.SubscribeResponse.SubscriptionReference.Address
				l(fmt.Sprintf("Set Broker Ref =  %s", newSub.BrokerRef), true)
				if newSub.BrokerRef != "" {
					newSub.Pathway = i.Pathway
					newSub.Expression = brokerSub.Expression
					newSub.Topic = tukcnst.DSUB_TOPIC_TYPE_CODE
					newsubs := tukdbint.Subscriptions{Action: tukcnst.INSERT}
					newsubs.Subscriptions = append(newsubs.Subscriptions, newSub)
					l("Registering Subscription with Event Service", true)
					if err := tukdbint.NewDBEvent(&newsubs); err != nil {
						l(err.Error(), false)
						return err
					}
					i.Subs.Subscriptions = append(i.Subs.Subscriptions, newSub)
					if i.Subs.LastInsertId < int64(newSub.Id) {
						i.Subs.LastInsertId = int64(newSub.Id)
					}
					i.Subs.Count = i.Subs.Count + 1
				} else {
					l("No Broker Reference found. Subscription not registered with Event Service", true)
				}
			}
		}
	}
	i.Response, _ = json.MarshalIndent(i.Subs, "", "  ")
	return nil
}
func (i *DSUBEvent) selectSubscriptions() error {
	i.Subs = tukdbint.Subscriptions{Action: tukcnst.SELECT}
	sub := tukdbint.Subscription{}
	if i.Pathway != "" {
		sub.Pathway = i.Pathway
	}
	i.Subs.Subscriptions = append(i.Subs.Subscriptions, sub)
	if err := tukdbint.NewDBEvent(&i.Subs); err != nil {
		return err
	}
	i.Response, _ = json.MarshalIndent(i.Subs, "", "  ")
	return nil
}
func l(msg string, debug bool) {
	if !debug {
		log.Println(msg)
	} else {
		if DebugMode {
			log.Println(msg)
		}
	}
}
