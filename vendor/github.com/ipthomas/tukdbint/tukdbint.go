package tukdbint

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"reflect"
	"strings"
	"time"

	"github.com/ipthomas/tukcnst"
	"github.com/ipthomas/tukhttp"

	_ "github.com/go-sql-driver/mysql"
)

var (
	DB_URL = ""
	DBConn *sql.DB
)

type TukDBConnection struct {
	DBUser        string
	DBPassword    string
	DBHost        string
	DBPort        string
	DBName        string
	DBTimeout     string
	DBReadTimeout string
	DB_URL        string
	DBReader_Only bool
}
type Subscription struct {
	Id         int    `json:"id"`
	Created    string `json:"created"`
	BrokerRef  string `json:"brokerref"`
	Pathway    string `json:"pathway"`
	Topic      string `json:"topic"`
	Expression string `json:"expression"`
}
type Subscriptions struct {
	Action        string         `json:"action"`
	LastInsertId  int64          `json:"lastinsertid"`
	Count         int            `json:"count"`
	Subscriptions []Subscription `json:"Subscriptions"`
}
type Event struct {
	EventId            int64  `json:"eventid"`
	Creationtime       string `json:"creationtime"`
	DocName            string `json:"docname"`
	ClassCode          string `json:"classcode"`
	ConfCode           string `json:"confcode"`
	FormatCode         string `json:"formatcode"`
	FacilityCode       string `json:"facilitycode"`
	PracticeCode       string `json:"practicecode"`
	Expression         string `json:"expression"`
	Authors            string `json:"authors"`
	XdsPid             string `json:"xdspid"`
	XdsDocEntryUid     string `json:"xdsdocentryuid"`
	RepositoryUniqueId string `json:"repositoryuniqueid"`
	NhsId              string `json:"nhsid"`
	User               string `json:"user"`
	Org                string `json:"org"`
	Role               string `json:"role"`
	Topic              string `json:"topic"`
	Pathway            string `json:"pathway"`
	Notes              string `json:"notes"`
	Version            string `json:"ver"`
	BrokerRef          string `json:"brokerref"`
}
type Events struct {
	Action       string  `json:"action"`
	LastInsertId int64   `json:"lastinsertid"`
	Count        int     `json:"count"`
	Events       []Event `json:"events"`
}
type Workflow struct {
	Id        int    `json:"id"`
	Created   string `json:"created"`
	XDW_Key   string `json:"xdw_key"`
	XDW_UID   string `json:"xdw_uid"`
	XDW_Doc   string `json:"xdw_doc"`
	XDW_Def   string `json:"xdw_def"`
	Version   int    `json:"version"`
	Published bool   `json:"published"`
}
type Workflows struct {
	Action       string     `json:"action"`
	LastInsertId int64      `json:"lastinsertid"`
	Count        int        `json:"count"`
	Workflows    []Workflow `json:"workflows"`
}
type XDWS struct {
	Action       string `json:"action"`
	LastInsertId int64  `json:"lastinsertid"`
	Count        int    `json:"count"`
	XDW          []XDW  `json:"xdws"`
}
type XDW struct {
	Id        int    `json:"id"`
	Name      string `json:"name"`
	IsXDSMeta bool   `json:"isxdsmeta"`
	XDW       string `json:"xdw"`
}
type IdMaps struct {
	Action       string
	LastInsertId int64
	Where        string
	Value        string
	Cnt          int
	LidMap       []IdMap
}
type IdMap struct {
	Id  int    `json:"id"`
	Lid string `json:"lid"`
	Mid string `json:"mid"`
}
type EventAcks struct {
	Action       string
	LastInsertId int64
	Where        string
	Value        string
	Cnt          int
	EventAck     []EventAck
}
type EventAck struct {
	Id           int    `json:"id"`
	CreationTime string `json:"creationtime"`
	SubRef       string `json:"subref"`
	EventID      int    `json:"eventid"`
}
type TUK_DB_Interface interface {
	newEvent() error
}

// NewDBEvent takes an Interface (struct Events, Workflows, Subscriptions, XDWS, TukDBConnection) and executes a mysql request. The response is poputlated into the interface struct.
// If DB_URL = "" the sql query is executed against the DBConn established using a DSN, if a value is present in DB_URL, the query is sent to the API Gateway URL (DB_URL) as a json string of the provided interface struct
func NewDBEvent(i TUK_DB_Interface) error {
	return i.newEvent()
}

