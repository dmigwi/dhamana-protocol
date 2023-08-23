// Copyright (c) 2023 Migwi Ndung'u
// See LICENSE for details.

package storage

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/dmigwi/dhamana-protocol/client/utils"
	_ "github.com/lib/pq"  // postgres
	_ "modernc.org/sqlite" // sqlite
)

const (
	// PostgresDriverName defines the postgres driver name.
	PostgresDriverName = "postgres"

	// SqliteDriverName defines the sqlite driver name.
	SqliteDriverName = "sqlite"

	// createTableBond is an sql statement creating a table with the name table_bond.
	// It creates the table if it doesn't exists.
	createTableBond = "CREATE TABLE IF NOT EXISTS table_bond (" +
		"bond_address VARCHAR(42) PRIMARY KEY," +
		"issuer_address VARCHAR(42)," +
		"holder_address VARCHAR(42)," +
		"created_time TIMESTAMPTZ," +
		"tx_hash VARCHAR(66)," +
		"created_block INTEGER," +
		"principal INTEGER," +
		"coupon_rate SMALLINT CHECK (coupon_rate BETWEEN 1 AND 100)," +
		"coupon_date TIMESTAMPTZ," +
		"maturity_date TIMESTAMPTZ," +
		"currency SMALLINT CHECK (currency BETWEEN 0 AND 50)," +
		"intro_msg TEXT," +
		"last_status SMALLINT CHECK (last_status BETWEEN 0 AND 10)," +
		"last_update TIMESTAMPTZ)"

	// createTableBondStatus is an sql statement creating a table with the name table_bond_status.
	// It creates the table if it doesn't exists.
	createTableBondStatus = "CREATE TABLE IF NOT EXISTS table_bond_status (" +
		"id SERIAL PRIMARY KEY," +
		"bond_status SMALLINT NOT NULL CHECK(bond_status BETWEEN 0 AND 10)," +
		"bond_address VARCHAR(42)," +
		"issuer_signed BOOLEAN," +
		"holder_signed BOOLEAN," +
		"update_time TIMESTAMPTZ)"

	// createChatTable is an sql statement creating a table with the name table_chat.
	// It creates the table if it doesn't exists.
	createChatTable = "CREATE TABLE IF NOT EXISTS table_chat (" +
		"id SERIAL PRIMARY KEY," +
		"sender VARCHAR(42)," +
		"chat_msg TEXT," +
		"recieved TIMESTAMPTZ)"

	// fetchBonds is an sql query that fetches all the bonds that are owned
	// by bond party or they are still in the negotiation stage.
	fetchBonds = "SELECT (bond_address,created_time,currency,last_status)" +
		"FROM table_bond WHERE issuer_address = ? OR last_status = 0 " +
		"Or holder_address = ? ORDER BY 'last_update' DESC LIMIT = ?"

	// fetchBondByAddress is an sql statement that returns a bond identified by
	// the provided address if the sender is a party to the bond or the bond
	// is still in the negotiation stage.
	fetchBondByAddress = "SELECT (bond_address,issuer_address,holder_address," +
		"created_time,tx_hash,created_block,principal,coupon_rate,coupon_date," +
		"maturity_date,currency,intro_msg,last_status,last_update)" +
		"FROM table_bond WHERE bond_address = ? AND (last_status = 0 OR " +
		"issuer_address = ? OR holder_address = ?) " +
		"ORDER BY 'last_update'"
)

// tablesToSQLStmt is an array of sql statements used to create the missing tables
// if they don't exist.
var tablesSQLStmt = []string{
	createTableBond,
	createTableBondStatus,
	createChatTable,
}

// reqToStmt matches the respective Methods supported to thier sql queries.
var reqToStmt = map[utils.Method]string{
	utils.GetBonds:         fetchBonds,
	utils.GetBondByAddress: fetchBondByAddress,
}

// DB defines the parameters needed to use a persistence db instance connect to.
type DB struct {
	*sql.DB
	ctx context.Context
}

// reader defines the method that read the sql row into the required interface
// type returned.
type reader interface {
	Read(row *sql.Row) interface{}
}

// NewDB returns an opened db instance whose connection has been tested with
// ping request. The driverName is required for specifying which db type to use.
// It generates the required tables if they don't exist.
func NewDB(ctx context.Context, port uint16,
	driverName, host, user, password, dbname string,
) (*DB, error) {
	connInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	db, err := sql.Open(driverName, connInfo)
	if err != nil {
		log.Errorf("unable to open to %s: err %v", driverName, err)
		return nil, err
	}

	if err = db.PingContext(ctx); err != nil {
		log.Errorf("connection to %s failed: err %v", driverName, err)
		return nil, err
	}

	for i, stmt := range tablesSQLStmt {
		if _, err := db.ExecContext(ctx, stmt); err != nil {
			log.Errorf("creating a table index (%d) failed Error: %v", i, err)
			return nil, err
		}
	}

	return &DB{
		DB:  db,
		ctx: ctx,
	}, nil
}

func (db *DB) QueryData(method utils.Method, r reader, params ...interface{}) ([]interface{}, error) {
	// meth
	return nil, nil
}
