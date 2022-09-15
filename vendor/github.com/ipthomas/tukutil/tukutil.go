package tukutil

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"text/template"
	"time"

	cnst "github.com/ipthomas/tukcnst"

	"github.com/google/uuid"
)

var (
	ServerName, _ = os.Hostname()
	SeedRoot      = "1.2.40.0.13.1.1.3542466645."
	IdSeed        = getIdIncrementSeed(5)
	CodeSystem    = make(map[string]string)
)

// TemplateFuncMap returns a functionMap of tukutils for use in templates
func TemplateFuncMap() template.FuncMap {
	return template.FuncMap{
		"dtday":          Tuk_Day,
		"dtmonth":        Tuk_Month,
		"dtyear":         Tuk_Year,
		"mapval":         GetCodeSystemVal,
		"prettytime":     PrettyTime,
		"newuuid":        NewUuid,
		"newid":          Newid,
		"tuktime":        Time_Now,
		"simpledatetime": SimpleDateTime,
		"isafternow":     IsAfterNow,
		"duration":       GetDuration,
		"durationsince":  GetDurationSince,
		"hasexpired":     OHT_HasExpired,
		"completedate":   OHT_CompleteByDate,
	}
}
func SetResponseHeaderServer(tukServerName string) {
	ServerName = tukServerName
}
func SimpleDateTime() string {
	return Tuk_Year() + Tuk_Month() + Tuk_Day() + Tuk_Hour() + Tuk_Min() + Tuk_Sec()
}

// SplitXDWKey takes a string input (xdwkey) and returns the pathway and nhs id for the xdw
func SplitXDWKey(xdwkey string) (string, string) {
	var pwy string
	var nhs string
	if len(xdwkey) > 10 {
		pwy = xdwkey[:len(xdwkey)-10]
		nhs = strings.TrimPrefix(xdwkey, pwy)
	}
	log.Printf("Pathway = %s NHS ID = %s", pwy, nhs)
	return pwy, nhs
}

func WriteResponseHeaders(fn http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		w.Header().Set("Server", ServerName)
		if r.Header.Get(cnst.ACCEPT) == cnst.APPLICATION_JSON {
			w.Header().Set(cnst.CONTENT_TYPE, cnst.APPLICATION_JSON)
		} else {
			if r.Header.Get(cnst.ACCEPT) == cnst.APPLICATION_XML {
				w.Header().Set(cnst.CONTENT_TYPE, cnst.APPLICATION_XML)
			} else {
				w.Header().Set(cnst.CONTENT_TYPE, cnst.TEXT_HTML)
			}

		}
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Headers", "accept, Content-Type")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")

		fn(w, r)
	}
}

// SetCodeSystem takes a map input and sets the codesystem map with the input
func SetCodeSystem(cs map[string]string) {
	CodeSystem = cs
	log.Printf("Loaded %v code system key values", len(CodeSystem))
	Log(CodeSystem)
}

// LoadCodeSystemFile loads the `codesystem.json` file from the `configs` folder
func LoadCodeSystemFile(codesystemFile string) error {
	file, err := os.Open(codesystemFile)
	if err != nil {
		log.Println(err.Error())
		return err
	}
	if err = json.NewDecoder(file).Decode(&CodeSystem); err != nil {
		log.Println(err.Error())
		return err
	}
	log.Printf("Loaded %v code system key values", len(CodeSystem))
	Log(CodeSystem)
	return nil
}

// GetCodeSystemVal takes a string input (key) and returns from the codesystem the string value corresponding to the input (key)
func GetCodeSystemVal(key string) string {
	val, ok := CodeSystem[key]
	if ok {
		return val
	}
	return key
}

