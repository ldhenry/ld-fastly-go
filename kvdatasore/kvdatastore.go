package kvdatastore

import (
	"github.com/launchdarkly/go-server-sdk/v7/subsystems"
)

type StoreBuilder[T any] struct {
	builderOptions builderOptions
	factory        func(*StoreBuilder[T], subsystems.ClientContext) (T, error)
}

type builderOptions struct {
	clientSideID string
	kvStoreName  string
}

func (b *StoreBuilder[T]) ClientSideID(clientSideID string) *StoreBuilder[T] {
	b.builderOptions.clientSideID = clientSideID
	return b
}

func (b *StoreBuilder[T]) KvStoreName(kvStoreName string) *StoreBuilder[T] {
	b.builderOptions.kvStoreName = kvStoreName
	return b
}

func (b *StoreBuilder[T]) Build(context subsystems.ClientContext) (T, error) {
	return b.factory(b, context)
}

func DataStore() *StoreBuilder[subsystems.PersistentDataStore] {
	return &StoreBuilder[subsystems.PersistentDataStore]{
		factory: createPersistentDataStore,
	}
}

func createPersistentDataStore(
	builder *StoreBuilder[subsystems.PersistentDataStore],
	clientContext subsystems.ClientContext,
) (subsystems.PersistentDataStore, error) {
	store := newFastlyKVDataStoreImpl(builder.builderOptions, clientContext.GetLogging().Loggers)
	return store, nil
}
