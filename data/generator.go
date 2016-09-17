// Package containing test data
package data

import . "github.com/BuJo/goneo/db"

type DatabaseGenerator interface {
	Generate() DatabaseService
}
