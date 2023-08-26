const chatContract = artifacts.require("../chat/ChatContract");
const truffleAssert = require("truffle-assertions");

// Bond fields
const currency = 2;
const principal = 12323;
const couponRate = 2;
const couponDate =  3;
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

// chat message types
const initChat = 0;
const intro = 1;
const security = 2;
const appendix = 3;

// While running sapphire-dev 2023-07-10-gitbacd168 (oasis-core: 22.2.8, sapphire-paratime: 0.5.2, oasis-web3-gateway: 3.3.0-gitbacd168)
// instances, I have noticed that revert message returned by the contract don't
// show the expected error message, instead take the format:

// Transaction: 0x2b33087437fe7bbb8b4bceb7d45ec88698b380e7efab1d98219797132abe3944 exited with an error (status 0).
//      Please check that the transaction:
//      - satisfies all conditions set by Solidity `require` statements.
//      - does not trigger a Solidity `revert` statement.

// As a permanent fix is awaited, I choose to check for unspecified reverts so
// as to guarrantee that tests will pass.

// Should a better solution be found in future, error messages check will be reverted.

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

            truffleAssert.eventEmitted(data, 'NewBondCreated', (ev) => {
                return ev.sender == owner;
            });
        });
    });

    describe("Test update bond body information", () => {
        console.log()
        beforeEach( async () => {
                this.chat = await chatContract.new({from: owner});
                let data = await this.chat.createBond();

                truffleAssert.eventEmitted(data, 'NewBondCreated', (ev) => {
                    return ev.sender == owner;
                });
                this.contractAddr = data.logs[0].args.bondAddress;
            }
        );

        it("should revert if the sender detected is not the bond issuer", async () => {
            await truffleAssert.reverts(
                this.chat.updateBodyInfo(this.contractAddr, principal, couponRate,
                    couponDate, maturityDate, currency, {from: fakeOwner})
                //  "Edits only added by Bond Issuer"
            );
        })

        it("should not revert when setting status during negotations or holderselection", async () => {
            let statusChange = await this.chat.updateBondStatus(this.contractAddr, holderSelection);
            truffleAssert.eventEmitted(statusChange, 'StatusChange');

            let bondterms = await this.chat.updateBodyInfo(this.contractAddr, principal, couponRate,
                    couponDate, maturityDate, currency);

            truffleAssert.eventEmitted(bondterms, 'BondBodyTerms', (ev) => {
                return ev.bondAddress == this.contractAddr &&
                    ev.principal == principal && ev.couponRate == couponRate &&
                    ev.couponDate == couponDate && ev.maturityDate == maturityDate &&
                    ev.currency == currency;
            });
        });

        it("should revert if editting is attempted when the bond is dispute", async () => {
            let statusChange = await this.chat.updateBondStatus(this.contractAddr, holderSelection);
            truffleAssert.eventEmitted(statusChange, 'StatusChange');

            let holderUpdate = await this.chat.updateBondHolder(this.contractAddr, holder1);
            truffleAssert.eventEmitted(holderUpdate, 'HolderUpdate');

            let bondterms = await this.chat.updateBodyInfo(this.contractAddr, principal, couponRate,
                couponDate, maturityDate, currency);

            truffleAssert.eventEmitted(bondterms, 'BondBodyTerms', (ev) => {
                return ev.bondAddress == this.contractAddr &&
                    ev.principal == principal && ev.couponRate == couponRate &&
                    ev.couponDate == couponDate && ev.maturityDate == maturityDate &&
                    ev.currency == currency;
            });

            await this.chat.updateBondStatus(this.contractAddr, bondInDispute);

            await truffleAssert.reverts(
                this.chat.updateBodyInfo(this.contractAddr, principal, couponRate,
                    couponDate, maturityDate, currency)
                // "Some Bond dispute(s) are pending"
            );
        });

        it("should revert if terms update are initiated past TermsAgreement status", async () => {
            let statusChange = await this.chat.updateBondStatus(this.contractAddr, holderSelection);
            truffleAssert.eventEmitted(statusChange, 'StatusChange');

            let holderUpdate = await this.chat.updateBondHolder(this.contractAddr, holder1);
            truffleAssert.eventEmitted(holderUpdate, 'HolderUpdate');

            let bondterms = await this.chat.updateBodyInfo(this.contractAddr, principal, couponRate,
                    couponDate, maturityDate, currency);
            truffleAssert.eventEmitted(bondterms, 'BondBodyTerms', (ev) => {
                return ev.bondAddress == this.contractAddr &&
                    ev.principal == principal && ev.couponRate == couponRate &&
                    ev.couponDate == couponDate && ev.maturityDate == maturityDate &&
                    ev.currency == currency;
            });

            statusChange = await this.chat.updateBondStatus(this.contractAddr, bondReselling);
            truffleAssert.eventEmitted(statusChange, 'StatusChange');

            await truffleAssert.reverts(
                this.chat.updateBodyInfo(this.contractAddr, principal, couponRate,
                    couponDate, maturityDate, currency)
                // "Bond terms update is disabled"
            );
        });
    });

    describe("Test update bond status", () => {

        beforeEach(async () => {
            this.chat = await chatContract.new({from: owner});
            let data = await this.chat.createBond();

            truffleAssert.eventEmitted(data, 'NewBondCreated');
            this.contractAddr = data.logs[0].args.bondAddress;
        });

        it("should revert if status is to be set past holderselection before setting the holder's address", async () => {
            let bondterms = await this.chat.updateBodyInfo(this.contractAddr, principal, couponRate,
                    couponDate, maturityDate, currency);
            truffleAssert.eventEmitted(bondterms, 'BondBodyTerms');

            // dispute cannot be initiated without two people.
            await truffleAssert.reverts(
                this.chat.updateBondStatus(this.contractAddr, bondInDispute)
               // "Missing bond holder address"
            );
        });

        it("should revert if status update past holderselection with some missing body fields is triggered", async () => {
            let statusChange = await this.chat.updateBondStatus(this.contractAddr, holderSelection);
            truffleAssert.eventEmitted(statusChange, 'StatusChange', (ev) => {
                return ev.sender == owner && ev.bondAddress == this.contractAddr &&
                    ev.status == holderSelection;
            });

            let holderUpdate = await this.chat.updateBondHolder(this.contractAddr, holder);
            truffleAssert.eventEmitted(holderUpdate, 'HolderUpdate');

            let bondterms = await this.chat.updateBodyInfo(this.contractAddr, principal, 0,
                    couponDate, 0, currency);
            truffleAssert.eventEmitted(bondterms, 'BondBodyTerms', (ev) => {
                return ev.bondAddress == this.contractAddr &&
                    ev.principal == principal && ev.couponRate == 0 &&
                    ev.couponDate == couponDate && ev.maturityDate == 0 &&
                    ev.currency == currency;
            });

            await truffleAssert.reverts(
                this.chat.updateBondStatus(this.contractAddr, termsAgreement)
                // "Empty Bond body fields exists"
            );
        });

        it("should block status change if bond dispute is yet to be resolved fully", async () => {
            let statusChange = await this.chat.updateBondStatus(this.contractAddr, holderSelection);
            truffleAssert.eventEmitted(statusChange, 'StatusChange', (ev) => {
                return ev.sender == owner && ev.bondAddress == this.contractAddr &&
                    ev.status == holderSelection;
            });

            let holderUpdate = await this.chat.updateBondHolder(this.contractAddr, holder);
            truffleAssert.eventEmitted(holderUpdate, 'HolderUpdate');

            let bondterms = await this.chat.updateBodyInfo(this.contractAddr, principal, couponRate,
                couponDate, maturityDate, currency);
            truffleAssert.eventEmitted(bondterms, 'BondBodyTerms');

            // mark the current bond as under dispute.
            let txInfo = await this.chat.updateBondStatus(this.contractAddr, bondInDispute);
            truffleAssert.eventEmitted(txInfo, 'StatusChange', (ev) => {
                return ev.sender == owner && ev.bondAddress == this.contractAddr &&
                    ev.status == bondInDispute;
            });

            // Attempting to change status without resolving the dispute first should fail.
            await truffleAssert.reverts(
                this.chat.updateBondStatus(this.contractAddr, termsAgreement)
                // "Some Bond dispute(s) are pending"
            );
        });

        it('should not revert if a bond previously under dispute is resolved and status updated', async () => {
            let statusChange = await this.chat.updateBondStatus(this.contractAddr, holderSelection);
            truffleAssert.eventEmitted(statusChange, 'StatusChange', (ev) => {
                return ev.sender == owner && ev.bondAddress == this.contractAddr &&
                    ev.status == holderSelection;
            });

            let holderUpdate = await this.chat.updateBondHolder(this.contractAddr, holder);
            truffleAssert.eventEmitted(holderUpdate, 'HolderUpdate');

            let bondterms = await this.chat.updateBodyInfo(this.contractAddr, principal, couponRate,
                couponDate, maturityDate, currency);
            truffleAssert.eventEmitted(bondterms, 'BondBodyTerms');

            // mark the current bond as under dispute.
            let txInfo = await this.chat.updateBondStatus(this.contractAddr, bondInDispute);
            truffleAssert.eventEmitted(txInfo, 'StatusChange', (ev) => {
                return ev.sender == owner && ev.bondAddress == this.contractAddr &&
                    ev.status == bondInDispute;
            });

            // Issuer resolves the dispute on his end.
            let issuerData = await this.chat.signBondStatus(this.contractAddr, {from: owner});
            truffleAssert.eventEmitted(issuerData, 'StatusSigned', (ev) => {
                return ev.sender == owner && ev.bondAddress == this.contractAddr &&
                    ev.status == bondInDispute;
            });

            // Attempting status change should still fail.
            await truffleAssert.reverts(
                this.chat.updateBondStatus(this.contractAddr, termsAgreement)
                // "Some Bond dispute(s) are pending"
            );

            // holder resolves the dispute on his end.
            let holderData = await this.chat.signBondStatus(this.contractAddr, {from: holder});
            truffleAssert.eventEmitted(holderData, 'StatusSigned', (ev) => {
                return ev.sender == holder && ev.bondAddress == this.contractAddr &&
                    ev.status == bondInDispute;
            });

            // Status change should succeed now.
            let change = await this.chat.updateBondStatus(this.contractAddr, termsAgreement);
            truffleAssert.eventEmitted(change, 'StatusChange', (ev) => {
                return ev.sender == owner && ev.bondAddress == this.contractAddr &&
                    ev.status == termsAgreement;
            });
        });

        it("should only allow setting the ContractSigned if TermsAgreement was the previous status is signed", async ()=>{
            let statusChange = await this.chat.updateBondStatus(this.contractAddr, holderSelection);
            truffleAssert.eventEmitted(statusChange, 'StatusChange', (ev) => {
                return ev.sender == owner && ev.bondAddress == this.contractAddr &&
                    ev.status == holderSelection;
            });

            let holderUpdate = await this.chat.updateBondHolder(this.contractAddr, holder);
            truffleAssert.eventEmitted(holderUpdate, 'HolderUpdate');

            let bondterms = await this.chat.updateBodyInfo(this.contractAddr, principal, couponRate,
                couponDate, maturityDate, currency);
            truffleAssert.eventEmitted(bondterms, 'BondBodyTerms');

            // Set the TermsAgreement bond status.
            statusChange =  await this.chat.updateBondStatus(this.contractAddr, termsAgreement);
            truffleAssert.eventEmitted(statusChange, 'StatusChange', (ev) => {
                return ev.sender == owner && ev.bondAddress == this.contractAddr &&
                    ev.status == termsAgreement;
            });

            // Setting the contract signed should fail since required signatures don't exist.
            await truffleAssert.reverts(
                this.chat.updateBondStatus(this.contractAddr, contractSigned)
                // "Terms agreed on not fully signed"
            );

            let change = await this.chat.signBondStatus(this.contractAddr); // Signed by issuer
            truffleAssert.eventEmitted(change, 'StatusSigned', (ev) => {
                return ev.sender == owner && ev.bondAddress == this.contractAddr &&
                    ev.status == termsAgreement;
            });

            change = await this.chat.signBondStatus(this.contractAddr, {from: holder});
            truffleAssert.eventEmitted(change, 'StatusSigned', (ev) => {
                return ev.sender == holder && ev.bondAddress == this.contractAddr &&
                    ev.status == termsAgreement;
            });

            // setting the contract signed status should work now.
            statusChange = await this.chat.updateBondStatus(this.contractAddr, contractSigned);
            truffleAssert.eventEmitted(statusChange, 'StatusChange', (ev) => {
                return ev.sender == owner && ev.bondAddress == this.contractAddr &&
                    ev.status == contractSigned;
            });
        });

        it("should revert if BondFinalised status is currently set", async () => {
            let statusChange = await this.chat.updateBondStatus(this.contractAddr, holderSelection);
            truffleAssert.eventEmitted(statusChange, 'StatusChange', (ev) => {
                return ev.sender == owner && ev.bondAddress == this.contractAddr &&
                    ev.status == holderSelection;
            });

            let holderUpdate = await this.chat.updateBondHolder(this.contractAddr, holder);
            truffleAssert.eventEmitted(holderUpdate, 'HolderUpdate');

            let bondterms = await this.chat.updateBodyInfo(this.contractAddr, principal, couponRate,
                couponDate, maturityDate, currency);
            truffleAssert.eventEmitted(bondterms, 'BondBodyTerms');

            // mark the current bond as finalised.
            statusChange = await this.chat.updateBondStatus(this.contractAddr, bondFinalised);
            truffleAssert.eventEmitted(statusChange, 'StatusChange', (ev) => {
                return ev.sender == owner && ev.bondAddress == this.contractAddr &&
                    ev.status == bondFinalised;
            });

            // finalised bond should not accept any more status changes.
            await truffleAssert.reverts(
                this.chat.updateBondStatus(this.contractAddr, holderSelection)
               // "Edits disabled on finalized Bond"
            );
        });
    });

    describe("Test add message to the bond chat", () => {

        beforeEach(async () => {
            this.chat = await chatContract.new({from: owner});
            let data = await this.chat.createBond();

            truffleAssert.eventEmitted(data, 'NewBondCreated');
            this.contractAddr = data.logs[0].args.bondAddress;

            let statusChange = await this.chat.updateBondStatus(this.contractAddr, holderSelection);
            truffleAssert.eventEmitted(statusChange, 'StatusChange');

            let holderUpdate = await this.chat.updateBondHolder(this.contractAddr, holder);
            truffleAssert.eventEmitted(holderUpdate, 'HolderUpdate');

            let bondterms = await this.chat.updateBodyInfo(this.contractAddr, principal, couponRate,
                couponDate, maturityDate, currency);
            truffleAssert.eventEmitted(bondterms, 'BondBodyTerms');

            // After update downgrade the status back to the default.
            statusChange = await this.chat.updateBondStatus(this.contractAddr, negotiating);
            truffleAssert.eventEmitted(statusChange, 'StatusChange');
        });

        it("should block potential bond holders from sending messages past negotiation stage", async () => {
            let statusChange = await this.chat.updateBondStatus(this.contractAddr, holderSelection);
            truffleAssert.eventEmitted(statusChange, 'StatusChange');

            await truffleAssert.reverts(
                this.chat.addMessage(this.contractAddr, initChat, chatMsg, {from: holder2})
                // "Only negotiation chat is general"
            );
        });

        it("should allow potential bond holders to only send messages during negotiation stage", async () => {
            let txInfo = await this.chat.addMessage(this.contractAddr, initChat, chatMsg, {from: holder1})

            truffleAssert.eventEmitted(txInfo, 'NewChatMessage', (ev) => {
                return ev.bondAddress == this.contractAddr &&
                    ev.sender == holder1 && ev.message == chatMsg;
            });
        })

        it("should allow both issuer and holder to chat messages before the bond is finalised", async () => {
            let statusChange = await this.chat.updateBondStatus(this.contractAddr, termsAgreement);
            truffleAssert.eventEmitted(statusChange, 'StatusChange');

            // issuer can send a message during TermsAgreement status.
            let issuerInfo = await this.chat.addMessage(this.contractAddr, initChat, chatMsg, {from: owner})

            truffleAssert.eventEmitted(issuerInfo, 'NewChatMessage', (ev) => {
                return ev.bondAddress == this.contractAddr && ev.sender == owner && ev.message == chatMsg;
            });

            let holderInfo = await this.chat.addMessage(this.contractAddr, initChat, chatMsg, {from: holder})

            truffleAssert.eventEmitted(holderInfo, 'NewChatMessage', (ev) => {
                return ev.bondAddress == this.contractAddr && ev.sender == holder && ev.message == chatMsg;
            });
        });

        it("should only allow the issuer to send the bond secure detail messages", async () => {
            let introInfo = await this.chat.addMessage(this.contractAddr, intro, chatMsg, {from: owner})

            truffleAssert.eventEmitted(introInfo, 'BondMotivation', (ev) => {
                return ev.bondAddress == this.contractAddr && ev.message == chatMsg;
            });

            await truffleAssert.passes(this.chat.addMessage(this.contractAddr, security, chatMsg, {from: owner}));
            await truffleAssert.passes(this.chat.addMessage(this.contractAddr, appendix, chatMsg, {from: owner}));

            await truffleAssert.reverts(
                this.chat.addMessage(this.contractAddr, intro, chatMsg, {from: holder})
                //  "Edits only added by Bond Issuer"
            );
            await truffleAssert.reverts(
                this.chat.addMessage(this.contractAddr, security, chatMsg, {from: holder})
                //  "Edits only added by Bond Issuer"
            );
            await truffleAssert.reverts(
                this.chat.addMessage(this.contractAddr, appendix, chatMsg, {from: holder}),
                //  "Edits only added by Bond Issuer"
            );
        });

        it("should revert if issuer chats/messages are sent once the bond has been finalised", async () => {
            let statusChange = await this.chat.updateBondStatus(this.contractAddr, bondFinalised);
            truffleAssert.eventEmitted(statusChange, 'StatusChange');

            await truffleAssert.reverts(
                this.chat.addMessage(this.contractAddr, security, chatMsg, {from: owner})
                // "Edits disabled on finalized Bond"
            );
        });

        it("should revert if issuer contact messages are sent once the bond is in dispute", async () => {
            let statusChange = await this.chat.updateBondStatus(this.contractAddr, bondInDispute);
            truffleAssert.eventEmitted(statusChange, 'StatusChange');

            await truffleAssert.reverts(
                this.chat.addMessage(this.contractAddr, security, chatMsg, {from: owner})
                // "Some Bond dispute(s) are pending"
            );
        });
    });

    describe("Test the update of bond holder address", () => {

        beforeEach(async () => {
            this.chat = await chatContract.new();
            let data = await this.chat.createBond();

            truffleAssert.eventEmitted(data, 'NewBondCreated');
            this.contractAddr = data.logs[0].args.bondAddress;

            let bondterms = await this.chat.updateBodyInfo(this.contractAddr, principal, couponRate,
                couponDate, maturityDate, currency);
            truffleAssert.eventEmitted(bondterms, 'BondBodyTerms');
        });

        it("should revert if the issuer attempts to set themselves as the holder", async () => {
            let statusChange = await this.chat.updateBondStatus(this.contractAddr, holderSelection);
            truffleAssert.eventEmitted(statusChange, 'StatusChange');

            await truffleAssert.reverts(
                this.chat.updateBondHolder(this.contractAddr, owner)
                // "Issuer & Holder must be separate"
            );
        });

        it("should revert if a holder is not set during holderSelection status", async () => {
            let statusChange = await this.chat.updateBondStatus(this.contractAddr, negotiating);
            truffleAssert.eventEmitted(statusChange, 'StatusChange');

            await truffleAssert.reverts(
                this.chat.updateBondHolder(this.contractAddr, holder)
                // "Holder is set on HolderSelection"
            );
        });

        it("should allow setting of an appropriate holder at the correct status", async () => {
            let statusChange = await this.chat.updateBondStatus(this.contractAddr, holderSelection);
            truffleAssert.eventEmitted(statusChange, 'StatusChange');

            let holderUpdate = await this.chat.updateBondHolder(this.contractAddr, holder);
            truffleAssert.eventEmitted(holderUpdate, 'HolderUpdate', (ev) => {
                return ev.bondAddress == this.contractAddr && ev.holder == holder;
            });
        });
    });

    describe("Test the signing of a status", () => {

        beforeEach(async () => {
            this.chat = await chatContract.new();
            let data = await this.chat.createBond();

            truffleAssert.eventEmitted(data, 'NewBondCreated');
            this.contractAddr = data.logs[0].args.bondAddress;

            let bondterms = await this.chat.updateBodyInfo(this.contractAddr, principal, couponRate,
                couponDate, maturityDate, currency);
            truffleAssert.eventEmitted(bondterms, 'BondBodyTerms');

            let statusChange = await this.chat.updateBondStatus(this.contractAddr, holderSelection);
            truffleAssert.eventEmitted(statusChange, 'StatusChange');

            let holderUpdate = await this.chat.updateBondHolder(this.contractAddr, holder);
            truffleAssert.eventEmitted(holderUpdate, 'HolderUpdate');
        });

        it("should emit the owner dispute resolved event once signed", async () => {
            let statusChange = await this.chat.updateBondStatus(this.contractAddr, bondInDispute);
            truffleAssert.eventEmitted(statusChange, 'StatusChange');

            let issuerData = await this.chat.signBondStatus(this.contractAddr);

            truffleAssert.eventEmitted(issuerData, "StatusSigned", (ev) => {
                return ev.sender = owner && ev.bondAddress == this.contractAddr &&
                    ev.status == bondInDispute;
            });
        });

        it("should emit event after signing the current bond status", async () => {
            let txInfo = await this.chat.signBondStatus(this.contractAddr, {from: holder});

            truffleAssert.eventEmitted(txInfo, "StatusSigned",(ev) => {
                return ev.sender = holder && ev.bondAddress == this.contractAddr &&
                    ev.status == holderSelection;
            });
        });
    });

    describe("Test fetching the bond secure details", () => {

        beforeEach(async () => {
            this.chat = await chatContract.new();
            let data = await this.chat.createBond();

            truffleAssert.eventEmitted(data, 'NewBondCreated');
            this.contractAddr = data.logs[0].args.bondAddress;

            let bondterms = await this.chat.updateBodyInfo(this.contractAddr, principal, couponRate,
                couponDate, maturityDate, currency);
            truffleAssert.eventEmitted(bondterms, 'BondBodyTerms');

            let statusChange = await this.chat.updateBondStatus(this.contractAddr, holderSelection);
            truffleAssert.eventEmitted(statusChange, 'StatusChange');

            let holderUpdate = await this.chat.updateBondHolder(this.contractAddr, holder);
            truffleAssert.eventEmitted(holderUpdate, 'HolderUpdate');

            let introInfo = await this.chat.addMessage(this.contractAddr, intro, chatMsg, {from: owner})
            truffleAssert.eventEmitted(introInfo, 'BondMotivation', (ev) => {
                return ev.bondAddress == this.contractAddr && ev.message == chatMsg;
            });

            await truffleAssert.passes(this.chat.addMessage(this.contractAddr, security, chatMsg, {from: owner}));
            await truffleAssert.passes(this.chat.addMessage(this.contractAddr, appendix, chatMsg, {from: owner}));
        });

        it("should revert if non-bond party attempts to fetch the details", async () => {
            await truffleAssert.reverts(
                this.chat.getBondSecureDetails(this.contractAddr, {from: holder1})
                // "Only accessed by parties to the bond"
            );
        });

        it("should return the secure details if fetched by the owner or holder", async () => {
            let a = await this.chat.getBondSecureDetails(this.contractAddr, {from: owner})
            assert.notEqual(a.security, "", "expected owner security response not be empty")
            assert.notEqual(a.appendix, "", "expected owner appendix response not be empty")

            let b = await this.chat.getBondSecureDetails(this.contractAddr, {from: holder})
            assert.notEqual(b.security, "", "expected holder security response not be empty")
            assert.notEqual(b.appendix, "", "expected holder appendix response not be empty")

            assert.equal(a.security, b.security, "expected owner security to match holder security")
            assert.equal(a.appendix, b.appendix, "expected owner appendix to match holder appendix")
        });
    });
});