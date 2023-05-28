// Package web interface for goneo.
package web

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/BuJo/goneo"
	goneodb "github.com/BuJo/goneo/db"
)

type (
	webHandler struct {
		routes []route

		db goneodb.DatabaseService
	}

	// NodeResponse is a representation of a node
	NodeResponse struct {
		Self                  string
		Property              string
		Properties            string
		Data                  map[string]string
		Labels                string
		OutgoingRelationships string `json:"Outgoing_Relationships"`
		IncomingRelationships string `json:"Incoming_Relationships"`
	}
	// RelationshipResponse is a representation of a relationship
	RelationshipResponse struct {
		Start      string
		Data       map[string]string
		Self       string
		Property   string
		Properties string
		Type       string
		End        string
	}
	// ErrorResponse is a representation of an error condition
	ErrorResponse struct {
		Message    string
		Exception  string
		FullName   string `json:"Fullname"`
		Stacktrace []string
	}
)

func (h *webHandler) baseNodeHandler(w http.ResponseWriter, req *http.Request) {
	http.Error(w, "Not Implemented", http.StatusNotImplemented)
}

// BUG(Jo): createNodeHandler can not create nodes with labels

func (h *webHandler) createNodeHandler(w http.ResponseWriter, req *http.Request) {
	var nodes map[string]string

	decoder := json.NewDecoder(req.Body)
	if decoder.Decode(&nodes) == nil {
		node := h.db.NewNode()
		for key, val := range nodes {
			node.SetProperty(key, val)
		}

		encoder := json.NewEncoder(w)
		_ = encoder.Encode(map[string]string{"status": "created node"})
		w.WriteHeader(http.StatusCreated)
	} else {
		http.Error(w, `{"status":"Bad Request"}`, http.StatusBadRequest)
	}
}
func (h *webHandler) relationshipHandler(w http.ResponseWriter, req *http.Request) {
	http.Error(w, "Not Implemented", http.StatusNotImplemented)
}
func (h *webHandler) nodeHandler(w http.ResponseWriter, req *http.Request) {
	nodeId, err := strconv.Atoi(getField(req, 0))
	if err != nil {
		http.Error(w, `{"status":"Bad ID"}`, http.StatusBadRequest)
		return
	}

	node, err := h.db.GetNode(nodeId)
	if err != nil {
		http.Error(w, `{"status":"Not Found"}`, http.StatusNotFound)
		return
	}

	res := NodeResponse{}
	res.Self = fmt.Sprintf("/db/data/node/%d", nodeId)
	res.OutgoingRelationships = fmt.Sprintf("/db/data/node/%d/direction/out", nodeId)

	res.Data = node.Properties()

	encoder := json.NewEncoder(w)
	_ = encoder.Encode(res)
	w.WriteHeader(http.StatusOK)
}
func (h *webHandler) nodeRelHandler(w http.ResponseWriter, req *http.Request) {
	direction := goneodb.Both
	if d, ok := req.URL.Query()["direction"]; ok {
		direction = goneodb.DirectionFromString(d[0])
	}

	nodeId, err := strconv.Atoi(getField(req, 0))
	if err != nil {
		http.Error(w, `{"status":"Not Found"}`, http.StatusNotFound)
		return
	}

	node, err := h.db.GetNode(nodeId)
	if err != nil {
		http.Error(w, `{"status":"Not Found"}`, http.StatusNotFound)
		return
	}

	res := make([]RelationshipResponse, 0, 5)

	for _, rel := range node.Relations(direction) {
		r := RelationshipResponse{}
		r.Start = "/db/data/node/" + strconv.Itoa(rel.Start().Id())
		r.Self = "/db/data/relationships/" + strconv.Itoa(rel.Id())
		r.End = "/db/data/node/" + strconv.Itoa(rel.End().Id())

		res = append(res, r)
	}

	encoder := json.NewEncoder(w)
	_ = encoder.Encode(res)
	w.WriteHeader(http.StatusOK)
}

// BUG(Jo): re-packaging result table into graph misses rel properties

func (h *webHandler) graphvizHandler(w http.ResponseWriter, req *http.Request) {
	req.ParseForm()

	gocy := req.FormValue("gocy")

	db := h.db
	if gocy != "" {
		// Execute query
		table, err := goneo.Evaluate(db, gocy)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		// Repackage the table from the query in a new in-memory DB for graphs sake
		newdb, _ := goneo.OpenDb("mem:temporary")

		nodemapping := make(map[int]goneodb.Node)
		// Create Nodes
		for i := 0; i < table.Len(); i++ {
			for _, col := range table.Columns() {
				// For each Node element copy the node
				if node, isnode := table.Get(i, col).(goneodb.Node); isnode {
					newnode := newdb.NewNode(node.Labels()...)
					nodemapping[node.Id()] = newnode
					for k, v := range node.Properties() {
						newnode.SetProperty(k, v)
					}
				}
			}
		}

		// Create Relations
		for i := 0; i < table.Len(); i++ {
			for _, col := range table.Columns() {
				if node, isnode := table.Get(i, col).(goneodb.Node); isnode {
					newnode := nodemapping[node.Id()]
					for _, rel := range node.Relations(goneodb.Outgoing) {
						// For each Relation check if the target is mapped as well and
						// create a new edge
						if newtarget, ismapped := nodemapping[rel.End().Id()]; ismapped {
							newnode.RelateTo(newtarget, rel.Type())
							// TODO: copy props
						}
					}
				}
			}
		}

		db = newdb
	}

	dot := goneo.DumpDot(db)

	w.Write([]byte(dot))
	w.WriteHeader(http.StatusOK)
}

func (h *webHandler) gocyTableHandler(w http.ResponseWriter, req *http.Request) {
	req.ParseForm()

	gocy := req.FormValue("gocy")
	if gocy != "" {
		table, err := goneo.Evaluate(h.db, gocy)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		w.Write([]byte(table.String()))
		w.WriteHeader(http.StatusOK)
		return
	}

	w.WriteHeader(http.StatusNotFound)
}

func NewGoneoServer(db goneodb.DatabaseService) http.Handler {
	s := &webHandler{
		db: db,
	}

	s.routes = []route{
		newRoute("GET", "/graphviz", s.graphvizHandler),
		newRoute("POST", "/graphviz", s.graphvizHandler),
		newRoute("POST", "/table", s.gocyTableHandler),
		newRoute("GET", "/db/data/node", s.baseNodeHandler),
		newRoute("POST", "/db/data/node", s.createNodeHandler),
		newRoute("GET", "/db/data/node/([^/]+)", s.nodeHandler),
		newRoute("GET", "/db/data/node/([^/]+)/relationships", s.nodeRelHandler),
		newRoute("GET", "/db/data/relationship/([^/]+)", s.relationshipHandler),
	}

	return s
}
