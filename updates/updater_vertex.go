package updates

import (
	"context"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path"

	"github.com/google/go-github/v50/github"
	"github.com/vertex-center/vertex/core/types"
	"github.com/vertex-center/vertex/pkg/log"
	"github.com/vertex-center/vertex/pkg/storage"
	"github.com/vertex-center/vlog"
)

type VertexUpdater struct {
	about types.About
}

func NewVertexUpdater(about types.About) VertexUpdater {
	return VertexUpdater{
		about: about,
	}
}

func (u VertexUpdater) CurrentVersion() (string, error) {
	return u.about.Version, nil
}

func (u VertexUpdater) Install(tag string) error {
	dir := path.Join(storage.Path, "updates", "vertex")

	log.Info("installing vertex", vlog.String("tag", tag))

	client := github.NewClient(nil)

	release, res, err := client.Repositories.GetReleaseByTag(context.Background(), "vertex-center", "vertex", tag)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	err = storage.DownloadGithubRelease(release, dir)
	if err != nil {
		return err
	}

	err = os.Rename("vertex", "vertex-old")
	if err != nil && !errors.Is(err, fs.ErrNotExist) {
		return fmt.Errorf("failed to rename old executable: %w", err)
	}

	err = os.Rename("vertex-kernel", "vertex-kernel-old")
	if err != nil && !errors.Is(err, fs.ErrNotExist) {
		return fmt.Errorf("failed to rename old executable: %w", err)
	}

	err = os.Rename(path.Join(dir, "vertex"), "vertex")
	if err != nil {
		return err
	}

	err = os.Rename(path.Join(dir, "vertex-kernel"), "vertex-kernel")
	if err != nil {
		return err
	}

	log.Warn("a new Vertex update has been installed. please restart Vertex to apply changes.")

	return nil
}

func (u VertexUpdater) IsInstalled() bool {
	return true
}

func (u VertexUpdater) ID() string {
	return "vertex"
}
