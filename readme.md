# MeiliSearch Go Client <!-- omit in toc -->

![GitHub Workflow Status](https://img.shields.io/github/workflow/status/meilisearch/meilisearch-go/Go)
[![Go Report Card](https://goreportcard.com/badge/github.com/meilisearch/meilisearch-go)](https://goreportcard.com/report/github.com/meilisearch/meilisearch-go)
[![Licence](https://img.shields.io/badge/licence-MIT-blue.svg)](https://img.shields.io/badge/licence-MIT-blue.svg)

The go client for MeiliSearch API.

MeiliSearch provides an ultra relevant and instant full-text search. Our solution is open-source and you can check out [our repository here](https://github.com/meilisearch/MeiliSearch).

Here is the [MeiliSearch documentation](https://docs.meilisearch.com/) 📖

## Table of Contents <!-- omit in toc -->

- [🔧 Installation](#-installation)
- [🚀 Getting started](#-getting-started)
- [🎬 Examples](#-examples)
  - [Indexes](#indexes)
  - [Documents](#documents)
  - [Update status](#update-status)
  - [Search](#search)
- [⚙️ Development Workflow](#️-development-workflow)
  - [Install Go](#install-go)
  - [Install dependencies](#install-dependencies)
  - [Tests and Linter](#tests-and-linter)
- [🤖 Compatibility with MeiliSearch](#-compatibility-with-meilisearch)

## 🔧 Installation

With `go get` in command line:
```bash
$ go get github.com/meilisearch/meilisearch-go
```

### Run MeiliSearch <!-- omit in toc -->

There are many easy ways to [download and run a MeiliSearch instance](https://docs.meilisearch.com/guides/advanced_guides/installation.html#download-and-launch).

For example, if you use Docker:
```bash
$ docker run -it --rm -p 7700:7700 getmeili/meilisearch:latest --master-key=masterKey
```

NB: you can also download MeiliSearch from **Homebrew** or **APT**.

## 🚀 Getting started

#### Add documents <!-- omit in toc -->

```go
package main

import (
    "fmt"
    "os"

    "github.com/meilisearch/meilisearch-go"
)

func main() {
    var client = meilisearch.NewClient(meilisearch.Config{
        Host: "http://127.0.0.1:7700",
        APIKey: "masterKey"
    })

    // Create an index if your index does not already exist
    _, err := client.Indexes().Create(meilisearch.CreateIndexRequest{
        UID: "books",
    })
    if err != nil {
        fmt.Println(err)
        os.Exit(1)
    }

    documents := []map[string]interface{}{
        {"book_id": 123,  "title": "Pride and Prejudice"},
        {"book_id": 456,  "title": "Le Petit Prince"},
        {"book_id": 1,    "title": "Alice In Wonderland"},
        {"book_id": 1344, "title": "The Hobbit"},
        {"book_id": 4,    "title": "Harry Potter and the Half-Blood Prince"},
        {"book_id": 42,   "title": "The Hitchhiker's Guide to the Galaxy"},
    }

    updateRes, err := client.Documents("books").AddOrUpdate(documents) // => { "updateId": 0 }
    if err != nil {
        fmt.Println(err)
        os.Exit(1)
    }

    fmt.Println(updateRes.UpdateID)
}
```

With the `updateId`, you can check the status (`processed` or `failed`) of your documents addition thanks to this [method](#update-status).

#### Search in index <!-- omit in toc -->

```go
package main

import (
    "fmt"
    "os"

    "github.com/meilisearch/meilisearch-go"
)

func main() {
    // MeiliSearch is typo-tolerant:
    searchRes, err := client.Search("books").Search(meilisearch.SearchRequest{
        Query: "harry pottre",
        Limit: 10,
    })
    if err != nil {
        fmt.Println(err)
        os.Exit(1)
    }

    fmt.Println(searchRes.Hits)
}
```

JSON output:
```json
{
  "hits": [{
    "book_id": 4,
    "title": "Harry Potter and the Half-Blood Prince"
  }],
  "offset": 0,
  "limit": 10,
  "processingTimeMs": 1,
  "query": "harry pottre"
}
```

## 🎬 Examples

All HTTP routes of MeiliSearch are accessible via methods in this SDK.</br>
You can check out [the API documentation](https://docs.meilisearch.com/references/).

### Indexes

#### Create an index <!-- omit in toc -->

```go
// Create an index with a specific uid (uid must be unique)
resp, err := client.Indexes().Create(meilisearch.CreateIndexRequest{
    UID: "books",
})
// Create an index with a primary key
resp, err := client.Indexes().Create(meilisearch.CreateIndexRequest{
    UID: "books",
    PrimaryKey: "book_id",
})
```

#### List all indexes <!-- omit in toc -->

```go
list, err := client.Indexes().List()
```

#### Get an index object <!-- omit in toc -->

```go
index, err := client.Indexes().Get("books")
```

### Documents

#### Fetch documents <!-- omit in toc -->

```go
// Get one document
var document map[int]interface{}
err := client.Documents("books").Get("123", &doc)
// Get documents by batch
var list []map[int]interface{}
err = client.Documents("books").List(ListDocumentsRequest{
    Offset: 0,
    Limit:  10,
}, &list)
```

#### Add documents <!-- omit in toc -->

```go
documents := []Book{
    {BookID: 90, Title: "Madame Bovary"},
}

upd_res, err := client.Documents("books").AddOrUpdate(documents)
```

Response:
```json
{
    "updateId": 1
}
```
With this `updateId` you can track your [operation update](#update-status).

#### Delete documents <!-- omit in toc -->

```go
// Delete one document
updateRes, err = client.Documents("books").Delete("123")
// Delete several documents
updateRes, err = client.Documents("books").Deletes([]string{"123", "456"})
// Delete all documents /!\
updateRes, err = client.Documents("books").DeleteAllDocuments()
```

### Update status

```go
// Get one update status
// Parameter: the updateId got after an asynchronous request (e.g. documents addition)
update, err := client.Updates("books").Get(1)
// Get all update satus
list, err := client.Updates("books").List()
```

### Search

#### Basic search <!-- omit in toc -->

```go
resp, err := client.Search(indexUID).Search(meilisearch.SearchRequest{
    Query: "prince",
    Limit: 10,
})
```

```json
{
    "hits": [
        {
            "book_id": 456,
            "title": "Le Petit Prince"
        },
        {
            "book_id": 4,
            "title": "Harry Potter and the Half-Blood Prince"
        }
    ],
    "offset": 0,
    "limit": 20,
    "processingTimeMs": 13,
    "query": "prince"
}
```

#### Custom search <!-- omit in toc -->

All the supported options are described in [this documentation section](https://docs.meilisearch.com/references/search.html#search-in-an-index).

```go
resp, err := client.Search(indexUID).Search(meilisearch.SearchRequest{
    Query: "harry pottre",
    AttributesToHighlight: "*",
})
```

```json
{
    "hits": [
        {
            "book_id": 456,
            "title": "Le Petit Prince",
            "_formatted": {
                "book_id": 456,
                "title": "Le Petit <em>Prince</em>"
            }
        }
    ],
    "offset": 0,
    "limit": 1,
    "processingTimeMs": 0,
    "query": "prince"
}
```

## ⚙️ Development Workflow

If you want to contribute, this section describes the steps to follow.

Thank you for your interest in a MeiliSearch tool! ♥️

### Install Go

Follow the official [tutorial](https://golang.org/doc/install)

### Install dependencies

```bash
$ go get -v -t -d ./...
```

### Tests and Linter

Each PR should pass the tests and the linter to be accepted.

```bash
# Tests
$ ./run_tests.sh
# Install golint
$ go get -u golang.org/x/lint/golint
# Use golint
$ golint
# Use gofmt
$ gofmt -w ./..
```

## 🤖 Compatibility with MeiliSearch

This gem works for MeiliSearch `v0.9.x`.
