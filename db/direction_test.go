package db

import (
	"fmt"
	"testing"
)

func ExampleDirectionFromString() {
	fmt.Println(DirectionFromString("out") == Outgoing)
	// Output: true
}

func testDir(t *testing.T, expected Direction, given string) {
	if DirectionFromString(given) != expected {
		t.Error("Should be " + expected.String() + ", was : " + DirectionFromString(given).String())
	}
}

func TestOutgoing(t *testing.T) {
	testDir(t, Outgoing, "out")
	testDir(t, Outgoing, "outgoing")
}

func TestIncoming(t *testing.T) {
	testDir(t, Incoming, "in")
	testDir(t, Incoming, "incoming")
}

func TestBoth(t *testing.T) {
	testDir(t, Both, "foo")
	testDir(t, Both, "bar")
	testDir(t, Both, "both")
}
