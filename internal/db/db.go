package db

import (
	"encoding/binary"
	"fmt"

	"github.com/dgraph-io/badger/v4"
	"github.com/merlinfuchs/vaultbin/internal/config"
)

type DB struct {
	bg *badger.DB
}

func New() (*DB, error) {
	db, err := badger.Open(badger.DefaultOptions(config.K.String("db_path")))
	if err != nil {
		return nil, fmt.Errorf("error opening badger db: %w", err)
	}

	return &DB{
		bg: db,
	}, nil
}

func uint64ToBytes(i uint64) []byte {
	var buf [8]byte
	binary.BigEndian.PutUint64(buf[:], i)
	return buf[:]
}

func bytesToUint64(b []byte) uint64 {
	return binary.BigEndian.Uint64(b)
}

func mergeOperatorAdd(existing, new []byte) []byte {
	return uint64ToBytes(bytesToUint64(existing) + bytesToUint64(new))
}