// CreateLog checks if the log folder exists and creates it if not. It then checks for a subfolder for the current year i.e. 2022 and creates it if it does not exist. It then checks for a log file with a name equal to the current day and month and extension .log i.e. 0905.log. If it exists log output is appended to the existing file otherwise a new log file is created.
func CreateLog(log_Folder string) *os.File {
	log.SetFlags(log.Ldate | log.Ltime | log.Lmicroseconds | log.Lshortfile)
	mdir := log_Folder
	if _, err := os.Stat(mdir); errors.Is(err, fs.ErrNotExist) {
		if e2 := os.Mkdir(mdir, 0700); e2 != nil {
			log.Println(err.Error())
			return nil
		}
	}
	dir := mdir + "/" + Tuk_Year()
	if _, err := os.Stat(dir); errors.Is(err, fs.ErrNotExist) {
		if e2 := os.Mkdir(dir, 0700); e2 != nil {
			log.Println(err.Error())
			return nil
		}
	}
	logFile, err := os.OpenFile(dir+"/"+Tuk_Month()+Tuk_Day()+".log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Println(err.Error())
		return nil
	}
	log.Println("Using log file - " + logFile.Name())
	log.SetOutput(logFile)
	log.Println("-----------------------------------------------------------------------------------")
	return logFile
}
func MonitorApp() {
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt, syscall.SIGTERM)
	go func() {
		signalType := <-ch
		signal.Stop(ch)
		fmt.Println("")
		fmt.Println("***")
		fmt.Println("Exit command received. Exiting...")
		exitcode := 0
		switch signalType {
		case os.Interrupt:
			log.Println("FATAL: CTRL+C pressed")
		case syscall.SIGTERM:
			log.Println("FATAL: SIGTERM detected")
			exitcode = 1
		}
		os.Exit(exitcode)
	}()
}

// Log takes any struc as input and logs out the struc as a json string
func Log(i interface{}) {
	b, _ := json.MarshalIndent(i, "", "  ")
	log.Println(string(b))
}

// returns unique id in format '1.2.40.0.13.1.1.3542466645.20211021090059143.32643'
// idroot constant - 1.2.40.0.13.1.1.3542466645.
// + datetime	   - 20211021090059143.
// + 5 digit seed  - 32643
// if state is maintained the seed is incremented after each call to newid() to ensure a unique id is generated.
// If state is not maintained the `new` datetime will ensure a unique id is generated.
func Newid() string {
	id := SeedRoot + dt_yyyyMMddhhmmSSsss() + "." + GetStringFromInt(IdSeed)
	IdSeed = IdSeed + 1
	return id
}
func OHT_HasExpired(startdate string, htDate string) bool {
	log.Printf("Calculating if Expired")
	if expireDate := OHT_CompleteByDate(startdate, htDate); expireDate == "NA" {
		return false
	} else {
		if exdate, err := time.Parse("2006-01-02 15:04:05", expireDate); err != nil {
			log.Println(err.Error())
			return false
		} else {
			log.Printf("Has Expired = %v", time.Now().After(exdate))
			return time.Now().After(exdate)
		}
	}
}

// OHT_CompletionByDate takes a 'start date' input as a rfc3339 formatted string and the number of days in the future as a string containing an OASIS Human Task xdw function `day(x)“ where x is number of days in the future from the `start date`
// It returns a string containing the future date. If there is an error it returns `NA`
//
//	For example - OHT_CompleteByDate("2022-09-06T23:06:54+01:00", "day(10)") returns `2022-09-16 23:06:54`
func OHT_CompleteByDate(startdate string, htDate string) string {
	var err error
	var futuredays int
	var st time.Time
	var ft time.Time
	if strings.Contains(htDate, "(") && strings.Contains(htDate, ")") {
		if futuredays, err = strconv.Atoi(strings.Split(strings.Split(htDate, "(")[1], ")")[0]); err == nil {
			log.Printf("Calculating date %v Days from %s", futuredays, startdate)
			if st, err = time.Parse(time.RFC3339, startdate); err == nil {
				ft = st.Add(time.Hour * 24 * time.Duration(futuredays))
				log.Printf("Complete By %s", PrettyTime(ft.String()))
				return PrettyTime(ft.String())
			}
		}
	}
	log.Println(err.Error())
	return "NA"
}

