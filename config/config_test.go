package config

import (
	"os"
	"testing"
)

func TestLoadConfig(t *testing.T) {
	// Create a temporary config file
	configContent := `{
        "node_id": "test-node",
        "port": 8081,
        "peer_nodes": ["localhost:8082"],
        "data_path": "./test-data",
        "secret_key": "test-key-32-bytes-long-exactly-ok",
        "shard_count": 3
    }`

	tmpfile, err := os.CreateTemp("", "config-*.json")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpfile.Name())

	if _, err := tmpfile.Write([]byte(configContent)); err != nil {
		t.Fatal(err)
	}
	if err := tmpfile.Close(); err != nil {
		t.Fatal(err)
	}

	config, err := LoadConfig(tmpfile.Name())
	if err != nil {
		t.Fatal(err)
	}

	if config.NodeID != "test-node" {
		t.Error("Wrong node ID")
	}
	if config.Port != 8081 {
		t.Error("Wrong port")
	}
	if len(config.PeerNodes) != 1 {
		t.Error("Wrong number of peer nodes")
	}
}
