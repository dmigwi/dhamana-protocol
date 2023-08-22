// Copyright (c) 2023 Migwi Ndung'u
// See LICENSE for details.

package storage

const (
	createDB        = "CREATE DATABASE %s;"
	createTableBond = `CREATE TABLE table_bond (
		bond_address VARCHAR(42) PRIMARY KEY,
		issuer_address VARCHAR(42),
		holder_address VARCHAR(42),
		created_time TIMESTAMP,
		tx_hash VARCHAR(66),
		created_block INTEGER,
		principal INTEGER,
		coupon_rate SMALLINT CHECK (coupon_rate BETWEEN 1 AND 100),
		coupon_date TIMESTAMP,
		maturity_date TIMESTAMP,
		currency SMALLINT CHECK (currency BETWEEN 0 AND 50),,
		intro_msg TEXT,
		last_status SMALLINT CHECK (last_status BETWEEN 0 AND 10),
		last_update TIMESTAMP,
	);`
	createTableBondStatus = `CREATE TABLE table_bond_status (
		id SERIAL PRIMARY KEY,
		bond_status SMALLINT NOT NULL (bond_status BETWEEN 0 AND 10),
		bond_address VARCHAR(42),
		issuer_signed BOOLEAN,
		holder_signed BOOLEAN,
		update_time TIMESTAMP
	);`
	createChatTable = `CREATE TABLE table_chat (
		id SERIAL PRIMARY KEY,
		sender VARCHAR(42),
		chat_msg TEXT,
		recieved TIMESTAMP
	);`
	alterUserPassword  = "ALTER USER %s WITH ENCRYPTED password %s"
	grantAllPrivileges = "GRANT ALL PRIVILEGES on DATABASE %s TO %s;"
)
