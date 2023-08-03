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
    event NewBondCreated(address sender, address bondAddress, uint timestamp);

    // FinalBondTerms creates a new bonds with final terms once the terms are
    // agreed upon by both parties.
    event FinalBondTerms(uint32 principal, uint8 couponRate, uint32 couponDate,
        uint32 maturityDate, BondContract.CurrencyType currency);

    // NewChatMessage creates a new event when a new message as part of the
    // negotiation chat is received.
    event NewChatMessage(address sender);

    // BondUnderDispute creates an event to mark the specified bond is under
    // dispute. This information is exposed to the world so as to prevent
    // malicious characters from misuing the bond system.
    event BondUnderDispute(address sender, address bondAddress);

    // BondDisputeResolved creates an event showing the sender has more disputes
    // to resolve in the said bond.
    event BondDisputeResolved(address sender, address bondAddress);

    // Creates a bonds mapping to their contract address.
    mapping (address => BondContract) private bonds;

    struct MessageInfo {
        address sender;         // Address of the message sender.
        string message;         // Actual encrypted message sent.
        uint256 timestamp;      // Time when message was received.
    }

    // conversation defines an array of messages info sent during bond negotiation.
    MessageInfo[] private conversation;

    // createBond creates a new bond associated with the user who calls it.
    function createBond() external {
        BondContract bond = new BondContract(msg.sender);
        address bondAddress = address(bond);

        // Append the new contract created.
        bonds[bondAddress] = bond;

        emit NewBondCreated(msg.sender, bondAddress, uint256(block.timestamp));
    }

    function updateBodyInfo(
        address _contract, uint32 _principal, uint8 _couponRate, uint32 _couponDate,
        uint32 _maturityDate, BondContract.CurrencyType _currency
    ) external {
        BondContract bond = bonds[_contract];

        bond.setBodyInfo(_currency, _principal, _couponRate, _couponDate, _maturityDate, msg.sender);
    }

    function updateBondStatus(address _contract, BondContract.StatusChoice _status) external {
        BondContract bondC = bonds[_contract];

        bondC.setStatus(_status);

        // Shares this publicly to prevent malicious freezing of the bond.
        if (_status == BondContract.StatusChoice.BondInDispute) {
            emit BondUnderDispute(msg.sender, _contract);
        }

        // If the terms have been agreed upon create an event displaying the
        // bond body information.
        if (_status == BondContract.StatusChoice.TermsAgreement) {
            (
                ,,,, uint32 principal, uint8 couponRate, uint32 couponDate,
                uint32 maturityDate, BondContract.CurrencyType currency,,
            ) = bondC.bond();

            emit FinalBondTerms(principal, couponRate, couponDate, maturityDate, currency);
        }
    }

    // addMessage handles the messages received that finally make up the bond.
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
            conversation.push(MessageInfo({sender: msg.sender, message: _message, timestamp: block.timestamp}));

            emit NewChatMessage(msg.sender);
        }

        // Only the bond issuer can make this encrypted messages edits below.

        if (_tag == MessageTag.Introduction) {
            bondC.setIntro(_message, msg.sender);
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
    }

    // signBondStatus allows the parties involved to sign the current bond status.
    function signBondStatus(address _contract) external {
        BondContract bondC = bonds[_contract];

        bondC.signBondStatus(payable(msg.sender));

        (,,, BondContract.StatusChoice status,,,,,,,) = bondC.bond();

        // If status signed is the BondInDispute emit its event.
        if (status == BondContract.StatusChoice.BondInDispute) {
            emit BondDisputeResolved(msg.sender, _contract);
        }
    }
}