// GetDurationSince takes a time as string input in RFC3339 format (yyyy-MM-ddThh:mm:ssZ) and returns the duration in days, hours and mins in a 'pretty format' eg '2 Days 0 Hrs 52 Mins' between the provided time and time.Now() as a string
func GetDurationSince(stime string) string {
	log.Println("Obtaining time Duration since - " + stime)
	st, err := time.Parse(time.RFC3339, stime)
	if err != nil {
		log.Println(err.Error())
		return "Not Available"
	}
	dur := time.Since(st)
	log.Printf("Duration - %v", dur.String())
	days := 0
	hrs := int(dur.Hours())
	min := int(dur.Minutes())

	if hrs > 24 {
		days = hrs / 24
		hrs = hrs % 24
	}
	daysstr := strconv.Itoa(days)
	hrsstr := strconv.Itoa(hrs)
	minstr := strconv.Itoa(min - (days * 24 * 60) - (hrs * 60))
	log.Println("Returning " + daysstr + " Days " + hrsstr + " Hrs " + minstr + " Mins")
	return daysstr + " Days " + hrsstr + " Hrs " + minstr + " Mins"
}

// GetDuration takes 2 times as string inputs in RFC3339 format (yyyy-MM-ddThh:mm:ssZ) and returns the duration in days, hours and mins in a 'pretty format' eg '2 Days 0 Hrs 52 Mins' between the provided times as a string
//
//	Example : GetDuration("2022-09-04T13:15:20Z", "2022-09-14T16:20:01Z") returns `10 Days 3 Hrs 4 Mins`
func GetDuration(stime string, etime string) string {
	log.Println("Obtaining time Duration between - " + stime + " and " + etime)
	st, err := time.Parse(time.RFC3339, stime)
	if err != nil {
		log.Println(err.Error())
		return "Not Available"
	}
	et, err := time.Parse(time.RFC3339, etime)
	if err != nil {
		log.Println(err.Error())
		return "Not Available"
	}
	dur := et.Sub(st)
	days := 0
	hrs := int(dur.Hours())
	min := int(dur.Minutes())
	if hrs > 24 {
		days = hrs / 24
		hrs = hrs % 24
	}
	daysstr := strconv.Itoa(days)
	hrsstr := strconv.Itoa(hrs)
	minstr := strconv.Itoa(min - (days * 24 * 60) - (hrs * 60))
	log.Println("Returning " + daysstr + " Days " + hrsstr + " Hrs " + minstr + " Mins")
	return daysstr + " Days " + hrsstr + " Hrs " + minstr + " Mins"
}

// IsAfterNow takes a time as a string input in RFC3339 format (yyyy-MM-ddThh:mm:ssZ) and returns true if the input time is after time.Now() and false if input time is before time.Now()
func IsAfterNow(inTime string) bool {
	log.Printf("Checking if %s is after the current time", inTime)
	it, err := time.Parse(time.RFC3339, inTime)
	if err != nil {
		log.Println(err.Error())
		return false
	}
	now := time.Now().Local()
	log.Println("Time Now - " + now.Local().String())
	log.Println("Start Time - " + it.Local().String())
	log.Printf("Time %s IsAfter(time.Now()) = %v", inTime, now.Before(it))
	return now.Before(it)
}

// Pretty_Time_Now returns a pretty version of the current time for location Europe/London (strips everything after the `.` in Tuk_Time)
func Pretty_Time_Now() string {
	return PrettyTime(Time_Now())
}

// Time_Now returns the current time for location Europe/London.
func Time_Now() string {
	location, err := time.LoadLocation("Europe/London")
	if err != nil {
		log.Println(err.Error())
		return time.Now().String()
	}
	return time.Now().In(location).String()
}

// PrettyTime fist splits the input based on sep =`.`, it takes index 0 of the split and replaces any T with a space then removes any trailing Z. It then splits the resulting string on sep = `+` returning index 0 of the split
func PrettyTime(time string) string {
	return strings.TrimSuffix(strings.Split(strings.TrimSuffix(strings.ReplaceAll(strings.Split(time, ".")[0], "T", " "), "Z"), "+")[0], " ")
}

// TUK_Hour returns the current hour as a 2 digit string
func Tuk_Hour() string {
	return fmt.Sprintf("%02d",
		time.Now().Local().Hour())
}

