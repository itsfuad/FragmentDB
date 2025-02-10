package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"fragmentdb/config"
	"fragmentdb/node"
	"fragmentdb/recovery"
)

func main() {
	cfg, err := config.LoadConfig("config.json")
	if err != nil {
		log.Fatal(err)
	}

	n := node.NewNode(cfg.NodeID, cfg.PeerNodes, cfg.ShardCount, cfg)
	syncMgr := recovery.NewSyncManager(n)

	// Setup context for graceful shutdown.
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Start sync manager with context.
	syncMgr.StartSync(ctx)

	// Start node server concurrently.
	go func() {
		if err := n.Start(cfg.Port); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server error: %v", err)
		}
	}()

	log.Printf("Starting node %s on port %d", cfg.NodeID, cfg.Port)

	// Wait for a termination signal.
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	<-sigs

	log.Println("Shutdown signal received")

	// Gracefully shutdown node server.
	shutdownCtx, shutdownCancel := context.WithTimeout(ctx, 10*time.Second)
	defer shutdownCancel()
	if err := n.Stop(shutdownCtx); err != nil {
		log.Printf("Error during shutdown: %v", err)
	}

	cancel() // Stop the sync job.
	log.Println("Shutdown complete")
}
