package goneo

import (
	"errors"
	. "github.com/BuJo/goneo/db"
	"github.com/BuJo/goneo/db/mem"
	"net/url"
)

type newdb func(db string, options map[string][]string) (DatabaseService, error)

var dbRegistrty map[string]newdb = map[string]newdb{
	"mem": mem.NewDb,
}

// Open a database by URI.
// Example:
// 	OpenDb("mem:testdb")
func OpenDb(dbUri string) (DatabaseService, error) {
	uri, uriErr := url.ParseRequestURI(dbUri)
	if uriErr != nil {
		return nil, uriErr
	}

	dbType := uri.Scheme
	dbInfo := uri.Opaque
	dbOpts := uri.Query()

	dbfunc, foundType := dbRegistrty[dbType]
	if !foundType {
		return nil, errors.New("Did not find DB type for " + dbType)
	}

	return dbfunc(dbInfo, dbOpts)
}
