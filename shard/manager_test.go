package shard

import (
	"bytes"
	"testing"
)

func TestShardManager_SplitData(t *testing.T) {
	sm := NewShardManager(3, "test-node")
	testData := []byte("Hello, World!")

	shards := sm.SplitData(testData)

	if len(shards) != 3 {
		t.Errorf("Expected 3 shards, got %d", len(shards))
	}

	// Reconstruct data
	var reconstructed []byte
	for i := 0; i < 3; i++ {
		reconstructed = append(reconstructed, shards[i]...)
	}

	if !bytes.Equal(testData, reconstructed) {
		t.Errorf("Data reconstruction failed")
	}
}

func TestShardManager_Encryption(t *testing.T) {
	sm := NewShardManager(3, "test-node")
	testData := []byte("Secret message")
	key := "12345678901234567890123456789012" // 32-byte key

	encrypted, err := sm.EncryptShard(testData, key)
	if err != nil {
		t.Fatal(err)
	}

	decrypted, err := sm.DecryptShard(encrypted, key)
	if err != nil {
		t.Fatal(err)
	}

	if !bytes.Equal(testData, decrypted) {
		t.Error("Encryption/decryption cycle failed")
	}
}
