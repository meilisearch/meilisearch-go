## Problem
Currently, the `FacetSearchRequest` struct defines the `Filter` field as a `string`.

```go
type FacetSearchRequest struct {
    // ...
    Filter string `json:"filter,omitempty"`
}
```

However, the Meilisearch API documentation (and behavior) allows `filter` to be either a **string** or an **array of strings** (allowing for complex AND/OR logic).

The `SearchRequest` struct already correctly handles this by using `interface{}`:
```go
type SearchRequest struct {
    // ...
    Filter interface{} `json:"filter,omitempty"`
}
```

## Solution
Change `FacetSearchRequest.Filter` type from `string` to `interface{}` to match `SearchRequest` and support slice input.
