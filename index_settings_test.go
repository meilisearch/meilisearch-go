package meilisearch

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestIndex_GetFilterableAttributes(t *testing.T) {
	type args struct {
		UID    string
		client *Client
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "TestIndexBasicGetFilterableAttributes",
			args: args{
				UID:    "indexUID",
				client: defaultClient,
			},
		},
		{
			name: "TestIndexGetFilterableAttributesWithCustomClient",
			args: args{
				UID:    "indexUID",
				client: customClient,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			SetUpIndexForFaceting()
			c := tt.args.client
			i := c.Index(tt.args.UID)

			gotResp, err := i.GetFilterableAttributes()
			require.NoError(t, err)
			require.Empty(t, gotResp)
		})
	}
}

func TestIndex_GetDisplayedAttributes(t *testing.T) {
	type args struct {
		UID    string
		client *Client
	}
	tests := []struct {
		name     string
		args     args
		wantResp *[]string
	}{
		{
			name: "TestIndexBasicGetDisplayedAttributes",
			args: args{
				UID:    "indexUID",
				client: defaultClient,
			},
			wantResp: &[]string{"*"},
		},
		{
			name: "TestIndexGetDisplayedAttributesWithCustomClient",
			args: args{
				UID:    "indexUID",
				client: customClient,
			},
			wantResp: &[]string{"*"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			SetUpIndexForFaceting()
			c := tt.args.client
			i := c.Index(tt.args.UID)

			gotResp, err := i.GetDisplayedAttributes()
			require.NoError(t, err)
			require.Equal(t, tt.wantResp, gotResp)
			deleteAllIndexes(c)
		})
	}
}

func TestIndex_GetDistinctAttribute(t *testing.T) {
	type args struct {
		UID    string
		client *Client
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "TestIndexBasicGetDistinctAttribute",
			args: args{
				UID:    "indexUID",
				client: defaultClient,
			},
		},
		{
			name: "TestIndexBasicGetDistinctAttribute",
			args: args{
				UID:    "indexUID",
				client: defaultClient,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			SetUpIndexForFaceting()
			c := tt.args.client
			i := c.Index(tt.args.UID)

			gotResp, err := i.GetDistinctAttribute()
			require.NoError(t, err)
			require.Empty(t, gotResp)
			deleteAllIndexes(c)
		})
	}
}

func TestIndex_GetRankingRules(t *testing.T) {
	type args struct {
		UID    string
		client *Client
	}
	tests := []struct {
		name     string
		args     args
		wantResp *[]string
	}{
		{
			name: "TestIndexBasicGetRankingRules",
			args: args{
				UID:    "indexUID",
				client: defaultClient,
			},
			wantResp: &[]string{"words", "typo", "proximity", "attribute", "exactness"},
		},
		{
			name: "TestIndexGetRankingRulesWithCustomClient",
			args: args{
				UID:    "indexUID",
				client: defaultClient,
			},
			wantResp: &[]string{"words", "typo", "proximity", "attribute", "exactness"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			SetUpIndexForFaceting()
			c := tt.args.client
			i := c.Index(tt.args.UID)

			gotResp, err := i.GetRankingRules()
			require.NoError(t, err)
			require.Equal(t, tt.wantResp, gotResp)
			deleteAllIndexes(c)
		})
	}
}

func TestIndex_GetSearchableAttributes(t *testing.T) {
	type args struct {
		UID    string
		client *Client
	}
	tests := []struct {
		name     string
		args     args
		wantResp *[]string
	}{
		{
			name: "TestIndexBasicGetSearchableAttributes",
			args: args{
				UID:    "indexUID",
				client: defaultClient,
			},
			wantResp: &[]string{"*"},
		},
		{
			name: "TestIndexGetSearchableAttributesCustomClient",
			args: args{
				UID:    "indexUID",
				client: defaultClient,
			},
			wantResp: &[]string{"*"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			SetUpIndexForFaceting()
			c := tt.args.client
			i := c.Index(tt.args.UID)

			gotResp, err := i.GetSearchableAttributes()
			require.NoError(t, err)
			require.Equal(t, tt.wantResp, gotResp)
			deleteAllIndexes(c)
		})
	}
}

func TestIndex_GetSettings(t *testing.T) {
	type args struct {
		UID    string
		client *Client
	}
	tests := []struct {
		name     string
		args     args
		wantResp *Settings
	}{
		{
			name: "TestIndexBasicGetSettings",
			args: args{
				UID:    "indexUID",
				client: defaultClient,
			},
			wantResp: &Settings{
				RankingRules: []string{
					"words", "typo", "proximity", "attribute", "exactness",
				},
				DistinctAttribute:    (*string)(nil),
				SearchableAttributes: []string{"*"},
				DisplayedAttributes:  []string{"*"},
				StopWords:            []string{},
				Synonyms:             map[string][]string(nil),
				FilterableAttributes: []string{},
			},
		},
		{
			name: "TestIndexGetSettingsWithCustomClient",
			args: args{
				UID:    "indexUID",
				client: customClient,
			},
			wantResp: &Settings{
				RankingRules: []string{
					"words", "typo", "proximity", "attribute", "exactness",
				},
				DistinctAttribute:    (*string)(nil),
				SearchableAttributes: []string{"*"},
				DisplayedAttributes:  []string{"*"},
				StopWords:            []string{},
				Synonyms:             map[string][]string(nil),
				FilterableAttributes: []string{},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			SetUpIndexForFaceting()
			c := tt.args.client
			i := c.Index(tt.args.UID)

			gotResp, err := i.GetSettings()
			require.NoError(t, err)
			require.Equal(t, tt.wantResp, gotResp)
			deleteAllIndexes(c)
		})
	}
}

