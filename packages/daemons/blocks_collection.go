// Copyright 2016 The go-daylight Authors
// This file is part of the go-daylight library.
//
// The go-daylight library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The go-daylight library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the go-daylight library. If not, see <http://www.gnu.org/licenses/>.

package daemons

import (
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"time"

	"github.com/GenesisKernel/go-genesis/packages/network"
	"github.com/GenesisKernel/go-genesis/packages/network/tcpclient"

	"github.com/GenesisKernel/go-genesis/packages/conf"
	"github.com/GenesisKernel/go-genesis/packages/conf/syspar"
	"github.com/GenesisKernel/go-genesis/packages/consts"
	"github.com/GenesisKernel/go-genesis/packages/model"
	"github.com/GenesisKernel/go-genesis/packages/parser"
	"github.com/GenesisKernel/go-genesis/packages/service"
	"github.com/GenesisKernel/go-genesis/packages/utils"

	log "github.com/sirupsen/logrus"
)

// ErrNodesUnavailable is returned when all nodes is unavailable
var ErrNodesUnavailable = errors.New("All nodes unavailable")

// BlocksCollection collects and parses blocks
func BlocksCollection(ctx context.Context, d *daemon) error {
	if ctx.Err() != nil {
		d.logger.WithFields(log.Fields{"type": consts.ContextError, "error": ctx.Err()}).Error("context error")
		return ctx.Err()
	}

	return blocksCollection(ctx, d)
}

func InitialLoad(logger *log.Entry) error {

	// check for initial load
	toLoad, err := needLoad(logger)
	if err != nil {
		return err
	}

	if toLoad {
		logger.Debug("start first block loading")

		if err := firstLoad(logger); err != nil {
			return err
		}
	}

	return nil
}

func blocksCollection(ctx context.Context, d *daemon) (err error) {
	host, maxBlockID, err := getHostWithMaxID(d.logger)
	if err != nil {
		d.logger.WithFields(log.Fields{"error": err}).Error("on checking best host")
		return err
	}

	infoBlock := &model.InfoBlock{}
	found, err := infoBlock.Get()
	if err != nil {
		log.WithFields(log.Fields{"type": consts.DBError, "error": err}).Error("Getting cur blockID")
		return err
	}
	if !found {
		log.WithFields(log.Fields{"type": consts.NotFound, "error": err}).Error("Info block not found")
		return errors.New("Info block not found")
	}

	if infoBlock.BlockID >= maxBlockID {
		log.WithFields(log.Fields{"blockID": infoBlock.BlockID, "maxBlockID": maxBlockID}).Debug("Max block is already in the host")
		return nil
	}

	DBLock()
	defer func() {
		DBUnlock()
		service.NodeDoneUpdatingBlockchain()
	}()

	// update our chain till maxBlockID from the host
	return UpdateChain(ctx, d, host, maxBlockID)
}

// UpdateChain load from host all blocks from our last block to maxBlockID
func UpdateChain(ctx context.Context, d *daemon, host string, maxBlockID int64) error {
	var (
		err   error
		count int
	)

	// get current block id from our blockchain
	curBlock := &model.InfoBlock{}
	if _, err = curBlock.Get(); err != nil {
		d.logger.WithFields(log.Fields{"type": consts.DBError, "error": err}).Error("Getting info block")
		return err
	}

	if ctx.Err() != nil {
		d.logger.WithFields(log.Fields{"type": consts.ContextError, "error": ctx.Err()}).Error("context error")
		return ctx.Err()
	}

	playRawBlock := func(rb []byte) error {

		block, err := parser.ProcessBlockWherePrevFromBlockchainTable(rb, true)
		defer func() {
			if err != nil {
				d.logger.WithFields(log.Fields{"error": err, "type": consts.BlockError}).Error("retrieving blockchain from node")
				banNode(host, block, err)
			}
		}()

		if err != nil {
			d.logger.WithFields(log.Fields{"error": err, "type": consts.BlockError}).Error("processing block")
			return err
		}

		// hash compare could be failed in the case of fork
		hashMatched, thisErrIsOk := block.CheckHash()
		if thisErrIsOk != nil {
			d.logger.WithFields(log.Fields{"error": err, "type": consts.BlockError}).Error("checking block hash")
		}

		if !hashMatched {
			//it should be fork, replace our previous blocks to ones from the host
			err = parser.GetBlocks(block.Header.BlockID-1, host)
			if err != nil {
				d.logger.WithFields(log.Fields{"error": err, "type": consts.ParserError}).Error("processing block")
				return err
			}
		}

		block.PrevHeader, err = parser.GetBlockDataFromBlockChain(block.Header.BlockID - 1)
		if err != nil {
			return utils.ErrInfo(fmt.Errorf("can't get block %d", block.Header.BlockID-1))
		}

		if err = block.CheckBlock(); err != nil {
			return err
		}

		if err = block.PlayBlockSafe(); err != nil {
			return err
		}

		return nil
	}

	st := time.Now()
	d.logger.Infof("starting downloading blocks from %d to %d (%d) \n", curBlock.BlockID, maxBlockID, maxBlockID-curBlock.BlockID)

	for blockID := curBlock.BlockID + 1; blockID <= maxBlockID; blockID += int64(network.BlocksPerRequest) {
		rawBlocksChan, err := tcpclient.GetBlocksBodies(host, blockID, false)
		if err != nil {
			d.logger.WithFields(log.Fields{"error": err, "type": consts.BlockError}).Error("getting block body")
			return err
		}

		for rawBlock := range rawBlocksChan {
			if err = playRawBlock(rawBlock); err != nil {
				d.logger.WithFields(log.Fields{"error": err, "type": consts.BlockError}).Error("playing raw block")
				return err
			}
			count++
		}

		d.logger.Infof("%d blocks was collected (%s) \n", count, time.Since(st).String())
	}
	return nil
}

