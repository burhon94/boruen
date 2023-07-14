package boruen

import (
	"fmt"
	"math/rand"
	"strconv"
	"testing"

	"github.com/alifcapital/easy"
	"github.com/jackc/pgx/v4"
)

var (
	TestRuleModule RuleModule
)

func TestMain(m *testing.M) {

	AddTagFunction("firstInHistory", isFirstTag)
	AddTagFunction("firstInHisotyrForRule", isFirstForRule)

	easy.New("file://./db")
	defer easy.Close()

	TestRuleModule = NewRuleEngine(fmt.Sprintf("postgres://mobi:test@localhost:%s/%s?sslmode=disable", "1234", "mobidb"))

	m.Run()
}

func TestCreate(t *testing.T) {
	for i := 0; i < rand.Int()%1000+11; i++ {
		err := TestRuleModule.Create(ruleGenerator())
		if err != nil && err != ErrMoreThan1RuleInSubset {
			t.Error("can not create rule", err)
			t.Fail()
		}
	}
}

// var Filters = map[string]string {
// 	"ById" :
// }

// add to easy package restart migrations!!!
// add to easy clean up tables

func TestGetRuleIds(t *testing.T) {
	_, err := TestRuleModule.getRuleIds()
	if err != nil && err != pgx.ErrNoRows {
		t.Error("can not get ids", err)
		t.Fail()
	}
}

func TestGetAll(t *testing.T) {

	//state1 GetById

	var rule Rule
	for i := 0; i < rand.Int()%1000+1; i++ {
		rule = ruleGenerator()
		err := TestRuleModule.Create(rule)
		if err != nil && err != ErrMoreThan1RuleInSubset {
			t.Error("can not create rule [TestGetAll]", err)
			t.Fail()
		}
	}

	ids, err := TestRuleModule.getRuleIds()
	if err != nil {
		t.Error("can not get ids", err)
		t.Fail()
	}

	for _, id := range ids {
		_, err := TestRuleModule.GetAll(ById(strconv.Itoa(id)))
		if err != nil {
			t.Error("can not get rule by id", err)
			t.Fail()
		}
		// fmt.Println("rules by id", rule)
	}
}

func TestFindRuleForTransaction(t *testing.T) {
	transaction := genTrx()

	_, err := TestRuleModule.FindRuleForTransaction(transaction, typeItems[rand.Int()%2], audience[rand.Int()%3])
	if err != nil && err != ErrNoRules {
		t.Error("error getting rule for trx", err)
		t.Fail()
	}
	//fmt.Println("transaction \n", transaction, "\n rule for transaction \n", rule, "\n", err, " TestFindRuleForTransaction")
}

func TestSpecialCasesForRulePerTransaction(t *testing.T) {
	rules, err := TestRuleModule.getByFilter()
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	for _, v := range rules {
		err = TestRuleModule.Delete(v.ID)
		if err != nil {
			t.Error(err)
			t.FailNow()
		}
	}

	rule := ruleGenerator()

	transaction := genTrx()

	rule.Audience = GenNullString("")
	rule.Type = GenNullString("")
	rule.SourceAccGate = GenNullInt32(-1)
	rule.SourceAccMethod = GenNullInt32(-1)
	rule.DestAccGate = GenNullInt32(-1)
	rule.DestAccMethod = GenNullInt32(-1)
	rule.OperationType = GenNullInt32(-1)
	rule.LocationID = GenNullInt32(-1)
	rule.MinAmount = GenNullInt32(0)
	rule.MaxAmount = GenNullInt32(100000000)
	rule.TerminalID = GenNullString("")
	rule.ProviderID = GenNullInt32(-1)
	rule.UsePeriod = GenNullInt32(10)
	rule.AwardType = GenNullString("percent")
	rule.Status = GenNullString("active")

	err = TestRuleModule.Create(rule)

	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	_, err = TestRuleModule.getByFilter()

	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	_, err = TestRuleModule.FindRuleForTransaction(transaction, rule.Type.String, rule.Audience.String)

	if err != nil {
		t.Error(err)
		t.FailNow()
	}
}

func TestValidatorRequiredType(t *testing.T) {

	var rule = ruleGenerator()
	rule.MinAmount.Valid = false

	err := validate(rule)
	if err == nil {
		t.Error(err)
		t.FailNow()
	}
}

func TestValidatorEnumType(t *testing.T) {
	var rule = ruleGenerator()
	rule.Status.String = "wrongType"

	err := validate(rule)

	if err == nil {
		t.Error(err)
		t.FailNow()
	}
}

func TestRuleForPassingValidator(t *testing.T) {
	var rule = ruleGenerator()
	err := validate(rule)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
}

func TestEdit(t *testing.T) {
	rule := ruleGenerator()
	if err := TestRuleModule.Create(rule); err != nil {
		t.Error(err)
		t.FailNow()
	}
	id, err := TestRuleModule.getRuleIds()
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	rules, count, err := TestRuleModule.GetList(RuleRequest{ID: NullInt32{Int32: int32(id[0]), Valid: true}})
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	if count != 1 {
		t.Error("logical error: there must be 1 rule")
		t.FailNow()
	}
	prevRule := rules[0]
	rules[0].Audience = NullString{String: "new", Valid: true}

	err = TestRuleModule.EditRule(rules[0])
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	rules, _, _ = TestRuleModule.GetList(RuleRequest{ID: NullInt32{Int32: int32(id[0]), Valid: true}})

	if rules[0].Audience.String == prevRule.Audience.String {
		t.Error("data was not edited")
		t.FailNow()
	}
}