func TestIndex_GetStopWords(t *testing.T) {
	type args struct {
		UID    string
		client *Client
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "TestIndexBasicGetStopWords",
			args: args{
				UID:    "indexUID",
				client: defaultClient,
			},
		},
		{
			name: "TestIndexGetStopWordsCustomClient",
			args: args{
				UID:    "indexUID",
				client: customClient,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			SetUpIndexForFaceting()
			c := tt.args.client
			i := c.Index(tt.args.UID)

			gotResp, err := i.GetStopWords()
			require.NoError(t, err)
			require.Empty(t, gotResp)
			deleteAllIndexes(c)
		})
	}
}

func TestIndex_GetSynonyms(t *testing.T) {
	type args struct {
		UID    string
		client *Client
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "TestIndexBasicGetSynonyms",
			args: args{
				UID:    "indexUID",
				client: defaultClient,
			},
		},
		{
			name: "TestIndexGetSynonymsWithCustomClient",
			args: args{
				UID:    "indexUID",
				client: customClient,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			SetUpIndexForFaceting()
			c := tt.args.client
			i := c.Index(tt.args.UID)

			gotResp, err := i.GetSynonyms()
			require.NoError(t, err)
			require.Empty(t, gotResp)
			deleteAllIndexes(c)
		})
	}
}

func TestIndex_ResetFilterableAttributes(t *testing.T) {
	type args struct {
		UID    string
		client *Client
	}
	tests := []struct {
		name       string
		args       args
		wantUpdate *AsyncUpdateID
	}{
		{
			name: "TestIndexBasicResetFilterableAttributes",
			args: args{
				UID:    "indexUID",
				client: defaultClient,
			},
			wantUpdate: &AsyncUpdateID{
				UpdateID: 1,
			},
		},
		{
			name: "TestIndexResetFilterableAttributesCustomClient",
			args: args{
				UID:    "indexUID",
				client: customClient,
			},
			wantUpdate: &AsyncUpdateID{
				UpdateID: 1,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			SetUpIndexForFaceting()
			c := tt.args.client
			i := c.Index(tt.args.UID)

			gotUpdate, err := i.ResetFilterableAttributes()
			require.NoError(t, err)
			require.Equal(t, tt.wantUpdate, gotUpdate)
			i.DefaultWaitForPendingUpdate(gotUpdate)

			gotResp, err := i.GetFilterableAttributes()
			require.NoError(t, err)
			require.Empty(t, gotResp)

			deleteAllIndexes(c)
		})
	}
}

func TestIndex_ResetDisplayedAttributes(t *testing.T) {
	type args struct {
		UID    string
		client *Client
	}
	tests := []struct {
		name       string
		args       args
		wantUpdate *AsyncUpdateID
		wantResp   *[]string
	}{
		{
			name: "TestIndexBasicResetDisplayedAttributes",
			args: args{
				UID:    "indexUID",
				client: defaultClient,
			},
			wantUpdate: &AsyncUpdateID{
				UpdateID: 1,
			},
			wantResp: &[]string{"*"},
		},
		{
			name: "TestIndexResetDisplayedAttributesCustomClient",
			args: args{
				UID:    "indexUID",
				client: customClient,
			},
			wantUpdate: &AsyncUpdateID{
				UpdateID: 1,
			},
			wantResp: &[]string{"*"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			SetUpIndexForFaceting()
			c := tt.args.client
			i := c.Index(tt.args.UID)

			gotUpdate, err := i.ResetDisplayedAttributes()
			require.NoError(t, err)
			require.GreaterOrEqual(t, gotUpdate.UpdateID, tt.wantUpdate.UpdateID)
			i.DefaultWaitForPendingUpdate(gotUpdate)

			gotResp, err := i.GetDisplayedAttributes()
			require.NoError(t, err)
			require.Equal(t, tt.wantResp, gotResp)

			deleteAllIndexes(c)
		})
	}
}

func TestIndex_ResetDistinctAttribute(t *testing.T) {
	type args struct {
		UID    string
		client *Client
	}
	tests := []struct {
		name       string
		args       args
		wantUpdate *AsyncUpdateID
	}{
		{
			name: "TestIndexBasicResetDistinctAttribute",
			args: args{
				UID:    "indexUID",
				client: defaultClient,
			},
			wantUpdate: &AsyncUpdateID{
				UpdateID: 1,
			},
		},
		{
			name: "TestIndexResetDistinctAttributeWithCustomClient",
			args: args{
				UID:    "indexUID",
				client: customClient,
			},
			wantUpdate: &AsyncUpdateID{
				UpdateID: 1,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			SetUpIndexForFaceting()
			c := tt.args.client
			i := c.Index(tt.args.UID)

			gotUpdate, err := i.ResetDistinctAttribute()
			require.NoError(t, err)
			require.GreaterOrEqual(t, gotUpdate.UpdateID, tt.wantUpdate.UpdateID)
			i.DefaultWaitForPendingUpdate(gotUpdate)

			gotResp, err := i.GetDistinctAttribute()
			require.NoError(t, err)
			require.Empty(t, gotResp)
			deleteAllIndexes(c)
		})
	}
}

