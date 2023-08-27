// Copyright (c) 2023 Migwi Ndung'u
// See LICENSE for details.

package storage

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/dmigwi/dhamana-protocol/client/utils"
	_ "github.com/lib/pq" // postgres
)

const (
	// semVersion holds the current semantic version requires for the tables.
	// The semantic version stored in the tables_version must match this value
	// otherwise the system will not be able initiate the db instance until the
	// user manually handles the data migration or creates a new
	// database to use with the new tables.
	semVersion = "v0.0.1"

	// createVersionTable enables version tables preventing tables with
	// incompatible schemas from being used.
	createVersionTable = "CREATE TABLE IF NOT EXISTS tables_version (" +
		"id SERIAL PRIMARY KEY," +
		"sem_version VARCHAR(10) UNIQUE," +
		"tables_created_on TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP)"

	// createTableBond is a prepared statement creating a table identified with the
	// name table_bond if it doesn't exists.
	createTableBond = "CREATE TABLE IF NOT EXISTS table_bond (" +
		"id SERIAL PRIMARY KEY," +
		"bond_address VARCHAR(42) UNIQUE NOT NULL," +
		"issuer_address VARCHAR(42) NOT NULL," +
		"holder_address VARCHAR(42)," +
		"created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP," +
		"created_at_block INTEGER NOT NULL," +
		"principal INTEGER," +
		"coupon_rate SMALLINT CHECK (coupon_rate BETWEEN 0 AND 100)," +
		"coupon_date SMALLINT CHECK (coupon_date BETWEEN 0 AND 50)," +
		"maturity_date TIMESTAMPTZ," +
		"currency SMALLINT CHECK (currency BETWEEN 0 AND 50)," +
		"intro_msg TEXT," +
		"last_status SMALLINT CHECK (last_status BETWEEN 0 AND 10)," +
		"last_update TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP," +
		"last_synced_block INTEGER NOT NULL)"

	// createTableBondStatus is a prepared statement creating a table identified
	// with the name table_status if it doesn't exists.
	createTableBondStatus = "CREATE TABLE IF NOT EXISTS table_status (" +
		"id SERIAL PRIMARY KEY," +
		"sender VARCHAR(42) NOT NULL," +
		"bond_address VARCHAR(42) NOT NULL," +
		"bond_status SMALLINT NOT NULL CHECK(bond_status BETWEEN 0 AND 10)," +
		"added_on TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP," +
		"last_synced_block INTEGER NOT NULL)"

	// createTableBondStatusSigned is a prepared statement creating a table
	// identified with the name table_status if it doesn't exists.
	createTableBondStatusSigned = "CREATE TABLE IF NOT EXISTS table_status_signed (" +
		"id SERIAL PRIMARY KEY," +
		"sender VARCHAR(42) NOT NULL," +
		"bond_address VARCHAR(42) NOT NULL," +
		"bond_status SMALLINT NOT NULL CHECK(bond_status BETWEEN 0 AND 10)," +
		"signed_on TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP," +
		"last_synced_block INTEGER NOT NULL)"

	// createChatTable is a prepared statement creating a table identified with
	// the name table_chat if it doesn't exists.
	createChatTable = "CREATE TABLE IF NOT EXISTS table_chat (" +
		"id SERIAL PRIMARY KEY," +
		"sender VARCHAR(42) NOT NULL," +
		"bond_address VARCHAR(42) NOT NULL," +
		"chat_msg TEXT NOT NULL," +
		"created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP," +
		"last_synced_block INTEGER NOT NULL)"

	// fetchBonds is a prepared statement that fetches all the bonds that owned
	// by the bond party with the address or they are still in the negotiation stage.
	fetchBonds = "SELECT bond_address,issuer_address,created_at,coupon_rate,currency,last_status " +
		"FROM table_bond WHERE issuer_address = $1 OR last_status = 0 " +
		"OR holder_address = $2 ORDER BY last_update DESC LIMIT $3 OFFSET $4"

	// fetchBondByAddress is a prepared statement that returns a bond identified by
	// the provided address if the sender is a party to the bond or the bond
	// is still in the negotiation stage.
	fetchBondByAddress = "SELECT bond_address,issuer_address,holder_address," +
		"created_at,created_at_block,principal,coupon_rate,coupon_date," +
		"maturity_date,currency,intro_msg,last_status,last_update,last_synced_block " +
		"FROM table_bond WHERE bond_address = $1 AND " +
		"(last_status = 0 OR issuer_address = $2 OR holder_address = $3)"

	// fetchChats is a prepared statement that fetches the conversation within
	// the bond identified by the provided address if the sender is a bond party
	// or its still in the negotiation stage.
	fetchChats = "SELECT c.sender, c.bond_address, c.chat_msg, c.created_at, " +
		"c.last_synced_block FROM table_chat as c LEFT JOIN table_bond as b " +
		"ON c.bond_address = b.bond_address WHERE b.bond_address = $1 AND " +
		"(b.issuer_address = $2 OR b.last_status = 0 OR b.holder_address = $3) " +
		"ORDER BY c.created_at DESC LIMIT $4 OFFSET $5"

	// fetchTableVersion fetches the last set tables version.
	fetchTableVersion = "SELECT sem_version,tables_created_on " +
		"FROM tables_version ORDER BY id DESC LIMIT 1"

	// fetchLastSyncBlock returns the last block to be synced on the table_bond.
	fetchLastSyncBlock = "SELECT last_synced_block FROM table_bond ORDER BY" +
		" last_synced_block DESC LIMIT 1"

	// setBondBodyTerms updates the table_bond with data from the BondBodyTerms event.
	setBondBodyTerms = "UPDATE table_bond SET principal = $1, coupon_rate = $2, " +
		"coupon_date = $3, maturity_date = $4, currency = $5, last_update = $6, " +
		"last_synced_block = $7 WHERE bond_address = $8"

	// setBondMotivation update the table_bond with data from BondMotivation event.
	setBondMotivation = "UPDATE table_bond SET intro_msg = $1, last_update = $2, " +
		"last_synced_block = $3 WHERE bond_address = $4"

	// setHolder updates table_bond with data from HolderUpdate event.
	setHolder = "UPDATE table_bond SET holder_address = $1, last_update = $2," +
		"last_synced_block = $3 WHERE bond_address = $4"

	// setLastStatus updates table_bond with data from StatusChange event.
	setLastStatus = "UPDATE table_bond SET last_status = $1, last_update = $2, " +
		"last_synced_block = $3 WHERE bond_address = $4"

	// addNewBondCreated inserts into table_bond new data from event NewBondCreated.
	addNewBondCreated = "INSERT INTO table_bond (bond_address, issuer_address, " +
		"created_at_block, last_synced_block) VALUES ($1, $2, $3, $4)"

	// ddNewChatMessage inserts into table_chat new data from event NewChatMessage.
	addNewChatMessage = "INSERT INTO table_chat (sender, bond_address, " +
		"chat_msg, last_synced_block) VALUES ($1, $2, $3, $4)"

	// addStatusChange inserts into table_status new data from event StatusChange.
	addStatusChange = "INSERT INTO table_status (sender, bond_address, " +
		"bond_status, last_synced_block) VALUES ($1, $2, $3, $4)"

	// addStatusSigned inserts into table_status_signed new data from event StatusSigned.
	addStatusSigned = "INSERT INTO table_status_signed (sender, bond_address, " +
		"bond_status, last_synced_block) VALUES ($1, $2, $3, $4)"

	// addTablesVersion inserts into tables_version the latest supported tables version.
	addTablesVersion = "INSERT INTO tables_version (sem_version) VALUES ($1)"

	dropTableBondRecords         = "DELETE FROM table_bond WHERE last_synced_block = $1"
	dropTableStatusRecords       = "DELETE FROM table_status WHERE last_synced_block = $1"
	dropTableStatusSignedRecords = "DELETE FROM table_status_signed WHERE last_synced_block = $1"
	dropTableChatRecords         = "DELETE FROM table_chat WHERE last_synced_block = $1"
)

