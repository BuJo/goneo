package web

import (
	"github.com/gin-gonic/gin"
	"goneo"
	goneodb "goneo/db"
	"net/http"
	"strconv"
)

type (
	GoneoServer struct {
		router  *gin.Engine
		binding string
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

func graphvizHandler(c *gin.Context) {
	db, _ := c.MustGet("db").(goneodb.DatabaseService)

	dot := goneo.DumpDot(db)

	c.String(http.StatusOK, dot)
}

func NewGoneoServer(db goneodb.DatabaseService) *GoneoServer {
	s := new(GoneoServer)

	s.router = gin.Default()
	s.binding = ":7474"

	s.router.Use(func(c *gin.Context) { c.Set("db", db) })

	datarouter := s.router.Group("/db/data")

	datarouter.GET("/graphviz", graphvizHandler)

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

	return s
}

func (s *GoneoServer) Bind(binding string) *GoneoServer {
	s.binding = binding

	return s
}

func (s *GoneoServer) Start() {

	s.router.Run(s.binding)
}
