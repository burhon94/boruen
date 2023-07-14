package boruen

import (
	"context"
	"reflect"

	"sort"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/lib/pq"
)

type RuleModule struct {
	db *pgxpool.Pool
}

func NewRuleEngine(dns string) RuleModule {

	var (
		RuleModule RuleModule
		err        error
	)

	// set getting connections from package
	RuleModule.db, err = pgxpool.Connect(context.Background(), dns)

	if err != nil {
		panic(err)
	}

	return RuleModule
}

func columns() string {

	return `id,
		type,
		transaction_type,
		source_acc_method,
		source_acc_gate,
		dest_acc_method,
		dest_acc_gate,
		operation_type,
		location_id,
		award_type,
		award_rate,
		min_amount,
		max_amount,
		up_limit,
		audience,
		condition_text,
		status,
		use_period,
		created_at::text,
		updated_at::text,
		priority_tags,
		terminal_id,
		provider_id `

}

func fields(r *Rule) []interface{} {

	return []interface{}{&r.ID,
		&r.Type,
		&r.TransactionType,
		&r.SourceAccMethod,
		&r.SourceAccGate,
		&r.DestAccMethod,
		&r.DestAccGate,
		&r.OperationType,
		&r.LocationID,
		&r.AwardType,
		&r.AwardRate,
		&r.MinAmount,
		&r.MaxAmount,
		&r.UppLimit,
		&r.Audience,
		&r.ConditionText,
		&r.Status,
		&r.UsePeriod,
		&r.CreatedAT,
		&r.UpdatedAT,
		&r.PriorityTags,
		&r.TerminalID,
		&r.ProviderID,
	}
}

// FindRuleForTransaction:
// rule tables will be sorted by count of null elements in cols by descending
// after that one by one transaction will be compared with rules
// if rule fits to transaction it rule will be returned
func (r *RuleModule) FindRuleForTransaction(transaction interface{}, ruleType string, audience string) (rule Rule, err error) {

	trx, ok := transaction.(MobiTransaction)
	if !ok {
		return Rule{}, ErrWrongTransactionFormat
	}

	query := `
		SELECT ` + columns() + ` FROM rules
		where 	status = 'active' 
				and created_at + interval '1' month * use_period >= $1 
				and created_at <= $1
				order by f_num_nulls('rules') asc, 
				array_length(priority_tags, 1) desc 
	`

	var (
		rules []Rule
	)

	rows, err := r.db.Query(context.Background(), query, trx.CreatedAT)

	if err != nil {
		return Rule{}, err
	}

	defer rows.Close()

	for rows.Next() {
		var i Rule
		err = rows.Scan(fields(&i)...)
		if err != nil {
			return Rule{}, err
		}
		rules = append(rules, i)
	}

	for _, v := range rules {

		if (int(v.SourceAccMethod.Int32) == trx.SourceAccountMethod || !v.SourceAccMethod.Valid) &&
			(int(v.SourceAccGate.Int32) == trx.SourceAccountGate || !v.SourceAccGate.Valid) &&
			(int(v.DestAccMethod.Int32) == trx.DestAccountMethod || !v.DestAccMethod.Valid) &&
			(int(v.DestAccGate.Int32) == trx.DestAccountGate || !v.DestAccGate.Valid) &&
			(int(v.LocationID.Int32) == trx.UserLocationID || !v.LocationID.Valid) &&
			(int(v.OperationType.Int32) == trx.OperationType || !v.OperationType.Valid) &&
			(int(v.MinAmount.Int32) <= trx.Amount && int(v.MaxAmount.Int32) >= trx.Amount) &&
			(v.TerminalID.String == trx.TerminalID || !v.TerminalID.Valid) &&
			(int(v.ProviderID.Int32) == trx.ProviderID || !v.ProviderID.Valid) &&
			checkTagFunctions(v, trx) &&
			(!v.Audience.Valid || v.Audience.String == audience || audience == "") &&
			(!v.Type.Valid || v.Type.String == ruleType || ruleType == "") {
			return v, nil
		}
	}

	return Rule{}, ErrNoRules
}

// TODO:: get by filter
func (r *RuleModule) GetList(filter RuleRequest) (rules []Rule, count int, err error) {

	var (
		query       = "SELECT " + columns() + " FROM rules WHERE 1 = 1 "
		queryCount  = "SELECT count(*) FROM rules WHERE 1 = 1 "
		len         = 0
		lenCount    = 0
		params      = make([]interface{}, 0)
		paramsCount = make([]interface{}, 0)
	)

	r.requestFiltering(filter, &query, &len, &params)

	rows, err := r.db.Query(context.Background(), query, params...)

	if err != nil {
		return nil, 0, err
	}

	defer rows.Close()

	for rows.Next() {
		var i Rule
		err := rows.Scan(fields(&i)...)
		if err != nil {
			return nil, 0, err
		}
		rules = append(rules, i)
	}

	// get count of rows for front
	r.requestFiltering(filter, &queryCount, &lenCount, &paramsCount)

	err = r.db.QueryRow(context.Background(), queryCount, paramsCount...).Scan(&count)
	if err != nil {
		return nil, 0, err
	}

	return rules, count, nil
}

