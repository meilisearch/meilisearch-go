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
				TypoTolerance:        &defaultTypoTolerance,
				Pagination:           &defaultPagination,
				Faceting:             &defaultFaceting,
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
				TypoTolerance:        &defaultTypoTolerance,
				Pagination:           &defaultPagination,
				Faceting:             &defaultFaceting,
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

func TestIndex_GetTypoTolerance(t *testing.T) {
	type args struct {
		UID    string
		client *Client
	}
	tests := []struct {
		name     string
		args     args
		wantResp *TypoTolerance
	}{
		{
			name: "TestIndexBasicGetTypoTolerance",
			args: args{
				UID:    "indexUID",
				client: defaultClient,
			},
			wantResp: &defaultTypoTolerance,
		},
		{
			name: "TestIndexGetTypoToleranceWithCustomClient",
			args: args{
				UID:    "indexUID",
				client: customClient,
			},
			wantResp: &defaultTypoTolerance,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			SetUpIndexForFaceting()
			c := tt.args.client
			i := c.Index(tt.args.UID)
			t.Cleanup(cleanup(c))

			gotResp, err := i.GetTypoTolerance()
			require.NoError(t, err)
			require.Equal(t, tt.wantResp, gotResp)
		})
	}
}

func TestIndex_GetPagination(t *testing.T) {
	type args struct {
		UID    string
		client *Client
	}
	tests := []struct {
		name     string
		args     args
		wantResp *Pagination
	}{
		{
			name: "TestIndexBasicGetPagination",
			args: args{
				UID:    "indexUID",
				client: defaultClient,
			},
			wantResp: &defaultPagination,
		},
		{
			name: "TestIndexGetPaginationWithCustomClient",
			args: args{
				UID:    "indexUID",
				client: customClient,
			},
			wantResp: &defaultPagination,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			SetUpIndexForFaceting()
			c := tt.args.client
			i := c.Index(tt.args.UID)
			t.Cleanup(cleanup(c))

			gotResp, err := i.GetPagination()
			require.NoError(t, err)
			require.Equal(t, tt.wantResp, gotResp)
		})
	}
}

