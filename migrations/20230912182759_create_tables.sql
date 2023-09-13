-- +goose Up
-- +goose StatementBegin
create table account (
    id serial primary key
);

create table schema (
    name text unique not null,
    id serial primary key
);

create table provider (
    id varchar(2) primary key,
    name text not null check ( name ~ '^[A-Z]{2}$' )
);

create table airline (
    code varchar(2) primary key check ( code ~ '^[0-9A-Z\u0401\u0451\u0410-\u044f]{2}$' ),
    name text not null
);

create table account_schema (
    account_id integer references account(id) on delete cascade primary key,
    schema_id integer references schema(id) on delete cascade
);

create table schema_provider (
    schema_id integer references schema(id) on delete cascade,
    provider_id varchar(2) references provider(id) on delete cascade
);

create table provider_airline (
    provider_id varchar(2) references provider(id),
    airline_code varchar(2) references airline(code)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table
    account,
    account_schema,
    schema,
    schema_provider,
    provider,
    provider_airline,
    airline;
-- +goose StatementEnd
