package port

import (
	"github.com/vertex-center/vertex/apps/containers/core/types"
	sqltypes "github.com/vertex-center/vertex/apps/sql/core/types"
)

type SqlService interface {
	Get(inst *types.Container) (sqltypes.DBMS, error)
	EnvCredentials(inst *types.Container, user string, pass string) (types.EnvVariables, error)
}
