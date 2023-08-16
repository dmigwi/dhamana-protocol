// Copyright (c) 2023 Migwi Ndung'u
// See LICENSE for details.

/*
 	Client holds an implementation of a node validator node that interacts with
	sapphire paratime when running the dhamana-protocol.
	
	dhamana-protocol is an implementation that allows people unknown to each other
	but approved by a reputable 3rd party to engage in a conversation on how
	one of them (bond issuer) can have his idea financed in a legally binding
	contract (bond).

	-------------------------------------DHAMANA-PROTOCOL OVERVIEW--------------------------------------



						David	 Ann	  Bob		   Alice	Steve		Anonymous      Users who want to sign up
						|		 |		   |			  |			|			|			to the platform must be
						|		 |		   |			  |			|			|			vetted by the 3rd party.
						v		 v		   v			  v			v			v
				+-------------------------------------------------------------------+
				|		Reputable 3rd Party - Does Users Vetting Only				|
				|																	|
				+-------------------------------------------------------------------+
						|		 |		   |			  |			|			
			*********************************************************************************  Vetted users are granted 
						|		 |		   |			  |			|							access to the dhamana
						v		 v		   v			  v			v				protocol via a POA (Point of Access).			
				+---------+  +-------+  +---------+	 +---------+  +---------+		POA can be a mobile app or a website.
				| David's |	 | Ann's |  |  Bob's  |  | Alice's |  | Steve's |
		 		|   POA   |	 | POA   |  |   POA   |	 |   POA   |  |   POA   |
				+---------+	 +-------+  +---------+	 +---------+  +---------+
											|
											| Communication between POA and the client is encrypted to mitigate man-in-the-middle attack.
											v
							+-------------------------------------------+   
							|  Bob creates a bond outlining why someone	|
							|  should fund it. He also offers matching	|
							|  securities as a guarantee of payback.	|
							|											|
````````````````````````````|```````````````````````````````````````````|``````````````````     All Users can express interest.
Negotiation Stage			|  David, Alice and Ann express an interest | 	<------- David		Only the issuer can select a 
				Ann ------> |	 of subscribing at Bob's bond.			|  	<------- Alice		potential bond holder.
							|											|
````````````````````````````|```````````````````````````````````````````|``````````````````
HolderSelection Stage		|	Bob considers all the interests			|						All other users apart from
							|	 expressed but finds Alice's terms		|	<------- Alice		Alice and Bob are locked
							|	more favourable. He selects Alice as	|						out of the bond activities.
							|	the potential holder.					|
							|											|
````````````````````````````|```````````````````````````````````````````|``````````````````
TermsAgreement Stage		|	Bob and Alice agree on the finer		|   
							|	details of the bond. Past this stage	|	<------- Alice
							|	further terms update is disabled.		|
							|											|
````````````````````````````|```````````````````````````````````````````|``````````````````
BondInDispute Stage			|	Should either Alice or Bob				|
							|	 (either parties to the bond) become 	|	<------- Alice
							|	disatisfied with the party insincere	|
							|	actions, they can move the bond into	|
							|	 dispute resolution phase. Further		|
							|	progress on the bond stages is blocked	|
							|	till all parties resolve the conflicts.	|
							|											|
````````````````````````````|```````````````````````````````````````````|``````````````````
ContractSigned Stage		|	After Bob and Alice agree on the final	|    
							|	details of the bond terms, a bond		|
							|	document is encrypted with holder's		|
							|	(Alice's) keys to show her ownership	|
							|											|
````````````````````````````|```````````````````````````````````````````|``````````````````
BondReselling Stage			|	Alice can choose to transfer all her 	|   
				Ann ------> | rights on the bond to someone else (Ann).	|  <------- Alice
							|											|	
````````````````````````````|```````````````````````````````````````````|``````````````````
BondFinalised Stage			|	Bob has fulfilled his obligation by		|   
							|	paying all the money owed in full.		|
							|	Bond isconsidered complete and further	|
							|		activity on it is disabled.			|
							+-------------------------------------------+
							


 */

package main
