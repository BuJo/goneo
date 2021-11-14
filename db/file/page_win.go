// +build windows

package file

import (
	"errors"
)

type PageStore struct {
}

func NewPageStore(filename string) (*PageStore, error) {
	return nil, errors.New("not implemented")
}

func (ps *PageStore) NumPages() int {
	return 0
}

func (ps *PageStore) GetFreePage() (int, error) {
	return 0, errors.New("not implemented")
}

func (ps *PageStore) GetPage(pgnum int) ([]byte, error) {
	return nil, errors.New("not implemented")
}

func (ps *PageStore) AddPage() (err error) {
	return errors.New("not implemented")
}
