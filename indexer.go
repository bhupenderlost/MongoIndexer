// Copyright (c) 2025 Bhupender Singh
// Use of this source code is governed by the MIT License
// that can be found in the LICENSE file.

package mongoindexer

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/mongo"
)

// CreateIndexes inspects the model and creates indexes on the provided collection.
func (i *indexer) CreateIndexes(ctx context.Context, coll CollectionAPI, model any) error {
	ims, err := parseModelIndexes(model)
	if err != nil {
		return fmt.Errorf("mongoindexer: parseModelIndexes failed: %w", err)
	}

	if len(ims) == 0 {
		return nil
	}

	// CreateMany accepts []mongo.IndexModel. If a user supplies identical index
	// definitions that already exist, the driver will return an error only if
	// they conflict. We simply forward the error.
	names, err := coll.Indexes().CreateMany(ctx, ims)
	if err != nil {
		return fmt.Errorf("mongoindexer: create many failed for %s: %w", coll.Name(), err)
	}

	// Friendly log - in real library we could swap this for a pluggable logger.
	fmt.Printf("mongoindexer: created indexes for %s: %v\n", coll.Name(), names)
	return nil
}

// Convenience: wrap mongo.Collection to our CollectionAPI without polluting callers.
func WrapMongoCollection(c *mongo.Collection) CollectionAPI { return WrapCollection(c) }
