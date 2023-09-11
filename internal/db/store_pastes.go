package db

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/dgraph-io/badger/v4"
	"github.com/merlinfuchs/vaultbin/internal/common"
	"github.com/merlinfuchs/vaultbin/internal/store"
)

func (db *DB) CreatePaste(content, language string, ttl time.Duration) (*store.Paste, error) {
	encryptionKey, err := common.GenerateEncryptionKey()
	if err != nil {
		return nil, fmt.Errorf("generating encryption key failed: %w", err)
	}

	id := common.EncodeEncryptionKey(encryptionKey)
	paste := &store.Paste{
		ID:        id,
		Content:   content,
		Language:  language,
		CreatedAt: time.Now().UTC(),
		ExpiresAt: time.Now().UTC().Add(ttl),
	}

	val, err := json.Marshal(paste)
	if err != nil {
		return nil, fmt.Errorf("json marshal failed: %w", err)
	}

	val, err = common.EncryptBytes(encryptionKey, val)
	if err != nil {
		return nil, fmt.Errorf("encryption failed: %w", err)
	}

	keyHash, err := common.HashEncryptionKey(encryptionKey)
	if err != nil {
		return nil, fmt.Errorf("hashing encryption key failed: %w", err)
	}

	key := fmt.Sprintf("pastes.%s", keyHash)
	err = db.bg.Update(func(txn *badger.Txn) error {
		e := badger.NewEntry([]byte(key), val).WithTTL(ttl)
		return txn.SetEntry(e)
	})
	if err != nil {
		return nil, fmt.Errorf("db transaction failed: %w", err)
	}

	return paste, nil
}

func (db *DB) Paste(id string) (*store.Paste, error) {
	encryptionKey, err := common.DecodeEncryptionKey(id)
	if err != nil {
		return nil, nil
	}

	var res *store.Paste

	keyHash, err := common.HashEncryptionKey(encryptionKey)
	if err != nil {
		return nil, fmt.Errorf("hashing encryption key failed: %w", err)
	}

	key := fmt.Sprintf("pastes.%s", keyHash)
	err = db.bg.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte(key))
		if err != nil {
			if err == badger.ErrKeyNotFound {
				return nil
			}
			return err
		}

		return item.Value(func(val []byte) error {
			val, err = common.DecryptBytes(encryptionKey, val)
			if err != nil {
				return err
			}
			return json.Unmarshal(val, &res)
		})
	})
	if err != nil {
		return nil, fmt.Errorf("db transaction failed: %w", err)
	}

	return res, nil
}

func (db *DB) CountPasteView(id string) (int, error) {
	encryptionKey, err := common.DecodeEncryptionKey(id)
	if err != nil {
		return 0, fmt.Errorf("failed decoding encryption key: %w", err)
	}

	keyHash, err := common.HashEncryptionKey(encryptionKey)
	if err != nil {
		return 0, fmt.Errorf("hashing encryption key failed: %w", err)
	}

	var res uint64

	key := fmt.Sprintf("pastes.views.%s", keyHash)
	err = db.bg.Update(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte(key))
		if err != nil {
			if err != badger.ErrKeyNotFound {
				return err
			}
		} else {
			err = item.Value(func(val []byte) error {
				res = bytesToUint64(val) + 1
				return nil
			})
			if err != nil {
				return err
			}
		}

		return txn.Set([]byte(key), uint64ToBytes(res))
	})
	if err != nil {
		return 0, fmt.Errorf("db transaction failed: %w", err)
	}

	return int(res), nil
}
