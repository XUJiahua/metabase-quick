package metabase

import (
	"fmt"
	"github.com/davecgh/go-spew/spew"
	"github.com/gin-gonic/gin"
	"github.com/xujiahua/metabase-quick/pkg/metabase/model"
	"net/http"
	"net/http/httputil"
	"os"
)

type Metadata struct {
	Database string
	Tables   []string
}

// Server mock metabase server
type Server struct {
	Databases []*model.Database
}

func New(metadata *Metadata, verbose bool) (*Server, error) {
	if verbose {
		spew.Dump(metadata)
	}

	// show databases;
	// show tables;
	database := model.NewDatabase(metadata.Database, 1)
	for i, table := range metadata.Tables {
		database.AddTable(table, i)
	}

	return &Server{Databases: []*model.Database{database}}, nil
}

func (s Server) Start() error {
	router := gin.Default()
	apiGroup := router.Group("/api")
	//const contentType = "application/json;charset=utf-8"
	{
		// /database?include=tables
		// /database?saved=true
		apiGroup.GET("/database", func(c *gin.Context) {
			c.JSON(200, s.Databases)
		})
	}

	// reverse proxy to metabase server
	router.NoRoute(ReverseProxy())
	port, _ := os.LookupEnv("PORT")
	return router.Run(fmt.Sprintf(":%s", port))
}

// debug mode, direct to metabase backend
func ReverseProxy() gin.HandlerFunc {
	target := "localhost:3000"

	return func(c *gin.Context) {
		director := func(req *http.Request) {
			req.URL.Scheme = "http"
			req.URL.Host = target
		}
		proxy := &httputil.ReverseProxy{Director: director}
		proxy.ServeHTTP(c.Writer, c.Request)
	}
}
