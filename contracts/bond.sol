// SPDX-License-Identifier: ISC
pragma solidity ^0.8.13;

/// @author dmigwi: (migwindungu0@gmail.com) @2023
/// @title Bond Contract.
contract BondContract {

    /// @param StatusChoice describes the various stages a bond must pass through.
    /// @param Negotiating => Open to the public, potential bond holder can their
    ///              counter proposal hoping to be selected.
    /// @param HolderSelection => The issuer select a potential holder from all
    ///              of those who have expessed interest in taking up the bond.
    /// @param BondInDispute => signals that a trusted entity between the parties
    ///                to the bond should lead the meditation process.
    ///                Once in this stage, all parties involved must approval any
    ///                further status change.
    /// @param TermsAgreement => Both the issuer and the selected holder agree
    ///               on the final terms of the bond. Only the two parties
    ///               have access to the bond. Security and Appendix section
    ///               can be added at this stage.
    /// @param ContractSigned => Once the final version of the terms are agreed, the final
    ///                 document is hashed with the holders keys. Unless disputes
    ///                 the bond might spend most of its life time on this status
    ///                 until its finalised
    /// @param BondReselling => The bond holder can re-advertise his bond to
    ///                anyone else.
    /// @param BondFinalised => The issuer has fulfilled his obligation to pay
    ///            all the amount in full as agreed on in the terms
    enum StatusChoice {
            Negotiating, HolderSelection, BondInDispute, TermsAgreement,
                ContractSigned, BondReselling, BondFinalised
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
    mapping(bytes32 => uint256) signedBondStatuses;

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

    // sapphire requires public variables to have an explicitly defined getter
    // so as to guarrantee privacy.
    Bond public bond;

    // Once created the issuer address cannot be changed. Accepts a sender address
    // because the default msg.sender points to the caller contract.
    constructor(address _sender) {
        bond.issuer = payable(_sender);
        bond.status = StatusChoice.Negotiating;
    }

    // getBond satisfies the requirement to have an explicitly defined getter
    // method for all public variables declared.
    function getBond() external view returns (Bond memory) {
        return bond;
    }

    // OnlyIssuerBondEditsAllowed is used to restricts bond changes to only be
    // introduced by the issuer on select functionality.
    error OnlyIssuerBondEditsAllowed();

    // ActivitiesDisabledOnFinalisedBond confirms that no more changes/activities
    // are allowed from any person, including the the holder and issuer.
    error ActivitiesDisabledOnFinalisedBond();

    // BondDetailsInDispute is used to show the current bond is in disputed between
    // the parties involved. This prevents further actions on the bond until
    // they reconcile by signing the BondInDispute status.
    error BondDetailsInDispute();

    // TermsUpdateDisabled prevents further terms edit since the version of the
    // terms have been agreed upon.
    // Terms edidt is only possible fully or partially only on this bond statuses only:
    // Negotiating, HolderSelection, BondInDispute, TermsAgreement
    error TermsUpdateDisabled();

    // MissingBondHolderAddress indicates that the bond holder address was expected
    // but it not yet set.
    error MissingBondHolderAddress();

    // EmptyBondBodyFieldsFound indicates that empty bond fields were not expected.
    error EmptyBondBodyFieldsFound();

    // BondIssuerCannotBeAHolder returned if the bond issuer tries to assign themselves
    // as a holder of their bond.
    error BondIssuerCannotBeAHolder();

    // HolderOnlySetOnHolderselection prevents reassignment of the holder unless
    // its on the Holderselection bond status.
    error HolderOnlySetOnHolderselection();

    // UnknownBondStatusSigner indicates the signer address used neither belongs
    // to the bond issuer or the holder thus not allowed to sign a status.
    error UnknownBondStatusSigner();

    // setBodyInfo updates the body information. Only the bond issuer can make
    // this change.
    function setBodyInfo(
        CurrencyType _currency, uint32 _principal, uint8 _couponRate,
        uint32 _couponDate, uint32 _maturityDate, address _sender
    ) public {
        if (_sender != bond.issuer) revert OnlyIssuerBondEditsAllowed();

        if (bond.status == StatusChoice.BondInDispute && !signaturesExists()) {
            revert BondDetailsInDispute();
        }

        if (uint(bond.status) > uint(BondContract.StatusChoice.TermsAgreement)) {
            revert TermsUpdateDisabled();
        }

        bond.currency = _currency;
        bond.principal = _principal;
        bond.couponRate = _couponRate;
        bond.couponDate = _couponDate;
        bond.maturityDate = _maturityDate;
    }

    // setStatus sets the bond status. Can be triggered by both the issuer and the
    // the holder.
    function setStatus(StatusChoice _status) public {
        if (bond.status == StatusChoice.BondInDispute && !signaturesExists()) {
            revert BondDetailsInDispute();
        }

        if (bond.status == StatusChoice.BondFinalised) {
            revert ActivitiesDisabledOnFinalisedBond();
        }

        // past the negotiating stage, a holder must be set.
        if (uint8(_status) > uint8(StatusChoice.HolderSelection) && bond.holder == address(0)) {
            revert MissingBondHolderAddress();
        }

        // past holder selection stage, all bond body information must be set.
        if (uint8(_status) > uint8(StatusChoice.HolderSelection) && isAnyBondBodyFieldEmpty()) {
            revert EmptyBondBodyFieldsFound();
        }

        // If setting the new status to ContractSigned, TermsAgreement must be
        // signed by both parties. To set ContractSigned status, the previous
        // status must be TermsAgreement and have it signed.
        if (_status == StatusChoice.ContractSigned) {
            if (uint(bond.status) > uint(BondContract.StatusChoice.TermsAgreement)) {
                revert TermsUpdateDisabled();
            }
        }

        // Removes the prev signatures for BondInDispute and TermsAgreement statuses if they exist.
        if (_status == StatusChoice.BondInDispute || _status == StatusChoice.TermsAgreement) {
            deleteSignatures();
        }

        bond.status = _status;
    }

    // setBondHolder sets the bond holder after negotiations are complete.
    // Once selected, the bond issuer cannot replace the holder.
    // only the bond issuer can select a bond holder from the interested people.
    function setBondHolder(address payable _holder, address payable _sender) public {
        if (_sender != bond.issuer) revert OnlyIssuerBondEditsAllowed();

        if (bond.status == StatusChoice.BondInDispute && !signaturesExists()) {
            revert BondDetailsInDispute();
        }

        if (bond.issuer == _holder) revert BondIssuerCannotBeAHolder();

        if (uint8(bond.status) != uint8(StatusChoice.HolderSelection)) {
            revert HolderOnlySetOnHolderselection();
        }

        bond.holder = _holder;
    }

    // setIntro sets the introduction description. Only the issuer can do this.
    function setIntro( string memory _intro, address _sender ) public {
        if (_sender != bond.issuer) revert OnlyIssuerBondEditsAllowed();

        if (bond.status == StatusChoice.BondInDispute && !signaturesExists()) {
            revert BondDetailsInDispute();
        }

        if (uint(bond.status) > uint(BondContract.StatusChoice.TermsAgreement)) {
                revert TermsUpdateDisabled();
        }

        bond.intro = _intro;
    }

    // setSecurity sets the security information. Only the issuer can do this.
    function setSecurity(string memory _security, address _sender) public {
        if (_sender != bond.issuer) revert OnlyIssuerBondEditsAllowed();

        if (bond.status == StatusChoice.BondInDispute && !signaturesExists()) {
            revert BondDetailsInDispute();
        }

        if (uint(bond.status) > uint(BondContract.StatusChoice.TermsAgreement)) {
                revert TermsUpdateDisabled();
        }

        bond.security = _security;
    }

    // setAppendix sets the Appendix information. Only the issuer can do this.
    function setAppendix(string memory _appendix, address _sender) public {
        if (_sender != bond.issuer) revert OnlyIssuerBondEditsAllowed();

        if (bond.status == StatusChoice.BondInDispute && !signaturesExists()) {
            revert BondDetailsInDispute();
        }

        if (uint(bond.status) > uint(BondContract.StatusChoice.TermsAgreement)) {
                revert TermsUpdateDisabled();
        }

        bond.appendix = _appendix;
    }

    // isAnyBondBodyFieldEmpty returns true if either of; principal, couponRate,
    // couponDate or maturityDate are empty.
    function isAnyBondBodyFieldEmpty() internal view  returns (bool) {
        return bond.principal == 0 || bond.couponRate == 0 ||
            bond.couponDate == 0 || bond.maturityDate == 0;
    }

    // signBondStatus returns the signed status by the sender. Can only be called
    // from outside this contract since it sets the signed status signature
    // against its block timestamp.
    function signBondStatus(address payable _sender) external {
        if (_sender != bond.issuer && _sender != bond.holder) {
            revert UnknownBondStatusSigner();
        }

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