// init first block from file or from embedded value
func loadFirstBlock(logger *log.Entry) error {
	newBlock, err := ioutil.ReadFile(conf.Config.FirstBlockPath)
	if err != nil {
		logger.WithFields(log.Fields{
			"type": consts.IOError, "error": err, "path": conf.Config.FirstBlockPath,
		}).Error("reading first block from file")
	}

	if err = parser.InsertBlockWOForks(newBlock, false, true); err != nil {
		logger.WithFields(log.Fields{"type": consts.ParserError, "error": err}).Error("inserting new block")
		return err
	}

	return nil
}

func firstLoad(logger *log.Entry) error {
	DBLock()
	defer DBUnlock()

	return loadFirstBlock(logger)
}

func needLoad(logger *log.Entry) (bool, error) {
	infoBlock := &model.InfoBlock{}
	_, err := infoBlock.Get()
	if err != nil {
		logger.WithFields(log.Fields{"error": err, "type": consts.DBError}).Error("getting info block")
		return false, err
	}
	// we have empty blockchain, we need to load blockchain from file or other source
	if infoBlock.BlockID == 0 {
		logger.Debug("blockchain should be loaded")
		return true, nil
	}
	return false, nil
}

func banNode(host string, block *parser.Block, err error) {
	var (
		reason             string
		blockId, blockTime int64
	)
	if err != nil {
		reason = err.Error()
	}

	if block != nil {
		blockId, blockTime = block.Header.BlockID, block.Header.Time
	} else {
		blockId, blockTime = -1, time.Now().Unix()
	}

	log.WithFields(log.Fields{"reason": reason, "host": host, "block_id": blockId, "block_time": blockTime}).Debug("ban node")

	n, err := syspar.GetNodeByHost(host)
	if err != nil {
		log.WithFields(log.Fields{"error": err}).Error("getting node by host")
		return
	}

	err = service.GetNodesBanService().RegisterBadBlock(n, blockId, blockTime, reason)
	if err != nil {
		log.WithFields(log.Fields{"error": err, "node": n.KeyID, "block": blockId}).Error("registering bad block from node")
	}
}

// GetHostWithMaxID returns host with maxBlockID
func getHostWithMaxID(logger *log.Entry) (host string, maxBlockID int64, err error) {

	nbs := service.GetNodesBanService()
	hosts, err := nbs.FilterBannedHosts(syspar.GetRemoteHosts())
	if err != nil {
		logger.WithFields(log.Fields{"error": err}).Error("on filtering banned hosts")
	}

	host, maxBlockID, err = tcpclient.HostWithMaxBlock(hosts)
	if err != nil && err == tcpclient.ErrNodesUnavailable {
		hosts, err := nbs.FilterBannedHosts(conf.GetNodesAddr())
		if err != nil {
			logger.WithFields(log.Fields{"error": err}).Error("on filtering banned hosts")
			return "", -1, err
		}
		return tcpclient.HostWithMaxBlock(hosts)
	}

	return
}
