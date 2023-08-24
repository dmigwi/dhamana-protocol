// Copyright (c) 2023 Migwi Ndung'u
// See LICENSE for details.

package storage

import (
	"context"
	"database/sql"
	"fmt"
	"regexp"

	"github.com/dmigwi/dhamana-protocol/client/utils"
	_ "github.com/lib/pq"  // postgres
	_ "modernc.org/sqlite" // sqlite
)

const (
	// createTableBond is a prepared statement creating a table identified with the
	// name table_bond if it doesn't exists.
	createTableBond = "CREATE TABLE IF NOT EXISTS table_bond (" +
		"bond_address VARCHAR(42) PRIMARY KEY," +
		"issuer_address VARCHAR(42)," +
		"holder_address VARCHAR(42)," +
		"created_at TIMESTAMPTZ," +
		"created_at_block INTEGER," +
		"principal INTEGER," +
		"coupon_rate SMALLINT CHECK (coupon_rate BETWEEN 1 AND 100)," +
		"coupon_date TIMESTAMPTZ," +
		"maturity_date TIMESTAMPTZ," +
		"currency SMALLINT CHECK (currency BETWEEN 0 AND 50)," +
		"intro_msg TEXT," +
		"last_status SMALLINT CHECK (last_status BETWEEN 0 AND 10)," +
		"last_update TIMESTAMPTZ," +
		"last_synced_block INTEGER)"

	// createTableBondStatus is a prepared statement creating a table identified
	// with the name table_status if it doesn't exists.
	createTableBondStatus = "CREATE TABLE IF NOT EXISTS table_status (" +
		"id SERIAL PRIMARY KEY," +
		"sender VARCHAR(42)," +
		"bond_address VARCHAR(42)," +
		"bond_status SMALLINT NOT NULL CHECK(bond_status BETWEEN 0 AND 10)," +
		"added_on TIMESTAMPTZ," +
		"last_synced_block INTEGER)"

	// createTableBondStatusSigned is a prepared statement creating a table
	// identified with the name table_status if it doesn't exists.
	createTableBondStatusSigned = "CREATE TABLE IF NOT EXISTS table_status_signed (" +
		"id SERIAL PRIMARY KEY," +
		"sender VARCHAR(42)," +
		"bond_address VARCHAR(42)," +
		"bond_status SMALLINT NOT NULL CHECK(bond_status BETWEEN 0 AND 10)," +
		"signed_on TIMESTAMPTZ," +
		"last_synced_block INTEGER)"

	// createChatTable is a prepared statement creating a table identified with
	// the name table_chat if it doesn't exists.
	createChatTable = "CREATE TABLE IF NOT EXISTS table_chat (" +
		"id SERIAL PRIMARY KEY," +
		"sender VARCHAR(42)," +
		"bond_address VARCHAR(42)," +
		"chat_msg TEXT," +
		"created_at TIMESTAMPTZ," +
		"last_synced_block INTEGER)"

	// fetchBonds is a prepared statement that fetches all the bonds that owned
	// by the bond party with the address or they are still in the negotiation stage.
	fetchBonds = "SELECT (bond_address,created_at,coupon_rate,currency,last_status)" +
		"FROM table_bond WHERE issuer_address = $1 OR last_status = 0 " +
		"Or holder_address = $2 ORDER BY last_update DESC LIMIT $3 OFFSET $4"

	// fetchBondByAddress is a prepared statement that returns a bond identified by
	// the provided address if the sender is a party to the bond or the bond
	// is still in the negotiation stage.
	fetchBondByAddress = "SELECT (bond_address,issuer_address,holder_address," +
		"created_at,created_at_block,principal,coupon_rate,coupon_date," +
		"maturity_date,currency,intro_msg,last_status,last_update,last_synced_block)" +
		"FROM table_bond WHERE bond_address = $1 AND " +
		"(last_status = 0 OR issuer_address = $2 OR holder_address = $3)"

	// fetchChats is a prepared statement that fetches the conversation within
	// the bond identified by the provided address if the sender is a bond party
	// or its still in the negotiation stage.
	fetchChats = "SELECT (c.sender, c.bond_address, c.chat_msg, c.created_at, c.last_synced_block)" +
		"FROM table_chat as c LEFT JOIN table_bond as b WHERE c.bond_address = b.bond_address AND " +
		"(b.issuer_address = $1 OR b.last_status = 0 Or b.holder_address = $2)" +
		" ORDER BY c.created_at DESC LIMIT $3 OFFSET $4"

	// fetchLastSyncBlock returns the last block to be synced on the table_bond.
	fetchLastSyncBlock = "SELECT last_synced_block FROM table_bond ORDER BY" +
		" last_synced_block DESC LIMIT 1"

	// setBondBodyTerms updates the table_bond with data from the BondBodyTerms event.
	setBondBodyTerms = "UPDATE table_bond SET principal = $1, coupon_rate = $2, " +
		"coupon_date = $3, maturity_date = $4, currency = $5, last_update = $6, " +
		"last_synced_block = $7 WHERE bond_address = $8"

	// setBondMotivation update the table_bond with data from BondMotivation event.
	setBondMotivation = "UPDATE table_bond SET intro_msg = $1, last_update = $2," +
		" last_synced_block = $3  WHERE bond_address = $4"

	// setHolder updates table_bond with data from HolderUpdate event.
	setHolder = "UPDATE table_bond SET holder_address = $1, last_update = $2, " +
		"last_synced_block = $3 WHERE bond_address = $4"

	// setLastStatus updates table_bond with data from StatusChange event.
	setLastStatus = "UPDATE table_bond SET last_status = $1, last_update = $2, " +
		"last_synced_block = $3 WHERE bond_address = $4"

	// addNewBondCreated inserts into table_bond new data from event NewBondCreated.
	addNewBondCreated = "INSERT INTO table_bond (bond_address, issuer_address, " +
		"last_update, last_synced_block) VALUES ($1, $2, $3, $4)"

	// ddNewChatMessage inserts into table_chat new data from event NewChatMessage.
	addNewChatMessage = "INSERT INTO table_chat (sender, bond_address, " +
		" chat_msg, created_at, last_synced_block) VALUES ($1, $2, $3, $4, $5)"

	// addStatusChange inserts into table_status new data from event StatusChange.
	addStatusChange = "INSERT INTO table_status (sender, bond_address, " +
		"bond_status, added_on, last_synced_block) VALUES ($1, $2, $3, $4, $5)"

	// addStatusSigned inserts into table_status_signed new data from event StatusSigned.
	addStatusSigned = "INSERT INTO table_status_signed (sender, bond_address, " +
		"bond_status, signed_on, last_synced_block) VALUES ($1, $2, $3, $4, $5)"

	dropTableBondRecords         = "DELETE * FROM table_bond WHERE last_synced_block = $1"
	dropTableStatusRecords       = "DELETE * FROM table_status WHERE last_synced_block = $1"
	dropTableStatusSignedRecords = "DELETE * FROM table_status_signed WHERE last_synced_block = $1"
	dropTableChatRecords         = "DELETE * FROM table_chat WHERE last_synced_block = $1"
)

