// Copyright (c) 2023 Migwi Ndung'u
// See LICENSE for details.

package utils

// ValueType defines supported request parameters.
type ValueType string

// MethodType defines type in relation to how the method is implemented
type MethodType int

const (
	// IntType defines all integer value types.
	IntType ValueType = "int"

	// StringType defines all string value types.
	StringType ValueType = "string"

	// FloatType defines all float value types.
	FloatType ValueType = "float"

	// UnsupportedType defines all other types not classified as int, float or string
	UnsupportedType ValueType = "unsupported"

	LocalType     MethodType = iota // Locally implemented
	ContractType                    // Implemented by the contracts
	ServerKeyType                   // Method for route /serverpubkey
	UnknownType                     // method not supported
)

var (
	// ContractMethods is a mapping of the supported contract methods with their respective
	// parameter types and count expected. Parameter types are placed at the
	// position they are expected.
	ContractMethods = map[string][]ValueType{
		// createBond creates a new bond instance owned by the method sender.
		// No user parameters are expected.
		"createBond": {},

		// addMessage is used to update the bond details and also send bond chats.
		// Parameters Required: bondAddress string, tag int, message string
		// bondAddress => Defines the address of the bond in question.
		// tag => Defines the type of message being sent
		//			tag 0: => General bond chat message.
		// 			tag 1: => Bond Intro by the bond issuer.
		// 			tag 2: => Bond Security by the issuer.
		// 			tag 3: => Bond Appendix by the issuer.
		// message => Defines the actual message being sent. Should be limited
		// to 1000 characters before encryption.
		"addMessage": {StringType, IntType, StringType},

		// signBondStatus is used to show the sender has approved changes to the
		// bond as they are in the current bond status. i.e. To approve the
		// TermsAgreement stage, the sender signs that status. To resolve
		// dispute the sender signs the BondInDispute status for the specific bond.
		// Parameters Required: bondAddress string
		// bondAddress => Defines the address of the bond in question.
		"signBondStatus": {StringType},

		// updateBodyInfo is used to update the body fields.
		// Parameter Required: bondAddress string, principal int, couponRate int,
		//  	couponDate int, maturityDate int, currency int
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
		"updateBodyInfo": {StringType, IntType, IntType, IntType, IntType, IntType},

		// updateBondHolder is used by the issuer during the HolderSelection stage to set
		// a potential bond holder.
		// Parameter Required: bondAddress string, holderAddress string
		// bondAddress => Defines the address of the bond in question.
		// holderAddress => Defines the address of potential holder choosen.
		"updateBondHolder": {StringType, StringType},

		// updateBondStatus is used to move the bond along the supported bond status
		// stages.
		// Parameter Required: bondAddress string, status int
		// bondAddress => Defines the address of the bond in question.
		// status => Defines the bond stages in its lifecycle
		// 	 		status: 0 => represents Negotiating
		// 	 		status: 1 => represents HolderSelection
		// 	 		status: 2 => represents BondInDispute
		// 	 		status: 3 => represents TermsAgreement
		// 	 		status: 4 => represents ContractSigned
		// 	 		status: 5 => represents BondReselling
		// 	 		status: 6 => represents BondFinalised
		"updateBondStatus": {StringType, IntType},
	}

	// LocalMethods is a mapping of the supported locally implemented methods with
	// their respective parameter types and count expected. Parameter types are
	// placed at the position they are expected to be.
	LocalMethods = map[string][]ValueType{
		// getBondByAddress returns a bond at any status if the request sender is
		// also the bond issuer otherwise only returns bond with status
		// Negotiating.
		// Parameter Required: bondAddress string
		// bondAddress => Defines the address of the bond in question.
		"getBondByAddress": {StringType},

		// getBondsByStatus returns all the bonds status Negotiating or only
		// bonds owned by the sender if any other status is used.
		// Parameter Required: status int
		// status => Defines the bond stages in its lifecycle
		// 	 		status: 0 => represents Negotiating
		// 	 		status: 1 => represents HolderSelection
		// 	 		status: 2 => represents BondInDispute
		// 	 		status: 3 => represents TermsAgreement
		// 	 		status: 4 => represents ContractSigned
		// 	 		status: 5 => represents BondReselling
		// 	 		status: 6 => represents BondFinalised
		"getBondsByStatus": {IntType},
	}

	// ServerKeyMethod defines the method used to query the server keys
	ServerKeyMethod = map[string][]ValueType{
		// getServerPubKey is used to query the session's server public key.
		// Parameter Required: clientPubkey string
		// The client provides its public key and in return the server sends
		// back its public key. Using diffie-hellman, a sharedkey developed
		// is used to communicate securely between the client and the server.
		"getServerPubKey": {StringType},
	}
)

// GetMethodParams returns the parameters of the method provided if supported.
func GetMethodParams(method string) (implementation MethodType, param []ValueType) {
	// contract implemented methods
	if data, ok := ContractMethods[method]; ok {
		return ContractType, data
	}

	// Local methods
	if data, ok := LocalMethods[method]; ok {
		return LocalType, data
	}

	// Server keys method
	if data, ok := ServerKeyMethod[method]; ok {
		return ServerKeyType, data
	}

	return UnknownType, nil
}

// GetParamType returns the type of the param passed.
func GetParamType(param interface{}) ValueType {
	switch param.(type) {
	case string:
		return StringType
	case int:
		return IntType
	case float32, float64:
		return FloatType
	default:
		return UnsupportedType
	}
}
