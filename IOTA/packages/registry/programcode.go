package registry

import (
	"fmt"

	"wasp/packages/hashing"
	"wasp/plugins/database"
	"wasp/plugins/publisher"
)

func dbkeyProgramCode(progHash *hashing.HashValue) []byte {
	return database.MakeKey(database.ObjectTypeProgramCode, progHash[:])
}

// TODO save program code in the smart contract state
func GetProgramCode(progHash *hashing.HashValue) ([]byte, error) {
	db := database.GetRegistryPartition()
	data, err := db.Get(dbkeyProgramCode(progHash))
	if err != nil {
		return nil, err
	}
	hash := hashing.HashData(data)
	if *hash != *progHash {
		return nil, fmt.Errorf("program code is corrupted. Expected: %s. Got: %s", progHash.String(), hash.String())
	}
	return data, nil
}

func SaveProgramCode(programCode []byte) (ret hashing.HashValue, err error) {
	progHash := hashing.HashData(programCode)
	db := database.GetRegistryPartition()
	if err = db.Set(dbkeyProgramCode(progHash), programCode); err != nil {
		return
	}
	ret = *progHash

	defer publisher.Publish("programcode", progHash.String())
	return
}
