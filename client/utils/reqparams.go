// Copyright (c) 2023 Migwi Ndung'u
// See LICENSE for details.

package utils

type (
	// ParamType defines supported request parameters.
	ParamType string

	// MethodType defines type in relation to how the method is implemented
	MethodType int

	// Method defines the specific method names implemented.
	Method string
)

const (
	// Uint8Type defines unsigned integer parameter value type of uint8.
	Uint8Type ParamType = "uint8"

	// Uint16Type defines unsigned integer parameter value type of uint16.
	Uint16Type ParamType = "uint16"

	// Uint32Type defines unsigned integer parameter value type of uint32.
	Uint32Type ParamType = "uint32"

	// AddressType defines the address type as supportted in ethereum types.
	AddressType ParamType = "address"

	// StringType defines a string value type.
	StringType ParamType = "string"

	// UnsupportedType defines all other types not classified as int, float or string
	UnsupportedType ParamType = "unsupported"

	LocalType     MethodType = iota // Locally implemented
	ContractType                    // Implemented by the contracts
	ServerKeyType                   // Method for route /serverpubkey
	UnknownType                     // method not supported

	// --- Server methods supported ---

	// contract type methods - Sent via the server

	CreateBond       Method = "createBond"
	AddMessage       Method = "addMessage"
	SignBondStatus   Method = "signBondStatus"
	UpdateBodyInfo   Method = "updateBodyInfo"
	UpdateBondHolder Method = "updateBondHolder"
	UpdateBondStatus Method = "updateBondStatus"

	// server key type method - Sent via the server

	GetServerPubKey Method = "getServerPubKey"

	// Local type methods - Sent via the server

	GetBonds         Method = "getBonds"
	GetBondByAddress Method = "getBondByAddress"

	// Local Utils Methods. Results not sent via the server

	GetLastSyncedBlock Method = "getLastSyncedBlock"

	UpdateBondBodyTerms  Method = "updateBondBodyTerms"
	UpdateBondMotivation Method = "updateBondMotivation"
	UpdateHolder         Method = "updateHolder"
	UpdateLastStatus     Method = "updateLastStatus"
	InsertNewBondCreated Method = "insertNewBondCreated"
	InsertNewChatMessage Method = "insertNewchatMsg"
	InsertStatusChange   Method = "insertStatusChange"
	InsertStatusSigned   Method = "insertStatusSigned"
)

