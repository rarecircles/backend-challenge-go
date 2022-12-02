package sql

var CreateTableTokens = `
CREATE TABLE IF NOT EXISTS
 tokens(
 	id serial primary key,
 	address varchar(255) unique,
 	name varchar(255),
 	symbol varchar(30),
	decimals int,
 	total_supply numeric
 )
`

// Note: Could change ID serial to uuid, snowflake, or simply rely on address as primary key.
