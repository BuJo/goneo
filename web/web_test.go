package web

import (
	"goneo/db/mem"
)

func ExampleNewGoneoServer() {
	db := mem.NewDb()

	NewGoneoServer(db).Bind(":7878").Start()
}
