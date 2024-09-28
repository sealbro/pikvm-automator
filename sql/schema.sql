CREATE TABLE commands (
    created     timestamp       DEFAULT CURRENT_TIMESTAMP,
    updated     timestamp       DEFAULT CURRENT_TIMESTAMP,
    id          varchar(100)    PRIMARY KEY,
    description varchar(200)    NOT NULL,
    expression  varchar(500)    NOT NULL
);