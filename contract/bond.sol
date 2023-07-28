// SPDX-License-Identifier: GPL-3.0
pragma solidity ^0.8.0;

/// @author dmigwi: (migwindungu0@gmail.com)
/// @title Bond Contract.
contract BondContract {

    // StatusChoice describes the various stages a bond must pass through.
    // Negotiating => Open to the public, potential bond holder can their
    //                counter proposal hoping to be selected.
    // HolderSelection => The issuer select a potential holder from all of those
    //                who have expessed interest in taking up the bond.
    // TermsAgreement => Both the issuer and the selected holder agree on the
    //                 final terms of the bond. On the two parties to the bond
    //                 have access to the bond. Security and Appendix section
    //                 are added at this stage.
    // ContractSigned => Once the final version of the terms are agreed, the final
    //                 document is hashed with the holders keys.
    // BondReselling => The bond holder can re-advertise his bond to anyone else.
    // BondFinalised => The issuer has fulfilled his obligation to pay all the amount
    //              in full as agreed on in the terms
    enum StatusChoice { 
            Negotiating, HolderSelection, TermsAgreement,
            ContractSigned, BondReselling, BondFinalised 
        }

    // CurrencyType defines the various types of currency types supported in the
    // bond declaration.
    // usd => represents the fiat type.
    // btc => represents Bitcoin.
    // eth => represents Ethereum.
    // etc => represents Ethereum Classic.
    // xrp => represents Ripple
    // usdt => represents Tether Coin.
    // dcr => represents Decred Coin.
    enum CurrencyType { usd, btc, eth, etc, xrp, usdt, dcr}

    // Bond describes the collection of information that make up a bond. 
    struct Bond {
        // Section: 1 
        // intro is added by the bond issuer where they describe the motivation
        // or what they hope  to achieve once the bond is subscribed to.
        string intro;

        // Section: 2 (body)     
        // issuer describes the address of the bond issuer. It should always be
        // set in all bonds.
        address payable issuer;
        // holder describes the address of the current bond holder. Once the
        // bond is resold, the new holder address should be set here.
        address payable holder;
        // StatusChoice describes the stages between issuance and completion of
        // a bond.
        StatusChoice status;
        // principal defines the initial amount in the specified currency the
        // issuer hopes to receive.
        uint32 principal;
        // couponRate describes the percentage of the interest the issuer
        // is willing to pay at the couponDate intervals.
        uint8 couponRate;
        // couponDate decribes the interval i.e. monthly, quartely, yearly etc
        // when the interest is due to be paid to the holder.
        uint32 couponDate;
        // maturityDate decribes date when the whole amount owed to the holder
        // should be cleared/paid in full.
        uint32 maturityDate;
        // CurrencyType describes the currency which the bond issue and holder
        // wishes to transact in.
        CurrencyType currency;

        // Section: 3
        // security describe some collataral or a form of a legal action a bond
        // holder can institute against the issuer should, money in the agreed
        // amount fail to get to him at the expected date. 
        // Only exposed to the bond issuer and holder only
        string security;

        // Section: 4
        // appendix describes any other information including details on how the
        // bond issuer and holder should make financial with each other.
        string appendix;
    }

    Bond public bond;

    // Once created the issuer address cannot be changed.
    constructor() {
        bond.issuer = payable(msg.sender);
        bond.status = StatusChoice.Negotiating;
    }

    // onlyIssuerAllowed restricts bond changes edits to only the issuer.
    modifier onlyIssuerAllowed {
       require (
            msg.sender == bond.issuer,
            "Only the bond issue can make changes"
        );
        _;
    }

    // setBodyInfo updates the body information. Only the bond issuer can make
    // this change. 
    function setBodyInfo(
        CurrencyType _currency, uint32 _principal, uint8 _couponRate,
        uint32 _couponDate, uint32 _maturityDate) public onlyIssuerAllowed {
        bond.currency = _currency;
        bond.principal = _principal;
        bond.couponRate = _couponRate;
        bond.couponDate = _couponDate;
        bond.maturityDate = _maturityDate;
    }

    // setStatus sets the bond status. Can be triggered by both the issuer and the
    // the holder.
    function setStatus(StatusChoice _status) public {
        bond.status = _status;
    }

    // setBondHolder sets the bond holder after negotiations are complete.
    function setBondHolder(address payable _holder) public {
        require(msg.sender != _holder, "Bond issuer cannot be the holder too");
        require(bond.status == StatusChoice.Negotiating, "Cannot set holder during negotiations");

        bond.holder = _holder;
    }

    // setIntro sets the introduction description. Only the issuer can do this.
    function setIntro(string memory _intro) public onlyIssuerAllowed {
        bond.intro = _intro;
    }

    // setSecurity sets the security information. Only the issuer can do this. 
    function setSecurity(string memory _security) public onlyIssuerAllowed {
        bond.security = _security;
    }

    // setAppendix sets the Appendix information. Only the issuer can do this.
    function setAppendix(string memory _appendix) public onlyIssuerAllowed {
        bond.appendix = _appendix;
    }
}