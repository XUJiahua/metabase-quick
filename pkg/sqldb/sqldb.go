package sqldb

import (
	sqle "github.com/dolthub/go-mysql-server"
	"github.com/dolthub/go-mysql-server/auth"
	"github.com/dolthub/go-mysql-server/memory"
	"github.com/dolthub/go-mysql-server/server"
	"github.com/dolthub/go-mysql-server/sql"
	"time"
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

func New(addr string) (*Server, error) {
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
		defaultDB: defaultDB,
		Server:    internalServer,
	}, nil
}

func (s *Server) ImportTable() {
	const (
		// TODO: table name is filename
		tableName = "mytable"
	)
	// TODO: schema inferred from?
	table := memory.NewTable(tableName, sql.Schema{
		{Name: "name", Type: sql.Text, Nullable: false, Source: tableName},
		{Name: "email", Type: sql.Text, Nullable: false, Source: tableName},
		{Name: "phone_numbers", Type: sql.JSON, Nullable: false, Source: tableName},
		{Name: "created_at", Type: sql.Timestamp, Nullable: false, Source: tableName},
	})

	s.defaultDB.AddTable(tableName, table)

	ctx := sql.NewEmptyContext()
	rows := []sql.Row{
		sql.NewRow("John Doe", "john@doe.com", []string{"555-555-555"}, time.Now()),
		sql.NewRow("John Doe", "johnalt@doe.com", []string{}, time.Now()),
		sql.NewRow("Jane Doe", "jane@doe.com", []string{}, time.Now()),
		sql.NewRow("Evil Bob", "evilbob@gmail.com", []string{"555-666-555", "666-666-666"}, time.Now()),
	}
	for _, row := range rows {
		err := table.Insert(ctx, row)
		if err != nil {
			return
		}
	}
}
