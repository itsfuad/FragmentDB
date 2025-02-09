package node

import (
	"encoding/json"
	"fmt"
	"fragmentdb/config"
	"fragmentdb/shard"
	"io"
	"net/http"
	"strings"
	"sync"
)

const (
	METHOD_NOT_ALLOWED = "Method not allowed"
)

type Node struct {
	ID       string
	shardMgr *shard.ShardManager
	Data     map[string][]byte
	Peers    []string
	Mu       sync.RWMutex
	config   *config.ServerConfig
}

func NewNode(id string, peers []string, shardCount int, cfg *config.ServerConfig) *Node {
	return &Node{
		ID:       id,
		shardMgr: shard.NewShardManager(shardCount, id),
		Data:     make(map[string][]byte),
		Peers:    peers,
		config:   cfg,
	}
}

func (n *Node) Start(port int) error {
	mux := http.NewServeMux()
	mux.HandleFunc("/put", n.handlePut)
	mux.HandleFunc("/get/", n.handleGet) // Note the trailing slash
	mux.HandleFunc("/sync", n.handleSync)
	// Add other CRUD handlers
	return http.ListenAndServe(fmt.Sprintf(":%d", port), mux)
}

func (n *Node) handleGet(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, METHOD_NOT_ALLOWED, http.StatusMethodNotAllowed)
		return
	}

	// Parse key from URL path: /get/keyname
	key := strings.TrimPrefix(r.URL.Path, "/get/")
	if key == "" {
		http.Error(w, "Key not provided", http.StatusBadRequest)
		return
	}

	var result []byte
	n.Mu.RLock()
	for i := 0; i < n.shardMgr.TotalShards; i++ {
		shardKey := fmt.Sprintf("%s:%d", key, i)
		if shard, exists := n.Data[shardKey]; exists {
			decrypted, err := n.shardMgr.DecryptShard(shard, n.config.SecretKey)
			if err != nil {
				n.Mu.RUnlock()
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			result = append(result, decrypted...)
		}
	}
	n.Mu.RUnlock()

	w.Header().Set("Content-Type", "application/json")
	// Return result as a string instead of raw bytes.
	json.NewEncoder(w).Encode(map[string]interface{}{
		"key":   key,
		"value": string(result),
	})
}

func (n *Node) handlePut(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, METHOD_NOT_ALLOWED, http.StatusMethodNotAllowed)
		return
	}

	// Change the Value field to string so JSON sends a normal string.
	var data struct {
		Key   string `json:"key"`
		Value string `json:"value"`
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	if err := json.Unmarshal(body, &data); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Convert value string to []byte
	valueBytes := []byte(data.Value)
	shards := n.shardMgr.SplitData(valueBytes)
	n.Mu.Lock()
	for shardID, shardData := range shards {
		encrypted, err := n.shardMgr.EncryptShard(shardData, n.config.SecretKey)
		if err != nil {
			n.Mu.Unlock()
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		n.Data[fmt.Sprintf("%s:%d", data.Key, shardID)] = encrypted
	}
	n.Mu.Unlock()

	w.WriteHeader(http.StatusOK)
}

func (n *Node) handleSync(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, METHOD_NOT_ALLOWED, http.StatusMethodNotAllowed)
		return
	}

	n.Mu.RLock()
	defer n.Mu.RUnlock()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(n.Data)
}