func TestIndex_ResetRankingRules(t *testing.T) {
	type args struct {
		UID    string
		client *Client
	}
	tests := []struct {
		name       string
		args       args
		wantUpdate *AsyncUpdateID
		wantResp   *[]string
	}{
		{
			name: "TestIndexBasicResetRankingRules",
			args: args{
				UID:    "indexUID",
				client: defaultClient,
			},
			wantUpdate: &AsyncUpdateID{
				UpdateID: 1,
			},
			wantResp: &[]string{"words", "typo", "proximity", "attribute", "exactness"},
		},
		{
			name: "TestIndexResetRankingRulesWithCustomClient",
			args: args{
				UID:    "indexUID",
				client: customClient,
			},
			wantUpdate: &AsyncUpdateID{
				UpdateID: 1,
			},
			wantResp: &[]string{"words", "typo", "proximity", "attribute", "exactness"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			SetUpIndexForFaceting()
			c := tt.args.client
			i := c.Index(tt.args.UID)

			gotUpdate, err := i.ResetRankingRules()
			require.NoError(t, err)
			require.GreaterOrEqual(t, gotUpdate.UpdateID, tt.wantUpdate.UpdateID)
			i.DefaultWaitForPendingUpdate(gotUpdate)

			gotResp, err := i.GetRankingRules()
			require.NoError(t, err)
			require.Equal(t, tt.wantResp, gotResp)
			deleteAllIndexes(c)
		})
	}
}

func TestIndex_ResetSearchableAttributes(t *testing.T) {
	type args struct {
		UID    string
		client *Client
	}
	tests := []struct {
		name       string
		args       args
		wantUpdate *AsyncUpdateID
		wantResp   *[]string
	}{
		{
			name: "TestIndexBasicResetSearchableAttributes",
			args: args{
				UID:    "indexUID",
				client: defaultClient,
			},
			wantUpdate: &AsyncUpdateID{
				UpdateID: 1,
			},
			wantResp: &[]string{"*"},
		},
		{
			name: "TestIndexResetSearchableAttributesWithCustomClient",
			args: args{
				UID:    "indexUID",
				client: customClient,
			},
			wantUpdate: &AsyncUpdateID{
				UpdateID: 1,
			},
			wantResp: &[]string{"*"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			SetUpIndexForFaceting()
			c := tt.args.client
			i := c.Index(tt.args.UID)

			gotUpdate, err := i.ResetSearchableAttributes()
			require.NoError(t, err)
			require.GreaterOrEqual(t, gotUpdate.UpdateID, tt.wantUpdate.UpdateID)
			i.DefaultWaitForPendingUpdate(gotUpdate)

			gotResp, err := i.GetSearchableAttributes()
			require.NoError(t, err)
			require.Equal(t, tt.wantResp, gotResp)
			deleteAllIndexes(c)
		})
	}
}

func TestIndex_ResetSettings(t *testing.T) {
	type args struct {
		UID    string
		client *Client
	}
	tests := []struct {
		name       string
		args       args
		wantUpdate *AsyncUpdateID
		wantResp   *Settings
	}{
		{
			name: "TestIndexBasicResetSettings",
			args: args{
				UID:    "indexUID",
				client: defaultClient,
			},
			wantUpdate: &AsyncUpdateID{
				UpdateID: 1,
			},
			wantResp: &Settings{
				RankingRules: []string{
					"words", "typo", "proximity", "attribute", "exactness",
				},
				DistinctAttribute:    (*string)(nil),
				SearchableAttributes: []string{"*"},
				DisplayedAttributes:  []string{"*"},
				StopWords:            []string{},
				Synonyms:             map[string][]string(nil),
				FilterableAttributes: []string{},
			},
		},
		{
			name: "TestIndexResetSettingsWithCustomClient",
			args: args{
				UID:    "indexUID",
				client: customClient,
			},
			wantUpdate: &AsyncUpdateID{
				UpdateID: 1,
			},
			wantResp: &Settings{
				RankingRules: []string{
					"words", "typo", "proximity", "attribute", "exactness",
				},
				DistinctAttribute:    (*string)(nil),
				SearchableAttributes: []string{"*"},
				DisplayedAttributes:  []string{"*"},
				StopWords:            []string{},
				Synonyms:             map[string][]string(nil),
				FilterableAttributes: []string{},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			SetUpIndexForFaceting()
			c := tt.args.client
			i := c.Index(tt.args.UID)

			gotUpdate, err := i.ResetSettings()
			require.NoError(t, err)
			require.GreaterOrEqual(t, gotUpdate.UpdateID, tt.wantUpdate.UpdateID)
			i.DefaultWaitForPendingUpdate(gotUpdate)

			gotResp, err := i.GetSettings()
			require.NoError(t, err)
			require.Equal(t, tt.wantResp, gotResp)
			deleteAllIndexes(c)
		})
	}
}

func TestIndex_ResetStopWords(t *testing.T) {
	type args struct {
		UID    string
		client *Client
	}
	tests := []struct {
		name       string
		args       args
		wantUpdate *AsyncUpdateID
	}{
		{
			name: "TestIndexBasicResetStopWords",
			args: args{
				UID:    "indexUID",
				client: defaultClient,
			},
			wantUpdate: &AsyncUpdateID{
				UpdateID: 1,
			},
		},
		{
			name: "TestIndexResetStopWordsWithCustomClient",
			args: args{
				UID:    "indexUID",
				client: customClient,
			},
			wantUpdate: &AsyncUpdateID{
				UpdateID: 1,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			SetUpIndexForFaceting()
			c := tt.args.client
			i := c.Index(tt.args.UID)

			gotUpdate, err := i.ResetStopWords()
			require.NoError(t, err)
			require.GreaterOrEqual(t, gotUpdate.UpdateID, tt.wantUpdate.UpdateID)
			i.DefaultWaitForPendingUpdate(gotUpdate)

			gotResp, err := i.GetStopWords()
			require.NoError(t, err)
			require.Empty(t, gotResp)
			deleteAllIndexes(c)
		})
	}
}

