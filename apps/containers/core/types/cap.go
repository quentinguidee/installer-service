package types

import "github.com/vertex-center/vertex/common/uuid"

type (
	Capabilities []Capability
	Capability   struct {
		ContainerID uuid.UUID `json:"container_id" db:"container_id" example:"d1fb743c-f937-4f3d-95b9-1a8475464591"`
		Name        string    `json:"name"         db:"name"`
	}
)
