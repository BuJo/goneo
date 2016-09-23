package file

import (
	"log"
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

func TestWritingPages(t *testing.T) {
	ps, err := NewPageStore("test-page0.db")
	if err != nil {
		t.Fatal(err.Error())
	}
	if ps == nil {
		t.Fatal("Should have a page store")
	}
	defer os.Remove("test-page0.db")

	for i := 0; i < 10; i++ {
		err = ps.AddPage()
		if err != nil {
			t.Fatal("Should be able to add page")
		}

		page, perr := ps.GetPage(i)
		if perr != nil {
			t.Fatal("Should be able to get added page")
		}
		log.Printf("page %d(%d->%d) len:%d, cap:%d, pagesize:%d", i, HEADER_SIZE+PAGE_SIZE*i, HEADER_SIZE+PAGE_SIZE*i+PAGE_SIZE, len(page), cap(page), PAGE_SIZE)

		for chunknr := 0; chunknr*128 < PAGE_SIZE-50; chunknr++ {
			log.Printf("page %d offset:%d len:%d, cap:%d, pagesize:%d", i, chunknr*128, len(page[chunknr*128:]), cap(page[chunknr*128:]), PAGE_SIZE)
			copy(page[chunknr*128:], "||>"+string([]byte{byte(39 + (chunknr % 51))})+"<||")
		}
	}
}
