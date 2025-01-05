<p align="center">
  <img src="https://raw.githubusercontent.com/meilisearch/integration-guides/main/assets/logos/meilisearch_go.svg" alt="Meilisearch-Go" width="200" height="200" />
</p>

<h1 align="center">Meilisearch Go</h1>

<h4 align="center">
  <a href="https://github.com/meilisearch/meilisearch">Meilisearch</a> |
<a href="https://www.meilisearch.com/cloud?utm_campaign=oss&utm_source=github&utm_medium=meilisearch-go">Meilisearch Cloud</a> |
  <a href="https://www.meilisearch.com/docs">Documentation</a> |
  <a href="https://discord.meilisearch.com">Discord</a> |
  <a href="https://roadmap.meilisearch.com/tabs/1-under-consideration">Roadmap</a> |
  <a href="https://www.meilisearch.com">Website</a> |
  <a href="https://www.meilisearch.com/docs/faq">FAQ</a>
</h4>

<p align="center">
  <a href="https://github.com/meilisearch/meilisearch-go/actions"><img src="https://github.com/meilisearch/meilisearch-go/workflows/Tests/badge.svg" alt="GitHub Workflow Status"></a>
  <a href="https://goreportcard.com/report/github.com/meilisearch/meilisearch-go"><img src="https://goreportcard.com/badge/github.com/meilisearch/meilisearch-go" alt="Test"></a>
  <a href="https://codecov.io/gh/meilisearch/meilisearch-go"><img src="https://codecov.io/gh/meilisearch/meilisearch-go/branch/main/graph/badge.svg?token=8N6N60D5UI" alt="CodeCov"></a>
  <a href="https://pkg.go.dev/github.com/meilisearch/meilisearch-go"><img src="https://pkg.go.dev/badge/github.com/meilisearch/meilisearch-go.svg" alt="Go Reference"></a>
  <a href="https://github.com/meilisearch/meilisearch-go/blob/main/LICENSE"><img src="https://img.shields.io/badge/license-MIT-informational" alt="License"></a>
  <a href="https://ms-bors.herokuapp.com/repositories/58"><img src="https://bors.tech/images/badge_small.svg" alt="Bors enabled"></a>
</p>

<p align="center">⚡ The Meilisearch API client written for Golang</p>

**Meilisearch Go** is the Meilisearch API client for Go developers.

