package recovery

import (
	"encoding/json"
	"fragmentdb/node"
	"net/http"
	"sync"
	"time"
)

type SyncManager struct {
	node     *node.Node
	interval time.Duration
	mu       sync.Mutex
}

func NewSyncManager(node *node.Node) *SyncManager {
	return &SyncManager{
		node:     node,
		interval: time.Minute * 5,
	}
}

func (sm *SyncManager) StartSync() {
	go func() {
		for {
			sm.synchronizeWithPeers()
			time.Sleep(sm.interval)
		}
	}()
}

func (sm *SyncManager) synchronizeWithPeers() {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	for _, peer := range sm.node.Peers {
		go func(peerAddr string) {
			resp, err := http.Get("http://" + peerAddr + "/sync")
			if err != nil {
				return
			}
			defer resp.Body.Close()

			var peerData map[string][]byte
			if err := json.NewDecoder(resp.Body).Decode(&peerData); err != nil {
				return
			}

			sm.node.Mu.Lock()
			for k, v := range peerData {
				if _, exists := sm.node.Data[k]; !exists {
					sm.node.Data[k] = v
				}
			}
			sm.node.Mu.Unlock()
		}(peer)
	}
}
