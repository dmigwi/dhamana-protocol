package storage

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/bitcomplete/sqltestutil"
	"github.com/dmigwi/dhamana-protocol/client/servertypes"
	"github.com/dmigwi/dhamana-protocol/client/utils"
	"github.com/ethereum/go-ethereum/common"
	"gotest.tools/assert/cmp"
)

var (
	db       *DB
	ctx      context.Context
	cancelFn context.CancelFunc

	pgContainer *sqltestutil.PostgresContainer
)

// TestMain sets up the test
func TestMain(m *testing.M) {
	ctx, cancelFn = context.WithCancel(context.Background())

	var err error
	var pgContainer *sqltestutil.PostgresContainer

	defer func() {
		if err != nil {
			cancelFn()
			log.Error(err)
			os.Exit(1)
		}

		// clean up the db after tests are complete
		pgContainer.Shutdown(ctx)
	}()

	// use a mocked postgres db to run tests.
	pgContainer, err = sqltestutil.StartPostgresContainer(ctx, "12")
	if err != nil {
		err = fmt.Errorf("initializing the postgres container failed: %v", err)
		return
	}

	db, err = NewDB(ctx, pgContainer.ConnectionString())
	if err != nil {
		return
	}

	// Insert the initial records for tests.
	err = insertTestData()
	if err != nil {
		err = fmt.Errorf("inserting initial data failed: %v", err)
		return
	}

	m.Run()
}

