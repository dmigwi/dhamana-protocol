// SPDX-License-Identifier: GPL-3.0
pragma solidity ^0.8.10;

/// @author dmigwi: (migwindungu0@gmail.com) @2023
/// @title Bond Contract.
contract BondContract {

    /// @param StatusChoice describes the various stages a bond must pass through.
    /// @param Negotiating => Open to the public, potential bond holder can their
    ///              counter proposal hoping to be selected.
    /// @param HolderSelection => The issuer select a potential holder from all
    ///              of those who have expessed interest in taking up the bond.
    /// @param TermsAgreement => Both the issuer and the selected holder agree
    ///               on the final terms of the bond. Only the two parties
    ///               have access to the bond. Security and Appendix section
    ///               can be added at this stage.
    /// @param BondInDispute => signals that a trusted entity between the parties
    ///                to the bond should lead the meditation process. 
    ///                Once in this stage, all parties involved must approval any
    ///                further status change.
    /// @param BondReselling => The bond holder can re-advertise his bond to
    ///                anyone else
    /// @param BondFinalised => The issuer has fulfilled his obligation to pay
    ///            all the amount in full as agreed on in the terms
    enum StatusChoice { 
            Negotiating, HolderSelection, TermsAgreement,
            BondInDispute, BondReselling, BondFinalised 
        }

    /// @param @param CurrencyType defines the various types of currency types 
    /// supported in the bond declaration.
    /// @param usd => represents the fiat type.
    /// @param btc => represents Bitcoin.
    /// @param eth => represents Ethereum.
    /// @param etc => represents Ethereum Classic.
    /// @param xrp => represents Ripple
    /// @param usdt => represents Tether Coin.
    /// @param dcr => represents Decred Coin.
    enum CurrencyType { usd, btc, eth, etc, xrp, usdt, dcr }

    // issuerKey defines the key used to identify data signed by the bond issuer.
    string issuerKey = "0x00000000f937";

    // holderKey defines the key used to identify data signed by the bond holder.
    string holderKey = "0x000000017d29";

    // signedBondStatuses holds the list of all bond statuses signed (approved)
    // by either of current party to this bond. Signed status to previous parties
    // are removed immediately they are detected.
    mapping(bytes32 signature => uint256 timestamp) signedBondStatuses;

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

    // Once created the issuer address cannot be changed. Accepts a sender address
    // because the default msg.sender points to the caller contract.
    constructor(address _sender) {
        bond.issuer = payable(_sender);
        bond.status = StatusChoice.Negotiating;
    }

    // onlyIssuerAllowed restricts bond changes edits to only the issuer.
    modifier onlyIssuerAllowed(address _sender) {
       require (
            _sender == bond.issuer,
            "Only the bond issuer can introduce changes"
        );
        _;
    }

    // bondAlreadyFinalised checks if checks if the Bond has been finalised. 
    // Once finalised no more changes are accepted.
    modifier bondAlreadyFinalised {
        require(
            bond.status != StatusChoice.BondFinalised,
            "Bond already finalised"
        );
        _;
    }

    // bondDetailsInDispute checks if the current bond is being disputed between
    // the parties involved. This prevents further actions on the bond until
    // they reconcile by signing the BondInDispute status.
    modifier bondDetailsInDispute {
        require(
            bond.status == StatusChoice.BondInDispute && !signaturesExists(),
            "Bond disputes not resolved by all the parties"
        );
        _;
    }

    // setBodyInfo updates the body information. Only the bond issuer can make
    // this change. 
    function setBodyInfo(
        CurrencyType _currency, uint32 _principal, uint8 _couponRate,
        uint32 _couponDate, uint32 _maturityDate, address _sender
    ) public onlyIssuerAllowed(_sender) bondDetailsInDispute bondAlreadyFinalised {
        bond.currency = _currency;
        bond.principal = _principal;
        bond.couponRate = _couponRate;
        bond.couponDate = _couponDate;
        bond.maturityDate = _maturityDate;
    }

    // setStatus sets the bond status. Can be triggered by both the issuer and the
    // the holder.
    function setStatus(StatusChoice _status) public bondDetailsInDispute bondAlreadyFinalised {
        // past the negotiating stage, a holder must be set.
        require(
            _status != StatusChoice.Negotiating && bond.issuer == address(0), 
            "Missing bond holder address"
        );

        // past holder selection stage, all bond body information must be set. 
        require(
            uint8(_status) >= uint8(StatusChoice.HolderSelection) && isBondBodySet(),
            "Bond body fields may contain empty values"
        );

        // removes the previous signatures if they existed.
        if (_status == StatusChoice.BondInDispute) {
            deleteSignatures();
        }
       
        bond.status = _status;
    }

    // setBondHolder sets the bond holder after negotiations are complete.
    // Once bond terms are agreed, the bond issuer cannot change the
    function setBondHolder(address payable _holder) public bondDetailsInDispute bondAlreadyFinalised {
        require(bond.issuer != _holder, "Bond issuer cannot be the holder too");
        require(bond.status == StatusChoice.Negotiating, "Cannot set holder during negotiations");

        bond.holder = _holder;
    }

    // setIntro sets the introduction description. Only the issuer can do this.
    function setIntro(
        string memory _intro, address _sender
    ) public bondDetailsInDispute onlyIssuerAllowed(_sender) bondAlreadyFinalised {
        bond.intro = _intro;
    }

    // setSecurity sets the security information. Only the issuer can do this. 
    function setSecurity(
        string memory _security, address _sender
    ) public bondDetailsInDispute onlyIssuerAllowed(_sender) bondAlreadyFinalised {
        bond.security = _security;
    }

    // setAppendix sets the Appendix information. Only the issuer can do this.
    function setAppendix(
        string memory _appendix, address _sender
    ) public bondDetailsInDispute onlyIssuerAllowed(_sender) bondAlreadyFinalised {
        bond.appendix = _appendix;
    }

    // isBondBodySet confirms that all the bond body information has been set
    // with non-empty values.
    function isBondBodySet() internal view  returns (bool) {
        return bond.principal > 0 && bond.couponRate > 0 &&
            bond.couponDate > 0 && bond.maturityDate > 0;
    }

    // signBondStatus returns the signed status by the sender. Can only be called 
    // from outside this contract since it sets the signed status signature 
    // against its block timestamp.
    function signBondStatus(address _sender) external {
        require (
            _sender != bond.issuer && _sender == bond.holder,
            "Unknown bond status signer address used"
        );

        bytes32 signature = bondStatusSignature(_sender, bond.status);
        signedBondStatuses[signature] = block.timestamp;
    }

    // bondStatusSignature does the actual status signing. Can only be accessed from
    // this contract.
    function bondStatusSignature(address _sender, StatusChoice _status) internal view returns (bytes32) {
         string memory encodingKey = "";
        if (_sender == bond.issuer) {
            encodingKey = issuerKey;
        } else {
            encodingKey = holderKey;
        }

        return keccak256(abi.encodeWithSignature(encodingKey, _sender, uint8(_status)+1));
    }

    // signaturesExists confirms if the required bond status has been signed
    // by all the parties involved.
    function signaturesExists() view internal returns (bool) {
        bytes32 issuerSignature = bondStatusSignature(bond.issuer, bond.status);
        bytes32 holderSignature = bondStatusSignature(bond.holder, bond.status);

        uint256 issuerSigningTime = signedBondStatuses[issuerSignature];
        uint256 holderSigningTime = signedBondStatuses[holderSignature];

        return issuerSigningTime > 0 && holderSigningTime > 0;
    }

    // deleteSignatures removes the previous signatures so that fresh ones
    // can be added. 
    function deleteSignatures() internal {
        bytes32 issuerSignature = bondStatusSignature(bond.issuer, bond.status);
        bytes32 holderSignature = bondStatusSignature(bond.holder, bond.status);

        signedBondStatuses[issuerSignature] = 0;
        signedBondStatuses[holderSignature] = 0;
    }
}