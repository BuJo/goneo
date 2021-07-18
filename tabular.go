package goneo

import (
	"bytes"
	"fmt"
	"text/tabwriter"

	"github.com/BuJo/goneo/log"
)

// Describes information in a tabular format.
type TabularData struct {
	line []map[string]interface{}
}

func (t *TabularData) String() string {
	if len(t.line) == 0 {
		return ""
	}
	b := new(bytes.Buffer)

	w := new(tabwriter.Writer)

	// Format in tab-separated columns with a tab stop of 8.
	w.Init(b, 0, 8, 0, '\t', 0)

	headers := make([]string, 0)

	for header := range t.line[0] {
		headers = append(headers, header)
	}

	for _, header := range headers {
		fmt.Fprint(w, header+"\t")
	}
	fmt.Fprintln(w, "")

	for _, line := range t.line {
		for _, header := range headers {
			if _, ok := line[header]; ok && line[header] != nil {
				fmt.Fprintf(w, "%s\t", line[header])
			}
		}
		fmt.Fprintln(w, "")
	}
	w.Flush()

	return b.String()
}

// Return count of lines
func (t *TabularData) Len() int {
	return len(t.line)
}

// Return column names
func (t *TabularData) Columns() []string {
	cols := make([]string, 0)
	for k := range t.line[0] {
		cols = append(cols, k)
	}
	return cols
}

// Return a single field
func (t *TabularData) Get(row int, column string) interface{} {
	return t.line[row][column]
}

// BUG(Jo): TabularData does not handle tables of different sizes.

// Merge two tables.
func (t *TabularData) Merge(t2 *TabularData) *TabularData {
	merged := new(TabularData)

	if t.Len() > 0 && t.Len() != t2.Len() {
		// TODO: product? unsure how to handle...
		log.Println("TODO: ignored differently sized tables: ", t.Len(), t2.Len())
	} else if t.Len() == 0 {
		merged.line = make([]map[string]interface{}, t2.Len(), t2.Len())

		for i := 0; i < t2.Len(); i += 1 {
			merged.line[i] = make(map[string]interface{})
			for k, v := range t2.line[i] {
				merged.line[i][k] = v
			}
		}
	} else {
		merged.line = make([]map[string]interface{}, t2.Len(), t2.Len())

		for i := 0; i < t2.Len(); i += 1 {
			merged.line[i] = make(map[string]interface{})
			for k, v := range t.line[i] {
				merged.line[i][k] = v
			}
			for k, v := range t2.line[i] {
				merged.line[i][k] = v
			}
		}
	}
	//log.Print("merged tables: ", merged)
	return merged
}
