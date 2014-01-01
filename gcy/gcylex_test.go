package gcy

import "testing"

func TestStartQuery(t *testing.T) {
	testItems(t, "start n=node(*)", itemStart, itemIdentifier, itemEqual, itemIdentifier, itemLParen, itemStar, itemRParen)

	testItems(t, "start n=relation(0)", itemStart, itemIdentifier, itemEqual, itemIdentifier, itemLParen, itemNumber, itemRParen)

	testItems(t, "start n=node(1, 5)", itemStart, itemIdentifier, itemEqual, itemIdentifier, itemLParen, itemNumber, itemComma, itemNumber, itemRParen)

	testItems(t, "start n=relation(0..5)", itemStart, itemIdentifier, itemEqual, itemIdentifier, itemLParen, itemNumber, itemRange, itemNumber, itemRParen)

	testItems(t, "start n=node(..6)", itemStart, itemIdentifier, itemEqual, itemIdentifier, itemLParen, itemRange, itemNumber, itemRParen)

}

func testItems(t *testing.T, query string, items ...itemType) {
	_, ch := lex("t", query)

	for i := range ch {
		if i.typ == itemError {
			t.Error(i.val)
			break
		}
		if i.typ == itemEOF {
			break
		}

		if len(items) == 0 {
			t.Error("not enough items got: ", i)
			break
		}

		expectedTyp := items[0]
		items = items[1:]

		if i.typ != expectedTyp {
			t.Error("unexpected typ: ", i, " expected: ", item{typ: expectedTyp})
		}
	}

	if len(items) != 0 {
		t.Error(len(items), "items left")
	}
}
