package gcy

import (
	"fmt"
	"strings"
	"testing"
)

func TestEofWhileParsing(t *testing.T) {
	_, err := Parse("goneo", "start n=node(")
	if err == nil {
		t.Fatal("Parsing should fail")
	}

	errors, iserr := err.(errorList)
	if !iserr {
		t.Fatal("Error should be a list of errors")
	}

	found := false
	for _, e := range errors {
		if strings.Contains(e.Error(), "EOF") {
			found = true
		}
	}
	if !found {
		t.Error("Errors should contain one with EOF")
	}
}

func TestInvalidOperator(t *testing.T) {
	_, err := Parse("goneo", "start n=node(*) return n.a + n")
	if err == nil {
		t.Fatal("Parsing should fail")
	}

	errors, iserr := err.(errorList)
	if !iserr {
		t.Fatal("Error should be a list of errors")
	}

	found := false
	for _, e := range errors {
		fmt.Printf("[invop] " + e.Error())
		if strings.Contains(e.Error(), "+") {
			found = true
		}
	}
	if !found {
		t.Error("Errors should contain one with +")
	}
}

func TestInvalidSquishedOpInRet(t *testing.T) {
	_, err := Parse("goneo", "match (e:Episode)-[:ARCS_TO]->(e2) return e.title+e2.title")
	if err == nil {
		t.Fatal("Parsing should fail")
	}

	errors, iserr := err.(errorList)
	if !iserr {
		t.Fatal("Error should be a list of errors")
	}

	found := false
	for _, e := range errors {
		fmt.Printf("[invop2] " + e.Error())
		if strings.Contains(e.Error(), "bad character") {
			found = true
		}
	}
	if !found {
		t.Error("Errors should contain one with +")
	}
}

func TestParse(t *testing.T) {
	q, err := Parse("goneo", "start n=node(*) return n as node")
	if err != nil {
		t.Fatal(err)
	}

	if len(q.Roots) != 1 {
		t.Error("should have 1 root")
	}

	if len(q.Deletes) != 0 {
		t.Error("should have 0 deletions")
	}

	if len(q.Creates) != 0 {
		t.Error("should have 0 creations")
	}
	if len(q.Returns) != 1 {
		t.Error("should have 1 return")
	}

	if q.Match != nil {
		t.Error("should have no match")
	}

	ret := q.Returns[0]
	if ret.Name != "n" {
		t.Error("return should be named n")

	}
	if ret.Alias != "node" {
		t.Error("return should have node alias")
	}

}

func TestParseCaseSensitivity(t *testing.T) {
	_, err := Parse("goneo", "START n=node(*) RETURN n AS node")
	if err != nil {
		t.Fatal(err)
	}

}