// tablesToSQLStmt is an array of sql statements used to create the missing tables
// if they don't exist.
var tablesSQLStmt = []string{
	createTableBond,
	createTableBondStatus,
	createTableBondStatusSigned,
	createChatTable,
}

// This are clean up methods employed if corrupt or dirty writes are made at
// a certain last synced block.
var cleanUpStmt = []string{
	dropTableBondRecords,
	dropTableStatusRecords,
	dropTableStatusSignedRecords,
	dropTableChatRecords,
}

// reqToStmt matches the respective local type Methods supported to their sql queries.
var reqToStmt = map[utils.Method]string{
	utils.GetBonds:         fetchBonds,
	utils.GetBondByAddress: fetchBondByAddress,
	utils.GetChats:         fetchChats,

	// method needed locally. Results are not sent via the server
	utils.GetLastSyncedBlock: fetchLastSyncBlock,

	utils.UpdateBondBodyTerms:  setBondBodyTerms,
	utils.UpdateBondMotivation: setBondMotivation,
	utils.UpdateHolder:         setHolder,
	utils.UpdateLastStatus:     setLastStatus,
	utils.InsertNewBondCreated: addNewBondCreated,
	utils.InsertNewChatMessage: addNewChatMessage,
	utils.InsertStatusChange:   addStatusChange,
	utils.InsertStatusSigned:   addStatusSigned,
}

// DB defines the parameters needed to use a persistence db instance connect to.
type DB struct {
	db     *sql.DB
	ctx    context.Context
	driver string
}

// Reader defines the method that reads the row fields into the require data interface.
// To read data, pass pointers to the expect field the parameter function.
type Reader interface {
	Read(fn func(fields ...any) error) (interface{}, error)
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
		db:     db,
		ctx:    ctx,
		driver: driverName,
	}, nil
}

// QueryLocalData executes the sql statement associated with the provided local
// method and uses the reader interface provided to read the row data result set.
// It then returns an array of data for each row read successfully otherwise
// an error is returned.
func (d *DB) QueryLocalData(method utils.Method, r Reader, sender string,
	params ...interface{},
) ([]interface{}, error) {
	stmt, ok := reqToStmt[method]
	if !ok {
		return nil, fmt.Errorf("missing query for method %q", method)
	}

	switch method {
	case utils.GetBondByAddress, utils.GetBonds:
		// The sender's address is used 3 times as an argument for all queries.
		params = append(params, []interface{}{sender, sender, sender}...)
	}

	rows, err := d.db.QueryContext(d.ctx, d.formatPreparedStmt(stmt), params)
	if err != nil {
		return nil, fmt.Errorf("fetching query for method %q failed: %v", method, err)
	}

	var data []interface{}
	for rows.Next() {
		row, err := r.Read(rows.Scan)
		if err != nil {
			return nil, err
		}

		data = append(data, row)
	}

	return data, nil
}

// SetLocalData inserts the provided data using the sql staements associated with
// method param provided.
func (d *DB) SetLocalData(method utils.Method, params ...interface{}) error {
	stmt, ok := reqToStmt[method]
	if !ok {
		return fmt.Errorf("missing query for method %q", method)
	}

	_, err := d.db.ExecContext(d.ctx, d.formatPreparedStmt(stmt), params...)
	if err != nil {
		err = fmt.Errorf("inserting data for method %q failed: %v", method, err)
		return err
	}

	return nil
}

// CleanUpLocalData removes any dirty writes that may have been written on a certain
// last synced block.
func (d *DB) CleanUpLocalData(lastSyncedBlock uint64) {
	for _, stmt := range cleanUpStmt {
		// if an error in one query occurs, do no stop.
		_, err := d.db.ExecContext(d.ctx, d.formatPreparedStmt(stmt), lastSyncedBlock)
		if err != nil {
			log.Errorf("query %q failed: %v", stmt, err)
		}
	}
}

// formatPreparedStmt sets the required blind placedholder in the prepared
// statement. By default they are set to postgres standard placeholder.
func (d *DB) formatPreparedStmt(stmt string) string {
	if d.driver != utils.PostgresDriverName {
		placeholderRegex := utils.SupportedDbDrivers[utils.PostgresDriverName]
		replacementStr, _ := utils.SupportedDbDrivers[d.driver]
		stmt = regexp.MustCompile(placeholderRegex).ReplaceAllString(stmt, replacementStr)
	}
	return stmt
}
