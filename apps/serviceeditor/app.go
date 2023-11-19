package serviceeditor

import (
	"github.com/vertex-center/vertex/apps/serviceeditor/core/service"
	"github.com/vertex-center/vertex/apps/serviceeditor/handler"
	apptypes "github.com/vertex-center/vertex/core/types/app"
	"github.com/vertex-center/vertex/pkg/router"
)

const (
	AppRoute = "/vx-devtools-service-editor"
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
		ID:          "vx-devtools-service-editor",
		Name:        "Vertex Service Editor",
		Description: "Create services for publishing.",
		Icon:        "frame_source",
		Category:    "devtools",
	}
}

func (a *App) Initialize(r *router.Group) error {
	editorService := service.NewEditorService()

	editorHandler := handler.NewEditorHandler(editorService)
	editor := r.Group("/editor")
	editor.POST("/to-yaml", editorHandler.ToYaml)

	return nil
}