func TestIndex_ResetSynonyms(t *testing.T) {
	type args struct {
		UID    string
		client *Client
	}
	tests := []struct {
		name       string
		args       args
		wantUpdate *AsyncUpdateID
	}{
		{
			name: "TestIndexBasicResetSynonyms",
			args: args{
				UID:    "indexUID",
				client: defaultClient,
			},
			wantUpdate: &AsyncUpdateID{
				UpdateID: 1,
			},
		},
		{
			name: "TestIndexResetSynonymsWithCustomClient",
			args: args{
				UID:    "indexUID",
				client: customClient,
			},
			wantUpdate: &AsyncUpdateID{
				UpdateID: 1,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			SetUpIndexForFaceting()
			c := tt.args.client
			i := c.Index(tt.args.UID)

			gotUpdate, err := i.ResetSynonyms()
			require.NoError(t, err)
			require.GreaterOrEqual(t, gotUpdate.UpdateID, tt.wantUpdate.UpdateID)
			i.DefaultWaitForPendingUpdate(gotUpdate)

			gotResp, err := i.GetSynonyms()
			require.NoError(t, err)
			require.Empty(t, gotResp)
			deleteAllIndexes(c)
		})
	}
}

func TestIndex_UpdateFilterableAttributes(t *testing.T) {
	type args struct {
		UID     string
		client  *Client
		request []string
	}
	tests := []struct {
		name       string
		args       args
		wantUpdate *AsyncUpdateID
	}{
		{
			name: "TestIndexBasicUpdateFilterableAttributes",
			args: args{
				UID:    "indexUID",
				client: defaultClient,
				request: []string{
					"title",
				},
			},
			wantUpdate: &AsyncUpdateID{
				UpdateID: 1,
			},
		},
		{
			name: "TestIndexUpdateFilterableAttributesWithCustomClient",
			args: args{
				UID:    "indexUID",
				client: customClient,
				request: []string{
					"title",
				},
			},
			wantUpdate: &AsyncUpdateID{
				UpdateID: 1,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			SetUpIndexForFaceting()
			c := tt.args.client
			i := c.Index(tt.args.UID)

			gotResp, err := i.GetFilterableAttributes()
			require.NoError(t, err)
			require.Empty(t, gotResp)

			gotUpdate, err := i.UpdateFilterableAttributes(&tt.args.request)
			require.NoError(t, err)
			require.Equal(t, gotUpdate.UpdateID, tt.wantUpdate.UpdateID)
			i.DefaultWaitForPendingUpdate(gotUpdate)

			gotResp, err = i.GetFilterableAttributes()
			require.NoError(t, err)
			require.Equal(t, &tt.args.request, gotResp)
			deleteAllIndexes(c)
		})
	}
}

func TestIndex_UpdateDisplayedAttributes(t *testing.T) {
	type args struct {
		UID     string
		client  *Client
		request []string
	}
	tests := []struct {
		name       string
		args       args
		wantUpdate *AsyncUpdateID
		wantResp   *[]string
	}{
		{
			name: "TestIndexBasicUpdateDisplayedAttributes",
			args: args{
				UID:    "indexUID",
				client: defaultClient,
				request: []string{
					"book_id", "tag", "title",
				},
			},
			wantUpdate: &AsyncUpdateID{
				UpdateID: 1,
			},
			wantResp: &[]string{"*"},
		},
		{
			name: "TestIndexUpdateDisplayedAttributesWithCustomClient",
			args: args{
				UID:    "indexUID",
				client: customClient,
				request: []string{
					"book_id", "tag", "title",
				},
			},
			wantUpdate: &AsyncUpdateID{
				UpdateID: 1,
			},
			wantResp: &[]string{"*"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			SetUpIndexForFaceting()
			c := tt.args.client
			i := c.Index(tt.args.UID)

			gotResp, err := i.GetDisplayedAttributes()
			require.NoError(t, err)
			require.Equal(t, tt.wantResp, gotResp)

			gotUpdate, err := i.UpdateDisplayedAttributes(&tt.args.request)
			require.NoError(t, err)
			require.GreaterOrEqual(t, gotUpdate.UpdateID, tt.wantUpdate.UpdateID)
			i.DefaultWaitForPendingUpdate(gotUpdate)

			gotResp, err = i.GetDisplayedAttributes()
			require.NoError(t, err)
			require.Equal(t, &tt.args.request, gotResp)
			deleteAllIndexes(c)
		})
	}
}

