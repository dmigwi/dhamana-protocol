// SPDX-License-Identifier: GPL-3.0
pragma solidity ^0.8.10;

import "./bond.sol";

/// @author dmigwi: (migwindungu0@gmail.com) @2023
/// @title Chat Contract.
contract chatContract {

    /// @param sectionTag describes the various message types are expected from all
    /// people interacting with the bond via chat messages. 
    /// @param InitConversation => describes the bond conversation during the
    ///                      negotiation stage for potential bond holders, and
    ///                      messages by the issuer till the terms are agreed upon. 
    /// @param Introduction => describes the message sent by the issuer describing their
    ///                  motivation to issue a bond.
    /// @param Security => describes the message sent by the issuer describing the
    ///              collateral the issuer is willing to commit as proof of them
    ///              honoring their end of the deal.
    enum sectionTag { InitConversation, Introduction, Security, Appendix }

    // newBondCreated creates new event showing the a new contract has been created.
    event newBondCreated(address contractAddress, uint timestamp);

    // finalBondTerms creates a new bonds with final terms once the terms are
    // agreed upon by both parties.
    event finalBondTerms(uint32 principal, uint8 couponRate, uint32 couponDate,
        uint32 maturityDate);

    // newChatMessage creates a new event when a new message as part of the
    // negotiation chat is received.
    event newChatMessage(address _sender);

    // Creates a bonds mapping to their contract address.
    mapping (address => BondContract) bonds;

    struct messageInfo {
        address sender;   // Address of the message sender.
        string message;   // Actual encrypted message sent.
        uint timestamp;   // Time when message was received.
    }

    // conversation defines an array of messages info sent during bond negotiation. 
    messageInfo[] conversation;

    // createBond creates a new bond associated with the user who calls it.
    function createBond() external {
        BondContract bond = new BondContract();
        address bondAddress = address(bond);

        // Append the new contract created. 
        bonds[bondAddress] = bond;

        emit newBondCreated(bondAddress, block.timestamp);
    }

    function updateBodyInfo(
        address _contract, BondContract.CurrencyType _currency,
        uint32 _principal, uint8 _couponRate, uint32 _couponDate,
        uint32 _maturityDate
    ) public {
        BondContract bond = bonds[_contract];
        require(bond != new BondContract(), "invalid contract address provided");

        bond.setBodyInfo(_currency, _principal, _couponRate, _couponDate, _maturityDate);
    }

    function updateBondStatus(address _contract, BondContract.StatusChoice _status) public {
        BondContract bondC = bonds[_contract];
        require(bondC != new BondContract(), "invalid contract address provided");

        bondC.setStatus(_status);

        // If the terms have been agreed upon create an event displaying the 
        // bond body information.
        if (_status == BondContract.StatusChoice.TermsAgreement) {
            (,,,,uint32 principal,uint8 couponRate, uint32 couponDate,uint32 maturityDate,,,) = bondC.bond();

            emit finalBondTerms(principal, couponRate,couponDate,maturityDate);
        }
    }

    // addMessage handles the messages received that finally make up the bond.
    function addMessage(address _contract, sectionTag _tag, string memory _message) public {
        BondContract bondC = bonds[_contract];
        require(bondC != new BondContract(), "invalid contract address provided");

        (,address payable issuer,, BondContract.StatusChoice status,,,,,,,) = bondC.bond();

        // Potential bond holders can only sent messages at negotiating stage.  
        require (
            status != BondContract.StatusChoice.Negotiating && issuer != msg.sender,
            "Only bond issuers can comment past negotiating stage"
        );

        // Bond issues can only make changes till the TermsAgreement stage. Past here
        // the contract is considered binding and cannot be editted.
        require (
            uint(status) > uint(BondContract.StatusChoice.TermsAgreement),
            "no more messages are allowed"
        );

        // Messages not part of the negotiations shouldn't get to the chat.
        if (_tag == sectionTag.InitConversation) {
            conversation.push(messageInfo({sender: msg.sender, message: _message, timestamp: block.timestamp}));

            emit newChatMessage(msg.sender);
        }

        if (issuer != msg.sender) {
            return;  // The bond issue did not send the current message, ignore further update.
        }

        if (_tag == sectionTag.Introduction) {
            bondC.setIntro(_message);
        } else if (_tag == sectionTag.Security) {
            bondC.setSecurity(_message);
        } else if (_tag == sectionTag.Appendix) {
            bondC.setAppendix(_message);
        }
    }

    // updateBondholder sets the bond holder address after the negotiation stage
    // is complete.
    function updateBondholder(address _contract, address payable _holder) public {
        BondContract bondC = bonds[_contract];
        require(bondC != new BondContract(), "invalid contract address provided");

        bondC.setBondHolder(_holder);
    }
}