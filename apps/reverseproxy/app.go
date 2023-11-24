package reverseproxy

import (
	"github.com/vertex-center/vertex/apps/reverseproxy/adapter"
	"github.com/vertex-center/vertex/apps/reverseproxy/core/port"
	"github.com/vertex-center/vertex/apps/reverseproxy/core/service"
	"github.com/vertex-center/vertex/apps/reverseproxy/handler"
	apptypes "github.com/vertex-center/vertex/core/types/app"
	"github.com/vertex-center/vertex/pkg/log"
	"github.com/vertex-center/vertex/pkg/router"
)

var (
	proxyFSAdapter port.ProxyAdapter

	proxyService port.ProxyService
)

type App struct {
	ctx   *apptypes.Context
	proxy *ProxyRouter
}

func NewApp() *App {
	return &App{}
}

func (a *App) Load(ctx *apptypes.Context) {
	a.ctx = ctx
}

func (a *App) Meta() apptypes.Meta {
	return apptypes.Meta{
		ID:          "vx-reverse-proxy",
		Name:        "Vertex Reverse Proxy",
		Description: "Redirect traffic to your containers.",
		Icon:        "router",
	}
}

func (a *App) Initialize(r *router.Group) error {
	proxyFSAdapter = adapter.NewProxyFSAdapter(nil)

	proxyService = service.NewProxyService(proxyFSAdapter)

	a.proxy = NewProxyRouter(proxyService)

	go func() {
		err := a.proxy.Start()
		if err != nil {
			log.Error(err)
		}
	}()

	proxyHandler := handler.NewProxyHandler(proxyService)
	// docapi:v route /app/vx-reverse-proxy/redirects vx_reverse_proxy_get_redirects
	r.GET("/redirects", proxyHandler.GetRedirects)
	// docapi:v route /app/vx-reverse-proxy/redirect vx_reverse_proxy_add_redirect
	r.POST("/redirect", proxyHandler.AddRedirect)
	// docapi:v route /app/vx-reverse-proxy/redirect/{id} vx_reverse_proxy_remove_redirect
	r.DELETE("/redirect/:id", proxyHandler.RemoveRedirect)

	return nil
}

func (a *App) Uninitialize() error {
	return a.proxy.Stop()
}
