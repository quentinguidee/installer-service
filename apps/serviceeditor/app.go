package serviceeditor

import (
	"github.com/vertex-center/vertex/apps/auth/middleware"
	"github.com/vertex-center/vertex/apps/serviceeditor/core/service"
	"github.com/vertex-center/vertex/apps/serviceeditor/handler"
	apptypes "github.com/vertex-center/vertex/core/types/app"
	"github.com/vertex-center/vertex/pkg/router"
)

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
	return apptypes.Meta{
		ID:          "devtools-service-editor",
		Name:        "Vertex Service Editor",
		Description: "Create services for publishing.",
		Icon:        "frame_source",
		Category:    "devtools",
	}
}

func (a *App) Initialize(r *router.Group) error {
	editorService := service.NewEditorService()

	editorHandler := handler.NewEditorHandler(editorService)
	editor := r.Group("/editor", middleware.Authenticated)
	// docapi:v route /app/devtools-service-editor/editor/to-yaml vx_devtools_service_editor_to_yaml
	editor.POST("/to-yaml", editorHandler.ToYaml)

	return nil
}