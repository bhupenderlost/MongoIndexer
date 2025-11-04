// Copyright (c) 2025 Bhupender Singh
// Use of this source code is governed by the MIT License
// that can be found in the LICENSE file.

package mongoindexer

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// mockIndexView allows us to simulate CreateMany behavior.
type mockIndexView struct {
	mock.Mock
}

func (m *mockIndexView) CreateMany(ctx context.Context, models []mongo.IndexModel, opts ...*options.CreateIndexesOptions) ([]string, error) {
	args := m.Called(ctx, models)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]string), args.Error(1)
}

// mockCollection implements CollectionAPI.
type mockCollection struct {
	mock.Mock
	iv   IndexView
	name string
}

func (m *mockCollection) Indexes() IndexView { return m.iv }
func (m *mockCollection) Name() string       { return m.name }

func TestCreateIndexes_Success(t *testing.T) {
	mi := &mockIndexView{}
	mc := &mockCollection{iv: mi, name: "users"}

	mi.On("CreateMany", mock.Anything, mock.MatchedBy(func(models []mongo.IndexModel) bool {
		// ensure at least one model present
		return len(models) >= 1
	})).Return([]string{"email_idx"}, nil).Once()

	i := New()
	err := i.CreateIndexes(context.Background(), mc, struct {
		Email string `bson:"email" mongoIndex:"unique,asc,name=email_idx"`
	}{},
	)

	assert.NoError(t, err)
	mi.AssertExpectations(t)
}

func TestCreateIndexes_ParseError(t *testing.T) {
	mc := &mockCollection{iv: &mockIndexView{}, name: "bad"}
	i := New()
	err := i.CreateIndexes(context.Background(), mc, struct {
		F string `bson:"f" mongoIndex:"ttl=bad"`
	}{},
	)
	assert.Error(t, err)
}

func TestCreateIndexes_CreateManyError(t *testing.T) {
	mi := &mockIndexView{}
	mc := &mockCollection{iv: mi, name: "users"}

	mi.On("CreateMany", mock.Anything, mock.Anything).Return(nil, errors.New("db fail")).Once()

	i := New()
	err := i.CreateIndexes(context.Background(), mc, struct {
		Email string `bson:"email" mongoIndex:"asc"`
	}{},
	)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "db fail")
	mi.AssertExpectations(t)
}

func TestCreateIndexes_NoIndexes(t *testing.T) {
	mc := &mockCollection{iv: &mockIndexView{}, name: "none"}
	i := New()
	err := i.CreateIndexes(context.Background(), mc, struct {
		ID string `bson:"_id"`
	}{},
	)
	assert.NoError(t, err)
}