// functions for mysql DB access via DSN

func (i *TukDBConnection) newEvent() error {
	var err error
	if i.DB_URL != "" {
		log.Println("Database API URL provided. Will connect to mysql instance via AWS API Gateway url " + i.DB_URL)
		DB_URL = i.DB_URL
	} else {
		if i.DBUser == "" {
			i.DBUser = "root"
		}
		if i.DBPassword == "" {
			i.DBPassword = "rootPass"
		}
		if i.DBHost == "" {
			i.DBHost = "localhost"
		}
		if i.DBPort == "" {
			i.DBPort = ":3306"
		} else {
			if !strings.HasPrefix(i.DBPort, ":") {
				i.DBPort = ":" + i.DBPort
			}
		}
		if i.DBName == "" {
			i.DBName = "tuk"
		}
		if i.DBTimeout == "" {
			i.DBTimeout = "60s"
		} else {
			if !strings.HasSuffix(i.DBTimeout, "s") {
				i.DBTimeout = i.DBTimeout + "s"
			}
		}
		if i.DBReadTimeout == "" {
			i.DBReadTimeout = "2s"
		} else {
			if !strings.HasSuffix(i.DBReadTimeout, "s") {
				i.DBReadTimeout = i.DBReadTimeout + "s"
			}
		}
		dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?parseTime=true&timeout=%s&readTimeout=%s",
			i.DBUser,
			i.DBPassword,
			i.DBHost+i.DBPort,
			i.DBName,
			i.DBTimeout,
			i.DBReadTimeout)
		log.Println("No Database API URL provided. Opening DB Connection to mysql instance via DSN - " + dsn)
		DBConn, err = sql.Open(tukcnst.MYSQL, dsn)
	}

	return err
}
func (i *Subscriptions) newEvent() error {
	if DB_URL != "" {
		return i.newAWSEvent()
	}
	var err error
	var stmntStr = tukcnst.SQL_DEFAULT_SUBSCRIPTIONS
	var rows *sql.Rows
	var vals []interface{}
	ctx, cancelCtx := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancelCtx()
	if len(i.Subscriptions) > 0 {
		if stmntStr, vals, err = createPreparedStmnt(i.Action, tukcnst.SUBSCRIPTIONS, reflectStruct(reflect.ValueOf(i.Subscriptions[0]))); err != nil {
			log.Println(err.Error())
			return err
		}
	}
	sqlStmnt, err := DBConn.PrepareContext(ctx, stmntStr)
	if err != nil {
		log.Println(err.Error())
		return err
	}
	defer sqlStmnt.Close()

	if i.Action == tukcnst.SELECT {
		rows, err = setRows(ctx, sqlStmnt, vals)
		if err != nil {
			log.Println(err.Error())
			return err
		}

		for rows.Next() {
			sub := Subscription{}
			if err := rows.Scan(&sub.Id, &sub.Created, &sub.BrokerRef, &sub.Pathway, &sub.Topic, &sub.Expression); err != nil {
				switch {
				case err == sql.ErrNoRows:
					return nil
				default:
					log.Println(err.Error())
					return err
				}
			}
			i.Subscriptions = append(i.Subscriptions, sub)
			i.Count = i.Count + 1
		}
	} else {
		i.LastInsertId, err = setLastID(ctx, sqlStmnt, vals)
	}
	return err
}
func (i *Events) newEvent() error {
	if DB_URL != "" {
		return i.newAWSEvent()
	}
	var err error
	var stmntStr = tukcnst.SQL_DEFAULT_EVENTS
	var rows *sql.Rows
	var vals []interface{}
	ctx, cancelCtx := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancelCtx()
	if len(i.Events) > 0 {
		if stmntStr, vals, err = createPreparedStmnt(i.Action, tukcnst.EVENTS, reflectStruct(reflect.ValueOf(i.Events[0]))); err != nil {
			log.Println(err.Error())
			return err
		}
	}
	sqlStmnt, err := DBConn.PrepareContext(ctx, stmntStr)
	if err != nil {
		log.Println(err.Error())
		return err
	}
	defer sqlStmnt.Close()

	if i.Action == tukcnst.SELECT {
		rows, err = setRows(ctx, sqlStmnt, vals)
		if err != nil {
			log.Println(err.Error())
			return err
		}

		for rows.Next() {
			ev := Event{}
			if err := rows.Scan(&ev.EventId, &ev.Creationtime, &ev.DocName, &ev.ClassCode, &ev.ConfCode, &ev.FormatCode, &ev.FacilityCode, &ev.PracticeCode, &ev.Expression, &ev.Authors, &ev.XdsPid, &ev.XdsDocEntryUid, &ev.RepositoryUniqueId, &ev.NhsId, &ev.User, &ev.Org, &ev.Role, &ev.Topic, &ev.Pathway, &ev.Notes, &ev.Version); err != nil {
				switch {
				case err == sql.ErrNoRows:
					return nil
				default:
					log.Println(err.Error())
					return err
				}
			}
			i.Events = append(i.Events, ev)
			i.Count = i.Count + 1
		}
	} else {
		i.LastInsertId, err = setLastID(ctx, sqlStmnt, vals)
	}
	return err
}
func (i *Workflows) newEvent() error {
	if DB_URL != "" {
		return i.newAWSEvent()
	}
	var err error
	var stmntStr = tukcnst.SQL_DEFAULT_WORKFLOWS
	var rows *sql.Rows
	var vals []interface{}
	ctx, cancelCtx := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancelCtx()
	if len(i.Workflows) > 0 {
		if stmntStr, vals, err = createPreparedStmnt(i.Action, tukcnst.WORKFLOWS, reflectStruct(reflect.ValueOf(i.Workflows[0]))); err != nil {
			log.Println(err.Error())
			return err
		}
	}
	sqlStmnt, err := DBConn.PrepareContext(ctx, stmntStr)
	if err != nil {
		log.Println(err.Error())
		return err
	}
	defer sqlStmnt.Close()

	if i.Action == tukcnst.SELECT {
		rows, err = setRows(ctx, sqlStmnt, vals)
		if err != nil {
			log.Println(err.Error())
			return err
		}
		for rows.Next() {
			workflow := Workflow{}
			if err := rows.Scan(&workflow.Id, &workflow.Created, &workflow.XDW_Key, &workflow.XDW_UID, &workflow.XDW_Doc, &workflow.XDW_Def, &workflow.Version, &workflow.Published); err != nil {
				switch {
				case err == sql.ErrNoRows:
					return nil
				default:
					log.Println(err.Error())
					return err
				}
			}
			i.Workflows = append(i.Workflows, workflow)
			i.Count = i.Count + 1
		}
	} else {
		i.LastInsertId, err = setLastID(ctx, sqlStmnt, vals)
	}
	return err
}
func (i *XDWS) newEvent() error {
	if DB_URL != "" {
		return i.newAWSEvent()
	}
	var err error
	var stmntStr = tukcnst.SQL_DEFAULT_XDWS
	var rows *sql.Rows
	var vals []interface{}
	ctx, cancelCtx := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancelCtx()
	if len(i.XDW) > 0 {
		if stmntStr, vals, err = createPreparedStmnt(i.Action, tukcnst.XDWS, reflectStruct(reflect.ValueOf(i.XDW[0]))); err != nil {
			log.Println(err.Error())
			return err
		}
	}
	sqlStmnt, err := DBConn.PrepareContext(ctx, stmntStr)
	if err != nil {
		log.Println(err.Error())
		return err
	}
	defer sqlStmnt.Close()

	if i.Action == tukcnst.SELECT {
		rows, err = setRows(ctx, sqlStmnt, vals)
		if err != nil {
			log.Println(err.Error())
			return err
		}
		for rows.Next() {
			xdw := XDW{}
			if err := rows.Scan(&xdw.Id, &xdw.Name, &xdw.IsXDSMeta, &xdw.XDW); err != nil {
				switch {
				case err == sql.ErrNoRows:
					return nil
				default:
					log.Println(err.Error())
					return err
				}
			}
			i.XDW = append(i.XDW, xdw)
			i.Count = i.Count + 1
		}
	} else {
		i.LastInsertId, err = setLastID(ctx, sqlStmnt, vals)
	}
	return err
}
func (i *IdMaps) newEvent() error {
	if DB_URL != "" {
		return i.newAWSEvent()
	}
	var err error
	var stmntStr = tukcnst.SQL_DEFAULT_IDMAPS
	var rows *sql.Rows
	var vals []interface{}
	ctx, cancelCtx := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancelCtx()
	if len(i.LidMap) > 0 {
		if stmntStr, vals, err = createPreparedStmnt(i.Action, tukcnst.ID_MAPS, reflectStruct(reflect.ValueOf(i.LidMap[0]))); err != nil {
			log.Println(err.Error())
			return err
		}
	}
	sqlStmnt, err := DBConn.PrepareContext(ctx, stmntStr)
	if err != nil {
		log.Println(err.Error())
		return err
	}
	defer sqlStmnt.Close()

	if i.Action == tukcnst.SELECT {
		rows, err = setRows(ctx, sqlStmnt, vals)
		if err != nil {
			log.Println(err.Error())
			return err
		}
		for rows.Next() {
			idmap := IdMap{}
			if err := rows.Scan(&idmap.Id, &idmap.Lid, &idmap.Mid); err != nil {
				switch {
				case err == sql.ErrNoRows:
					return nil
				default:
					log.Println(err.Error())
					return err
				}
			}
			i.LidMap = append(i.LidMap, idmap)
			i.Cnt = i.Cnt + 1
		}
	} else {
		i.LastInsertId, err = setLastID(ctx, sqlStmnt, vals)
	}
	return err
}
func (i *EventAcks) newEvent() error {
	if DB_URL != "" {
		return i.newAWSEvent()
	}
	var err error
	var stmntStr = tukcnst.SQL_DEFAULT_EVENT_ACKS
	var rows *sql.Rows
	var vals []interface{}
	ctx, cancelCtx := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancelCtx()
	if len(i.EventAck) > 0 {
		if stmntStr, vals, err = createPreparedStmnt(i.Action, tukcnst.EVENT_ACKS, reflectStruct(reflect.ValueOf(i.EventAck[0]))); err != nil {
			log.Println(err.Error())
			return err
		}
	}
	sqlStmnt, err := DBConn.PrepareContext(ctx, stmntStr)
	if err != nil {
		log.Println(err.Error())
		return err
	}
	defer sqlStmnt.Close()

	if i.Action == tukcnst.SELECT {
		rows, err = setRows(ctx, sqlStmnt, vals)
		if err != nil {
			log.Println(err.Error())
			return err
		}
		for rows.Next() {
			eventack := EventAck{}
			if err := rows.Scan(&eventack.Id, &eventack.CreationTime, &eventack.SubRef, &eventack.EventID); err != nil {
				switch {
				case err == sql.ErrNoRows:
					return nil
				default:
					log.Println(err.Error())
					return err
				}
			}
			i.EventAck = append(i.EventAck, eventack)
			i.Cnt = i.Cnt + 1
		}
	} else {
		i.LastInsertId, err = setLastID(ctx, sqlStmnt, vals)
	}
	return err
}
func reflectStruct(i reflect.Value) map[string]interface{} {
	params := make(map[string]interface{})
	structType := i.Type()
	for f := 0; f < i.NumField(); f++ {
		if structType.Field(f).Name == "Id" {
			tid := i.Field(f).Interface().(int)
			if tid > 0 {
				params[strings.ToLower(structType.Field(f).Name)] = tid
			}
		}
		if structType.Field(f).Name != "Id" && i.Field(f).Interface() != "" {
			log.Printf("Reflecting Field %s Value %v", structType.Field(f).Name, i.Field(f).Interface())
			params[strings.ToLower(structType.Field(f).Name)] = i.Field(f).Interface()
		}
	}
	log.Printf("Obtained %v Key Values - %s", len(params), params)
	return params
}
func createPreparedStmnt(action string, table string, params map[string]interface{}) (string, []interface{}, error) {
	var vals []interface{}
	stmntStr := "SELECT * FROM " + table
	if len(params) > 0 {
		switch action {
		case tukcnst.SELECT:
			var paramStr string
			stmntStr = stmntStr + " WHERE "
			for param, val := range params {
				paramStr = paramStr + param + "= ? AND "
				vals = append(vals, val)
			}
			paramStr = strings.TrimSuffix(paramStr, " AND ")
			stmntStr = stmntStr + paramStr
		case tukcnst.INSERT:
			var paramStr string
			var qStr string
			stmntStr = "INSERT INTO " + table + " ("
			for param, val := range params {
				paramStr = paramStr + param + ", "
				qStr = qStr + "?, "
				vals = append(vals, val)
			}
			paramStr = strings.TrimSuffix(paramStr, ", ") + ") VALUES ("
			qStr = strings.TrimSuffix(qStr, ", ")
			stmntStr = stmntStr + paramStr + qStr + ")"
		case tukcnst.DEPRECATE:
			switch table {
			case tukcnst.WORKFLOWS:
				stmntStr = "UPDATE workflows SET version = version + 1 WHERE xdw_key=?"
				vals = append(vals, params["xdw_key"])
			case tukcnst.EVENTS:
				stmntStr = "UPDATE events SET version = version + 1 WHERE pathway=? AND nhsid=?"
				vals = append(vals, params["pathway"])
				vals = append(vals, params["nhsid"])
			}
		case tukcnst.UPDATE:
			switch table {
			case tukcnst.WORKFLOWS:
				stmntStr = "UPDATE workflows SET xdw_doc = ?, published = ? WHERE xdw_key = ? AND version = 0"
				vals = append(vals, params["xdw_doc"])
				vals = append(vals, params["published"])
				vals = append(vals, params["xdw_key"])
			}
		case tukcnst.DELETE:
			stmntStr = "DELETE FROM " + table + " WHERE "
			var paramStr string
			for param, val := range params {
				paramStr = paramStr + param + "= ? AND "
				vals = append(vals, val)
			}
			paramStr = strings.TrimSuffix(paramStr, " AND ")
			stmntStr = stmntStr + paramStr
		}
		log.Printf("Created Prepared Statement %s - Values %s", stmntStr, vals)
	}
	return stmntStr, vals, nil
}
func setRows(ctx context.Context, sqlStmnt *sql.Stmt, vals []interface{}) (*sql.Rows, error) {
	if len(vals) > 0 {
		return sqlStmnt.QueryContext(ctx, vals...)
	} else {
		return sqlStmnt.QueryContext(ctx)
	}
}
func setLastID(ctx context.Context, sqlStmnt *sql.Stmt, vals []interface{}) (int64, error) {
	if len(vals) > 0 {
		sqlrslt, err := sqlStmnt.ExecContext(ctx, vals...)
		if err != nil {
			log.Println(err.Error())
			return 0, err
		}
		id, err := sqlrslt.LastInsertId()
		if err != nil {
			log.Println(err.Error())
			return 0, err
		} else {
			return id, nil
		}
	} else {
		return 0, nil
	}
}

