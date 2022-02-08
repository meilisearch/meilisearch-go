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
			t.Cleanup(cleanup(c))

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
			t.Cleanup(cleanup(c))

			gotResp, err := i.GetDisplayedAttributes()
			require.NoError(t, err)
			require.Equal(t, tt.wantResp, gotResp)
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
			t.Cleanup(cleanup(c))

			gotResp, err := i.GetDistinctAttribute()
			require.NoError(t, err)
			require.Empty(t, gotResp)
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
			wantResp: &defaultRankingRules,
		},
		{
			name: "TestIndexGetRankingRulesWithCustomClient",
			args: args{
				UID:    "indexUID",
				client: defaultClient,
			},
			wantResp: &defaultRankingRules,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			SetUpIndexForFaceting()
			c := tt.args.client
			i := c.Index(tt.args.UID)
			t.Cleanup(cleanup(c))

			gotResp, err := i.GetRankingRules()
			require.NoError(t, err)
			require.Equal(t, tt.wantResp, gotResp)
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
			t.Cleanup(cleanup(c))

			gotResp, err := i.GetSearchableAttributes()
			require.NoError(t, err)
			require.Equal(t, tt.wantResp, gotResp)
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
				RankingRules:         defaultRankingRules,
				DistinctAttribute:    (*string)(nil),
				SearchableAttributes: []string{"*"},
				DisplayedAttributes:  []string{"*"},
				StopWords:            []string{},
				Synonyms:             map[string][]string(nil),
				FilterableAttributes: []string{},
				SortableAttributes:   []string{},
			},
		},
		{
			name: "TestIndexGetSettingsWithCustomClient",
			args: args{
				UID:    "indexUID",
				client: customClient,
			},
			wantResp: &Settings{
				RankingRules:         defaultRankingRules,
				DistinctAttribute:    (*string)(nil),
				SearchableAttributes: []string{"*"},
				DisplayedAttributes:  []string{"*"},
				StopWords:            []string{},
				Synonyms:             map[string][]string(nil),
				FilterableAttributes: []string{},
				SortableAttributes:   []string{},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			SetUpIndexForFaceting()
			c := tt.args.client
			i := c.Index(tt.args.UID)
			t.Cleanup(cleanup(c))

			gotResp, err := i.GetSettings()
			require.NoError(t, err)
			require.Equal(t, tt.wantResp, gotResp)
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
			t.Cleanup(cleanup(c))

			gotResp, err := i.GetStopWords()
			require.NoError(t, err)
			require.Empty(t, gotResp)
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
			t.Cleanup(cleanup(c))

			gotResp, err := i.GetSynonyms()
			require.NoError(t, err)
			require.Empty(t, gotResp)
		})
	}
}

func TestIndex_GetSortableAttributes(t *testing.T) {
	type args struct {
		UID    string
		client *Client
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "TestIndexBasicGetSortableAttributes",
			args: args{
				UID:    "indexUID",
				client: defaultClient,
			},
		},
		{
			name: "TestIndexGetSortableAttributesWithCustomClient",
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
			t.Cleanup(cleanup(c))

			gotResp, err := i.GetSortableAttributes()
			require.NoError(t, err)
			require.Empty(t, gotResp)
		})
	}
}

func TestIndex_ResetFilterableAttributes(t *testing.T) {
	type args struct {
		UID    string
		client *Client
	}
	tests := []struct {
		name     string
		args     args
		wantTask *Task
	}{
		{
			name: "TestIndexBasicResetFilterableAttributes",
			args: args{
				UID:    "indexUID",
				client: defaultClient,
			},
			wantTask: &Task{
				UID: 1,
			},
		},
		{
			name: "TestIndexResetFilterableAttributesCustomClient",
			args: args{
				UID:    "indexUID",
				client: customClient,
			},
			wantTask: &Task{
				UID: 1,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			SetUpIndexForFaceting()
			c := tt.args.client
			i := c.Index(tt.args.UID)
			t.Cleanup(cleanup(c))

			gotTask, err := i.ResetFilterableAttributes()
			require.NoError(t, err)
			require.GreaterOrEqual(t, gotTask.UID, tt.wantTask.UID)
			testWaitForTask(t, i, gotTask)

			gotResp, err := i.GetFilterableAttributes()
			require.NoError(t, err)
			require.Empty(t, gotResp)
		})
	}
}

