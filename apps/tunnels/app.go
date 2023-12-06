package tunnels

import (
	authmeta "github.com/vertex-center/vertex/apps/auth/meta"
	"github.com/vertex-center/vertex/apps/auth/middleware"
	containersmeta "github.com/vertex-center/vertex/apps/containers/meta"
	"github.com/vertex-center/vertex/apps/tunnels/handler"
	apptypes "github.com/vertex-center/vertex/core/types/app"
	"github.com/vertex-center/vertex/pkg/router"
)

// docapi:tunnels title Vertex Tunnels
// docapi:tunnels description A tunnel manager.
// docapi:tunnels version 0.0.0
// docapi:tunnels filename tunnels

// docapi:tunnels url http://{ip}:{port-kernel}/api
// docapi:tunnels urlvar ip localhost The IP address of the server.
// docapi:tunnels urlvar port-kernel 7514 The port of the server.

var Meta = apptypes.Meta{
	ID:          "tunnels",
	Name:        "Vertex Tunnels",
	Description: "Create and manage tunnels.",
	Icon:        "subway",
	DefaultPort: "7514",
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

	providerHandler := handler.NewProviderHandler()
	// docapi:tunnels route /app/tunnels/provider/{provider}/install vx_tunnels_install_provider
	r.POST("/provider/:provider/install", middleware.Authenticated, providerHandler.Install)

	return nil
}
