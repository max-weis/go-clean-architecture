BEGIN;

create table products
(
    id          VARCHAR(36) PRIMARY KEY,
    title       VARCHAR UNIQUE NOT NULL,
    description VARCHAR,
    price       integer        NOT NULL,
    created_at  timestamp,
    modified_at timestamp
);

COMMIT;