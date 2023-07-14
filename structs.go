package boruen

type Filter func(string) string

type Rule struct {
	ID int `json:"id"`

	Type NullString `json:"type"` //

	// template

	// TransactionType string `json:"transaction_type" validate:"enum:income|expense required"`
	//Country         string `json:"country"`
	SourceAccMethod NullInt32 `json:"sourceAccMethod"`
	SourceAccGate   NullInt32 `json:"sourceAccGate"`

	DestAccMethod NullInt32 `json:"destAccMethod"`
	DestAccGate   NullInt32 `json:"destAccGate"`

	OperationType NullInt32 `json:"operationType"`

	LocationID NullInt32 `json:"locationId"`

	// award and period options for template
	AwardType NullString `json:"awardType" validate:"required"` //
	AwardRate NullInt32  `json:"awardRate" validate:"required"` //

	MinAmount NullInt32 `json:"minAmount" validate:"required"` //
	MaxAmount NullInt32 `json:"maxAmount" validate:"required"` //

	UppLimit NullInt32 `json:"upLimit" validate:"required"` //

	Audience NullString `json:"audience"` //

	ConditionText string `json:"conditionText"` //

	Status NullString `json:"status" validate:"required enum:active|inactive"` // active, inactive

	UsePeriod NullInt32 `json:"usePeriod" validate:"required"` //

	CreatedAT string `json:"createdAT"` //
	UpdatedAT string `json:"updatedAT"` //

	PriorityTags []string `json:"priorityTags"` //

	TerminalID      NullString `json:"terminalId"`                                              //
	ProviderID      NullInt32  `json:"providerId"`                                              //
	TransactionType NullString `json:"transactionType" validate:"required enum:income|expense"` //

}

type RuleRequest struct {
	ID   NullInt32  `json:"id" db:"id"`
	Type NullString `json:"type" db:"type"`

	SourceAccMethod NullInt32 `json:"sourceAccMethod" db:"source_acc_method"`
	SourceAccGate   NullInt32 `json:"sourceAccGate" db:"source_acc_gate"`
	DestAccMethod   NullInt32 `json:"destAccMethod" db:"dest_acc_method"`
	DestAccGate     NullInt32 `json:"destAccGate" db:"dest_acc_gate"`
	OperationType   NullInt32 `json:"operationType" db:"operation_type"`

	LocationID NullInt32 `json:"location_id" db:"location_id"`

	AwardRate NullInt32  `json:"award_rate" db:"award_rate"`
	AwardType NullString `json:"award_type" db:"award_type"`

	MinAmount NullInt32 `json:"minAmount" db:"min_amount"`
	MaxAmount NullInt32 `json:"maxAmount" db:"max_amount"`
	UppLimit  NullInt32 `json:"upLimit" db:"up_limit"`

	Audience   NullString `json:"audience" db:"audience"`
	ProviderID NullInt32  `json:"provider_id" db:"provider_id"`

	ConditionText NullString `json:"conditionText" db:"condition_text"`

	Status NullString `json:"status" db:"status"` //active, inactive

	UsePeriod NullInt32 `json:"usePeriod" db:"use_period"`

	DateFrom NullString `json:"date_from" db:"date_from"`
	DateTo   NullString `json:"date_to" db:"date_to"`
}

type MobiTransaction struct {
	ID         int
	ExternalID int `json:"id"`

	UserID         string `json:"userID"`
	UserLocationID int    `json:"userLocationID"`

	SourceAccount           string `json:"sourceAccount"`
	SourceAccountName       string `json:"sourceAccountName"`
	SourceAccountMethod     int    `json:"sourceAccountMethod"`
	SourceAccountMethodType string `json:"sourceAccountMethodType"`
	SourceAccountGate       int    `json:"sourceAccountGate"`

	DestAccount           string `json:"destAccount"`     //StoreID,
	DestAccountName       string `json:"destAccountName"` //merchant, provider
	DestAccountMethod     int    `json:"destAccountMethod"`
	DestAccountMethodType string `json:"destAccountMethodType"`
	DestAccountGate       int    `json:"destAccountGate"`

	ProviderID int    `json:"providerID"`
	TerminalID string `json:"terminalID"`

	Type      string `json:"type"`
	Processed bool   `json:"processed"`

	Amount        int    `json:"amount"`
	Commission    int    `json:"commission"`
	Comment       string `json:"comment"`
	OperationType int    `json:"operationType"` //payment, topup
	MobiStatus    string `json:"status"`
	ServiceStatus string `json:"serviceStatus"`

	ExtraInfo string `json:"extraInfo"`
	CreatedAT string `json:"createdAT"`
	UpdatedAT string `json:"updatedAT"`
}

type TagRequest struct {
	ID  int    `json:"id"`
	Tag string `json:"tag"`
}
