// Copyright (c) 2023 Migwi Ndung'u
// See LICENSE for details.

package server

import (
	"fmt"
	"math/big"
	"strings"
	"time"

	"github.com/dmigwi/dhamana-protocol/client/contracts"
	"github.com/dmigwi/dhamana-protocol/client/utils"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

const (
	// blocksFilterInterval defines the number of blocks that can be filtered at ago.
	blocksFilterInterval int64 = 100

	// loggingInterval describes the intervals at which historical events
	// sync progress is logged.
	loggingInterval = 10 * time.Second

	// pollinginterval describes the intervals at which future events are polled.
	pollinginterval = 2 * time.Minute
)

var (
	// bondBodyTermsChan defines the BondBodyTerms event listener channel.
	bondBodyTermsChan = make(chan *contracts.ChatBondBodyTerms)
	// bondMotivationChan defines the BondMotivation event listener channel.
	bondMotivationChan = make(chan *contracts.ChatBondMotivation)
	// holderUpdateChan defines the HolderUpdate event listener channel.
	holderUpdateChan = make(chan *contracts.ChatHolderUpdate)
	// newBondCreatedChan defines the NewBondCreated event listener channel.
	newBondCreatedChan = make(chan *contracts.ChatNewBondCreated)
	// newChatMessageChan defines the NewChatMessage event listener channel.
	newChatMessageChan = make(chan *contracts.ChatNewChatMessage)
	// statusChangeChan defines the StatusChange event listener channel.
	statusChangeChan = make(chan *contracts.ChatStatusChange)
	// statusSignedChan defines the StatusSigned event listener channel.
	statusSignedChan = make(chan *contracts.ChatStatusSigned)

	// quit is used to indicate that a shutdown request was recieved and the
	// loop or goroutine should exit too.
	quit = make(chan struct{})

	// quitWithErr allows the error that initiates listeners and loops shutdown
	// to be sent via it.
	quitWithErr = make(chan error)

	// eventNames defines a list of all event names currently supported.
	// If a new event is introduced, it must be added here otherwise the system
	// will exit with an error when parsing the logs.
	eventNames = []string{
		"NewBondCreated", "NewChatMessage", "StatusChange", "StatusSigned",
		"BondBodyTerms", "BondMotivation", "HolderUpdate",
	}
)

// eventData contains data packed from each event recieved.
type eventData struct {
	method  utils.Method
	blockNo uint64
	params  []interface{}
}

// SyncData polls for the historical events data in a blocking operation before
// shifting to poll for future blocks asynchronously.
func (s *ServerConfig) SyncData() error {
	// fetch the block at which the contract was deployed
	deployedBlock := int64(getDeployedBlock(s.network))

	// fetch the last synced block from the database.
	lastSyncedBlock, _ := s.db.QueryLocalData(utils.GetLastSyncedBlock,
		new(lastSyncedBlockResp), "")

	var syncedBlock int64
	if len(lastSyncedBlock) > 0 {
		// To start on the next block yet to be synced add 1.
		syncedBlock = lastSyncedBlock[0].(int64) + 1
	}

	// compare the two blocks and pick the latest one.
	if deployedBlock > syncedBlock {
		syncedBlock = deployedBlock
	}

	// Fetch the events search parameters.
	topics, err := s.fetchTopics()
	if err != nil {
		return err
	}

	targetBlock, err := s.bestBlock()
	if err != nil {
		// Exit the sync, if fetching the best block failed.
		return err
	}

	endBlock := syncedBlock + blocksFilterInterval
	filterOpts := ethereum.FilterQuery{
		FromBlock: big.NewInt(syncedBlock),
		ToBlock:   big.NewInt(endBlock),
		Addresses: []common.Address{getContractAddress(s.network)},
		Topics:    topics,
	}

	// filterLogsFunc requests the filtered logs using the set filter query options.
	// It returns a count of the processed filtered logs events.
	filterLogsFunc := func() (int, error) {
		logs, err := s.backend.FilterLogs(s.ctx, filterOpts)
		if err != nil {
			return 0, fmt.Errorf("fetching logs between block %d and %d failed: %v",
				syncedBlock, endBlock, err)
		}

		return len(logs), s.parseEvents(logs)
	}

	// Process all the recieved events asynchronously by piping them into their
	// respectived channel types for further processing.
	go s.processEvents()

	// ---- Process all the historical events data in a blocking operation ----

	// Create a logging ticker timer.
	ticker := time.NewTicker(loggingInterval)

	var eventCounter, totalEvents int

	// Block till the blocks are synced to the target block.
	for endBlock <= targetBlock {
		select {
		case <-quit:
			// shutdown request was received, so exit.
			return nil
		case <-s.ctx.Done():
			// If context is shut during the looping, exit
			return nil
		case <-ticker.C:
			log.Infof("Syncing data from block=%s To target block=%d, events previously processed=%d",
				filterOpts.FromBlock.Int64(), targetBlock, eventCounter)

			totalEvents += eventCounter
			eventCounter = 0 // reset the events counter.
		default:
			// no shutdown or ticker event received.
		}

		counter, err := filterLogsFunc()
		if err != nil {
			return err
		}

		filterOpts.FromBlock = big.NewInt(endBlock)
		endBlock += blocksFilterInterval
		filterOpts.ToBlock = big.NewInt(endBlock)

		eventCounter += counter
	}

	log.Infof("Processed events=%d from start block=%d to target block=%d",
		totalEvents, syncedBlock, targetBlock)

	// ---Process asynchronously all the future events data, till shutdown ----

	// Reset the ticker timer to be used in polling the future events data.
	ticker.Reset(pollinginterval)

	syncedBlock = filterOpts.ToBlock.Int64()
	filterOpts.FromBlock.SetInt64(syncedBlock)

	// At the end of the polling interval a variable number of blocks could have
	// been added. ToBlock is set to nil so that the returned logs can have the
	// latest best block logs.
	filterOpts.ToBlock = nil

	go func() {
		for {
			select {
			case <-quit:
				// shutdown request recieved
				return
			case <-s.ctx.Done():
				// context is already cancelled.
				return
			case <-ticker.C:
				counter, err := filterLogsFunc()
				if err != nil {
					quitWithErr <- err
					return
				}

				currentBestBlock, err := s.bestBlock()
				if err != nil {
					quitWithErr <- err
					return
				}

				filterOpts.FromBlock.SetInt64(currentBestBlock)

				log.Infof("Processed events=%d upto the current best block=%d",
					counter, currentBestBlock)
			}
		}
	}()

	return nil
}

func (s *ServerConfig) processEvents() {
	// eventsData is a buffered chan that allows each event to submit its data
	// atleast twice before its considered full.
	eventsData := make(chan *eventData, len(eventNames)*2)

	go func() {
		for {
			select {
			case err := <-quitWithErr:
				log.Info("Sync shutdown request recieved")

				// trigger all other loops and listeners to close too.
				close(quit)

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

			case data := <-newBondCreatedChan:
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
			eventsDataLoop:
				for info := range eventsData {
					if err := s.db.SetLocalData(info.method, info.params...); err != nil {

						// clean dirty writes made before on the current block.
						// This is a safeguard built in here but execution should never
						// here under normal circumstances.
						s.db.CleanUpLocalData(info.blockNo)

						quitWithErr <- err

						// exit the for loop so that the shutdown error can be processed.
						break eventsDataLoop
					}
				}
			}
		}
	}()
}

// parseEvents attempts to match the returned logs with one of the event parsers.
// If none of the parsers was a postive match then an error is returned to indicate
// presence of an unsupported event.
func (s *ServerConfig) parseEvents(logs []types.Log) error {
	var eventsCount int
	for _, eventLog := range logs {
		select {
		case <-quit:
			break // Breaks the for-loop
		case <-s.ctx.Done():
			break // Breaks the for-loop
		default:
			// No shutdown request received.
		}

		eventsCount++

		newBondCreated, _ := s.bondChat.ChatFilterer.ParseNewBondCreated(eventLog)
		if newBondCreated != nil {
			newBondCreatedChan <- newBondCreated
			continue
		}

		newChatMessage, _ := s.bondChat.ChatFilterer.ParseNewChatMessage(eventLog)
		if newChatMessage != nil {
			newChatMessageChan <- newChatMessage
			continue
		}

		statusChange, _ := s.bondChat.ChatFilterer.ParseStatusChange(eventLog)
		if statusChange != nil {
			statusChangeChan <- statusChange
			continue
		}

		statusSigned, _ := s.bondChat.ChatFilterer.ParseStatusSigned(eventLog)
		if statusSigned != nil {
			statusSignedChan <- statusSigned
			continue
		}

		bondBodyTerms, _ := s.bondChat.ChatFilterer.ParseBondBodyTerms(eventLog)
		if bondBodyTerms != nil {
			bondBodyTermsChan <- bondBodyTerms
			continue
		}

		bondMotivation, _ := s.bondChat.ChatFilterer.ParseBondMotivation(eventLog)
		if bondMotivation != nil {
			bondMotivationChan <- bondMotivation
			continue
		}

		holderUpdate, _ := s.bondChat.ChatFilterer.ParseHolderUpdate(eventLog)
		if holderUpdate != nil {
			holderUpdateChan <- holderUpdate
			continue
		}

		// If one of the parsers failed to return a positive match then there
		// must be an unsupported event in the returned logs.
		return fmt.Errorf("unsupported event at contract address: %v and Block No: %v ",
			eventLog.Address, eventLog.BlockNumber)
	}
	return nil
}

// bestBlock returns the current chain best block. In case of an error,
// -1 is returned.
func (s *ServerConfig) bestBlock() (int64, error) {
	targetHeader, err := s.backend.HeaderByNumber(s.ctx, nil)
	if err != nil {
		return -1, fmt.Errorf("fetching the current bestblock failed: %v", err)
	}
	return targetHeader.Number.Int64(), nil
}

// fetchTopics returns the search parameters for all the supported events,
// this helps to narrow down the events topics search to only the supported
// events.
func (s *ServerConfig) fetchTopics() ([][]common.Hash, error) {
	chatABI, err := abi.JSON(strings.NewReader(contracts.ChatABI))
	if err != nil {
		return nil, fmt.Errorf("unable to parse the ABI interface: %v", err)
	}

	var query [][]interface{}
	for _, n := range eventNames {
		query = append([][]interface{}{{chatABI.Events[n].ID}}, query...)
	}

	topics, err := abi.MakeTopics(query...)
	if err != nil {
		return nil, fmt.Errorf("unable generate event topics: %v", err)
	}

	return topics, nil
}
