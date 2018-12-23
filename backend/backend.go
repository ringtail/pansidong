package backend

import (
	"github.com/ringtail/pansidong/types"
)

func BackendFactory(bc types.BackendConfig) types.BackendStore {
	switch bc.Name() {
	case "boltdb":
		return NewBoltdbBackend(bc.Config().(*types.BoltDBConfig))
	}
	return nil
}


