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
		_, err := sv.UpdateSearchRule(uid, &meilisearch.SearchRulesRequest{
			Description: fmt.Sprintf("Rule %d for %s", i, uid),
			Precedence:  intPtr(i + 1),
			Active:      boolPtr(true),
			Conditions: &meilisearch.Condition{
				Query: &meilisearch.QueryCondition{
					IsEmpty: boolPtr(true),
				},
				Time: &meilisearch.TimeCondition{
					Start: &start,
					End:   &end,
				},
			},
			Actions: []meilisearch.Action{
				{
					Selector: meilisearch.Selector{
						IndexUid: "products",
						ID:       fmt.Sprintf("%d", i+1),
					},
					Action: meilisearch.ActionDef{
						Type:     "pin",
						Position: 1,
					},
				},
			},
		})
		require.NoError(t, err)
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
			assert.NotEmpty(t, rule.Actions)
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
				AttributePatterns: []string{"categories"},
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
		rule, err := sv.UpdateSearchRule(uid, &meilisearch.SearchRulesRequest{
			Description: "Promotional campaign rules",
			Precedence:    intPtr(10),
			Active:      boolPtr(true),
			Conditions: &meilisearch.Condition{
				Query: &meilisearch.QueryCondition{
					IsEmpty: boolPtr(true),
				},
				Time: &meilisearch.TimeCondition{
					Start: &start,
					End:   &end,
				},
			},
			Actions: []meilisearch.Action{
				{
					Selector: meilisearch.Selector{
						IndexUid: "products",
						ID:       "456",
					},
					Action: meilisearch.ActionDef{
						Type:     "pin",
						Position: 1,
					},
				},
			},
		})
		require.NoError(t, err)
		require.NotNil(t, rule)
		assert.Equal(t, uid, rule.Uid)
		assert.Equal(t, "Promotional campaign rules", rule.Description)
		assert.Equal(t, 10, rule.Precedence)
		assert.True(t, rule.Active)
		assert.Len(t, rule.Conditions, 2)
		assert.Len(t, rule.Actions, 1)
	})

	t.Run("update existing rule", func(t *testing.T) {
		updatedRule, err := sv.UpdateSearchRule(uid, &meilisearch.SearchRulesRequest{
			Description: "Updated promotional campaign rules",
			Precedence:    intPtr(8),
			Active:      boolPtr(false),
			Conditions: &meilisearch.Condition{
				Query: &meilisearch.QueryCondition{
					IsEmpty: boolPtr(false),
				},
			},
			Actions: []meilisearch.Action{
				{
					Selector: meilisearch.Selector{
						IndexUid: "products",
						ID:       "789",
					},
					Action: meilisearch.ActionDef{
						Type:     "pin",
						Position: 2,
					},
				},
				{
					Selector: meilisearch.Selector{
						IndexUid: "categories",
						ID:       "001",
					},
					Action: meilisearch.ActionDef{
						Type:     "pin",
						Position: 1,
					},
				},
			},
		})
		require.NoError(t, err)
		require.NotNil(t, updatedRule)
		assert.Equal(t, uid, updatedRule.Uid)
		assert.Equal(t, "Updated promotional campaign rules", updatedRule.Description)
		assert.Equal(t, 8, updatedRule.Precedence)
		assert.False(t, updatedRule.Active)
		assert.Len(t, updatedRule.Conditions, 1)
		assert.Len(t, updatedRule.Actions, 2)

		// Verify update persisted
		fetchedRule, err := sv.GetSearchRule(uid)
		require.NoError(t, err)
		assert.Equal(t, updatedRule.Description, fetchedRule.Description)
		assert.Equal(t, updatedRule.Precedence, fetchedRule.Precedence)
		assert.Equal(t, updatedRule.Active, fetchedRule.Active)
	})

	t.Run("create rule with filter condition", func(t *testing.T) {
		rule, err := sv.UpdateSearchRule("test-filter", &meilisearch.SearchRulesRequest{
			Description: "Filter condition rule",
			Precedence: intPtr(1),
			Active:     boolPtr(true),
			Conditions: &meilisearch.Condition{
				Filter: &meilisearch.FilterCondition{
					Values: map[string]any{
						"color":    "red",
						"category": "shirt",
					},
				},
			},
			Actions: []meilisearch.Action{
				{
					Selector: meilisearch.Selector{
						IndexUid: "products",
						ID:       "123",
					},
					Action: meilisearch.ActionDef{
						Type:     "pin",
						Position: 1,
					},
				},
			},
		})

		require.NoError(t, err)
		require.NotNil(t, rule)

		assert.NotNil(t, rule.Conditions.Filter)
		assert.Equal(t, "red", rule.Conditions.Filter.Values["color"])
		assert.Equal(t, "shirt", rule.Conditions.Filter.Values["category"])
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

	want, err := sv.UpdateSearchRule(uid, &meilisearch.SearchRulesRequest{
		Description: "Black Friday 2025 rules",
		Precedence:    intPtr(5),
		Active:      boolPtr(true),
		Conditions: &meilisearch.Condition{
			Query: &meilisearch.QueryCondition{
				IsEmpty: boolPtr(true),
			},
			Time: &meilisearch.TimeCondition{
				Start: &start,
				End:   &end,
			},
		},
		Actions: []meilisearch.Action{
			{
				Selector: meilisearch.Selector{
					IndexUid: "products",
					ID:       "123",
				},
				Action: meilisearch.ActionDef{
					Type:     "pin",
					Position: 1,
				},
			},
		},
	})
	require.NoError(t, err)

	got, err := sv.GetSearchRule(uid)
	require.NoError(t, err)
	assert.Equal(t, want, got)
	assert.False(t, got.LastUpdatedAt.IsZero())
	assert.Equal(t, want.LastUpdatedAt, got.LastUpdatedAt)
}

func Test_DeleteSearchRule(t *testing.T) {
	sv := setup(t, "")
	t.Cleanup(cleanupSearchRules(sv))

	resp, err := sv.ExperimentalFeatures().SetDynamicSearchRules(true).Update()
	require.NoError(t, err)
	require.NotNil(t, resp)
	require.True(t, resp.DynamicSearchRules)

	uid := "black-friday"
	start := time.Now()
	end := start.Add(time.Hour * 1)
	_, err = sv.UpdateSearchRule(uid, &meilisearch.SearchRulesRequest{
		Description: "Black Friday 2025 rules",
		Precedence:    intPtr(5),
		Active:      boolPtr(true),
		Conditions: &meilisearch.Condition{
			Query: &meilisearch.QueryCondition{
				IsEmpty: boolPtr(true),
			},
			Time: &meilisearch.TimeCondition{
				Start: &start,
				End:   &end,
			},
		},
		Actions: []meilisearch.Action{
			{
				Selector: meilisearch.Selector{
					IndexUid: "products",
					ID:       "123",
				},
				Action: meilisearch.ActionDef{
					Type:     "pin",
					Position: 1,
				},
			},
		},
	})
	require.NoError(t, err)

	err = sv.DeleteSearchRule(uid)
	require.NoError(t, err)

	got, err := sv.GetSearchRule(uid)
	require.Error(t, err)
	assert.Nil(t, got)
}
