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
