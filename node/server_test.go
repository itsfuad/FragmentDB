package node

import (
	"encoding/json"
	"fragmentdb/config"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func setupTestNode() *Node {
	cfg := &config.ServerConfig{
		NodeID:     "test-node",
		Port:       8081,
		PeerNodes:  []string{},
		ShardCount: 3,
		SecretKey:  "12345678901234567890123456789012",
	}
	return NewNode(cfg.NodeID, cfg.PeerNodes, cfg.ShardCount, cfg)
}

func TestNodeHandlePut(t *testing.T) {

	node := setupTestNode()

	payload := `{"key":"test-key","value":"test-value"}`
	req := httptest.NewRequest("POST", "/put", strings.NewReader(payload))
	w := httptest.NewRecorder()

	node.handlePut(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status OK, got %v", w.Code)
	}
}

func TestNodeHandleGet(t *testing.T) {
	node := setupTestNode()

	// First put some data
	putPayload := `{"key":"test-key","value":"test-value"}`
	putReq := httptest.NewRequest("POST", "/put", strings.NewReader(putPayload))
	putW := httptest.NewRecorder()
	node.handlePut(putW, putReq)

	// Then try to get it
	req := httptest.NewRequest("GET", "/get/test-key", nil)
	w := httptest.NewRecorder()
	node.handleGet(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status OK, got %v", w.Code)
	}

	var response map[string]interface{}
	if err := json.NewDecoder(w.Body).Decode(&response); err != nil {
		t.Fatal(err)
	}

	if response["key"] != "test-key" {
		t.Error("Wrong key in response")
	}
}
