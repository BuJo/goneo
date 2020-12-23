package web

import (
	"github.com/BuJo/goneo"
)

func ExampleNewGoneoServer() {
	db, _ := goneo.OpenDb("mem:test")

	NewGoneoServer(db).Bind(":7878").Start()
}
