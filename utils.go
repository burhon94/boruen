package boruen

import (
	"database/sql"
	"encoding/json"
	"reflect"
	"strconv"
)

type FieldGlobalInterface interface {
	GetValidField() bool
	GetValue() interface{}
}

// for comfortable using and handling null types from datastore

type NullString sql.NullString

func (n NullString) GetValidField() bool {
	return n.Valid
}

func (n NullString) GetValue() interface{} {
	return n.String
}

func (n *NullString) MarshalJSON() ([]byte, error) {
	if !n.Valid {
		return []byte("null"), nil
	}

	return json.Marshal(n.String)
}

func (n *NullString) UnmarshalJSON(data []byte) error {
	//todo::null or not filled
	if string(data) == "null" {
		n.String, n.Valid = "", false
		return nil
	}

	n.String, n.Valid = string(data), true
	return nil
}

func (n *NullString) Scan(value interface{}) error {
	var i sql.NullString
	if err := i.Scan(value); err != nil {
		return err
	}
	if reflect.TypeOf(value) == nil {
		*n = NullString{
			Valid:  false,
			String: i.String,
		}
	} else {
		*n = NullString{
			Valid:  true,
			String: i.String,
		}
	}
	return nil
}

type NullInt32 sql.NullInt32

func (ni NullInt32) GetValidField() bool {
	return ni.Valid
}
func (ni NullInt32) GetValue() interface{} {
	return ni.Int32
}

func (ni *NullInt32) MarshalJSON() ([]byte, error) {
	if !ni.Valid {
		return []byte("null"), nil
	}

	return json.Marshal(ni.Int32)
}

func (ni *NullInt32) UnmarshalJSON(data []byte) error {

	if string(data) == "null" || string(data) == "" {
		ni.Int32, ni.Valid = 0, false
		return nil
	}

	intVal, err := strconv.Atoi(string(data))

	if err != nil {
		return err
	}

	ni.Int32, ni.Valid = int32(intVal), true
	return nil
}

func (ni *NullInt32) Scan(value interface{}) error {
	var i sql.NullInt32
	if err := i.Scan(value); err != nil {
		return err
	}

	if reflect.TypeOf(value) == nil {
		*ni = NullInt32{
			Valid: false,
			Int32: i.Int32,
		}
	} else {
		*ni = NullInt32{
			Valid: true,
			Int32: i.Int32,
		}
	}

	return nil
}

// for testing
var (
	ById = func(val string) string {
		return " and id = " + val
	}
	ByType = func(val string) string {
		return " and type = '" + val + "'"
	}
	BySourceAccMethod = func(val string) string {
		return " and source_acc_method = " + val
	}
	BySourceAccGate = func(val string) string {
		return " and source_acc_gate = " + val
	}
	ByDestinationAccMethod = func(val string) string {
		return " and dest_acc_method = " + val
	}
	ByDestinationAccGate = func(val string) string {
		return " and dest_acc_gate = " + val
	}
	ByOperationType = func(val string) string {
		return " and dest_acc_gate = " + val
	}
	ByLocationID = func(val string) string {
		return " and operation_type = " + val
	}
	ByAudience = func(val string) string {
		return " and audience = '" + val + "'"
	}
	ByProviderID = func(val string) string {
		return " and provider_id = " + val
	}
	ByTerminalID = func(val string) string {
		return " and terminal_id = '" + val + "'"
	}
	ByTransactionType = func(val string) string {
		return " and transaction_type = '" + val + "'"
	}
	ByAwardRate = func(val string) string {
		return " and award_rate = " + val
	}
	ByAwardType = func(val string) string {
		return " and award_type = '" + val + "'"
	}
	ByMinAmount = func(amount string) string {
		return " and min_amount <= " + amount
	}
	ByMaxAmount = func(amount string) string {
		return " and max_amount >= " + amount
	}
	ByUppLimit = func(amount string) string {
		//
		// some of fileds allways must be 1
		return ` and up_limit >= (select case 
				when award_type = 'percent' then award_rate * ` + amount +
			"\n else award_rate )"
	}

	ByConditionText = func(val string) string {
		return " and condition_text like '%" + val + "%'"
	}
	ByStatus = func(val string) string {
		return " and status = '" + val + "'"
	}
	ByUsePeriod = func(val string) string {
		return " and use_period = " + val
	}
	ByDateFrom = func(date string) string {
		return " and to_date('" + date + "', 'DD.MM.YYYY') <= created_at"
	}
	ByDateTo = func(date string) string {
		return " and to_date('" + date + "', 'DD.MM.YYYY') + interval '1' day > created_at"
	}
)

var (
	typeItems = map[int]string{
		0: "momentaly",
		1: "monthly",
		2: "",
	}
	audience = map[int]string{
		0: "user",
		1: "referrer",
		2: "referral",
		3: "",
	}
	transactionType = map[int]string{
		0: "income",
		1: "expense",
	}
	awardType = map[int]string{
		0: "percent",
		1: "sum",
	}

	statusOfTrx = map[int]string{
		0: "approved",
		1: "canceled",
		2: "pending",
		3: "error",
	}

	ruleStatus = map[int]string{
		1: "active",
		0: "inactive",
	}
)