func TestIndex_UpdateDistinctAttribute(t *testing.T) {
	type args struct {
		UID     string
		client  *Client
		request string
	}
	tests := []struct {
		name       string
		args       args
		wantUpdate *AsyncUpdateID
	}{
		{
			name: "TestIndexBasicUpdateDistinctAttribute",
			args: args{
				UID:     "indexUID",
				client:  defaultClient,
				request: "movie_id",
			},
			wantUpdate: &AsyncUpdateID{
				UpdateID: 1,
			},
		},
		{
			name: "TestIndexUpdateDistinctAttributeWithCustomClient",
			args: args{
				UID:     "indexUID",
				client:  customClient,
				request: "movie_id",
			},
			wantUpdate: &AsyncUpdateID{
				UpdateID: 1,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			SetUpIndexForFaceting()
			c := tt.args.client
			i := c.Index(tt.args.UID)

			gotResp, err := i.GetDistinctAttribute()
			require.NoError(t, err)
			require.Empty(t, gotResp)

			gotUpdate, err := i.UpdateDistinctAttribute(tt.args.request)
			require.NoError(t, err)
			require.GreaterOrEqual(t, gotUpdate.UpdateID, tt.wantUpdate.UpdateID)
			i.DefaultWaitForPendingUpdate(gotUpdate)

			gotResp, err = i.GetDistinctAttribute()
			require.NoError(t, err)
			require.Equal(t, &tt.args.request, gotResp)
			deleteAllIndexes(c)
		})
	}
}

func TestIndex_UpdateRankingRules(t *testing.T) {
	type args struct {
		UID     string
		client  *Client
		request []string
	}
	tests := []struct {
		name       string
		args       args
		wantUpdate *AsyncUpdateID
		wantResp   *[]string
	}{
		{
			name: "TestIndexBasicUpdateRankingRules",
			args: args{
				UID:    "indexUID",
				client: defaultClient,
				request: []string{
					"typo", "words",
				},
			},
			wantUpdate: &AsyncUpdateID{
				UpdateID: 1,
			},
			wantResp: &[]string{"words", "typo", "proximity", "attribute", "exactness"},
		},
		{
			name: "TestIndexUpdateRankingRulesWithCustomClient",
			args: args{
				UID:    "indexUID",
				client: customClient,
				request: []string{
					"typo", "words",
				},
			},
			wantUpdate: &AsyncUpdateID{
				UpdateID: 1,
			},
			wantResp: &[]string{"words", "typo", "proximity", "attribute", "exactness"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			SetUpIndexForFaceting()
			c := tt.args.client
			i := c.Index(tt.args.UID)

			gotResp, err := i.GetRankingRules()
			require.NoError(t, err)
			require.Equal(t, tt.wantResp, gotResp)

			gotUpdate, err := i.UpdateRankingRules(&tt.args.request)
			require.NoError(t, err)
			require.GreaterOrEqual(t, gotUpdate.UpdateID, tt.wantUpdate.UpdateID)
			i.DefaultWaitForPendingUpdate(gotUpdate)

			gotResp, err = i.GetRankingRules()
			require.NoError(t, err)
			require.Equal(t, &tt.args.request, gotResp)
			deleteAllIndexes(c)
		})
	}
}

func TestIndex_UpdateSearchableAttributes(t *testing.T) {
	type args struct {
		UID     string
		client  *Client
		request []string
	}
	tests := []struct {
		name       string
		args       args
		wantUpdate *AsyncUpdateID
		wantResp   *[]string
	}{
		{
			name: "TestIndexBasicUpdateSearchableAttributes",
			args: args{
				UID:    "indexUID",
				client: defaultClient,
				request: []string{
					"title", "tag",
				},
			},
			wantUpdate: &AsyncUpdateID{
				UpdateID: 1,
			},
			wantResp: &[]string{"*"},
		},
		{
			name: "TestIndexUpdateSearchableAttributesWithCustomClient",
			args: args{
				UID:    "indexUID",
				client: customClient,
				request: []string{
					"title", "tag",
				},
			},
			wantUpdate: &AsyncUpdateID{
				UpdateID: 1,
			},
			wantResp: &[]string{"*"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			SetUpIndexForFaceting()
			c := tt.args.client
			i := c.Index(tt.args.UID)

			gotResp, err := i.GetSearchableAttributes()
			require.NoError(t, err)
			require.Equal(t, tt.wantResp, gotResp)

			gotUpdate, err := i.UpdateSearchableAttributes(&tt.args.request)
			require.NoError(t, err)
			require.GreaterOrEqual(t, gotUpdate.UpdateID, tt.wantUpdate.UpdateID)
			i.DefaultWaitForPendingUpdate(gotUpdate)

			gotResp, err = i.GetSearchableAttributes()
			require.NoError(t, err)
			require.Equal(t, &tt.args.request, gotResp)
			deleteAllIndexes(c)
		})
	}
}