func (r *RuleModule) GetAll(filters ...string) (rules []Rule, err error) {
	return r.getByFilter(filters...)
}

// when user creates rule this method will check for intersection with other rules
// if there is some rule that intersects with current will return error
func (r *RuleModule) RuleValidation(rule *Rule) (ruleSubset []Rule, err error) {

	if rule.MinAmount.Int32 > rule.MaxAmount.Int32 {
		return nil, ErrMinMaxRate
	}

	sort.Strings(rule.PriorityTags)

	var (
		count int = 0
	)

	queryCheck := `
		SELECT count(*) FROM rules WHERE 
			type = $1 
			AND source_acc_method = $2 
			AND source_acc_gate = $3 
			AND dest_acc_method = $4
			AND dest_acc_gate = $5 
			AND operation_type = $6
			AND location_id = $7
			AND audience = $8
			AND priority_tags = $9
			AND (($10 >= min_amount and  $10 <= max_amount) 
					or ($11 >= min_amount and $11 <= max_amount))
			AND status = 'active'
	`

	var (
		SourceAccountMethod interface{}
		SourceAccGate       interface{}
		DestAccountMethod   interface{}
		DestAccountGate     interface{}
		OpertaionType       interface{}
		LocationID          interface{}
		Audience            interface{}
	)

	if rule.SourceAccMethod.Valid {
		SourceAccountMethod = rule.SourceAccMethod.Int32
	}
	if rule.SourceAccGate.Valid {
		SourceAccGate = rule.SourceAccGate.Int32
	}

	if rule.DestAccMethod.Valid {
		DestAccountMethod = rule.DestAccMethod.Int32
	}

	if rule.DestAccGate.Valid {
		DestAccountGate = rule.DestAccGate.Int32
	}

	if rule.OperationType.Valid {
		OpertaionType = rule.OperationType.Int32
	}

	if rule.LocationID.Valid {
		LocationID = rule.LocationID.Int32
	}
	if rule.Audience.Valid {
		Audience = rule.Audience.String
	}

	err = r.db.QueryRow(context.Background(), queryCheck,
		rule.Type.String,
		SourceAccountMethod,
		SourceAccGate,
		DestAccountMethod,
		DestAccountGate,
		OpertaionType,
		LocationID,
		Audience,
		rule.PriorityTags,
		rule.MinAmount.Int32,
		rule.MaxAmount.Int32,
	).Scan(&count)

	if err != nil {
		return nil, err
	}

	if count > 0 {
		return nil, ErrRuleIntersection
	}

	queryGetAll := `SELECT ` + columns() + `FROM rules where status = 'active'`

	rows, err := r.db.Query(context.Background(), queryGetAll)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var i Rule

		err := rows.Scan(fields(&i)...)

		if err != nil {
			return nil, err
		}

		if (rule.Type.String == i.Type.String || !i.Type.Valid || !rule.Type.Valid) &&
			(rule.SourceAccGate.Int32 == i.SourceAccGate.Int32 || !i.SourceAccGate.Valid || !rule.SourceAccGate.Valid) &&
			(rule.SourceAccMethod.Int32 == i.SourceAccMethod.Int32 || !i.SourceAccMethod.Valid || !rule.SourceAccMethod.Valid) &&
			(rule.DestAccGate.Int32 == i.DestAccGate.Int32 || !i.DestAccGate.Valid || !rule.DestAccGate.Valid) &&
			(rule.DestAccMethod.Int32 == i.DestAccMethod.Int32 || !i.DestAccMethod.Valid || !rule.DestAccMethod.Valid) &&
			(rule.OperationType.Int32 == i.OperationType.Int32 || !i.OperationType.Valid || !rule.OperationType.Valid) &&
			(rule.LocationID.Int32 == i.LocationID.Int32 || !i.LocationID.Valid || !rule.LocationID.Valid) &&
			(rule.Audience.String == i.Audience.String || !i.Audience.Valid || !rule.Audience.Valid) {

			ruleSubset = append(ruleSubset, *rule)
		}
	}

	if len(ruleSubset) != 0 {
		return ruleSubset, ErrMoreThan1RuleInSubset
	}

	return nil, nil
}

