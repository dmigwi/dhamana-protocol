package storage

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/bitcomplete/sqltestutil"
	"github.com/btcsuite/btclog"
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

type testLogger struct {
	btclog.Logger
}

func (tl *testLogger) Debugf(format string, params ...interface{}) {
	fmt.Printf("DEBUG: %v %v \n", format, params)
}

func (tl *testLogger) Infof(format string, params ...interface{}) {
	fmt.Printf("INFO: %v %v \n", format, params)
}

func (tl *testLogger) Warnf(format string, params ...interface{}) {
	fmt.Printf("WARN: %v %v \n", format, params)
}

func (tl *testLogger) Errorf(format string, params ...interface{}) {
	fmt.Printf("ERROR: %v %v \n", format, params)
}

func (tl *testLogger) Info(v ...interface{}) {
	fmt.Printf("INFO: %v \n", v)
}

func (tl *testLogger) Warn(v ...interface{}) {
	fmt.Printf("WARN: %v \n", v)
}

func (tl *testLogger) Error(v ...interface{}) {
	fmt.Printf("ERROR: %v \n", v)
}

func (tl *testLogger) Debug(v ...interface{}) {
	fmt.Printf("DEBUG: %v \n", v)
}

// TestMain sets up the test
func TestMain(m *testing.M) {
	ctx, cancelFn = context.WithCancel(context.Background())

	// Assign a test log instance.
	log = new(testLogger)

	var pgContainer *sqltestutil.PostgresContainer

	var err error
	processError := func() {
		if err != nil {
			fmt.Println("Cleaning up the tables: ", err)

			// clean up the db after tests are complete
			pgContainer.Shutdown(ctx)

			cancelFn()
			os.Exit(1)
		}
	}

	// use a mocked postgres db to run tests.
	pgContainer, err = sqltestutil.StartPostgresContainer(ctx, "12")
	processError()

	db, err = NewDB(ctx, pgContainer.ConnectionString()+"?sslmode=disable")
	processError()

	// Insert the initial records for tests.
	err = insertTestData()
	processError()

	m.Run()
}