func TestIndex_ResetDisplayedAttributes(t *testing.T) {
	type args struct {
		UID    string
		client *Client
	}
	tests := []struct {
		name     string
		args     args
		wantTask *Task
		wantResp *[]string
	}{
		{
			name: "TestIndexBasicResetDisplayedAttributes",
			args: args{
				UID:    "indexUID",
				client: defaultClient,
			},
			wantTask: &Task{
				UID: 1,
			},
			wantResp: &[]string{"*"},
		},
		{
			name: "TestIndexResetDisplayedAttributesCustomClient",
			args: args{
				UID:    "indexUID",
				client: customClient,
			},
			wantTask: &Task{
				UID: 1,
			},
			wantResp: &[]string{"*"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			SetUpIndexForFaceting()
			c := tt.args.client
			i := c.Index(tt.args.UID)
			t.Cleanup(cleanup(c))

			gotTask, err := i.ResetDisplayedAttributes()
			require.NoError(t, err)
			require.GreaterOrEqual(t, gotTask.UID, tt.wantTask.UID)
			testWaitForTask(t, i, gotTask)

			gotResp, err := i.GetDisplayedAttributes()
			require.NoError(t, err)
			require.Equal(t, tt.wantResp, gotResp)
		})
	}
}

func TestIndex_ResetDistinctAttribute(t *testing.T) {
	type args struct {
		UID    string
		client *Client
	}
	tests := []struct {
		name     string
		args     args
		wantTask *Task
	}{
		{
			name: "TestIndexBasicResetDistinctAttribute",
			args: args{
				UID:    "indexUID",
				client: defaultClient,
			},
			wantTask: &Task{
				UID: 1,
			},
		},
		{
			name: "TestIndexResetDistinctAttributeWithCustomClient",
			args: args{
				UID:    "indexUID",
				client: customClient,
			},
			wantTask: &Task{
				UID: 1,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			SetUpIndexForFaceting()
			c := tt.args.client
			i := c.Index(tt.args.UID)
			t.Cleanup(cleanup(c))

			gotTask, err := i.ResetDistinctAttribute()
			require.NoError(t, err)
			require.GreaterOrEqual(t, gotTask.UID, tt.wantTask.UID)
			testWaitForTask(t, i, gotTask)

			gotResp, err := i.GetDistinctAttribute()
			require.NoError(t, err)
			require.Empty(t, gotResp)
		})
	}
}

func TestIndex_ResetRankingRules(t *testing.T) {
	type args struct {
		UID    string
		client *Client
	}
	tests := []struct {
		name     string
		args     args
		wantTask *Task
		wantResp *[]string
	}{
		{
			name: "TestIndexBasicResetRankingRules",
			args: args{
				UID:    "indexUID",
				client: defaultClient,
			},
			wantTask: &Task{
				UID: 1,
			},
			wantResp: &defaultRankingRules,
		},
		{
			name: "TestIndexResetRankingRulesWithCustomClient",
			args: args{
				UID:    "indexUID",
				client: customClient,
			},
			wantTask: &Task{
				UID: 1,
			},
			wantResp: &defaultRankingRules,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			SetUpIndexForFaceting()
			c := tt.args.client
			i := c.Index(tt.args.UID)
			t.Cleanup(cleanup(c))

			gotTask, err := i.ResetRankingRules()
			require.NoError(t, err)
			require.GreaterOrEqual(t, gotTask.UID, tt.wantTask.UID)
			testWaitForTask(t, i, gotTask)

			gotResp, err := i.GetRankingRules()
			require.NoError(t, err)
			require.Equal(t, tt.wantResp, gotResp)
		})
	}
}