**Meilisearch** is an open-source search engine. [Learn more about Meilisearch.](https://github.com/meilisearch/Meilisearch)

## Table of Contents

- [📖 Documentation](#-documentation)
- [🔧 Installation](#-installation)
- [🚀 Getting started](#-getting-started)
  - [Add documents](#add-documents)
  - [Basic search](#basic-search)
  - [Custom search](#custom-search)
  - [Custom search with filter](#custom-search-with-filters)
  - [Customize client](#customize-client)
- [🤖 Compatibility with Meilisearch](#-compatibility-with-meilisearch)
- [⚡️ Benchmark performance](#-benchmark-performance)
- [💡 Learn more](#-learn-more)
- [⚙️ Contributing](#️-contributing)

## 📖 Documentation

This readme contains all the documentation you need to start using this Meilisearch SDK.

For general information on how to use Meilisearch—such as our API reference, tutorials, guides, and in-depth articles—refer to our [main documentation website](https://www.meilisearch.com/docs/).


## 🔧 Installation

With `go get` in command line:
```bash
go get github.com/meilisearch/meilisearch-go
```

### Run Meilisearch <!-- omit in toc -->

⚡️ **Launch, scale, and streamline in minutes with Meilisearch Cloud**—no maintenance, no commitment, cancel anytime. [Try it free now](https://cloud.meilisearch.com/login?utm_campaign=oss&utm_source=github&utm_medium=meilisearch-go).

🪨  Prefer to self-host? [Download and deploy](https://www.meilisearch.com/docs/learn/self_hosted/getting_started_with_self_hosted_meilisearch?utm_campaign=oss&utm_source=github&utm_medium=meilisearch-go) our fast, open-source search engine on your own infrastructure.

## 🚀 Getting started

#### Add documents

```go
package main

import (
	"fmt"
	"os"

	"github.com/meilisearch/meilisearch-go"
)

func main() {
	client := meilisearch.New("http://localhost:7700", meilisearch.WithAPIKey("foobar"))

	// An index is where the documents are stored.
	index := client.Index("movies")

	// If the index 'movies' does not exist, Meilisearch creates it when you first add the documents.
	documents := []map[string]interface{}{
        { "id": 1, "title": "Carol", "genres": []string{"Romance", "Drama"} },
        { "id": 2, "title": "Wonder Woman", "genres": []string{"Action", "Adventure"} },
        { "id": 3, "title": "Life of Pi", "genres": []string{"Adventure", "Drama"} },
        { "id": 4, "title": "Mad Max: Fury Road", "genres": []string{"Adventure", "Science Fiction"} },
        { "id": 5, "title": "Moana", "genres": []string{"Fantasy", "Action"} },
        { "id": 6, "title": "Philadelphia", "genres": []string{"Drama"} },
	}
	task, err := index.AddDocuments(documents)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Println(task.TaskUID)
}
```

With the `taskUID`, you can check the status (`enqueued`, `canceled`, `processing`, `succeeded` or `failed`) of your documents addition using the [task endpoint](https://www.meilisearch.com/docs/reference/api/tasks).

#### Basic Search

```go
package main

import (
    "fmt"
    "os"

    "github.com/meilisearch/meilisearch-go"
)

func main() {
    // Meilisearch is typo-tolerant:
    searchRes, err := client.Index("movies").Search("philoudelphia",
        &meilisearch.SearchRequest{
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
    "id": 6,
    "title": "Philadelphia",
    "genres": ["Drama"]
  }],
  "offset": 0,
  "limit": 10,
  "processingTimeMs": 1,
  "query": "philoudelphia"
}
```

#### Custom Search

All the supported options are described in the [search parameters](https://www.meilisearch.com/docs/reference/api/search#search-parameters) section of the documentation.

```go
func main() {
    searchRes, err := client.Index("movies").Search("wonder",
        &meilisearch.SearchRequest{
            AttributesToHighlight: []string{"*"},
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
    "hits": [
        {
            "id": 2,
            "title": "Wonder Woman",
            "genres": ["Action", "Adventure"],
            "_formatted": {
                "id": 2,
                "title": "<em>Wonder</em> Woman"
            }
        }
    ],
    "offset": 0,
    "limit": 20,
    "processingTimeMs": 0,
    "query": "wonder"
}
```

#### Custom Search With Filters

If you want to enable filtering, you must add your attributes to the `filterableAttributes` index setting.

```go
task, err := index.UpdateFilterableAttributes(&[]string{"id", "genres"})
```

You only need to perform this operation once.

Note that Meilisearch will rebuild your index whenever you update `filterableAttributes`. Depending on the size of your dataset, this might take time. You can track the process using the [task status](https://www.meilisearch.com/docs/learn/advanced/asynchronous_operations).

Then, you can perform the search:

```go
searchRes, err := index.Search("wonder",
    &meilisearch.SearchRequest{
        Filter: "id > 1 AND genres = Action",
    })
```

```json
{
  "hits": [
    {
      "id": 2,
      "title": "Wonder Woman",
      "genres": ["Action","Adventure"]
    }
  ],
  "offset": 0,
  "limit": 20,
  "estimatedTotalHits": 1,
  "processingTimeMs": 0,
  "query": "wonder"
}
```

#### Customize Client

The client supports many customization options:

- `WithCustomClient` sets a custom `http.Client`.
- `WithCustomClientWithTLS` enables TLS for the HTTP client.
- `WithAPIKey` sets the API key or master [key](https://www.meilisearch.com/docs/reference/api/keys).
- `WithContentEncoding` configures [content encoding](https://www.meilisearch.com/docs/reference/api/overview#content-encoding) for requests and responses. Currently, gzip, deflate, and brotli are supported.
- `WithCustomRetries` customizes retry behavior based on specific HTTP status codes (`retryOnStatus`, defaults to 502, 503, and 504) and allows setting the maximum number of retries.
- `DisableRetries` disables the retry logic. By default, retries are enabled.

```go
package main

import (
    "net/http"
    "github.com/meilisearch/meilisearch-go"
)

func main() {
	client := meilisearch.New("http://localhost:7700",
        meilisearch.WithAPIKey("foobar"),
        meilisearch.WithCustomClient(http.DefaultClient),
        meilisearch.WithContentEncoding(meilisearch.GzipEncoding, meilisearch.BestCompression),
        meilisearch.WithCustomRetries([]int{502}, 20),
    )
}
```

## 🤖 Compatibility with Meilisearch

This package guarantees compatibility with [version v1.x of Meilisearch](https://github.com/meilisearch/meilisearch/releases/latest), but some features may not be present. Please check the [issues](https://github.com/meilisearch/meilisearch-go/issues?q=is%3Aissue+is%3Aopen+label%3A%22good+first+issue%22+label%3Aenhancement) for more info.

## ⚡️ Benchmark Performance

The Meilisearch client performance was tested in [client_bench_test.go](/client_bench_test.go).

```shell
goos: linux
goarch: amd64
pkg: github.com/meilisearch/meilisearch-go
cpu: AMD Ryzen 7 5700U with Radeon Graphics
```

**Results**

```shell
Benchmark_ExecuteRequest-16                  	   10000	    105880 ns/op	    7241 B/op	      87 allocs/op
Benchmark_ExecuteRequestWithEncoding-16      	    2716	    455548 ns/op	 1041998 B/op	     169 allocs/op
Benchmark_ExecuteRequestWithoutRetries-16    	       1	3002787257 ns/op	   56528 B/op	     332 allocs/op
```

## 💡 Learn more

The following sections in our main documentation website may interest you:

- **Manipulate documents**: see the [API references](https://www.meilisearch.com/docs/reference/api/documents) or read more about [documents](https://www.meilisearch.com/docs/learn/core_concepts/documents).
- **Search**: see the [API references](https://www.meilisearch.com/docs/reference/api/search) or follow our guide on [search parameters](https://www.meilisearch.com/docs/reference/api/search#search-parameters).
- **Manage the indexes**: see the [API references](https://www.meilisearch.com/docs/reference/api/indexes) or read more about [indexes](https://www.meilisearch.com/docs/learn/core_concepts/indexes).
- **ClientConfigure the index settings**: see the [API references](https://www.meilisearch.com/docs/reference/api/settings) or follow our guide on [settings parameters](https://www.meilisearch.com/docs/reference/api/settings#settings_parameters).

## ⚙️ Contributing

Any new contribution is more than welcome in this project!

If you want to know more about the development workflow or want to contribute, please visit our [contributing guidelines](/CONTRIBUTING.md) for detailed instructions!

<hr>

**Meilisearch** provides and maintains many **SDKs and Integration tools** like this one. We want to provide everyone with an **amazing search experience for any kind of project**. If you want to contribute, make suggestions, or just know what's going on right now, visit us in the [integration-guides](https://github.com/meilisearch/integration-guides) repository.
