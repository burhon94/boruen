# Boruen - abbreviation for Bonus Rule Engine. 
The main goal of package is to analyze transaction and do some actions like accural bonus for bonus-transaction.


Dont forget that the package is in private repo and you have to import private package.

## Database requiremens

There is some migration files that contain sql queries for creating table and function that is need for boruen package to work.

For manipulating with database migrations you have to use [go-migrate](https://github.com/golang-migrate/migrate) package.

### Rule table

```go-migrate create -ext sql -dir path/to/migrations/dir -seq create_table_boruen.rules```

And past in ```######_craete_table_boruen.rules.up``` folowong code:

```sql
CREATE TABLE IF NOT EXISTS rules
(
    id                BIGSERIAL PRIMARY KEY,
    type              VARCHAR(30),
    transaction_type VARCHAR(30),
    source_acc_method BIGINT ,
    source_acc_gate   BIGINT,
    dest_acc_method   BIGINT,
    dest_acc_gate     BIGINT,
    operation_type    BIGINT,
    location_id       BIGINT,
--     country          VARCHAR(30) NOT NULL,
    award_type        VARCHAR(20)      NOT NULL,
    award_rate        BIGINT           NOT NULL,
    min_amount        BIGINT,
    max_amount        BIGINT,
    up_limit          BIGINT,
    audience          VARCHAR(30)      ,

    provider_id       bigint  ,
    terminal_id       varchar(200) ,

    status            varchar(30)      NOT NULL,

    use_period        BIGINT           NOT NULL,

    priority_tags     text[],

    condition_text    VARCHAR(500),
    created_at        TIMESTAMP        NOT NULL,
    updated_at        TIMESTAMP        NOT NULL
);
```

And in ```######_craete_table_boruen.rules.down```: 

```sql
DROP TABLE IF EXISTS rules;
```

### Num of Null columns postgres function:

Like aboce described? you have to use go-migrate to create function postgres

Function: 

```sql
CREATE OR REPLACE FUNCTION f_num_nulls(_tbl regclass)
  RETURNS SETOF int AS
$func$
BEGIN
   RETURN QUERY EXECUTE format(
      'SELECT num_nulls(%s) FROM %s'
    , (SELECT string_agg(quote_ident(attname), ', ')  -- column list
       FROM   pg_attribute
       WHERE  attrelid = _tbl
       AND    NOT attisdropped    -- no dropped (dead) columns
       AND    attnum > 0)         -- no system columns
    , _tbl
   );
END
$func$  LANGUAGE plpgsql;
```

## Usage

### First of all get ruleModule that is main engine to work with rules:

```go
import boruen

func main() {

  dns := "postgres://user:password@host:port/database"
  ruleModule, err := boruen.NewRuleEngine(dns)
  
}

``` 

## FindRuleForTransaction 

- transaction - is transaction struct object that was received
- ruleType variable depended on servcie property that can be used for recognizing purposes (for example in referral service monthly/momentaly)
- audience - wich user type can get bonus (example referral service referrer/referral)

```go
rule , err := ruleModule.FindRuleForTransaction(transaction, ruleType, audience)
if err != nil {
  return err
}
```

## Administration monitogring and manipulating with rules

Add Rule form Admin UI : 

```go 
    func (this *boruenManager) CreateRule(w http.ResponseWriter, r *http.Request) {

    var (
        response reply.Response
        request  boruen.Rule
        ctx      = r.Context()
    )

    span, ctx := this.tracer.StartSpanFromContext(ctx, "CreateRule")
    defer span.Finish()

    defer response.WriteJson(&w)

    err := httpprocessor.GetBodyJson(r, &request)
    if err != nil {
        this.logger.ErrorWithSpan(span, "cab not parse request body err", zap.Error(err))
        response = reply.BadRequest
        return
    }

    err = this.ruleModule.Create(request)
    if err != nil {
        this.logger.ErrorWithSpan(span, "can not create rule", zap.Error(err))
        this.sentry.CaptureException(err)
        response = reply.InternalError
        return
    }

    response = reply.OkResponse
    return
}
```
## Set of transactions

Note that for sum of transactions you have to generate one comon transaction by grouping all input transactions and call FindRuleForTransaction method with your ruleType property

__Example__:

(create some comont key for grouping. For transactions it can be sum of some properties of MobiTransaction struct)
```go

type Result struct {
    Rule boruen.Rule
    GroupedTransactions []MobiTransactions
}

func GroupTransactionsForFindingMonthlyRule(trx ...MobiTransactions) ([]Result, err error) {
    mp := make(mao[string][]MobiTransactions)
    res := make([]Result, 0)
    
    // group transactions
    
    for i := range trx {
        mp[strings.Itoa(trx.SourceAccountMethod) + strings.Itoa(trx.DestinationAccountMethod)] = 
            append(mp[strings.Itoa(trx.SourceAccountMethod) + strings.Itoa(trx.DestinationAccountMethod)], trx[i])
    }
    
    for i := range mp {
    
        sum := 0
        
        for j := range mp[i] {
            sum += mp[i][j].Amount
        }
        
        transaction := MobiTransaction {
            SourceAccountMethod : mp[i][0].SourceAccountMethod,
            DestinationAccountMethod: mp[i][0].DestinationAccountMethod,
        }
        
        rule, err := boruen.FindRuleForTransaction(transaction)
        
        if err != nil {
            return err
        }
        
        res = append(res, Result{
            rule,
            mp[i],
        })
    }
    
    return res
}
```
# boruen