var (
	// contractMethods is a mapping of the supported contract methods with their respective
	// parameter types and count expected. Parameter types are placed at the
	// position they are expected.
	contractMethods = map[Method][]ParamType{
		// createBond creates a new bond instance owned by the method sender.
		// No user parameters are expected.
		CreateBond: {},

		// addMessage is used to update the bond details and also send bond chats.
		// Parameters Required: bondAddress address, tag uint8, message string
		// bondAddress => Defines the address of the bond in question.
		// tag => Defines the type of message being sent
		//			tag 0: => General bond chat message.
		// 			tag 1: => Bond Intro by the bond issuer.
		// 			tag 2: => Bond Security by the issuer.
		// 			tag 3: => Bond Appendix by the issuer.
		// message => Defines the actual message being sent. Should be limited
		// to 1000 characters before encryption.
		AddMessage: {AddressType, Uint8Type, StringType},

		// signBondStatus is used to show the sender has approved changes to the
		// bond as they are in the current bond status. i.e. To approve the
		// TermsAgreement stage, the sender signs that status. To resolve
		// dispute the sender signs the BondInDispute status for the specific bond.
		// Parameters Required: bondAddress address
		// bondAddress => Defines the address of the bond in question.
		SignBondStatus: {AddressType},

		// updateBodyInfo is used to update the body fields.
		// Parameter Required: bondAddress address, principal uint32,
		// 		couponRate uint8, couponDate uint32, maturityDate uint32, currency uint8
		// bondAddress => Defines the address of the bond in question.
		// principal => Defines the asking amount in the currency type supported.
		// couponRate => Defines the percentage of the interest payable.
		// couponDate => Defines interval of time when the interest payment is due.
		// 			CouponDate: 0 => Hourly
		// 			CouponDate: 1 => Daily
		// 			CouponDate: 2 => Weekly
		// 			CouponDate: 3 => Every-Forty-Night
		// 			CouponDate: 4 => Monthly
		// 			CouponDate: 5 => Quarterly
		// 			CouponDate: 6 => Yearly
		// 			CouponDate: 7 => Bi-Anually
		// 			CouponDate: 8 => Every-3-Years
		// 			CouponDate: 9 => Every-4-Years
		// 			CouponDate: 10 => Every-5-Years ..e.t.c.
		// maturityDate => Defines time in seconds when the issuer should finalise paying it.
		// Currency => Defines the currency type supported.
		// 			Currency: 0  => represents the fiat type.
		// 			Currency: 1 => represents Bitcoin.
		// 			Currency: 2 => represents Ethereum.
		// 			Currency: 3 => represents Ethereum Classic.
		// 			Currency: 4 => represents Ripple
		// 			Currency: 5 => represents Tether Coin.
		// 			Currency: 6 => represents Decred Coin.
		UpdateBodyInfo: {AddressType, Uint32Type, Uint8Type, Uint32Type, Uint32Type, Uint8Type},

		// updateBondHolder is used by the issuer during the HolderSelection stage to set
		// a potential bond holder.
		// Parameter Required: bondAddress string, holderAddress string
		// bondAddress => Defines the address of the bond in question.
		// holderAddress => Defines the address of potential holder choosen.
		UpdateBondHolder: {AddressType, AddressType},

		// updateBondStatus is used to move the bond along the supported bond status
		// stages.
		// Parameter Required: bondAddress string, status uint8
		// bondAddress => Defines the address of the bond in question.
		// status => Defines the bond stages in its lifecycle
		// 	 		status: 0 => represents Negotiating
		// 	 		status: 1 => represents HolderSelection
		// 	 		status: 2 => represents BondInDispute
		// 	 		status: 3 => represents TermsAgreement
		// 	 		status: 4 => represents ContractSigned
		// 	 		status: 5 => represents BondReselling
		// 	 		status: 6 => represents BondFinalised
		UpdateBondStatus: {AddressType, Uint8Type},
	}

	// localMethods is a mapping of the supported locally implemented methods with
	// their respective parameter types and count expected. Parameter types are
	// placed at the position they are expected to be.
	localMethods = map[Method][]ParamType{
		// getBondByAddress returns a bond at any status if the request sender is
		// also the bond issuer otherwise only returns bond with status
		// Negotiating.
		// Parameter Required: bondAddress string
		// bondAddress => Defines the address of the bond in question.
		GetBondByAddress: {AddressType},

		// getBonds returns all the bonds with status Negotiating or owned by
		// the sender if their current status status is past Negotiating stage.
		// No parameter Required
		GetBonds: {},
	}

	// serverKeyMethod defines the method used to query the server keys
	serverKeyMethod = map[Method][]ParamType{
		// getServerPubKey is used to query the session's server public key.
		// Parameter Required: clientPubkey string
		// The client provides its public key and in return the server sends
		// back its public key. Using diffie-hellman, a sharedkey developed
		// is used to communicate securely between the client and the server.
		GetServerPubKey: {StringType},
	}
)

// GetMethodParams returns the parameters of the method provided if supported.
func GetMethodParams(method Method) (implementation MethodType, param []ParamType) {
	// contract implemented methods
	if data, ok := contractMethods[method]; ok {
		return ContractType, data
	}

	// Local methods
	if data, ok := localMethods[method]; ok {
		return LocalType, data
	}

	// Server keys method
	if data, ok := serverKeyMethod[method]; ok {
		return ServerKeyType, data
	}

	return UnknownType, nil
}
