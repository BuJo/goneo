package goneo

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"io"
	"log"
	"net/http"
	"strconv"
	"time"
)

type (
	GoneoServer struct {
		binding string

		router *mux.Router

		db *DatabaseService
	}

	ServiceResponse struct {
		Node   string
		Cypher string
	}
	NodeResponse struct {
		Self                   string
		Property               string
		Properties             string
		Data                   map[string]string
		Labels                 string
		Outgoing_Relationships string
		Incoming_Relationships string
	}
	RelationshipResponse struct {
		Start      string
		Data       map[string]string
		Self       string
		Property   string
		Properties string
		Type       string
		End        string
	}
	ErrorResponse struct {
		Message    string
		Exception  string
		Fullname   string
		Stacktrace []string
	}
)

var (
	currentServer *GoneoServer = nil
)

func baseNodeHandler(w http.ResponseWriter, req *http.Request) {

}
func getUrl(urlName string, vars ...string) (string, error) {
	url, err := currentServer.router.Get(urlName).URL(vars...)
	if err != nil {
		fmt.Println("error fetching url:", err)
		return "", err
	}
	return url.String(), nil
}
func createNodeHandler(w http.ResponseWriter, req *http.Request) {
	var data map[string]string

	decoder := json.NewDecoder(req.Body)
	for {
		if err := decoder.Decode(&data); err == io.EOF {
			break
		} else if err != nil {
			log.Fatal(err)
			return
		}

		node := currentServer.db.NewNode()
		for key, val := range data {
			node.SetProperty(key, val)
		}
	}

}
func relationshipHandler(w http.ResponseWriter, req *http.Request) {
}
func nodeHandler(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)

	res := NodeResponse{}

	res.Self, _ = getUrl("Node", "id", vars["id"])
	res.Outgoing_Relationships, _ = getUrl("NodeRelationships", "id", vars["id"], "direction", "out")

	id, _ := strconv.Atoi(vars["id"])
	node, noderr := currentServer.db.GetNode(id)
	if noderr != nil {
		fmt.Println(noderr)
		return
	}

	res.Data = node.Properties()

	b, err := json.MarshalIndent(res, " ", "  ")
	if err != nil {
		fmt.Println("error:", err)
	}
	w.Write(b)
}
func nodeRelHandler(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)

	res := make([]RelationshipResponse, 0, 5)

	nodeId, _ := strconv.Atoi(vars["id"])
	node, _ := currentServer.db.GetNode(nodeId)

	for _, rel := range node.Relations(DirectionFromString(vars["direction"])) {
		r := RelationshipResponse{}
		r.Start, _ = getUrl("Node", "id", strconv.Itoa(rel.Start.id))
		r.Self, _ = getUrl("Relationship", "id", strconv.Itoa(rel.id))
		r.End, _ = getUrl("Node", "id", strconv.Itoa(rel.End.id))

		res = append(res, r)
	}
	b, err := json.MarshalIndent(res, " ", "  ")
	if err != nil {
		fmt.Println("error:", err)
	}
	w.Write(b)
}
func baseHandler(w http.ResponseWriter, req *http.Request) {
	res := ServiceResponse{
		Node: "/db/data/node",
	}

	b, err := json.MarshalIndent(res, " ", "  ")
	if err != nil {
		fmt.Println("error:", err)
	}
	w.Write(b)
}

func NewGoneoServer(db *DatabaseService) *GoneoServer {
	s := new(GoneoServer)

	s.db = db
	s.binding = ":7474"

	currentServer = s

	return s
}

func (s *GoneoServer) Bind(binding string) *GoneoServer {
	s.binding = binding

	return s
}

func (s *GoneoServer) Start() {

	s.router = mux.NewRouter()

	router := s.router.PathPrefix("/db/data").Subrouter()
	router.HandleFunc("/", baseHandler)
	router.HandleFunc("/node", baseNodeHandler).Methods("GET")
	router.HandleFunc("/node", createNodeHandler).Methods("POST")
	router.HandleFunc("/node/{id}", nodeHandler).Name("Node")
	router.HandleFunc("/node/{id}/relationships/{direction}", nodeRelHandler).Name("NodeRelationships")
	router.HandleFunc("/relationship/{id}", relationshipHandler).Name("Relationship")

	srv := &http.Server{
		Addr:           s.binding,
		Handler:        s.router,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	log.Fatal(srv.ListenAndServe())
}
