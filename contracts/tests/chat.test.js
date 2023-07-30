const chatContract = artifacts.require("../chat/chatContract");
const truffleAssert = require("truffle-assertions");
const web3 = require("web3");

const currency = 2;
const principal = 12323;
const couponRate = 2;
const couponDate = 1690670000;
const maturityDate = 1690869606;
const chatMsg = "This text is a placeholder of an encrypted message";

contract("Test create a new Bond", (accounts) => {
    it("should receive a newBondCreated event", async () => {
        const owner = accounts[0];
        const chat = await chatContract.new({from: owner});
        let data = await chat.createBond(); 

        truffleAssert.eventEmitted(data, 'newBondCreated', (ev) => { return true });
    });
});

contract("Test update bond body information", (accounts) => {
    const owner = accounts[0];
    beforeEach(
        async () => {
            this.chat = await chatContract.new({from: owner});

            let data = await this.chat.createBond();
            truffleAssert.eventEmitted(data, 'newBondCreated');
    
            this.contractAddr = data.logs[0].args.contractAddress;
        }
    );
    
    it("should revert if the sender detected is not the bond issuer", async () => {
        const fakeOwner = accounts[4];
        await truffleAssert.reverts(
            this.chat.updateBodyInfo(this.contractAddr, principal, couponRate,
                couponDate, maturityDate, currency, {from: fakeOwner}),
            "Only the bond issuer can introduce changes"
        );
    })

    it("should not return a revert when setting status during negotations or holderselection", async () => {
        await this.chat.updateBondStatus(this.contractAddr, 1); // 1 is index of HolderSelection
        await truffleAssert.passes(
            this.chat.updateBodyInfo(this.contractAddr, principal, couponRate,
                couponDate, maturityDate, currency)
        );
    });

    it("should revert if editting is attempted when the bond is dispute", async () => {
        await this.chat.updateBondStatus(this.contractAddr, 1); // 1 is index of HolderSelection
        await this.chat.updateBondHolder(this.contractAddr, accounts[2]);
        await truffleAssert.passes(
            this.chat.updateBodyInfo(this.contractAddr, principal, couponRate,
                couponDate, maturityDate, currency)
        );

        await this.chat.updateBondStatus(this.contractAddr, 2); // 2 is index of BondInDispute status

        await truffleAssert.reverts(
            this.chat.updateBodyInfo(this.contractAddr, principal, couponRate,
                couponDate, maturityDate, currency),
            "Bond disputes not yet resolved by all the parties"
        );
    });

    it("should revert if terms update are initiated past TermsAgreement status", async () => {
        await this.chat.updateBondStatus(this.contractAddr, 1); // 1 is index of HolderSelection
        await this.chat.updateBondHolder(this.contractAddr, accounts[2]);
        await truffleAssert.passes(
            this.chat.updateBodyInfo(this.contractAddr, principal, couponRate,
                couponDate, maturityDate, currency)
        );

        await this.chat.updateBondStatus(this.contractAddr, 5); // 5 is index of BondReselling status

        await truffleAssert.reverts(
            this.chat.updateBodyInfo(this.contractAddr, principal, couponRate,
                couponDate, maturityDate, currency),
            "Bond terms update is disabled"
        );
    });
});

