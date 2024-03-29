package sqldb

import (
	"bufio"
	sqle "github.com/dolthub/go-mysql-server"
	"github.com/dolthub/go-mysql-server/auth"
	"github.com/dolthub/go-mysql-server/memory"
	"github.com/dolthub/go-mysql-server/server"
	"github.com/dolthub/go-mysql-server/sql"
	"github.com/go-gota/gota/dataframe"
	"github.com/go-gota/gota/series"
	"github.com/sirupsen/logrus"
	"github.com/xujiahua/metabase-quick/pkg/util"
	"io"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

type Server struct {
	// default db that new imported table will attach
	DefaultDB *memory.Database
	// internal server
	*server.Server
	mu sync.Mutex
}

func New(addr, defaultDBUser, defaultDBPass, defaultDBName string) (*Server, error) {
	defaultDB := memory.NewDatabase(defaultDBName)

	engine := sqle.NewDefault()
	engine.AddDatabase(defaultDB)

	config := server.Config{
		Protocol: "tcp",
		Address:  addr,
		Auth:     auth.NewNativeSingle(defaultDBUser, defaultDBPass, auth.AllPermissions),
	}

	internalServer, err := server.NewDefaultServer(config, engine)
	if err != nil {
		return nil, err
	}
	return &Server{
		DefaultDB: defaultDB,
		Server:    internalServer,
	}, nil
}

// TODO: support more types, such as timestamp
var typeMapping = map[series.Type]sql.Type{
	series.String: sql.Text,
	series.Int:    sql.Int64,
	series.Float:  sql.Float64,
	series.Bool:   sql.Boolean,
}

func simplifyName(name string) string {
	name = strings.ToLower(name)
	return strings.ReplaceAll(name, " ", "_")
}

func (s *Server) Import(filenames []string, hasHeader bool) error {
	res := make(chan error, len(filenames))
	for _, filename := range filenames {
		logrus.Debugf("loading file %s", filename)
		go func(filename string) {
			res <- s.ImportTable(filename, hasHeader)
		}(filename)
	}

	for i := 0; i < len(filenames); i++ {
		err := <-res
		if err != nil {
			return err
		}
	}

	return nil
}

func jsonArray(r io.Reader) (io.Reader, error) {
	scanner := bufio.NewScanner(r)
	var lines []string
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}

	s := "[" + strings.Join(lines, ",") + "]"
	return strings.NewReader(s), nil
}

func (s *Server) ImportTable(filename string, hasHeader bool) error {
	begin := time.Now()
	defer func() {
		logrus.Infof("load file %s in %v seconds", filename, time.Now().Sub(begin).Seconds())
	}()
	tableName := util.GetFilenameWithExt(filename)
	ext := filepath.Ext(filename)

	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	var dataFrame dataframe.DataFrame
	// infer csv or json from ext, by default, csv
	// schema inferred from dataframe package
	if ext == ".json" {
		// support json array
		r, err := jsonArray(file)
		if err != nil {
			return err
		}
		dataFrame = dataframe.ReadJSON(r)
	} else {
		dataFrame = dataframe.ReadCSV(file, dataframe.HasHeader(hasHeader))
	}
	if dataFrame.Error() != nil {
		return dataFrame.Error()
	}
	var schema sql.Schema
	for _, colName := range dataFrame.Names() {
		schema = append(schema, &sql.Column{
			Name:     simplifyName(colName),
			Type:     typeMapping[dataFrame.Col(colName).Type()],
			Nullable: true,
			Source:   tableName,
		})
	}
	table := memory.NewTable(tableName, schema)

	ctx := sql.NewEmptyContext()
	//inserter := table.Inserter(ctx)
	//defer inserter.Close(ctx)
	for i := 0; i < dataFrame.Nrow(); i++ {
		var row []interface{}
		for _, colName := range dataFrame.Names() {
			row = append(row, dataFrame.Col(colName).Elem(i).Val())
		}
		// TODO: insert is time consuming
		//err = inserter.Insert(ctx, sql.NewRow(row...))
		err = table.Insert(ctx, sql.NewRow(row...))
		if err != nil {
			return err
		}
	}

	s.mu.Lock()
	defer s.mu.Unlock()
	// attach to default db
	// TODO: maybe duplicate table name
	s.DefaultDB.AddTable(tableName, table)

	return nil
}
