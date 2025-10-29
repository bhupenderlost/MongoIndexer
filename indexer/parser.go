package mongoindexer

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"
	"strings"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// parseModelIndexes inspects a struct type and builds mongo.IndexModel slice based on `mongoIndex` tags.
func parseModelIndexes(model any) ([]mongo.IndexModel, error) {
	if model == nil {
		return nil, fmt.Errorf("mongoindexer: model is nil")
	}

	t := reflect.TypeOf(model)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	if t.Kind() != reflect.Struct {
		return nil, fmt.Errorf("mongoindexer: model must be a struct, got %s", t.Kind())
	}

	var results []mongo.IndexModel

	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)
		tag := f.Tag.Get("mongoIndex")
		if tag == "" {
			continue
		}

		bsonTag := f.Tag.Get("bson")
		if bsonTag == "" {
			bsonTag = strings.ToLower(f.Name)
		} else {
			bsonTag = strings.Split(bsonTag, ",")[0]
		}

		im, err := buildIndexFromTag(bsonTag, tag)
		if err != nil {
			return nil, fmt.Errorf("field %s: %w", f.Name, err)
		}

		results = append(results, im)
	}

	return results, nil
}

// buildIndexFromTag converts a single field+tag into mongo.IndexModel.
// Tag examples:
//
//	"unique,asc,name=email_idx"
//	"desc,ttl=3600"
//	"asc,partial={\"age\":{\"$gt\":18}}"
func buildIndexFromTag(field, tag string) (mongo.IndexModel, error) {
	opts := options.Index()
	direction := int32(1)

	parts := strings.Split(tag, ",")
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p == "" {
			continue
		}

		switch {
		case p == "unique":
			opts.SetUnique(true)
		case p == "asc":
			direction = 1
		case p == "desc":
			direction = -1
		case p == "sparse":
			opts.SetSparse(true)
		case strings.HasPrefix(p, "name="):
			opts.SetName(strings.TrimPrefix(p, "name="))
		case strings.HasPrefix(p, "ttl="):
			ttlStr := strings.TrimPrefix(p, "ttl=")
			ttlVal, err := strconv.Atoi(ttlStr)
			if err != nil {
				return mongo.IndexModel{}, fmt.Errorf("invalid ttl value: %s", ttlStr)
			}
			opts.SetExpireAfterSeconds(int32(ttlVal))
		case strings.HasPrefix(p, "partial="):
			// partial expression expects a JSON object encoded inside the tag
			re := strings.TrimPrefix(p, "partial=")
			// try to decode JSON into bson.M
			var m bson.M
			if err := json.Unmarshal([]byte(re), &m); err != nil {
				return mongo.IndexModel{}, fmt.Errorf("invalid partial filter expression: %w", err)
			}
			opts.SetPartialFilterExpression(m)
		default:
			// unknown token - ignore but don't fail
		}
	}

	return mongo.IndexModel{
		Keys:    bson.D{{Key: field, Value: direction}},
		Options: opts,
	}, nil
}
