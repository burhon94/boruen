package boruen

import (
	"math/rand"
	"strconv"
)

func ruleGenerator() Rule {

	var tags []string

	for k, _ := range mapOfTags {
		tags = append(tags, k)
	}

	var RandomPriorityTags []string

	for i := 0; i < rand.Int()%3; i++ {
		RandomPriorityTags = append(RandomPriorityTags, tags[i])
	}

	rule := Rule{
		Type:            GenNullString(typeItems[rand.Int()%3]),
		SourceAccMethod: GenNullInt32(rand.Int() % 3),
		SourceAccGate:   GenNullInt32(rand.Int() % 3),
		DestAccMethod:   GenNullInt32(rand.Int() % 3),
		DestAccGate:     GenNullInt32(rand.Int() % 3),
		OperationType:   GenNullInt32(rand.Int() % 3),
		LocationID:      GenNullInt32(rand.Int() % 3),
		Audience:        GenNullString(audience[rand.Int()%4]),
		ProviderID:      GenNullInt32(rand.Int() % 2),
		TerminalID:      GenNullString(GenString(rand.Int() % 3)),
		TransactionType: GenNullString(transactionType[rand.Int()%2]),
		AwardRate:       GenNullInt32(rand.Int() % 1000 * 100),
		AwardType:       GenNullString(awardType[rand.Int()%2]),
		MinAmount:       GenNullInt32(rand.Int() % 100000 * 100),
		MaxAmount:       GenNullInt32(rand.Int() % 100000 * 100),
		UppLimit:        GenNullInt32(rand.Int() % 100000 * 100),
		ConditionText:   GenString(rand.Int()%18 + 10),
		PriorityTags:    RandomPriorityTags,
		Status:          GenNullString("active"), //RuleStatus[rand.Int()%2],
		UsePeriod:       GenNullInt32(rand.Int()%14 + 5),
		CreatedAT:       "now()", //"now() + " + strconv.Itoa(rand.Int()%1200) + " * interval '1' hour",
		UpdatedAT:       "now()", //"(now() + interval '1' hour * mod( (random() * 100)::integer, 12))::timestamp ",
	}
	if rule.MinAmount.Int32 > rule.MaxAmount.Int32 {
		rule.MinAmount, rule.MaxAmount = rule.MaxAmount, rule.MinAmount
	}

	return rule
}

func genTrx() MobiTransaction {
	TransactionType := map[int]string{
		0: "income",
		1: "expense",
	}
	return MobiTransaction{
		ExternalID:              rand.Int(),
		UserID:                  GenString(10),
		UserLocationID:          rand.Int() % 2,
		SourceAccount:           GenString(rand.Int() % 2),
		SourceAccountMethod:     rand.Int() % 2,
		SourceAccountGate:       rand.Int() % 2,
		SourceAccountMethodType: GenString(rand.Int() % 2),
		DestAccount:             GenString(rand.Int()%2 + 1),
		DestAccountGate:         rand.Int() % 2,
		DestAccountMethod:       rand.Int() % 2,
		DestAccountMethodType:   GenString(rand.Int() % 2),
		ProviderID:              rand.Int() % 2,
		TerminalID:              GenString(20),
		Type:                    TransactionType[rand.Int()%2+1],
		Processed:               rand.Int()%2 == 1,
		Amount:                  (rand.Int() % 10000) * 100,
		Commission:              rand.Int() % 100,
		Comment:                 GenString(20),
		OperationType:           rand.Int() % 2,
		MobiStatus:              "accepted",
		ServiceStatus:           "approved",
		ExtraInfo:               GenString(30),
		CreatedAT:               "now()", //"now() + " + strconv.Itoa(rand.Int()%12) + " * interval '1' hour",
		UpdatedAT:               "now() + " + strconv.Itoa(rand.Int()%12) + " * interval '1' hour",
	}

}

var glbdin = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890"

func GenString(ln int) (s string) {
	for i := 0; i < ln; i++ {
		s = s + string(glbdin[rand.Int()%len(glbdin)])
	}
	return s
}

// generate datas
func GenNullString(data string) NullString {
	if data == "" {
		return NullString{
			Valid:  false,
			String: "",
		}
	}

	return NullString{
		Valid:  true,
		String: data,
	}
}

func GenNullInt32(data int) NullInt32 {
	if data == -1 {
		return NullInt32{
			Valid: false,
			Int32: 0,
		}
	}
	return NullInt32{
		Valid: true,
		Int32: int32(data),
	}
}

// "firstInAllHistory":     isFirstTag,
// "firstInHistoryForRule": isFirstForRule,
// }

// write sql qeuery checkers
func isFirstTag(rule Rule, trx interface{}) (ok bool) {
	ok = true
	return ok
}

func isFirstForRule(rule Rule, trx interface{}) (ok bool) {
	ok = true
	// use Rule
	return ok
}