// insertTestData inserts sample data into the tables.
func insertTestData() error {
	tableBondStmt := "INSERT INTO table_bond (" +
		"bond_address, issuer_address, holder_address, created_at_block, " +
		"principal, coupon_rate, coupon_date, maturity_date, currency, " +
		"intro_msg,last_status, last_synced_block) " +
		"VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)"

	tableBondData := [][]interface{}{
		{
			"0xc02aaa39b223fe8d0a0e5c4f27ead9083c756dba", // bond_address
			"0xf977814e90da44bfa03b6295a0616a897441aadd", // issuer_address
			"0x47ac0fb4f2d84898e4d9e7b4dab3c24507a6dbod", // holder_address
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
			"0xc61b9bb3a7a0767e3179713f3a5c7a9aedce1dbb", // bond_address
			"0x47ac0fb4f2d84898e4d9e7b4dab3c24507a6dbod",
			"", 72, 0, 0, 0, "2024-08-02 00:00:00.501361+03", 0, "", 0, 72,
		},
		{ // Data when holder is selected and bond terms updated.
			"0xc61b9bb3a7a0767e3179713f3a5c7a9aedce1dbc", // bond_address
			"0x2b6ed29a95753c3ad948348e3e7b1a251080fadd",
			"0x47ac0fb4f2d84898e4d9e7b4dab3c24507a6dbod",
			72, 167000, 7, 2, "2024-08-02 00:00:00.501361+03", 1,
			"", 1, 75,
		},
		{ // Data when intro_msg is updated by the bond issuer.
			"0xc61b9bb3a7a0767e3179713f3a5c7a9aedce1dbd", // bond_address
			"0x2b6ed29a95753c3ad948348e3e7b1a251080fadd",
			"0x47ac0fb4f2d84898e4d9e7b4dab3c24507a6dbod",
			72, 167000, 7, 2, "2024-08-02 00:00:00.501361+03",
			1, "This is an encrypted message", 3, 80,
		},
	}

	tableStatusStmt := "INSERT INTO table_status(" +
		"sender, bond_address, bond_status, last_synced_block" +
		") VALUES ($1, $2, $3, $4)"

	tableStatusData := [][]interface{}{
		{ // Data when on setting holder address during status HolderUpdate. Update made by the bond Issuer.
			"0xf977814e90da44bfa03b6295a0616a897441aadd", // sender
			"0xc02aaa39b223fe8d0a0e5c4f27ead9083c756dba", // bond_address
			1, 76,
		},
		{ // Data when the bond moved to status TermsAgreement. Update made by the bond holder.
			"0x47ac0fb4f2d84898e4d9e7b4dab3c24507a6dbod",
			"0xc02aaa39b223fe8d0a0e5c4f27ead9083c756dba",
			2, 80,
		},
		{ // Data when the bond moved to status bondInDiputed. Update made by the bond Issuer.
			"0xf977814e90da44bfa03b6295a0616a897441aadd",
			"0xc02aaa39b223fe8d0a0e5c4f27ead9083c756dba",
			3, 85,
		},
	}

	tableStatusSignedStmt := "INSERT INTO table_status_signed(" +
		"sender, bond_address, bond_status, last_synced_block" +
		") VALUES ($1, $2, $3, $4)"

	tableStatusSignedData := [][]interface{}{
		{ // Data when the bond Holder signed status bondInDiputed. Sent by the bond Holder.
			"0x47ac0fb4f2d84898e4d9e7b4dab3c24507a6dbod", // sender
			"0xc02aaa39b223fe8d0a0e5c4f27ead9083c756dba", // bond_address
			3, 89,
		},
	}

	tableChatStmt := "INSERT INTO table_chat(" +
		"sender, bond_address, chat_msg, last_synced_block" +
		") VALUES ($1, $2, $3, $4)"

	tableChatData := [][]interface{}{
		{ // Data when a potential bond Holder expressed interest. Sent by the bond Holder.
			"0x47ac0fb4f2d84898e4d9e7b4dab3c24507a6dbod", // sender
			"0x47ac0fb4f2d84898e4d9e7b4dab3c24507567845", // bond_address
			"xxxxxxx-encrypted", 74,
		},
		{ // Data when the bond Holder accepted the Issuer bond terms. Sent by the bond Holder.
			"0x47ac0fb4f2d84898e4d9e7b4dab3c24507a6dbod",
			"0xc02aaa39b223fe8d0a0e5c4f27ead9083c756dba",
			"xxxxxxx-encrypted", 76,
		},
		{ // Data when the bond Issuer explain why they moved the bond to status bondInDiputed.
			"0xf977814e90da44bfa03b6295a0616a897441aadd",
			"0xc02aaa39b223fe8d0a0e5c4f27ead9083c756dba",
			"uywteuyrw ddhhdyugdhjdna", 90,
		},
	}

	tablesdata := map[string][][]interface{}{
		tableBondStmt:         tableBondData,
		tableStatusStmt:       tableStatusData,
		tableStatusSignedStmt: tableStatusSignedData,
		tableChatStmt:         tableChatData,
	}

	for query, data := range tablesdata {
		for i, v := range data {
			_, err := db.db.ExecContext(ctx, query, v...)
			if err != nil {
				err = fmt.Errorf("query at index %d method: %v failed with error: %v", i, query, err)
				return err
			}
		}
	}
	return nil
}

