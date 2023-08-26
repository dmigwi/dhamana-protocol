// SPDX-License-Identifier: ISC
pragma solidity ^0.8.13;

import { BondContract } from "./bond.sol";

/// @author dmigwi: (migwindungu0@gmail.com) @2023
/// @title Chat Contract.
contract ChatContract {

    /// @param MessageTag describes the various message types are expected from all
    ///                      people interacting with the bond via chat messages.
    /// @param InitConversation => describes the bond conversation during the
    ///                      negotiation stage for potential bond holders, and
    ///                      messages by the issuer till the terms are agreed upon.
    /// @param Introduction => describes the message sent by the issuer describing their
    ///                  motivation to issue a bond.
    /// @param Security => describes the message sent by the issuer describing the
    ///              collateral the issuer is willing to commit as proof of them
    ///              honoring their end of the deal.
    /// @param Appendix => describes any other information that is crucial to
    ///               complete. This information my include but limited to financial
    ///               financial transaction information between the two parties.
    enum MessageTag { InitConversation, Introduction, Security, Appendix }

    // NewBondCreated creates new event showing the a new contract has been created.
    event NewBondCreated(address sender, address bondAddress);

    // BondBodyTerms creates a new bonds with the terms body terms update made.
    event BondBodyTerms(address bondAddress, uint32 principal, uint8 couponRate, uint8 couponDate,
        uint32 maturityDate, BondContract.CurrencyType currency);

    // BondMotivation creates a new event that comprises of the bond introduction
    // message sent by the bond issuer.
    event BondMotivation(address bondAddress, string message);

    // Statuschange creates a new event sent every time the bond status is
    // changed by either the issuer or the holder.
    event StatusChange(address sender, address bondAddress, BondContract.StatusChoice status);

    // StatusUpdate creates a new event sent every time the issuer or the holder
    // sign the current status as a sign of consensus.
    event StatusSigned(address sender, address bondAddress, BondContract.StatusChoice status);

    // HolderUpdate creates an event when the issuer selects the potential holder
    // they would like to deal with.
    event HolderUpdate(address bondAddress, address holder);

    // Creates a bonds mapping to their contract address.
    mapping (address => BondContract) private bonds;

    // NewChatMessage creates a new event when a new message as part of the
    // negotiation chat is received. Past negotaiation stage, events will still
    // be sent but only the bond issuer and holder can send those messages.
    // message should always be encrypted depending on the information sensitivity.
    event NewChatMessage(address bondAddress,  address sender, string message);

    // createBond creates a new bond associated with the user who calls it.
    function createBond() external {
        BondContract bond = new BondContract(msg.sender);
        address bondAddress = address(bond);

        // Append the new contract created.
        bonds[bondAddress] = bond;

        emit NewBondCreated(msg.sender, bondAddress);
    }

    // updateBodyInfo updates the bond body details and emits an event.
    // Only the bond issuer allowed to make this changes.
    function updateBodyInfo(
        address _contract, uint32 _principal, uint8 _couponRate, uint8 _couponDate,
        uint32 _maturityDate, BondContract.CurrencyType _currency
    ) external {
        BondContract bond = bonds[_contract];

        bond.setBodyInfo(_currency, _principal, _couponRate, _couponDate, _maturityDate, msg.sender);

        // BondBodyTerms emits an event showing the latest update in bond body terms.
        emit BondBodyTerms(_contract, _principal, _couponRate, _couponDate, _maturityDate, _currency);
    }

    // updateBondStatus moves the bond along its various lifecycle stages (statuses).
    // For every status change an event is emitted recording the status change.
    function updateBondStatus(address _contract, BondContract.StatusChoice _status) external {
        BondContract bondC = bonds[_contract];

        bondC.setStatus(_status);

        // StatusChange event is emitted for every status change.
        emit StatusChange(msg.sender, _contract, _status);
    }

    // addMessage handles the messages received that finally make up the bond.
    // Events are emitted for the message tags apart from security and appendix
    // information accessed by the bond parties on demand.
    function addMessage(address _contract, MessageTag _tag, string memory _message) external {
        BondContract bondC = bonds[_contract];

        (,address payable issuer,address payable holder, BondContract.StatusChoice status,,,,,,,) = bondC.bond();

         // The choosen bond holder and the issuer can send messages till the bond is finalised.
        if (msg.sender != issuer && msg.sender != holder) {
            // Potential bond holders cannot comment past negotiating stage in the chat.
            // Potential is used to refer to people who have shown interest in the
            // bond via commenting in the chat but haven't been selected by the issuer.
            require (
                status == BondContract.StatusChoice.Negotiating,
                "Only negotiation chat is general"
            );
        }

        require (
            status != BondContract.StatusChoice.BondFinalised,
            "Edits disabled on finalized Bond"
        );

        // Messages not part of the negotiations shouldn't get to the chat.
        if (_tag == MessageTag.InitConversation) {
            // Chat are messages not stored because they have very little effect on the
            // final bond. Also storing them would create unnecessary data stored in
            // the contracts. The same data can be obtained from the event logs.
            emit NewChatMessage(_contract, msg.sender, _message);
        }

        // Only the bond issuer can make this encrypted messages edits below.
        if (_tag == MessageTag.Introduction) {
            bondC.setIntro(_message, msg.sender);

            // Introduction message is accessible to everyone via events but only
            // editted by the bond issuer.
            emit BondMotivation(_contract, _message);
        } else if (_tag == MessageTag.Security) {
            bondC.setSecurity(_message, msg.sender);
        } else if (_tag == MessageTag.Appendix) {
            bondC.setAppendix(_message, msg.sender);
        }
    }

    // updateBondHolder allows setting of the bond holder address after
    // the negotiation stage is complete. Only the bond hold can select the
    // a holder from the interested people.
    function updateBondHolder(address _contract, address _holder) external {
        BondContract bondC = bonds[_contract];

        bondC.setBondHolder(payable(_holder), payable(msg.sender));

        // HolderUpdate emits the event for every successful status change.
        emit HolderUpdate(_contract, _holder);
    }

    // signBondStatus allows the parties involved to sign the current bond status.
    function signBondStatus(address _contract) external {
        BondContract bondC = bonds[_contract];

        bondC.signBondStatus(payable(msg.sender));

        (,,, BondContract.StatusChoice status,,,,,,,) = bondC.bond();

        // for every successful status signed, StatusSigned event is emitted.
        emit StatusSigned(msg.sender, _contract, status);
    }

    // getBondSecureDetails returns the bonds private messages on accessible to
    // bond holder and bond issuer.
    function getBondSecureDetails(address _contract) public view returns (
        string memory _security, string memory _appendix
    ) {
        BondContract bondC = bonds[_contract];

        (,address payable issuer,address payable holder,,,,,,,
            string memory security, string memory appendix) = bondC.bond();
        require (
            !((msg.sender != issuer && msg.sender != holder)),
            "Only used by the bond parties"
        );

        _security = security;
        _appendix = appendix;
        return (_security, _appendix);
    }
}