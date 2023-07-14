package boruen

import "errors"

var (
	ErrNoRules                = errors.New("no rules found for transaction")
	ErrMoreThan1RuleInSubset  = errors.New("warning : more than 1 rule in the subset of rules")
	ErrWrongTransactionFormat = errors.New("wrong transaction format")
	ErrMinMaxRate             = errors.New("min_amount is greater than max_amount")
	ErrRuleIntersection       = errors.New("there is exists rule with such parametres")
	NoRowsAffected            = errors.New("no rows affected")
)
