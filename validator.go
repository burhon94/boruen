package boruen

import (
	"errors"
	"reflect"
	"strings"
)

func validate(s interface{}) error {

	tp := reflect.TypeOf(s)

	val := reflect.ValueOf(s)

	for i := 0; i < tp.NumField(); i++ {
		field := tp.Field(i)
		fieldVal := val.Field(i)

		err := validateField(field, fieldVal)
		if err != nil {
			return err
		}

	}

	return nil
}

func validateField(field reflect.StructField, val reflect.Value) error {

	if field.Tag.Get("validate") == "" {
		return nil
	}

	ruleBook := strings.Split(field.Tag.Get("validate"), " ")

	for _, v := range ruleBook {
		parsedRuleBook := strings.Split(v, ":")

		rule := parsedRuleBook[0]
		var options = []string{}

		if len(parsedRuleBook) > 2 {
			return errors.New("wrong foramt")
		}

		if len(parsedRuleBook) == 2 {
			options = strings.Split(parsedRuleBook[1], "|")
		}

		err := mapOfValidatorCheckers[rule](val, options...)

		if err != nil {
			return err
		}
	}

	return nil
}

var (
	mapOfValidatorCheckers = map[string]func(val reflect.Value, options ...string) error{
		"required": validatorRequiredChecker,
		"enum":     validatorEnumChecker,
	}
)

func validatorRequiredChecker(val reflect.Value, options ...string) error {
	if !val.Interface().(FieldGlobalInterface).GetValidField() {
		return errors.New("required but filed is null")
	}
	return nil
}

func validatorEnumChecker(val reflect.Value, options ...string) error {
	if !val.Interface().(FieldGlobalInterface).GetValidField() {
		return nil
	}

	for _, v := range options {
		if v == val.Interface().(FieldGlobalInterface).GetValue() {
			return nil
		}
	}

	return errors.New("not found value in enum")
}
