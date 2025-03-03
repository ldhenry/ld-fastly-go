package kvdatastore

import (
	"encoding/json"
	"fmt"

	"github.com/fastly/compute-sdk-go/kvstore"
	"github.com/launchdarkly/go-sdk-common/v3/ldlog"
	"github.com/launchdarkly/go-server-sdk/v7/subsystems"
	"github.com/launchdarkly/go-server-sdk/v7/subsystems/ldstoretypes"
)

const ENV_KEY_PREFIX = "LD-Env-"

type fastlyKVDataStoreImpl struct {
	loggers      ldlog.Loggers
	clientSideID string
	kvStoreName  string
}

var _ subsystems.PersistentDataStore = (*fastlyKVDataStoreImpl)(nil)

func newFastlyKVDataStoreImpl(
	builder builderOptions,
	loggers ldlog.Loggers,
) *fastlyKVDataStoreImpl {
	impl := &fastlyKVDataStoreImpl{
		loggers:      loggers,
		clientSideID: builder.clientSideID,
		kvStoreName:  builder.kvStoreName,
	}
	impl.loggers.SetPrefix("FastlyKVDataStore:")

	return impl
}

func (store *fastlyKVDataStoreImpl) getKV() (*kvstore.Store, error) {
	return kvstore.Open(store.kvStoreName)
}

func (store *fastlyKVDataStoreImpl) Close() error {
	store.loggers.Info("Closing FastlyKVDataStore")
	return nil
}

func (store *fastlyKVDataStoreImpl) Init(allData []ldstoretypes.SerializedCollection) error {
	store.loggers.Info("Initializing FastlyKVDataStore")
	_, err := store.getKV()
	if err != nil {
		return err
	}
	return nil
}

func (store *fastlyKVDataStoreImpl) IsInitialized() bool {
	store.loggers.Info("Checking if FastlyKVDataStore is initialized")
	_, err := store.getKV()
	return err == nil
}
func (store *fastlyKVDataStoreImpl) IsStoreAvailable() bool {
	store.loggers.Info("Checking if FastlyKVDataStore is available")
	_, err := store.getKV()
	return err == nil
}

type allDataStruct struct {
	Flags map[string]interface{} `json:"flags"`
}

func (store *fastlyKVDataStoreImpl) Get(kind ldstoretypes.DataKind, key string) (ldstoretypes.SerializedItemDescriptor, error) {
	store.loggers.Info("Getting item from FastlyKVDataStore")
	o, err := store.getKV()
	if err != nil {
		return ldstoretypes.SerializedItemDescriptor{}, err
	}

	v, err := o.Lookup(ENV_KEY_PREFIX + store.clientSideID)
	if err != nil {
		return ldstoretypes.SerializedItemDescriptor{}, err
	}

	var allData allDataStruct
	err = json.Unmarshal([]byte(v.String()), &allData)
	if err != nil {
		return ldstoretypes.SerializedItemDescriptor{}, err
	}

	flag, ok := allData.Flags["animal"]
	if !ok {
		return ldstoretypes.SerializedItemDescriptor{}, fmt.Errorf("flag not found")
	}

	flagBytes, err := json.Marshal(flag)
	if err != nil {
		return ldstoretypes.SerializedItemDescriptor{}, err
	}

	return ldstoretypes.SerializedItemDescriptor{
		SerializedItem: flagBytes,
	}, nil
}

func (store *fastlyKVDataStoreImpl) GetAll(kind ldstoretypes.DataKind) ([]ldstoretypes.KeyedSerializedItemDescriptor, error) {
	store.loggers.Info("Getting all items from FastlyKVDataStore")
	return nil, nil
}

func (store *fastlyKVDataStoreImpl) Upsert(
	kind ldstoretypes.DataKind,
	key string,
	newItem ldstoretypes.SerializedItemDescriptor,
) (bool, error) {
	store.loggers.Info("Upserting item to FastlyKVDataStore")
	return false, nil
}
