// Copyright (c) 2023 Migwi Ndung'u
// See LICENSE for details.

package server

import (
	"fmt"
	"time"

	"github.com/dmigwi/dhamana-protocol/client/contracts"
	"github.com/dmigwi/dhamana-protocol/client/utils"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/event"
)

// eventData contains data packed from each event recieved.
type eventData struct {
	method  utils.Method
	blockNo uint64
	params  []interface{}
}

// SyncData creates a listener that waits for events data to be received then
// updates it to the local data persistence storage.
func (s *ServerConfig) SyncData() error {
	// fetch the block at which the contract was deployed
	deployedBlock := getDeployedBlock(s.network)

	// fetch the last synced block from the database.
	lastSyncedBlock, _ := s.db.QueryLocalData(utils.GetLastSyncedBlock,
		new(lastSyncedBlockResp), "")

	var syncedBlock uint64
	if len(lastSyncedBlock) > 0 {
		// To start on the next block yet to be synced add 1.
		syncedBlock = lastSyncedBlock[0].(uint64) + 1
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

	// Event listener channels.
	bondBodyTermsChan := make(chan *contracts.ChatBondBodyTerms)
	bondMotivationChan := make(chan *contracts.ChatBondMotivation)
	holderUpdateChan := make(chan *contracts.ChatHolderUpdate)
	newbondCreatedChan := make(chan *contracts.ChatNewBondCreated)
	newChatMessageChan := make(chan *contracts.ChatNewChatMessage)
	statusChangeChan := make(chan *contracts.ChatStatusChange)
	statusSignedChan := make(chan *contracts.ChatStatusSigned)

	log.Infof("Subscribing to all contract events listeners starting from block %d", syncedBlock)
	bondBodyTerms, err := s.bondChat.ChatFilterer.WatchBondBodyTerms(watchOpts, bondBodyTermsChan)
	if err != nil {
		err = fmt.Errorf("subscribing to event BondBodyTerms failed Error: %v", err)
		return err
	}

	bondMotivation, err := s.bondChat.ChatFilterer.WatchBondMotivation(watchOpts, bondMotivationChan)
	if err != nil {
		err = fmt.Errorf("subscribing to event BondMotivation failed Error: %v", err)
		return err
	}

	holderUpdate, err := s.bondChat.ChatFilterer.WatchHolderUpdate(watchOpts, holderUpdateChan)
	if err != nil {
		err = fmt.Errorf("subscribing to event HolderUpdate failed Error: %v", err)
		return err
	}

	newBondCreated, err := s.bondChat.ChatFilterer.WatchNewBondCreated(watchOpts, newbondCreatedChan)
	if err != nil {
		err = fmt.Errorf("subscribing to event NewBondCreated failed Error: %v", err)
		return err
	}

	newchatMsg, err := s.bondChat.ChatFilterer.WatchNewChatMessage(watchOpts, newChatMessageChan)
	if err != nil {
		err = fmt.Errorf("subscribing to event NewChatMessage failed Error: %v", err)
		return err
	}

	statusChange, err := s.bondChat.ChatFilterer.WatchStatusChange(watchOpts, statusChangeChan)
	if err != nil {
		err = fmt.Errorf("subscribing to event StatusChange failed Error: %v", err)
		return err
	}

	statusSigned, err := s.bondChat.ChatFilterer.WatchStatusSigned(watchOpts, statusSignedChan)
	if err != nil {
		err = fmt.Errorf("subscribing to event StatusSigned failed Error: %v", err)
		return err
	}

	// eventsSubscriptions maps the event name with their respective subscription instance.
	eventsSubscriptions := map[string]event.Subscription{
		"BondBodyTerms":  bondBodyTerms,
		"BondMotivation": bondMotivation,
		"HolderUpdate":   holderUpdate,
		"NewBondCreated": newBondCreated,
		"NewChatMessage": newchatMsg,
		"StatusChange":   statusChange,
		"StatusSigned":   statusSigned,
	}

	// quitWithErr receives the error that shutdowns the listeners and events subscription.
	quitWithErr := make(chan error)
	// eventsExit shutdowns the events subscription.
	eventsExit := make(chan struct{})

	for name, ev := range eventsSubscriptions {
		// Launches several goroutines whose purpose is to listen to errors from
		// events and pipe them into a single channel.
		go func(n string, evSub event.Subscription) {
			for {
				select {
				case <-eventsExit:
					return
				case err := <-evSub.Err():
					quitWithErr <- fmt.Errorf("event: %s error: %v", n, err)
				}
			}
		}(name, ev)
	}

	// eventsData is an buffered chan that allows each event to submit its data
	// atleast twice before its considered full.
	eventsData := make(chan *eventData, len(eventsSubscriptions)*2)

	go func() {
		for {
			select {
			case err := <-quitWithErr:
				log.Info("Sync shutdown request recieved")
				// Shutdowns the error listeners immediately the context is cancelled
				close(eventsExit)

				for _, event := range eventsSubscriptions {
					event.Unsubscribe()
				}

				// This attempts to empty all app pending write requests.
			loop:
				for {
					select {
					case <-eventsData:
						break loop
					default:
					}
				}

				if err != nil {
					log.Errorf("events data syncing ended with an error: %v", err)
				}

			case <-s.ctx.Done():
				close(quitWithErr)

			case data := <-bondBodyTermsChan:
				eventsData <- &eventData{
					method:  utils.UpdateBondBodyTerms,
					blockNo: data.Raw.BlockNumber,
					params: []interface{}{
						data.Principal, data.CouponRate, data.CouponDate,
						data.MaturityDate, data.Currency, time.Now().UTC(),
						data.Raw.BlockNumber, data.BondAddress,
					},
				}

			case data := <-bondMotivationChan:
				eventsData <- &eventData{
					method:  utils.UpdateBondMotivation,
					blockNo: data.Raw.BlockNumber,
					params: []interface{}{
						data.Message, time.Now().UTC(),
						data.Raw.BlockNumber, data.BondAddress,
					},
				}

			case data := <-holderUpdateChan:
				eventsData <- &eventData{
					method:  utils.UpdateHolder,
					blockNo: data.Raw.BlockNumber,
					params: []interface{}{
						data.Holder, time.Now().UTC(),
						data.Raw.BlockNumber, data.BondAddress,
					},
				}

			case data := <-newbondCreatedChan:
				eventsData <- &eventData{
					method:  utils.InsertNewBondCreated,
					blockNo: data.Raw.BlockNumber,
					params: []interface{}{
						data.BondAddress, data.Sender, time.Now().UTC(),
						data.Raw.BlockNumber,
					},
				}

			case data := <-newChatMessageChan:
				eventsData <- &eventData{
					method:  utils.InsertNewChatMessage,
					blockNo: data.Raw.BlockNumber,
					params: []interface{}{
						data.Sender, data.BondAddress, data.Message,
						time.Now().UTC(), data.Raw.BlockNumber,
					},
				}

			case data := <-statusChangeChan:
				eventsData <- &eventData{
					method:  utils.InsertStatusChange,
					blockNo: data.Raw.BlockNumber,
					params: []interface{}{
						data.Sender, data.BondAddress, data.Status,
						time.Now().UTC(), data.Raw.BlockNumber,
					},
				}

				eventsData <- &eventData{
					method:  utils.UpdateLastStatus,
					blockNo: data.Raw.BlockNumber,
					params: []interface{}{
						data.Status, time.Now().UTC(),
						data.Raw.BlockNumber, data.BondAddress,
					},
				}

			case data := <-statusSignedChan:
				eventsData <- &eventData{
					method:  utils.InsertStatusSigned,
					blockNo: data.Raw.BlockNumber,
					params: []interface{}{
						data.Sender, data.BondAddress, data.Status,
						time.Now().UTC(), data.Raw.BlockNumber,
					},
				}

			default:
				for info := range eventsData {
					if err := s.db.SetLocalData(info.method, info.params...); err != nil {

						// clean dirty writes made before on the current block.
						// This is a safeguard built in here but execution should never
						// here under normal circumstances.
						s.db.CleanUpLocalData(info.blockNo)

						quitWithErr <- err

						// exit the for loop so that the shutdown error can be processed.
						break
					}
				}
			}
		}
	}()

	return nil
}