// functions for AWS Aurora DB Access via AWS API GW URL

func (i *Subscriptions) newAWSEvent() error {
	body, _ := json.Marshal(i)
	awsreq := aws_APIRequest(i.Action, tukcnst.SUBSCRIPTIONS, body)
	if err := tukhttp.NewRequest(&awsreq); err != nil {
		return err
	}
	return json.Unmarshal(awsreq.Response, &i)
}
func (i *Events) newAWSEvent() error {
	body, _ := json.Marshal(i)
	awsreq := aws_APIRequest(i.Action, tukcnst.EVENTS, body)
	if err := tukhttp.NewRequest(&awsreq); err != nil {
		return err
	}
	return json.Unmarshal(awsreq.Response, &i)
}
func (i *Workflows) newAWSEvent() error {
	body, _ := json.Marshal(i)
	awsreq := aws_APIRequest(i.Action, tukcnst.WORKFLOWS, body)
	if err := tukhttp.NewRequest(&awsreq); err != nil {
		return err
	}
	return json.Unmarshal(awsreq.Response, &i)
}
func (i *XDWS) newAWSEvent() error {
	body, _ := json.Marshal(i)
	awsreq := aws_APIRequest(i.Action, tukcnst.XDWS, body)
	if err := tukhttp.NewRequest(&awsreq); err != nil {
		return err
	}
	return json.Unmarshal(awsreq.Response, &i)
}
func (i *IdMaps) newAWSEvent() error {
	body, _ := json.Marshal(i)
	awsreq := aws_APIRequest(i.Action, tukcnst.ID_MAPS, body)
	if err := tukhttp.NewRequest(&awsreq); err != nil {
		return err
	}
	return json.Unmarshal(awsreq.Response, &i)
}
func (i *EventAcks) newAWSEvent() error {
	body, _ := json.Marshal(i)
	awsreq := aws_APIRequest(i.Action, tukcnst.EVENT_ACKS, body)
	if err := tukhttp.NewRequest(&awsreq); err != nil {
		return err
	}
	return json.Unmarshal(awsreq.Response, &i)
}
func aws_APIRequest(action string, resource string, body []byte) tukhttp.AWS_APIRequest {
	return tukhttp.AWS_APIRequest{
		URL:      DB_URL,
		Act:      action,
		Resource: resource,
		Timeout:  5,
		Body:     body,
	}
}
