package txdb

import (
	"github.com/coreos/bbolt"
)

// CoinsViewDB -
type CoinsViewDB struct {
	DB bolt.DB
}

// BlockTreeDB -
type BlockTreeDB struct{}
