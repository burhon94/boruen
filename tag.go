package boruen

import (
	"errors"
)

var (
	ErrKeyAlreadyInUse = errors.New("key is already in use")
)

type TagFunction func(rule Rule, trx interface{}) bool

var mapOfTags = make(map[string]TagFunction)

func AddTagFunction(key string, tagFunction TagFunction) (err error) {
	if _, ok := mapOfTags[key]; ok {
		return ErrKeyAlreadyInUse
	}
	mapOfTags[key] = tagFunction
	return nil
}

func UpdateTagFunction(key string, newTagFunction TagFunction) {
	mapOfTags[key] = newTagFunction
}

func GetTagFunctionsList() map[string]TagFunction {
	return mapOfTags
}

func checkTagFunctions(rule Rule, trx MobiTransaction) bool {
	if len(rule.PriorityTags) == 0 {
		return true
	}

	for _, v := range rule.PriorityTags {
		if !mapOfTags[v](rule, trx) {
			return false
		}
	}
	return true
}
