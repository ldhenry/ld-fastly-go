package kvdatastore

import (
	"encoding/json"
	"errors"

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
	store.loggers.Debug("Closing FastlyKVDataStore. This is a no-op.")
	return nil
}

func (store *fastlyKVDataStoreImpl) Init(allData []ldstoretypes.SerializedCollection) error {
	store.loggers.Debug("Init called. This is not supported for FastlyKVDataStore.")
	return errors.New("init not supported for FastlyKVDataStore")
}

func (store *fastlyKVDataStoreImpl) IsInitialized() bool {
	store.loggers.Debug("IsInitialized called. This is not supported for FastlyKVDataStore.")
	return false
}

func (store *fastlyKVDataStoreImpl) IsStoreAvailable() bool {
	store.loggers.Debug("Checking if FastlyKVDataStore is available")
	_, err := store.getKV()
	return err == nil
}

type allDataStruct struct {
	Flags map[string]interface{} `json:"flags"`
}

func (store *fastlyKVDataStoreImpl) getAllFlagData() (allDataStruct, error) {
	o, err := store.getKV()
	if err != nil {
		return allDataStruct{}, err
	}

	v, err := o.Lookup(ENV_KEY_PREFIX + store.clientSideID)
	if err != nil {
		return allDataStruct{}, err
	}

	var allData allDataStruct
	err = json.Unmarshal([]byte(v.String()), &allData)
	if err != nil {
		return allDataStruct{}, err
	}

	return allData, nil
}

func (store *fastlyKVDataStoreImpl) Get(kind ldstoretypes.DataKind, key string) (ldstoretypes.SerializedItemDescriptor, error) {
	store.loggers.Debug("Getting item from FastlyKVDataStore")

	allData, err := store.getAllFlagData()
	if err != nil {
		return ldstoretypes.SerializedItemDescriptor{}, err
	}

	flag, ok := allData.Flags[key]
	if !ok {
		return ldstoretypes.SerializedItemDescriptor{}.NotFound(), nil
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
	store.loggers.Debug("Getting all items from FastlyKVDataStore")

	allData, err := store.getAllFlagData()
	if err != nil {
		return nil, err
	}

	flagDescriptors := make([]ldstoretypes.KeyedSerializedItemDescriptor, 0, len(allData.Flags))
	for k, v := range allData.Flags {
		flagBytes, err := json.Marshal(v)
		if err != nil {
			return nil, err
		}
		flagDescriptors = append(flagDescriptors, ldstoretypes.KeyedSerializedItemDescriptor{
			Key: k,
			Item: ldstoretypes.SerializedItemDescriptor{
				SerializedItem: flagBytes,
			},
		})
	}
	return flagDescriptors, nil
}

// Upsert should not be needed since the KV store is updated by LaunchDarkly directly
func (store *fastlyKVDataStoreImpl) Upsert(
	kind ldstoretypes.DataKind,
	key string,
	newItem ldstoretypes.SerializedItemDescriptor,
) (bool, error) {
	store.loggers.Debug("Upsert called. This is not supported for FastlyKVDataStore.")
	return false, errors.New("upsert not supported for FastlyKVDataStore")
}
