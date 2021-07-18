// Simple file backed DatabaseService.
package simplefile

import (
	"encoding/binary"
	"os"

	. "github.com/BuJo/goneo/db"
	"github.com/BuJo/goneo/db/mem"
	"github.com/BuJo/goneo/log"
)

type databaseService struct {
	mem      DatabaseService
	filename string
}

// Create a DB instance of a simple file backed graph DB
func NewDb(name string, options map[string][]string) (DatabaseService, error) {
	db, err := mem.NewDb(name, options)
	if err != nil {
		return nil, err
	}
	filename := name

	log.Printf("trying to load from file: %s", filename)
	if _, err := os.Stat(filename); err == nil {
		err = loadDbFile(db, filename)
	}

	return &databaseService{db, filename}, err
}

func (db *databaseService) NewNode(labels ...string) Node { return db.mem.NewNode(labels...) }

func (db *databaseService) GetNode(id int) (Node, error) { return db.mem.GetNode(id) }
func (db *databaseService) GetAllNodes() []Node          { return db.mem.GetAllNodes() }

func (db *databaseService) GetRelation(id int) (Relation, error) { return db.mem.GetRelation(id) }
func (db *databaseService) GetAllRelations() []Relation          { return db.mem.GetAllRelations() }

func (db *databaseService) FindPath(start, end Node) Path { return db.mem.FindPath(start, end) }

func (db *databaseService) FindNodeByProperty(prop, value string) []Node {
	return db.mem.FindNodeByProperty(prop, value)
}

func (db *databaseService) Close() {
	log.Print("Saving to " + db.filename + ".tmp")
	err := saveDbFile(db.filename+".tmp", db.mem)
	if err == nil {
		moveDbFile(db.filename, db.filename+".tmp")
	}

	db.mem.Close()
}

func saveDbFile(tempfile string, db DatabaseService) (err error) {

	file, err := os.OpenFile(tempfile, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0640)
	if err != nil {
		log.Print("Could not open db file")
		return err
	}

	log.Printf("Writing %d nodes", len(db.GetAllNodes()))

	buf := make([]byte, 4, 4)
	for _, n := range db.GetAllNodes() {
		binary.LittleEndian.PutUint32(buf, uint32(n.Id()))
		file.Write(buf)
	}

	err = file.Close()

	return err
}

func moveDbFile(filename, tempfile string) error {
	log.Printf("Moving new db file (%s) over db (%s)", tempfile, filename)
	return os.Rename(tempfile, filename)
}

func loadDbFile(db DatabaseService, filename string) (err error) {
	file, err := os.OpenFile(filename, os.O_RDONLY, 0640)
	if err != nil {
		log.Print("Could not open db file")
		return err
	}
	log.Print("Loading db from file")
	buf := make([]byte, 4, 4)
	i := 0
	for n, e := file.Read(buf); e == nil && i < 50; i++ {
		id := binary.LittleEndian.Uint32(buf)
		log.Printf("Read %d bytes: id: %d", n, id)
		if e != nil {
			err = e
			break
		}
		//id := binary.LittleEndian.Uint32(buf)
		db.NewNode()
	}
	log.Printf("Read %d nodes from file", i)
	return
}
