// Copyright (c) 2025 Bhupender Singh
// Use of this source code is governed by the MIT License
// that can be found in the LICENSE file.

package mongoindexer

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// CollectionAPI is a minimal interface representing what we need from a collection.
// This allows easy mocking in tests.
type CollectionAPI interface {
	Indexes() IndexView
	Name() string
}

// IndexView is a minimal interface matching the driver's index view used by this package.
// We only need CreateMany here but keeping it small makes testing simple.
type IndexView interface {
	CreateMany(ctx context.Context, models []mongo.IndexModel, opts ...*options.CreateIndexesOptions) ([]string, error)
}

// WrapCollection adapts *mongo.Collection to our CollectionAPI.
func WrapCollection(c *mongo.Collection) CollectionAPI {
	return &mongoCollection{c}
}

type mongoCollection struct {
	c *mongo.Collection
}

func (m *mongoCollection) Indexes() IndexView {
	return m.c.Indexes()
}

func (m *mongoCollection) Name() string { return m.c.Name() }

// Indexer is the public interface consumers will use.
type Indexer interface {
	// CreateIndexes inspects `model` (struct or pointer to struct), reads
	// `mongoIndex` tags, builds index models and creates them on the provided collection.
	CreateIndexes(ctx context.Context, coll CollectionAPI, model any) error
}

// New returns a default Indexer instance.
func New() Indexer {
	return &indexer{}
}

// concrete implementation
type indexer struct{}
