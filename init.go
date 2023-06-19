package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/vertex-center/vertex/config"
	"github.com/vertex-center/vertex/pkg/logger"
	"github.com/vertex-center/vertex/pkg/storage"
	"github.com/vertex-center/vertex/router"
	"github.com/vertex-center/vertex/services"
	"github.com/vertex-center/vertex/types"
)

// version, commit and date will be overridden by goreleaser
var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

func main() {
	logger.CreateDefaultLogger()
	defer logger.DefaultLogger.CloseLogFiles()

	logger.Log("Vertex starting...").Print()

	parseArgs()

	err := os.MkdirAll(storage.PathInstances, os.ModePerm)
	if err != nil && !os.IsExist(err) {
		logger.Error(fmt.Errorf("failed to create directory: %v", err)).
			AddKeyValue("message", "failed to create directory").
			AddKeyValue("path", storage.PathInstances).
			Print()

		return
	}

	err = setupDependencies()
	if err != nil {
		logger.Error(fmt.Errorf("failed to setup dependencies: %v", err)).Print()
		return
	}

	r := router.Create(router.About{
		Version: version,
		Commit:  commit,
		Date:    date,
	})
	defer router.Unload()

	logger.Log("Vertex started.").Print()

	err = r.Run(config.Current.Host)
	if err != nil {
		logger.Error(fmt.Errorf("error while starting server: %v", err)).Print()
	}

	logger.Log("Vertex stopped.").Print()
}

func parseArgs() {
	flagVersion := flag.Bool("version", false, "Print vertex version")
	flagV := flag.Bool("v", false, "Print vertex version")
	flagDate := flag.Bool("date", false, "Print the release date")
	flagCommit := flag.Bool("commit", false, "Print the commit hash")

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
}

func setupDependencies() error {
	dependencies := []struct {
		dir string
		dep types.Dependency
	}{
		{storage.PathClient, services.DependencyClient},
		{storage.PathServices, services.DependencyServices},
		{storage.PathPackages, services.DependencyPackages},
	}

	for _, d := range dependencies {
		err := setupDependency(d.dir, d.dep)
		if err != nil {
			logger.Error(err).Print()
			os.Exit(1)
		}
	}
	return nil
}

func setupDependency(dir string, dependency types.Dependency) error {
	err := os.Mkdir(dir, os.ModePerm)
	if os.IsExist(err) {
		// The dependency already exists.
		return nil
	}
	if err != nil {
		return err
	}

	// download
	_, err = dependency.CheckForUpdate()
	if err != nil && err != services.ErrDependencyNotInstalled {
		return err
	}
	return dependency.InstallUpdate()
}
