package main

import (
	"flag"
	"fmt"
	"os"
	"runtime"

	"github.com/vertex-center/vertex/config"
	"github.com/vertex-center/vertex/migration"
	"github.com/vertex-center/vertex/pkg/log"
	"github.com/vertex-center/vertex/pkg/storage"
	"github.com/vertex-center/vertex/router"
	"github.com/vertex-center/vertex/types"
)

// version, commit and date will be overridden by goreleaser
var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

func main() {
	defer log.Default.Close()

	log.Info("Vertex starting...")

	postMigrationCommands, err := migration.NewMigrationTool(storage.Path).Migrate()
	if err != nil {
		panic(err)
	}

	parseArgs()

	checkNotRoot()

	r := router.NewRouter(types.About{
		Version: version,
		Commit:  commit,
		Date:    date,

		OS:   runtime.GOOS,
		Arch: runtime.GOARCH,
	}, postMigrationCommands)

	// Logs
	url := config.Current.VertexURL()
	fmt.Printf("\n-- Vertex Client :: %s\n\n", url)

	r.Start(fmt.Sprintf(":%s", config.Current.Port))
}

func parseArgs() {
	flagVersion := flag.Bool("version", false, "Print vertex version")
	flagV := flag.Bool("v", false, "Print vertex version")
	flagDate := flag.Bool("date", false, "Print the release date")
	flagCommit := flag.Bool("commit", false, "Print the commit hash")

	var (
		flagHost = flag.String("host", config.Current.Host, "The Vertex access url")

		flagPort           = flag.String("port", config.Current.Port, "The Vertex port")
		flagPortKernel     = flag.String("port-kernel", config.Current.PortKernel, "The Vertex Kernel port")
		flagPortProxy      = flag.String("port-proxy", config.Current.PortProxy, "The Vertex Proxy port")
		flagPortPrometheus = flag.String("port-prometheus", config.Current.PortPrometheus, "The Prometheus port")
	)

	flag.Parse()

	if *flagVersion || *flagV {
		fmt.Println(version)
		os.Exit(0)
	}
	if *flagDate {
		fmt.Println(date)
		os.Exit(0)
	}
	if *flagCommit {
		fmt.Println(commit)
		os.Exit(0)
	}
	config.Current.Host = *flagHost
	config.Current.Port = *flagPort
	config.Current.PortKernel = *flagPortKernel
	config.Current.PortProxy = *flagPortProxy
	config.Current.PortPrometheus = *flagPortPrometheus
}

func checkNotRoot() {
	if os.Getuid() == 0 {
		log.Warn("while vertex-kernel must be run as root, the vertex user should not be root")
	}
}
