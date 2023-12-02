package port

import (
	"github.com/vertex-center/vertex/apps/admin/core/types"
	"github.com/vertex-center/vertex/pkg/user"
)

type (
	SshAdapter interface {
		GetAll() ([]types.PublicKey, error)
		Add(key string, username string) error
		Remove(fingerprint string, username string) error
		GetUsers() ([]user.User, error)
	}

	SshKernelAdapter interface {
		GetAll(users []user.User) ([]types.PublicKey, error)
		Add(key string, user user.User) error
		Remove(fingerprint string, user user.User) error
		GetUsers() ([]user.User, error)
	}
)