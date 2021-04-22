package metabase

import (
	"fmt"
	"github.com/davecgh/go-spew/spew"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/xujiahua/metabase-quick/pkg/metabase/model"
	"github.com/xujiahua/metabase-quick/pkg/sqlclient"
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
	sqlClient *sqlclient.Client
}

func New(metadata *Metadata, verbose bool, client *sqlclient.Client) (*Server, error) {
	if verbose {
		spew.Dump(metadata)
	}

	// show databases;
	// show tables;
	database := model.NewDatabase(metadata.Database, 1)
	for i, table := range metadata.Tables {
		database.AddTable(table, i)
	}

	return &Server{
		Databases: []*model.Database{database},
		sqlClient: client,
	}, nil
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

		apiGroup.POST("/dataset", func(c *gin.Context) {
			var request model.DataSetRequest
			err := c.ShouldBindJSON(&request)
			if err != nil {
				logrus.Error(err)
				return
			}
			spew.Dump(request.Native.Query)

			rows, columns, err := s.sqlClient.RowsAndCols(request.Native.Query)
			if err != nil {
				logrus.Error(err)
				return
			}

			response := &model.DataSetResponse{
				Data: &model.Data{
					Rows: rows,
					Cols: columns,
				},
			}
			c.JSON(200, response)
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
