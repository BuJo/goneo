package file

import (
	"testing"
)

func TestOpeningAndClosingDb(t *testing.T) {
	db := NewDb("file.db", nil)
	defer db.Close()
	if db == nil {
		t.Fatal("Db should have been created")
	}
}
