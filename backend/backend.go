package backend

import (
	"github.com/ringtail/pansidong/types"
	"encoding/json"
	"log"
)

func BackendFactory(bc types.BackendConfig) types.BackendStore {
	switch bc.Name() {
	case "boltdb":
		return NewBoltdbBackend(bc.Config().(*types.BoltDBConfig))
	}
	return nil
}

func CreateBackendConfig(backendType string, j_string string) types.BackendConfig {
	switch backendType {
	case "boltdb":
		bc := &types.BoltDBConfig{}
		err := json.Unmarshal([]byte(j_string), bc)
		if err != nil {
			log.Fatal("Failed to create %s backend config,because of %s", backendType, err.Error())
		}
		return bc
	}
	return nil
}