func TestIndex_ResetSearchableAttributes(t *testing.T) {
	type args struct {
		UID    string
		client *Client
	}
	tests := []struct {
		name     string
		args     args
		wantTask *Task
		wantResp *[]string
	}{
		{
			name: "TestIndexBasicResetSearchableAttributes",
			args: args{
				UID:    "indexUID",
				client: defaultClient,
			},
			wantTask: &Task{
				UID: 1,
			},
			wantResp: &[]string{"*"},
		},
		{
			name: "TestIndexResetSearchableAttributesWithCustomClient",
			args: args{
				UID:    "indexUID",
				client: customClient,
			},
			wantTask: &Task{
				UID: 1,
			},
			wantResp: &[]string{"*"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			SetUpIndexForFaceting()
			c := tt.args.client
			i := c.Index(tt.args.UID)
			t.Cleanup(cleanup(c))

			gotTask, err := i.ResetSearchableAttributes()
			require.NoError(t, err)
			require.GreaterOrEqual(t, gotTask.UID, tt.wantTask.UID)
			testWaitForTask(t, i, gotTask)

			gotResp, err := i.GetSearchableAttributes()
			require.NoError(t, err)
			require.Equal(t, tt.wantResp, gotResp)
		})
	}
}

func TestIndex_ResetSettings(t *testing.T) {
	type args struct {
		UID    string
		client *Client
	}
	tests := []struct {
		name     string
		args     args
		wantTask *Task
		wantResp *Settings
	}{
		{
			name: "TestIndexBasicResetSettings",
			args: args{
				UID:    "indexUID",
				client: defaultClient,
			},
			wantTask: &Task{
				UID: 1,
			},
			wantResp: &Settings{
				RankingRules:         defaultRankingRules,
				DistinctAttribute:    (*string)(nil),
				SearchableAttributes: []string{"*"},
				DisplayedAttributes:  []string{"*"},
				StopWords:            []string{},
				Synonyms:             map[string][]string(nil),
				FilterableAttributes: []string{},
				SortableAttributes:   []string{},
			},
		},
		{
			name: "TestIndexResetSettingsWithCustomClient",
			args: args{
				UID:    "indexUID",
				client: customClient,
			},
			wantTask: &Task{
				UID: 1,
			},
			wantResp: &Settings{
				RankingRules:         defaultRankingRules,
				DistinctAttribute:    (*string)(nil),
				SearchableAttributes: []string{"*"},
				DisplayedAttributes:  []string{"*"},
				StopWords:            []string{},
				Synonyms:             map[string][]string(nil),
				FilterableAttributes: []string{},
				SortableAttributes:   []string{},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			SetUpIndexForFaceting()
			c := tt.args.client
			i := c.Index(tt.args.UID)
			t.Cleanup(cleanup(c))

			gotTask, err := i.ResetSettings()
			require.NoError(t, err)
			require.GreaterOrEqual(t, gotTask.UID, tt.wantTask.UID)
			testWaitForTask(t, i, gotTask)

			gotResp, err := i.GetSettings()
			require.NoError(t, err)
			require.Equal(t, tt.wantResp, gotResp)
		})
	}
}

func TestIndex_ResetStopWords(t *testing.T) {
	type args struct {
		UID    string
		client *Client
	}
	tests := []struct {
		name     string
		args     args
		wantTask *Task
	}{
		{
			name: "TestIndexBasicResetStopWords",
			args: args{
				UID:    "indexUID",
				client: defaultClient,
			},
			wantTask: &Task{
				UID: 1,
			},
		},
		{
			name: "TestIndexResetStopWordsWithCustomClient",
			args: args{
				UID:    "indexUID",
				client: customClient,
			},
			wantTask: &Task{
				UID: 1,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			SetUpIndexForFaceting()
			c := tt.args.client
			i := c.Index(tt.args.UID)
			t.Cleanup(cleanup(c))

			gotTask, err := i.ResetStopWords()
			require.NoError(t, err)
			require.GreaterOrEqual(t, gotTask.UID, tt.wantTask.UID)
			testWaitForTask(t, i, gotTask)

			gotResp, err := i.GetStopWords()
			require.NoError(t, err)
			require.Empty(t, gotResp)
		})
	}
}

func TestIndex_ResetSynonyms(t *testing.T) {
	type args struct {
		UID    string
		client *Client
	}
	tests := []struct {
		name     string
		args     args
		wantTask *Task
	}{
		{
			name: "TestIndexBasicResetSynonyms",
			args: args{
				UID:    "indexUID",
				client: defaultClient,
			},
			wantTask: &Task{
				UID: 1,
			},
		},
		{
			name: "TestIndexResetSynonymsWithCustomClient",
			args: args{
				UID:    "indexUID",
				client: customClient,
			},
			wantTask: &Task{
				UID: 1,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			SetUpIndexForFaceting()
			c := tt.args.client
			i := c.Index(tt.args.UID)
			t.Cleanup(cleanup(c))

			gotTask, err := i.ResetSynonyms()
			require.NoError(t, err)
			require.GreaterOrEqual(t, gotTask.UID, tt.wantTask.UID)
			testWaitForTask(t, i, gotTask)

			gotResp, err := i.GetSynonyms()
			require.NoError(t, err)
			require.Empty(t, gotResp)
		})
	}
}

