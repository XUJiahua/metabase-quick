package metabase

import (
	"encoding/json"
	"fmt"
	"github.com/davecgh/go-spew/spew"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"github.com/xujiahua/metabase-quick/pkg/metabase/model"
	"github.com/xujiahua/metabase-quick/pkg/sqlclient"
	"net/http"
	"net/http/httputil"
	"os"
	"time"
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

func JSON(w http.ResponseWriter, code int, obj interface{}) {
	w.WriteHeader(code)
	_ = json.NewEncoder(w).Encode(obj)
}

func (s Server) Start(dev bool) error {
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

		rows, columns, err := s.sqlClient.Query(request.Native.Query)
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
		JSON(w, 200, response)
	})

	if dev {
		// reverse proxy to metabase server
		r.PathPrefix("/").HandlerFunc(ReverseProxy("localhost:3000"))
	} else {
		// TODO: refresh on frontend generated url, 404
		r.PathPrefix("/").Handler(http.StripPrefix("/", http.FileServer(http.Dir("./frontend_client"))))
	}

	port, _ := os.LookupEnv("PORT")
	srv := &http.Server{
		Addr: fmt.Sprintf(":%s", port),
		// Good practice to set timeouts to avoid Slowloris attacks.
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
		Handler:      r, // Pass our instance of gorilla/mux in.
	}

	return srv.ListenAndServe()
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