func TestIndex_GetFaceting(t *testing.T) {
	type args struct {
		UID    string
		client *Client
	}
	tests := []struct {
		name     string
		args     args
		wantResp *Faceting
	}{
		{
			name: "TestIndexBasicGetFaceting",
			args: args{
				UID:    "indexUID",
				client: defaultClient,
			},
			wantResp: &defaultFaceting,
		},
		{
			name: "TestIndexGetFacetingWithCustomClient",
			args: args{
				UID:    "indexUID",
				client: customClient,
			},
			wantResp: &defaultFaceting,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			SetUpIndexForFaceting()
			c := tt.args.client
			i := c.Index(tt.args.UID)
			t.Cleanup(cleanup(c))

			gotResp, err := i.GetFaceting()
			require.NoError(t, err)
			require.Equal(t, tt.wantResp, gotResp)
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
		wantTask *TaskInfo
	}{
		{
			name: "TestIndexBasicResetFilterableAttributes",
			args: args{
				UID:    "indexUID",
				client: defaultClient,
			},
			wantTask: &TaskInfo{
				TaskUID: 1,
			},
		},
		{
			name: "TestIndexResetFilterableAttributesCustomClient",
			args: args{
				UID:    "indexUID",
				client: customClient,
			},
			wantTask: &TaskInfo{
				TaskUID: 1,
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
			require.GreaterOrEqual(t, gotTask.TaskUID, tt.wantTask.TaskUID)
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
		wantTask *TaskInfo
		wantResp *[]string
	}{
		{
			name: "TestIndexBasicResetDisplayedAttributes",
			args: args{
				UID:    "indexUID",
				client: defaultClient,
			},
			wantTask: &TaskInfo{
				TaskUID: 1,
			},
			wantResp: &[]string{"*"},
		},
		{
			name: "TestIndexResetDisplayedAttributesCustomClient",
			args: args{
				UID:    "indexUID",
				client: customClient,
			},
			wantTask: &TaskInfo{
				TaskUID: 1,
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
			require.GreaterOrEqual(t, gotTask.TaskUID, tt.wantTask.TaskUID)
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
		wantTask *TaskInfo
	}{
		{
			name: "TestIndexBasicResetDistinctAttribute",
			args: args{
				UID:    "indexUID",
				client: defaultClient,
			},
			wantTask: &TaskInfo{
				TaskUID: 1,
			},
		},
		{
			name: "TestIndexResetDistinctAttributeWithCustomClient",
			args: args{
				UID:    "indexUID",
				client: customClient,
			},
			wantTask: &TaskInfo{
				TaskUID: 1,
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
			require.GreaterOrEqual(t, gotTask.TaskUID, tt.wantTask.TaskUID)
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
		wantTask *TaskInfo
		wantResp *[]string
	}{
		{
			name: "TestIndexBasicResetRankingRules",
			args: args{
				UID:    "indexUID",
				client: defaultClient,
			},
			wantTask: &TaskInfo{
				TaskUID: 1,
			},
			wantResp: &defaultRankingRules,
		},
		{
			name: "TestIndexResetRankingRulesWithCustomClient",
			args: args{
				UID:    "indexUID",
				client: customClient,
			},
			wantTask: &TaskInfo{
				TaskUID: 1,
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
			require.GreaterOrEqual(t, gotTask.TaskUID, tt.wantTask.TaskUID)
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
		wantTask *TaskInfo
		wantResp *[]string
	}{
		{
			name: "TestIndexBasicResetSearchableAttributes",
			args: args{
				UID:    "indexUID",
				client: defaultClient,
			},
			wantTask: &TaskInfo{
				TaskUID: 1,
			},
			wantResp: &[]string{"*"},
		},
		{
			name: "TestIndexResetSearchableAttributesWithCustomClient",
			args: args{
				UID:    "indexUID",
				client: customClient,
			},
			wantTask: &TaskInfo{
				TaskUID: 1,
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
			require.GreaterOrEqual(t, gotTask.TaskUID, tt.wantTask.TaskUID)
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
		wantTask *TaskInfo
		wantResp *Settings
	}{
		{
			name: "TestIndexBasicResetSettings",
			args: args{
				UID:    "indexUID",
				client: defaultClient,
			},
			wantTask: &TaskInfo{
				TaskUID: 1,
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
				TypoTolerance:        &defaultTypoTolerance,
				Pagination:           &defaultPagination,
				Faceting:             &defaultFaceting,
			},
		},
		{
			name: "TestIndexResetSettingsWithCustomClient",
			args: args{
				UID:    "indexUID",
				client: customClient,
			},
			wantTask: &TaskInfo{
				TaskUID: 1,
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
				TypoTolerance:        &defaultTypoTolerance,
				Pagination:           &defaultPagination,
				Faceting:             &defaultFaceting,
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
			require.GreaterOrEqual(t, gotTask.TaskUID, tt.wantTask.TaskUID)
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
		wantTask *TaskInfo
	}{
		{
			name: "TestIndexBasicResetStopWords",
			args: args{
				UID:    "indexUID",
				client: defaultClient,
			},
			wantTask: &TaskInfo{
				TaskUID: 1,
			},
		},
		{
			name: "TestIndexResetStopWordsWithCustomClient",
			args: args{
				UID:    "indexUID",
				client: customClient,
			},
			wantTask: &TaskInfo{
				TaskUID: 1,
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
			require.GreaterOrEqual(t, gotTask.TaskUID, tt.wantTask.TaskUID)
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
		wantTask *TaskInfo
	}{
		{
			name: "TestIndexBasicResetSynonyms",
			args: args{
				UID:    "indexUID",
				client: defaultClient,
			},
			wantTask: &TaskInfo{
				TaskUID: 1,
			},
		},
		{
			name: "TestIndexResetSynonymsWithCustomClient",
			args: args{
				UID:    "indexUID",
				client: customClient,
			},
			wantTask: &TaskInfo{
				TaskUID: 1,
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
			require.GreaterOrEqual(t, gotTask.TaskUID, tt.wantTask.TaskUID)
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
		wantTask *TaskInfo
	}{
		{
			name: "TestIndexBasicResetSortableAttributes",
			args: args{
				UID:    "indexUID",
				client: defaultClient,
			},
			wantTask: &TaskInfo{
				TaskUID: 1,
			},
		},
		{
			name: "TestIndexResetSortableAttributesCustomClient",
			args: args{
				UID:    "indexUID",
				client: customClient,
			},
			wantTask: &TaskInfo{
				TaskUID: 1,
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
			require.GreaterOrEqual(t, gotTask.TaskUID, tt.wantTask.TaskUID)
			testWaitForTask(t, i, gotTask)

			gotResp, err := i.GetSortableAttributes()
			require.NoError(t, err)
			require.Empty(t, gotResp)
		})
	}
}

func TestIndex_ResetTypoTolerance(t *testing.T) {
	type args struct {
		UID    string
		client *Client
	}
	tests := []struct {
		name     string
		args     args
		wantTask *TaskInfo
		wantResp *TypoTolerance
	}{
		{
			name: "TestIndexBasicResetTypoTolerance",
			args: args{
				UID:    "indexUID",
				client: defaultClient,
			},
			wantTask: &TaskInfo{
				TaskUID: 1,
			},
			wantResp: &defaultTypoTolerance,
		},
		{
			name: "TestIndexResetTypoToleranceWithCustomClient",
			args: args{
				UID:    "indexUID",
				client: customClient,
			},
			wantTask: &TaskInfo{
				TaskUID: 1,
			},
			wantResp: &defaultTypoTolerance,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			SetUpIndexForFaceting()
			c := tt.args.client
			i := c.Index(tt.args.UID)
			t.Cleanup(cleanup(c))

			gotTask, err := i.ResetTypoTolerance()
			require.NoError(t, err)
			require.GreaterOrEqual(t, gotTask.TaskUID, tt.wantTask.TaskUID)
			testWaitForTask(t, i, gotTask)

			gotResp, err := i.GetTypoTolerance()
			require.NoError(t, err)
			require.Equal(t, tt.wantResp, gotResp)
		})
	}
}

func TestIndex_ResetPagination(t *testing.T) {
	type args struct {
		UID    string
		client *Client
	}
	tests := []struct {
		name     string
		args     args
		wantTask *TaskInfo
		wantResp *Pagination
	}{
		{
			name: "TestIndexBasicResetPagination",
			args: args{
				UID:    "indexUID",
				client: defaultClient,
			},
			wantTask: &TaskInfo{
				TaskUID: 1,
			},
			wantResp: &defaultPagination,
		},
		{
			name: "TestIndexResetPaginationWithCustomClient",
			args: args{
				UID:    "indexUID",
				client: customClient,
			},
			wantTask: &TaskInfo{
				TaskUID: 1,
			},
			wantResp: &defaultPagination,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			SetUpIndexForFaceting()
			c := tt.args.client
			i := c.Index(tt.args.UID)
			t.Cleanup(cleanup(c))

			gotTask, err := i.ResetPagination()
			require.NoError(t, err)
			require.GreaterOrEqual(t, gotTask.TaskUID, tt.wantTask.TaskUID)
			testWaitForTask(t, i, gotTask)

			gotResp, err := i.GetPagination()
			require.NoError(t, err)
			require.Equal(t, tt.wantResp, gotResp)
		})
	}
}

func TestIndex_ResetFaceting(t *testing.T) {
	type args struct {
		UID    string
		client *Client
	}
	tests := []struct {
		name     string
		args     args
		wantTask *TaskInfo
		wantResp *Faceting
	}{
		{
			name: "TestIndexBasicResetFaceting",
			args: args{
				UID:    "indexUID",
				client: defaultClient,
			},
			wantTask: &TaskInfo{
				TaskUID: 1,
			},
			wantResp: &defaultFaceting,
		},
		{
			name: "TestIndexResetFacetingWithCustomClient",
			args: args{
				UID:    "indexUID",
				client: customClient,
			},
			wantTask: &TaskInfo{
				TaskUID: 1,
			},
			wantResp: &defaultFaceting,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			SetUpIndexForFaceting()
			c := tt.args.client
			i := c.Index(tt.args.UID)
			t.Cleanup(cleanup(c))

			gotTask, err := i.ResetFaceting()
			require.NoError(t, err)
			require.GreaterOrEqual(t, gotTask.TaskUID, tt.wantTask.TaskUID)
			testWaitForTask(t, i, gotTask)

			gotResp, err := i.GetFaceting()
			require.NoError(t, err)
			require.Equal(t, tt.wantResp, gotResp)
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
		wantTask *TaskInfo
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
			wantTask: &TaskInfo{
				TaskUID: 1,
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
			wantTask: &TaskInfo{
				TaskUID: 1,
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
			require.GreaterOrEqual(t, gotTask.TaskUID, tt.wantTask.TaskUID)
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
		wantTask *TaskInfo
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
			wantTask: &TaskInfo{
				TaskUID: 1,
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
			wantTask: &TaskInfo{
				TaskUID: 1,
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
			require.GreaterOrEqual(t, gotTask.TaskUID, tt.wantTask.TaskUID)
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
		wantTask *TaskInfo
	}{
		{
			name: "TestIndexBasicUpdateDistinctAttribute",
			args: args{
				UID:     "indexUID",
				client:  defaultClient,
				request: "movie_id",
			},
			wantTask: &TaskInfo{
				TaskUID: 1,
			},
		},
		{
			name: "TestIndexUpdateDistinctAttributeWithCustomClient",
			args: args{
				UID:     "indexUID",
				client:  customClient,
				request: "movie_id",
			},
			wantTask: &TaskInfo{
				TaskUID: 1,
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
			require.GreaterOrEqual(t, gotTask.TaskUID, tt.wantTask.TaskUID)
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
		wantTask *TaskInfo
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
			wantTask: &TaskInfo{
				TaskUID: 1,
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
			wantTask: &TaskInfo{
				TaskUID: 1,
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
			wantTask: &TaskInfo{
				TaskUID: 1,
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
			require.GreaterOrEqual(t, gotTask.TaskUID, tt.wantTask.TaskUID)
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
		wantTask *TaskInfo
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
			wantTask: &TaskInfo{
				TaskUID: 1,
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
			wantTask: &TaskInfo{
				TaskUID: 1,
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
			require.GreaterOrEqual(t, gotTask.TaskUID, tt.wantTask.TaskUID)
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
		wantTask *TaskInfo
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
					TypoTolerance: &TypoTolerance{
						Enabled: true,
						MinWordSizeForTypos: MinWordSizeForTypos{
							OneTypo:  7,
							TwoTypos: 10,
						},
						DisableOnWords:      []string{},
						DisableOnAttributes: []string{},
					},
					Pagination: &Pagination{
						MaxTotalHits: 1200,
					},
					Faceting: &Faceting{
						MaxValuesPerFacet: 200,
					},
				},
			},
			wantTask: &TaskInfo{
				TaskUID: 1,
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
				TypoTolerance:        &defaultTypoTolerance,
				Pagination:           &defaultPagination,
				Faceting:             &defaultFaceting,
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
					TypoTolerance: &TypoTolerance{
						Enabled: true,
						MinWordSizeForTypos: MinWordSizeForTypos{
							OneTypo:  7,
							TwoTypos: 10,
						},
						DisableOnWords:      []string{},
						DisableOnAttributes: []string{},
					},
					Pagination: &Pagination{
						MaxTotalHits: 1200,
					},
					Faceting: &Faceting{
						MaxValuesPerFacet: 200,
					},
				},
			},
			wantTask: &TaskInfo{
				TaskUID: 1,
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
				TypoTolerance:        &defaultTypoTolerance,
				Pagination:           &defaultPagination,
				Faceting:             &defaultFaceting,
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
			require.GreaterOrEqual(t, gotTask.TaskUID, tt.wantTask.TaskUID)
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
		wantTask *TaskInfo
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
					TypoTolerance:        &defaultTypoTolerance,
					Pagination:           &defaultPagination,
					Faceting:             &defaultFaceting,
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
					TypoTolerance:        &defaultTypoTolerance,
					Pagination:           &defaultPagination,
					Faceting:             &defaultFaceting,
				},
			},
			wantTask: &TaskInfo{
				TaskUID: 1,
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
				TypoTolerance:        &defaultTypoTolerance,
				Pagination:           &defaultPagination,
				Faceting:             &defaultFaceting,
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
					TypoTolerance:        &defaultTypoTolerance,
					Pagination:           &defaultPagination,
					Faceting:             &defaultFaceting,
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
					TypoTolerance:        &defaultTypoTolerance,
					Pagination:           &defaultPagination,
					Faceting:             &defaultFaceting,
				},
			},
			wantTask: &TaskInfo{
				TaskUID: 1,
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
				TypoTolerance:        &defaultTypoTolerance,
				Pagination:           &defaultPagination,
				Faceting:             &defaultFaceting,
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
					TypoTolerance:        &defaultTypoTolerance,
					Pagination:           &defaultPagination,
					Faceting:             &defaultFaceting,
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
					TypoTolerance:        &defaultTypoTolerance,
					Pagination:           &defaultPagination,
					Faceting:             &defaultFaceting,
				},
			},
			wantTask: &TaskInfo{
				TaskUID: 1,
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
				TypoTolerance:        &defaultTypoTolerance,
				Pagination:           &defaultPagination,
				Faceting:             &defaultFaceting,
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
					TypoTolerance:        &defaultTypoTolerance,
					Pagination:           &defaultPagination,
					Faceting:             &defaultFaceting,
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
					TypoTolerance:        &defaultTypoTolerance,
					Pagination:           &defaultPagination,
					Faceting:             &defaultFaceting,
				},
			},
			wantTask: &TaskInfo{
				TaskUID: 1,
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
				TypoTolerance:        &defaultTypoTolerance,
				Pagination:           &defaultPagination,
				Faceting:             &defaultFaceting,
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
					TypoTolerance:        &defaultTypoTolerance,
					Pagination:           &defaultPagination,
					Faceting:             &defaultFaceting,
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
					TypoTolerance:        &defaultTypoTolerance,
					Pagination:           &defaultPagination,
					Faceting:             &defaultFaceting,
				},
			},
			wantTask: &TaskInfo{
				TaskUID: 1,
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
				TypoTolerance:        &defaultTypoTolerance,
				Pagination:           &defaultPagination,
				Faceting:             &defaultFaceting,
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
					TypoTolerance:      &defaultTypoTolerance,
					Pagination:         &defaultPagination,
					Faceting:           &defaultFaceting,
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
					TypoTolerance:      &defaultTypoTolerance,
					Pagination:         &defaultPagination,
					Faceting:           &defaultFaceting,
				},
			},
			wantTask: &TaskInfo{
				TaskUID: 1,
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
				TypoTolerance:        &defaultTypoTolerance,
				Pagination:           &defaultPagination,
				Faceting:             &defaultFaceting,
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
					TypoTolerance: &defaultTypoTolerance,
					Pagination:    &defaultPagination,
					Faceting:      &defaultFaceting,
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
					TypoTolerance: &defaultTypoTolerance,
					Pagination:    &defaultPagination,
					Faceting:      &defaultFaceting,
				},
			},
			wantTask: &TaskInfo{
				TaskUID: 1,
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
				TypoTolerance:        &defaultTypoTolerance,
				Pagination:           &defaultPagination,
				Faceting:             &defaultFaceting,
			},
		},
		{
			name: "TestIndexUpdateJustTypoTolerance",
			args: args{
				UID:    "indexUID",
				client: defaultClient,
				firstRequest: Settings{
					RankingRules: []string{
						"typo", "words",
					},
					TypoTolerance: &TypoTolerance{
						Enabled: true,
						MinWordSizeForTypos: MinWordSizeForTypos{
							OneTypo:  7,
							TwoTypos: 10,
						},
						DisableOnWords:      []string{},
						DisableOnAttributes: []string{},
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
					TypoTolerance: &TypoTolerance{
						Enabled: true,
						MinWordSizeForTypos: MinWordSizeForTypos{
							OneTypo:  7,
							TwoTypos: 10,
						},
						DisableOnWords:      []string{},
						DisableOnAttributes: []string{},
					},
					FilterableAttributes: []string{},
					SortableAttributes:   []string{},
					Pagination:           &defaultPagination,
					Faceting:             &defaultFaceting,
				},
				secondRequest: Settings{
					TypoTolerance: &TypoTolerance{
						Enabled: true,
						MinWordSizeForTypos: MinWordSizeForTypos{
							OneTypo:  7,
							TwoTypos: 10,
						},
						DisableOnWords: []string{
							"and",
						},
						DisableOnAttributes: []string{
							"year",
						},
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
					SortableAttributes:   []string{},
					TypoTolerance: &TypoTolerance{
						Enabled: true,
						MinWordSizeForTypos: MinWordSizeForTypos{
							OneTypo:  7,
							TwoTypos: 10,
						},
						DisableOnWords: []string{
							"and",
						},
						DisableOnAttributes: []string{
							"year",
						},
					},
					Pagination: &defaultPagination,
					Faceting:   &defaultFaceting,
				},
			},
			wantTask: &TaskInfo{
				TaskUID: 1,
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
				TypoTolerance:        &defaultTypoTolerance,
				Pagination:           &defaultPagination,
				Faceting:             &defaultFaceting,
			},
		},
		{
			name: "TestIndexUpdateJustPagination",
			args: args{
				UID:    "indexUID",
				client: defaultClient,
				firstRequest: Settings{
					RankingRules: []string{
						"typo", "words",
					},
					Pagination: &Pagination{
						MaxTotalHits: 1200,
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
					SortableAttributes:   []string{},
					TypoTolerance:        &defaultTypoTolerance,
					Pagination: &Pagination{
						MaxTotalHits: 1200,
					},
					Faceting: &defaultFaceting,
				},
				secondRequest: Settings{
					Pagination: &Pagination{
						MaxTotalHits: 1200,
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
					SortableAttributes:   []string{},
					TypoTolerance:        &defaultTypoTolerance,
					Pagination: &Pagination{
						MaxTotalHits: 1200,
					},
					Faceting: &defaultFaceting,
				},
			},
			wantTask: &TaskInfo{
				TaskUID: 1,
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
				TypoTolerance:        &defaultTypoTolerance,
				Pagination:           &defaultPagination,
				Faceting:             &defaultFaceting,
			},
		},
		{
			name: "TestIndexUpdateJustFaceting",
			args: args{
				UID:    "indexUID",
				client: defaultClient,
				firstRequest: Settings{
					RankingRules: []string{
						"typo", "words",
					},
					Faceting: &Faceting{
						MaxValuesPerFacet: 200,
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
					SortableAttributes:   []string{},
					TypoTolerance:        &defaultTypoTolerance,
					Pagination:           &defaultPagination,
					Faceting: &Faceting{
						MaxValuesPerFacet: 200,
					},
				},
				secondRequest: Settings{
					Faceting: &Faceting{
						MaxValuesPerFacet: 200,
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
					SortableAttributes:   []string{},
					TypoTolerance:        &defaultTypoTolerance,
					Pagination:           &defaultPagination,
					Faceting: &Faceting{
						MaxValuesPerFacet: 200,
					},
				},
			},
			wantTask: &TaskInfo{
				TaskUID: 1,
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
				TypoTolerance:        &defaultTypoTolerance,
				Pagination:           &defaultPagination,
				Faceting:             &defaultFaceting,
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
			require.GreaterOrEqual(t, gotTask.TaskUID, tt.wantTask.TaskUID)
			testWaitForTask(t, i, gotTask)

			gotResp, err = i.GetSettings()
			require.NoError(t, err)
			require.Equal(t, &tt.args.firstResponse, gotResp)

			gotTask, err = i.UpdateSettings(&tt.args.secondRequest)
			require.NoError(t, err)
			require.GreaterOrEqual(t, gotTask.TaskUID, tt.wantTask.TaskUID)

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
		wantTask *TaskInfo
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
			wantTask: &TaskInfo{
				TaskUID: 1,
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
			wantTask: &TaskInfo{
				TaskUID: 1,
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
			require.GreaterOrEqual(t, gotTask.TaskUID, tt.wantTask.TaskUID)
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
		wantTask *TaskInfo
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
			wantTask: &TaskInfo{
				TaskUID: 1,
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
			wantTask: &TaskInfo{
				TaskUID: 1,
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
			require.GreaterOrEqual(t, gotTask.TaskUID, tt.wantTask.TaskUID)
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
		wantTask *TaskInfo
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
			wantTask: &TaskInfo{
				TaskUID: 1,
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
			wantTask: &TaskInfo{
				TaskUID: 1,
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
			require.GreaterOrEqual(t, gotTask.TaskUID, tt.wantTask.TaskUID)
			testWaitForTask(t, i, gotTask)

			gotResp, err = i.GetSortableAttributes()
			require.NoError(t, err)
			require.Equal(t, &tt.args.request, gotResp)
		})
	}
}

func TestIndex_UpdateTypoTolerance(t *testing.T) {
	type args struct {
		UID     string
		client  *Client
		request TypoTolerance
	}
	tests := []struct {
		name     string
		args     args
		wantTask *TaskInfo
		wantResp *TypoTolerance
	}{
		{
			name: "TestIndexBasicUpdateTypoTolerance",
			args: args{
				UID:    "indexUID",
				client: defaultClient,
				request: TypoTolerance{
					Enabled: true,
					MinWordSizeForTypos: MinWordSizeForTypos{
						OneTypo:  7,
						TwoTypos: 10,
					},
					DisableOnWords:      []string{},
					DisableOnAttributes: []string{},
				},
			},
			wantTask: &TaskInfo{
				TaskUID: 1,
			},
			wantResp: &defaultTypoTolerance,
		},
		{
			name: "TestIndexUpdateTypoToleranceWithCustomClient",
			args: args{
				UID:    "indexUID",
				client: customClient,
				request: TypoTolerance{
					Enabled: true,
					MinWordSizeForTypos: MinWordSizeForTypos{
						OneTypo:  7,
						TwoTypos: 10,
					},
					DisableOnWords:      []string{},
					DisableOnAttributes: []string{},
				},
			},
			wantTask: &TaskInfo{
				TaskUID: 1,
			},
			wantResp: &defaultTypoTolerance,
		},
		{
			name: "TestIndexUpdateTypoToleranceWithDisableOnWords",
			args: args{
				UID:    "indexUID",
				client: defaultClient,
				request: TypoTolerance{
					Enabled: true,
					MinWordSizeForTypos: MinWordSizeForTypos{
						OneTypo:  7,
						TwoTypos: 10,
					},
					DisableOnWords: []string{
						"and",
					},
					DisableOnAttributes: []string{},
				},
			},
			wantTask: &TaskInfo{
				TaskUID: 1,
			},
			wantResp: &defaultTypoTolerance,
		},
		{
			name: "TestIndexUpdateTypoToleranceWithDisableOnAttributes",
			args: args{
				UID:    "indexUID",
				client: defaultClient,
				request: TypoTolerance{
					Enabled: true,
					MinWordSizeForTypos: MinWordSizeForTypos{
						OneTypo:  7,
						TwoTypos: 10,
					},
					DisableOnWords: []string{},
					DisableOnAttributes: []string{
						"year",
					},
				},
			},
			wantTask: &TaskInfo{
				TaskUID: 1,
			},
			wantResp: &defaultTypoTolerance,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			SetUpIndexForFaceting()
			c := tt.args.client
			i := c.Index(tt.args.UID)
			t.Cleanup(cleanup(c))

			gotResp, err := i.GetTypoTolerance()
			require.NoError(t, err)
			require.Equal(t, tt.wantResp, gotResp)

			gotTask, err := i.UpdateTypoTolerance(&tt.args.request)
			require.NoError(t, err)
			require.GreaterOrEqual(t, gotTask.TaskUID, tt.wantTask.TaskUID)
			testWaitForTask(t, i, gotTask)

			gotResp, err = i.GetTypoTolerance()
			require.NoError(t, err)
			require.Equal(t, &tt.args.request, gotResp)
		})
	}
}

func TestIndex_UpdatePagination(t *testing.T) {
	type args struct {
		UID     string
		client  *Client
		request Pagination
	}
	tests := []struct {
		name     string
		args     args
		wantTask *TaskInfo
		wantResp *Pagination
	}{
		{
			name: "TestIndexBasicUpdatePagination",
			args: args{
				UID:    "indexUID",
				client: defaultClient,
				request: Pagination{
					MaxTotalHits: 1200,
				},
			},
			wantTask: &TaskInfo{
				TaskUID: 1,
			},
			wantResp: &defaultPagination,
		},
		{
			name: "TestIndexUpdatePaginationWithCustomClient",
			args: args{
				UID:    "indexUID",
				client: customClient,
				request: Pagination{
					MaxTotalHits: 1200,
				},
			},
			wantTask: &TaskInfo{
				TaskUID: 1,
			},
			wantResp: &defaultPagination,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			SetUpIndexForFaceting()
			c := tt.args.client
			i := c.Index(tt.args.UID)
			t.Cleanup(cleanup(c))

			gotResp, err := i.GetPagination()
			require.NoError(t, err)
			require.Equal(t, tt.wantResp, gotResp)

			gotTask, err := i.UpdatePagination(&tt.args.request)
			require.NoError(t, err)
			require.GreaterOrEqual(t, gotTask.TaskUID, tt.wantTask.TaskUID)
			testWaitForTask(t, i, gotTask)

			gotResp, err = i.GetPagination()
			require.NoError(t, err)
			require.Equal(t, &tt.args.request, gotResp)
		})
	}
}

func TestIndex_UpdateFaceting(t *testing.T) {
	type args struct {
		UID     string
		client  *Client
		request Faceting
	}
	tests := []struct {
		name     string
		args     args
		wantTask *TaskInfo
		wantResp *Faceting
	}{
		{
			name: "TestIndexBasicUpdateFaceting",
			args: args{
				UID:    "indexUID",
				client: defaultClient,
				request: Faceting{
					MaxValuesPerFacet: 200,
				},
			},
			wantTask: &TaskInfo{
				TaskUID: 1,
			},
			wantResp: &defaultFaceting,
		},
		{
			name: "TestIndexUpdateFacetingWithCustomClient",
			args: args{
				UID:    "indexUID",
				client: customClient,
				request: Faceting{
					MaxValuesPerFacet: 200,
				},
			},
			wantTask: &TaskInfo{
				TaskUID: 1,
			},
			wantResp: &defaultFaceting,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			SetUpIndexForFaceting()
			c := tt.args.client
			i := c.Index(tt.args.UID)
			t.Cleanup(cleanup(c))

			gotResp, err := i.GetFaceting()
			require.NoError(t, err)
			require.Equal(t, tt.wantResp, gotResp)

			gotTask, err := i.UpdateFaceting(&tt.args.request)
			require.NoError(t, err)
			require.GreaterOrEqual(t, gotTask.TaskUID, tt.wantTask.TaskUID)
			testWaitForTask(t, i, gotTask)

			gotResp, err = i.GetFaceting()
			require.NoError(t, err)
			require.Equal(t, &tt.args.request, gotResp)
		})
	}
}
