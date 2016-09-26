package file

import (
	"encoding/binary"
	"errors"
	. "github.com/BuJo/goneo/db"
	"log"
	"sync"
	"sync/atomic"
)

type nodeentry struct {
	pagenum, offset int
	rts, wts        uint64
	prev            *nodeentry
}
type nodeentrymap map[uint64]*nodeentry
type filedb struct {
	name    string
	options map[string][]string

	pagestore *PageStore

	idmap nodeentrymap

	wrmutex  sync.Mutex
	nextfree struct{ pagenum, offset int }
	nextid   uint64

	tid uint64
}

func NewDb(name string, options map[string][]string) (DatabaseService, error) {
	db := new(filedb)

	db.name = name
	db.options = options

	var err error
	db.pagestore, err = NewPageStore(db.name)

	if db.pagestore.NumPages() > 0 {
		err = db.loadNodes()
	} else {
		err = db.initializeNodestructure()
	}

	return db, err
}

func (db *filedb) initializeNodestructure() error {
	db.idmap = make(map[uint64]*nodeentry)

	// id mapping page
	if err := db.pagestore.AddPage(); err != nil {
		return err
	}

	// first node page
	if err := db.pagestore.AddPage(); err != nil {
		return err
	}

	db.nextfree.pagenum = 1
	db.nextfree.offset = 0

	log.Printf("Initialized nodes, nextfree:%v", db.nextfree)

	return nil
}

func (db *filedb) saveNodes() error {
	idpage, err := db.getIdPage()
	if err != nil {
		return err
	}

	db.wrmutex.Lock()
	defer db.wrmutex.Unlock()

	log.Printf("Saving %d nodes, nfree: %v", len(db.idmap), db.nextfree)

	// next free
	binary.LittleEndian.PutUint64(idpage[0:], uint64(db.nextfree.pagenum))
	binary.LittleEndian.PutUint64(idpage[8:], uint64(db.nextfree.offset))

	idlist := make([]uint64, 0, len(db.idmap))
	for id, _ := range db.idmap {
		idlist = append(idlist, id)
	}
	log.Printf("Writing nodes to page: %d, entries left: %d, space left: %d", 0, len(idlist), len(idpage[16:]))
	db.writeNodesPageRec(idpage[16:], idlist)

	return nil
}

func (db *filedb) writeNodesPageRec(space []byte, idlist []uint64) error {
	i := 0
	// next idpage
	binary.LittleEndian.PutUint64(space[0:], 0)
	i += 8

	// node entries
	binary.LittleEndian.PutUint64(space[8:], uint64(len(idlist)))
	i += 8
	for entrynr, id := range idlist {
		if len(space[i:]) < 16 {
			newpage, _ := db.pagestore.GetFreePage()
			newspace, _ := db.pagestore.GetPage(newpage)
			binary.LittleEndian.PutUint64(space[0:], uint64(entrynr-1))
			log.Printf("Continue writing nodes to page: %d, entries left: %d", newpage)

			return db.writeNodesPageRec(newspace, idlist[entrynr:])
		}

		entry := db.idmap[id]
		log.Printf("Persising node(%d) at: %d", id, i)
		binary.LittleEndian.PutUint64(space[i:], uint64(entry.pagenum))
		i += 8
		binary.LittleEndian.PutUint64(space[i:], uint64(entry.offset))
		i += 8

	}

	return nil
}

func (db *filedb) loadNodes() error {
	idpage, err := db.getIdPage()
	if err != nil {
		return err
	}

	db.wrmutex.Lock()
	defer db.wrmutex.Unlock()

	// next free
	db.nextfree.pagenum = int(binary.LittleEndian.Uint64(idpage[0:]))
	db.nextfree.offset = int(binary.LittleEndian.Uint64(idpage[8:]))
	db.idmap = make(nodeentrymap)
	db.loadNodesPageRec(idpage[16:])

	log.Printf("Loaded %d nodes, nextfree:%v", len(db.idmap), db.nextfree)

	return nil
}

func (db *filedb) loadNodesPageRec(space []byte) error {
	i := 0
	// next idpage
	nextpage := int(binary.LittleEndian.Uint64(space[i:]))
	i += 8

	// node entries
	nrentries := binary.LittleEndian.Uint64(space[i:])

	i += 8
	for id := uint64(0); id < nrentries; id++ {
		entry := new(nodeentry)

		entry.pagenum = int(binary.LittleEndian.Uint64(space[i:]))
		i += 8
		entry.offset = int(binary.LittleEndian.Uint64(space[i:]))
		i += 8

		db.idmap[id] = entry
	}
	if nextpage > 0 {
		newspace, _ := db.pagestore.GetPage(nextpage)
		return db.loadNodesPageRec(newspace)
	}
	return nil
}

func (db *filedb) getIdPage() ([]byte, error) {

	idpage, err := db.pagestore.GetPage(0)
	if err != nil {
		return nil, err
	}

	return idpage, nil
}

func (db *filedb) NewNode(labels ...string) Node {

	db.wrmutex.Lock()
	defer db.wrmutex.Unlock()

	wts := atomic.AddUint64(&db.tid, 1)
	id := atomic.AddUint64(&db.nextid, 1)

	pagenum, offset := db.nextfree.pagenum, db.nextfree.offset

	entry := &nodeentry{pagenum, offset, wts, wts, nil}
	db.idmap[id] = entry

	n := &node{db, int(id), labels}

	// saving
	page, err := db.pagestore.GetPage(entry.pagenum)
	if err != nil {
		// aborting transaction
		return nil
	}
	space := page[offset:]

	i := 0
	// labels
	binary.LittleEndian.PutUint16(space[i:], uint16(len(labels)))
	i += 2
	for _, label := range labels {
		binary.LittleEndian.PutUint16(space[i:], uint16(len(label)))
		i += 2
		i += copy(space[i:], label)
	}

	db.nextfree.pagenum, db.nextfree.offset = pagenum, offset+i

	return n
}

func (db *filedb) GetNode(id int) (Node, error) {
	entry, ok := db.idmap[uint64(id)]
	if !ok {
		return nil, errors.New("did not find node for id")
	}

	rts := atomic.AddUint64(&db.tid, 1)
	entry.rts = rts

	n := new(node)
	n.db = db
	n.id = id

	// loading
	page, err := db.pagestore.GetPage(entry.pagenum)
	if err != nil {
		// aborting transaction
		return nil, errors.New("could not get storage for id")
	}
	space := page[entry.offset:]

	i := 0
	// labels
	nlabels := int(binary.LittleEndian.Uint16(space[i:]))
	n.labels = make([]string, 0, nlabels)
	i += 2
	for ; nlabels > 0; nlabels-- {
		lbllen := int(binary.LittleEndian.Uint16(space[i:]))
		i += 2
		n.labels = append(n.labels, string(space[i:lbllen]))
		i += lbllen
	}

	return n, nil
}

func (db *filedb) GetAllNodes() []Node                          { return nil }
func (db *filedb) GetRelation(id int) (Relation, error)         { return nil, nil }
func (db *filedb) GetAllRelations() []Relation                  { return nil }
func (db *filedb) FindPath(start, end Node) Path                { return nil }
func (db *filedb) FindNodeByProperty(prop, value string) []Node { return nil }

func (db *filedb) Close() {
	db.saveNodes()
}