// insertTestData inserts sample data into the tables.
func insertTestData() error {
	tableBondStmt := "INSERT INTO table_bond(" +
		"bond_address,issuer_address,holder_address," +
		"created_at_block,principal,coupon_rate,coupon_date,maturity_date," +
		"currency,intro_msg,last_status,last_synced_block)" +
		"VALUES (?,?,?,?,?,?,?,?,?,?,?,?)"

	tableBondData := [][]interface{}{
		{
			"0xc02aaa39b223fe8d0a0e5c4f27ead9083c756dha", // bond_address
			"0xf977814e90da44bfa03b6295a0616a897441aadd", // issuer_address
			"0x47ac0fb4f2d84898e4d9e7b4dab3c24507a6dhod", // holder_address
			71,                              // created_at_block
			14000,                           // principal
			7,                               // coupon_rate
			2,                               // coupon_date
			"2024-08-02 00:00:00.501361+03", // maturity_date
			0,                               // currency
			"This is an encrypted message",  // intro_msg
			3,                               // last_status
			89,                              // last_synced_block
		},
		{ // Data when a bond is created.
			"0xc61b9bb3a7a0767e3179713f3a5c7a9aedce1dhb", // bond_address
			"0x2b6ed29a95753c3ad948348e3e7b1a251080fadd",
			"", 72, 0, 0, 0, "", 0, "", 0, 72,
		},
		{ // Data when holder is selected and bond terms updated.
			"0xc61b9bb3a7a0767e3179713f3a5c7a9aedce1dhc", // bond_address
			"0x2b6ed29a95753c3ad948348e3e7b1a251080fadd",
			"0x47ac0fb4f2d84898e4d9e7b4dab3c24507a6dhod",
			72, 167000, 7, 2, "2024-08-02 00:00:00.501361+03", 1,
			"", 1, 75,
		},
		{ // Data when intro_msg is updated by the bond issuer.
			"0xc61b9bb3a7a0767e3179713f3a5c7a9aedce1dhd", // bond_address
			"0x2b6ed29a95753c3ad948348e3e7b1a251080fadd",
			"0x47ac0fb4f2d84898e4d9e7b4dab3c24507a6dhod",
			72, 167000, 7, 2, "2024-08-02 00:00:00.501361+03",
			1, "This is an encrypted message", 3, 80,
		},
	}

	tableStatusStmt := "INSERT INTO table_status(" +
		"sender,bond_address,bond_status,last_synced_block" +
		") VALUES (?,?,?,?,?)"

	tableStatusData := [][]interface{}{
		{ // Data when on setting holder address during status HolderUpdate. Update made by the bond Issuer.
			"0xf977814e90da44bfa03b6295a0616a897441aadd", // sender
			"0xc02aaa39b223fe8d0a0e5c4f27ead9083c756dha", // bond_address
			1, 76,
		},
		{ // Data when the bond moved to status TermsAgreement. Update made by the bond holder.
			"0x47ac0fb4f2d84898e4d9e7b4dab3c24507a6dhod",
			"0xc02aaa39b223fe8d0a0e5c4f27ead9083c756dha",
			2, 80,
		},
		{ // Data when the bond moved to status bondInDiputed. Update made by the bond Issuer.
			"0xf977814e90da44bfa03b6295a0616a897441aadd",
			"0xc02aaa39b223fe8d0a0e5c4f27ead9083c756dha",
			3, 85,
		},
	}

	tableStatusSignedStmt := "INSERT INTO table_status_signed(" +
		"sender,bond_address,bond_status,last_synced_block" +
		") VALUES (?,?,?,?)"

	tableStatusSignedData := [][]interface{}{
		{ // Data when the bond Holder signed status bondInDiputed. Sent by the bond Holder.
			"0x47ac0fb4f2d84898e4d9e7b4dab3c24507a6dhod", // sender
			"0xc02aaa39b223fe8d0a0e5c4f27ead9083c756dha", // bond_address
			3, 89,
		},
	}

	tableChatStmt := "INSERT INTO table_chat(" +
		"sender,bond_address,chat_msg,last_synced_block" +
		") VALUES (?,?,?,?)"

	tableChatData := [][]interface{}{
		{ // Data when a potential bond Holder expressed interest. Sent by the bond Holder.
			"0x47ac0fb4f2d84898e4d9e7b4dab3c24507a6dhod", // sender
			"0x47ac0fb4f2d84898e4d9e7b4dab3c24507a6dhod", // bond_address
			"xxxxxxx-encrypted", 74,
		},
		{ // Data when the bond Holder accepted the Issuer bond terms. Sent by the bond Holder.
			"0x47ac0fb4f2d84898e4d9e7b4dab3c24507a6dhod",
			"0xc02aaa39b223fe8d0a0e5c4f27ead9083c756dha",
			"xxxxxxx-encrypted", 76,
		},
		{ // Data when the bond Issuer explain why they moved the bond to status bondInDiputed.
			"0xf977814e90da44bfa03b6295a0616a897441aadd",
			"0xc02aaa39b223fe8d0a0e5c4f27ead9083c756dha",
			3, 90,
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
	limit := 5
	offset := 0

	expectedData := []servertypes.BondResp{
		{
			BondAddress: common.HexToAddress("0xc61b9bb3a7a0767e3179713f3a5c7a9aedce1dhd"), // bond_address
			Issuer:      common.HexToAddress("0x2b6ed29a95753c3ad948348e3e7b1a251080fadd"), // issuer_address
			CreatedTime: time.Time{},                                                       // Auto generated by postgres
			CouponRate:  7,
			Currency:    1,
			LastStatus:  3,
		},
		{
			BondAddress: common.HexToAddress("0xc02aaa39b223fe8d0a0e5c4f27ead9083c756dha"),
			Issuer:      common.HexToAddress("0xf977814e90da44bfa03b6295a0616a897441aadd"),
			CreatedTime: time.Time{}, // Auto generated by postgres
			CouponRate:  7,
			Currency:    0,
			LastStatus:  3,
		},
		{
			BondAddress: common.HexToAddress("0xc61b9bb3a7a0767e3179713f3a5c7a9aedce1dhc"),
			Issuer:      common.HexToAddress("0x2b6ed29a95753c3ad948348e3e7b1a251080fadd"),
			CreatedTime: time.Time{}, // Auto generated by postgres
			CouponRate:  7,
			Currency:    1,
			LastStatus:  1,
		},
		{
			BondAddress: common.HexToAddress("0xc61b9bb3a7a0767e3179713f3a5c7a9aedce1dhb"),
			Issuer:      common.HexToAddress(""),
			CreatedTime: time.Time{}, // Auto generated by postgres
			CouponRate:  0,
			Currency:    0,
			LastStatus:  0,
		},
	}

	// Compare two structs.
	compare := func(returned, expected interface{}) {
		if !cmp.Equal(returned, expected)().Success() {
			r, _ := json.Marshal(returned)
			e, _ := json.Marshal(expected)
			t.Fatalf("expected returned data %s to match %s but it didn't", string(r), string(e))
		}
	}

	t.Run("Test GetBonds results", func(t *testing.T) {
		data, err := db.QueryLocalData(utils.GetBonds, new(servertypes.BondResp), sender, limit, offset)
		if err != nil {
			t.Fatalf("expected no error but found: %v", err)
		}

		for i, res := range data {
			// Compares the two structs.
			compare(res, expectedData[i])
		}
	})

	maturityDate, _ := time.Parse(utils.PgDateFormat, "2024-08-02 00:00:00.501361+03")
	dataExp := servertypes.BondByAddressResp{
		BondResp: &servertypes.BondResp{
			BondAddress: common.HexToAddress("0xc02aaa39b223fe8d0a0e5c4f27ead9083c756dha"),
			Issuer:      common.HexToAddress("0xf977814e90da44bfa03b6295a0616a897441aadd"),
			CreatedTime: time.Time{}, // Auto generated by postgres
			CouponRate:  7,
			Currency:    0,
			LastStatus:  3,
		},
		Holder:          common.HexToAddress("0x47ac0fb4f2d84898e4d9e7b4dab3c24507a6dhod"),
		CreatedAtBlock:  71,
		Principal:       14000,
		CouponDate:      2,
		MaturityDate:    maturityDate,
		IntroMessage:    "This is an encrypted message",
		LastUpdate:      time.Time{}, // Auto generated by postgres
		LastSyncedBlock: 89,
	}

	t.Run("Test GetBondByAddress result", func(t *testing.T) {
		data, err := db.QueryLocalData(utils.GetBonds, new(servertypes.BondResp), sender, limit, offset)
		if err != nil {
			t.Fatalf("expected no error but found: %v", err)
		}

		compare(data, dataExp)
	})
}