// TUK_Min returns the current minute as a 2 digit string
func Tuk_Min() string {
	return fmt.Sprintf("%02d", time.Now().Local().Minute())
}

// TUK_Sec returns the current second as a 2 digit string
func Tuk_Sec() string {
	return fmt.Sprintf("%02d",
		time.Now().Local().Second())
}

// TUK_MilliSec returns the current milliseconds as a 3 digit int
func Tuk_MilliSec() int {
	return GetIntFromString(Substr(GetStringFromInt(time.Now().Nanosecond()), 0, 3))
}

// TUK_Day returns the current day as a 2 digit string
func Tuk_Day() string {
	return fmt.Sprintf("%02d",
		time.Now().Local().Day())
}

// TUK_Year returns the current year as a 4 digit string
func Tuk_Year() string {
	return fmt.Sprintf("%d",
		time.Now().Local().Year())
}

// TUK_Month returns the current month as a 2 digit string
func Tuk_Month() string {
	return fmt.Sprintf("%02d",
		time.Now().Local().Month())
}

// NewUuid returns a random UUID as a string
func NewUuid() string {
	u := uuid.New()
	return u.String()
}

// GetStringFromInt takes a int input and returns a string of that value.
func GetStringFromInt(i int) string {
	return strconv.Itoa(i)
}

// GetIntFromString takes a string input with an integer value and returns an int of that value. If the input is not numeric, 0 is returned
func GetIntFromString(input string) int {
	i, _ := strconv.Atoi(input)
	return i
}

// Substr takes a string input and returns the rune (Substring) defined by the start and length in th start and length input values
func Substr(input string, start int, length int) string {
	asRunes := []rune(input)
	if start >= len(asRunes) {
		return ""
	}
	if start+length > len(asRunes) {
		length = len(asRunes) - start
	}
	return string(asRunes[start : start+length])
}

// GetXMLNodeList takes an xml message as input and returns the xml node list as a string for the node input value provide
func GetXMLNodeList(message string, node string) string {
	if strings.Contains(message, node) {
		var nodeopen = "<" + node
		var nodeclose = "</" + node + ">"
		log.Println("Searching for XML Element: " + nodeopen + ">")
		var start = strings.Index(message, nodeopen)
		var end = strings.Index(message, nodeclose) + len(nodeclose)
		m := message[start:end]
		log.Println("Extracted XML Element Nodelist")
		return m
	}
	log.Println("Message does not contain Element : " + node)
	return ""
}

// PrettyAuthorInstitution takes a string input (XDS Author.Institution format) and returns a string with just the Institution name
func PrettyAuthorInstitution(institution string) string {
	if strings.Contains(institution, "^") {
		return strings.Split(institution, "^")[0] + ","
	}
	return institution
}

// PrettyAuthorPerson takes a string input (XDS Author.Person format) and returns a string with the person last and first names
func PrettyAuthorPerson(author string) string {
	if strings.Contains(author, "^") {
		authorsplit := strings.Split(author, "^")
		if len(authorsplit) > 2 {
			return authorsplit[1] + " " + authorsplit[2]
		}
		if len(authorsplit) > 1 {
			return authorsplit[1]
		}
	}
	return author
}

// GetFolderFiles takes a string input of the complete folder path and returns a fs.DirEntry
func GetFolderFiles(folder string) ([]fs.DirEntry, error) {
	var err error
	var f *os.File
	var fileInfo []fs.DirEntry
	f, err = os.Open(folder)
	if err != nil {
		log.Println(err)
		return fileInfo, err
	}
	fileInfo, err = f.ReadDir(-1)
	f.Close()
	if err != nil {
		log.Println(err)
	}
	return fileInfo, err
}
func getIdIncrementSeed(len int) int {
	return GetIntFromString(Substr(GetStringFromInt(time.Now().Nanosecond()), 0, len))
}
func dt_yyyyMMddhhmmSSsss() string {
	return Tuk_Year() + Tuk_Month() + Tuk_Day() + Tuk_Hour() + Tuk_Min() + Tuk_Sec() + strconv.Itoa(Tuk_MilliSec())
}
