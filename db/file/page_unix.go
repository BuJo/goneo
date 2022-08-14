//go:build linux || darwin

package file

// #include <unistd.h>
import "C"

import (
	"encoding/binary"
	"errors"
	"reflect"
	"syscall"
	"unsafe"

	"github.com/BuJo/goneo/log"
)

const (
	HEADER_SIZE = 512
	PAGE_SIZE   = 128
)

type PageStore struct {
	size    int64
	fd      int
	backing uintptr

	freepage int
}

func NewPageStore(filename string) (*PageStore, error) {
	ps := new(PageStore)

	oerr := ps.openFile(filename)
	if oerr != nil {
		return nil, errors.New("Could not open file: " + oerr.Error())
	}

	if ps.size == 0 {
		fterr := ps.resizeFile(HEADER_SIZE)
		if fterr != nil {
			return nil, errors.New("Could not initialize new file: " + fterr.Error())
		}
	}

	maperr := ps.mapFile()
	if maperr != nil {
		return nil, errors.New("Could not mmap: " + maperr.Error())
	}

	return ps, nil
}

func (ps *PageStore) NumPages() int {
	log.Printf("Have %d pages", int((ps.size-HEADER_SIZE)/PAGE_SIZE))

	return int((ps.size - HEADER_SIZE) / PAGE_SIZE)
}

func (ps *PageStore) GetFreePage() (int, error) {
	if !ps.haveFreePage() {
		ps.AddPage()
	}
	return ps.getNextFreePage()
}

func (ps *PageStore) haveFreePage() bool {
	return ps.freepage > 0
}

func (ps *PageStore) getNextFreePage() (int, error) {
	if ps.haveFreePage() {
		freepage := ps.freepage
		p, err := ps.GetPage(freepage)
		if err != nil {
			return -1, err
		}

		ps.freepage = int(binary.LittleEndian.Uint64(p))

		return freepage, nil
	}
	return -1, nil
}

func (ps *PageStore) GetPage(pgnum int) ([]byte, error) {
	if pgnum < 0 {
		return nil, errors.New("page number must be greater than zero")
	}
	if numpages := ps.NumPages(); pgnum > numpages-1 {
		return nil, errors.New("page number too high")
	}

	var s []byte
	sp := (*reflect.SliceHeader)(unsafe.Pointer(&s))
	sp.Data = ps.backing + HEADER_SIZE + PAGE_SIZE*uintptr(pgnum)
	sp.Len, sp.Cap = PAGE_SIZE, PAGE_SIZE

	return s, nil
}

func (ps *PageStore) AddPage() (err error) {
	err = ps.resizeFile(PAGE_SIZE)
	if err == nil {
		err = ps.remapFile()
		if err != nil {
			ps.resizeFile(-PAGE_SIZE)
		}

		var page []byte
		var freepage = ps.NumPages() - 1

		page, err = ps.GetPage(freepage)
		binary.LittleEndian.PutUint64(page, uint64(ps.freepage))
		ps.freepage = freepage
	}
	return err
}

func (ps *PageStore) mapFile() error {

	addr, errno := mmap(0, ps.size, syscall.PROT_READ|syscall.PROT_WRITE, syscall.MAP_SHARED|syscall.MAP_FILE, ps.fd, 0)
	if errno != nil {
		return errno
	}

	ps.backing = addr

	log.Printf("Mapped 0x%08x length %d", ps.backing, ps.size)

	return nil
}

func (ps *PageStore) remapFile() error {

	addr, errno := mmap(ps.backing, ps.size, syscall.PROT_READ|syscall.PROT_WRITE, syscall.MAP_SHARED|syscall.MAP_FILE|syscall.MAP_FIXED, ps.fd, 0)
	if errno != nil {
		return errno
	}
	if ps.backing != addr {
		return errors.New("can not change addr")
	}

	log.Printf("Remapped 0x%08x to length %d", ps.backing, ps.size)

	return nil
}

func (ps *PageStore) openFile(filename string) error {
	fd, fderr := syscall.Open(filename, syscall.O_CREAT|syscall.O_RDWR, 0600)
	if fderr != nil {
		return fderr
	}
	ps.fd = fd

	var stat syscall.Stat_t
	staterr := syscall.Fstat(fd, &stat)
	if staterr != nil {
		return staterr
	}

	ps.size = stat.Size

	return nil
}

func (ps *PageStore) resizeFile(size int64) error {
	log.Printf("resizing paged file(0x%08x) from %d to %d", ps.backing, ps.size, ps.size+size)

	fterr := syscall.Ftruncate(ps.fd, ps.size+size)
	if fterr != nil {
		return fterr
	}
	ps.size = ps.size + size
	return nil
}

func mmap(addr uintptr, length int64, prot, flags uintptr, fd int, offset uintptr) (uintptr, error) {

	if pgsiz := uintptr(C.getpagesize()); flags&syscall.MAP_FIXED > 0 && addr%pgsiz != 0 {
		return 0, errors.New("addr should be page aligned")
	}
	if flags&(syscall.MAP_PRIVATE|syscall.MAP_SHARED) == 0 {
		return 0, errors.New("flags should include either anon or shared")
	}
	if length <= 0 {
		return 0, errors.New("len should be > 0")
	}
	if pgsiz := uintptr(C.getpagesize()); offset%pgsiz != 0 {
		return 0, errors.New("offset should be page aligned")
	}
	if flags&syscall.MAP_ANON > 0 && fd != 0 {
		return 0, errors.New("anonymous mapping and no support for vm tags")
	}

	xaddr, _, errno := syscall.Syscall6(syscall.SYS_MMAP, addr, uintptr(length), prot, flags, uintptr(fd), offset)
	if errno != 0 {
		return 0, errno
	}
	return xaddr, nil
}

func munmap(addr, len uintptr) error {
	_, _, errno := syscall.Syscall(syscall.SYS_MUNMAP, addr, len, 0)
	if errno != 0 {
		return errno
	}
	return nil
}
