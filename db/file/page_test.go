package file

import (
	"os"
	"testing"
)

func TestPage(t *testing.T) {
	ps, err := NewPageStore("test-page0.db")
	if err != nil {
		t.Fatal(err.Error())
	}
	if ps == nil {
		t.Fatal("Should have a page store")
	}
	defer os.Remove("test-page0.db")

	page, perr := ps.GetPage(0)
	if perr == nil {
		t.Fatal("Should not yet have a page")
	}

	aerr := ps.AddPage()
	if aerr != nil {
		t.Fatal("adding page should work: " + aerr.Error())
	}

	page, perr = ps.GetPage(0)
	if perr != nil {
		t.Fatal(perr.Error())
	}
	if page == nil {
		t.Fatal("Should have a page")
	}
}
