package goneo

import (
	"fmt"
	. "github.com/BuJo/goneo/db"
	"strings"
	"testing"
)

func TestBasicStartQuery(t *testing.T) {
	db := NewUniverseGenerator().Generate()
	count := len(db.GetAllNodes())

	table, err := Evaluate(db, "start n=node(*) return n as node")
	NewTableTester(t, table, err).HasColumns("node").HasLen(count)
}
func TestUniverse(t *testing.T) {
	db := NewUniverseGenerator().Generate()

	creators := db.FindNodeByProperty("creator", "Joss Whedon")
	if len(creators) != 1 {
		t.Fail()
	}

}

func TestTagged(t *testing.T) {
	db := NewUniverseGenerator().Generate()

	table, err := Evaluate(db, "match (n:Tag)<-[:IS_TAGGED]-(v) return v")
	NewTableTester(t, table, err).HasLen(6).HasColumns("v")
}

func TestTwoReturns(t *testing.T) {
	db := NewUniverseGenerator().Generate()

	table, err := Evaluate(db, "match (n:Tag)<-[r:IS_TAGGED]-(v) return v, n")
	NewTableTester(t, table, err).HasLen(6).HasColumns("v", "n")
}

func TestPropertyRetMatch(t *testing.T) {
	db := NewUniverseGenerator().Generate()

	table, err := Evaluate(db, "match (n:Person {actor: \"Joss Whedon\"}) return n.actor as actor")
	NewTableTester(t, table, err).Has("actor", "Joss Whedon")
}

func TestStartMatch(t *testing.T) {
	db := NewUniverseGenerator().Generate()

	creator, _ := db.GetNode(0)
	fmt.Println(creator, creator.Relations(Both))

	table, err := Evaluate(db, "start joss=node(0) match (joss)-->(o) return o.series")
	NewTableTester(t, table, err).Has("o.series", "Firefly")
}

func TestLongPathMatch(t *testing.T) {
	db := NewUniverseGenerator().Generate()

	table, err := Evaluate(db, "match (e1:Episode)<-[:APPEARED_IN]-(niska {character: \"Adelai Niska\"})-[:APPEARED_IN]->(e2:Episode) return e1, e2")
	if err != nil {
		t.Error(err)
		return
	}
	if table.Len() < 1 {
		t.Error("Evaluation not implemented")
		return
	}
}

func TestMultiMatch(t *testing.T) {
	db := NewUniverseGenerator().Generate()

	table, err := Evaluate(db, "match (e1)-[:ARCS_TO]->(e2), (e1)<-[:APPEARED_IN]-(niska {character: \"Adelai Niska\"})-[:APPEARED_IN]->(e2) return e1.episode, e2.episode")
	NewTableTester(t, table, err).Has("e1.episode", "2")
}

func TestPathVariable(t *testing.T) {
	db := NewUniverseGenerator().Generate()

	table, err := Evaluate(db, "start joss=node(0) match path = (joss)-->(o) return path")
	if err != nil {
		t.Error(err)
		return
	}
	if table.Len() < 1 {
		t.Skip("Saving paths not implemented")
		return
	}
}

func TestFunctionCount(t *testing.T) {
	db := NewUniverseGenerator().Generate()

	table, err := Evaluate(db, "match (e:Episode)-[:ARCS_TO]->(e2) return count(e) as nrArcs")
	NewTableTester(t, table, err).Has("nrArcs", 1)
}

func TestErrorBehaviour(t *testing.T) {
	db := NewUniverseGenerator().Generate()

	table, err := Evaluate(db, "match (e:Episode)-[:ARCS_TO]->(e2) return e.title+e2.title")
	if err == nil || table != nil {
		t.Fatal("Should break with error")
	}

	if !strings.Contains(err.Error(), "bad") {
		t.Fatal("Should contain bad character error")
	}
}

type TableTester struct {
	t          *testing.T
	table      *TabularData
	currentRow int
	err        error
}

func NewTableTester(t *testing.T, table *TabularData, err error) *TableTester {
	tester := new(TableTester)
	tester.t = t
	tester.table = table
	tester.currentRow = -1
	tester.err = err

	if err != nil {
		t.Error(err)
	}

	return tester
}
func (t *TableTester) Has(column string, expected interface{}) *TableTester {
	t.currentRow += 1

	if t.err != nil {
		return t
	}

	if t.table.Len() < t.currentRow+1 {
		t.t.Error("Table not big enough, want: ", t.currentRow+1)
		return t
	}

	if actual := t.table.Get(t.currentRow, column); actual != expected {
		t.t.Error("Bad cell, expected ", expected, " got ", actual)
	}

	return t
}
func (t *TableTester) HasLen(l int) *TableTester {
	if t.err != nil {
		return t
	}

	if t.table.Len() != l {
		t.t.Error("Bad table length, expected ", l, " got ", t.table.Len())
	}
	return t
}
func (t *TableTester) HasColumns(cols ...string) *TableTester {
	if t.err != nil {
		return t
	}

	if len(t.table.Columns()) != len(cols) {
		t.t.Error("Bad number of columns, expected ", len(cols), " got ", len(t.table.Columns()), ": ", t.table.Columns())
	}

	for _, expectedCol := range cols {
		found := false
		for _, actualCol := range t.table.Columns() {
			if expectedCol == actualCol {
				found = true
			}
		}
		if !found {
			t.t.Error("Bad column length, expected ", expectedCol, " in ", t.table.Columns())
		}
	}
	return t
}
