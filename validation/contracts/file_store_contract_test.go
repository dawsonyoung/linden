//go:build contracts

package contracts

import (
	"testing"

	"github.com/dawsonyoung/linden/storage"
)

func Test_FileStore_SatisfiesStoreContract(t *testing.T) {
	runStoreContract(t, func(t *testing.T) storage.Store {
		tmpDir := t.TempDir()
		store, err := storage.NewFileStore(tmpDir)
		if err != nil {
			t.Fatalf("failed to create FileStore: %v", err)
		}
		return store
	})
}
