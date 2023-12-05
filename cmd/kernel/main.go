package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"os/user"
	"strconv"
	"time"

	"github.com/vertex-center/vertex/apps/admin"
	"github.com/vertex-center/vertex/apps/auth"
	"github.com/vertex-center/vertex/apps/containers"
	"github.com/vertex-center/vertex/apps/monitoring"
	"github.com/vertex-center/vertex/apps/reverseproxy"
	"github.com/vertex-center/vertex/apps/serviceeditor"
	"github.com/vertex-center/vertex/apps/sql"
	"github.com/vertex-center/vertex/apps/tunnels"
	"github.com/vertex-center/vertex/config"
	"github.com/vertex-center/vertex/core/service"
	"github.com/vertex-center/vertex/core/types"
	"github.com/vertex-center/vertex/core/types/app"
	"github.com/vertex-center/vertex/core/types/server"
	"github.com/vertex-center/vertex/pkg/log"
	"github.com/vertex-center/vertex/pkg/netcap"
	"github.com/vertex-center/vlog"
)

var (
	srv *server.Server
	ctx *types.VertexContext
)

// docapi:k title Vertex Kernel
// docapi:k description A platform to manage your self-hosted server.
// docapi:k version 0.0.0
// docapi:k filename kernel

// docapi:k url http://{ip}:{port-kernel}/api
// docapi:k urlvar ip localhost The IP address of the kernel.
// docapi:k urlvar port-kernel 6131 The port of the server.

func main() {
	ensureRoot()
	parseArgs()

	// If go.mod is there, build vertex first.
	_, err := os.Stat("go.mod")
	if err == nil {
		log.Info("init.go found. Building vertex...")
		buildVertex()
	}

	err = netcap.AllowPortsManagement("vertex")
	if err != nil {
		log.Error(err)
	}

	ctx = types.NewVertexContext(types.About{}, true)
	addr := fmt.Sprintf(":%s", config.KernelCurrent.Ports["VERTEX_KERNEL"])

	srv = server.New("kernel", addr, ctx)
	initServices()

	ctx.DispatchEvent(types.EventServerLoad{})
	ctx.DispatchEvent(types.EventServerStart{})

	exitKernelChan := srv.StartAsync()
	exitVertexChan := make(chan error)

	var vertex *exec.Cmd
	go func() {
		defer close(exitVertexChan)

		var err error
		vertex, err = runVertex()
		if err != nil {
			exitVertexChan <- err
			return
		}
		exitVertexChan <- vertex.Wait()
	}()

	for {
		select {
		case err := <-exitKernelChan:
			if err != nil {
				log.Error(err)
			}
			if vertex != nil && vertex.Process != nil {
				_ = vertex.Process.Signal(os.Interrupt)
				_, _ = vertex.Process.Wait()
			}
		case err := <-exitVertexChan:
			if err != nil {
				log.Error(err)
			}
		}
	}
}

func ensureRoot() {
	if os.Getuid() != 0 {
		log.Warn("vertex-kernel must be run as root to work properly")
	}
}

func parseArgs() {
	var (
		flagUsername = flag.String("user", "", "username of the unprivileged user")
		flagUID      = flag.Uint("uid", 0, "uid of the unprivileged user")
		flagGID      = flag.Uint("gid", 0, "gid of the unprivileged user")
		flagHost     = flag.String("host", config.Current.Host, "The Vertex access url")
	)

	flag.Parse()

	config.KernelCurrent.Host = *flagHost

	if *flagUsername == "" {
		*flagUsername = os.Getenv("USER")
		log.Warn("no username specified; trying to retrieve username from env", vlog.String("user", *flagUsername))
	}

	if *flagUsername != "" {
		u, err := user.Lookup(*flagUsername)
		if err != nil {
			log.Error(err)
			os.Exit(1)
		}

		uid, err := strconv.ParseInt(u.Uid, 10, 32)
		if err != nil {
			log.Error(err)
			os.Exit(1)
		}

		gid, err := strconv.ParseInt(u.Gid, 10, 32)
		if err != nil {
			log.Error(err)
			os.Exit(1)
		}

		config.KernelCurrent.Uid = uint32(uid)
		config.KernelCurrent.Gid = uint32(gid)
		return
	}

	config.KernelCurrent.Uid = uint32(*flagUID)
	config.KernelCurrent.Gid = uint32(*flagGID)
}

func buildVertex() {
	log.Info("Building vertex")

	start := time.Now()
	cmd := exec.Command("go", "build", "-o", "vertex", "cmd/main/main.go")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	end := time.Now()

	if err != nil {
		log.Error(err)
		os.Exit(1)
	}

	log.Info("Build completed in " + end.Sub(start).String())
}

func initServices() {
	service.NewAppsService(ctx, true, []app.Interface{
		admin.NewApp(),
		auth.NewApp(),
		sql.NewApp(),
		tunnels.NewApp(),
		monitoring.NewApp(),
		containers.NewApp(),
		reverseproxy.NewApp(),
		serviceeditor.NewApp(),
	})
}
