package meilisearch

import (
	"crypto/tls"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestIndex_GetFilterableAttributes(t *testing.T) {
	meili := setup(t, "")
	customMeili := setup(t, "", WithCustomClientWithTLS(&tls.Config{
		InsecureSkipVerify: true,
	}))

	type args struct {
		UID    string
		client ServiceManager
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
	customMeili := setup(t, "", WithCustomClientWithTLS(&tls.Config{
		InsecureSkipVerify: true,
	}))

	type args struct {
		UID    string
		client ServiceManager
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
		client ServiceManager
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
		client ServiceManager
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
		client ServiceManager
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
	customMeili := setup(t, "", WithCustomClientWithTLS(&tls.Config{
		InsecureSkipVerify: true,
	}))

	type args struct {
		UID    string
		client ServiceManager
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
				client: meili,
			},
			wantResp: &Settings{
				RankingRules:         defaultRankingRules,
				DistinctAttribute:    (*string)(nil),
				SearchableAttributes: []string{"*"},
				SearchCutoffMs:       0,
				ProximityPrecision:   ByWord,
				DisplayedAttributes:  []string{"*"},
				StopWords:            []string{},
				Synonyms:             map[string][]string(nil),
				FilterableAttributes: []string{},
				SortableAttributes:   []string{},
				TypoTolerance:        &defaultTypoTolerance,
				Pagination:           &defaultPagination,
				Faceting:             &defaultFaceting,
				SeparatorTokens:      make([]string, 0),
				NonSeparatorTokens:   make([]string, 0),
				Dictionary:           make([]string, 0),
				LocalizedAttributes:  nil,
			},
		},
		{
			name: "TestIndexGetSettingsWithCustomClient",
			args: args{
				UID:    "indexUID",
				client: customMeili,
			},
			wantResp: &Settings{
				RankingRules:         defaultRankingRules,
				DistinctAttribute:    (*string)(nil),
				SearchableAttributes: []string{"*"},
				SearchCutoffMs:       0,
				ProximityPrecision:   ByWord,
				DisplayedAttributes:  []string{"*"},
				StopWords:            []string{},
				Synonyms:             map[string][]string(nil),
				FilterableAttributes: []string{},
				SortableAttributes:   []string{},
				TypoTolerance:        &defaultTypoTolerance,
				Pagination:           &defaultPagination,
				Faceting:             &defaultFaceting,
				SeparatorTokens:      make([]string, 0),
				NonSeparatorTokens:   make([]string, 0),
				Dictionary:           make([]string, 0),
				LocalizedAttributes:  nil,
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
	customMeili := setup(t, "", WithCustomClientWithTLS(&tls.Config{
		InsecureSkipVerify: true,
	}))

	type args struct {
		UID    string
		client ServiceManager
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
	customMeili := setup(t, "", WithCustomClientWithTLS(&tls.Config{
		InsecureSkipVerify: true,
	}))

	type args struct {
		UID    string
		client ServiceManager
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
	customMeili := setup(t, "", WithCustomClientWithTLS(&tls.Config{
		InsecureSkipVerify: true,
	}))

	type args struct {
		UID    string
		client ServiceManager
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
	customMeili := setup(t, "", WithCustomClientWithTLS(&tls.Config{
		InsecureSkipVerify: true,
	}))

	type args struct {
		UID    string
		client ServiceManager
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
	customMeili := setup(t, "", WithCustomClientWithTLS(&tls.Config{
		InsecureSkipVerify: true,
	}))

	type args struct {
		UID    string
		client ServiceManager
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
	customMeili := setup(t, "", WithCustomClientWithTLS(&tls.Config{
		InsecureSkipVerify: true,
	}))

	type args struct {
		UID    string
		client ServiceManager
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
	customMeili := setup(t, "", WithCustomClientWithTLS(&tls.Config{
		InsecureSkipVerify: true,
	}))

	type args struct {
		UID    string
		client ServiceManager
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
				client: meili,
			},
			wantTask: &TaskInfo{
				TaskUID: 1,
			},
		},
		{
			name: "TestIndexResetFilterableAttributesCustomClient",
			args: args{
				UID:    "indexUID",
				client: customMeili,
			},
			wantTask: &TaskInfo{
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
	customMeili := setup(t, "", WithCustomClientWithTLS(&tls.Config{
		InsecureSkipVerify: true,
	}))

	type args struct {
		UID    string
		client ServiceManager
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
				client: meili,
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
				client: customMeili,
			},
			wantTask: &TaskInfo{
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
	customMeili := setup(t, "", WithCustomClientWithTLS(&tls.Config{
		InsecureSkipVerify: true,
	}))

	type args struct {
		UID    string
		client ServiceManager
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
				client: meili,
			},
			wantTask: &TaskInfo{
				TaskUID: 1,
			},
		},
		{
			name: "TestIndexResetDistinctAttributeWithCustomClient",
			args: args{
				UID:    "indexUID",
				client: customMeili,
			},
			wantTask: &TaskInfo{
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
	customMeili := setup(t, "", WithCustomClientWithTLS(&tls.Config{
		InsecureSkipVerify: true,
	}))

	type args struct {
		UID    string
		client ServiceManager
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
				client: meili,
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
				client: customMeili,
			},
			wantTask: &TaskInfo{
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
	customMeili := setup(t, "", WithCustomClientWithTLS(&tls.Config{
		InsecureSkipVerify: true,
	}))

	type args struct {
		UID    string
		client ServiceManager
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
				client: meili,
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
				client: customMeili,
			},
			wantTask: &TaskInfo{
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
	customMeili := setup(t, "", WithCustomClientWithTLS(&tls.Config{
		InsecureSkipVerify: true,
	}))

	type args struct {
		UID    string
		client ServiceManager
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
				client: meili,
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
				ProximityPrecision:   ByWord,
				SeparatorTokens:      make([]string, 0),
				NonSeparatorTokens:   make([]string, 0),
				Dictionary:           make([]string, 0),
				LocalizedAttributes:  nil,
			},
		},
		{
			name: "TestIndexResetSettingsWithCustomClient",
			args: args{
				UID:    "indexUID",
				client: customMeili,
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
				ProximityPrecision:   ByWord,
				SeparatorTokens:      make([]string, 0),
				NonSeparatorTokens:   make([]string, 0),
				Dictionary:           make([]string, 0),
				LocalizedAttributes:  nil,
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
	customMeili := setup(t, "", WithCustomClientWithTLS(&tls.Config{
		InsecureSkipVerify: true,
	}))

	type args struct {
		UID    string
		client ServiceManager
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
				client: meili,
			},
			wantTask: &TaskInfo{
				TaskUID: 1,
			},
		},
		{
			name: "TestIndexResetStopWordsWithCustomClient",
			args: args{
				UID:    "indexUID",
				client: customMeili,
			},
			wantTask: &TaskInfo{
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
	customMeili := setup(t, "", WithCustomClientWithTLS(&tls.Config{
		InsecureSkipVerify: true,
	}))

	type args struct {
		UID    string
		client ServiceManager
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
				client: meili,
			},
			wantTask: &TaskInfo{
				TaskUID: 1,
			},
		},
		{
			name: "TestIndexResetSynonymsWithCustomClient",
			args: args{
				UID:    "indexUID",
				client: customMeili,
			},
			wantTask: &TaskInfo{
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
	customMeili := setup(t, "", WithCustomClientWithTLS(&tls.Config{
		InsecureSkipVerify: true,
	}))

	type args struct {
		UID    string
		client ServiceManager
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
				client: meili,
			},
			wantTask: &TaskInfo{
				TaskUID: 1,
			},
		},
		{
			name: "TestIndexResetSortableAttributesCustomClient",
			args: args{
				UID:    "indexUID",
				client: customMeili,
			},
			wantTask: &TaskInfo{
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
	customMeili := setup(t, "", WithCustomClientWithTLS(&tls.Config{
		InsecureSkipVerify: true,
	}))

	type args struct {
		UID    string
		client ServiceManager
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
				client: meili,
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
				client: customMeili,
			},
			wantTask: &TaskInfo{
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
	customMeili := setup(t, "", WithCustomClientWithTLS(&tls.Config{
		InsecureSkipVerify: true,
	}))

	type args struct {
		UID    string
		client ServiceManager
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
				client: meili,
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
				client: customMeili,
			},
			wantTask: &TaskInfo{
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
	customMeili := setup(t, "", WithCustomClientWithTLS(&tls.Config{
		InsecureSkipVerify: true,
	}))

	type args struct {
		UID    string
		client ServiceManager
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
				client: meili,
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
				client: customMeili,
			},
			wantTask: &TaskInfo{
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
	customMeili := setup(t, "", WithCustomClientWithTLS(&tls.Config{
		InsecureSkipVerify: true,
	}))

	type args struct {
		UID     string
		client  ServiceManager
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
				client: meili,
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
				client: customMeili,
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
	customMeili := setup(t, "", WithCustomClientWithTLS(&tls.Config{
		InsecureSkipVerify: true,
	}))

	type args struct {
		UID     string
		client  ServiceManager
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
				client: meili,
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
				client: customMeili,
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
	customMeili := setup(t, "", WithCustomClientWithTLS(&tls.Config{
		InsecureSkipVerify: true,
	}))

	type args struct {
		UID     string
		client  ServiceManager
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
				client:  meili,
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
				client:  customMeili,
				request: "movie_id",
			},
			wantTask: &TaskInfo{
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
	customMeili := setup(t, "", WithCustomClientWithTLS(&tls.Config{
		InsecureSkipVerify: true,
	}))

	type args struct {
		UID     string
		client  ServiceManager
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
				client: meili,
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
				client: customMeili,
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
				client: meili,
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
	customMeili := setup(t, "", WithCustomClientWithTLS(&tls.Config{
		InsecureSkipVerify: true,
	}))

	type args struct {
		UID     string
		client  ServiceManager
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
				client: meili,
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
				client: customMeili,
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
	customMeili := setup(t, "", WithCustomClientWithTLS(&tls.Config{
		InsecureSkipVerify: true,
	}))

	type args struct {
		UID     string
		client  ServiceManager
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
				client: meili,
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
						SortFacetValuesBy: map[string]SortFacetType{
							"*": SortFacetTypeAlpha,
						},
					},
					SearchCutoffMs:     150,
					ProximityPrecision: ByAttribute,
					SeparatorTokens:    make([]string, 0),
					NonSeparatorTokens: make([]string, 0),
					Dictionary:         make([]string, 0),
					LocalizedAttributes: []*LocalizedAttributes{
						{
							Locales:           []string{"jpn", "eng"},
							AttributePatterns: []string{"*_ja"},
						},
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
				SearchCutoffMs:       150,
				ProximityPrecision:   ByAttribute,
				SeparatorTokens:      make([]string, 0),
				NonSeparatorTokens:   make([]string, 0),
				Dictionary:           make([]string, 0),
				LocalizedAttributes: []*LocalizedAttributes{
					{
						Locales:           []string{"jpn", "eng"},
						AttributePatterns: []string{"*_ja"},
					},
				},
			},
		},
		{
			name: "TestIndexUpdateSettingsWithCustomClient",
			args: args{
				UID:    "indexUID",
				client: customMeili,
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
						SortFacetValuesBy: map[string]SortFacetType{
							"*": SortFacetTypeAlpha,
						},
					},
					SearchCutoffMs:     150,
					ProximityPrecision: ByWord,
					SeparatorTokens:    make([]string, 0),
					NonSeparatorTokens: make([]string, 0),
					Dictionary:         make([]string, 0),
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
				SearchCutoffMs:       150,
				ProximityPrecision:   ByWord,
				SeparatorTokens:      make([]string, 0),
				NonSeparatorTokens:   make([]string, 0),
				Dictionary:           make([]string, 0),
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
	customMeili := setup(t, "", WithCustomClientWithTLS(&tls.Config{
		InsecureSkipVerify: true,
	}))

	type args struct {
		UID            string
		client         ServiceManager
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
				client: meili,
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
					ProximityPrecision:   ByWord,
					SeparatorTokens:      make([]string, 0),
					NonSeparatorTokens:   make([]string, 0),
					Dictionary:           make([]string, 0),
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
					ProximityPrecision:   ByWord,
					SeparatorTokens:      make([]string, 0),
					NonSeparatorTokens:   make([]string, 0),
					Dictionary:           make([]string, 0),
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
				ProximityPrecision:   ByWord,
				SeparatorTokens:      make([]string, 0),
				NonSeparatorTokens:   make([]string, 0),
				Dictionary:           make([]string, 0),
			},
		},
		{
			name: "TestIndexUpdateJustSynonymsWithCustomClient",
			args: args{
				UID:    "indexUID",
				client: customMeili,
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
					ProximityPrecision:   ByWord,
					SeparatorTokens:      make([]string, 0),
					NonSeparatorTokens:   make([]string, 0),
					Dictionary:           make([]string, 0),
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
					ProximityPrecision:   ByWord,
					SeparatorTokens:      make([]string, 0),
					NonSeparatorTokens:   make([]string, 0),
					Dictionary:           make([]string, 0),
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
				ProximityPrecision:   ByWord,
				SeparatorTokens:      make([]string, 0),
				NonSeparatorTokens:   make([]string, 0),
				Dictionary:           make([]string, 0),
			},
		},
		{
			name: "TestIndexUpdateJustSearchableAttributes",
			args: args{
				UID:    "indexUID",
				client: meili,
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
					ProximityPrecision:   ByWord,
					SeparatorTokens:      make([]string, 0),
					NonSeparatorTokens:   make([]string, 0),
					Dictionary:           make([]string, 0),
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
					ProximityPrecision:   ByWord,
					SeparatorTokens:      make([]string, 0),
					NonSeparatorTokens:   make([]string, 0),
					Dictionary:           make([]string, 0),
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
				ProximityPrecision:   ByWord,
				SeparatorTokens:      make([]string, 0),
				NonSeparatorTokens:   make([]string, 0),
				Dictionary:           make([]string, 0),
			},
		},
		{
			name: "TestIndexUpdateJustDisplayedAttributes",
			args: args{
				UID:    "indexUID",
				client: meili,
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
					ProximityPrecision:   ByWord,
					SeparatorTokens:      make([]string, 0),
					NonSeparatorTokens:   make([]string, 0),
					Dictionary:           make([]string, 0),
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
					ProximityPrecision:   ByWord,
					SeparatorTokens:      make([]string, 0),
					NonSeparatorTokens:   make([]string, 0),
					Dictionary:           make([]string, 0),
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
				ProximityPrecision:   ByWord,
				SeparatorTokens:      make([]string, 0),
				NonSeparatorTokens:   make([]string, 0),
				Dictionary:           make([]string, 0),
			},
		},
		{
			name: "TestIndexUpdateJustStopWords",
			args: args{
				UID:    "indexUID",
				client: meili,
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
					ProximityPrecision:   ByWord,
					SeparatorTokens:      make([]string, 0),
					NonSeparatorTokens:   make([]string, 0),
					Dictionary:           make([]string, 0),
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
					ProximityPrecision:   ByWord,
					SeparatorTokens:      make([]string, 0),
					NonSeparatorTokens:   make([]string, 0),
					Dictionary:           make([]string, 0),
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
				ProximityPrecision:   ByWord,
				SeparatorTokens:      make([]string, 0),
				NonSeparatorTokens:   make([]string, 0),
				Dictionary:           make([]string, 0),
			},
		},
		{
			name: "TestIndexUpdateJustFilterableAttributes",
			args: args{
				UID:    "indexUID",
				client: meili,
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
					ProximityPrecision: ByWord,
					SeparatorTokens:    make([]string, 0),
					NonSeparatorTokens: make([]string, 0),
					Dictionary:         make([]string, 0),
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
					ProximityPrecision: ByWord,
					SeparatorTokens:    make([]string, 0),
					NonSeparatorTokens: make([]string, 0),
					Dictionary:         make([]string, 0),
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
				ProximityPrecision:   ByWord,
				SeparatorTokens:      make([]string, 0),
				NonSeparatorTokens:   make([]string, 0),
				Dictionary:           make([]string, 0),
			},
		},
		{
			name: "TestIndexUpdateJustSortableAttributes",
			args: args{
				UID:    "indexUID",
				client: meili,
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
					TypoTolerance:      &defaultTypoTolerance,
					Pagination:         &defaultPagination,
					Faceting:           &defaultFaceting,
					ProximityPrecision: ByWord,
					SeparatorTokens:    make([]string, 0),
					NonSeparatorTokens: make([]string, 0),
					Dictionary:         make([]string, 0),
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
					TypoTolerance:      &defaultTypoTolerance,
					Pagination:         &defaultPagination,
					Faceting:           &defaultFaceting,
					ProximityPrecision: ByWord,
					SeparatorTokens:    make([]string, 0),
					NonSeparatorTokens: make([]string, 0),
					Dictionary:         make([]string, 0),
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
				ProximityPrecision:   ByWord,
				SeparatorTokens:      make([]string, 0),
				NonSeparatorTokens:   make([]string, 0),
				Dictionary:           make([]string, 0),
			},
		},
		{
			name: "TestIndexUpdateJustTypoTolerance",
			args: args{
				UID:    "indexUID",
				client: meili,
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
					ProximityPrecision:   ByWord,
					SeparatorTokens:      make([]string, 0),
					NonSeparatorTokens:   make([]string, 0),
					Dictionary:           make([]string, 0),
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
					Pagination:         &defaultPagination,
					Faceting:           &defaultFaceting,
					ProximityPrecision: ByWord,
					SeparatorTokens:    make([]string, 0),
					NonSeparatorTokens: make([]string, 0),
					Dictionary:         make([]string, 0),
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
				ProximityPrecision:   ByWord,
				SeparatorTokens:      make([]string, 0),
				NonSeparatorTokens:   make([]string, 0),
				Dictionary:           make([]string, 0),
			},
		},
		{
			name: "TestIndexUpdateJustPagination",
			args: args{
				UID:    "indexUID",
				client: meili,
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
					Faceting:           &defaultFaceting,
					ProximityPrecision: ByWord,
					SeparatorTokens:    make([]string, 0),
					NonSeparatorTokens: make([]string, 0),
					Dictionary:         make([]string, 0),
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
					Faceting:           &defaultFaceting,
					ProximityPrecision: ByWord,
					SeparatorTokens:    make([]string, 0),
					NonSeparatorTokens: make([]string, 0),
					Dictionary:         make([]string, 0),
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
				ProximityPrecision:   ByWord,
				SeparatorTokens:      make([]string, 0),
				NonSeparatorTokens:   make([]string, 0),
				Dictionary:           make([]string, 0),
			},
		},
		{
			name: "TestIndexUpdateJustFaceting",
			args: args{
				UID:    "indexUID",
				client: meili,
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
					ProximityPrecision:   ByWord,
					Faceting: &Faceting{
						MaxValuesPerFacet: 200,
						SortFacetValuesBy: map[string]SortFacetType{
							"*": SortFacetTypeAlpha,
						},
					},
					SeparatorTokens:    make([]string, 0),
					NonSeparatorTokens: make([]string, 0),
					Dictionary:         make([]string, 0),
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
					ProximityPrecision:   ByWord,
					Faceting: &Faceting{
						MaxValuesPerFacet: 200,
						SortFacetValuesBy: map[string]SortFacetType{
							"*": SortFacetTypeAlpha,
						},
					},
					SeparatorTokens:    make([]string, 0),
					NonSeparatorTokens: make([]string, 0),
					Dictionary:         make([]string, 0),
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
				ProximityPrecision:   ByWord,
				SeparatorTokens:      make([]string, 0),
				NonSeparatorTokens:   make([]string, 0),
				Dictionary:           make([]string, 0),
			},
		},
		{
			name: "TestIndexUpdateJustProximityPrecision",
			args: args{
				UID:    "indexUID",
				client: meili,
				firstRequest: Settings{
					RankingRules: []string{
						"typo", "words",
					},
					ProximityPrecision: ByAttribute,
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
					ProximityPrecision:   ByAttribute,
					Faceting: &Faceting{
						MaxValuesPerFacet: 100,
						SortFacetValuesBy: map[string]SortFacetType{
							"*": SortFacetTypeAlpha,
						},
					},
					SeparatorTokens:    make([]string, 0),
					NonSeparatorTokens: make([]string, 0),
					Dictionary:         make([]string, 0),
				},
				secondRequest: Settings{
					RankingRules: []string{
						"typo", "words",
					},
					ProximityPrecision: ByWord,
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
					ProximityPrecision:   ByWord,
					Faceting: &Faceting{
						MaxValuesPerFacet: 100,
						SortFacetValuesBy: map[string]SortFacetType{
							"*": SortFacetTypeAlpha,
						},
					},
					SeparatorTokens:    make([]string, 0),
					NonSeparatorTokens: make([]string, 0),
					Dictionary:         make([]string, 0),
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
				ProximityPrecision:   ByWord,
				SeparatorTokens:      make([]string, 0),
				NonSeparatorTokens:   make([]string, 0),
				Dictionary:           make([]string, 0),
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
	meili := setup(t, "")
	customMeili := setup(t, "", WithCustomClientWithTLS(&tls.Config{
		InsecureSkipVerify: true,
	}))

	type args struct {
		UID     string
		client  ServiceManager
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
				client: meili,
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
				client: customMeili,
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
	customMeili := setup(t, "", WithCustomClientWithTLS(&tls.Config{
		InsecureSkipVerify: true,
	}))

	type args struct {
		UID     string
		client  ServiceManager
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
				client: meili,
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
				client: customMeili,
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
	customMeili := setup(t, "", WithCustomClientWithTLS(&tls.Config{
		InsecureSkipVerify: true,
	}))

	type args struct {
		UID     string
		client  ServiceManager
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
				client: meili,
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
				client: customMeili,
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
	customMeili := setup(t, "", WithCustomClientWithTLS(&tls.Config{
		InsecureSkipVerify: true,
	}))

	type args struct {
		UID     string
		client  ServiceManager
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
				client: meili,
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
				client: customMeili,
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
				client: meili,
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
				client: meili,
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
		{
			name: "TestIndexDisableTypoTolerance",
			args: args{
				UID:    "indexUID",
				client: meili,
				request: TypoTolerance{
					Enabled: false,
					MinWordSizeForTypos: MinWordSizeForTypos{
						OneTypo:  5,
						TwoTypos: 9,
					},
					DisableOnWords:      []string{},
					DisableOnAttributes: []string{},
				},
			},
			wantTask: &TaskInfo{
				TaskUID: 1,
			},
			wantResp: &TypoTolerance{
				Enabled: false,
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
			require.Equal(t, &tt.args.request, gotResp)
		})
	}
}

func TestIndex_UpdatePagination(t *testing.T) {
	meili := setup(t, "")
	customMeili := setup(t, "", WithCustomClientWithTLS(&tls.Config{
		InsecureSkipVerify: true,
	}))

	type args struct {
		UID     string
		client  ServiceManager
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
				client: meili,
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
				client: customMeili,
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
	customMeili := setup(t, "", WithCustomClientWithTLS(&tls.Config{
		InsecureSkipVerify: true,
	}))

	type args struct {
		UID     string
		client  ServiceManager
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
				client: meili,
				request: Faceting{
					MaxValuesPerFacet: 200,
					SortFacetValuesBy: map[string]SortFacetType{
						"*": SortFacetTypeAlpha,
					},
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
				client: customMeili,
				request: Faceting{
					MaxValuesPerFacet: 200,
					SortFacetValuesBy: map[string]SortFacetType{
						"*": SortFacetTypeAlpha,
					},
				},
			},
			wantTask: &TaskInfo{
				TaskUID: 1,
			},
			wantResp: &defaultFaceting,
		},
		{
			name: "TestIndexGetStartedFaceting",
			args: args{
				UID:    "indexUID",
				client: meili,
				request: Faceting{
					MaxValuesPerFacet: 2,
					SortFacetValuesBy: map[string]SortFacetType{
						"*": SortFacetTypeCount,
					},
				},
			},
			wantTask: &TaskInfo{
				TaskUID: 1,
			},
			wantResp: &Faceting{
				MaxValuesPerFacet: 2,
				SortFacetValuesBy: map[string]SortFacetType{
					"*": SortFacetTypeCount,
				},
			},
		},
		{
			name: "TestIndexSortFacetValuesByCountFaceting",
			args: args{
				UID:    "indexUID",
				client: meili,
				request: Faceting{
					SortFacetValuesBy: map[string]SortFacetType{
						"*":        SortFacetTypeAlpha,
						"indexUID": SortFacetTypeCount,
					},
				},
			},
			wantTask: &TaskInfo{
				TaskUID: 1,
			},
			wantResp: &Faceting{
				SortFacetValuesBy: map[string]SortFacetType{
					"*":        SortFacetTypeAlpha,
					"indexUID": SortFacetTypeCount,
				},
			},
		},
		{
			name: "TestIndexSortFacetValuesAllIndexFaceting",
			args: args{
				UID:    "indexUID",
				client: meili,
				request: Faceting{
					SortFacetValuesBy: map[string]SortFacetType{
						"*": SortFacetTypeCount,
					},
				},
			},
			wantTask: &TaskInfo{
				TaskUID: 1,
			},
			wantResp: &Faceting{
				SortFacetValuesBy: map[string]SortFacetType{
					"*": SortFacetTypeCount,
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
		client   ServiceManager
		request  Settings
		newIndex bool
	}
	tests := []struct {
		name          string
		args          args
		wantTask      *TaskInfo
		wantEmbedders map[string]Embedder
		wantErr       string
	}{
		{
			name: "TestIndexUpdateSettingsEmbeddersErr",
			args: args{
				UID:    "indexUID",
				client: meili,
				request: Settings{
					Embedders: map[string]Embedder{
						"default": {
							Source:           "openAi",
							DocumentTemplate: "{{doc.foobar}}",
						},
					},
				},
			},
			wantTask: &TaskInfo{
				TaskUID: 1,
			},
			wantErr: "foobar",
		},
		{
			name: "TestIndexUpdateSettingsEmbeddersErr",
			args: args{
				UID:    "indexUID",
				client: meili,
				request: Settings{
					Embedders: map[string]Embedder{
						"default": {
							Source:           "openAi",
							APIKey:           "xxx",
							Model:            "text-embedding-3-small",
							DocumentTemplate: "A movie titled '{{doc.title}}'",
						},
					},
				},
			},
			wantTask: &TaskInfo{
				TaskUID: 1,
			},
			wantErr: "Incorrect API key",
		},
		{
			name: "TestIndexUpdateSettingsEmbeddersUserProvided",
			args: args{
				newIndex: true,
				UID:      "newIndexUID",
				client:   meili,
				request: Settings{
					Embedders: map[string]Embedder{
						"default": {
							Source:     "userProvided",
							Dimensions: 3,
						},
					},
				},
			},
			wantTask: &TaskInfo{
				TaskUID: 1,
			},
		},
		{
			name: "TestIndexUpdateSettingsEmbeddersWithRestSource",
			args: args{
				newIndex: true,
				UID:      "newIndexUID",
				client:   meili,
				request: Settings{
					Embedders: map[string]Embedder{
						"default": {
							Source:           "rest",
							URL:              "https://api.openai.com/v1/embeddings",
							APIKey:           "<your-openai-api-key>",
							Dimensions:       1536,
							DocumentTemplate: "A movie titled '{{doc.title}}' whose description starts with {{doc.overview|truncatewords: 20}}",
							Distribution: &Distribution{
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
			wantTask: &TaskInfo{
				TaskUID: 1,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := tt.args.client
			i, err := setUpIndexWithVector(c.(*meilisearch), tt.args.UID)
			require.NoError(t, err)
			t.Cleanup(cleanup(c))

			gotTask, err := i.UpdateSettings(&tt.args.request)
			require.NoError(t, err)
			require.GreaterOrEqual(t, gotTask.TaskUID, tt.wantTask.TaskUID)
			testWaitForTask(t, i, gotTask)

			gotResp, err := i.GetSettings()
			require.NoError(t, err)
			require.NotNil(t, gotResp)
		})
	}
}

func TestIndex_GetEmbedders(t *testing.T) {
	c := setup(t, "")
	t.Cleanup(cleanup(c))

	indexID := "newIndexUID"
	i := c.Index(indexID)
	task, err := c.CreateIndex(&IndexConfig{Uid: indexID})
	require.NoError(t, err)
	testWaitForTask(t, i, task)

	expected := map[string]Embedder{
		"default": {
			Source:     "userProvided",
			Dimensions: 3,
		},
	}
	task, err = i.UpdateSettings(&Settings{
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
	taskInfo, err := c.CreateIndex(&IndexConfig{Uid: indexID})
	require.NoError(t, err)
	testWaitForTask(t, i, taskInfo)

	embedders := map[string]Embedder{
		"someEmbbeder": {
			Source:     "userProvided",
			Dimensions: 3,
		},
	}
	taskInfo, err = i.UpdateSettings(&Settings{
		Embedders: embedders,
	})
	require.NoError(t, err)
	testWaitForTask(t, i, taskInfo)

	updated := map[string]Embedder{
		"someEmbbeder": {
			Source:     "userProvided",
			Dimensions: 5,
		},
	}

	taskInfo, err = i.UpdateEmbedders(updated)
	require.NoError(t, err)
	task, err := i.WaitForTask(taskInfo.TaskUID, 0)
	require.NoError(t, err)
	require.Equal(t, TaskStatusSucceeded, task.Status)

	got, err := i.GetEmbedders()
	require.NoError(t, err)
	require.Equal(t, updated, got)
}

func TestIndex_ResetEmbedders(t *testing.T) {
	c := setup(t, "")
	t.Cleanup(cleanup(c))

	indexID := "newIndexUID"
	i := c.Index(indexID)
	taskInfo, err := c.CreateIndex(&IndexConfig{Uid: indexID})
	require.NoError(t, err)
	testWaitForTask(t, i, taskInfo)

	taskInfo, err = i.UpdateSettings(&Settings{
		Embedders: map[string]Embedder{
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
	require.Equal(t, TaskStatusSucceeded, task.Status)

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
	taskInfo, err := c.CreateIndex(&IndexConfig{Uid: indexID})
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
	require.Equal(t, ByWord, got)

	task, err := i.UpdateProximityPrecision(ByAttribute)
	require.NoError(t, err)
	testWaitForTask(t, i, task)

	got, err = i.GetProximityPrecision()
	require.NoError(t, err)
	require.Equal(t, ByAttribute, got)

	task, err = i.ResetProximityPrecision()
	require.NoError(t, err)
	testWaitForTask(t, i, task)

	got, err = i.GetProximityPrecision()
	require.NoError(t, err)
	require.Equal(t, ByWord, got)
}

func Test_LocalizedAttributes(t *testing.T) {
	c := setup(t, "")
	t.Cleanup(cleanup(c))

	indexID := "newIndexUID"
	i := c.Index(indexID)
	taskInfo, err := c.CreateIndex(&IndexConfig{Uid: indexID})
	require.NoError(t, err)
	testWaitForTask(t, i, taskInfo)

	defer t.Cleanup(cleanup(c))

	t.Run("Test valid locate", func(t *testing.T) {
		got, err := i.GetLocalizedAttributes()
		require.NoError(t, err)
		require.Len(t, got, 0)

		localized := &LocalizedAttributes{
			Locales:           []string{"jpn", "eng"},
			AttributePatterns: []string{"*_ja"},
		}

		task, err := i.UpdateLocalizedAttributes([]*LocalizedAttributes{localized})
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
		invalidLocalized := &LocalizedAttributes{
			Locales:           []string{"foo"},
			AttributePatterns: []string{"*_ja"},
		}

		_, err := i.UpdateLocalizedAttributes([]*LocalizedAttributes{invalidLocalized})
		require.Error(t, err)
	})
}
