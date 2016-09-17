// Web interface for goneo.
package web

import (
	"github.com/BuJo/goneo"
	goneodb "github.com/BuJo/goneo/db"
	"github.com/BuJo/goneo/db/mem"
	"github.com/gin-gonic/gin"
	Graphite "github.com/marpaia/graphite-golang"
	"log"
	"net/http"
	"strconv"
	"time"
)

type GoneoServer interface {
	Bind(binding string) GoneoServer // default: :7474
	EnableGraphite(host string, port int, prefix string, notlogged ...string) GoneoServer
	Start()
}

type (
	goneoServer struct {
		router  *gin.Engine
		binding string
	}

	// JSON representation of a node
	NodeResponse struct {
		Self                   string
		Property               string
		Properties             string
		Data                   map[string]string
		Labels                 string
		Outgoing_Relationships string
		Incoming_Relationships string
	}
	// JSON representation of a relationship
	RelationshipResponse struct {
		Start      string
		Data       map[string]string
		Self       string
		Property   string
		Properties string
		Type       string
		End        string
	}
	// JSON representation of an error condition
	ErrorResponse struct {
		Message    string
		Exception  string
		Fullname   string
		Stacktrace []string
	}
)

func baseNodeHandler(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{"status": "not implemented"})
}

func createNodeHandler(c *gin.Context) {
	db, _ := c.MustGet("db").(goneodb.DatabaseService)

	var json map[string]string
	if c.BindJSON(&json) == nil {
		node := db.NewNode()
		for key, val := range json {
			node.SetProperty(key, val)
		}
		c.JSON(http.StatusCreated, gin.H{"status": "created node"})
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"status": "bad request"})
	}
}
func relationshipHandler(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{"status": "not implemented"})
}
func nodeHandler(c *gin.Context) {
	db, _ := c.MustGet("db").(goneodb.DatabaseService)

	nodeId, idErr := strconv.Atoi(c.Param("id"))
	if idErr != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "bad id"})
		return
	}

	node, nodeerr := db.GetNode(nodeId)
	if nodeerr != nil {
		c.JSON(http.StatusNotFound, gin.H{"status": "not found"})
		return
	}

	res := NodeResponse{}
	res.Self = "/db/data/node/" + c.Param("id")
	res.Outgoing_Relationships = "/db/data/node/" + c.Param("id") + "/direction/out"

	res.Data = node.Properties()

	c.JSON(http.StatusOK, res)
}
func nodeRelHandler(c *gin.Context) {
	db, _ := c.MustGet("db").(goneodb.DatabaseService)

	direction := goneodb.DirectionFromString(c.Param("direction"))
	nodeId, idErr := strconv.Atoi(c.Param("id"))
	if idErr != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "bad id"})
		return
	}

	node, nodeerr := db.GetNode(nodeId)
	if nodeerr != nil {
		c.JSON(http.StatusNotFound, gin.H{"status": "not found"})
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

	c.JSON(http.StatusOK, res)
}

// BUG(Jo): re-packaging result table into graph misses rel properties

func graphvizHandler(c *gin.Context) {
	db, _ := c.MustGet("db").(goneodb.DatabaseService)

	gocy := c.PostForm("gocy")
	if gocy != "" {
		// Execute query
		table, err := goneo.Evaluate(db, gocy)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"status": err.Error()})
			return
		}

		// Repackage the table from the query in a new in-memory DB for graphs sake
		newdb := mem.NewDb()

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
					newnode, _ := nodemapping[node.Id()]
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

	c.String(http.StatusOK, dot)
}

func gocyTableHandler(c *gin.Context) {
	db, _ := c.MustGet("db").(goneodb.DatabaseService)

	gocy := c.PostForm("gocy")
	if gocy != "" {
		table, err := goneo.Evaluate(db, gocy)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"status": err.Error()})
			return
		}

		c.String(http.StatusOK, table.String())
	}

	c.String(http.StatusOK, "")
}

func NewGoneoServer(db goneodb.DatabaseService) GoneoServer {
	s := new(goneoServer)

	s.router = gin.Default()
	s.binding = ":7474"

	s.router.Use(func(c *gin.Context) { c.Set("db", db) })

	return s
}

func (s *goneoServer) Bind(binding string) GoneoServer {
	s.binding = binding

	return s
}

// Enable sending Graphite metrics
// Output:
// 	<prefix>.http.time
func (s *goneoServer) EnableGraphite(host string, port int, prefix string, notlogged ...string) GoneoServer {
	graphite := &Graphite.Graphite{Host: host, Port: port, Protocol: "udp", Prefix: prefix}
	err := graphite.Connect()

	if err != nil {
		log.Print("failed to connect to graphite", err.Error())
		graphite = Graphite.NewGraphiteNop(host, port)
	}

	var skip map[string]struct{}

	if length := len(notlogged); length > 0 {
		skip = make(map[string]struct{}, length)

		for _, path := range notlogged {
			skip[path] = struct{}{}
		}
	}

	s.router.Use(func(c *gin.Context) {
		// Start timer
		start := time.Now()
		path := c.Request.URL.Path

		log.Print("starting timer for graphite")

		// Process request
		c.Next()

		// Log only when path is not being skipped
		if _, ok := skip[path]; !ok {
			// Stop timer
			log.Print("stopping timer for graphite")

			graphite.SimpleSend("http.time", strconv.FormatFloat(time.Now().Sub(start).Seconds(), 'f', -1, 32))
		}
	})

	return s
}

func (s *goneoServer) Start() {

	datarouter := s.router.Group("/db/data")

	s.router.GET("/graphviz", graphvizHandler)
	s.router.POST("/graphviz", graphvizHandler)
	s.router.POST("/table", gocyTableHandler)

	noderouter := datarouter.Group("/node")
	{
		noderouter.GET("/", baseNodeHandler)
		noderouter.POST("/", createNodeHandler)
		noderouter.GET("/:id", nodeHandler)
		noderouter.GET("/:id/relationships/:direction", nodeRelHandler)
	}
	relrouter := datarouter.Group("/relationship")
	{
		relrouter.GET("/:id", relationshipHandler)
	}

	s.router.Run(s.binding)
}