func TestIndex_UpdateSettings(t *testing.T) {
	type args struct {
		UID     string
		client  *Client
		request Settings
	}
	tests := []struct {
		name       string
		args       args
		wantUpdate *AsyncUpdateID
		wantResp   *Settings
	}{
		{
			name: "TestIndexBasicUpdateSettings",
			args: args{
				UID:    "indexUID",
				client: defaultClient,
				request: Settings{
					RankingRules: []string{
						"typo", "words",
					},
					DistinctAttribute: (*string)(nil),
					SearchableAttributes: []string{
						"title", "tag",
					},
					DisplayedAttributes: []string{
						"book_id", "tag", "title",
					},
					StopWords: []string{
						"of", "the",
					},
					Synonyms: map[string][]string{
						"hp": {"harry potter"},
					},
					FilterableAttributes: []string{
						"title",
					},
				},
			},
			wantUpdate: &AsyncUpdateID{
				UpdateID: 1,
			},
			wantResp: &Settings{
				RankingRules: []string{
					"words", "typo", "proximity", "attribute", "exactness",
				},
				DistinctAttribute:    (*string)(nil),
				SearchableAttributes: []string{"*"},
				DisplayedAttributes:  []string{"*"},
				StopWords:            []string{},
				Synonyms:             map[string][]string(nil),
				FilterableAttributes: []string{},
			},
		},
		{
			name: "TestIndexUpdateSettingsWithCustomClient",
			args: args{
				UID:    "indexUID",
				client: customClient,
				request: Settings{
					RankingRules: []string{
						"typo", "words",
					},
					DistinctAttribute: (*string)(nil),
					SearchableAttributes: []string{
						"title", "tag",
					},
					DisplayedAttributes: []string{
						"book_id", "tag", "title",
					},
					StopWords: []string{
						"of", "the",
					},
					Synonyms: map[string][]string{
						"hp": {"harry potter"},
					},
					FilterableAttributes: []string{
						"title",
					},
				},
			},
			wantUpdate: &AsyncUpdateID{
				UpdateID: 1,
			},
			wantResp: &Settings{
				RankingRules: []string{
					"words", "typo", "proximity", "attribute", "exactness",
				},
				DistinctAttribute:    (*string)(nil),
				SearchableAttributes: []string{"*"},
				DisplayedAttributes:  []string{"*"},
				StopWords:            []string{},
				Synonyms:             map[string][]string(nil),
				FilterableAttributes: []string{},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			SetUpIndexForFaceting()
			c := tt.args.client
			i := c.Index(tt.args.UID)

			gotResp, err := i.GetSettings()
			require.NoError(t, err)
			require.Equal(t, tt.wantResp, gotResp)

			gotUpdate, err := i.UpdateSettings(&tt.args.request)
			require.NoError(t, err)
			require.GreaterOrEqual(t, gotUpdate.UpdateID, tt.wantUpdate.UpdateID)
			i.DefaultWaitForPendingUpdate(gotUpdate)

			gotResp, err = i.GetSettings()
			require.NoError(t, err)
			require.Equal(t, &tt.args.request, gotResp)
			deleteAllIndexes(c)
		})
	}
}

