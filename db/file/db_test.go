package file

import (
	"testing"
)

func TestNewDb(t *testing.T) {
	db := NewDb("file.db", nil)
	if db == nil {
		t.Fatal("Db should have been created")
	}
}
