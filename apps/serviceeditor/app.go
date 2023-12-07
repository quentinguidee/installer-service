package serviceeditor

import (
	authmeta "github.com/vertex-center/vertex/apps/auth/meta"
	"github.com/vertex-center/vertex/apps/auth/middleware"
	"github.com/vertex-center/vertex/apps/serviceeditor/core/service"
	"github.com/vertex-center/vertex/apps/serviceeditor/handler"
	apptypes "github.com/vertex-center/vertex/core/types/app"
	"github.com/vertex-center/vertex/pkg/router"
)

var Meta = apptypes.Meta{
	ID:          "devtools-service-editor",
	Name:        "Vertex Service Editor",
	Description: "Create services for publishing.",
	Icon:        "frame_source",
	Category:    "devtools",
	DefaultPort: "7510",
	Dependencies: []*apptypes.Meta{
		&authmeta.Meta,
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

	var (
		editorService = service.NewEditorService()
		editorHandler = handler.NewEditorHandler(editorService)
		editor        = r.Group("/editor", "Editor", "Service editor routes", middleware.Authenticated)
	)

	editor.POST("/to-yaml", editorHandler.ToYamlInfo(), editorHandler.ToYaml)

	return nil
}