func TestIndex_UpdateSettingsOneByOne(t *testing.T) {
	type args struct {
		UID            string
		client         *Client
		firstRequest   Settings
		firstResponse  Settings
		secondRequest  Settings
		secondResponse Settings
	}
	tests := []struct {
		name       string
		args       args
		wantUpdate *AsyncUpdateID
		wantResp   *Settings
	}{
		{
			name: "TestIndexUpdateJustSynonyms",
			args: args{
				UID:    "indexUID",
				client: defaultClient,
				firstRequest: Settings{
					RankingRules: []string{
						"typo", "words",
					},
					Synonyms: map[string][]string{
						"hp": {"harry potter"},
					},
				},
				firstResponse: Settings{
					RankingRules: []string{
						"typo", "words",
					},
					DistinctAttribute:    (*string)(nil),
					SearchableAttributes: []string{"*"},
					DisplayedAttributes:  []string{"*"},
					StopWords:            []string{},
					Synonyms: map[string][]string{
						"hp": {"harry potter"},
					},
					FilterableAttributes: []string{},
				},
				secondRequest: Settings{
					Synonyms: map[string][]string{
						"al": {"alice"},
					},
				},
				secondResponse: Settings{
					RankingRules: []string{
						"typo", "words",
					},
					DistinctAttribute:    (*string)(nil),
					SearchableAttributes: []string{"*"},
					DisplayedAttributes:  []string{"*"},
					StopWords:            []string{},
					Synonyms: map[string][]string{
						"al": {"alice"},
					},
					FilterableAttributes: []string{},
				},
			},
			wantUpdate: &AsyncUpdateID{
				UpdateID: 1,
			},
			wantResp: &Settings{
				RankingRules: []string{
					"words", "typo", "proximity", "attribute", "exactness",
				},
				DistinctAttribute:    (*string)(nil),
				SearchableAttributes: []string{"*"},
				DisplayedAttributes:  []string{"*"},
				StopWords:            []string{},
				Synonyms:             map[string][]string(nil),
				FilterableAttributes: []string{},
			},
		},
		{
			name: "TestIndexUpdateJustSynonymsWithCustomClient",
			args: args{
				UID:    "indexUID",
				client: customClient,
				firstRequest: Settings{
					RankingRules: []string{
						"typo", "words",
					},
					Synonyms: map[string][]string{
						"hp": {"harry potter"},
					},
				},
				firstResponse: Settings{
					RankingRules: []string{
						"typo", "words",
					},
					DistinctAttribute:    (*string)(nil),
					SearchableAttributes: []string{"*"},
					DisplayedAttributes:  []string{"*"},
					StopWords:            []string{},
					Synonyms: map[string][]string{
						"hp": {"harry potter"},
					},
					FilterableAttributes: []string{},
				},
				secondRequest: Settings{
					Synonyms: map[string][]string{
						"al": {"alice"},
					},
				},
				secondResponse: Settings{
					RankingRules: []string{
						"typo", "words",
					},
					DistinctAttribute:    (*string)(nil),
					SearchableAttributes: []string{"*"},
					DisplayedAttributes:  []string{"*"},
					StopWords:            []string{},
					Synonyms: map[string][]string{
						"al": {"alice"},
					},
					FilterableAttributes: []string{},
				},
			},
			wantUpdate: &AsyncUpdateID{
				UpdateID: 1,
			},
			wantResp: &Settings{
				RankingRules: []string{
					"words", "typo", "proximity", "attribute", "exactness",
				},
				DistinctAttribute:    (*string)(nil),
				SearchableAttributes: []string{"*"},
				DisplayedAttributes:  []string{"*"},
				StopWords:            []string{},
				Synonyms:             map[string][]string(nil),
				FilterableAttributes: []string{},
			},
		},
		{
			name: "TestIndexUpdateJustSearchableAttributes",
			args: args{
				UID:    "indexUID",
				client: defaultClient,
				firstRequest: Settings{
					RankingRules: []string{
						"typo", "words",
					},
					SearchableAttributes: []string{
						"tag",
					},
				},
				firstResponse: Settings{
					RankingRules: []string{
						"typo", "words",
					},
					DistinctAttribute: (*string)(nil),
					SearchableAttributes: []string{
						"tag",
					},
					DisplayedAttributes:  []string{"*"},
					StopWords:            []string{},
					Synonyms:             map[string][]string(nil),
					FilterableAttributes: []string{},
				},
				secondRequest: Settings{
					SearchableAttributes: []string{
						"title",
					},
				},
				secondResponse: Settings{
					RankingRules: []string{
						"typo", "words",
					},
					DistinctAttribute: (*string)(nil),
					SearchableAttributes: []string{
						"title",
					},
					DisplayedAttributes:  []string{"*"},
					StopWords:            []string{},
					Synonyms:             map[string][]string(nil),
					FilterableAttributes: []string{},
				},
			},
			wantUpdate: &AsyncUpdateID{
				UpdateID: 1,
			},
			wantResp: &Settings{
				RankingRules: []string{
					"words", "typo", "proximity", "attribute", "exactness",
				},
				DistinctAttribute:    (*string)(nil),
				SearchableAttributes: []string{"*"},
				DisplayedAttributes:  []string{"*"},
				StopWords:            []string{},
				Synonyms:             map[string][]string(nil),
				FilterableAttributes: []string{},
			},
		},
		{
			name: "TestIndexUpdateJustDisplayedAttributes",
			args: args{
				UID:    "indexUID",
				client: defaultClient,
				firstRequest: Settings{
					RankingRules: []string{
						"typo", "words",
					},
					DisplayedAttributes: []string{
						"book_id", "tag", "title",
					},
				},
				firstResponse: Settings{
					RankingRules: []string{
						"typo", "words",
					},
					DistinctAttribute:    (*string)(nil),
					SearchableAttributes: []string{"*"},
					DisplayedAttributes: []string{
						"book_id", "tag", "title",
					},
					StopWords:            []string{},
					Synonyms:             map[string][]string(nil),
					FilterableAttributes: []string{},
				},
				secondRequest: Settings{
					DisplayedAttributes: []string{
						"book_id", "tag",
					},
				},
				secondResponse: Settings{
					RankingRules: []string{
						"typo", "words",
					},
					DistinctAttribute:    (*string)(nil),
					SearchableAttributes: []string{"*"},
					DisplayedAttributes: []string{
						"book_id", "tag",
					},
					StopWords:            []string{},
					Synonyms:             map[string][]string(nil),
					FilterableAttributes: []string{},
				},
			},
			wantUpdate: &AsyncUpdateID{
				UpdateID: 1,
			},
			wantResp: &Settings{
				RankingRules: []string{
					"words", "typo", "proximity", "attribute", "exactness",
				},
				DistinctAttribute:    (*string)(nil),
				SearchableAttributes: []string{"*"},
				DisplayedAttributes:  []string{"*"},
				StopWords:            []string{},
				Synonyms:             map[string][]string(nil),
				FilterableAttributes: []string{},
			},
		},
		{
			name: "TestIndexUpdateJustStopWords",
			args: args{
				UID:    "indexUID",
				client: defaultClient,
				firstRequest: Settings{
					RankingRules: []string{
						"typo", "words",
					},
					StopWords: []string{
						"of", "the",
					},
				},
				firstResponse: Settings{
					RankingRules: []string{
						"typo", "words",
					},
					DistinctAttribute:    (*string)(nil),
					SearchableAttributes: []string{"*"},
					DisplayedAttributes:  []string{"*"},
					StopWords: []string{
						"of", "the",
					},
					Synonyms:             map[string][]string(nil),
					FilterableAttributes: []string{},
				},
				secondRequest: Settings{
					StopWords: []string{
						"of", "the",
					},
				},
				secondResponse: Settings{
					RankingRules: []string{
						"typo", "words",
					},
					DistinctAttribute:    (*string)(nil),
					SearchableAttributes: []string{"*"},
					DisplayedAttributes:  []string{"*"},
					StopWords: []string{
						"of", "the",
					},
					Synonyms:             map[string][]string(nil),
					FilterableAttributes: []string{},
				},
			},
			wantUpdate: &AsyncUpdateID{
				UpdateID: 1,
			},
			wantResp: &Settings{
				RankingRules: []string{
					"words", "typo", "proximity", "attribute", "exactness",
				},
				DistinctAttribute:    (*string)(nil),
				SearchableAttributes: []string{"*"},
				DisplayedAttributes:  []string{"*"},
				StopWords:            []string{},
				Synonyms:             map[string][]string(nil),
				FilterableAttributes: []string{},
			},
		},
		{
			name: "TestIndexUpdateJustFilterableAttributes",
			args: args{
				UID:    "indexUID",
				client: defaultClient,
				firstRequest: Settings{
					RankingRules: []string{
						"typo", "words",
					},
					FilterableAttributes: []string{
						"title",
					},
				},
				firstResponse: Settings{
					RankingRules: []string{
						"typo", "words",
					},
					DistinctAttribute:    (*string)(nil),
					SearchableAttributes: []string{"*"},
					DisplayedAttributes:  []string{"*"},
					StopWords:            []string{},
					Synonyms:             map[string][]string(nil),
					FilterableAttributes: []string{
						"title",
					},
				},
				secondRequest: Settings{
					FilterableAttributes: []string{
						"title",
					},
				},
				secondResponse: Settings{
					RankingRules: []string{
						"typo", "words",
					},
					DistinctAttribute:    (*string)(nil),
					SearchableAttributes: []string{"*"},
					DisplayedAttributes:  []string{"*"},
					StopWords:            []string{},
					Synonyms:             map[string][]string(nil),
					FilterableAttributes: []string{
						"title",
					},
				},
			},
			wantUpdate: &AsyncUpdateID{
				UpdateID: 1,
			},
			wantResp: &Settings{
				RankingRules: []string{
					"words", "typo", "proximity", "attribute", "exactness",
				},
				DistinctAttribute:    (*string)(nil),
				SearchableAttributes: []string{"*"},
				DisplayedAttributes:  []string{"*"},
				StopWords:            []string{},
				Synonyms:             map[string][]string(nil),
				FilterableAttributes: []string{},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			SetUpIndexForFaceting()
			c := tt.args.client
			i := c.Index(tt.args.UID)

			gotResp, err := i.GetSettings()
			require.NoError(t, err)
			require.Equal(t, tt.wantResp, gotResp)

			gotUpdate, err := i.UpdateSettings(&tt.args.firstRequest)
			require.NoError(t, err)
			require.GreaterOrEqual(t, gotUpdate.UpdateID, tt.wantUpdate.UpdateID)
			i.DefaultWaitForPendingUpdate(gotUpdate)

			gotResp, err = i.GetSettings()
			require.NoError(t, err)
			require.Equal(t, &tt.args.firstResponse, gotResp)

			gotUpdate, err = i.UpdateSettings(&tt.args.secondRequest)
			require.NoError(t, err)
			require.GreaterOrEqual(t, gotUpdate.UpdateID, tt.wantUpdate.UpdateID)
			i.DefaultWaitForPendingUpdate(gotUpdate)

			gotResp, err = i.GetSettings()
			require.NoError(t, err)
			require.Equal(t, &tt.args.secondResponse, gotResp)
			deleteAllIndexes(c)
		})
	}
}

