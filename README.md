# MongoIndexer

_Automatic MongoDB Index Creation from Go Struct Tags_

`mongoindexer` is a lightweight and idiomatic Go package that automatically creates MongoDB indexes based on struct tags.
It’s designed to be **generic, testable, and easy to integrate** with any MongoDB project — no manual index management required.

---

## Features

- Automatically parses struct tags to generate indexes.
- Supports **unique**, **ascending/descending**, and **TTL** indexes.
- Simple, interface-based design for testing and extensibility.
- Fully covered with unit tests using `testify`.

---

## Installation

```bash
go get github.com/bhupenderlost/mongoindexer
```

---

## Usage

### 1. Define your struct

Add the `mongoIndex` tag to specify index configuration.

```go
package main

type User struct {
    Email     string `bson:"email" mongoIndex:"unique,asc,name=email_idx"`
    CreatedAt int64  `bson:"created_at" mongoIndex:"ttl=3600"`
}
```

---

### 2. Create indexes

```go
package main

import (
    "context"
    "log"

    "github.com/bhupenderlost/mongoindexer"
    "go.mongodb.org/mongo-driver/mongo"
    "go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
    ctx := context.Background()

    client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))
    if err != nil {
        log.Fatal(err)
    }

    db := client.Database("testdb")
    coll := db.Collection("users")

    indexer := mongoindexer.New()
  	wrapped := mongoindexer.WrapCollection(coll)

    if err := indexer.CreateIndexes(ctx, wrapped, User{}); err != nil {
        log.Fatal("Index creation failed:", err)
    }

    log.Println("Indexes created successfully!")
}
```

---

## Tag Syntax

| Option         | Description                          | Example                              |
| -------------- | ------------------------------------ | ------------------------------------ |
| `asc` / `desc` | Index order (ascending / descending) | `mongoIndex:"asc"`                   |
| `unique`       | Creates a unique index               | `mongoIndex:"unique,asc"`            |
| `ttl=<sec>`    | Sets a TTL (Time-to-Live) index      | `mongoIndex:"ttl=3600"`              |
| `name=<idx>`   | Custom index name                    | `mongoIndex:"unique,asc,name=email"` |

---

## Testing

This package comes with a complete set of unit tests using `testify` and `mock`.

Run tests with:

```bash
go test ./...
```

Example test:

```go
func TestCreateIndexes_Success(t *testing.T) {
    // mock IndexView, verify that CreateMany was called
}
```

---

## Project Structure

```
indexer/
│
├── indexer.go            # Core logic for parsing and creating indexes
├── parser.go             # Tag parsing and validation
├── interfaces.go         # Mongo interfaces for abstraction
└── indexer_test.go       # Unit tests using testify mocks for indexer
└── parser_test.go        # Unit tests using testify mocks for parser

```

---

## Design Philosophy

- **Interfaces-first**: Easily mockable for testing (`CollectionAPI`, `IndexView`).
- **Generic**: Works with any MongoDB collection that implements the minimal interface.
- **Declarative**: Indexes are defined once in the struct — no manual duplication.

---

## Example Output

When executed with the example `User` struct, the following indexes are created:

```json
[
  {
    "key": { "email": 1 },
    "name": "email_idx",
    "unique": true
  },
  {
    "key": { "created_at": 1 },
    "expireAfterSeconds": 3600
  }
]
```

---

## Contributing

1. Fork the repository.
2. Create a feature branch:

   ```bash
   git checkout -b feature/awesome-index
   ```

3. Commit your changes:

   ```bash
   git commit -m "feat: add new index rule"
   ```

4. Push and open a PR

---

## License

Licensed under the MIT License.  
See [LICENSE](./LICENSE) for more information.
