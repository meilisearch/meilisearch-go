package meilisearch

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSearchRequest_InlineRule_JSON(t *testing.T) {
	t.Run("omitted when nil", func(t *testing.T) {
		data, err := json.Marshal(&SearchRequest{Query: "x"})
		require.NoError(t, err)
		require.NotContains(t, string(data), "inlineRule")
	})

	t.Run("builder sets inline rule", func(t *testing.T) {
		r := &SearchRequestInlineRule{Uid: "r"}
		req := NewSearchRequestBuilder().WithInlineRule(r).Build()
		require.Same(t, r, req.InlineRule)
	})
}

func TestIndex_Search_InlineRule(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, "/indexes/movies/search", r.URL.Path)
		body, _ := io.ReadAll(r.Body)
		var got map[string]any
		require.NoError(t, json.Unmarshal(body, &got))
		require.Equal(t, "user-likes-batman", got["inlineRule"].(map[string]any)["uid"])
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"hits":[],"query":"super hero","processingTimeMs":1}`))
	}))
	defer srv.Close()

	resp, err := New(srv.URL).Index("movies").Search(context.Background(), "super hero", &SearchRequest{
		InlineRule: &SearchRequestInlineRule{
			Uid:     "user-likes-batman",
			Actions: &Actions{Scale: []Scale{{Filter: "hero = batman", Weight: 1.2}}},
		},
	})
	require.NoError(t, err)
	require.Nil(t, resp.Metadata)
}

// headerTransport adds a fixed header to every request.
type headerTransport struct {
	key, value string
	base       http.RoundTripper
}

func (h *headerTransport) RoundTrip(r *http.Request) (*http.Response, error) {
	r = r.Clone(r.Context())
	r.Header.Set(h.key, h.value)
	return h.base.RoundTrip(r)
}

func TestMultiSearch_InlineRule_ConditionOutcomes(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, "/multi-search", r.URL.Path)
		require.Equal(t, "true", r.Header.Get("Meili-Include-Metadata"))
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"results":[{"hits":[]},{"hits":[],"metadata":{
			"queryUid":"q1","indexUid":"products",
			"inlineRuleConditionOutcomes":{
				"satisfiesActiveCondition":true,
				"satisfiesQueryWordsCondition":{"missingWord":{"word":"pants"}}}}}]}`))
	}))
	defer srv.Close()

	cli := New(srv.URL, WithCustomClient(&http.Client{Transport: &headerTransport{
		key: "Meili-Include-Metadata", value: "true", base: http.DefaultTransport,
	}}))
	resp, err := cli.MultiSearch(context.Background(), &MultiSearchRequest{Queries: []*SearchRequest{
		{IndexUID: "products", Query: "fashion"},
		{IndexUID: "products", Query: "fashion", InlineRule: &SearchRequestInlineRule{
			Uid:        "special-sales",
			Conditions: &Conditions{Query: &QueryCondition{Words: "pants"}},
		}},
	}})
	require.NoError(t, err)
	require.Len(t, resp.Results, 2)
	require.Nil(t, resp.Results[0].Metadata)

	md := resp.Results[1].Metadata
	require.NotNil(t, md)
	require.Equal(t, "products", md.IndexUID)
	require.True(t, *md.InlineRuleConditionOutcomes.SatisfiesActiveCondition)
	require.JSONEq(t, `{"missingWord":{"word":"pants"}}`, string(md.InlineRuleConditionOutcomes.SatisfiesQueryWordsCondition))
}