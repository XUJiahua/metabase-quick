package sqldb

import (
	sqle "github.com/dolthub/go-mysql-server"
	"github.com/dolthub/go-mysql-server/auth"
	"github.com/dolthub/go-mysql-server/memory"
	"github.com/dolthub/go-mysql-server/server"
	"github.com/dolthub/go-mysql-server/sql"
	"github.com/go-gota/gota/dataframe"
	"github.com/go-gota/gota/series"
	"github.com/sirupsen/logrus"
	"github.com/xujiahua/metabase-quick/pkg/util"
	"os"
	"strings"
	"time"
)

type Server struct {
	// default db that new imported table will attach
	DefaultDB *memory.Database
	// internal server
	*server.Server
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

func (s *Server) ImportTable(filename string, hasHeader bool) error {
	begin := time.Now()
	defer func() {
		logrus.Infof("load file %s in %v seconds", filename, time.Now().Sub(begin).Seconds())
	}()
	tableName := util.GetFilenameWithExt(filename)

	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	// schema inferred from dataframe package
	dataFrame := dataframe.ReadCSV(file, dataframe.HasHeader(hasHeader))
	var schema sql.Schema
	for _, colName := range dataFrame.Names() {
		schema = append(schema, &sql.Column{
			Name:     simplifyName(colName),
			Type:     typeMapping[dataFrame.Col(colName).Type()],
			Nullable: true,
			Source:   tableName,
		})
	}
	// attach to default db
	table := memory.NewTable(tableName, schema)
	// TODO: maybe duplicate table name
	s.DefaultDB.AddTable(tableName, table)

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

	return nil
}
