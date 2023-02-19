CREATE TABLE tokens (
    id serial not null primary key,
    name varchar(255) not null,
    symbol varchar(255) not null,
    address varchar(255) not null,
    decimals varchar(255) not null,
    total_supply varchar(255) not null
);