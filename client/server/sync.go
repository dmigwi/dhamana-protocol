// Copyright (c) 2023 Migwi Ndung'u
// See LICENSE for details.

package server

import (
	"github.com/dmigwi/dhamana-protocol/client/contracts"
	"github.com/dmigwi/dhamana-protocol/client/utils"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
)

func (s *ServerConfig) SyncHistoricalData() error {
	// fetch the block at which the contract was deployed
	deployedBlock := getDeployedBlock(s.network)

	// fetch the last synced block from the database.
	lastSyncedBlock, _ := s.db.QueryLocalData(utils.GetLastSyncedBlock, new(lastSyncedBlockResp), "")
	var syncedBlock uint32
	if len(lastSyncedBlock) > 0 {
		syncedBlock = lastSyncedBlock[0].(uint32)
	}

	// compare the two blocks and pick the latests one.
	if deployedBlock > syncedBlock {
		syncedBlock = deployedBlock
	}

	start := uint64(syncedBlock)
	watchOpts := &bind.WatchOpts{
		Context: s.ctx,
		Start:   &start,
	}

	bondBodyTermsChan := make(chan *contracts.ChatBondBodyTerms)
	bondMotivationChan := make(chan *contracts.ChatBondMotivation)
	bondsignedStatusChan := make(chan *contracts.ChatBondStatusSigned)
	holderUpdateChan := make(chan *contracts.ChatHolderUpdate)
	newbondCreatedChan := make(chan *contracts.ChatNewBondCreated)
	newChatMessageChan := make(chan *contracts.ChatNewChatMessage)
	statusChangeChan := make(chan *contracts.ChatStatusChange)
	statusSignedChan := make(chan *contracts.ChatStatusSigned)

	s.bondChat.ChatFilterer.WatchBondBodyTerms(watchOpts, bondBodyTermsChan)
	s.bondChat.ChatFilterer.WatchBondMotivation(watchOpts, bondMotivationChan)
	s.bondChat.ChatFilterer.WatchBondStatusSigned(watchOpts, bondsignedStatusChan)
	s.bondChat.ChatFilterer.WatchHolderUpdate(watchOpts, holderUpdateChan)
	s.bondChat.ChatFilterer.WatchNewBondCreated(watchOpts, newbondCreatedChan)
	s.bondChat.ChatFilterer.WatchNewChatMessage(watchOpts, newChatMessageChan)
	s.bondChat.ChatFilterer.WatchStatusChange(watchOpts, statusChangeChan)
	s.bondChat.ChatFilterer.WatchStatusSigned(watchOpts, statusSignedChan)

	return nil
}
