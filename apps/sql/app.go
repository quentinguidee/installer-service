package sql

import (
	authmeta "github.com/vertex-center/vertex/apps/auth/meta"
	"github.com/vertex-center/vertex/apps/auth/middleware"
	containersmeta "github.com/vertex-center/vertex/apps/containers/meta"
	"github.com/vertex-center/vertex/apps/sql/core/port"
	"github.com/vertex-center/vertex/apps/sql/core/service"
	"github.com/vertex-center/vertex/apps/sql/handler"
	apptypes "github.com/vertex-center/vertex/core/types/app"
	"github.com/vertex-center/vertex/pkg/router"
)

// docapi:sql title Vertex SQL
// docapi:sql description A SQL database manager.
// docapi:sql version 0.0.0
// docapi:sql filename sql

// docapi:sql url http://{ip}:{port-kernel}/api
// docapi:sql urlvar ip localhost The IP address of the server.
// docapi:sql urlvar port-kernel 7512 The port of the server.

var (
	sqlService port.SqlService
)

var Meta = apptypes.Meta{
	ID:          "sql",
	Name:        "Vertex SQL",
	Description: "Create and manage SQL databases.",
	Icon:        "database",
	DefaultPort: "7512",
	Dependencies: []*apptypes.Meta{
		&authmeta.Meta,
		&containersmeta.Meta,
	},
}

type App struct {
	ctx *apptypes.Context
}

func NewApp() *App {
	return &App{}
}

func (a *App) Load(ctx *apptypes.Context) {
	a.ctx = ctx
}

func (a *App) Meta() apptypes.Meta {
	return Meta
}

func (a *App) Initialize(r *router.Group) error {
	r.Use(middleware.ReadAuth)

	sqlService = service.New(a.ctx)

	dbmsHandler := handler.NewDBMSHandler(sqlService)
	// docapi:sql route /container/{container_uuid} vx_sql_get_dbms
	r.GET("/container/:container_uuid", middleware.Authenticated, dbmsHandler.Get)
	// docapi:sql route /dbms/{dbms}/install vx_sql_install_dbms
	r.POST("/dbms/:dbms/install", middleware.Authenticated, dbmsHandler.Install)

	return nil
}
