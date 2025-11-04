// Copyright (c) 2025 Bhupender Singh
// Use of this source code is governed by the MIT License
// that can be found in the LICENSE file.

package mongoindexer

import (
	"testing"

	"go.mongodb.org/mongo-driver/bson"
)

func TestBuildIndexFromTag_Valid(t *testing.T) {
	im, err := buildIndexFromTag("email", "unique,asc,name=email_idx")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if im.Keys == nil {
		t.Fatalf("keys nil")
	}

	d := im.Keys.(bson.D)
	if d[0].Key != "email" || d[0].Value != int32(1) {
		t.Fatalf("unexpected keys: %v", d)
	}

	if im.Options == nil || im.Options.Name == nil || *im.Options.Name != "email_idx" {
		t.Fatalf("expected name=email_idx")
	}

	if im.Options.Unique == nil || *im.Options.Unique != true {
		t.Fatalf("expected unique=true")
	}
}

func TestBuildIndexFromTag_TTL(t *testing.T) {
	im, err := buildIndexFromTag("createdAt", "desc,ttl=3600")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if im.Options == nil || im.Options.ExpireAfterSeconds == nil || *im.Options.ExpireAfterSeconds != int32(3600) {
		t.Fatalf("expected ttl=3600")
	}
}

func TestBuildIndexFromTag_Partial(t *testing.T) {
	// partial expects JSON object
	partial := `{"age":{"$gt":18}}`
	im, err := buildIndexFromTag("age", "asc,partial="+partial)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if im.Options == nil || im.Options.PartialFilterExpression == nil {
		t.Fatalf("expected partial filter expression set")
	}
}

func TestBuildIndexFromTag_BadTTL(t *testing.T) {
	_, err := buildIndexFromTag("f", "ttl=notint")
	if err == nil {
		t.Fatalf("expected error for bad ttl")
	}
}
