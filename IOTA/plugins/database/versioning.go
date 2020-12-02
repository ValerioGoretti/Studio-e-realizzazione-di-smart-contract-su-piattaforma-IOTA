package database

import (
	"bytes"
	"errors"
	"fmt"
	"wasp/packages/hashing"

	"github.com/iotaledger/goshimmer/dapps/valuetransfers/packages/address"

	"github.com/iotaledger/hive.go/kvstore"
)

const (
	// DBVersion defines the version of the database schema this version of Wasp supports.
	// Every time there's a breaking change regarding the stored data, this version flag should be adjusted.
	DBVersion = 0
)

var (
	// ErrDBVersionIncompatible is returned when the database has an unexpected version.
	ErrDBVersionIncompatible = errors.New("database version is not compatible. please delete your database folder and restart")
)

// checks whether the database is compatible with the current schema version.
// also automatically sets the version if the database if new.
// version is stored in niladdr partition.
// it consists of one byte of version and the hash (checksum) of that one byte
func checkDatabaseVersion() error {
	var niladdr address.Address
	db := GetPartition(&niladdr)
	ver, err := db.Get(MakeKey(ObjectTypeDBSchemaVersion))

	var versiondata [1 + hashing.HashSize]byte
	versiondata[0] = DBVersion
	copy(versiondata[1:], hashing.HashStrings(fmt.Sprintf("dbversion = %d", DBVersion)).Bytes())

	if err == kvstore.ErrKeyNotFound {
		// set the version in an empty DB
		return db.Set(MakeKey(ObjectTypeDBSchemaVersion), versiondata[:])
	}
	if err != nil {
		return err
	}
	if len(ver) == 0 {
		return fmt.Errorf("%w: no database version was persisted", ErrDBVersionIncompatible)
	}
	if !bytes.Equal(ver, versiondata[:]) {
		return fmt.Errorf("%w: supported version: %d, version of database: %d", ErrDBVersionIncompatible, DBVersion, ver[0])
	}
	return nil
}
