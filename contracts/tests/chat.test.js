const chatContract = artifacts.require("../chat/chatContract");
const assert = require("truffle-assertions");
const web3 = require("web3");

contract("Test create a new Bond", () => {
    it("newBondCreated event should be received", async () => {
        const chat = await chatContract.new();
        let data = await chat.createBond(); 

        assert.eventEmitted(data, 'newBondCreated', (ev) => { return true });
    });
});

contract("Test update the body information", () => {
    beforeEach(
        async () => {
            this.chat = await chatContract.new();
        }
    );

    //it()
});

contract("Test update bond status", () => {

});

contract("Test add message to the bond chat", () => {

});

contract("Test update the bond holder address", () => {

});