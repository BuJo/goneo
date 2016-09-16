package web

import (
	"github.com/gin-gonic/gin"
	goneodb "goneo/db"
	"net/http"
	"strconv"
)

type (
	GoneoServer struct {
		binding string
		db      goneodb.DatabaseService
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

func baseNodeHandler(c *gin.Context) {

}

func createNodeHandler(c *gin.Context) {

	var json map[string]string
	if c.BindJSON(&json) == nil {
		node := currentServer.db.NewNode()
		for key, val := range json {
			node.SetProperty(key, val)
		}
		c.JSON(http.StatusCreated, gin.H{"status": "created node"})
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"status": "bad request"})
	}
}
func relationshipHandler(c *gin.Context) {
}
func nodeHandler(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))

	node, nodeerr := currentServer.db.GetNode(id)
	if nodeerr != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": nodeerr})
		return
	}

	res := NodeResponse{}
	res.Self = "/db/data/node/" + c.Param("id")
	res.Outgoing_Relationships = "/db/data/node/" + c.Param("id") + "/direction/out"

	res.Data = node.Properties()

	c.JSON(http.StatusOK, res)
}
func nodeRelHandler(c *gin.Context) {
	direction := goneodb.DirectionFromString(c.Param("direction"))
	nodeId, _ := strconv.Atoi(c.Param("id"))
	node, _ := currentServer.db.GetNode(nodeId)

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

func NewGoneoServer(db goneodb.DatabaseService) *GoneoServer {
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

	router := gin.Default()

	noderouter := router.Group("/db/data/node")
	{
		noderouter.GET("/", baseNodeHandler)
		noderouter.POST("/", createNodeHandler)
		noderouter.GET("/:id", nodeHandler)
		noderouter.GET("/:id/relationships/:direction", nodeRelHandler)
	}
	relrouter := router.Group("/db/data/relationship")
	{
		relrouter.GET("/relationship/:id", relationshipHandler)
	}

	router.Run(s.binding)
}