// TestQueryLocalData tests the functionality of QueryLocalData method.
func TestQueryLocalData(t *testing.T) {
	sender := "0x47ac0fb4f2d84898e4d9e7b4dab3c24507a6dbod"
	limit := 5
	offset := 0

	bondExp := []servertypes.BondResp{
		{
			BondAddress: common.HexToAddress("0xc61b9bb3a7a0767e3179713f3a5c7a9aedce1dbd"), // bond_address
			Issuer:      common.HexToAddress("0x2b6ed29a95753c3ad948348e3e7b1a251080fadd"), // issuer_address
			CreatedTime: time.Time{},                                                       // Auto generated by postgres
			CouponRate:  7,
			Currency:    1,
			LastStatus:  3,
		},
		{
			BondAddress: common.HexToAddress("0xc61b9bb3a7a0767e3179713f3a5c7a9aedce1dbc"),
			Issuer:      common.HexToAddress("0x2b6ed29a95753c3ad948348e3e7b1a251080fadd"),
			CreatedTime: time.Time{}, // Auto generated by postgres
			CouponRate:  7,
			Currency:    1,
			LastStatus:  1,
		},
		{
			BondAddress: common.HexToAddress("0xc61b9bb3a7a0767e3179713f3a5c7a9aedce1dbb"),
			Issuer:      common.HexToAddress("0x47ac0fb4f2d84898e4d9e7b4dab3c24507a6dbod"),
			CreatedTime: time.Time{}, // Auto generated by postgres
			CouponRate:  0,
			Currency:    0,
			LastStatus:  0,
		},
		{
			BondAddress: common.HexToAddress("0xc02aaa39b223fe8d0a0e5c4f27ead9083c756dba"),
			Issuer:      common.HexToAddress("0xf977814e90da44bfa03b6295a0616a897441aadd"),
			CreatedTime: time.Time{}, // Auto generated by postgres
			CouponRate:  7,
			Currency:    0,
			LastStatus:  3,
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

		if len(data) != len(bondExp) {
			t.Fatalf("expected %d  records but found %d records", len(bondExp), len(data))
		}

		for i, res := range data {
			// check returned data type.
			ex, ok := res.(*servertypes.BondResp)
			if !ok {
				t.Fatalf("expected the returned data to be of type *servertypes.BondResp but it wasn't")
			}

			if ex.CreatedTime.IsZero() {
				t.Fatalf("expected the db generated timestamp not to have a zero value")
			}

			// set to zero the db  auto-filled created_at field.
			ex.CreatedTime = time.Time{}

			// Compares the two structs.
			compare(*ex, bondExp[i])
		}
	})

	maturityDate, _ := time.Parse(utils.PgDateFormat, "2024-08-02 00:00:00.501361+03")
	dataExp := servertypes.BondByAddressResp{
		BondResp: &servertypes.BondResp{
			BondAddress: common.HexToAddress("0xc02aaa39b223fe8d0a0e5c4f27ead9083c756dba"),
			Issuer:      common.HexToAddress("0xf977814e90da44bfa03b6295a0616a897441aadd"),
			CreatedTime: time.Time{}, // Auto generated by postgres
			CouponRate:  7,
			Currency:    0,
			LastStatus:  3,
		},
		Holder:          common.HexToAddress("0x47ac0fb4f2d84898e4d9e7b4dab3c24507a6dbod"),
		CreatedAtBlock:  71,
		Principal:       14000,
		CouponDate:      2,
		MaturityDate:    maturityDate,
		IntroMessage:    "This is an encrypted message",
		LastUpdate:      time.Time{}, // Auto generated by postgres
		LastSyncedBlock: 89,
	}

	t.Run("Test GetBondByAddress result", func(t *testing.T) {
		data, err := db.QueryLocalData(utils.GetBondByAddress, new(servertypes.BondByAddressResp), sender,
			common.HexToAddress("0xc02aaa39b223fe8d0a0e5c4f27ead9083c756dba"))
		if err != nil {
			t.Fatalf("expected no error but found: %v", err)
		}

		if len(data) != 0 {
			t.Fatalf("expected one record returned but found %v records", len(data))
		}

		res, ok := data[0].(*servertypes.BondByAddressResp)
		if !ok {
			t.Fatalf("expected the returned data to be of type *servertypes.BondByAddressResp but it wasn't")
		}

		if res.CreatedTime.IsZero() || res.LastUpdate.IsZero() {
			t.Fatalf("expected the db generated timestamps not to have zero values")
		}

		// set to zero the db  auto-filled created_at field.
		res.CreatedTime = time.Time{}
		res.LastUpdate = time.Time{}

		compare(res, &dataExp)
	})

	chatsExp := []servertypes.ChatMsgsResp{
		{
			Sender:          common.HexToAddress("0xf977814e90da44bfa03b6295a0616a897441aadd"),
			BondAddress:     common.HexToAddress("0xc02aaa39b223fe8d0a0e5c4f27ead9083c756dba"),
			Message:         "uywteuyrw ddhhdyugdhjdna",
			CreatedTime:     time.Time{},
			LastSyncedBlock: 90,
		},
		{
			Sender:          common.HexToAddress("0x47ac0fb4f2d84898e4d9e7b4dab3c24507a6dbod"),
			BondAddress:     common.HexToAddress("0xc02aaa39b223fe8d0a0e5c4f27ead9083c756dba"),
			Message:         "xxxxxxx-encrypted",
			CreatedTime:     time.Time{},
			LastSyncedBlock: 76,
		},
	}

	t.Run("Test GetChats results", func(t *testing.T) {
		data, err := db.QueryLocalData(utils.GetChats, new(servertypes.ChatMsgsResp), sender,
			common.HexToAddress("0xc02aaa39b223fe8d0a0e5c4f27ead9083c756dba"), limit, offset)
		if err != nil {
			t.Fatalf("expected no error but found: %v", err)
		}

		if len(data) != len(chatsExp) {
			t.Fatalf("expected %d  records but found %d records", len(chatsExp), len(data))
		}

		for i, res := range data {
			// check returned data type.
			ex, ok := res.(*servertypes.ChatMsgsResp)
			if !ok {
				t.Fatalf("expected the returned data to be of type *servertypes.ChatMsgsResp but it wasn't")
			}

			if ex.CreatedTime.IsZero() {
				t.Fatalf("expected the db generated timestamp not to have a zero value")
			}

			// set to zero the db  auto-filled created_at field.
			ex.CreatedTime = time.Time{}

			// Compares the two structs.
			compare(ex, &chatsExp[i])
		}
	})

	blockExp := servertypes.LastSyncedBlockResp(90)

	t.Run("Test GetLastSyncedBlock result", func(t *testing.T) {
		data, err := db.QueryLocalData(utils.GetLastSyncedBlock, new(servertypes.LastSyncedBlockResp), "")
		if err != nil {
			t.Fatalf("expected no error but found: %v", err)
		}

		if len(data) != 0 {
			t.Fatalf("expected one record returned but found %v records", len(data))
		}

		compare(data, &blockExp)
	})
}