func (r *RuleModule) Create(rule Rule) error {

	_, err := r.RuleValidation(&rule)

	if err != nil {
		return err
	}

	query := `INSERT INTO rules (type,
		source_acc_method,
		source_acc_gate,
		dest_acc_method,
		dest_acc_gate,
		operation_type,
		location_id,
		award_type,
		award_rate,
		min_amount,
		max_amount,
		up_limit,
		audience,
		condition_text,
		status,
		use_period,
		created_at,
		updated_at,
		priority_tags,
		terminal_id,
		provider_id,
		transaction_type) 
									
				VALUES($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16,  $17, $18, $19, $20, $21, $22);`

	var fields []interface{}

	s := reflect.ValueOf(&rule).Elem()
	typeOfRule := s.Type()

	for i := 0; i < s.NumField(); i++ {
		f := s.Field(i)

		if typeOfRule.Field(i).Name == "ID" {
			continue
		}

		switch f.Type().String() {
		case "boruen.NullString":
			ok := (f.Interface()).(NullString)
			switch ok.Valid {
			case true:
				fields = append(fields, ok.String)
			default:
				fields = append(fields, nil)
			}
		case "boruen.NullInt32":
			val := (f.Interface()).(NullInt32)
			switch val.Valid {
			case true:
				fields = append(fields, val.Int32)
			default:
				fields = append(fields, nil)
			}

		default:
			fields = append(fields, f.Interface())
		}
	}

	_, err = r.db.Exec(context.Background(), query, fields...)

	return err
}

// delete by id
func (r *RuleModule) Delete(id int) error {

	query := `DELETE FROM rules WHERE id = $1`

	_, err := r.db.Exec(context.Background(), query, id)

	if err != nil {
		return err
	}

	return err
}

func (r *RuleModule) GetConditions() (conditions []string, err error) {

	query := "select coalesce(array_agg(condition_text),ARRAY[]::text[]) from rules;"

	err = r.db.QueryRow(context.Background(), query).Scan(pq.Array(&conditions))
	if err != nil {
		return nil, err
	}

	return conditions, nil
}

func (r *RuleModule) EditRule(ruleRequest Rule) error {

	query := `UPDATE rules SET  type = $1,
								source_acc_method = $2,
								source_acc_gate = $3,
								dest_acc_method = $4,
								dest_acc_gate = $5,
								operation_type = $6,
								location_id = $7,
								award_type = $8,
								award_rate = $9,
								min_amount = $10,
								max_amount = $11,
								up_limit = $12,
								audience = $13,
								condition_text = $14,
								status = $15,
								use_period = $16,
								created_at = $17,
								updated_at = $18,
								priority_tags = $19,
								terminal_id = $20,
								provider_id = $21,
								transaction_type = $22
			WHERE id = $23`

	var fields []interface{}

	elems := reflect.ValueOf(ruleRequest)
	typeOfRuleRequest := reflect.TypeOf(ruleRequest)

	for i := 0; i < elems.NumField(); i++ {
		field := elems.Field(i)
		tp := typeOfRuleRequest.Field(i)

		if tp.Name == "ID" {
			continue
		}
		switch field.Type().String() {
		case "boruen.NullString":
			if !field.Interface().(NullString).Valid {
				fields = append(fields, nil)
			} else {
				fields = append(fields, field.Interface().(NullString).String)
			}
		case "boruen.NullInt32":
			if !field.Interface().(NullInt32).Valid {
				fields = append(fields, nil)
			} else {
				fields = append(fields, field.Interface().(NullInt32).Int32)
			}
		default:
			fields = append(fields, field.Interface())
		}
	}

	fields = append(fields, ruleRequest.ID)

	_, err := r.db.Exec(context.Background(), query, fields...)

	if err != nil {
		return err
	}
	return nil
}

func (r RuleModule) getRuleIds() (ids []int, err error) {
	query := `select id from rules`
	rows, err := r.db.Query(context.Background(), query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var i int
		err = rows.Scan(&i)
		if err != nil {
			return nil, err
		}
		ids = append(ids, i)
	}

	return ids, nil
}

func (r *RuleModule) getRuleSubSets() (RuleSubSets [][]Rule, err error) {
	var ()
	return nil, err
}

func (r *RuleModule) AddRuleTag(request TagRequest) error {
	query := `UPDATE rules SET priority_tags = array_append(priority_tags, $2) WHERE id = $1`

	exec, err := r.db.Exec(context.Background(), query, request.ID, request.Tag)
	if err != nil {
		return err
	}

	if exec.RowsAffected() == 0 {
		return NoRowsAffected
	}

	return nil
}

func (r *RuleModule) RemoveRuleTag(request TagRequest) error {
	query := `UPDATE rules SET priority_tags = array_remove(priority_tags, $2) WHERE id = $1`

	exec, err := r.db.Exec(context.Background(), query, request.ID, request.Tag)
	if err != nil {
		return err
	}

	if exec.RowsAffected() == 0 {
		return NoRowsAffected
	}

	return nil
}
