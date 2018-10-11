package sync

import (
	"testing"
)

func TestSync(t *testing.T) {
	node := NewNode()
	node.Sync()
}