func TestIndex_UpdateStopWords(t *testing.T) {
	type args struct {
		UID     string
		client  *Client
		request []string
	}
	tests := []struct {
		name       string
		args       args
		wantUpdate *AsyncUpdateID
	}{
		{
			name: "TestIndexBasicUpdateStopWords",
			args: args{
				UID:    "indexUID",
				client: defaultClient,
				request: []string{
					"of", "the", "to",
				},
			},
			wantUpdate: &AsyncUpdateID{
				UpdateID: 1,
			},
		},
		{
			name: "TestIndexUpdateStopWordsWithCustomClient",
			args: args{
				UID:    "indexUID",
				client: customClient,
				request: []string{
					"of", "the", "to",
				},
			},
			wantUpdate: &AsyncUpdateID{
				UpdateID: 1,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			SetUpIndexForFaceting()
			c := tt.args.client
			i := c.Index(tt.args.UID)

			gotResp, err := i.GetStopWords()
			require.NoError(t, err)
			require.Empty(t, gotResp)

			gotUpdate, err := i.UpdateStopWords(&tt.args.request)
			require.NoError(t, err)
			require.GreaterOrEqual(t, gotUpdate.UpdateID, tt.wantUpdate.UpdateID)
			i.DefaultWaitForPendingUpdate(gotUpdate)

			gotResp, err = i.GetStopWords()
			require.NoError(t, err)
			require.Equal(t, &tt.args.request, gotResp)
			deleteAllIndexes(c)
		})
	}
}

func TestIndex_UpdateSynonyms(t *testing.T) {
	type args struct {
		UID     string
		client  *Client
		request map[string][]string
	}
	tests := []struct {
		name       string
		args       args
		wantUpdate *AsyncUpdateID
	}{
		{
			name: "TestIndexBasicUpdateSynonyms",
			args: args{
				UID:    "indexUID",
				client: defaultClient,
				request: map[string][]string{
					"wolverine": []string{"logan", "xmen"},
				},
			},
			wantUpdate: &AsyncUpdateID{
				UpdateID: 1,
			},
		},
		{
			name: "TestIndexUpdateSynonymsWithCustomClient",
			args: args{
				UID:    "indexUID",
				client: customClient,
				request: map[string][]string{
					"wolverine": []string{"logan", "xmen"},
				},
			},
			wantUpdate: &AsyncUpdateID{
				UpdateID: 1,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			SetUpIndexForFaceting()
			c := tt.args.client
			i := c.Index(tt.args.UID)

			gotResp, err := i.GetSynonyms()
			require.NoError(t, err)
			require.Empty(t, gotResp)

			gotUpdate, err := i.UpdateSynonyms(&tt.args.request)
			require.NoError(t, err)
			require.GreaterOrEqual(t, gotUpdate.UpdateID, tt.wantUpdate.UpdateID)
			i.DefaultWaitForPendingUpdate(gotUpdate)

			gotResp, err = i.GetSynonyms()
			require.NoError(t, err)
			require.Equal(t, &tt.args.request, gotResp)
			deleteAllIndexes(c)
		})
	}
}
