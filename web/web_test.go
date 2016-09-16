package web

import (
	"github.com/BuJo/goneo/db/mem"
)

func ExampleNewGoneoServer() {
	db := mem.NewDb()

	NewGoneoServer(db).Bind(":7878").Start()
}
