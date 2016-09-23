package file

import (
	"flag"
	"log"
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	flag.Parse()

	log.SetFlags(log.Ldate | log.Lmicroseconds | log.Lshortfile | log.LUTC)

	os.Exit(m.Run())
}