contract("Test update bond status", (accounts) => {
    const owner = accounts[0];
    beforeEach(async () => {
        this.chat = await chatContract.new({from: owner});

        let data = await this.chat.createBond();
        truffleAssert.eventEmitted(data, 'newBondCreated');

        this.contractAddr = data.logs[0].args.contractAddress;
    });

    it("should revert if status is to be set past holderselection before setting the holder's address", async () => {
        await truffleAssert.passes(
            this.chat.updateBodyInfo(this.contractAddr, principal, couponRate,
                couponDate, maturityDate, currency)
        );

        // dispute cannot be initiated without two people.
        await truffleAssert.reverts(
            this.chat.updateBondStatus(this.contractAddr, 2), // 2 is index of BondInDispute status
            "Missing bond holder address"
        );
    });

    it("should revert if status update past holderselection with some missing body fields is triggered", async () => {
        await this.chat.updateBondStatus(this.contractAddr, 1); // 1 is index of HolderSelection
        await this.chat.updateBondHolder(this.contractAddr, accounts[2]);
        await truffleAssert.passes(
            this.chat.updateBodyInfo(this.contractAddr, principal, 0,
                couponDate, 0, currency)
        );

        await truffleAssert.reverts(
            this.chat.updateBondStatus(this.contractAddr, 3), // 3 is index of TermsAgreement status
            "Bond body fields may contain empty values"
        );
    });

    it("should revert if bond dispute is yet to be resolved", async () => {
        await this.chat.updateBondStatus(this.contractAddr, 1); // 1 is index of HolderSelection
        await truffleAssert.passes(this.chat.updateBondHolder(this.contractAddr, accounts[2]));
        await this.chat.updateBodyInfo(this.contractAddr, principal, couponRate,
            couponDate, maturityDate, currency)

        // mark the current bond as under dispute. 
        let txInfo = await this.chat.updateBondStatus(this.contractAddr, 2); // 2 is index of BondInDispute status

        truffleAssert.eventEmitted(txInfo, 'bondUnderDispute', (ev) => {
            return ev.sender == owner &&  ev.bondAddress == this.contractAddr
        });

        // Attempting to change status without resolving the dispute first should fail.
        await truffleAssert.reverts(
            this.chat.updateBondStatus(this.contractAddr, 3), // 3 is index of TermsAgreement status
            "Bond disputes not yet resolved by all the parties"
        );
    });

    it('should not revert if a bond previously under dispute is resolved and status updated', async () => {
        await this.chat.updateBondStatus(this.contractAddr, 1); // 1 is index of HolderSelection
        await truffleAssert.passes(this.chat.updateBondHolder(this.contractAddr, accounts[2]));
        await this.chat.updateBodyInfo(this.contractAddr, principal, couponRate,
            couponDate, maturityDate, currency)

        // mark the current bond as under dispute. 
        let txInfo = await this.chat.updateBondStatus(this.contractAddr, 2);

        truffleAssert.eventEmitted(txInfo, 'bondUnderDispute', (ev) => {
            return ev.sender == owner &&  ev.bondAddress == this.contractAddr
        });

        // Issuer resolves the dispute on his end.
        let issuerData = await this.chat.signBondStatus(this.contractAddr, {from: owner});
        truffleAssert.eventEmitted(issuerData, 'bondDisputeResolved', (ev) => {
            return ev.sender == owner && ev.bondAddress == this.contractAddr
        });

        // Attempting status change should still fail.
        await truffleAssert.reverts(
            this.chat.updateBondStatus(this.contractAddr, 3),
            "Bond disputes not yet resolved by all the parties"
        );

        // holder resolves the dispute on his end.
        let holderData = await this.chat.signBondStatus(this.contractAddr, {from: accounts[2]});
        truffleAssert.eventEmitted(holderData, 'bondDisputeResolved', (ev) => {
            return ev.sender == accounts[2] && ev.bondAddress == this.contractAddr;
        });

        // Status change should succeed now.
        await truffleAssert.passes(this.chat.updateBondStatus(this.contractAddr, 3));
    });

    it("should allow emit finalBondTerms event on TermsAgreement status", async () => {
        await this.chat.updateBondStatus(this.contractAddr, 1); // 1 is index of HolderSelection
        await this.chat.updateBondHolder(this.contractAddr, accounts[2]);
        await truffleAssert.passes(
            this.chat.updateBodyInfo(this.contractAddr, principal, couponRate,
                couponDate, maturityDate, currency)
        );

        // Set the TermsAgreement bond status.
       let newdata =  await this.chat.updateBondStatus(this.contractAddr, 3);

        truffleAssert.eventEmitted(newdata, 'finalBondTerms', (ev) => {
            return ev.principal == principal && ev.couponRate == couponRate &&
                ev.couponDate == couponDate && ev.maturityDate == maturityDate && 
                ev.currency == currency;
        });
    });

    it("should revert if BondFinalised status is currently set", async () => {
        await this.chat.updateBondStatus(this.contractAddr, 1); // 1 is index of HolderSelection
        await this.chat.updateBondHolder(this.contractAddr, accounts[2]);
        await this.chat.updateBodyInfo(this.contractAddr, principal, couponRate,
            couponDate, maturityDate, currency);

        // mark the current bond as finalised. 
        await truffleAssert.passes(this.chat.updateBondStatus(this.contractAddr, 6));

        // finalised bond should not accept any more status changes.
        await truffleAssert.reverts(
            this.chat.updateBondStatus(this.contractAddr, 1),
            "Bond already finalised"
        );
    });
});

