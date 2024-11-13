package composite

import (
	"github.com/sourcenetwork/raccoondb/marshal"
	"github.com/sourcenetwork/raccoondb/stores"
)

const objsPrefix = "objs"
const idxsPrefix = "indexes"

type Index[T any, I any] struct {
	Name      string
	Extractor IndexValueExtractor[T, I]
	Marshaler marshal.Marshaler[I]
	store     *ObjectIndexStore[T, I]
}

func (i *Index[T, I]) getObjectIndex() *ObjectIndexStore[T, I] {
	return i.store
}

func (i *Index[T, I]) setObjectIndex(store *ObjectIndexStore[T, I]) {
	i.store = store
}

type IndexedObjectStoreSchema[T any] struct {
	Marshaler marshal.Marshaler[T]
	Indexes   []*Index[T, any]
}

func NewStoreManager[T any](schema *IndexedObjectStoreSchema[T]) *StoreManager[T] {
	return &StoreManager[T]{}
}

type StoreManager[T any] struct {
	schema  IndexedObjectStoreSchema[T]
	store   *IndexedObjectStore[T]
	indexes []*ObjectIndexStore[T, any]
}

func (m *StoreManager[T]) Initialize(kv stores.KVStore) {
	indexesKv := stores.NewPrefixedKV(kv, []byte(idxsPrefix))
	var idxs []*ObjectIndexStore[T, any]
	for _, def := range m.schema.Indexes {
		idxKv := stores.NewPrefixedKV(indexesKv, []byte(def.Name))
		fieldIdx := stores.NewFieldIndexStore(idxKv)
		idx := newObjectIndexStore(def.Name, &fieldIdx, def.Extractor, def.Marshaler)

		def.store = idx
		idxs = append(idxs, idx)
	}

	objKv := stores.NewPrefixedKV(kv, []byte(objsPrefix))
	keyObjStore := stores.NewKeyObjectStore(objKv, m.schema.Marshaler)

	objStores := newIndexedObjectStore(&keyObjStore, idxs)

	m.store = &objStores
	m.indexes = idxs
}

func (m *StoreManager[T]) GetIndexedObjectStore() *IndexedObjectStore[T] {
	return m.store
}

func GetIndex[T, I any](manager *StoreManager[T], idx *Index[T, I]) *ObjectIndexStore[T, I] {
	s := idx.getObjectIndex()
	if s == nil {
		panic("manager does not manage given index")
	}
	return s
}
