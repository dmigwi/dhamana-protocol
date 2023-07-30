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
        }
    );

    const fakeOwner = accounts[4];
   
    it("should revert if the sender detected is not the bond issuer", async () => {
        let data = await this.chat.createBond();
        truffleAssert.eventEmitted(data, 'newBondCreated');

        // fetch the bond address from the event.
        const contractAddr = data.logs[0].args.contractAddress;

        await truffleAssert.reverts(
            this.chat.updateBodyInfo(contractAddr, principal, couponRate,
                couponDate, maturityDate, currency, {from: fakeOwner}),
            "Only the bond issuer can introduce changes"
        );
    })

    it("should revert if BondFinalised status is currently set", async () => {
        let data = await this.chat.createBond();
        truffleAssert.eventEmitted(data, 'newBondCreated');

        const contractAddr = data.logs[0].args.contractAddress;
        await this.chat.updateBondStatus(contractAddr, 5); // 5 is index of BondFinalised status

        await truffleAssert.reverts(
            this.chat.updateBodyInfo(contractAddr, principal, couponRate,
                couponDate, maturityDate, currency),
            "Bond already finalised"
        );
    });

    it("should not return a revert if all requirements are met", async () => {
        let data = await this.chat.createBond();
        truffleAssert.eventEmitted(data, 'newBondCreated');

        const contractAddr = data.logs[0].args.contractAddress;
        await truffleAssert.passes(
            this.chat.updateBodyInfo(contractAddr, principal, couponRate,
                couponDate, maturityDate, currency)
        );
    });
});

contract("Test update bond status", () => {
    it("should allow emit finalBondTerms event on final terms agreement", async () => {
        let chat = await chatContract.deployed();

        let data = await chat.createBond();
        truffleAssert.eventEmitted(data, 'newBondCreated');

        this.contractAddr = data.logs[0].args.contractAddress;

        // set the bond body info
        await truffleAssert.passes(
            chat.updateBodyInfo(this.contractAddr, principal, couponRate,
                couponDate, maturityDate, currency)
        );

        // Set the TermsAgreement bond status.
       let newdata =  await chat.updateBondStatus(this.contractAddr, 2);

        truffleAssert.eventEmitted(newdata, 'finalBondTerms', (ev) => {
            return ev.principal == principal && ev.couponRate == couponRate &&
                ev.couponDate == couponDate && ev.maturityDate == maturityDate && 
                ev.currency == currency
        });
    });

    it("should allow update of any bond status", async () => {
        let chat = await chatContract.deployed();

        await truffleAssert.passes(chat.updateBondStatus(this.contractAddr, 5));
    });

    it("should revert if BondFinalised status is currently set", async () => {
        let chat = await chatContract.deployed();

        // finalised bond should not accept any more status change.
        await truffleAssert.reverts(
            chat.updateBondStatus(this.contractAddr, 1),
            "Bond already finalised"
        );
    });
});

contract("Test add message to the bond chat", (accounts) => {
    const owner = accounts[0];
    beforeEach(
        async () => {
            this.chat = await chatContract.new({from: owner});

            let data = await chat.createBond();
            truffleAssert.eventEmitted(data, 'newBondCreated');
        }
    );

    const potentialBondHolder1 = accounts[5];

    it("should only allow potential bond holders to only send messages during negotiation stage", async (accounts) => {

    })

});

contract("Test update the bond holder address", () => {

});