contract("Test add message to the bond chat", (accounts) => {
    const owner = accounts[0];
    const holder = accounts[2];

    beforeEach(async () => {
        this.chat = await chatContract.new({from: owner});

        let data = await this.chat.createBond();
        truffleAssert.eventEmitted(data, 'newBondCreated');

        this.contractAddr = data.logs[0].args.contractAddress;

        await truffleAssert.passes(this.chat.updateBondStatus(this.contractAddr, 1));
        await truffleAssert.passes(this.chat.updateBondHolder(this.contractAddr, holder));
        await this.chat.updateBodyInfo(this.contractAddr, principal, couponRate,
            couponDate, maturityDate, currency);
        
        await this.chat.updateBondStatus(this.contractAddr, 0); // revert back to the default status
    });

    const potentialBondHolder1 = accounts[4];
    const potentialBondHolder2 = accounts[5];

    const initChat = 0;
    const intro = 1;
    const security = 2;
    const appendix = 3;
    
    it("should block potential bond holders from sending messages past negotiation stage", async () => {
        await this.chat.updateBondStatus(this.contractAddr, 1)
        await truffleAssert.reverts(
            this.chat.addMessage(this.contractAddr, initChat, chatMsg, {from: potentialBondHolder2}),
            "Potential holders cannot comment past negotiating stage"
        );
    });

    it("should allow potential bond holders to only send messages during negotiation stage", async () => {
        let txInfo = await this.chat.addMessage(this.contractAddr, initChat, chatMsg, {from: potentialBondHolder1})

        truffleAssert.eventEmitted(txInfo, 'newChatMessage', (ev) => {
            return ev.sender == potentialBondHolder1;
        });
    })

    it("should allow both issuer and holder to chat messages till the bond is finalised", async () => {
        await truffleAssert.passes(this.chat.updateBondStatus(this.contractAddr, 3));

        // issuer can send a message during TermsAgreement status.
        let issuerInfo = await this.chat.addMessage(this.contractAddr, initChat, chatMsg, {from: owner})

        truffleAssert.eventEmitted(issuerInfo, 'newChatMessage', (ev) => {
            return ev.sender == owner;
        });

        let holderInfo = await this.chat.addMessage(this.contractAddr, initChat, chatMsg, {from: holder})

        truffleAssert.eventEmitted(holderInfo, 'newChatMessage', (ev) => {
            return ev.sender == holder;
        });
    });

    it("should only allow the issuer to send the contract specific details messages", async () => {
        await truffleAssert.passes(this.chat.addMessage(this.contractAddr, intro, chatMsg, {from: owner}));
        await truffleAssert.passes(this.chat.addMessage(this.contractAddr, security, chatMsg, {from: owner}));
        await truffleAssert.passes(this.chat.addMessage(this.contractAddr, appendix, chatMsg, {from: owner}));

        await truffleAssert.reverts(
            this.chat.addMessage(this.contractAddr, intro, chatMsg, {from: holder}),
            "Only the bond issuer can introduce changes"
        );
        await truffleAssert.reverts(
            this.chat.addMessage(this.contractAddr, security, chatMsg, {from: holder}),
            "Only the bond issuer can introduce changes"
        );
        await truffleAssert.reverts(
            this.chat.addMessage(this.contractAddr, appendix, chatMsg, {from: holder}),
            "Only the bond issuer can introduce changes"
        );
    });

    it("should revert if issuer chats/messages are sent once the bond has been finalised", async () => {
        await this.chat.updateBondStatus(this.contractAddr, 6);

        await truffleAssert.reverts(
            this.chat.addMessage(this.contractAddr, security, chatMsg, {from: owner}),
            "No more comments are allowed after the bond is finalised"
        );
    });

    it("should revert if issuer contact messages are sent once the bond is in dispute", async () => {
        await this.chat.updateBondStatus(this.contractAddr, 2);

        await truffleAssert.reverts(
            this.chat.addMessage(this.contractAddr, security, chatMsg, {from: owner}),
            "Bond disputes not yet resolved by all the parties"
        );
    });

});

contract("Test update the bond holder address", (accounts) => {
    const owner = accounts[0];
    const holder = accounts[2];

    beforeEach(async () => {
        this.chat = await chatContract.new({from: owner});

        let data = await this.chat.createBond();
        truffleAssert.eventEmitted(data, 'newBondCreated');

        this.contractAddr = data.logs[0].args.contractAddress;

        await this.chat.updateBodyInfo(this.contractAddr, principal, couponRate,
            couponDate, maturityDate, currency);
    });

    it("should revert if the issuer attempts to set themselves as the holder", async () => {
        await this.chat.updateBondStatus(this.contractAddr, 1);
        await truffleAssert.reverts(
            this.chat.updateBondHolder(this.contractAddr, owner),
            "Bond issuer cannot be the holder too"
        );
    });

    it("should revert if a holder is not set during holderSelection status", async () => {
        await this.chat.updateBondStatus(this.contractAddr, 0);
        await truffleAssert.reverts(
            this.chat.updateBondHolder(this.contractAddr, holder),
            "Only set during the HolderSelection status"
        );
    });

    it("should allow setting of an appropriate holder at the correct status", async () => {
        await this.chat.updateBondStatus(this.contractAddr, 1);
        await truffleAssert.passes( this.chat.updateBondHolder(this.contractAddr, holder));
    });
});