func TestIndex_ResetSortableAttributes(t *testing.T) {
	type args struct {
		UID    string
		client *Client
	}
	tests := []struct {
		name     string
		args     args
		wantTask *Task
	}{
		{
			name: "TestIndexBasicResetSortableAttributes",
			args: args{
				UID:    "indexUID",
				client: defaultClient,
			},
			wantTask: &Task{
				UID: 1,
			},
		},
		{
			name: "TestIndexResetSortableAttributesCustomClient",
			args: args{
				UID:    "indexUID",
				client: customClient,
			},
			wantTask: &Task{
				UID: 1,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			SetUpIndexForFaceting()
			c := tt.args.client
			i := c.Index(tt.args.UID)
			t.Cleanup(cleanup(c))

			gotTask, err := i.ResetSortableAttributes()
			require.NoError(t, err)
			require.GreaterOrEqual(t, gotTask.UID, tt.wantTask.UID)
			testWaitForTask(t, i, gotTask)

			gotResp, err := i.GetSortableAttributes()
			require.NoError(t, err)
			require.Empty(t, gotResp)
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
		name     string
		args     args
		wantTask *Task
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
			wantTask: &Task{
				UID: 1,
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
			wantTask: &Task{
				UID: 1,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			SetUpIndexForFaceting()
			c := tt.args.client
			i := c.Index(tt.args.UID)
			t.Cleanup(cleanup(c))

			gotResp, err := i.GetFilterableAttributes()
			require.NoError(t, err)
			require.Empty(t, gotResp)

			gotTask, err := i.UpdateFilterableAttributes(&tt.args.request)
			require.NoError(t, err)
			require.GreaterOrEqual(t, gotTask.UID, tt.wantTask.UID)
			testWaitForTask(t, i, gotTask)

			gotResp, err = i.GetFilterableAttributes()
			require.NoError(t, err)
			require.Equal(t, &tt.args.request, gotResp)
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
		name     string
		args     args
		wantTask *Task
		wantResp *[]string
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
			wantTask: &Task{
				UID: 1,
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
			wantTask: &Task{
				UID: 1,
			},
			wantResp: &[]string{"*"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			SetUpIndexForFaceting()
			c := tt.args.client
			i := c.Index(tt.args.UID)
			t.Cleanup(cleanup(c))

			gotResp, err := i.GetDisplayedAttributes()
			require.NoError(t, err)
			require.Equal(t, tt.wantResp, gotResp)

			gotTask, err := i.UpdateDisplayedAttributes(&tt.args.request)
			require.NoError(t, err)
			require.GreaterOrEqual(t, gotTask.UID, tt.wantTask.UID)
			testWaitForTask(t, i, gotTask)

			gotResp, err = i.GetDisplayedAttributes()
			require.NoError(t, err)
			require.Equal(t, &tt.args.request, gotResp)
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
		name     string
		args     args
		wantTask *Task
	}{
		{
			name: "TestIndexBasicUpdateDistinctAttribute",
			args: args{
				UID:     "indexUID",
				client:  defaultClient,
				request: "movie_id",
			},
			wantTask: &Task{
				UID: 1,
			},
		},
		{
			name: "TestIndexUpdateDistinctAttributeWithCustomClient",
			args: args{
				UID:     "indexUID",
				client:  customClient,
				request: "movie_id",
			},
			wantTask: &Task{
				UID: 1,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			SetUpIndexForFaceting()
			c := tt.args.client
			i := c.Index(tt.args.UID)
			t.Cleanup(cleanup(c))

			gotResp, err := i.GetDistinctAttribute()
			require.NoError(t, err)
			require.Empty(t, gotResp)

			gotTask, err := i.UpdateDistinctAttribute(tt.args.request)
			require.NoError(t, err)
			require.GreaterOrEqual(t, gotTask.UID, tt.wantTask.UID)
			testWaitForTask(t, i, gotTask)

			gotResp, err = i.GetDistinctAttribute()
			require.NoError(t, err)
			require.Equal(t, &tt.args.request, gotResp)
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
		name     string
		args     args
		wantTask *Task
		wantResp *[]string
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
			wantTask: &Task{
				UID: 1,
			},
			wantResp: &defaultRankingRules,
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
			wantTask: &Task{
				UID: 1,
			},
			wantResp: &defaultRankingRules,
		},
		{
			name: "TestIndexUpdateRankingRulesAscending",
			args: args{
				UID:    "indexUID",
				client: defaultClient,
				request: []string{
					"BookID:asc",
				},
			},
			wantTask: &Task{
				UID: 1,
			},
			wantResp: &defaultRankingRules,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			SetUpIndexForFaceting()
			c := tt.args.client
			i := c.Index(tt.args.UID)
			t.Cleanup(cleanup(c))

			gotResp, err := i.GetRankingRules()
			require.NoError(t, err)
			require.Equal(t, tt.wantResp, gotResp)

			gotTask, err := i.UpdateRankingRules(&tt.args.request)
			require.NoError(t, err)
			require.GreaterOrEqual(t, gotTask.UID, tt.wantTask.UID)
			testWaitForTask(t, i, gotTask)

			gotResp, err = i.GetRankingRules()
			require.NoError(t, err)
			require.Equal(t, &tt.args.request, gotResp)
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
		name     string
		args     args
		wantTask *Task
		wantResp *[]string
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
			wantTask: &Task{
				UID: 1,
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
			wantTask: &Task{
				UID: 1,
			},
			wantResp: &[]string{"*"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			SetUpIndexForFaceting()
			c := tt.args.client
			i := c.Index(tt.args.UID)
			t.Cleanup(cleanup(c))

			gotResp, err := i.GetSearchableAttributes()
			require.NoError(t, err)
			require.Equal(t, tt.wantResp, gotResp)

			gotTask, err := i.UpdateSearchableAttributes(&tt.args.request)
			require.NoError(t, err)
			require.GreaterOrEqual(t, gotTask.UID, tt.wantTask.UID)
			testWaitForTask(t, i, gotTask)

			gotResp, err = i.GetSearchableAttributes()
			require.NoError(t, err)
			require.Equal(t, &tt.args.request, gotResp)
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
		name     string
		args     args
		wantTask *Task
		wantResp *Settings
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
					SortableAttributes: []string{
						"title",
					},
				},
			},
			wantTask: &Task{
				UID: 1,
			},
			wantResp: &Settings{
				RankingRules:         defaultRankingRules,
				DistinctAttribute:    (*string)(nil),
				SearchableAttributes: []string{"*"},
				DisplayedAttributes:  []string{"*"},
				StopWords:            []string{},
				Synonyms:             map[string][]string(nil),
				FilterableAttributes: []string{},
				SortableAttributes:   []string{},
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
					SortableAttributes: []string{
						"title",
					},
				},
			},
			wantTask: &Task{
				UID: 1,
			},
			wantResp: &Settings{
				RankingRules:         defaultRankingRules,
				DistinctAttribute:    (*string)(nil),
				SearchableAttributes: []string{"*"},
				DisplayedAttributes:  []string{"*"},
				StopWords:            []string{},
				Synonyms:             map[string][]string(nil),
				FilterableAttributes: []string{},
				SortableAttributes:   []string{},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			SetUpIndexForFaceting()
			c := tt.args.client
			i := c.Index(tt.args.UID)
			t.Cleanup(cleanup(c))

			gotResp, err := i.GetSettings()
			require.NoError(t, err)
			require.Equal(t, tt.wantResp, gotResp)

			gotTask, err := i.UpdateSettings(&tt.args.request)
			require.NoError(t, err)
			require.GreaterOrEqual(t, gotTask.UID, tt.wantTask.UID)
			testWaitForTask(t, i, gotTask)

			gotResp, err = i.GetSettings()
			require.NoError(t, err)
			require.Equal(t, &tt.args.request, gotResp)
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
		name     string
		args     args
		wantTask *Task
		wantResp *Settings
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
					SortableAttributes:   []string{},
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
					SortableAttributes:   []string{},
				},
			},
			wantTask: &Task{
				UID: 1,
			},
			wantResp: &Settings{
				RankingRules:         defaultRankingRules,
				DistinctAttribute:    (*string)(nil),
				SearchableAttributes: []string{"*"},
				DisplayedAttributes:  []string{"*"},
				StopWords:            []string{},
				Synonyms:             map[string][]string(nil),
				FilterableAttributes: []string{},
				SortableAttributes:   []string{},
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
					SortableAttributes:   []string{},
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
					SortableAttributes:   []string{},
				},
			},
			wantTask: &Task{
				UID: 1,
			},
			wantResp: &Settings{
				RankingRules:         defaultRankingRules,
				DistinctAttribute:    (*string)(nil),
				SearchableAttributes: []string{"*"},
				DisplayedAttributes:  []string{"*"},
				StopWords:            []string{},
				Synonyms:             map[string][]string(nil),
				FilterableAttributes: []string{},
				SortableAttributes:   []string{},
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
					SortableAttributes:   []string{},
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
					SortableAttributes:   []string{},
				},
			},
			wantTask: &Task{
				UID: 1,
			},
			wantResp: &Settings{
				RankingRules:         defaultRankingRules,
				DistinctAttribute:    (*string)(nil),
				SearchableAttributes: []string{"*"},
				DisplayedAttributes:  []string{"*"},
				StopWords:            []string{},
				Synonyms:             map[string][]string(nil),
				FilterableAttributes: []string{},
				SortableAttributes:   []string{},
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
					SortableAttributes:   []string{},
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
					SortableAttributes:   []string{},
				},
			},
			wantTask: &Task{
				UID: 1,
			},
			wantResp: &Settings{
				RankingRules:         defaultRankingRules,
				DistinctAttribute:    (*string)(nil),
				SearchableAttributes: []string{"*"},
				DisplayedAttributes:  []string{"*"},
				StopWords:            []string{},
				Synonyms:             map[string][]string(nil),
				FilterableAttributes: []string{},
				SortableAttributes:   []string{},
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
					SortableAttributes:   []string{},
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
					SortableAttributes:   []string{},
				},
			},
			wantTask: &Task{
				UID: 1,
			},
			wantResp: &Settings{
				RankingRules:         defaultRankingRules,
				DistinctAttribute:    (*string)(nil),
				SearchableAttributes: []string{"*"},
				DisplayedAttributes:  []string{"*"},
				StopWords:            []string{},
				Synonyms:             map[string][]string(nil),
				FilterableAttributes: []string{},
				SortableAttributes:   []string{},
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
					SortableAttributes: []string{},
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
					SortableAttributes: []string{},
				},
			},
			wantTask: &Task{
				UID: 1,
			},
			wantResp: &Settings{
				RankingRules:         defaultRankingRules,
				DistinctAttribute:    (*string)(nil),
				SearchableAttributes: []string{"*"},
				DisplayedAttributes:  []string{"*"},
				StopWords:            []string{},
				Synonyms:             map[string][]string(nil),
				FilterableAttributes: []string{},
				SortableAttributes:   []string{},
			},
		},
		{
			name: "TestIndexUpdateJustSortableAttributes",
			args: args{
				UID:    "indexUID",
				client: defaultClient,
				firstRequest: Settings{
					RankingRules: []string{
						"typo", "words",
					},
					SortableAttributes: []string{
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
					FilterableAttributes: []string{},
					SortableAttributes: []string{
						"title",
					},
				},
				secondRequest: Settings{
					SortableAttributes: []string{
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
					FilterableAttributes: []string{},
					SortableAttributes: []string{
						"title",
					},
				},
			},
			wantTask: &Task{
				UID: 1,
			},
			wantResp: &Settings{
				RankingRules:         defaultRankingRules,
				DistinctAttribute:    (*string)(nil),
				SearchableAttributes: []string{"*"},
				DisplayedAttributes:  []string{"*"},
				StopWords:            []string{},
				Synonyms:             map[string][]string(nil),
				FilterableAttributes: []string{},
				SortableAttributes:   []string{},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			SetUpIndexForFaceting()
			c := tt.args.client
			i := c.Index(tt.args.UID)
			t.Cleanup(cleanup(c))

			gotResp, err := i.GetSettings()
			require.NoError(t, err)
			require.Equal(t, tt.wantResp, gotResp)

			gotTask, err := i.UpdateSettings(&tt.args.firstRequest)
			require.NoError(t, err)
			require.GreaterOrEqual(t, gotTask.UID, tt.wantTask.UID)
			testWaitForTask(t, i, gotTask)

			gotResp, err = i.GetSettings()
			require.NoError(t, err)
			require.Equal(t, &tt.args.firstResponse, gotResp)

			gotTask, err = i.UpdateSettings(&tt.args.secondRequest)
			require.NoError(t, err)
			require.GreaterOrEqual(t, gotTask.UID, tt.wantTask.UID)

			testWaitForTask(t, i, gotTask)

			gotResp, err = i.GetSettings()
			require.NoError(t, err)
			require.Equal(t, &tt.args.secondResponse, gotResp)
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
		name     string
		args     args
		wantTask *Task
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
			wantTask: &Task{
				UID: 1,
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
			wantTask: &Task{
				UID: 1,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			SetUpIndexForFaceting()
			c := tt.args.client
			i := c.Index(tt.args.UID)
			t.Cleanup(cleanup(c))

			gotResp, err := i.GetStopWords()
			require.NoError(t, err)
			require.Empty(t, gotResp)

			gotTask, err := i.UpdateStopWords(&tt.args.request)
			require.NoError(t, err)
			require.GreaterOrEqual(t, gotTask.UID, tt.wantTask.UID)
			testWaitForTask(t, i, gotTask)

			gotResp, err = i.GetStopWords()
			require.NoError(t, err)
			require.Equal(t, &tt.args.request, gotResp)
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
		name     string
		args     args
		wantTask *Task
	}{
		{
			name: "TestIndexBasicUpdateSynonyms",
			args: args{
				UID:    "indexUID",
				client: defaultClient,
				request: map[string][]string{
					"wolverine": {"logan", "xmen"},
				},
			},
			wantTask: &Task{
				UID: 1,
			},
		},
		{
			name: "TestIndexUpdateSynonymsWithCustomClient",
			args: args{
				UID:    "indexUID",
				client: customClient,
				request: map[string][]string{
					"wolverine": {"logan", "xmen"},
				},
			},
			wantTask: &Task{
				UID: 1,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			SetUpIndexForFaceting()
			c := tt.args.client
			i := c.Index(tt.args.UID)
			t.Cleanup(cleanup(c))

			gotResp, err := i.GetSynonyms()
			require.NoError(t, err)
			require.Empty(t, gotResp)

			gotTask, err := i.UpdateSynonyms(&tt.args.request)
			require.NoError(t, err)
			require.GreaterOrEqual(t, gotTask.UID, tt.wantTask.UID)
			testWaitForTask(t, i, gotTask)

			gotResp, err = i.GetSynonyms()
			require.NoError(t, err)
			require.Equal(t, &tt.args.request, gotResp)
		})
	}
}

func TestIndex_UpdateSortableAttributes(t *testing.T) {
	type args struct {
		UID     string
		client  *Client
		request []string
	}
	tests := []struct {
		name     string
		args     args
		wantTask *Task
	}{
		{
			name: "TestIndexBasicUpdateSortableAttributes",
			args: args{
				UID:    "indexUID",
				client: defaultClient,
				request: []string{
					"title",
				},
			},
			wantTask: &Task{
				UID: 1,
			},
		},
		{
			name: "TestIndexUpdateSortableAttributesWithCustomClient",
			args: args{
				UID:    "indexUID",
				client: customClient,
				request: []string{
					"title",
				},
			},
			wantTask: &Task{
				UID: 1,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			SetUpIndexForFaceting()
			c := tt.args.client
			i := c.Index(tt.args.UID)
			t.Cleanup(cleanup(c))

			gotResp, err := i.GetSortableAttributes()
			require.NoError(t, err)
			require.Empty(t, gotResp)

			gotTask, err := i.UpdateSortableAttributes(&tt.args.request)
			require.NoError(t, err)
			require.GreaterOrEqual(t, gotTask.UID, tt.wantTask.UID)
			testWaitForTask(t, i, gotTask)

			gotResp, err = i.GetSortableAttributes()
			require.NoError(t, err)
			require.Equal(t, &tt.args.request, gotResp)
		})
	}
}
