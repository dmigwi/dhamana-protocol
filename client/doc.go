// Copyright (c) 2023 Migwi Ndung'u
// See LICENSE for details.

/*
    Client implements the dhamana-protocol on the sapphire paratime.

    dhamana-protocol is an implementation that allows people unknown to each other
    but vetted by a Trust Organisation to engage in a conversation on how
    one of them (bond issuer) can have their idea financed in a legally binding
    contract (bond).

    The integrity of the system is maintained by trust earned through referals.

    -------------------------------------DHAMANA-PROTOCOL OVERVIEW--------------------------------------


                        David     Ann         Bob           Alice      Steve        Anonymous      Users who want to sign up
                        |         |           |              |            |            |            to the platform must be
                        |         |           |              |            |            |            vetted by the 3rd party.
                        v         v           v              v            v            v
                    +---------------------------------------------------------------------+
                    |        Trust Organisations - Does Users Vetting Only                |
                    |                                                                     |
                    +---------------------------------------------------------------------+
                        |         |           |              |            |
            *********************************************************************************  Vetted users are granted
                        |         |           |              |            |                    access to the dhamana
                        v         v           v              v            v          protocol via a POA (Point of Access).
                +---------+  +-------+  +---------+  +---------+  +---------+        POA can be an app or a website.
                | David's |  | Ann's |  |  Bob's  |  | Alice's |  | Steve's |
                |   POA   |  | POA   |  |   POA   |  |   POA   |  |   POA   |
                +---------+  +-------+  +---------+  +---------+  +---------+
                                            |
                                            | Communication between POA and the client is encrypted to mitigate
                                            v                                            man-in-the-middle attack.
                            +-------------------------------------------+
(Bond Lifecycle)            |  Bob creates a bond outlining why someone |
                            |  should fund it. He also offers matching  |
                            |  securities as a guarantee of payback.    |
                            |                                           |
````````````````````````````|```````````````````````````````````````````|``````````````````     All Users can express interest.
1. Negotiation Stage        |  David, Alice and Ann express interest    |    <------- David     Only the issuer can select a
                Ann ------> |     in subscribing to Bob's bond.         |    <------- Alice     potential bond holder.
                            |                                           |
````````````````````````````|```````````````````````````````````````````|``````````````````
2. HolderSelection Stage    |    Bob considers all the interests        |                        All other users apart from
                            |     expressed but finds Alice's terms     |    <------- Alice        Alice and Bob are locked
                            |    more favourable. He selects Alice as   |                        out of the bond activities.
                            |    the potential holder.                  |
                            |                                           |
````````````````````````````|```````````````````````````````````````````|``````````````````
3. TermsAgreement Stage     |    Bob and Alice agree on the finer       |
                            |    details of the bond. Past this stage   |    <------- Alice
                            |    further terms update is disabled.      |
                            |                                           |
````````````````````````````|```````````````````````````````````````````|``````````````````
4. BondInDispute Stage      |    Should either Alice or Bob             |
                            |     (any party to the bond) become        |    <------- Alice
                            |    disatisfied with the other party's     |
                            |    actions, they can move the bond into   |
                            |   dispute resolution phase. Further       |
                            |  progress on the bond stages is blocked   |
                            | till all parties resolve the conflicts.   |
                            |                                           |
````````````````````````````|```````````````````````````````````````````|``````````````````
5. ContractSigned Stage     |  After Bob and Alice agree on the final   |
                            |  details of the bond terms, a bond        |
                            |  document is encrypted with holder's      |
                            |   (Alice's) keys to show her ownership    |
                            |                                           |
````````````````````````````|```````````````````````````````````````````|``````````````````
6. BondReselling Stage      |    Alice can choose to transfer all her   |
                Ann ------> | rights on the bond to someone else (Ann). |  <------- Alice
                            |                                           |
````````````````````````````|```````````````````````````````````````````|``````````````````
7. BondFinalised Stage      |  Bob has fulfilled his obligation by      |
                            |    paying all the money owed in full.     |
                            |  Bond isconsidered complete and further   |
                            |    activity on it is disabled.            |
                            +-------------------------------------------+

        ****************************************************************************************

                                        SAPPHIRE PARATIME LAYER

    -------------------------------------DHAMANA-PROTOCOL END--------------------------------------

The Truth organisation concept introduces the delegated managing of critical user identifying information
in a distributed manner. Users interact with each other based on the trust owned by the vetting trust organisations.
A user's trust is a strong and solid as the trust organisation that vetted them. The trustworthiness of a truth
organisation is directly proportional to the users in the system that use it as expected. There is no single source
of truth (user identity) but multiple sources of the same truth.
If the integrity of one truth organisation is compromised, users can still rely on the other safe truth organisations.
The effective functioning of dhamana-protocol rely on the integrity of the truth organisations. The system is meant
to work with users anonymously at the behest of the truth organisations.
If a malicious person penetrates into the system, it affects the truth organisation and users can be warned
to avoid the said truth organisation.

Dhamana-protocol is a trust based system where all users join the system expecting no malicious person/bot has been intentionally let in.
*/

package main
