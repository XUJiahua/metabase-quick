package sqldb

import (
	sqle "github.com/dolthub/go-mysql-server"
	"github.com/dolthub/go-mysql-server/auth"
	"github.com/dolthub/go-mysql-server/memory"
	"github.com/dolthub/go-mysql-server/server"
	"github.com/dolthub/go-mysql-server/sql"
	"github.com/go-gota/gota/dataframe"
	"github.com/go-gota/gota/series"
	"github.com/xujiahua/csvvisual/pkg/util"
	"os"
)

type Server struct {
	// default db that new imported table will attach
	defaultDB *memory.Database
	// internal server
	*server.Server
}

const defaultDBName = "default"
const defaultDBUser = "root"
const defaultDBPass = ""
const defaultDBAddr = "localhost:3306"

func New(addr string) (*Server, error) {
	defaultDB := memory.NewDatabase(defaultDBName)

	engine := sqle.NewDefault()
	engine.AddDatabase(defaultDB)

	if addr == "" {
		addr = defaultDBAddr
	}
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
		defaultDB: defaultDB,
		Server:    internalServer,
	}, nil
}

var typeMapping = map[series.Type]sql.Type{
	series.String: sql.Text,
	series.Int:    sql.Int64,
	series.Float:  sql.Float64,
	series.Bool:   sql.Boolean,
}

func (s *Server) ImportTable(filename string, hasHeader bool) error {
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
			Name:     colName,
			Type:     typeMapping[dataFrame.Col(colName).Type()],
			Nullable: true,
			Source:   tableName,
		})
	}
	// attach to default db
	table := memory.NewTable(tableName, schema)
	s.defaultDB.AddTable(tableName, table)

	ctx := sql.NewEmptyContext()
	for i := 0; i < dataFrame.Nrow(); i++ {
		var row []interface{}
		for _, colName := range dataFrame.Names() {
			row = append(row, dataFrame.Col(colName).Elem(i).Val())
		}
		err := table.Insert(ctx, sql.NewRow(row...))
		if err != nil {
			return err
		}
	}

	return nil
}