// TestSetLocalData tests if the inserts and update queries execute without
// returning an error.
func TestSetLocalData(t *testing.T) {
	testData := map[utils.Method][]interface{}{
		utils.InsertNewBondCreated: {
			"0xc61b9bb3a7a0767e317971000000000000001dbd", // bond_address
			"0x47ac0fb4f2d84898e4d9e7b4dab3c24507a6dbod", // issuer_address
			100, // last_synced_block
		},
		utils.InsertNewChatMessage: {
			"0x47ac0fb4f2d84898e4d9e7b4dab3c24507a6dbod", // sender
			"0xc61b9bb3a7a0767e317971000000000000001dbd", // bond_address
			"yuqteuqteuqeqe", // chat_msg
			120,              // last_synced_block
		},
		utils.InsertStatusChange: {
			"0x47ac0fb4f2d84898e4d9e7b4dab3c24507a6dbod", // sender
			"0xc61b9bb3a7a0767e317971000000000000001dbd", // bond_address
			1,   // bond_status
			120, // last_synced_block
		},
		utils.InsertStatusSigned: {
			"0x47ac0fb4f2d84898e4d9e7b4dab3c24507a6dbod", // sender
			"0xc61b9bb3a7a0767e317971000000000000001dbd", // bond_address
			3,   // bond_status
			120, // last_synced_block
		},
		utils.UpdateBondBodyTerms: {
			41564316,                        // principal
			8,                               // coupon_rate
			3,                               // coupon_date
			"2024-08-02 00:00:00.501361+03", // maturity_date
			2,                               // currency
			"2023-09-01 00:00:00.501361+03", // last_update
			120,                             // last_synced_block
			"0xc61b9bb3a7a0767e317971000000000000001dbd", // bond_address
		},
		utils.UpdateBondMotivation: {
			"xxxx",                          // intro_msg
			"2023-09-01 01:00:00.501361+03", // last_update
			120,                             // last_synced_block
			"0xc61b9bb3a7a0767e317971000000000000001dbd", // bond_address
		},
		utils.UpdateHolder: {
			"0xf97781467250000000000095a0616a8974422222", // holder
			"2023-09-01 02:00:00.501361+03",              // last_update
			120,                                          // last_synced_block
			"0xc61b9bb3a7a0767e317971000000000000001dbd", // bond_address
		},
		utils.UpdateLastStatus: {
			1,                               // last_status
			"2023-09-01 02:45:00.501361+03", // last_update
			120,                             // last_synced_block
			"0xc61b9bb3a7a0767e317971000000000000001dbd", // bond_address
		},
	}

	for mthd, td := range testData {
		t.Run(fmt.Sprintf("Test %v insert", mthd), func(t *testing.T) {
			err := db.SetLocalData(mthd, td...)
			if err != nil {
				t.Fatal(err)
			}
		})
	}
}

// TestCleanUpLocalData test if records on a certain synced block can be deleted
// on all tables.
func TestCleanUpLocalData(t *testing.T) {
	t.Run("Test CleanUpLocalData", func(t *testing.T) {
		var lastSyncedBlock uint64

		err := db.db.QueryRow(fetchLastSyncBlock).Scan(&lastSyncedBlock)
		if err != nil {
			t.Fatal(err)
		}

		db.CleanUpLocalData(lastSyncedBlock)

		var newLastSyncedBlock uint64

		err = db.db.QueryRow(fetchLastSyncBlock).Scan(&newLastSyncedBlock)
		if err != nil {
			t.Fatal(err)
		}

		if lastSyncedBlock > newLastSyncedBlock {
			t.Fatalf("after cleaning data on block %d last block synced shouldn't be %d",
				lastSyncedBlock, newLastSyncedBlock)
		}
	})
}
