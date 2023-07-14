package boruen

import (
	"context"
	"reflect"
	"strconv"
)

func (r *RuleModule) getByFilter(FullWroteFilters ...string) (rules []Rule, err error) {

	query := "SELECT " + columns() + " FROM rules WHERE 1 = 1 "
	for _, v := range FullWroteFilters {
		query += v
	}

	rows, err := r.db.Query(context.Background(), query)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var i Rule

		err = rows.Scan(fields(&i)...)

		if err != nil {
			return nil, err
		}

		rules = append(rules, i)
	}

	return rules, nil
}

func (r *RuleModule) requestFiltering(filter RuleRequest, query *string, len *int, params *[]interface{}) {

	elems := reflect.ValueOf(&filter).Elem()
	typeOfFilter := reflect.TypeOf(filter)

	typeOfElems := elems.Type()
	for i := 0; i < elems.NumField(); i++ {
		field := elems.Field(i)
		typeOfField := typeOfFilter.Field(i)

		switch field.Type().String() {
		case "boruen.Nullstring":
			val := field.Interface().(NullString)
			if val.Valid {
				*len++
				*params = append(*params, val.String)
				if typeOfElems.Field(i).Name == "DateFrom" {
					*query += " AND created_at >= $" + strconv.Itoa(*len)
				} else if typeOfElems.Field(i).Name == "DateTo" {
					*query += " AND created_at < $" + strconv.Itoa(*len) + " + interval 1 'day'"
				} else {
					*query += " AND " + typeOfField.Tag.Get("db") + " = $" + strconv.Itoa(*len)
				}
			}
		case "boruen.NullInt32":
			val := field.Interface().(NullInt32)
			if val.Valid {
				*len++
				*params = append(*params, val.Int32)
				*query += " AND " + typeOfField.Tag.Get("db") + " = $" + strconv.Itoa(*len)
			}
		}
	}
}
