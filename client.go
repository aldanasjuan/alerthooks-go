package alerthooks

import (
	"fmt"
	"net"
	"net/url"
	"strings"
	"time"

	jsoniter "github.com/json-iterator/go"
	"github.com/valyala/fasthttp"
)

const (
	STATUS_OPEN      RecordStatus = "open"
	STATUS_RETRYING  RecordStatus = "retrying"
	STATUS_CANCELED  RecordStatus = "canceled"
	STATUS_FAILED    RecordStatus = "failed"
	STATUS_SUCCEDED  RecordStatus = "succeded"
	RECORD_ONE_TIME  RecordType   = "one_time"
	RECORD_RECURRING RecordType   = "recurring"
)

const API = "https://alerts.rubbey.app"
const NEW_RECORD_URL = API + "/records"

var key string
var validMethods = []string{"GET", "POST", "PUT", "PATCH", "DELETE"}
var validTypes = []RecordType{RECORD_ONE_TIME, RECORD_RECURRING}

type RecordStatus string
type RecordType string

type Record struct {
	ID          string                 `json:"id"`
	Method      string                 `json:"method"`
	Endpoint    string                 `json:"endpoint"`
	Type        RecordType             `json:"type"`
	Status      RecordStatus           `json:"status"`
	DueDate     int64                  `json:"due_date"`
	CreatedAt   int64                  `json:"created_at"`
	CompletedAt int64                  `json:"completed_at"`
	Done        bool                   `json:"done"`
	Data        map[string]interface{} `json:"data"`
	Recurring   *Recurring             `json:"recurring"`
}

type NewRecordParams struct {
	Method    string                 `json:"method"`
	Endpoint  string                 `json:"endpoint"`
	Type      RecordType             `json:"type"`
	DueDate   int64                  `json:"due_date"`
	Data      map[string]interface{} `json:"data"`
	Recurring *Recurring             `json:"recurring"`
}

type Recurring struct {
	Minutes []int `json:"minutes,omitempty"` //0-59
	Hours   []int `json:"hours,omitempty"`   //0-23
	Days    []int `json:"days,omitempty"`    //1-31
	Months  []int `json:"months,omitempty"`  //1-12
}

func SetKey(k string) {
	key = k
}

func NewRecord(params *NewRecordParams) (*Record, error) {
	err := params.validate()
	if err != nil {
		return nil, err
	}
	b, err := jsoniter.Marshal(params)
	if err != nil {
		return nil, err
	}
	req := fasthttp.AcquireRequest()
	defer fasthttp.ReleaseRequest(req)
	defer req.Reset()
	req.Header.SetMethod("POST")
	req.SetRequestURI(NEW_RECORD_URL)
	req.Header.Set("Key", key)
	req.Header.Set("Content-Type", "application/json")
	req.SetBody(b)

	c := fasthttp.Client{}
	res := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseResponse(res)
	defer res.Reset()

	err = c.Do(req, res)
	if err != nil {
		return nil, err
	}
	if res.StatusCode() != 200 {
		return nil, fmt.Errorf(`status: %v error: %v`, res.StatusCode(), string(res.Body()))
	}
	record := &Record{}
	err = jsoniter.Unmarshal(res.Body(), record)
	if err != nil {
		return nil, err
	}
	return record, nil
}

func UpdateRecord(params *NewRecordParams) error {

	return nil
}
func DeleteRecord(params *NewRecordParams) error {

	return nil
}

func (r *NewRecordParams) validate() error {
	if !validateEndpoint(r.Endpoint) {
		return ErrInvalidURL
	}
	if !validateMethod(r.Method) {
		return ErrInvalidMethod
	}
	if !validateType(r.Type) {
		return ErrInvalidType
	}
	if r.Type == RECORD_RECURRING {
		return r.Recurring.validate()
	} else {
		if r.DueDate <= time.Now().Unix() {
			return ErrInvalidDueDate
		}
	}
	return nil
}

func validateMethod(method string) bool {
	method = strings.ToUpper(method)
	for _, m := range validMethods {
		if method == m {
			return true
		}
	}
	return false
}
func validateType(t RecordType) bool {
	for _, tp := range validTypes {
		if tp == t {
			return true
		}
	}
	return false
}

func validateEndpoint(endpoint string) bool {
	u, err := url.Parse(endpoint)
	if err != nil {
		return false
	}
	if !u.IsAbs() || u.Scheme == "" || u.Host == "" {
		return false
	}

	if strings.Contains(u.Host, "localhost") || net.ParseIP(u.Host) != nil || u.Port() != "" || u.Scheme != "https" {
		return false
	}
	_, err = net.LookupHost(u.Host)
	return err == nil
}

func (r *Recurring) validate() error {

	if r == nil {
		return ErrNilRecurring
	}

	if len(r.Minutes) < 1 || len(r.Hours) < 1 || len(r.Days) < 1 || len(r.Months) < 1 {
		return ErrMissingRecurring
	}

	if r.Minutes[0] != -1 {
		for _, v := range r.Minutes {
			if v > 59 || v < 0 {
				return ErrBadMinutes
			}
		}
	}
	if r.Hours[0] != -1 {
		for _, v := range r.Hours {
			if v > 23 || v < 0 {
				return ErrBadHours
			}
		}
	}

	if r.Days[0] != -1 {
		for _, v := range r.Days {
			if v > 31 || v < 1 {
				return ErrBadDays
			}
		}
	}
	if r.Months[0] != -1 {
		for _, v := range r.Months {
			if v > 12 || v < 1 {
				return ErrBadMonths
			}
		}
	}

	return nil
}
