package adapter

import (
	"context"
	"database/sql"
	"errors"
	"strings"

	"github.com/vertex-center/vertex/apps/containers/core/port"
	"github.com/vertex-center/vertex/apps/containers/core/types"
	"github.com/vertex-center/vertex/common/storage"
	"github.com/vertex-center/vertex/common/uuid"
)

type containerDBAdapter struct {
	db storage.DB
}

func NewContainerDBAdapter(db storage.DB) port.ContainerAdapter {
	return &containerDBAdapter{db}
}

func (a *containerDBAdapter) GetContainer(ctx context.Context, id uuid.UUID) (*types.Container, error) {
	var container types.Container
	err := a.db.Get(&container, `
		SELECT * FROM containers
		WHERE id = $1
	`, id)
	if errors.Is(err, sql.ErrNoRows) {
		return &container, nil
	}
	return &container, err
}

func (a *containerDBAdapter) GetContainers(ctx context.Context) (types.Containers, error) {
	return a.GetContainersWithFilters(ctx, types.ContainerFilters{})
}

func (a *containerDBAdapter) GetContainersWithFilters(ctx context.Context, filters types.ContainerFilters) (types.Containers, error) {
	var containers types.Containers
	query := `SELECT containers.* FROM containers`
	var args []interface{}
	if filters.Tags != nil {
		query += ` INNER JOIN container_tags ct on containers.id = ct.container_id`
		query += ` INNER JOIN tags t on ct.tag_id = t.id`
		query += ` WHERE t.name IN (?)`
		tags := strings.Join(*filters.Tags, ", ")
		args = append(args, tags)
	}
	err := a.db.Select(&containers, query, args...)
	if errors.Is(err, sql.ErrNoRows) {
		return containers, nil
	}
	return containers, err
}

func (a *containerDBAdapter) CreateContainer(ctx context.Context, c types.Container) error {
	_, err := a.db.NamedExec(`
		INSERT INTO containers (id, service_id, user_id, image, image_tag, status, launch_on_startup, name, description, color, icon, command)
		VALUES (:id, :service_id, :user_id, :image, :image_tag, :status, :launch_on_startup, :name, :description, :color, :icon, :command)
	`, c)
	return err
}

func (a *containerDBAdapter) UpdateContainer(ctx context.Context, c types.Container) error {
	_, err := a.db.NamedExec(`
		UPDATE containers
		SET service_id = :service_id,
			user_id = :user_id,
			image = :image,
			image_tag = :image_tag,
			status = :status,
			launch_on_startup = :launch_on_startup,
			name = :name,
			description = :description,
			color = :color,
			icon = :icon,
			command = :command
		WHERE id = :id
	`, c)
	return err
}

func (a *containerDBAdapter) DeleteContainer(ctx context.Context, id uuid.UUID) error {
	_, err := a.db.Exec(`
		DELETE FROM containers
		WHERE id = $1
	`, id)
	return err
}

func (a *containerDBAdapter) GetContainerTags(ctx context.Context, id uuid.UUID) (types.Tags, error) {
	var tags types.Tags
	err := a.db.Select(&tags, `
		SELECT tags.* FROM tags
		INNER JOIN container_tags ct on tags.id = ct.tag_id
		WHERE ct.container_id = $1
	`, id)
	if errors.Is(err, sql.ErrNoRows) {
		return tags, nil
	}
	return tags, err
}

func (a *containerDBAdapter) AddTag(ctx context.Context, id uuid.UUID, tagID types.TagID) error {
	_, err := a.db.Exec(`
		INSERT INTO container_tags (container_id, tag_id)
		VALUES ($1, $2)
	`, id, tagID)
	return err
}

func (a *containerDBAdapter) DeleteTags(ctx context.Context, id uuid.UUID) error {
	_, err := a.db.Exec(`
		DELETE FROM container_tags
		WHERE container_id = $1
	`, id)
	return err
}

func (a *containerDBAdapter) SetStatus(ctx context.Context, id uuid.UUID, status string) error {
	_, err := a.db.Exec(`
		UPDATE containers
		SET status = $1
		WHERE id = $2
	`, status, id)
	return err
}
