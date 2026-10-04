package integration

import (
	"fmt"
	"testing"
	"time"

	"github.com/meilisearch/meilisearch-go"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_ListSearchRule(t *testing.T) {
	sv := setup(t, "")
	t.Cleanup(cleanupSearchRules(sv))

	resp, err := sv.ExperimentalFeatures().SetDynamicSearchRules(true).Update()
	require.NoError(t, err)
	require.NotNil(t, resp)
	require.True(t, resp.DynamicSearchRules)

	uids := []string{"black-friday", "christmas-sale", "summer-deals"}
	start := time.Now().UTC().Truncate(time.Second)
	end := start.Add(time.Hour * 24)

	for i, uid := range uids {
		task, err := sv.UpdateSearchRule(uid, &meilisearch.SearchRulesRequest{
			Description: fmt.Sprintf("Rule %d for %s", i, uid),
			Precedence:  intPtr(i + 1),
			Active:      boolPtr(true),
			Conditions: &meilisearch.Conditions{
				Query: &meilisearch.QueryCondition{
					IsEmpty: boolPtr(true),
				},
				Time: &meilisearch.TimeCondition{
					Start: &start,
					End:   &end,
				},
			},
			Actions: &meilisearch.Actions{
				Pin: []meilisearch.Pin{{
					IndexUid: "products",
					ID:       fmt.Sprintf("%d", i+1),
					Position: 1,
				}},
				Scale: []meilisearch.Scale{{
					IndexUid: "products",
					Ids:      []string{fmt.Sprintf("%d", i+1)},
					Weight:   0.5,
				}},
			},
		})
		require.NoError(t, err)
		testWaitForTask(t, sv, task)
	}

	t.Run("list all rules with pagination", func(t *testing.T) {
		results, err := sv.ListSearchRules(&meilisearch.SearchRulesParams{
			Offset: 0,
			Limit:  20,
		})
		require.NoError(t, err)
		assert.NotNil(t, results)

		assert.Equal(t, results.Total, int64(len(uids)))
		assert.Equal(t, int64(0), results.Offset)
		assert.Equal(t, int64(20), results.Limit)
		assert.Len(t, results.Results, len(uids))

		// Verify results contain expected UIDs
		foundUIDs := make(map[string]bool)
		for _, rule := range results.Results {
			foundUIDs[rule.Uid] = true
			assert.NotEmpty(t, rule.Uid)
			assert.NotEmpty(t, rule.Description)
			assert.True(t, rule.Active)
			assert.Greater(t, rule.Precedence, 0)
			assert.NotEmpty(t, rule.Conditions)
			require.Len(t, rule.Actions.Pin, 1)
			assert.Equal(t, "products", rule.Actions.Pin[0].IndexUid)
			assert.Equal(t, 1, rule.Actions.Pin[0].Position)
			require.Len(t, rule.Actions.Scale, 1)
			assert.Equal(t, "products", rule.Actions.Scale[0].IndexUid)
			assert.Equal(t, []string{rule.Actions.Pin[0].ID}, rule.Actions.Scale[0].Ids)
			assert.Equal(t, 0.5, rule.Actions.Scale[0].Weight)
		}
		for _, uid := range uids {
			assert.True(t, foundUIDs[uid], "Expected to find rule with UID: %s", uid)
		}
	})

	t.Run("list rules with filter", func(t *testing.T) {
		results, err := sv.ListSearchRules(&meilisearch.SearchRulesParams{
			Offset: 0,
			Limit:  20,
			Filter: &meilisearch.SearchRulesFilter{
				Active: boolPtr(false),
			},
		})
		require.NoError(t, err)
		assert.NotNil(t, results)
		assert.Zero(t, results.Total)
	})

	t.Run("list rules with attribute patterns filter", func(t *testing.T) {
		results, err := sv.ListSearchRules(&meilisearch.SearchRulesParams{
			Offset: 0,
			Limit:  20,
			Filter: &meilisearch.SearchRulesFilter{
				Query: "categories",
			},
		})
		require.NoError(t, err)
		assert.NotNil(t, results)
		assert.Equal(t, len(results.Results), 0)
	})
}

func Test_UpdateSearchRule(t *testing.T) {
	sv := setup(t, "")
	t.Cleanup(cleanupSearchRules(sv))

	resp, err := sv.ExperimentalFeatures().SetDynamicSearchRules(true).Update()
	require.NoError(t, err)
	require.NotNil(t, resp)
	require.True(t, resp.DynamicSearchRules)

	uid := "promo-rule"
	start := time.Now().Truncate(time.Second)
	end := start.Add(time.Hour * 2)

	t.Run("create new rule", func(t *testing.T) {
		task, err := sv.UpdateSearchRule(uid, &meilisearch.SearchRulesRequest{
			Description: "Promotional campaign rules",
			Precedence:  intPtr(10),
			Active:      boolPtr(true),
			Conditions: &meilisearch.Conditions{
				Query: &meilisearch.QueryCondition{
					IsEmpty: boolPtr(true),
				},
				Time: &meilisearch.TimeCondition{
					Start: &start,
					End:   &end,
				},
			},
			Actions: &meilisearch.Actions{
				Pin: []meilisearch.Pin{{
					IndexUid: "products",
					ID:       "456",
					Position: 0,
				}},
			},
		})
		require.NoError(t, err)

		testWaitForTask(t, sv, task)

		rule, err := sv.GetSearchRule(uid)
		require.NoError(t, err)

		assert.Equal(t, uid, rule.Uid)
		assert.Equal(t, "Promotional campaign rules", rule.Description)
		assert.Equal(t, 10, rule.Precedence)
		assert.True(t, rule.Active)
		assert.Equal(t, []meilisearch.Pin{{IndexUid: "products", ID: "456", Position: 0}}, rule.Actions.Pin)
		assert.Empty(t, rule.Actions.Scale)
	})

	t.Run("update existing rule", func(t *testing.T) {
		task, err := sv.UpdateSearchRule(uid, &meilisearch.SearchRulesRequest{
			Description: "Updated promotional campaign rules",
			Precedence:  intPtr(8),
			Active:      boolPtr(false),
			Conditions: &meilisearch.Conditions{
				Query: &meilisearch.QueryCondition{
					IsEmpty: boolPtr(true),
				},
			},
			Actions: &meilisearch.Actions{
				Pin: []meilisearch.Pin{
					{IndexUid: "products", ID: "789", Position: 2},
					{IndexUid: "categories", ID: "001", Position: 1},
				},
			},
		})
		require.NoError(t, err)

		testWaitForTask(t, sv, task)

		rule, err := sv.GetSearchRule(uid)
		require.NoError(t, err)

		require.NotNil(t, rule)
		assert.Equal(t, uid, rule.Uid)
		assert.Equal(t, "Updated promotional campaign rules", rule.Description)
		assert.Equal(t, 8, rule.Precedence)
		assert.False(t, rule.Active)
		assert.Equal(t, []meilisearch.Pin{
			{IndexUid: "products", ID: "789", Position: 2},
			{IndexUid: "categories", ID: "001", Position: 1},
		}, rule.Actions.Pin)
		assert.Empty(t, rule.Actions.Scale)
	})

	t.Run("update description preserves actions", func(t *testing.T) {
		before, err := sv.GetSearchRule(uid)
		require.NoError(t, err)

		task, err := sv.UpdateSearchRule(uid, &meilisearch.SearchRulesRequest{
			Description: "Description-only update",
		})
		require.NoError(t, err)
		testWaitForTask(t, sv, task)

		rule, err := sv.GetSearchRule(uid)
		require.NoError(t, err)
		assert.Equal(t, "Description-only update", rule.Description)
		assert.Equal(t, before.Actions, rule.Actions)
	})

	t.Run("replace actions with scale actions", func(t *testing.T) {
		cases := []struct {
			name    string
			actions meilisearch.Actions
		}{
			{
				name: "boost by ids without index uid",
				actions: meilisearch.Actions{
					Scale: []meilisearch.Scale{{Weight: 2.5, Ids: []string{"123", "456"}}},
				},
			},
			{
				name: "demote by filter",
				actions: meilisearch.Actions{
					Scale: []meilisearch.Scale{{Weight: 0.5, Filter: "availability = out_of_stock", IndexUid: "products"}},
				},
			},
			{
				name: "hide by ids and array filter",
				actions: meilisearch.Actions{
					Scale: []meilisearch.Scale{{
						Weight: 0,
						Ids:    []string{"123"},
						Filter: []interface{}{[]interface{}{"availability = discontinued", "availability = archived"}},
					}},
				},
			},
			{
				name: "pin and scale together",
				actions: meilisearch.Actions{
					Pin: []meilisearch.Pin{{ID: "123", Position: 0}},
					Scale: []meilisearch.Scale{
						{Weight: 0, Ids: []string{"123"}},
						{Weight: 1.5, Filter: "series = batman", IndexUid: "products"},
					},
				},
			},
		}

		for _, tc := range cases {
			t.Run(tc.name, func(t *testing.T) {
				task, err := sv.UpdateSearchRule(uid, &meilisearch.SearchRulesRequest{Actions: &tc.actions})
				require.NoError(t, err)
				testWaitForTask(t, sv, task)

				rule, err := sv.GetSearchRule(uid)
				require.NoError(t, err)
				assert.Equal(t, tc.actions, rule.Actions)
			})
		}
	})

	t.Run("clear actions", func(t *testing.T) {
		task, err := sv.UpdateSearchRule(uid, &meilisearch.SearchRulesRequest{Actions: &meilisearch.Actions{}})
		require.NoError(t, err)
		testWaitForTask(t, sv, task)

		rule, err := sv.GetSearchRule(uid)
		require.NoError(t, err)
		assert.Empty(t, rule.Actions.Pin)
		assert.Empty(t, rule.Actions.Scale)
	})

	t.Run("reject reserved rule uid", func(t *testing.T) {
		task, err := sv.UpdateSearchRule("__meilisearch_metadata", &meilisearch.SearchRulesRequest{
			Actions: &meilisearch.Actions{Pin: []meilisearch.Pin{{ID: "123", Position: 0}}},
		})
		require.Error(t, err)
		assert.Nil(t, task)
		assert.Contains(t, err.Error(), "reserved")
	})
}

func Test_GetSearchRule(t *testing.T) {
	sv := setup(t, "")
	t.Cleanup(cleanupSearchRules(sv))

	resp, err := sv.ExperimentalFeatures().SetDynamicSearchRules(true).Update()
	require.NoError(t, err)
	require.NotNil(t, resp)
	require.True(t, resp.DynamicSearchRules)

	uid := "black-friday"
	start := time.Now()
	end := start.Add(time.Hour * 1)

	task, err := sv.UpdateSearchRule(uid, &meilisearch.SearchRulesRequest{
		Description: "Black Friday 2025 rules",
		Precedence:  intPtr(5),
		Active:      boolPtr(true),
		Conditions: &meilisearch.Conditions{
			Query: &meilisearch.QueryCondition{
				IsEmpty: boolPtr(true),
			},
			Time: &meilisearch.TimeCondition{
				Start: &start,
				End:   &end,
			},
		},
		Actions: &meilisearch.Actions{
			Pin: []meilisearch.Pin{{IndexUid: "products", ID: "123", Position: 1}},
		},
	})
	require.NoError(t, err)

	testWaitForTask(t, sv, task)

	rule, err := sv.GetSearchRule(uid)
	require.NoError(t, err)
	require.NotNil(t, rule)
	assert.Equal(t, uid, rule.Uid)
	assert.Equal(t, []meilisearch.Pin{{IndexUid: "products", ID: "123", Position: 1}}, rule.Actions.Pin)
	assert.Empty(t, rule.Actions.Scale)
}

func Test_DeleteSearchRule(t *testing.T) {
	sv := setup(t, "")
	t.Cleanup(cleanupSearchRules(sv))

	resp, err := sv.ExperimentalFeatures().SetDynamicSearchRules(true).Update()
	require.NoError(t, err)
	require.NotNil(t, resp)
	require.True(t, resp.DynamicSearchRules)

	t.Run("delete with a uid", func(t *testing.T) {
		uid := "black-friday"
		start := time.Now()
		end := start.Add(time.Hour * 1)

		task, err := sv.UpdateSearchRule(uid, &meilisearch.SearchRulesRequest{
			Description: "Black Friday 2025 rules",
			Precedence:  intPtr(5),
			Active:      boolPtr(true),
			Conditions: &meilisearch.Conditions{
				Query: &meilisearch.QueryCondition{
					IsEmpty: boolPtr(true),
				},
				Time: &meilisearch.TimeCondition{
					Start: &start,
					End:   &end,
				},
			},
			Actions: &meilisearch.Actions{
				Pin: []meilisearch.Pin{{IndexUid: "products", ID: "123", Position: 1}},
			},
		})
		require.NoError(t, err)

		testWaitForTask(t, sv, task)

		task, err = sv.DeleteSearchRule(&uid)
		require.NoError(t, err)

		testWaitForTask(t, sv, task)

		got, err := sv.GetSearchRule(uid)
		require.Error(t, err)
		assert.Nil(t, got)
	})

	t.Run("bulk delete", func(t *testing.T) {
		task, err := sv.DeleteSearchRule(nil)
		require.NoError(t, err)
		testWaitForTask(t, sv, task)
	})
}
