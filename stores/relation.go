package stores

import (
	"context"

	"github.com/sourcenetwork/raccoondb/iterator"
)

// maybe a relation store which keeps an index of related keys in both directions
// this is definitely necessary, i dig it and it's pretty simple
// it can be as simple as a kv storing all the node keys, an index store in one direction and an index in the reverse direction
// plus a catalogue which contains all this data.

type RelationStore struct {
	nodes     CountedKVStore
	sucessors FieldIndexStore
	ancestors FieldIndexStore
}

func (s *RelationStore) Create(ctx context.Context, src, dst []byte) (RecordCreated, error) {
	_, err := s.nodes.Set(ctx, src, src)
	if err != nil {
		return false, err
	}
	_, err = s.nodes.Set(ctx, dst, dst)
	if err != nil {
		return false, err
	}

	err = s.sucessors.IndexValue(ctx, src, dst)
	if err != nil {
		return false, err
	}

	err = s.ancestors.IndexValue(ctx, dst, src)
	if err != nil {
		return false, err
	}

	return false, nil // TODO
}

func (s *RelationStore) Delete(ctx context.Context, src, dst []byte) error {
	// TODO check if node needs removing

	/*
		err := s.sucessors.Delete(ctx, src, dst)
		if err != nil {
			return err
		}

		err = s.ancestors.Delete(ctx, dst, src)
		if err != nil {
			return err
		}
	*/

	return nil
}

func (s *RelationStore) IterateSucessors(ctx context.Context, node []byte) (iterator.BytesIterator, error) {
	return s.sucessors.GetBucketValues(ctx, node)
}

func (s *RelationStore) IterateAncestors(ctx context.Context, node []byte) (iterator.BytesIterator, error) {
	return s.ancestors.GetBucketValues(ctx, node)

}

func (s *RelationStore) IterateNodes(ctx context.Context) (iterator.BytesIterator, error) {
	//opts := types.NewOpenIterator()
	//return s.nodes.Iterate(ctx, opts)
	return nil, nil
}

func (s *RelationStore) Has(ctx context.Context, src, dst []byte) (bool, error) {
	return s.sucessors.Has(ctx, src, dst)
}

func (s *RelationStore) GetNodeCount(ctx context.Context) (uint64, error) {
	return s.nodes.GetCount(ctx)
}
