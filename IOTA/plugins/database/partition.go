package database

import (
	"bytes"
	"sync"
	"wasp/packages/parameters"

	"github.com/iotaledger/goshimmer/dapps/valuetransfers/packages/address"
	"github.com/iotaledger/goshimmer/packages/database"
	"github.com/iotaledger/hive.go/kvstore"
	"github.com/iotaledger/hive.go/logger"
)

// database is structured with 34 byte long prefixes 'address' || 'object type byte
// '
const (
	ObjectTypeDBSchemaVersion byte = iota
	ObjectTypeBootupData
	ObjectTypeDistributedKeyData
	ObjectTypeSolidState
	ObjectTypeStateUpdateBatch
	ObjectTypeProcessedRequestId
	ObjectTypeSolidStateIndex
	ObjectTypeStateVariable
	ObjectTypeProgramMetadata
	ObjectTypeProgramCode
)

type Partition struct {
	kvstore.KVStore
	mut *sync.RWMutex
}

var (
	// to be able to work with MapsDB
	partitions      = make(map[address.Address]*Partition)
	partitionsMutex sync.RWMutex
)

// storeInstance returns the KVStore instance.
func storeInstance() kvstore.KVStore {
	storeOnce.Do(createStore)
	return store
}

// storeRealm is a factory method for a different realm backed by the KVStore instance.
func storeRealm(realm kvstore.Realm) kvstore.KVStore {
	return storeInstance().WithRealm(realm)
}

// Partition returns store prefixed with the smart contract address
// Wasp ledger is partitioned by smart contract addresses
// cached to be able to work with MapsDB TODO
func GetPartition(addr *address.Address) *Partition {
	partitionsMutex.RLock()
	ret, ok := partitions[*addr]
	if ok {
		defer partitionsMutex.RUnlock()
		return ret
	}
	// switching to write lock
	partitionsMutex.RUnlock()
	partitionsMutex.Lock()
	defer partitionsMutex.Unlock()

	partitions[*addr] = &Partition{
		KVStore: storeRealm(addr[:]),
		mut:     &sync.RWMutex{},
	}
	return partitions[*addr]
}

func GetRegistryPartition() kvstore.KVStore {
	var niladdr address.Address
	return GetPartition(&niladdr)
}

func (part *Partition) RLock() {
	part.mut.RLock()
}

func (part *Partition) RUnlock() {
	part.mut.RUnlock()
}

func (part *Partition) Lock() {
	part.mut.Lock()
}

func (part *Partition) Unlock() {
	part.mut.Unlock()
}

// MakeKey makes key within the partition. It consists to one byte for object type
// and arbitrary byte fragments concatenated together
func MakeKey(objType byte, keyBytes ...[]byte) []byte {
	var buf bytes.Buffer
	buf.WriteByte(objType)
	for _, b := range keyBytes {
		buf.Write(b)
	}
	return buf.Bytes()
}

func createStore() {
	log = logger.NewLogger(PluginName)

	var err error
	if parameters.GetBool(parameters.DatabaseInMemory) {
		log.Infof("IN MEMORY DATABASE")
		db, err = database.NewMemDB()
	} else {
		dbDir := parameters.GetString(parameters.DatabaseDir)
		db, err = database.NewDB(dbDir)
	}
	if err != nil {
		log.Fatal(err)
	}

	store = db.NewStore()
}
