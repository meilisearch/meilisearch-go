package integration

import (
	"crypto/tls"
	"testing"

	"github.com/meilisearch/meilisearch-go"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestIndex_GetFilterableAttributes(t *testing.T) {
	meili := setup(t, "")
	customMeili := setup(t, "", meilisearch.WithCustomClientWithTLS(&tls.Config{
		InsecureSkipVerify: true,
	}))

	type args struct {
		UID    string
		client meilisearch.ServiceManager
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "TestIndexBasicGetFilterableAttributes",
			args: args{
				UID:    "indexUID",
				client: meili,
			},
		},
		{
			name: "TestIndexGetFilterableAttributesWithCustomClient",
			args: args{
				UID:    "indexUID",
				client: customMeili,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			setUpIndexForFaceting(tt.args.client)
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
	meili := setup(t, "")
	customMeili := setup(t, "", meilisearch.WithCustomClientWithTLS(&tls.Config{
		InsecureSkipVerify: true,
	}))

	type args struct {
		UID    string
		client meilisearch.ServiceManager
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
				client: meili,
			},
			wantResp: &[]string{"*"},
		},
		{
			name: "TestIndexGetDisplayedAttributesWithCustomClient",
			args: args{
				UID:    "indexUID",
				client: customMeili,
			},
			wantResp: &[]string{"*"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			setUpIndexForFaceting(tt.args.client)
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
	meili := setup(t, "")

	type args struct {
		UID    string
		client meilisearch.ServiceManager
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "TestIndexBasicGetDistinctAttribute",
			args: args{
				UID:    "indexUID",
				client: meili,
			},
		},
		{
			name: "TestIndexBasicGetDistinctAttribute",
			args: args{
				UID:    "indexUID",
				client: meili,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			setUpIndexForFaceting(tt.args.client)
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
	meili := setup(t, "")

	type args struct {
		UID    string
		client meilisearch.ServiceManager
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
				client: meili,
			},
			wantResp: &defaultRankingRules,
		},
		{
			name: "TestIndexGetRankingRulesWithCustomClient",
			args: args{
				UID:    "indexUID",
				client: meili,
			},
			wantResp: &defaultRankingRules,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			setUpIndexForFaceting(tt.args.client)
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
	meili := setup(t, "")

	type args struct {
		UID    string
		client meilisearch.ServiceManager
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
				client: meili,
			},
			wantResp: &[]string{"*"},
		},
		{
			name: "TestIndexGetSearchableAttributesCustomClient",
			args: args{
				UID:    "indexUID",
				client: meili,
			},
			wantResp: &[]string{"*"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			setUpIndexForFaceting(tt.args.client)
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
	meili := setup(t, "")
	customMeili := setup(t, "", meilisearch.WithCustomClientWithTLS(&tls.Config{
		InsecureSkipVerify: true,
	}))

	type args struct {
		UID    string
		client meilisearch.ServiceManager
	}
	tests := []struct {
		name     string
		args     args
		wantResp *meilisearch.Settings
	}{
		{
			name: "TestIndexBasicGetSettings",
			args: args{
				UID:    "indexUID",
				client: meili,
			},
			wantResp: &meilisearch.Settings{
				RankingRules:         defaultRankingRules,
				DistinctAttribute:    (*string)(nil),
				SearchableAttributes: []string{"*"},
				SearchCutoffMs:       0,
				ProximityPrecision:   meilisearch.ByWord,
				DisplayedAttributes:  []string{"*"},
				StopWords:            []string{},
				Synonyms:             map[string][]string{},
				FilterableAttributes: []string{},
				SortableAttributes:   []string{},
				TypoTolerance:        &defaultTypoTolerance,
				Pagination:           &defaultPagination,
				Faceting:             &defaultFaceting,
				SeparatorTokens:      make([]string, 0),
				NonSeparatorTokens:   make([]string, 0),
				Dictionary:           make([]string, 0),
				LocalizedAttributes:  nil,
				PrefixSearch:         stringPtr("indexingTime"),
				FacetSearch:          true,
				Embedders:            make(map[string]meilisearch.Embedder),
				Chat:                 &defaultChat,
			},
		},
		{
			name: "TestIndexGetSettingsWithCustomClient",
			args: args{
				UID:    "indexUID",
				client: customMeili,
			},
			wantResp: &meilisearch.Settings{
				RankingRules:         defaultRankingRules,
				DistinctAttribute:    (*string)(nil),
				SearchableAttributes: []string{"*"},
				SearchCutoffMs:       0,
				ProximityPrecision:   meilisearch.ByWord,
				DisplayedAttributes:  []string{"*"},
				StopWords:            []string{},
				Synonyms:             map[string][]string{},
				FilterableAttributes: []string{},
				SortableAttributes:   []string{},
				TypoTolerance:        &defaultTypoTolerance,
				Pagination:           &defaultPagination,
				Faceting:             &defaultFaceting,
				SeparatorTokens:      make([]string, 0),
				NonSeparatorTokens:   make([]string, 0),
				Dictionary:           make([]string, 0),
				LocalizedAttributes:  nil,
				PrefixSearch:         stringPtr("indexingTime"),
				FacetSearch:          true,
				Embedders:            make(map[string]meilisearch.Embedder),
				Chat:                 &defaultChat,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			setUpIndexForFaceting(tt.args.client)
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
	meili := setup(t, "")
	customMeili := setup(t, "", meilisearch.WithCustomClientWithTLS(&tls.Config{
		InsecureSkipVerify: true,
	}))

	type args struct {
		UID    string
		client meilisearch.ServiceManager
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "TestIndexBasicGetStopWords",
			args: args{
				UID:    "indexUID",
				client: meili,
			},
		},
		{
			name: "TestIndexGetStopWordsCustomClient",
			args: args{
				UID:    "indexUID",
				client: customMeili,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			setUpIndexForFaceting(tt.args.client)
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
	meili := setup(t, "")
	customMeili := setup(t, "", meilisearch.WithCustomClientWithTLS(&tls.Config{
		InsecureSkipVerify: true,
	}))

	type args struct {
		UID    string
		client meilisearch.ServiceManager
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "TestIndexBasicGetSynonyms",
			args: args{
				UID:    "indexUID",
				client: meili,
			},
		},
		{
			name: "TestIndexGetSynonymsWithCustomClient",
			args: args{
				UID:    "indexUID",
				client: customMeili,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			setUpIndexForFaceting(tt.args.client)
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
	meili := setup(t, "")
	customMeili := setup(t, "", meilisearch.WithCustomClientWithTLS(&tls.Config{
		InsecureSkipVerify: true,
	}))

	type args struct {
		UID    string
		client meilisearch.ServiceManager
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "TestIndexBasicGetSortableAttributes",
			args: args{
				UID:    "indexUID",
				client: meili,
			},
		},
		{
			name: "TestIndexGetSortableAttributesWithCustomClient",
			args: args{
				UID:    "indexUID",
				client: customMeili,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			setUpIndexForFaceting(tt.args.client)
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
	meili := setup(t, "")
	customMeili := setup(t, "", meilisearch.WithCustomClientWithTLS(&tls.Config{
		InsecureSkipVerify: true,
	}))

	type args struct {
		UID    string
		client meilisearch.ServiceManager
	}
	tests := []struct {
		name     string
		args     args
		wantResp *meilisearch.TypoTolerance
	}{
		{
			name: "TestIndexBasicGetTypoTolerance",
			args: args{
				UID:    "indexUID",
				client: meili,
			},
			wantResp: &defaultTypoTolerance,
		},
		{
			name: "TestIndexGetTypoToleranceWithCustomClient",
			args: args{
				UID:    "indexUID",
				client: customMeili,
			},
			wantResp: &defaultTypoTolerance,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			setUpIndexForFaceting(tt.args.client)
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
	meili := setup(t, "")
	customMeili := setup(t, "", meilisearch.WithCustomClientWithTLS(&tls.Config{
		InsecureSkipVerify: true,
	}))

	type args struct {
		UID    string
		client meilisearch.ServiceManager
	}
	tests := []struct {
		name     string
		args     args
		wantResp *meilisearch.Pagination
	}{
		{
			name: "TestIndexBasicGetPagination",
			args: args{
				UID:    "indexUID",
				client: meili,
			},
			wantResp: &defaultPagination,
		},
		{
			name: "TestIndexGetPaginationWithCustomClient",
			args: args{
				UID:    "indexUID",
				client: customMeili,
			},
			wantResp: &defaultPagination,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			setUpIndexForFaceting(tt.args.client)
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
	meili := setup(t, "")
	customMeili := setup(t, "", meilisearch.WithCustomClientWithTLS(&tls.Config{
		InsecureSkipVerify: true,
	}))

	type args struct {
		UID    string
		client meilisearch.ServiceManager
	}
	tests := []struct {
		name     string
		args     args
		wantResp *meilisearch.Faceting
	}{
		{
			name: "TestIndexBasicGetFaceting",
			args: args{
				UID:    "indexUID",
				client: meili,
			},
			wantResp: &defaultFaceting,
		},
		{
			name: "TestIndexGetFacetingWithCustomClient",
			args: args{
				UID:    "indexUID",
				client: customMeili,
			},
			wantResp: &defaultFaceting,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			setUpIndexForFaceting(tt.args.client)
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
	meili := setup(t, "")
	customMeili := setup(t, "", meilisearch.WithCustomClientWithTLS(&tls.Config{
		InsecureSkipVerify: true,
	}))

	type args struct {
		UID    string
		client meilisearch.ServiceManager
	}
	tests := []struct {
		name     string
		args     args
		wantTask *meilisearch.TaskInfo
	}{
		{
			name: "TestIndexBasicResetFilterableAttributes",
			args: args{
				UID:    "indexUID",
				client: meili,
			},
			wantTask: &meilisearch.TaskInfo{
				TaskUID: 1,
			},
		},
		{
			name: "TestIndexResetFilterableAttributesCustomClient",
			args: args{
				UID:    "indexUID",
				client: customMeili,
			},
			wantTask: &meilisearch.TaskInfo{
				TaskUID: 1,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			setUpIndexForFaceting(tt.args.client)
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
	meili := setup(t, "")
	customMeili := setup(t, "", meilisearch.WithCustomClientWithTLS(&tls.Config{
		InsecureSkipVerify: true,
	}))

	type args struct {
		UID    string
		client meilisearch.ServiceManager
	}
	tests := []struct {
		name     string
		args     args
		wantTask *meilisearch.TaskInfo
		wantResp *[]string
	}{
		{
			name: "TestIndexBasicResetDisplayedAttributes",
			args: args{
				UID:    "indexUID",
				client: meili,
			},
			wantTask: &meilisearch.TaskInfo{
				TaskUID: 1,
			},
			wantResp: &[]string{"*"},
		},
		{
			name: "TestIndexResetDisplayedAttributesCustomClient",
			args: args{
				UID:    "indexUID",
				client: customMeili,
			},
			wantTask: &meilisearch.TaskInfo{
				TaskUID: 1,
			},
			wantResp: &[]string{"*"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			setUpIndexForFaceting(tt.args.client)
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
	meili := setup(t, "")
	customMeili := setup(t, "", meilisearch.WithCustomClientWithTLS(&tls.Config{
		InsecureSkipVerify: true,
	}))

	type args struct {
		UID    string
		client meilisearch.ServiceManager
	}
	tests := []struct {
		name     string
		args     args
		wantTask *meilisearch.TaskInfo
	}{
		{
			name: "TestIndexBasicResetDistinctAttribute",
			args: args{
				UID:    "indexUID",
				client: meili,
			},
			wantTask: &meilisearch.TaskInfo{
				TaskUID: 1,
			},
		},
		{
			name: "TestIndexResetDistinctAttributeWithCustomClient",
			args: args{
				UID:    "indexUID",
				client: customMeili,
			},
			wantTask: &meilisearch.TaskInfo{
				TaskUID: 1,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			setUpIndexForFaceting(tt.args.client)
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
	meili := setup(t, "")
	customMeili := setup(t, "", meilisearch.WithCustomClientWithTLS(&tls.Config{
		InsecureSkipVerify: true,
	}))

	type args struct {
		UID    string
		client meilisearch.ServiceManager
	}
	tests := []struct {
		name     string
		args     args
		wantTask *meilisearch.TaskInfo
		wantResp *[]string
	}{
		{
			name: "TestIndexBasicResetRankingRules",
			args: args{
				UID:    "indexUID",
				client: meili,
			},
			wantTask: &meilisearch.TaskInfo{
				TaskUID: 1,
			},
			wantResp: &defaultRankingRules,
		},
		{
			name: "TestIndexResetRankingRulesWithCustomClient",
			args: args{
				UID:    "indexUID",
				client: customMeili,
			},
			wantTask: &meilisearch.TaskInfo{
				TaskUID: 1,
			},
			wantResp: &defaultRankingRules,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			setUpIndexForFaceting(tt.args.client)
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
	meili := setup(t, "")
	customMeili := setup(t, "", meilisearch.WithCustomClientWithTLS(&tls.Config{
		InsecureSkipVerify: true,
	}))

	type args struct {
		UID    string
		client meilisearch.ServiceManager
	}
	tests := []struct {
		name     string
		args     args
		wantTask *meilisearch.TaskInfo
		wantResp *[]string
	}{
		{
			name: "TestIndexBasicResetSearchableAttributes",
			args: args{
				UID:    "indexUID",
				client: meili,
			},
			wantTask: &meilisearch.TaskInfo{
				TaskUID: 1,
			},
			wantResp: &[]string{"*"},
		},
		{
			name: "TestIndexResetSearchableAttributesWithCustomClient",
			args: args{
				UID:    "indexUID",
				client: customMeili,
			},
			wantTask: &meilisearch.TaskInfo{
				TaskUID: 1,
			},
			wantResp: &[]string{"*"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			setUpIndexForFaceting(tt.args.client)
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
	meili := setup(t, "")
	customMeili := setup(t, "", meilisearch.WithCustomClientWithTLS(&tls.Config{
		InsecureSkipVerify: true,
	}))

	type args struct {
		UID    string
		client meilisearch.ServiceManager
	}
	tests := []struct {
		name     string
		args     args
		wantTask *meilisearch.TaskInfo
		wantResp *meilisearch.Settings
	}{
		{
			name: "TestIndexBasicResetSettings",
			args: args{
				UID:    "indexUID",
				client: meili,
			},
			wantTask: &meilisearch.TaskInfo{
				TaskUID: 1,
			},
			wantResp: &meilisearch.Settings{
				RankingRules:         defaultRankingRules,
				DistinctAttribute:    (*string)(nil),
				SearchableAttributes: []string{"*"},
				DisplayedAttributes:  []string{"*"},
				StopWords:            []string{},
				Synonyms:             map[string][]string{},
				FilterableAttributes: []string{},
				SortableAttributes:   []string{},
				TypoTolerance:        &defaultTypoTolerance,
				Pagination:           &defaultPagination,
				Faceting:             &defaultFaceting,
				ProximityPrecision:   meilisearch.ByWord,
				SeparatorTokens:      make([]string, 0),
				NonSeparatorTokens:   make([]string, 0),
				Dictionary:           make([]string, 0),
				LocalizedAttributes:  nil,
				PrefixSearch:         stringPtr("indexingTime"),
				FacetSearch:          true,
				Embedders:            make(map[string]meilisearch.Embedder),
				Chat:                 &defaultChat,
			},
		},
		{
			name: "TestIndexResetSettingsWithCustomClient",
			args: args{
				UID:    "indexUID",
				client: customMeili,
			},
			wantTask: &meilisearch.TaskInfo{
				TaskUID: 1,
			},
			wantResp: &meilisearch.Settings{
				RankingRules:         defaultRankingRules,
				DistinctAttribute:    (*string)(nil),
				SearchableAttributes: []string{"*"},
				DisplayedAttributes:  []string{"*"},
				StopWords:            []string{},
				Synonyms:             map[string][]string{},
				FilterableAttributes: []string{},
				SortableAttributes:   []string{},
				TypoTolerance:        &defaultTypoTolerance,
				Pagination:           &defaultPagination,
				Faceting:             &defaultFaceting,
				ProximityPrecision:   meilisearch.ByWord,
				SeparatorTokens:      make([]string, 0),
				NonSeparatorTokens:   make([]string, 0),
				Dictionary:           make([]string, 0),
				LocalizedAttributes:  nil,
				PrefixSearch:         stringPtr("indexingTime"),
				FacetSearch:          true,
				Embedders:            make(map[string]meilisearch.Embedder),
				Chat:                 &defaultChat,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			setUpIndexForFaceting(tt.args.client)
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
	meili := setup(t, "")
	customMeili := setup(t, "", meilisearch.WithCustomClientWithTLS(&tls.Config{
		InsecureSkipVerify: true,
	}))

	type args struct {
		UID    string
		client meilisearch.ServiceManager
	}
	tests := []struct {
		name     string
		args     args
		wantTask *meilisearch.TaskInfo
	}{
		{
			name: "TestIndexBasicResetStopWords",
			args: args{
				UID:    "indexUID",
				client: meili,
			},
			wantTask: &meilisearch.TaskInfo{
				TaskUID: 1,
			},
		},
		{
			name: "TestIndexResetStopWordsWithCustomClient",
			args: args{
				UID:    "indexUID",
				client: customMeili,
			},
			wantTask: &meilisearch.TaskInfo{
				TaskUID: 1,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			setUpIndexForFaceting(tt.args.client)
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
	meili := setup(t, "")
	customMeili := setup(t, "", meilisearch.WithCustomClientWithTLS(&tls.Config{
		InsecureSkipVerify: true,
	}))

	type args struct {
		UID    string
		client meilisearch.ServiceManager
	}
	tests := []struct {
		name     string
		args     args
		wantTask *meilisearch.TaskInfo
	}{
		{
			name: "TestIndexBasicResetSynonyms",
			args: args{
				UID:    "indexUID",
				client: meili,
			},
			wantTask: &meilisearch.TaskInfo{
				TaskUID: 1,
			},
		},
		{
			name: "TestIndexResetSynonymsWithCustomClient",
			args: args{
				UID:    "indexUID",
				client: customMeili,
			},
			wantTask: &meilisearch.TaskInfo{
				TaskUID: 1,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			setUpIndexForFaceting(tt.args.client)
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
	meili := setup(t, "")
	customMeili := setup(t, "", meilisearch.WithCustomClientWithTLS(&tls.Config{
		InsecureSkipVerify: true,
	}))

	type args struct {
		UID    string
		client meilisearch.ServiceManager
	}
	tests := []struct {
		name     string
		args     args
		wantTask *meilisearch.TaskInfo
	}{
		{
			name: "TestIndexBasicResetSortableAttributes",
			args: args{
				UID:    "indexUID",
				client: meili,
			},
			wantTask: &meilisearch.TaskInfo{
				TaskUID: 1,
			},
		},
		{
			name: "TestIndexResetSortableAttributesCustomClient",
			args: args{
				UID:    "indexUID",
				client: customMeili,
			},
			wantTask: &meilisearch.TaskInfo{
				TaskUID: 1,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			setUpIndexForFaceting(tt.args.client)
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
	meili := setup(t, "")
	customMeili := setup(t, "", meilisearch.WithCustomClientWithTLS(&tls.Config{
		InsecureSkipVerify: true,
	}))

	type args struct {
		UID    string
		client meilisearch.ServiceManager
	}
	tests := []struct {
		name     string
		args     args
		wantTask *meilisearch.TaskInfo
		wantResp *meilisearch.TypoTolerance
	}{
		{
			name: "TestIndexBasicResetTypoTolerance",
			args: args{
				UID:    "indexUID",
				client: meili,
			},
			wantTask: &meilisearch.TaskInfo{
				TaskUID: 1,
			},
			wantResp: &defaultTypoTolerance,
		},
		{
			name: "TestIndexResetTypoToleranceWithCustomClient",
			args: args{
				UID:    "indexUID",
				client: customMeili,
			},
			wantTask: &meilisearch.TaskInfo{
				TaskUID: 1,
			},
			wantResp: &defaultTypoTolerance,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			setUpIndexForFaceting(tt.args.client)
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
	meili := setup(t, "")
	customMeili := setup(t, "", meilisearch.WithCustomClientWithTLS(&tls.Config{
		InsecureSkipVerify: true,
	}))

	type args struct {
		UID    string
		client meilisearch.ServiceManager
	}
	tests := []struct {
		name     string
		args     args
		wantTask *meilisearch.TaskInfo
		wantResp *meilisearch.Pagination
	}{
		{
			name: "TestIndexBasicResetPagination",
			args: args{
				UID:    "indexUID",
				client: meili,
			},
			wantTask: &meilisearch.TaskInfo{
				TaskUID: 1,
			},
			wantResp: &defaultPagination,
		},
		{
			name: "TestIndexResetPaginationWithCustomClient",
			args: args{
				UID:    "indexUID",
				client: customMeili,
			},
			wantTask: &meilisearch.TaskInfo{
				TaskUID: 1,
			},
			wantResp: &defaultPagination,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			setUpIndexForFaceting(tt.args.client)
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
	meili := setup(t, "")
	customMeili := setup(t, "", meilisearch.WithCustomClientWithTLS(&tls.Config{
		InsecureSkipVerify: true,
	}))

	type args struct {
		UID    string
		client meilisearch.ServiceManager
	}
	tests := []struct {
		name     string
		args     args
		wantTask *meilisearch.TaskInfo
		wantResp *meilisearch.Faceting
	}{
		{
			name: "TestIndexBasicResetFaceting",
			args: args{
				UID:    "indexUID",
				client: meili,
			},
			wantTask: &meilisearch.TaskInfo{
				TaskUID: 1,
			},
			wantResp: &defaultFaceting,
		},
		{
			name: "TestIndexResetFacetingWithCustomClient",
			args: args{
				UID:    "indexUID",
				client: customMeili,
			},
			wantTask: &meilisearch.TaskInfo{
				TaskUID: 1,
			},
			wantResp: &defaultFaceting,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			setUpIndexForFaceting(tt.args.client)
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
	meili := setup(t, "")
	customMeili := setup(t, "", meilisearch.WithCustomClientWithTLS(&tls.Config{
		InsecureSkipVerify: true,
	}))

	type args struct {
		UID     string
		client  meilisearch.ServiceManager
		request []interface{}
	}
	tests := []struct {
		name     string
		args     args
		wantTask *meilisearch.TaskInfo
	}{
		{
			name: "TestIndexBasicUpdateFilterableAttributes",
			args: args{
				UID:    "indexUID",
				client: meili,
				request: []interface{}{
					"title",
				},
			},
			wantTask: &meilisearch.TaskInfo{
				TaskUID: 1,
			},
		},
		{
			name: "TestIndexUpdateFilterableAttributesWithCustomClient",
			args: args{
				UID:    "indexUID",
				client: customMeili,
				request: []interface{}{
					"title",
				},
			},
			wantTask: &meilisearch.TaskInfo{
				TaskUID: 1,
			},
		},
		{
			name: "TestIndexUpdateFilterableAttributesMixedRawAndObject",
			args: args{
				UID:    "indexUID",
				client: meili,
				request: []interface{}{
					"tag",
					map[string]interface{}{
						"attributePatterns": []interface{}{"year"},
						"features": map[string]interface{}{
							"facetSearch": false,
							"filter": map[string]interface{}{
								"equality":   true,
								"comparison": true,
							},
						},
					},
				},
			},
			wantTask: &meilisearch.TaskInfo{
				TaskUID: 1,
			},
		},
		{
			name: "TestIndexUpdateFilterableAttributesOnlyObject",
			args: args{
				UID:    "indexUID",
				client: meili,
				request: []interface{}{
					map[string]interface{}{
						"attributePatterns": []interface{}{"year"},
						"features": map[string]interface{}{
							"facetSearch": false,
							"filter": map[string]interface{}{
								"equality":   true,
								"comparison": true,
							},
						},
					},
				},
			},
			wantTask: &meilisearch.TaskInfo{
				TaskUID: 1,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			setUpIndexForFaceting(tt.args.client)
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
	meili := setup(t, "")
	customMeili := setup(t, "", meilisearch.WithCustomClientWithTLS(&tls.Config{
		InsecureSkipVerify: true,
	}))

	type args struct {
		UID     string
		client  meilisearch.ServiceManager
		request []string
	}
	tests := []struct {
		name     string
		args     args
		wantTask *meilisearch.TaskInfo
		wantResp *[]string
	}{
		{
			name: "TestIndexBasicUpdateDisplayedAttributes",
			args: args{
				UID:    "indexUID",
				client: meili,
				request: []string{
					"book_id", "tag", "title",
				},
			},
			wantTask: &meilisearch.TaskInfo{
				TaskUID: 1,
			},
			wantResp: &[]string{"*"},
		},
		{
			name: "TestIndexUpdateDisplayedAttributesWithCustomClient",
			args: args{
				UID:    "indexUID",
				client: customMeili,
				request: []string{
					"book_id", "tag", "title",
				},
			},
			wantTask: &meilisearch.TaskInfo{
				TaskUID: 1,
			},
			wantResp: &[]string{"*"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			setUpIndexForFaceting(tt.args.client)
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
	meili := setup(t, "")
	customMeili := setup(t, "", meilisearch.WithCustomClientWithTLS(&tls.Config{
		InsecureSkipVerify: true,
	}))

	type args struct {
		UID     string
		client  meilisearch.ServiceManager
		request string
	}
	tests := []struct {
		name     string
		args     args
		wantTask *meilisearch.TaskInfo
	}{
		{
			name: "TestIndexBasicUpdateDistinctAttribute",
			args: args{
				UID:     "indexUID",
				client:  meili,
				request: "movie_id",
			},
			wantTask: &meilisearch.TaskInfo{
				TaskUID: 1,
			},
		},
		{
			name: "TestIndexUpdateDistinctAttributeWithCustomClient",
			args: args{
				UID:     "indexUID",
				client:  customMeili,
				request: "movie_id",
			},
			wantTask: &meilisearch.TaskInfo{
				TaskUID: 1,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			setUpIndexForFaceting(tt.args.client)
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
	meili := setup(t, "")
	customMeili := setup(t, "", meilisearch.WithCustomClientWithTLS(&tls.Config{
		InsecureSkipVerify: true,
	}))

	type args struct {
		UID     string
		client  meilisearch.ServiceManager
		request []string
	}
	tests := []struct {
		name     string
		args     args
		wantTask *meilisearch.TaskInfo
		wantResp *[]string
	}{
		{
			name: "TestIndexBasicUpdateRankingRules",
			args: args{
				UID:    "indexUID",
				client: meili,
				request: []string{
					"typo", "words",
				},
			},
			wantTask: &meilisearch.TaskInfo{
				TaskUID: 1,
			},
			wantResp: &defaultRankingRules,
		},
		{
			name: "TestIndexUpdateRankingRulesWithCustomClient",
			args: args{
				UID:    "indexUID",
				client: customMeili,
				request: []string{
					"typo", "words",
				},
			},
			wantTask: &meilisearch.TaskInfo{
				TaskUID: 1,
			},
			wantResp: &defaultRankingRules,
		},
		{
			name: "TestIndexUpdateRankingRulesAscending",
			args: args{
				UID:    "indexUID",
				client: meili,
				request: []string{
					"BookID:asc",
				},
			},
			wantTask: &meilisearch.TaskInfo{
				TaskUID: 1,
			},
			wantResp: &defaultRankingRules,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			setUpIndexForFaceting(tt.args.client)
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
	meili := setup(t, "")
	customMeili := setup(t, "", meilisearch.WithCustomClientWithTLS(&tls.Config{
		InsecureSkipVerify: true,
	}))

	type args struct {
		UID     string
		client  meilisearch.ServiceManager
		request []string
	}
	tests := []struct {
		name     string
		args     args
		wantTask *meilisearch.TaskInfo
		wantResp *[]string
	}{
		{
			name: "TestIndexBasicUpdateSearchableAttributes",
			args: args{
				UID:    "indexUID",
				client: meili,
				request: []string{
					"title", "tag",
				},
			},
			wantTask: &meilisearch.TaskInfo{
				TaskUID: 1,
			},
			wantResp: &[]string{"*"},
		},
		{
			name: "TestIndexUpdateSearchableAttributesWithCustomClient",
			args: args{
				UID:    "indexUID",
				client: customMeili,
				request: []string{
					"title", "tag",
				},
			},
			wantTask: &meilisearch.TaskInfo{
				TaskUID: 1,
			},
			wantResp: &[]string{"*"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			setUpIndexForFaceting(tt.args.client)
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
	meili := setup(t, "")
	customMeili := setup(t, "", meilisearch.WithCustomClientWithTLS(&tls.Config{
		InsecureSkipVerify: true,
	}))

	type args struct {
		UID     string
		client  meilisearch.ServiceManager
		request meilisearch.Settings
	}
	tests := []struct {
		name     string
		args     args
		wantTask *meilisearch.TaskInfo
		wantResp *meilisearch.Settings
	}{
		{
			name: "TestIndexBasicUpdateSettings",
			args: args{
				UID:    "indexUID",
				client: meili,
				request: meilisearch.Settings{
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
					TypoTolerance: &meilisearch.TypoTolerance{
						Enabled: true,
						MinWordSizeForTypos: meilisearch.MinWordSizeForTypos{
							OneTypo:  7,
							TwoTypos: 10,
						},
						DisableOnWords:      []string{},
						DisableOnAttributes: []string{},
						DisableOnNumbers:    true,
					},
					Pagination: &meilisearch.Pagination{
						MaxTotalHits: 1200,
					},
					Faceting: &meilisearch.Faceting{
						MaxValuesPerFacet: 200,
						SortFacetValuesBy: map[string]meilisearch.SortFacetType{
							"*": meilisearch.SortFacetTypeAlpha,
						},
					},
					SearchCutoffMs:     150,
					ProximityPrecision: meilisearch.ByAttribute,
					SeparatorTokens:    make([]string, 0),
					NonSeparatorTokens: make([]string, 0),
					Dictionary:         make([]string, 0),
					LocalizedAttributes: []*meilisearch.LocalizedAttributes{
						{
							Locales:           []string{"jpn", "eng"},
							AttributePatterns: []string{"*_ja"},
						},
					},
					PrefixSearch: stringPtr("indexingTime"),
					FacetSearch:  true,
					Embedders:    make(map[string]meilisearch.Embedder),
					Chat:         &defaultChat,
				},
			},
			wantTask: &meilisearch.TaskInfo{
				TaskUID: 1,
			},
			wantResp: &meilisearch.Settings{
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
				SearchCutoffMs:       150,
				ProximityPrecision:   meilisearch.ByAttribute,
				SeparatorTokens:      make([]string, 0),
				NonSeparatorTokens:   make([]string, 0),
				Dictionary:           make([]string, 0),
				LocalizedAttributes: []*meilisearch.LocalizedAttributes{
					{
						Locales:           []string{"jpn", "eng"},
						AttributePatterns: []string{"*_ja"},
					},
				},
				PrefixSearch: stringPtr("indexingTime"),
				FacetSearch:  true,
				Embedders:    make(map[string]meilisearch.Embedder),
				Chat:         &defaultChat,
			},
		},
		{
			name: "TestIndexUpdateSettingsWithCustomClient",
			args: args{
				UID:    "indexUID",
				client: customMeili,
				request: meilisearch.Settings{
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
					TypoTolerance: &meilisearch.TypoTolerance{
						Enabled: true,
						MinWordSizeForTypos: meilisearch.MinWordSizeForTypos{
							OneTypo:  7,
							TwoTypos: 10,
						},
						DisableOnWords:      []string{},
						DisableOnAttributes: []string{},
						DisableOnNumbers:    true,
					},
					Pagination: &meilisearch.Pagination{
						MaxTotalHits: 1200,
					},
					Faceting: &meilisearch.Faceting{
						MaxValuesPerFacet: 200,
						SortFacetValuesBy: map[string]meilisearch.SortFacetType{
							"*": meilisearch.SortFacetTypeAlpha,
						},
					},
					SearchCutoffMs:     150,
					ProximityPrecision: meilisearch.ByWord,
					SeparatorTokens:    make([]string, 0),
					NonSeparatorTokens: make([]string, 0),
					Dictionary:         make([]string, 0),
					PrefixSearch:       stringPtr("indexingTime"),
					FacetSearch:        true,
					Embedders:          make(map[string]meilisearch.Embedder),
					Chat:               &defaultChat,
				},
			},
			wantTask: &meilisearch.TaskInfo{
				TaskUID: 1,
			},
			wantResp: &meilisearch.Settings{
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
				SearchCutoffMs:       150,
				ProximityPrecision:   meilisearch.ByWord,
				SeparatorTokens:      make([]string, 0),
				NonSeparatorTokens:   make([]string, 0),
				Dictionary:           make([]string, 0),
				PrefixSearch:         stringPtr("indexingTime"),
				FacetSearch:          true,
				Embedders:            make(map[string]meilisearch.Embedder),
				Chat:                 &defaultChat,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			setUpIndexForFaceting(tt.args.client)
			c := tt.args.client
			i := c.Index(tt.args.UID)
			t.Cleanup(cleanup(c))

			gotTask, err := i.UpdateSettings(&tt.args.request)
			require.NoError(t, err)
			require.GreaterOrEqual(t, gotTask.TaskUID, tt.wantTask.TaskUID)
			testWaitForTask(t, i, gotTask)

			gotResp, err := i.GetSettings()
			require.NoError(t, err)
			require.Equal(t, &tt.args.request, gotResp)
		})
	}
}

func TestIndex_UpdateSettingsOneByOne(t *testing.T) {
	meili := setup(t, "")
	customMeili := setup(t, "", meilisearch.WithCustomClientWithTLS(&tls.Config{
		InsecureSkipVerify: true,
	}))

	type args struct {
		UID            string
		client         meilisearch.ServiceManager
		firstRequest   meilisearch.Settings
		firstResponse  meilisearch.Settings
		secondRequest  meilisearch.Settings
		secondResponse meilisearch.Settings
	}
	tests := []struct {
		name     string
		args     args
		wantTask *meilisearch.TaskInfo
		wantResp *meilisearch.Settings
	}{
		{
			name: "TestIndexUpdateJustSynonyms",
			args: args{
				UID:    "indexUID",
				client: meili,
				firstRequest: meilisearch.Settings{
					RankingRules: []string{
						"typo", "words",
					},
					Synonyms: map[string][]string{
						"hp": {"harry potter"},
					},
				},
				firstResponse: meilisearch.Settings{
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
					ProximityPrecision:   meilisearch.ByWord,
					SeparatorTokens:      make([]string, 0),
					NonSeparatorTokens:   make([]string, 0),
					Dictionary:           make([]string, 0),
					PrefixSearch:         stringPtr("indexingTime"),
					FacetSearch:          true,
					Embedders:            make(map[string]meilisearch.Embedder),
					Chat:                 &defaultChat,
				},
				secondRequest: meilisearch.Settings{
					Synonyms: map[string][]string{
						"al": {"alice"},
					},
				},
				secondResponse: meilisearch.Settings{
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
					ProximityPrecision:   meilisearch.ByWord,
					SeparatorTokens:      make([]string, 0),
					NonSeparatorTokens:   make([]string, 0),
					Dictionary:           make([]string, 0),
					PrefixSearch:         stringPtr("indexingTime"),
					FacetSearch:          true,
					Embedders:            make(map[string]meilisearch.Embedder),
					Chat:                 &defaultChat,
				},
			},
			wantTask: &meilisearch.TaskInfo{
				TaskUID: 1,
			},
			wantResp: &meilisearch.Settings{
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
				ProximityPrecision:   meilisearch.ByWord,
				SeparatorTokens:      make([]string, 0),
				NonSeparatorTokens:   make([]string, 0),
				Dictionary:           make([]string, 0),
				PrefixSearch:         stringPtr("indexingTime"),
				FacetSearch:          true,
				Embedders:            make(map[string]meilisearch.Embedder),
				Chat:                 &defaultChat,
			},
		},
		{
			name: "TestIndexUpdateJustSynonymsWithCustomClient",
			args: args{
				UID:    "indexUID",
				client: customMeili,
				firstRequest: meilisearch.Settings{
					RankingRules: []string{
						"typo", "words",
					},
					Synonyms: map[string][]string{
						"hp": {"harry potter"},
					},
				},
				firstResponse: meilisearch.Settings{
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
					ProximityPrecision:   meilisearch.ByWord,
					SeparatorTokens:      make([]string, 0),
					NonSeparatorTokens:   make([]string, 0),
					Dictionary:           make([]string, 0),
					PrefixSearch:         stringPtr("indexingTime"),
					FacetSearch:          true,
					Embedders:            make(map[string]meilisearch.Embedder),
					Chat:                 &defaultChat,
				},
				secondRequest: meilisearch.Settings{
					Synonyms: map[string][]string{
						"al": {"alice"},
					},
				},
				secondResponse: meilisearch.Settings{
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
					ProximityPrecision:   meilisearch.ByWord,
					SeparatorTokens:      make([]string, 0),
					NonSeparatorTokens:   make([]string, 0),
					Dictionary:           make([]string, 0),
					PrefixSearch:         stringPtr("indexingTime"),
					FacetSearch:          true,
					Embedders:            make(map[string]meilisearch.Embedder),
					Chat:                 &defaultChat,
				},
			},
			wantTask: &meilisearch.TaskInfo{
				TaskUID: 1,
			},
			wantResp: &meilisearch.Settings{
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
				ProximityPrecision:   meilisearch.ByWord,
				SeparatorTokens:      make([]string, 0),
				NonSeparatorTokens:   make([]string, 0),
				Dictionary:           make([]string, 0),
				PrefixSearch:         stringPtr("indexingTime"),
				FacetSearch:          true,
				Embedders:            make(map[string]meilisearch.Embedder),
				Chat:                 &defaultChat,
			},
		},
		{
			name: "TestIndexUpdateJustSearchableAttributes",
			args: args{
				UID:    "indexUID",
				client: meili,
				firstRequest: meilisearch.Settings{
					RankingRules: []string{
						"typo", "words",
					},
					SearchableAttributes: []string{
						"tag",
					},
				},
				firstResponse: meilisearch.Settings{
					RankingRules: []string{
						"typo", "words",
					},
					DistinctAttribute: (*string)(nil),
					SearchableAttributes: []string{
						"tag",
					},
					DisplayedAttributes:  []string{"*"},
					StopWords:            []string{},
					Synonyms:             make(map[string][]string),
					FilterableAttributes: []string{},
					SortableAttributes:   []string{},
					TypoTolerance:        &defaultTypoTolerance,
					Pagination:           &defaultPagination,
					Faceting:             &defaultFaceting,
					ProximityPrecision:   meilisearch.ByWord,
					SeparatorTokens:      make([]string, 0),
					NonSeparatorTokens:   make([]string, 0),
					Dictionary:           make([]string, 0),
					PrefixSearch:         stringPtr("indexingTime"),
					FacetSearch:          true,
					Embedders:            make(map[string]meilisearch.Embedder),
					Chat:                 &defaultChat,
				},
				secondRequest: meilisearch.Settings{
					SearchableAttributes: []string{
						"title",
					},
				},
				secondResponse: meilisearch.Settings{
					RankingRules: []string{
						"typo", "words",
					},
					DistinctAttribute: (*string)(nil),
					SearchableAttributes: []string{
						"title",
					},
					DisplayedAttributes:  []string{"*"},
					StopWords:            []string{},
					Synonyms:             make(map[string][]string),
					FilterableAttributes: []string{},
					SortableAttributes:   []string{},
					TypoTolerance:        &defaultTypoTolerance,
					Pagination:           &defaultPagination,
					Faceting:             &defaultFaceting,
					ProximityPrecision:   meilisearch.ByWord,
					SeparatorTokens:      make([]string, 0),
					NonSeparatorTokens:   make([]string, 0),
					Dictionary:           make([]string, 0),
					PrefixSearch:         stringPtr("indexingTime"),
					FacetSearch:          true,
					Embedders:            make(map[string]meilisearch.Embedder),
					Chat:                 &defaultChat,
				},
			},
			wantTask: &meilisearch.TaskInfo{
				TaskUID: 1,
			},
			wantResp: &meilisearch.Settings{
				RankingRules:         defaultRankingRules,
				DistinctAttribute:    (*string)(nil),
				SearchableAttributes: []string{"*"},
				DisplayedAttributes:  []string{"*"},
				StopWords:            []string{},
				Synonyms:             make(map[string][]string),
				FilterableAttributes: []string{},
				SortableAttributes:   []string{},
				TypoTolerance:        &defaultTypoTolerance,
				Pagination:           &defaultPagination,
				Faceting:             &defaultFaceting,
				ProximityPrecision:   meilisearch.ByWord,
				SeparatorTokens:      make([]string, 0),
				NonSeparatorTokens:   make([]string, 0),
				Dictionary:           make([]string, 0),
				PrefixSearch:         stringPtr("indexingTime"),
				FacetSearch:          true,
				Embedders:            make(map[string]meilisearch.Embedder),
				Chat:                 &defaultChat,
			},
		},
		{
			name: "TestIndexUpdateJustDisplayedAttributes",
			args: args{
				UID:    "indexUID",
				client: meili,
				firstRequest: meilisearch.Settings{
					RankingRules: []string{
						"typo", "words",
					},
					DisplayedAttributes: []string{
						"book_id", "tag", "title",
					},
				},
				firstResponse: meilisearch.Settings{
					RankingRules: []string{
						"typo", "words",
					},
					DistinctAttribute:    (*string)(nil),
					SearchableAttributes: []string{"*"},
					DisplayedAttributes: []string{
						"book_id", "tag", "title",
					},
					StopWords:            []string{},
					Synonyms:             make(map[string][]string),
					FilterableAttributes: []string{},
					SortableAttributes:   []string{},
					TypoTolerance:        &defaultTypoTolerance,
					Pagination:           &defaultPagination,
					Faceting:             &defaultFaceting,
					ProximityPrecision:   meilisearch.ByWord,
					SeparatorTokens:      make([]string, 0),
					NonSeparatorTokens:   make([]string, 0),
					Dictionary:           make([]string, 0),
					PrefixSearch:         stringPtr("indexingTime"),
					FacetSearch:          true,
					Embedders:            make(map[string]meilisearch.Embedder),
					Chat:                 &defaultChat,
				},
				secondRequest: meilisearch.Settings{
					DisplayedAttributes: []string{
						"book_id", "tag",
					},
				},
				secondResponse: meilisearch.Settings{
					RankingRules: []string{
						"typo", "words",
					},
					DistinctAttribute:    (*string)(nil),
					SearchableAttributes: []string{"*"},
					DisplayedAttributes: []string{
						"book_id", "tag",
					},
					StopWords:            []string{},
					Synonyms:             make(map[string][]string),
					FilterableAttributes: []string{},
					SortableAttributes:   []string{},
					TypoTolerance:        &defaultTypoTolerance,
					Pagination:           &defaultPagination,
					Faceting:             &defaultFaceting,
					ProximityPrecision:   meilisearch.ByWord,
					SeparatorTokens:      make([]string, 0),
					NonSeparatorTokens:   make([]string, 0),
					Dictionary:           make([]string, 0),
					PrefixSearch:         stringPtr("indexingTime"),
					FacetSearch:          true,
					Embedders:            make(map[string]meilisearch.Embedder),
					Chat:                 &defaultChat,
				},
			},
			wantTask: &meilisearch.TaskInfo{
				TaskUID: 1,
			},
			wantResp: &meilisearch.Settings{
				RankingRules:         defaultRankingRules,
				DistinctAttribute:    (*string)(nil),
				SearchableAttributes: []string{"*"},
				DisplayedAttributes:  []string{"*"},
				StopWords:            []string{},
				Synonyms:             make(map[string][]string),
				FilterableAttributes: []string{},
				SortableAttributes:   []string{},
				TypoTolerance:        &defaultTypoTolerance,
				Pagination:           &defaultPagination,
				Faceting:             &defaultFaceting,
				ProximityPrecision:   meilisearch.ByWord,
				SeparatorTokens:      make([]string, 0),
				NonSeparatorTokens:   make([]string, 0),
				Dictionary:           make([]string, 0),
				PrefixSearch:         stringPtr("indexingTime"),
				FacetSearch:          true,
				Embedders:            make(map[string]meilisearch.Embedder),
				Chat:                 &defaultChat,
			},
		},
		{
			name: "TestIndexUpdateJustStopWords",
			args: args{
				UID:    "indexUID",
				client: meili,
				firstRequest: meilisearch.Settings{
					RankingRules: []string{
						"typo", "words",
					},
					StopWords: []string{
						"of", "the",
					},
				},
				firstResponse: meilisearch.Settings{
					RankingRules: []string{
						"typo", "words",
					},
					DistinctAttribute:    (*string)(nil),
					SearchableAttributes: []string{"*"},
					DisplayedAttributes:  []string{"*"},
					StopWords: []string{
						"of", "the",
					},
					Synonyms:             make(map[string][]string),
					FilterableAttributes: []string{},
					SortableAttributes:   []string{},
					TypoTolerance:        &defaultTypoTolerance,
					Pagination:           &defaultPagination,
					Faceting:             &defaultFaceting,
					ProximityPrecision:   meilisearch.ByWord,
					SeparatorTokens:      make([]string, 0),
					NonSeparatorTokens:   make([]string, 0),
					Dictionary:           make([]string, 0),
					PrefixSearch:         stringPtr("indexingTime"),
					FacetSearch:          true,
					Embedders:            make(map[string]meilisearch.Embedder),
					Chat:                 &defaultChat,
				},
				secondRequest: meilisearch.Settings{
					StopWords: []string{
						"of", "the",
					},
				},
				secondResponse: meilisearch.Settings{
					RankingRules: []string{
						"typo", "words",
					},
					DistinctAttribute:    (*string)(nil),
					SearchableAttributes: []string{"*"},
					DisplayedAttributes:  []string{"*"},
					StopWords: []string{
						"of", "the",
					},
					Synonyms:             make(map[string][]string),
					FilterableAttributes: []string{},
					SortableAttributes:   []string{},
					TypoTolerance:        &defaultTypoTolerance,
					Pagination:           &defaultPagination,
					Faceting:             &defaultFaceting,
					ProximityPrecision:   meilisearch.ByWord,
					SeparatorTokens:      make([]string, 0),
					NonSeparatorTokens:   make([]string, 0),
					Dictionary:           make([]string, 0),
					PrefixSearch:         stringPtr("indexingTime"),
					FacetSearch:          true,
					Embedders:            make(map[string]meilisearch.Embedder),
					Chat:                 &defaultChat,
				},
			},
			wantTask: &meilisearch.TaskInfo{
				TaskUID: 1,
			},
			wantResp: &meilisearch.Settings{
				RankingRules:         defaultRankingRules,
				DistinctAttribute:    (*string)(nil),
				SearchableAttributes: []string{"*"},
				DisplayedAttributes:  []string{"*"},
				StopWords:            []string{},
				Synonyms:             make(map[string][]string),
				FilterableAttributes: []string{},
				SortableAttributes:   []string{},
				TypoTolerance:        &defaultTypoTolerance,
				Pagination:           &defaultPagination,
				Faceting:             &defaultFaceting,
				ProximityPrecision:   meilisearch.ByWord,
				SeparatorTokens:      make([]string, 0),
				NonSeparatorTokens:   make([]string, 0),
				Dictionary:           make([]string, 0),
				PrefixSearch:         stringPtr("indexingTime"),
				FacetSearch:          true,
				Embedders:            make(map[string]meilisearch.Embedder),
				Chat:                 &defaultChat,
			},
		},
		{
			name: "TestIndexUpdateJustFilterableAttributes",
			args: args{
				UID:    "indexUID",
				client: meili,
				firstRequest: meilisearch.Settings{
					RankingRules: []string{
						"typo", "words",
					},
					FilterableAttributes: []string{
						"title",
					},
				},
				firstResponse: meilisearch.Settings{
					RankingRules: []string{
						"typo", "words",
					},
					DistinctAttribute:    (*string)(nil),
					SearchableAttributes: []string{"*"},
					DisplayedAttributes:  []string{"*"},
					StopWords:            []string{},
					Synonyms:             make(map[string][]string),
					FilterableAttributes: []string{
						"title",
					},
					SortableAttributes: []string{},
					TypoTolerance:      &defaultTypoTolerance,
					Pagination:         &defaultPagination,
					Faceting:           &defaultFaceting,
					ProximityPrecision: meilisearch.ByWord,
					SeparatorTokens:    make([]string, 0),
					NonSeparatorTokens: make([]string, 0),
					Dictionary:         make([]string, 0),
					PrefixSearch:       stringPtr("indexingTime"),
					FacetSearch:        true,
					Embedders:          make(map[string]meilisearch.Embedder),
					Chat:               &defaultChat,
				},
				secondRequest: meilisearch.Settings{
					FilterableAttributes: []string{
						"title",
					},
				},
				secondResponse: meilisearch.Settings{
					RankingRules: []string{
						"typo", "words",
					},
					DistinctAttribute:    (*string)(nil),
					SearchableAttributes: []string{"*"},
					DisplayedAttributes:  []string{"*"},
					StopWords:            []string{},
					Synonyms:             make(map[string][]string),
					FilterableAttributes: []string{
						"title",
					},
					SortableAttributes: []string{},
					TypoTolerance:      &defaultTypoTolerance,
					Pagination:         &defaultPagination,
					Faceting:           &defaultFaceting,
					ProximityPrecision: meilisearch.ByWord,
					SeparatorTokens:    make([]string, 0),
					NonSeparatorTokens: make([]string, 0),
					Dictionary:         make([]string, 0),
					PrefixSearch:       stringPtr("indexingTime"),
					FacetSearch:        true,
					Embedders:          make(map[string]meilisearch.Embedder),
					Chat:               &defaultChat,
				},
			},
			wantTask: &meilisearch.TaskInfo{
				TaskUID: 1,
			},
			wantResp: &meilisearch.Settings{
				RankingRules:         defaultRankingRules,
				DistinctAttribute:    (*string)(nil),
				SearchableAttributes: []string{"*"},
				DisplayedAttributes:  []string{"*"},
				StopWords:            []string{},
				Synonyms:             make(map[string][]string),
				FilterableAttributes: []string{},
				SortableAttributes:   []string{},
				TypoTolerance:        &defaultTypoTolerance,
				Pagination:           &defaultPagination,
				Faceting:             &defaultFaceting,
				ProximityPrecision:   meilisearch.ByWord,
				SeparatorTokens:      make([]string, 0),
				NonSeparatorTokens:   make([]string, 0),
				Dictionary:           make([]string, 0),
				PrefixSearch:         stringPtr("indexingTime"),
				FacetSearch:          true,
				Embedders:            make(map[string]meilisearch.Embedder),
				Chat:                 &defaultChat,
			},
		},
		{
			name: "TestIndexUpdateJustSortableAttributes",
			args: args{
				UID:    "indexUID",
				client: meili,
				firstRequest: meilisearch.Settings{
					RankingRules: []string{
						"typo", "words",
					},
					SortableAttributes: []string{
						"title",
					},
				},
				firstResponse: meilisearch.Settings{
					RankingRules: []string{
						"typo", "words",
					},
					DistinctAttribute:    (*string)(nil),
					SearchableAttributes: []string{"*"},
					DisplayedAttributes:  []string{"*"},
					StopWords:            []string{},
					Synonyms:             make(map[string][]string),
					FilterableAttributes: []string{},
					SortableAttributes: []string{
						"title",
					},
					TypoTolerance:      &defaultTypoTolerance,
					Pagination:         &defaultPagination,
					Faceting:           &defaultFaceting,
					ProximityPrecision: meilisearch.ByWord,
					SeparatorTokens:    make([]string, 0),
					NonSeparatorTokens: make([]string, 0),
					Dictionary:         make([]string, 0),
					PrefixSearch:       stringPtr("indexingTime"),
					FacetSearch:        true,
					Embedders:          make(map[string]meilisearch.Embedder),
					Chat:               &defaultChat,
				},
				secondRequest: meilisearch.Settings{
					SortableAttributes: []string{
						"title",
					},
				},
				secondResponse: meilisearch.Settings{
					RankingRules: []string{
						"typo", "words",
					},
					DistinctAttribute:    (*string)(nil),
					SearchableAttributes: []string{"*"},
					DisplayedAttributes:  []string{"*"},
					StopWords:            []string{},
					Synonyms:             make(map[string][]string),
					FilterableAttributes: []string{},
					SortableAttributes: []string{
						"title",
					},
					TypoTolerance:      &defaultTypoTolerance,
					Pagination:         &defaultPagination,
					Faceting:           &defaultFaceting,
					ProximityPrecision: meilisearch.ByWord,
					SeparatorTokens:    make([]string, 0),
					NonSeparatorTokens: make([]string, 0),
					Dictionary:         make([]string, 0),
					PrefixSearch:       stringPtr("indexingTime"),
					FacetSearch:        true,
					Embedders:          make(map[string]meilisearch.Embedder),
					Chat:               &defaultChat,
				},
			},
			wantTask: &meilisearch.TaskInfo{
				TaskUID: 1,
			},
			wantResp: &meilisearch.Settings{
				RankingRules:         defaultRankingRules,
				DistinctAttribute:    (*string)(nil),
				SearchableAttributes: []string{"*"},
				DisplayedAttributes:  []string{"*"},
				StopWords:            []string{},
				Synonyms:             make(map[string][]string),
				FilterableAttributes: []string{},
				SortableAttributes:   []string{},
				TypoTolerance:        &defaultTypoTolerance,
				Pagination:           &defaultPagination,
				Faceting:             &defaultFaceting,
				ProximityPrecision:   meilisearch.ByWord,
				SeparatorTokens:      make([]string, 0),
				NonSeparatorTokens:   make([]string, 0),
				Dictionary:           make([]string, 0),
				PrefixSearch:         stringPtr("indexingTime"),
				FacetSearch:          true,
				Embedders:            make(map[string]meilisearch.Embedder),
				Chat:                 &defaultChat,
			},
		},
		{
			name: "TestIndexUpdateJustTypoTolerance",
			args: args{
				UID:    "indexUID",
				client: meili,
				firstRequest: meilisearch.Settings{
					RankingRules: []string{
						"typo", "words",
					},
					TypoTolerance: &meilisearch.TypoTolerance{
						Enabled: true,
						MinWordSizeForTypos: meilisearch.MinWordSizeForTypos{
							OneTypo:  7,
							TwoTypos: 10,
						},
						DisableOnWords:      []string{},
						DisableOnAttributes: []string{},
						DisableOnNumbers:    false,
					},
				},
				firstResponse: meilisearch.Settings{
					RankingRules: []string{
						"typo", "words",
					},
					DistinctAttribute:    (*string)(nil),
					SearchableAttributes: []string{"*"},
					DisplayedAttributes:  []string{"*"},
					StopWords:            []string{},
					Synonyms:             make(map[string][]string),
					TypoTolerance: &meilisearch.TypoTolerance{
						Enabled: true,
						MinWordSizeForTypos: meilisearch.MinWordSizeForTypos{
							OneTypo:  7,
							TwoTypos: 10,
						},
						DisableOnWords:      []string{},
						DisableOnAttributes: []string{},
						DisableOnNumbers:    false,
					},
					FilterableAttributes: []string{},
					SortableAttributes:   []string{},
					Pagination:           &defaultPagination,
					Faceting:             &defaultFaceting,
					ProximityPrecision:   meilisearch.ByWord,
					SeparatorTokens:      make([]string, 0),
					NonSeparatorTokens:   make([]string, 0),
					Dictionary:           make([]string, 0),
					PrefixSearch:         stringPtr("indexingTime"),
					FacetSearch:          true,
					Embedders:            make(map[string]meilisearch.Embedder),
					Chat:                 &defaultChat,
				},
				secondRequest: meilisearch.Settings{
					TypoTolerance: &meilisearch.TypoTolerance{
						Enabled: true,
						MinWordSizeForTypos: meilisearch.MinWordSizeForTypos{
							OneTypo:  7,
							TwoTypos: 10,
						},
						DisableOnWords: []string{
							"and",
						},
						DisableOnAttributes: []string{
							"year",
						},
						DisableOnNumbers: true,
					},
				},
				secondResponse: meilisearch.Settings{
					RankingRules: []string{
						"typo", "words",
					},
					DistinctAttribute:    (*string)(nil),
					SearchableAttributes: []string{"*"},
					DisplayedAttributes:  []string{"*"},
					StopWords:            []string{},
					Synonyms:             make(map[string][]string),
					FilterableAttributes: []string{},
					SortableAttributes:   []string{},
					TypoTolerance: &meilisearch.TypoTolerance{
						Enabled: true,
						MinWordSizeForTypos: meilisearch.MinWordSizeForTypos{
							OneTypo:  7,
							TwoTypos: 10,
						},
						DisableOnWords: []string{
							"and",
						},
						DisableOnAttributes: []string{
							"year",
						},
						DisableOnNumbers: true,
					},
					Pagination:         &defaultPagination,
					Faceting:           &defaultFaceting,
					ProximityPrecision: meilisearch.ByWord,
					SeparatorTokens:    make([]string, 0),
					NonSeparatorTokens: make([]string, 0),
					Dictionary:         make([]string, 0),
					PrefixSearch:       stringPtr("indexingTime"),
					FacetSearch:        true,
					Embedders:          make(map[string]meilisearch.Embedder),
					Chat:               &defaultChat,
				},
			},
			wantTask: &meilisearch.TaskInfo{
				TaskUID: 1,
			},
			wantResp: &meilisearch.Settings{
				RankingRules:         defaultRankingRules,
				DistinctAttribute:    (*string)(nil),
				SearchableAttributes: []string{"*"},
				DisplayedAttributes:  []string{"*"},
				StopWords:            []string{},
				Synonyms:             make(map[string][]string),
				FilterableAttributes: []string{},
				SortableAttributes:   []string{},
				TypoTolerance:        &defaultTypoTolerance,
				Pagination:           &defaultPagination,
				Faceting:             &defaultFaceting,
				ProximityPrecision:   meilisearch.ByWord,
				SeparatorTokens:      make([]string, 0),
				NonSeparatorTokens:   make([]string, 0),
				Dictionary:           make([]string, 0),
				PrefixSearch:         stringPtr("indexingTime"),
				FacetSearch:          true,
				Embedders:            make(map[string]meilisearch.Embedder),
				Chat:                 &defaultChat,
			},
		},
		{
			name: "TestIndexUpdateJustPagination",
			args: args{
				UID:    "indexUID",
				client: meili,
				firstRequest: meilisearch.Settings{
					RankingRules: []string{
						"typo", "words",
					},
					Pagination: &meilisearch.Pagination{
						MaxTotalHits: 1200,
					},
				},
				firstResponse: meilisearch.Settings{
					RankingRules: []string{
						"typo", "words",
					},
					DistinctAttribute:    (*string)(nil),
					SearchableAttributes: []string{"*"},
					DisplayedAttributes:  []string{"*"},
					StopWords:            []string{},
					Synonyms:             make(map[string][]string),
					FilterableAttributes: []string{},
					SortableAttributes:   []string{},
					TypoTolerance:        &defaultTypoTolerance,
					Pagination: &meilisearch.Pagination{
						MaxTotalHits: 1200,
					},
					Faceting:           &defaultFaceting,
					ProximityPrecision: meilisearch.ByWord,
					SeparatorTokens:    make([]string, 0),
					NonSeparatorTokens: make([]string, 0),
					Dictionary:         make([]string, 0),
					PrefixSearch:       stringPtr("indexingTime"),
					FacetSearch:        true,
					Embedders:          make(map[string]meilisearch.Embedder),
					Chat:               &defaultChat,
				},
				secondRequest: meilisearch.Settings{
					Pagination: &meilisearch.Pagination{
						MaxTotalHits: 1200,
					},
				},
				secondResponse: meilisearch.Settings{
					RankingRules: []string{
						"typo", "words",
					},
					DistinctAttribute:    (*string)(nil),
					SearchableAttributes: []string{"*"},
					DisplayedAttributes:  []string{"*"},
					StopWords:            []string{},
					Synonyms:             make(map[string][]string),
					FilterableAttributes: []string{},
					SortableAttributes:   []string{},
					TypoTolerance:        &defaultTypoTolerance,
					Pagination: &meilisearch.Pagination{
						MaxTotalHits: 1200,
					},
					Faceting:           &defaultFaceting,
					ProximityPrecision: meilisearch.ByWord,
					SeparatorTokens:    make([]string, 0),
					NonSeparatorTokens: make([]string, 0),
					Dictionary:         make([]string, 0),
					PrefixSearch:       stringPtr("indexingTime"),
					FacetSearch:        true,
					Embedders:          make(map[string]meilisearch.Embedder),
					Chat:               &defaultChat,
				},
			},
			wantTask: &meilisearch.TaskInfo{
				TaskUID: 1,
			},
			wantResp: &meilisearch.Settings{
				RankingRules:         defaultRankingRules,
				DistinctAttribute:    (*string)(nil),
				SearchableAttributes: []string{"*"},
				DisplayedAttributes:  []string{"*"},
				StopWords:            []string{},
				Synonyms:             make(map[string][]string),
				FilterableAttributes: []string{},
				SortableAttributes:   []string{},
				TypoTolerance:        &defaultTypoTolerance,
				Pagination:           &defaultPagination,
				Faceting:             &defaultFaceting,
				ProximityPrecision:   meilisearch.ByWord,
				SeparatorTokens:      make([]string, 0),
				NonSeparatorTokens:   make([]string, 0),
				Dictionary:           make([]string, 0),
				PrefixSearch:         stringPtr("indexingTime"),
				FacetSearch:          true,
				Embedders:            make(map[string]meilisearch.Embedder),
				Chat:                 &defaultChat,
			},
		},
		{
			name: "TestIndexUpdateJustFaceting",
			args: args{
				UID:    "indexUID",
				client: meili,
				firstRequest: meilisearch.Settings{
					RankingRules: []string{
						"typo", "words",
					},
					Faceting: &meilisearch.Faceting{
						MaxValuesPerFacet: 200,
					},
				},
				firstResponse: meilisearch.Settings{
					RankingRules: []string{
						"typo", "words",
					},
					DistinctAttribute:    (*string)(nil),
					SearchableAttributes: []string{"*"},
					DisplayedAttributes:  []string{"*"},
					StopWords:            []string{},
					Synonyms:             make(map[string][]string),
					FilterableAttributes: []string{},
					SortableAttributes:   []string{},
					TypoTolerance:        &defaultTypoTolerance,
					Pagination:           &defaultPagination,
					ProximityPrecision:   meilisearch.ByWord,
					Faceting: &meilisearch.Faceting{
						MaxValuesPerFacet: 200,
						SortFacetValuesBy: map[string]meilisearch.SortFacetType{
							"*": meilisearch.SortFacetTypeAlpha,
						},
					},
					SeparatorTokens:    make([]string, 0),
					NonSeparatorTokens: make([]string, 0),
					Dictionary:         make([]string, 0),
					PrefixSearch:       stringPtr("indexingTime"),
					FacetSearch:        true,
					Embedders:          make(map[string]meilisearch.Embedder),
					Chat:               &defaultChat,
				},
				secondRequest: meilisearch.Settings{
					Faceting: &meilisearch.Faceting{
						MaxValuesPerFacet: 200,
					},
				},
				secondResponse: meilisearch.Settings{
					RankingRules: []string{
						"typo", "words",
					},
					DistinctAttribute:    (*string)(nil),
					SearchableAttributes: []string{"*"},
					DisplayedAttributes:  []string{"*"},
					StopWords:            []string{},
					Synonyms:             make(map[string][]string),
					FilterableAttributes: []string{},
					SortableAttributes:   []string{},
					TypoTolerance:        &defaultTypoTolerance,
					Pagination:           &defaultPagination,
					ProximityPrecision:   meilisearch.ByWord,
					Faceting: &meilisearch.Faceting{
						MaxValuesPerFacet: 200,
						SortFacetValuesBy: map[string]meilisearch.SortFacetType{
							"*": meilisearch.SortFacetTypeAlpha,
						},
					},
					SeparatorTokens:    make([]string, 0),
					NonSeparatorTokens: make([]string, 0),
					Dictionary:         make([]string, 0),
					PrefixSearch:       stringPtr("indexingTime"),
					FacetSearch:        true,
					Embedders:          make(map[string]meilisearch.Embedder),
					Chat:               &defaultChat,
				},
			},
			wantTask: &meilisearch.TaskInfo{
				TaskUID: 1,
			},
			wantResp: &meilisearch.Settings{
				RankingRules:         defaultRankingRules,
				DistinctAttribute:    (*string)(nil),
				SearchableAttributes: []string{"*"},
				DisplayedAttributes:  []string{"*"},
				StopWords:            []string{},
				Synonyms:             make(map[string][]string),
				FilterableAttributes: []string{},
				SortableAttributes:   []string{},
				TypoTolerance:        &defaultTypoTolerance,
				Pagination:           &defaultPagination,
				Faceting:             &defaultFaceting,
				ProximityPrecision:   meilisearch.ByWord,
				SeparatorTokens:      make([]string, 0),
				NonSeparatorTokens:   make([]string, 0),
				Dictionary:           make([]string, 0),
				PrefixSearch:         stringPtr("indexingTime"),
				FacetSearch:          true,
				Embedders:            make(map[string]meilisearch.Embedder),
				Chat:                 &defaultChat,
			},
		},
		{
			name: "TestIndexUpdateJustProximityPrecision",
			args: args{
				UID:    "indexUID",
				client: meili,
				firstRequest: meilisearch.Settings{
					RankingRules: []string{
						"typo", "words",
					},
					ProximityPrecision: meilisearch.ByAttribute,
				},
				firstResponse: meilisearch.Settings{
					RankingRules: []string{
						"typo", "words",
					},
					DistinctAttribute:    (*string)(nil),
					SearchableAttributes: []string{"*"},
					DisplayedAttributes:  []string{"*"},
					StopWords:            []string{},
					Synonyms:             make(map[string][]string),
					FilterableAttributes: []string{},
					SortableAttributes:   []string{},
					TypoTolerance:        &defaultTypoTolerance,
					Pagination:           &defaultPagination,
					ProximityPrecision:   meilisearch.ByAttribute,
					Faceting: &meilisearch.Faceting{
						MaxValuesPerFacet: 100,
						SortFacetValuesBy: map[string]meilisearch.SortFacetType{
							"*": meilisearch.SortFacetTypeAlpha,
						},
					},
					SeparatorTokens:    make([]string, 0),
					NonSeparatorTokens: make([]string, 0),
					Dictionary:         make([]string, 0),
					PrefixSearch:       stringPtr("indexingTime"),
					FacetSearch:        true,
					Embedders:          make(map[string]meilisearch.Embedder),
					Chat:               &defaultChat,
				},
				secondRequest: meilisearch.Settings{
					RankingRules: []string{
						"typo", "words",
					},
					ProximityPrecision: meilisearch.ByWord,
				},
				secondResponse: meilisearch.Settings{
					RankingRules: []string{
						"typo", "words",
					},
					DistinctAttribute:    (*string)(nil),
					SearchableAttributes: []string{"*"},
					DisplayedAttributes:  []string{"*"},
					StopWords:            []string{},
					Synonyms:             make(map[string][]string),
					FilterableAttributes: []string{},
					SortableAttributes:   []string{},
					TypoTolerance:        &defaultTypoTolerance,
					Pagination:           &defaultPagination,
					ProximityPrecision:   meilisearch.ByWord,
					Faceting: &meilisearch.Faceting{
						MaxValuesPerFacet: 100,
						SortFacetValuesBy: map[string]meilisearch.SortFacetType{
							"*": meilisearch.SortFacetTypeAlpha,
						},
					},
					SeparatorTokens:    make([]string, 0),
					NonSeparatorTokens: make([]string, 0),
					Dictionary:         make([]string, 0),
					PrefixSearch:       stringPtr("indexingTime"),
					FacetSearch:        true,
					Embedders:          make(map[string]meilisearch.Embedder),
					Chat:               &defaultChat,
				},
			},
			wantTask: &meilisearch.TaskInfo{
				TaskUID: 1,
			},
			wantResp: &meilisearch.Settings{
				RankingRules:         defaultRankingRules,
				DistinctAttribute:    (*string)(nil),
				SearchableAttributes: []string{"*"},
				DisplayedAttributes:  []string{"*"},
				StopWords:            []string{},
				Synonyms:             make(map[string][]string),
				FilterableAttributes: []string{},
				SortableAttributes:   []string{},
				TypoTolerance:        &defaultTypoTolerance,
				Pagination:           &defaultPagination,
				Faceting:             &defaultFaceting,
				ProximityPrecision:   meilisearch.ByWord,
				SeparatorTokens:      make([]string, 0),
				NonSeparatorTokens:   make([]string, 0),
				Dictionary:           make([]string, 0),
				PrefixSearch:         stringPtr("indexingTime"),
				FacetSearch:          true,
				Embedders:            make(map[string]meilisearch.Embedder),
				Chat:                 &defaultChat,
			},
		},
		{
			name: "TestIndexUpdateJustChat",
			args: args{
				UID:    "indexUID",
				client: meili,
				firstRequest: meilisearch.Settings{
					Chat: &meilisearch.Chat{
						Description:              "A comprehensive movie database containing titles, descriptions, genres, and release dates to help users find movies",
						DocumentTemplate:         "{% for field in fields %}{% if field.is_searchable and field.value != nil %}{{ field.name }}: {{ field.value }}\n{% endif %}{% endfor %}",
						DocumentTemplateMaxBytes: 400,
					},
				},
				firstResponse: meilisearch.Settings{
					RankingRules:         defaultRankingRules,
					DistinctAttribute:    (*string)(nil),
					SearchableAttributes: []string{"*"},
					DisplayedAttributes:  []string{"*"},
					StopWords:            []string{},
					Synonyms:             make(map[string][]string),
					FilterableAttributes: []string{},
					SortableAttributes:   []string{},
					TypoTolerance:        &defaultTypoTolerance,
					Pagination:           &defaultPagination,
					Faceting:             &defaultFaceting,
					ProximityPrecision:   meilisearch.ByWord,
					SeparatorTokens:      make([]string, 0),
					NonSeparatorTokens:   make([]string, 0),
					Dictionary:           make([]string, 0),
					PrefixSearch:         stringPtr("indexingTime"),
					FacetSearch:          true,
					Embedders:            make(map[string]meilisearch.Embedder),
					Chat: &meilisearch.Chat{
						Description:              "A comprehensive movie database containing titles, descriptions, genres, and release dates to help users find movies",
						DocumentTemplate:         "{% for field in fields %}{% if field.is_searchable and field.value != nil %}{{ field.name }}: {{ field.value }}\n{% endif %}{% endfor %}",
						DocumentTemplateMaxBytes: 400,
						SearchParameters:         defaultChat.SearchParameters,
					},
				},
				secondRequest: meilisearch.Settings{
					Chat: &meilisearch.Chat{
						Description:              "A comprehensive movie database containing titles, descriptions, genres, and release dates to help users find movies",
						DocumentTemplate:         "{% for field in fields %}{% if field.is_searchable and field.value != nil %}{{ field.name }}: {{ field.value }}\n{% endif %}{% endfor %}",
						DocumentTemplateMaxBytes: 400,
						SearchParameters: &meilisearch.SearchParameters{
							Limit: 20,
						},
					},
				},
				secondResponse: meilisearch.Settings{
					RankingRules:         defaultRankingRules,
					DistinctAttribute:    (*string)(nil),
					SearchableAttributes: []string{"*"},
					DisplayedAttributes:  []string{"*"},
					StopWords:            []string{},
					Synonyms:             make(map[string][]string),
					FilterableAttributes: []string{},
					SortableAttributes:   []string{},
					TypoTolerance:        &defaultTypoTolerance,
					Pagination:           &defaultPagination,
					Faceting:             &defaultFaceting,
					ProximityPrecision:   meilisearch.ByWord,
					SeparatorTokens:      make([]string, 0),
					NonSeparatorTokens:   make([]string, 0),
					Dictionary:           make([]string, 0),
					PrefixSearch:         stringPtr("indexingTime"),
					FacetSearch:          true,
					Embedders:            make(map[string]meilisearch.Embedder),
					Chat: &meilisearch.Chat{
						Description:              "A comprehensive movie database containing titles, descriptions, genres, and release dates to help users find movies",
						DocumentTemplate:         "{% for field in fields %}{% if field.is_searchable and field.value != nil %}{{ field.name }}: {{ field.value }}\n{% endif %}{% endfor %}",
						DocumentTemplateMaxBytes: 400,
						SearchParameters: &meilisearch.SearchParameters{
							Limit: 20,
						},
					},
				},
			},
			wantTask: &meilisearch.TaskInfo{
				TaskUID: 1,
			},
			wantResp: &meilisearch.Settings{
				RankingRules:         defaultRankingRules,
				DistinctAttribute:    (*string)(nil),
				SearchableAttributes: []string{"*"},
				DisplayedAttributes:  []string{"*"},
				StopWords:            []string{},
				Synonyms:             make(map[string][]string),
				FilterableAttributes: []string{},
				SortableAttributes:   []string{},
				TypoTolerance:        &defaultTypoTolerance,
				Pagination:           &defaultPagination,
				Faceting:             &defaultFaceting,
				ProximityPrecision:   meilisearch.ByWord,
				SeparatorTokens:      make([]string, 0),
				NonSeparatorTokens:   make([]string, 0),
				Dictionary:           make([]string, 0),
				PrefixSearch:         stringPtr("indexingTime"),
				FacetSearch:          true,
				Embedders:            make(map[string]meilisearch.Embedder),
				Chat:                 &defaultChat,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			setUpIndexForFaceting(tt.args.client)
			c := tt.args.client
			i := c.Index(tt.args.UID)
			t.Cleanup(cleanup(c))

			_, err := i.GetSettings()
			require.NoError(t, err)

			gotTask, err := i.UpdateSettings(&tt.args.firstRequest)
			require.NoError(t, err)
			require.GreaterOrEqual(t, gotTask.TaskUID, tt.wantTask.TaskUID)
			testWaitForTask(t, i, gotTask)

			gotResp, err := i.GetSettings()
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
	meili := setup(t, "")
	customMeili := setup(t, "", meilisearch.WithCustomClientWithTLS(&tls.Config{
		InsecureSkipVerify: true,
	}))

	type args struct {
		UID     string
		client  meilisearch.ServiceManager
		request []string
	}
	tests := []struct {
		name     string
		args     args
		wantTask *meilisearch.TaskInfo
	}{
		{
			name: "TestIndexBasicUpdateStopWords",
			args: args{
				UID:    "indexUID",
				client: meili,
				request: []string{
					"of", "the", "to",
				},
			},
			wantTask: &meilisearch.TaskInfo{
				TaskUID: 1,
			},
		},
		{
			name: "TestIndexUpdateStopWordsWithCustomClient",
			args: args{
				UID:    "indexUID",
				client: customMeili,
				request: []string{
					"of", "the", "to",
				},
			},
			wantTask: &meilisearch.TaskInfo{
				TaskUID: 1,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			setUpIndexForFaceting(tt.args.client)
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
	meili := setup(t, "")
	customMeili := setup(t, "", meilisearch.WithCustomClientWithTLS(&tls.Config{
		InsecureSkipVerify: true,
	}))

	type args struct {
		UID     string
		client  meilisearch.ServiceManager
		request map[string][]string
	}
	tests := []struct {
		name     string
		args     args
		wantTask *meilisearch.TaskInfo
	}{
		{
			name: "TestIndexBasicUpdateSynonyms",
			args: args{
				UID:    "indexUID",
				client: meili,
				request: map[string][]string{
					"wolverine": {"logan", "xmen"},
				},
			},
			wantTask: &meilisearch.TaskInfo{
				TaskUID: 1,
			},
		},
		{
			name: "TestIndexUpdateSynonymsWithCustomClient",
			args: args{
				UID:    "indexUID",
				client: customMeili,
				request: map[string][]string{
					"wolverine": {"logan", "xmen"},
				},
			},
			wantTask: &meilisearch.TaskInfo{
				TaskUID: 1,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			setUpIndexForFaceting(tt.args.client)
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
	meili := setup(t, "")
	customMeili := setup(t, "", meilisearch.WithCustomClientWithTLS(&tls.Config{
		InsecureSkipVerify: true,
	}))

	type args struct {
		UID     string
		client  meilisearch.ServiceManager
		request []string
	}
	tests := []struct {
		name     string
		args     args
		wantTask *meilisearch.TaskInfo
	}{
		{
			name: "TestIndexBasicUpdateSortableAttributes",
			args: args{
				UID:    "indexUID",
				client: meili,
				request: []string{
					"title",
				},
			},
			wantTask: &meilisearch.TaskInfo{
				TaskUID: 1,
			},
		},
		{
			name: "TestIndexUpdateSortableAttributesWithCustomClient",
			args: args{
				UID:    "indexUID",
				client: customMeili,
				request: []string{
					"title",
				},
			},
			wantTask: &meilisearch.TaskInfo{
				TaskUID: 1,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			setUpIndexForFaceting(tt.args.client)
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
	meili := setup(t, "")
	customMeili := setup(t, "", meilisearch.WithCustomClientWithTLS(&tls.Config{
		InsecureSkipVerify: true,
	}))

	type args struct {
		UID     string
		client  meilisearch.ServiceManager
		request meilisearch.TypoTolerance
	}
	tests := []struct {
		name     string
		args     args
		wantTask *meilisearch.TaskInfo
		wantResp *meilisearch.TypoTolerance
	}{
		{
			name: "TestIndexBasicUpdateTypoTolerance",
			args: args{
				UID:    "indexUID",
				client: meili,
				request: meilisearch.TypoTolerance{
					Enabled: true,
					MinWordSizeForTypos: meilisearch.MinWordSizeForTypos{
						OneTypo:  7,
						TwoTypos: 10,
					},
					DisableOnWords:      []string{},
					DisableOnAttributes: []string{},
					DisableOnNumbers:    true,
				},
			},
			wantTask: &meilisearch.TaskInfo{
				TaskUID: 1,
			},
			wantResp: &meilisearch.TypoTolerance{
				Enabled: true,
				MinWordSizeForTypos: meilisearch.MinWordSizeForTypos{
					OneTypo:  7,
					TwoTypos: 10,
				},
				DisableOnWords:      []string{},
				DisableOnAttributes: []string{},
				DisableOnNumbers:    true,
			},
		},
		{
			name: "TestIndexUpdateTypoToleranceWithCustomClient",
			args: args{
				UID:    "indexUID",
				client: customMeili,
				request: meilisearch.TypoTolerance{
					Enabled: true,
					MinWordSizeForTypos: meilisearch.MinWordSizeForTypos{
						OneTypo:  7,
						TwoTypos: 10,
					},
					DisableOnWords:      []string{},
					DisableOnAttributes: []string{},
					DisableOnNumbers:    false,
				},
			},
			wantTask: &meilisearch.TaskInfo{
				TaskUID: 1,
			},
			wantResp: &meilisearch.TypoTolerance{
				Enabled: true,
				MinWordSizeForTypos: meilisearch.MinWordSizeForTypos{
					OneTypo:  7,
					TwoTypos: 10,
				},
				DisableOnWords:      []string{},
				DisableOnAttributes: []string{},
				DisableOnNumbers:    false,
			},
		},
		{
			name: "TestIndexUpdateTypoToleranceWithDisableOnWords",
			args: args{
				UID:    "indexUID",
				client: meili,
				request: meilisearch.TypoTolerance{
					Enabled: true,
					MinWordSizeForTypos: meilisearch.MinWordSizeForTypos{
						OneTypo:  7,
						TwoTypos: 10,
					},
					DisableOnWords: []string{
						"and",
					},
					DisableOnAttributes: []string{},
					DisableOnNumbers:    true,
				},
			},
			wantTask: &meilisearch.TaskInfo{
				TaskUID: 1,
			},
			wantResp: &meilisearch.TypoTolerance{
				Enabled: true,
				MinWordSizeForTypos: meilisearch.MinWordSizeForTypos{
					OneTypo:  7,
					TwoTypos: 10,
				},
				DisableOnWords:      []string{"and"},
				DisableOnAttributes: []string{},
				DisableOnNumbers:    true,
			},
		},
		{
			name: "TestIndexUpdateTypoToleranceWithDisableOnAttributes",
			args: args{
				UID:    "indexUID",
				client: meili,
				request: meilisearch.TypoTolerance{
					Enabled: true,
					MinWordSizeForTypos: meilisearch.MinWordSizeForTypos{
						OneTypo:  7,
						TwoTypos: 10,
					},
					DisableOnWords: []string{},
					DisableOnAttributes: []string{
						"year",
					},
					DisableOnNumbers: true,
				},
			},
			wantTask: &meilisearch.TaskInfo{
				TaskUID: 1,
			},
			wantResp: &meilisearch.TypoTolerance{
				Enabled: true,
				MinWordSizeForTypos: meilisearch.MinWordSizeForTypos{
					OneTypo:  7,
					TwoTypos: 10,
				},
				DisableOnWords:      []string{},
				DisableOnAttributes: []string{"year"},
				DisableOnNumbers:    true,
			},
		},
		{
			name: "TestIndexDisableTypoTolerance",
			args: args{
				UID:    "indexUID",
				client: meili,
				request: meilisearch.TypoTolerance{
					Enabled: false,
					MinWordSizeForTypos: meilisearch.MinWordSizeForTypos{
						OneTypo:  5,
						TwoTypos: 9,
					},
					DisableOnWords:      []string{},
					DisableOnAttributes: []string{},
					DisableOnNumbers:    false,
				},
			},
			wantTask: &meilisearch.TaskInfo{
				TaskUID: 1,
			},
			wantResp: &meilisearch.TypoTolerance{
				Enabled: false,
				MinWordSizeForTypos: meilisearch.MinWordSizeForTypos{
					OneTypo:  5,
					TwoTypos: 9,
				},
				DisableOnWords:      []string{},
				DisableOnAttributes: []string{},
				DisableOnNumbers:    false,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			setUpIndexForFaceting(tt.args.client)
			c := tt.args.client
			i := c.Index(tt.args.UID)
			t.Cleanup(cleanup(c))

			gotTask, err := i.UpdateTypoTolerance(&tt.args.request)
			require.NoError(t, err)
			require.GreaterOrEqual(t, gotTask.TaskUID, tt.wantTask.TaskUID)
			testWaitForTask(t, i, gotTask)

			gotResp, err := i.GetTypoTolerance()
			require.NoError(t, err)
			require.Equal(t, tt.wantResp, gotResp)
		})
	}
}

func TestIndex_UpdatePagination(t *testing.T) {
	meili := setup(t, "")
	customMeili := setup(t, "", meilisearch.WithCustomClientWithTLS(&tls.Config{
		InsecureSkipVerify: true,
	}))

	type args struct {
		UID     string
		client  meilisearch.ServiceManager
		request meilisearch.Pagination
	}
	tests := []struct {
		name     string
		args     args
		wantTask *meilisearch.TaskInfo
		wantResp *meilisearch.Pagination
	}{
		{
			name: "TestIndexBasicUpdatePagination",
			args: args{
				UID:    "indexUID",
				client: meili,
				request: meilisearch.Pagination{
					MaxTotalHits: 1200,
				},
			},
			wantTask: &meilisearch.TaskInfo{
				TaskUID: 1,
			},
			wantResp: &defaultPagination,
		},
		{
			name: "TestIndexUpdatePaginationWithCustomClient",
			args: args{
				UID:    "indexUID",
				client: customMeili,
				request: meilisearch.Pagination{
					MaxTotalHits: 1200,
				},
			},
			wantTask: &meilisearch.TaskInfo{
				TaskUID: 1,
			},
			wantResp: &defaultPagination,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			setUpIndexForFaceting(tt.args.client)
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
	meili := setup(t, "")
	customMeili := setup(t, "", meilisearch.WithCustomClientWithTLS(&tls.Config{
		InsecureSkipVerify: true,
	}))

	type args struct {
		UID     string
		client  meilisearch.ServiceManager
		request meilisearch.Faceting
	}
	tests := []struct {
		name     string
		args     args
		wantTask *meilisearch.TaskInfo
		wantResp *meilisearch.Faceting
	}{
		{
			name: "TestIndexBasicUpdateFaceting",
			args: args{
				UID:    "indexUID",
				client: meili,
				request: meilisearch.Faceting{
					MaxValuesPerFacet: 200,
					SortFacetValuesBy: map[string]meilisearch.SortFacetType{
						"*": meilisearch.SortFacetTypeAlpha,
					},
				},
			},
			wantTask: &meilisearch.TaskInfo{
				TaskUID: 1,
			},
			wantResp: &defaultFaceting,
		},
		{
			name: "TestIndexUpdateFacetingWithCustomClient",
			args: args{
				UID:    "indexUID",
				client: customMeili,
				request: meilisearch.Faceting{
					MaxValuesPerFacet: 200,
					SortFacetValuesBy: map[string]meilisearch.SortFacetType{
						"*": meilisearch.SortFacetTypeAlpha,
					},
				},
			},
			wantTask: &meilisearch.TaskInfo{
				TaskUID: 1,
			},
			wantResp: &defaultFaceting,
		},
		{
			name: "TestIndexGetStartedFaceting",
			args: args{
				UID:    "indexUID",
				client: meili,
				request: meilisearch.Faceting{
					MaxValuesPerFacet: 2,
					SortFacetValuesBy: map[string]meilisearch.SortFacetType{
						"*": meilisearch.SortFacetTypeCount,
					},
				},
			},
			wantTask: &meilisearch.TaskInfo{
				TaskUID: 1,
			},
			wantResp: &meilisearch.Faceting{
				MaxValuesPerFacet: 2,
				SortFacetValuesBy: map[string]meilisearch.SortFacetType{
					"*": meilisearch.SortFacetTypeCount,
				},
			},
		},
		{
			name: "TestIndexSortFacetValuesByCountFaceting",
			args: args{
				UID:    "indexUID",
				client: meili,
				request: meilisearch.Faceting{
					SortFacetValuesBy: map[string]meilisearch.SortFacetType{
						"*":        meilisearch.SortFacetTypeAlpha,
						"indexUID": meilisearch.SortFacetTypeCount,
					},
				},
			},
			wantTask: &meilisearch.TaskInfo{
				TaskUID: 1,
			},
			wantResp: &meilisearch.Faceting{
				SortFacetValuesBy: map[string]meilisearch.SortFacetType{
					"*":        meilisearch.SortFacetTypeAlpha,
					"indexUID": meilisearch.SortFacetTypeCount,
				},
			},
		},
		{
			name: "TestIndexSortFacetValuesAllIndexFaceting",
			args: args{
				UID:    "indexUID",
				client: meili,
				request: meilisearch.Faceting{
					SortFacetValuesBy: map[string]meilisearch.SortFacetType{
						"*": meilisearch.SortFacetTypeCount,
					},
				},
			},
			wantTask: &meilisearch.TaskInfo{
				TaskUID: 1,
			},
			wantResp: &meilisearch.Faceting{
				SortFacetValuesBy: map[string]meilisearch.SortFacetType{
					"*": meilisearch.SortFacetTypeCount,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			setUpIndexForFaceting(tt.args.client)
			c := tt.args.client
			i := c.Index(tt.args.UID)
			t.Cleanup(cleanup(c))

			gotTask, err := i.UpdateFaceting(&tt.args.request)
			require.NoError(t, err)
			require.GreaterOrEqual(t, gotTask.TaskUID, tt.wantTask.TaskUID)
			testWaitForTask(t, i, gotTask)

			gotResp, err := i.GetFaceting()
			require.NoError(t, err)
			require.Equal(t, &tt.args.request, gotResp)
		})
	}
}

func TestIndex_UpdateSettingsEmbedders(t *testing.T) {
	meili := setup(t, "")

	type args struct {
		UID      string
		client   meilisearch.ServiceManager
		request  meilisearch.Settings
		newIndex bool
	}
	tests := []struct {
		name          string
		args          args
		wantTask      *meilisearch.TaskInfo
		wantEmbedders map[string]meilisearch.Embedder
		wantErr       string
	}{
		{
			name: "TestIndexUpdateSettingsEmbeddersErr",
			args: args{
				UID:    "indexUID",
				client: meili,
				request: meilisearch.Settings{
					Embedders: map[string]meilisearch.Embedder{
						"default": {
							Source:           meilisearch.OpenaiEmbedderSource,
							DocumentTemplate: "{{doc.foobar}}",
						},
					},
				},
			},
			wantTask: &meilisearch.TaskInfo{
				TaskUID: 1,
			},
			wantErr: "foobar",
		},
		{
			name: "TestIndexUpdateSettingsEmbeddersErr",
			args: args{
				UID:    "indexUID",
				client: meili,
				request: meilisearch.Settings{
					Embedders: map[string]meilisearch.Embedder{
						"default": {
							Source:           meilisearch.OpenaiEmbedderSource,
							APIKey:           "xxx",
							Model:            "text-embedding-3-small",
							DocumentTemplate: "A movie titled '{{doc.title}}'",
						},
					},
				},
			},
			wantTask: &meilisearch.TaskInfo{
				TaskUID: 1,
			},
			wantErr: "Incorrect API key",
		},
		{
			name: "TestIndexUpdateSettingsEmbeddersHuggingFace",
			args: args{
				newIndex: true,
				UID:      "newIndexUID",
				client:   meili,
				request: meilisearch.Settings{
					Embedders: map[string]meilisearch.Embedder{
						"default": {
							Source:           meilisearch.HuggingFaceEmbedderSource,
							Model:            "sentence-transformers/paraphrase-multilingual-MiniLM-L12-v2",
							DocumentTemplate: "A movie titled '{{doc.title}}' whose description starts with {{doc.overview|truncatewords: 20}}",
							Distribution: &meilisearch.Distribution{
								Mean:  0.7,
								Sigma: 0.3,
							},
							Pooling:                  meilisearch.UseModelEmbedderPooling,
							DocumentTemplateMaxBytes: 500,
							BinaryQuantized:          false,
						},
					},
				},
			},
			wantTask: &meilisearch.TaskInfo{
				TaskUID: 1,
			},
		},
		{
			name: "TestIndexUpdateSettingsEmbeddersUserProvided",
			args: args{
				newIndex: true,
				UID:      "newIndexUID",
				client:   meili,
				request: meilisearch.Settings{
					Embedders: map[string]meilisearch.Embedder{
						"default": {
							Source:     meilisearch.UserProvidedEmbedderSource,
							Dimensions: 3,
						},
					},
				},
			},
			wantTask: &meilisearch.TaskInfo{
				TaskUID: 1,
			},
		},
		{
			name: "TestIndexUpdateSettingsEmbeddersWithRestSource",
			args: args{
				newIndex: true,
				UID:      "newIndexUID",
				client:   meili,
				request: meilisearch.Settings{
					Embedders: map[string]meilisearch.Embedder{
						"default": {
							Source:           meilisearch.RestEmbedderSource,
							URL:              "https://api.openai.com/v1/embeddings",
							APIKey:           "<your-openai-api-key>",
							Dimensions:       1536,
							DocumentTemplate: "A movie titled '{{doc.title}}' whose description starts with {{doc.overview|truncatewords: 20}}",
							Distribution: &meilisearch.Distribution{
								Mean:  0.7,
								Sigma: 0.3,
							},
							Request: map[string]interface{}{
								"model": "text-embedding-3-small",
								"input": []string{"{{text}}", "{{..}}"},
							},
							Response: map[string]interface{}{
								"data": []interface{}{
									map[string]interface{}{
										"embedding": "{{embedding}}",
									},
									"{{..}}",
								},
							},
							Headers: map[string]string{
								"Custom-Header": "CustomValue",
							},
						},
					},
				},
			},
			wantTask: &meilisearch.TaskInfo{
				TaskUID: 1,
			},
		},
		{
			name: "TestIndexUpdateSettingsEmbeddersOllama",
			args: args{
				newIndex: true,
				UID:      "newIndexUID",
				client:   meili,
				request: meilisearch.Settings{
					Embedders: map[string]meilisearch.Embedder{
						"default": {
							Source:           meilisearch.OllamaEmbedderSource,
							URL:              "http://localhost:11434/api/embeddings",
							APIKey:           "<your-ollama-api-key>",
							Model:            "nomic-embed-text",
							DocumentTemplate: "blabla",
							Distribution: &meilisearch.Distribution{
								Mean:  0.7,
								Sigma: 0.3,
							},
							Dimensions:               512,
							DocumentTemplateMaxBytes: 500,
							BinaryQuantized:          false,
						},
					},
				},
			},
			wantTask: &meilisearch.TaskInfo{
				TaskUID: 1,
			},
		},
		{
			name: "TestIndexUpdateSettingsEmbeddersComposite",
			args: args{
				newIndex: true,
				UID:      "newIndexUID",
				client:   meili,
				request: meilisearch.Settings{
					Embedders: map[string]meilisearch.Embedder{
						"default": {
							Source: meilisearch.CompositeEmbedderSource,
							SearchEmbedder: &meilisearch.Embedder{
								Source:  meilisearch.HuggingFaceEmbedderSource,
								Model:   "sentence-transformers/paraphrase-multilingual-MiniLM-L12-v2",
								Pooling: meilisearch.UseModelEmbedderPooling,
							},
							IndexingEmbedder: &meilisearch.Embedder{
								Source:                   meilisearch.HuggingFaceEmbedderSource,
								Model:                    "sentence-transformers/paraphrase-multilingual-MiniLM-L12-v2",
								DocumentTemplate:         "{{doc.title}}",
								Pooling:                  meilisearch.UseModelEmbedderPooling,
								DocumentTemplateMaxBytes: 500,
							},
						},
					},
				},
			},
			wantTask: &meilisearch.TaskInfo{
				TaskUID: 1,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := tt.args.client
			i, err := setUpIndexWithVector(c, tt.args.UID)
			require.NoError(t, err)
			t.Cleanup(cleanup(c))

			feat := c.ExperimentalFeatures().SetCompositeEmbedders(true)
			resp, err := feat.Update()
			require.NoError(t, err)
			require.True(t, resp.CompositeEmbedders)

			gotTask, err := i.UpdateSettings(&tt.args.request)
			require.NoError(t, err)
			require.GreaterOrEqual(t, gotTask.TaskUID, tt.wantTask.TaskUID)
			testWaitForTask(t, i, gotTask)

			gotResp, err := i.GetEmbedders()
			require.NoError(t, err)
			require.NotNil(t, gotResp)
			assert.Equal(t, gotResp["default"].Source, tt.args.request.Embedders["default"].Source)
		})
	}
}

func TestIndex_GetEmbedders(t *testing.T) {
	c := setup(t, "")
	t.Cleanup(cleanup(c))

	indexID := "newIndexUID"
	i := c.Index(indexID)
	task, err := c.CreateIndex(&meilisearch.IndexConfig{Uid: indexID})
	require.NoError(t, err)
	testWaitForTask(t, i, task)

	expected := map[string]meilisearch.Embedder{
		"default": {
			Source:     "userProvided",
			Dimensions: 3,
		},
	}
	task, err = i.UpdateSettings(&meilisearch.Settings{
		Embedders: expected,
	})
	require.NoError(t, err)
	testWaitForTask(t, i, task)

	got, err := i.GetEmbedders()
	require.NoError(t, err)
	require.Equal(t, expected, got)
}

func TestIndex_UpdateEmbedders(t *testing.T) {
	c := setup(t, "")
	t.Cleanup(cleanup(c))

	indexID := "newIndexUID"
	i := c.Index(indexID)
	taskInfo, err := c.CreateIndex(&meilisearch.IndexConfig{Uid: indexID})
	require.NoError(t, err)
	testWaitForTask(t, i, taskInfo)

	embedders := map[string]meilisearch.Embedder{
		"someEmbbeder": {
			Source:     "userProvided",
			Dimensions: 3,
		},
	}
	taskInfo, err = i.UpdateSettings(&meilisearch.Settings{
		Embedders: embedders,
	})
	require.NoError(t, err)
	testWaitForTask(t, i, taskInfo)

	updated := map[string]meilisearch.Embedder{
		"someEmbbeder": {
			Source:     "userProvided",
			Dimensions: 5,
		},
	}

	taskInfo, err = i.UpdateEmbedders(updated)
	require.NoError(t, err)
	task, err := i.WaitForTask(taskInfo.TaskUID, 0)
	require.NoError(t, err)
	require.Equal(t, meilisearch.TaskStatusSucceeded, task.Status)

	got, err := i.GetEmbedders()
	require.NoError(t, err)
	require.Equal(t, updated, got)
}

func TestIndex_ResetEmbedders(t *testing.T) {
	c := setup(t, "")
	t.Cleanup(cleanup(c))

	indexID := "newIndexUID"
	i := c.Index(indexID)
	taskInfo, err := c.CreateIndex(&meilisearch.IndexConfig{Uid: indexID})
	require.NoError(t, err)
	testWaitForTask(t, i, taskInfo)

	taskInfo, err = i.UpdateSettings(&meilisearch.Settings{
		Embedders: map[string]meilisearch.Embedder{
			"default": {
				Source:     "userProvided",
				Dimensions: 3,
			},
		},
	})
	require.NoError(t, err)
	testWaitForTask(t, i, taskInfo)

	taskInfo, err = i.ResetEmbedders()
	require.NoError(t, err)
	task, err := i.WaitForTask(taskInfo.TaskUID, 0)
	require.NoError(t, err)
	require.Equal(t, meilisearch.TaskStatusSucceeded, task.Status)

	got, err := i.GetEmbedders()
	require.NoError(t, err)
	require.Empty(t, got)
}

func Test_Dictionary(t *testing.T) {
	c := setup(t, "")
	t.Cleanup(cleanup(c))

	indexID := "newIndexUID"
	i := c.Index(indexID)

	words := []string{"J. R. R.", "W. E. B."}

	task, err := i.UpdateDictionary(words)
	require.NoError(t, err)
	testWaitForTask(t, i, task)

	got, err := i.GetDictionary()
	require.NoError(t, err)
	require.Equal(t, words, got)

	task, err = i.ResetDictionary()
	require.NoError(t, err)
	testWaitForTask(t, i, task)

	got, err = i.GetDictionary()
	require.NoError(t, err)
	require.Equal(t, got, []string{})
}

func Test_SearchCutoffMs(t *testing.T) {
	c := setup(t, "")
	t.Cleanup(cleanup(c))

	indexID := "newIndexUID"
	i := c.Index(indexID)
	taskInfo, err := c.CreateIndex(&meilisearch.IndexConfig{Uid: indexID})
	require.NoError(t, err)
	testWaitForTask(t, i, taskInfo)

	n := int64(250)

	task, err := i.UpdateSearchCutoffMs(n)
	require.NoError(t, err)
	testWaitForTask(t, i, task)

	got, err := i.GetSearchCutoffMs()
	require.NoError(t, err)
	require.Equal(t, n, got)

	task, err = i.ResetSearchCutoffMs()
	require.NoError(t, err)
	testWaitForTask(t, i, task)

	got, err = i.GetSearchCutoffMs()
	require.NoError(t, err)
	require.Equal(t, int64(0), got)
}

func Test_SeparatorTokens(t *testing.T) {
	c := setup(t, "")

	indexID := "newIndexUID"
	i := c.Index(indexID)

	tokens := []string{"|", "&hellip;"}

	task, err := i.UpdateSeparatorTokens(tokens)
	require.NoError(t, err)
	testWaitForTask(t, i, task)

	got, err := i.GetSeparatorTokens()
	require.NoError(t, err)
	require.ElementsMatchf(t, tokens, got, "tokens is not match with got")

	task, err = i.ResetSeparatorTokens()
	require.NoError(t, err)
	testWaitForTask(t, i, task)

	got, err = i.GetSeparatorTokens()
	require.NoError(t, err)
	require.Equal(t, got, []string{})
}

func Test_NonSeparatorTokens(t *testing.T) {
	c := setup(t, "")

	indexID := "newIndexUID"
	i := c.Index(indexID)

	tokens := []string{"@", "#"}

	task, err := i.UpdateNonSeparatorTokens(tokens)
	require.NoError(t, err)
	testWaitForTask(t, i, task)

	got, err := i.GetNonSeparatorTokens()
	require.NoError(t, err)
	require.ElementsMatchf(t, tokens, got, "tokens is not match with got")

	task, err = i.ResetNonSeparatorTokens()
	require.NoError(t, err)
	testWaitForTask(t, i, task)

	got, err = i.GetNonSeparatorTokens()
	require.NoError(t, err)
	require.Equal(t, got, []string{})
}

func Test_ProximityPrecision(t *testing.T) {
	c := setup(t, "")
	t.Cleanup(cleanup(c))

	indexID := "newIndexUID"
	i := c.Index(indexID)

	got, err := i.GetProximityPrecision()
	require.NoError(t, err)
	require.Equal(t, meilisearch.ByWord, got)

	task, err := i.UpdateProximityPrecision(meilisearch.ByAttribute)
	require.NoError(t, err)
	testWaitForTask(t, i, task)

	got, err = i.GetProximityPrecision()
	require.NoError(t, err)
	require.Equal(t, meilisearch.ByAttribute, got)

	task, err = i.ResetProximityPrecision()
	require.NoError(t, err)
	testWaitForTask(t, i, task)

	got, err = i.GetProximityPrecision()
	require.NoError(t, err)
	require.Equal(t, meilisearch.ByWord, got)
}

func Test_LocalizedAttributes(t *testing.T) {
	c := setup(t, "")
	t.Cleanup(cleanup(c))

	indexID := "newIndexUID"
	i := c.Index(indexID)
	taskInfo, err := c.CreateIndex(&meilisearch.IndexConfig{Uid: indexID})
	require.NoError(t, err)
	testWaitForTask(t, i, taskInfo)

	defer t.Cleanup(cleanup(c))

	t.Run("Test valid locate", func(t *testing.T) {
		got, err := i.GetLocalizedAttributes()
		require.NoError(t, err)
		require.Len(t, got, 0)

		localized := &meilisearch.LocalizedAttributes{
			Locales:           []string{"jpn", "eng"},
			AttributePatterns: []string{"*_ja"},
		}

		task, err := i.UpdateLocalizedAttributes([]*meilisearch.LocalizedAttributes{localized})
		require.NoError(t, err)
		testWaitForTask(t, i, task)

		got, err = i.GetLocalizedAttributes()
		require.NoError(t, err)
		require.NotNil(t, got)

		require.Equal(t, localized.Locales, got[0].Locales)
		require.Equal(t, localized.AttributePatterns, got[0].AttributePatterns)

		task, err = i.ResetLocalizedAttributes()
		require.NoError(t, err)
		testWaitForTask(t, i, task)

		got, err = i.GetLocalizedAttributes()
		require.NoError(t, err)
		require.Len(t, got, 0)
	})

	t.Run("Test invalid locate", func(t *testing.T) {
		invalidLocalized := &meilisearch.LocalizedAttributes{
			Locales:           []string{"foo"},
			AttributePatterns: []string{"*_ja"},
		}

		_, err := i.UpdateLocalizedAttributes([]*meilisearch.LocalizedAttributes{invalidLocalized})
		require.Error(t, err)
	})
}

func TestIndex_GetPrefixSearch(t *testing.T) {
	meili := setup(t, "")
	customMeili := setup(t, "", meilisearch.WithCustomClientWithTLS(&tls.Config{
		InsecureSkipVerify: true,
	}))

	type args struct {
		UID    string
		client meilisearch.ServiceManager
	}
	tests := []struct {
		name     string
		args     args
		wantResp *string
	}{
		{
			name: "TestIndexBasicGetPrefixSearch",
			args: args{
				UID:    "indexUID",
				client: meili,
			},
			wantResp: stringPtr("indexingTime"),
		},
		{
			name: "TestIndexGetPrefixSearchWithCustomClient",
			args: args{
				UID:    "indexUID",
				client: customMeili,
			},
			wantResp: stringPtr("indexingTime"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			setUpIndexForFaceting(tt.args.client)
			c := tt.args.client
			i := c.Index(tt.args.UID)
			t.Cleanup(cleanup(c))

			gotResp, err := i.GetPrefixSearch()
			require.NoError(t, err)
			require.Equal(t, tt.wantResp, gotResp)
		})
	}
}

func TestIndex_UpdatePrefixSearch(t *testing.T) {
	meili := setup(t, "")
	customMeili := setup(t, "", meilisearch.WithCustomClientWithTLS(&tls.Config{
		InsecureSkipVerify: true,
	}))

	type args struct {
		UID     string
		client  meilisearch.ServiceManager
		request string
	}
	tests := []struct {
		name     string
		args     args
		wantTask *meilisearch.TaskInfo
		wantResp *string
	}{
		{
			name: "TestIndexBasicUpdatePrefixSearch",
			args: args{
				UID:     "indexUID",
				client:  meili,
				request: "disabled",
			},
			wantTask: &meilisearch.TaskInfo{TaskUID: 1},
			wantResp: stringPtr("disabled"),
		},
		{
			name: "TestIndexUpdatePrefixSearchWithCustomClient",
			args: args{
				UID:     "indexUID",
				client:  customMeili,
				request: "disabled",
			},
			wantTask: &meilisearch.TaskInfo{TaskUID: 1},
			wantResp: stringPtr("disabled"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			setUpIndexForFaceting(tt.args.client)
			c := tt.args.client
			i := c.Index(tt.args.UID)
			t.Cleanup(cleanup(c))

			gotTask, err := i.UpdatePrefixSearch(tt.args.request)
			require.NoError(t, err)
			require.GreaterOrEqual(t, gotTask.TaskUID, tt.wantTask.TaskUID)
			testWaitForTask(t, i, gotTask)

			gotResp, err := i.GetPrefixSearch()
			require.NoError(t, err)
			require.Equal(t, tt.wantResp, gotResp)
		})
	}
}

func TestIndex_ResetPrefixSearch(t *testing.T) {
	meili := setup(t, "")
	customMeili := setup(t, "", meilisearch.WithCustomClientWithTLS(&tls.Config{
		InsecureSkipVerify: true,
	}))

	type args struct {
		UID    string
		client meilisearch.ServiceManager
	}
	tests := []struct {
		name     string
		args     args
		wantTask *meilisearch.TaskInfo
		wantResp *string
	}{
		{
			name: "TestIndexBasicResetPrefixSearch",
			args: args{
				UID:    "indexUID",
				client: meili,
			},
			wantTask: &meilisearch.TaskInfo{TaskUID: 1},
			wantResp: stringPtr("indexingTime"),
		},
		{
			name: "TestIndexResetPrefixSearchWithCustomClient",
			args: args{
				UID:    "indexUID",
				client: customMeili,
			},
			wantTask: &meilisearch.TaskInfo{TaskUID: 1},
			wantResp: stringPtr("indexingTime"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			setUpIndexForFaceting(tt.args.client)
			c := tt.args.client
			i := c.Index(tt.args.UID)
			t.Cleanup(cleanup(c))

			// First update to a non-default value
			_, err := i.UpdatePrefixSearch("disabled")
			require.NoError(t, err)

			gotTask, err := i.ResetPrefixSearch()
			require.NoError(t, err)
			require.GreaterOrEqual(t, gotTask.TaskUID, tt.wantTask.TaskUID)
			testWaitForTask(t, i, gotTask)

			gotResp, err := i.GetPrefixSearch()
			require.NoError(t, err)
			require.Equal(t, tt.wantResp, gotResp)
		})
	}
}

func TestIndex_GetFacetSearch(t *testing.T) {
	meili := setup(t, "")
	customMeili := setup(t, "", meilisearch.WithCustomClientWithTLS(&tls.Config{
		InsecureSkipVerify: true,
	}))

	type args struct {
		UID    string
		client meilisearch.ServiceManager
	}
	tests := []struct {
		name     string
		args     args
		wantResp bool
	}{
		{
			name: "TestIndexBasicGetFacetSearch",
			args: args{
				UID:    "indexUID",
				client: meili,
			},
			wantResp: boolPtr(true),
		},
		{
			name: "TestIndexGetFacetSearchWithCustomClient",
			args: args{
				UID:    "indexUID",
				client: customMeili,
			},
			wantResp: boolPtr(true),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			setUpIndexForFaceting(tt.args.client)
			c := tt.args.client
			i := c.Index(tt.args.UID)
			t.Cleanup(cleanup(c))

			gotResp, err := i.GetFacetSearch()
			require.NoError(t, err)
			require.Equal(t, tt.wantResp, gotResp)
		})
	}
}

func TestIndex_UpdateFacetSearch(t *testing.T) {
	meili := setup(t, "")
	customMeili := setup(t, "", meilisearch.WithCustomClientWithTLS(&tls.Config{
		InsecureSkipVerify: true,
	}))

	type args struct {
		UID     string
		client  meilisearch.ServiceManager
		request bool
	}
	tests := []struct {
		name     string
		args     args
		wantTask *meilisearch.TaskInfo
		wantResp bool
	}{
		{
			name: "TestIndexBasicUpdateFacetSearch",
			args: args{
				UID:     "indexUID",
				client:  meili,
				request: false,
			},
			wantTask: &meilisearch.TaskInfo{TaskUID: 1},
			wantResp: boolPtr(false),
		},
		{
			name: "TestIndexUpdateFacetSearchWithCustomClient",
			args: args{
				UID:     "indexUID",
				client:  customMeili,
				request: false,
			},
			wantTask: &meilisearch.TaskInfo{TaskUID: 1},
			wantResp: boolPtr(false),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			setUpIndexForFaceting(tt.args.client)
			c := tt.args.client
			i := c.Index(tt.args.UID)
			t.Cleanup(cleanup(c))

			gotTask, err := i.UpdateFacetSearch(tt.args.request)
			require.NoError(t, err)
			require.GreaterOrEqual(t, gotTask.TaskUID, tt.wantTask.TaskUID)
			testWaitForTask(t, i, gotTask)

			gotResp, err := i.GetFacetSearch()
			require.NoError(t, err)
			require.Equal(t, tt.wantResp, gotResp)
		})
	}
}

func TestIndex_ResetFacetSearch(t *testing.T) {
	meili := setup(t, "")
	customMeili := setup(t, "", meilisearch.WithCustomClientWithTLS(&tls.Config{
		InsecureSkipVerify: true,
	}))

	type args struct {
		UID    string
		client meilisearch.ServiceManager
	}
	tests := []struct {
		name     string
		args     args
		wantTask *meilisearch.TaskInfo
		wantResp bool
	}{
		{
			name: "TestIndexBasicResetFacetSearch",
			args: args{
				UID:    "indexUID",
				client: meili,
			},
			wantTask: &meilisearch.TaskInfo{TaskUID: 1},
			wantResp: boolPtr(true),
		},
		{
			name: "TestIndexResetFacetSearchWithCustomClient",
			args: args{
				UID:    "indexUID",
				client: customMeili,
			},
			wantTask: &meilisearch.TaskInfo{TaskUID: 1},
			wantResp: boolPtr(true),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			setUpIndexForFaceting(tt.args.client)
			c := tt.args.client
			i := c.Index(tt.args.UID)
			t.Cleanup(cleanup(c))

			// First update to a non-default value
			_, err := i.UpdateFacetSearch(false)
			require.NoError(t, err)

			gotTask, err := i.ResetFacetSearch()
			require.NoError(t, err)
			require.GreaterOrEqual(t, gotTask.TaskUID, tt.wantTask.TaskUID)
			testWaitForTask(t, i, gotTask)

			gotResp, err := i.GetFacetSearch()
			require.NoError(t, err)
			require.Equal(t, tt.wantResp, gotResp)
		})
	}
}

// Helper functions
func stringPtr(s string) *string {
	return &s
}

func boolPtr(b bool) bool {
	return b
}