// tablesToSQLStmt is an array of sql statements used to create the missing tables
// if they don't exist.
var tablesSQLStmt = []string{
	createVersionTable,
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
	db  *sql.DB
	ctx context.Context
}

// Reader defines the method that reads the row fields into the require data interface.
// To read data, pass pointers to the expect field the parameter function.
type Reader interface {
	Read(fn func(fields ...any) error) (interface{}, error)
}

// ConnectionString returns on the connection string format supported by postgres.
// https://www.postgresql.org/docs/current/libpq-connect.html#LIBPQ-CONNSTRING
func ConnectionString(port uint16, host, user, password, dbname string) string {
	return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)
}

// NewDB returns an opened db instance whose connection has been tested with
// ping request. It generates the required tables if they don't exist.
func NewDB(ctx context.Context, connInfo string) (*DB, error) {
	db, err := sql.Open("postgres", connInfo)
	if err != nil {
		log.Errorf("unable to open to postgres db: err %v", err)
		return nil, err
	}

	if err = db.PingContext(ctx); err != nil {
		log.Errorf("connection to postgres db failed: err %v", err)
		return nil, err
	}

	log.Info("Confirming that all the database tables exists")
	for i, stmt := range tablesSQLStmt {
		if _, err := db.ExecContext(ctx, stmt); err != nil {
			log.Errorf("creating a table index (%d) failed Error: %v", i, err)
			return nil, err
		}
	}

	dbInstance := &DB{
		db:  db,
		ctx: ctx,
	}

	// -- Confirm the semantic version match the required on --

	var tableversion string
	var createdDate time.Time

	err = db.QueryRowContext(ctx, fetchTableVersion).Scan(&tableversion, &createdDate)
	if err != nil && err != sql.ErrNoRows {
		log.Errorf("unable to fetch tables versions: %v", err)
		return nil, err
	}

	switch tableversion {
	case "":
		log.Infof("Versioning the newly created tables with version=%s", semVersion)

		if _, err = db.ExecContext(ctx, addTablesVersion, semVersion); err != nil {
			log.Errorf("unable version the newly created tables : %v", err)
			return nil, err
		}

	case semVersion:
		// The correct tables version was found.
		log.Infof("Confirmed all the %d versioned tables exists", len(tablesSQLStmt))

	default:
		// versions mismatch found. Exit till the issue is resolved.
		err = fmt.Errorf("expected the tables version %s but found version %s created on = %v",
			semVersion, tableversion, createdDate)
		return nil, err
	}

	return dbInstance, nil
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
	case utils.GetBondByAddress:
		params = append(params, []interface{}{sender, sender}...)

	case utils.GetBonds, utils.GetChats:
		// The last two param values are always {limit, offset}.
		// Append sender params before those two params.
		n := len(params)
		var data []interface{}
		if len(params[:n-2]) > 0 {
			data = append(data, params[:n-2]...)
		}
		data = append(data, []interface{}{sender, sender}...)
		data = append(data, params[n-2:]...) // {limit, offset}.

		params = data
	}

	rows, err := d.db.QueryContext(d.ctx, stmt, params...)
	if err != nil {
		return nil, fmt.Errorf("fetching query for method %q failed: %v", method, err)
	}

	defer rows.Close()

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

	if _, err := d.db.ExecContext(d.ctx, stmt, params...); err != nil {
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
		if _, err := d.db.ExecContext(d.ctx, stmt, lastSyncedBlock); err != nil {
			log.Errorf("query %q failed: %v", stmt, err)
		}
	}
}
