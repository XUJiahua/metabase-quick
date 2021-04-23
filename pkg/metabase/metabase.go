package metabase

import (
	"context"
	"embed"
	"encoding/json"
	"github.com/davecgh/go-spew/spew"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"github.com/xujiahua/metabase-quick/pkg/metabase/model"
	"github.com/xujiahua/metabase-quick/pkg/sqlclient"
	"io/fs"
	"net/http"
	"net/http/httputil"
	"time"
)

//go:embed frontend_client
var staticFiles embed.FS

// Server mock metabase server
type Server struct {
	Databases []*model.Database
	sqlClient *sqlclient.Client
	srv       *http.Server
}

func New(client *sqlclient.Client, addr string, dev bool) (*Server, error) {
	tables, err := client.GetTables()
	if err != nil {
		return nil, err
	}

	database := model.NewDatabase(client.DBName, 1)
	for i, table := range tables {
		database.AddTable(table, i)
	}

	s := &Server{
		Databases: []*model.Database{database},
		sqlClient: client,
	}

	r := mux.NewRouter()
	api := r.PathPrefix("/api").Subrouter()
	// /database?include=tables
	// /database?saved=true
	api.HandleFunc("/database", func(w http.ResponseWriter, r *http.Request) {
		JSON(w, 200, s.Databases)
	})
	api.HandleFunc("/dataset", func(w http.ResponseWriter, r *http.Request) {
		var request model.DataSetRequest
		err := json.NewDecoder(r.Body).Decode(&request)
		if err != nil {
			logrus.Error(err)
			return
		}
		spew.Dump(request.Native.Query)

		response := &model.DataSetResponse{}
		rows, columns, err := s.sqlClient.Query(request.Native.Query)
		if err != nil {
			logrus.Error(err)
			response.Error = err.Error()
			response.Data = &model.Data{
				Rows: [][]interface{}{},
				Cols: []*model.Column{},
			}
		} else {
			response.Data = &model.Data{
				Rows: rows,
				Cols: columns,
			}
		}

		JSON(w, 200, response)
	})

	if dev {
		// reverse proxy to metabase server
		r.PathPrefix("/").HandlerFunc(ReverseProxy("localhost:3000"))
	} else {
		// TODO: refresh on frontend generated url, 404
		staticFiles, err := fs.Sub(staticFiles, "frontend_client")
		if err != nil {
			panic(err)
		}

		r.PathPrefix("/").Handler(http.FileServer(http.FS(staticFiles)))
	}

	s.srv = &http.Server{
		Addr: addr,
		// Good practice to set timeouts to avoid Slowloris attacks.
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
		Handler:      r, // Pass our instance of gorilla/mux in.
	}
	return s, nil
}

func JSON(w http.ResponseWriter, code int, obj interface{}) {
	w.WriteHeader(code)
	_ = json.NewEncoder(w).Encode(obj)
}

func (s Server) Start() error {
	if err := s.srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		return err
	}
	return nil
}

func (s Server) Close() error {
	ctx := context.Background()
	return s.srv.Shutdown(ctx)
}

// ReverseProxy debug mode, direct to metabase backend
func ReverseProxy(target string) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		director := func(req *http.Request) {
			req.URL.Scheme = "http"
			req.URL.Host = target
		}
		proxy := &httputil.ReverseProxy{Director: director}
		proxy.ServeHTTP(writer, request)
	}
}
