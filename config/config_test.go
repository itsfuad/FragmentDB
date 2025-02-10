package config

import (
	"os"
	"testing"
)

func TestLoadConfigValidConfig(t *testing.T) {
	configContent := `{
		"node_id": "test-node",
		"port": 8081,
		"peer_nodes": ["localhost:8082"],
		"data_path": "./test-data",
		"secret_key": "12345678901234567890123456789012",
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

func TestLoadConfigInvalidJSON(t *testing.T) {
	tmpfile, err := os.CreateTemp("", "invalid-config-*.json")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpfile.Name())

	invalidJSON := `{"node_id": "test-node", "port": "not-a-number"`
	if _, err := tmpfile.Write([]byte(invalidJSON)); err != nil {
		t.Fatal(err)
	}
	if err := tmpfile.Close(); err != nil {
		t.Fatal(err)
	}

	_, err = LoadConfig(tmpfile.Name())
	if err == nil {
		t.Error("Expected error for invalid JSON, got nil")
	}
}

func TestLoadConfigMissingSecretKey(t *testing.T) {
	configContent := `{
		"node_id": "test-node",
		"port": 8081,
		"peer_nodes": [],
		"data_path": "./test-data",
		"secret_key": "",
		"shard_count": 3
	}`
	tmpfile, err := os.CreateTemp("", "missing-secretkey-*.json")
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

	_, err = LoadConfig(tmpfile.Name())
	if err == nil {
		t.Error("Expected error for missing secret_key, got nil")
	}
}
