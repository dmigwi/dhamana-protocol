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
        await this.chat.updateBondHolder(this.contractAddr, accounts[2]);
        await truffleAssert.passes(
            this.chat.updateBodyInfo(this.contractAddr, principal, couponRate,
                couponDate, maturityDate, currency)
        );

        await this.chat.updateBondStatus(this.contractAddr, 4); // 4 is index of BondReselling status

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
                ev.currency == currency
        });
    });


    it("should revert if BondFinalised status is currently set", async () => {
        await this.chat.updateBondHolder(this.contractAddr, accounts[2]);
        await this.chat.updateBodyInfo(this.contractAddr, principal, couponRate,
            couponDate, maturityDate, currency)

        // mark the current bond as finalised. 
        await truffleAssert.passes(this.chat.updateBondStatus(this.contractAddr, 5));

        // finalised bond should not accept any more status changes.
        await truffleAssert.reverts(
            this.chat.updateBondStatus(this.contractAddr, 1),
            "Bond already finalised"
        );
    });
});

// contract("Test add message to the bond chat", (accounts) => {
//     const owner = accounts[0];
//     beforeEach(
//         async () => {
//             this.chat = await chatContract.new({from: owner});

//             let data = await this.chat.createBond();
//             truffleAssert.eventEmitted(data, 'newBondCreated');
//         }
//     );

    
//     it("should only allow potential bond holders to only send messages during negotiation stage", async (accounts) => {
//         const potentialBondHolder1 = accounts[5];

//     })

// });

// contract("Test update the bond holder address", () => {

// });