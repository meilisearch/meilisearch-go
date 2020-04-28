<p align="center">
  <img src="https://res.cloudinary.com/meilisearch/image/upload/v1587402338/SDKs/meilisearch_go.svg" alt="MeiliSearch-Go" width="200" height="200" />
</p>

<h1 align="center">MeiliSearch Go</h1>

<h4 align="center">
  <a href="https://github.com/meilisearch/MeiliSearch">MeiliSearch</a> |
  <a href="https://www.meilisearch.com">Website</a> |
  <a href="https://blog.meilisearch.com">Blog</a> |
  <a href="https://twitter.com/meilisearch">Twitter</a> |
  <a href="https://docs.meilisearch.com">Documentation</a> |
  <a href="https://docs.meilisearch.com/faq">FAQ</a>
</h4>

<p align="center">
  <a href="https://github.com/meilisearch/meilisearch-go/actions"><img src="https://img.shields.io/github/workflow/status/meilisearch/meilisearch-go/Go" alt="GitHub Workflow Status"></a>
  <a href="https://goreportcard.com/report/github.com/meilisearch/meilisearch-go"><img src="https://goreportcard.com/badge/github.com/meilisearch/meilisearch-go" alt="Test"></a>
  <a href="https://github.com/meilisearch/meilisearch-go/blob/master/LICENSE"><img src="https://img.shields.io/badge/license-MIT-informational" alt="License"></a>
  <a href="https://slack.meilisearch.com"><img src="https://img.shields.io/badge/slack-MeiliSearch-blue.svg?logo=slack" alt="Slack"></a>
</p>

<p align="center">⚡ Lightning Fast, Ultra Relevant, and Typo-Tolerant Search Engine MeiliSearch client written in Go</p>

**MeiliSearch Go** is a client for **MeiliSearch** written in Go. **MeiliSearch** is a powerful, fast, open-source, easy to use and deploy search engine. Both searching and indexing are highly customizable. Features such as typo-tolerance, filters, and synonyms are provided out-of-the-box.

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
        APIKey: "masterKey",
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
documents := []map[string]interface{}{
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
    AttributesToHighlight: []string{"*"},
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

This module works for MeiliSearch `v0.9.x`.
