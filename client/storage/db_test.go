package storage

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/dmigwi/dhamana-protocol/client/servertypes"
	"github.com/dmigwi/dhamana-protocol/client/utils"
	"golang.org/x/text/currency"
)

var (
	db       *DB
	ctx      context.Context
	cancelFn context.CancelFunc
)

// TestMain sets up the test
func TestMain(m *testing.M) {
	dir, err := os.MkdirTemp("", "test-")
	if err != nil {
		log.Errorf("unable to create a temporary dir: %v", err)
		os.Exit(1)
	}

	ctx, cancelFn = context.WithCancel(context.Background())

	db, err = NewDB(ctx, 0, utils.SqliteDriverName, "", "", "", "", filepath.Join(dir, "db"))
	if err != nil {
		cancelFn()
		log.Error(err)
		os.Exit(1)
	}

	// Insert the initial records for tests.
	err = insertTestData()
	if err != nil {
		cancelFn()
		log.Errorf("inserting initial data failed: %v", err)
		os.Exit(1)
	}

	m.Run()

	// clean up the db after tests are complete
	os.RemoveAll(dir)
}

// insertTestData inserts sample data into the tables.
func insertTestData() error {
	tableBondStmt := "INSERT INTO table_bond(" +
		"bond_address,issuer_address,holder_address,created_at," +
		"created_at_block,principal,coupon_rate,coupon_date,maturity_date," +
		"currency,intro_msg,last_status,last_update,last_synced_block)" +
		"VALUES (?,?,?,?,?,?,?,?,?,?,?,?,?,?)"

	tableBondData := [][]interface{}{
		{
			"0xc02aaa39b223fe8d0a0e5c4f27ead9083c756dha", // bond_address
			"0xf977814e90da44bfa03b6295a0616a897441aadd", // issuer_address
			"0x47ac0fb4f2d84898e4d9e7b4dab3c24507a6dhod", // holder_address
			"2023-08-02 19:10:25-00",                     // created_at
			71,                                           // created_at_block
			14000,                                        // principal
			7,                                            // coupon_rate
			1693048167,                                   // coupon_date
			1693149167,                                   // maturity_date
			0,                                            // currency
			"This is an encrypted message",               // intro_msg
			3,                                            // last_status
			"2023-08-04 08:53:21-00",                     // last_update
			89,                                           // last_synced_block
		},
		{ // Data when a bond is created.
			"0xc61b9bb3a7a0767e3179713f3a5c7a9aedce1dhb", // bond_address
			"0x2b6ed29a95753c3ad948348e3e7b1a251080fadd",
			"", "2023-08-02 19:10:25-07", 72, 0, 0,
			0, 0, 0, "", 0, "2023-08-02 19:10:25-00", 72,
		},
		{ // Data when holder is selected and bond terms updated.
			"0xc61b9bb3a7a0767e3179713f3a5c7a9aedce1dhc", // bond_address
			"0x2b6ed29a95753c3ad948348e3e7b1a251080fadd",
			"0x47ac0fb4f2d84898e4d9e7b4dab3c24507a6dhod",
			"2023-08-02 19:10:25-07", 72, 167000, 7,
			1693048167, 1693149167, 1, "", 1, "2023-08-03 10:10:31-00", 75,
		},
		{ // Data when intro_msg is updated by the bond issuer.
			"0xc61b9bb3a7a0767e3179713f3a5c7a9aedce1dhd", // bond_address
			"0x2b6ed29a95753c3ad948348e3e7b1a251080fadd",
			"0x47ac0fb4f2d84898e4d9e7b4dab3c24507a6dhod",
			"2023-08-02 19:10:25-07", 72, 167000, 7, 1693048167,
			1693149167, 1, "This is an encrypted message",
			3, "2023-08-04 17:54:25-00", 80,
		},
	}

	tableStatusStmt := "INSERT INTO table_status(" +
		"sender,bond_address,bond_status,added_on,last_synced_block" +
		") VALUES (?,?,?,?,?)"

	tableStatusData := [][]interface{}{
		{ // Data when on setting holder address during status HolderUpdate. Update made by the bond Issuer.
			"0xf977814e90da44bfa03b6295a0616a897441aadd", // sender
			"0xc02aaa39b223fe8d0a0e5c4f27ead9083c756dha", // bond_address
			1, "2023-08-03 09:10:05-00", 76,
		},
		{ // Data when the bond moved to status TermsAgreement. Update made by the bond holder.
			"0x47ac0fb4f2d84898e4d9e7b4dab3c24507a6dhod",
			"0xc02aaa39b223fe8d0a0e5c4f27ead9083c756dha",
			2, "2023-08-03 16:10:55-00", 80,
		},
		{ // Data when the bond moved to status bondInDiputed. Update made by the bond Issuer.
			"0xf977814e90da44bfa03b6295a0616a897441aadd",
			"0xc02aaa39b223fe8d0a0e5c4f27ead9083c756dha",
			3, "2023-08-04 02:33:21-00", 85,
		},
	}

	tableStatusSignedStmt := "INSERT INTO table_status_signed(" +
		"sender,bond_address,bond_status,signed_on,last_synced_block" +
		") VALUES (?,?,?,?,?)"

	tableStatusSignedData := [][]interface{}{
		{ // Data when the bond Holder signed status bondInDiputed. Sent by the bond Holder.
			"0x47ac0fb4f2d84898e4d9e7b4dab3c24507a6dhod", // sender
			"0xc02aaa39b223fe8d0a0e5c4f27ead9083c756dha", // bond_address
			3, "2023-08-04 08:53:21-00", 89,
		},
	}

	tableChatStmt := "INSERT INTO table_chat(" +
		"sender,bond_address,chat_msg,created_at,last_synced_block" +
		") VALUES (?,?,?,?,?)"

	tableChatData := [][]interface{}{
		{ // Data when a potential bond Holder expressed interest. Sent by the bond Holder.
			"0x47ac0fb4f2d84898e4d9e7b4dab3c24507a6dhod", // sender
			"0x47ac0fb4f2d84898e4d9e7b4dab3c24507a6dhod", // bond_address
			"xxxxxxx-encrypted", "2023-08-02 23:10:25-00", 74,
		},
		{ // Data when the bond Holder accepted the Issuer bond terms. Sent by the bond Holder.
			"0x47ac0fb4f2d84898e4d9e7b4dab3c24507a6dhod",
			"0xc02aaa39b223fe8d0a0e5c4f27ead9083c756dha",
			"xxxxxxx-encrypted", "2023-08-03 09:10:05-00", 76,
		},
		{ // Data when the bond Issuer explain why they moved the bond to status bondInDiputed.
			"0xf977814e90da44bfa03b6295a0616a897441aadd",
			"0xc02aaa39b223fe8d0a0e5c4f27ead9083c756dha",
			3, "2023-08-04 08:53:21-00", 90,
		},
	}

	tablesdata := map[string][][]interface{}{
		tableBondStmt:         tableBondData,
		tableStatusStmt:       tableStatusData,
		tableStatusSignedStmt: tableStatusSignedData,
		tableChatStmt:         tableChatData,
	}

	for query, data := range tablesdata {
		for _, v := range data {
			_, err := db.db.ExecContext(ctx, query, v...)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

// TestQueryLocalData tests the functionality of QueryLocalData method.
func TestQueryLocalData(t *testing.T) {
	sender := "0x47ac0fb4f2d84898e4d9e7b4dab3c24507a6dhod"

	expectedData := []servertypes.BondResp{
		{
			BondAddress:   "0xc02aaa39b223fe8d0a0e5c4f27ead9083c756dha", // bond_address
			Issuer: "0xf977814e90da44bfa03b6295a0616a897441aadd", // issuer_address
			CreatedTime: time.Time{},
			CouponRate: 0,
			Currency: 0,
			LastStatus: 0,
		}, {
			BondAddress:   "0xc61b9bb3a7a0767e3179713f3a5c7a9aedce1dhb",
			Issuer: "0x2b6ed29a95753c3ad948348e3e7b1a251080fadd",
			CreatedTime: time.Time{},
			CouponRate: 0,
			Currency: 0,
			LastStatus: 0,
		}, {
			BondAddress:   "0xc61b9bb3a7a0767e3179713f3a5c7a9aedce1dhc",
			Issuer: "0x2b6ed29a95753c3ad948348e3e7b1a251080fadd",
			CreatedTime: time.Time{},
			CreatedTime: time.Time{},
			CouponRate: 0,
			Currency: 0,
			LastStatus: 0,
		},
	}

	db.QueryLocalData(method utils.Method, r Reader, sender string, params ...interface{})
}
