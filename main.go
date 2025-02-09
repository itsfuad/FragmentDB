package main

import (
	"log"

	"fragmentdb/config"
	"fragmentdb/node"
	"fragmentdb/recovery"
)

func main() {
	cfg, err := config.LoadConfig("config.json")
	if err != nil {
		log.Fatal(err)
	}

	node := node.NewNode(cfg.NodeID, cfg.PeerNodes, cfg.ShardCount, cfg)
	syncMgr := recovery.NewSyncManager(node)
	syncMgr.StartSync()

	log.Printf("Starting node %s on port %d", cfg.NodeID, cfg.Port)
	if err := node.Start(cfg.Port); err != nil {
		log.Fatal(err)
	}
}
