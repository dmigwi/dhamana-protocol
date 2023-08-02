const chatContract = artifacts.require("../chat/chatContract");
const truffleAssert = require("truffle-assertions");

// Bond fields
const currency = 2;
const principal = 12323;
const couponRate = 2;
const couponDate = 1690670000;
const maturityDate = 1690869606;
const chatMsg = "This text is a placeholder of an encrypted message";

// Bond Status
const negotiating = 0;
const holderSelection = 1;
const bondInDispute = 2;
const termsAgreement = 3;
const contractSigned = 4
const bondReselling = 5;
const bondFinalised = 6


contract("ChatContract",  (accounts) => {

    // accounts
    const owner = accounts[0];
    const holder = accounts[1];
    const holder1 = accounts[2];
    const holder2 = accounts[3];
    const fakeOwner = accounts[4];

    describe("Test create a new Bond", () => {
        it("should receive a newBondCreated event", async () => {
            const chat = await chatContract.new({from: owner});
            let data = await chat.createBond();

            truffleAssert.eventEmitted(data, 'newBondCreated', (ev) => { return true });
        });
    });

    describe("Test update bond body information", () => {
        beforeEach(async () => {
                this.chat = await chatContract.new({from: owner});
                let data = await this.chat.createBond();

                truffleAssert.eventEmitted(data, 'newBondCreated');
                this.contractAddr = data.logs[0].args.contractAddress;
            }
        );

        it("should revert if the sender detected is not the bond issuer", async () => {
            await truffleAssert.reverts(
                this.chat.updateBodyInfo(this.contractAddr, principal, couponRate,
                    couponDate, maturityDate, currency, {from: fakeOwner}),
                "Only the bond issuer can introduce changes"
            );
        })

        it("should not return a revert when setting status during negotations or holderselection", async () => {
            await this.chat.updateBondStatus(this.contractAddr, holderSelection);
            await truffleAssert.passes(
                this.chat.updateBodyInfo(this.contractAddr, principal, couponRate,
                    couponDate, maturityDate, currency)
            );
        });

        it("should revert if editting is attempted when the bond is dispute", async () => {
            await this.chat.updateBondStatus(this.contractAddr, holderSelection);
            await this.chat.updateBondHolder(this.contractAddr, holder1);
            await truffleAssert.passes(
                this.chat.updateBodyInfo(this.contractAddr, principal, couponRate,
                    couponDate, maturityDate, currency)
            );

            await this.chat.updateBondStatus(this.contractAddr, bondInDispute);

            await truffleAssert.reverts(
                this.chat.updateBodyInfo(this.contractAddr, principal, couponRate,
                    couponDate, maturityDate, currency),
                "Bond disputes not yet resolved by all the parties"
            );
        });

        it("should revert if terms update are initiated past TermsAgreement status", async () => {
            await this.chat.updateBondStatus(this.contractAddr, holderSelection);
            await this.chat.updateBondHolder(this.contractAddr, holder1);
            await truffleAssert.passes(
                this.chat.updateBodyInfo(this.contractAddr, principal, couponRate,
                    couponDate, maturityDate, currency)
            );

            await this.chat.updateBondStatus(this.contractAddr, bondReselling);

            await truffleAssert.reverts(
                this.chat.updateBodyInfo(this.contractAddr, principal, couponRate,
                    couponDate, maturityDate, currency),
                "Bond terms update is disabled"
            );
        });
    });

    describe("Test update bond status", () => {

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
                this.chat.updateBondStatus(this.contractAddr, bondInDispute),
                "Missing bond holder address"
            );
        });

        it("should revert if status update past holderselection with some missing body fields is triggered", async () => {
            await this.chat.updateBondStatus(this.contractAddr, holderSelection);
            await this.chat.updateBondHolder(this.contractAddr, holder);
            await truffleAssert.passes(
                this.chat.updateBodyInfo(this.contractAddr, principal, 0,
                    couponDate, 0, currency)
            );

            await truffleAssert.reverts(
                this.chat.updateBondStatus(this.contractAddr, termsAgreement),
                "Bond body fields may contain empty values"
            );
        });

        it("should revert if bond dispute is yet to be resolved", async () => {
            await this.chat.updateBondStatus(this.contractAddr, holderSelection);
            await truffleAssert.passes(this.chat.updateBondHolder(this.contractAddr, holder));
            await this.chat.updateBodyInfo(this.contractAddr, principal, couponRate,
                couponDate, maturityDate, currency)

            // mark the current bond as under dispute.
            let txInfo = await this.chat.updateBondStatus(this.contractAddr, bondInDispute);

            truffleAssert.eventEmitted(txInfo, 'bondUnderDispute', (ev) => {
                return ev.sender == owner &&  ev.bondAddress == this.contractAddr
            });

            // Attempting to change status without resolving the dispute first should fail.
            await truffleAssert.reverts(
                this.chat.updateBondStatus(this.contractAddr, termsAgreement),
                "Bond disputes not yet resolved by all the parties"
            );
        });

        it('should not revert if a bond previously under dispute is resolved and status updated', async () => {
            await this.chat.updateBondStatus(this.contractAddr, holderSelection);
            await truffleAssert.passes(this.chat.updateBondHolder(this.contractAddr, holder));
            await this.chat.updateBodyInfo(this.contractAddr, principal, couponRate,
                couponDate, maturityDate, currency)

            // mark the current bond as under dispute.
            let txInfo = await this.chat.updateBondStatus(this.contractAddr, bondInDispute);

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
                this.chat.updateBondStatus(this.contractAddr, termsAgreement),
                "Bond disputes not yet resolved by all the parties"
            );

            // holder resolves the dispute on his end.
            let holderData = await this.chat.signBondStatus(this.contractAddr, {from: holder});
            truffleAssert.eventEmitted(holderData, 'bondDisputeResolved', (ev) => {
                return ev.sender == holder && ev.bondAddress == this.contractAddr;
            });

            // Status change should succeed now.
            await truffleAssert.passes(this.chat.updateBondStatus(this.contractAddr, termsAgreement));
        });

        it("should allow emit finalBondTerms event on TermsAgreement status", async () => {
            await this.chat.updateBondStatus(this.contractAddr, holderSelection);
            await this.chat.updateBondHolder(this.contractAddr, holder);
            await truffleAssert.passes(
                this.chat.updateBodyInfo(this.contractAddr, principal, couponRate,
                    couponDate, maturityDate, currency)
            );

            // Set the TermsAgreement bond status.
            let newdata =  await this.chat.updateBondStatus(this.contractAddr, termsAgreement);

            truffleAssert.eventEmitted(newdata, 'finalBondTerms', (ev) => {
                return ev.principal == principal && ev.couponRate == couponRate &&
                    ev.couponDate == couponDate && ev.maturityDate == maturityDate &&
                    ev.currency == currency;
            });
        });

        it("should only allow setting the ContractSigned and TermsAgreement as the previous status is signed", async ()=>{
            await this.chat.updateBondStatus(this.contractAddr, holderSelection);
            await this.chat.updateBondHolder(this.contractAddr, holder);
            await truffleAssert.passes(
                this.chat.updateBodyInfo(this.contractAddr, principal, couponRate,
                    couponDate, maturityDate, currency)
            );

            // Set the TermsAgreement bond status.
            let newdata =  await this.chat.updateBondStatus(this.contractAddr, termsAgreement);
            truffleAssert.eventEmitted(newdata, 'finalBondTerms');

            // Setting the contract signed should fail since required signatures don't exist.
            await truffleAssert.reverts(
                this.chat.updateBondStatus(this.contractAddr, contractSigned),
                "TermsAgreement must be signed before setting ContractSigned"
            );

            await truffleAssert.passes(this.chat.signBondStatus(this.contractAddr)); // Signed by issuer
            await truffleAssert.passes(this.chat.signBondStatus(this.contractAddr, {from: holder}));

            // setting the contract signed status should work now.
            await truffleAssert.passes(this.chat.updateBondStatus(this.contractAddr, contractSigned));
        });

        it("should revert if BondFinalised status is currently set", async () => {
            await this.chat.updateBondStatus(this.contractAddr, holderSelection);
            await this.chat.updateBondHolder(this.contractAddr, holder);
            await this.chat.updateBodyInfo(this.contractAddr, principal, couponRate,
                couponDate, maturityDate, currency);

            // mark the current bond as finalised.
            await truffleAssert.passes(this.chat.updateBondStatus(this.contractAddr, bondFinalised));

            // finalised bond should not accept any more status changes.
            await truffleAssert.reverts(
                this.chat.updateBondStatus(this.contractAddr, holderSelection),
                "Bond already finalised"
            );
        });
    });

    describe("Test add message to the bond chat", () => {

        beforeEach(async () => {
            this.chat = await chatContract.new({from: owner});
            let data = await this.chat.createBond();

            truffleAssert.eventEmitted(data, 'newBondCreated');
            this.contractAddr = data.logs[0].args.contractAddress;

            await truffleAssert.passes(this.chat.updateBondStatus(this.contractAddr, holderSelection));
            await truffleAssert.passes(this.chat.updateBondHolder(this.contractAddr, holder));
            await this.chat.updateBodyInfo(this.contractAddr, principal, couponRate,
                couponDate, maturityDate, currency);

            await this.chat.updateBondStatus(this.contractAddr, negotiating);
        });

        const initChat = 0;
        const intro = 1;
        const security = 2;
        const appendix = 3;

        it("should block potential bond holders from sending messages past negotiation stage", async () => {
            await this.chat.updateBondStatus(this.contractAddr, holderSelection)
            await truffleAssert.reverts(
                this.chat.addMessage(this.contractAddr, initChat, chatMsg, {from: holder2}),
                "Potential holders cannot comment past negotiating stage"
            );
        });

        it("should allow potential bond holders to only send messages during negotiation stage", async () => {
            let txInfo = await this.chat.addMessage(this.contractAddr, initChat, chatMsg, {from: holder1})

            truffleAssert.eventEmitted(txInfo, 'newChatMessage', (ev) => {
                return ev.sender == holder1;
            });
        })

        it("should allow both issuer and holder to chat messages till the bond is finalised", async () => {
            await truffleAssert.passes(this.chat.updateBondStatus(this.contractAddr, termsAgreement));

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
            await this.chat.updateBondStatus(this.contractAddr, bondInDispute);

            await truffleAssert.reverts(
                this.chat.addMessage(this.contractAddr, security, chatMsg, {from: owner}),
                "Bond disputes not yet resolved by all the parties"
            );
        });
    });

    describe("Test the update of bond holder address", () => {

        beforeEach(async () => {
            this.chat = await chatContract.new({from: owner});
            let data = await this.chat.createBond();

            truffleAssert.eventEmitted(data, 'newBondCreated');
            this.contractAddr = data.logs[0].args.contractAddress;

            await this.chat.updateBodyInfo(this.contractAddr, principal, couponRate,
                couponDate, maturityDate, currency);
        });

        it("should revert if the issuer attempts to set themselves as the holder", async () => {
            await this.chat.updateBondStatus(this.contractAddr, holderSelection);
            await truffleAssert.reverts(
                this.chat.updateBondHolder(this.contractAddr, owner),
                "Bond issuer cannot be the holder too"
            );
        });

        it("should revert if a holder is not set during holderSelection status", async () => {
            await this.chat.updateBondStatus(this.contractAddr, negotiating);
            await truffleAssert.reverts(
                this.chat.updateBondHolder(this.contractAddr, holder),
                "Only set during the HolderSelection status"
            );
        });

        it("should allow setting of an appropriate holder at the correct status", async () => {
            await this.chat.updateBondStatus(this.contractAddr, holderSelection);
            await truffleAssert.passes( this.chat.updateBondHolder(this.contractAddr, holder));
        });
    });

    describe("Test the signing of a status", () => {

        beforeEach(async () => {
            this.chat = await chatContract.new({from: owner});
            let data = await this.chat.createBond();

            truffleAssert.eventEmitted(data, 'newBondCreated');
            this.contractAddr = data.logs[0].args.contractAddress;

            await this.chat.updateBodyInfo(this.contractAddr, principal, couponRate,
                couponDate, maturityDate, currency);

            await this.chat.updateBondStatus(this.contractAddr, holderSelection);
            await truffleAssert.passes(this.chat.updateBondHolder(this.contractAddr, holder));
        });

        it("should emit the dispute resolved event once signed", async () => {
            await this.chat.updateBondStatus(this.contractAddr, bondInDispute);

            let issuerData = await this.chat.signBondStatus(this.contractAddr);

            truffleAssert.eventEmitted(issuerData, "bondDisputeResolved", (ev) => {
                return ev.sender = owner && ev.bondAddress == this.contractAddr;
            });
        });

        it("should not emit dispute resolved event when not signing BondInDispute Status", async () => {
            let txInfo = await this.chat.signBondStatus(this.contractAddr, {from: holder});

            truffleAssert.eventNotEmitted(txInfo, "bondDisputeResolved");
        });
    });
});