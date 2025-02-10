package recovery

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"fragmentdb/config"
	"fragmentdb/node"
)

// fakeSyncHandler returns predefined JSON data for a sync request.
func fakeSyncHandler(data map[string][]byte) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(data)
	}
}

func TestSynchronizeWithPeers(t *testing.T) {
	// Prepare fake data to be returned by the peer server.
	dummyData := map[string][]byte{"test-key:0": []byte("partial data")}
	peerServer := httptest.NewServer(fakeSyncHandler(dummyData))
	defer peerServer.Close()

	// Remove "http://" prefix to simulate peer address.
	peerAddr := strings.TrimPrefix(peerServer.URL, "http://")

	// Setup a test configuration and node.
	cfg := &config.ServerConfig{
		NodeID:     "test-node",
		Port:       8080,
		PeerNodes:  []string{peerAddr},
		ShardCount: 1,
		SecretKey:  "12345678901234567890123456789012",
	}
	n := node.NewNode(cfg.NodeID, cfg.PeerNodes, cfg.ShardCount, cfg)
	syncMgr := NewSyncManager(n)

	// Call the sync method.
	syncMgr.synchronizeWithPeers()

	// Wait briefly to allow goroutines to complete.
	time.Sleep(100 * time.Millisecond)

	// Verify that data from the peer was merged into node data.
	n.Mu.RLock()
	defer n.Mu.RUnlock()
	if _, exists := n.Data["test-key:0"]; !exists {
		t.Error("Data not synced from peer")
	}
}
