package integration

import (
	"fmt"
	"testing"
	"time"

	"github.com/meilisearch/meilisearch-go"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func updateSearchRuleAndWait(
	t *testing.T,
	sv meilisearch.ServiceManager,
	uid string,
	request *meilisearch.SearchRulesRequest,
) *meilisearch.TaskInfo {
	t.Helper()

	task, err := sv.UpdateSearchRule(uid, request)
	require.NoError(t, err)
	require.NotNil(t, task)
	require.Equal(t, meilisearch.TaskStatusEnqueued, task.Status)
	require.False(t, task.EnqueuedAt.IsZero())

	_, err = sv.WaitForTask(task.TaskUID, 0)
	require.NoError(t, err)
	return task
}

func Test_ListSearchRule(t *testing.T) {
	sv := setup(t, "")
	t.Cleanup(cleanupSearchRules(sv))

	resp, err := sv.ExperimentalFeatures().SetDynamicSearchRules(true).Update()
	require.NoError(t, err)
	require.NotNil(t, resp)
	require.True(t, resp.DynamicSearchRules)

	uids := []string{"black-friday", "christmas-sale", "summer-deals"}
	start := time.Now().UTC().Truncate(time.Second)
	end := start.Add(24 * time.Hour)

	for i, uid := range uids {
		words := fmt.Sprintf("%s promotion", uid)
		updateSearchRuleAndWait(t, sv, uid, &meilisearch.SearchRulesRequest{
			Description: fmt.Sprintf("Rule %d for %s", i, uid),
			Precedence:  intPtr(i + 1),
			Active:      boolPtr(true),
			Conditions: &meilisearch.SearchRuleConditions{
				Query: &meilisearch.QueryCondition{
					IsEmpty: boolPtr(false),
					Words:   &words,
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
	}

	t.Run("list all rules", func(t *testing.T) {
		results, err := sv.ListSearchRules(&meilisearch.SearchRulesParams{
			Offset: 0,
			Limit:  20,
		})
		require.NoError(t, err)
		require.NotNil(t, results)
		assert.Equal(t, int64(len(uids)), results.Total)
		assert.Equal(t, int64(0), results.Offset)
		assert.Equal(t, int64(20), results.Limit)
		assert.Len(t, results.Results, len(uids))

		foundUIDs := make(map[string]bool)
		for _, rule := range results.Results {
			foundUIDs[rule.Uid] = true
			assert.NotEmpty(t, rule.Description)
			assert.True(t, rule.Active)
			assert.Greater(t, rule.Precedence, 0)
			require.NotNil(t, rule.Conditions.Query)
			assert.NotNil(t, rule.Conditions.Query.Words)
			assert.NotNil(t, rule.Conditions.Time)
			assert.NotEmpty(t, rule.Actions)
		}
		for _, uid := range uids {
			assert.True(t, foundUIDs[uid], "expected to find rule with UID: %s", uid)
		}
	})

	t.Run("paginate rules", func(t *testing.T) {
		results, err := sv.ListSearchRules(&meilisearch.SearchRulesParams{
			Offset: 1,
			Limit:  1,
		})
		require.NoError(t, err)
		require.NotNil(t, results)
		assert.Equal(t, int64(1), results.Offset)
		assert.Equal(t, int64(1), results.Limit)
		assert.Equal(t, int64(len(uids)), results.Total)
		assert.Len(t, results.Results, 1)
	})

	t.Run("filter rules by active status", func(t *testing.T) {
		results, err := sv.ListSearchRules(&meilisearch.SearchRulesParams{
			Filter: &meilisearch.SearchRulesFilter{Active: boolPtr(false)},
		})
		require.NoError(t, err)
		require.NotNil(t, results)
		assert.Zero(t, results.Total)
	})

	t.Run("filter rules by query", func(t *testing.T) {
		results, err := sv.ListSearchRules(&meilisearch.SearchRulesParams{
			Filter: &meilisearch.SearchRulesFilter{Query: "christmas"},
		})
		require.NoError(t, err)
		require.NotNil(t, results)
		require.Len(t, results.Results, 1)
		assert.Equal(t, "christmas-sale", results.Results[0].Uid)
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
	start := time.Now().UTC().Truncate(time.Second)
	end := start.Add(2 * time.Hour)
	createWords := "special offer"

	createTask := updateSearchRuleAndWait(t, sv, uid, &meilisearch.SearchRulesRequest{
		Description: "Promotional campaign rules",
		Precedence:  intPtr(10),
		Active:      boolPtr(true),
		Conditions: &meilisearch.SearchRuleConditions{
			Query: &meilisearch.QueryCondition{
				IsEmpty: boolPtr(false),
				Words:   &createWords,
			},
			Time: &meilisearch.TimeCondition{Start: &start, End: &end},
		},
		Actions: []meilisearch.Action{
			{
				Selector: meilisearch.Selector{IndexUid: "products", ID: "456"},
				Action:   meilisearch.ActionDef{Type: "pin", Position: 1},
			},
		},
	})
	assert.Equal(t, meilisearch.TaskStatusEnqueued, createTask.Status)

	createdRule, err := sv.GetSearchRule(uid)
	require.NoError(t, err)
	require.NotNil(t, createdRule)
	assert.Equal(t, "Promotional campaign rules", createdRule.Description)
	assert.Equal(t, 10, createdRule.Precedence)
	assert.True(t, createdRule.Active)
	require.NotNil(t, createdRule.Conditions.Query)
	assert.Equal(t, createWords, *createdRule.Conditions.Query.Words)
	assert.Len(t, createdRule.Actions, 1)

	updateWords := "updated promotion"
	updateTask := updateSearchRuleAndWait(t, sv, uid, &meilisearch.SearchRulesRequest{
		Description: "Updated promotional campaign rules",
		Precedence:  intPtr(8),
		Active:      boolPtr(false),
		Conditions: &meilisearch.SearchRuleConditions{
			Query: &meilisearch.QueryCondition{
				IsEmpty: boolPtr(false),
				Words:   &updateWords,
			},
		},
		Actions: []meilisearch.Action{
			{
				Selector: meilisearch.Selector{IndexUid: "products", ID: "789"},
				Action:   meilisearch.ActionDef{Type: "pin", Position: 2},
			},
			{
				Selector: meilisearch.Selector{IndexUid: "categories", ID: "001"},
				Action:   meilisearch.ActionDef{Type: "pin", Position: 1},
			},
		},
	})
	assert.Equal(t, meilisearch.TaskStatusEnqueued, updateTask.Status)

	updatedRule, err := sv.GetSearchRule(uid)
	require.NoError(t, err)
	require.NotNil(t, updatedRule)
	assert.Equal(t, "Updated promotional campaign rules", updatedRule.Description)
	assert.Equal(t, 8, updatedRule.Precedence)
	assert.False(t, updatedRule.Active)
	require.NotNil(t, updatedRule.Conditions.Query)
	assert.Equal(t, updateWords, *updatedRule.Conditions.Query.Words)
	assert.Len(t, updatedRule.Actions, 2)
}

func Test_GetSearchRule(t *testing.T) {
	sv := setup(t, "")
	t.Cleanup(cleanupSearchRules(sv))

	resp, err := sv.ExperimentalFeatures().SetDynamicSearchRules(true).Update()
	require.NoError(t, err)
	require.NotNil(t, resp)
	require.True(t, resp.DynamicSearchRules)

	uid := "black-friday"
	updateSearchRuleAndWait(t, sv, uid, &meilisearch.SearchRulesRequest{
		Description: "Black Friday 2025 rules",
		Precedence:  intPtr(5),
		Active:      boolPtr(true),
		Conditions: &meilisearch.SearchRuleConditions{
			Query: &meilisearch.QueryCondition{IsEmpty: boolPtr(true)},
		},
		Actions: []meilisearch.Action{
			{
				Selector: meilisearch.Selector{IndexUid: "products", ID: "123"},
				Action:   meilisearch.ActionDef{Type: "pin", Position: 1},
			},
		},
	})

	got, err := sv.GetSearchRule(uid)
	require.NoError(t, err)
	require.NotNil(t, got)
	assert.Equal(t, uid, got.Uid)
	assert.Equal(t, "Black Friday 2025 rules", got.Description)
	assert.Equal(t, 5, got.Precedence)
	require.NotNil(t, got.Conditions.Query)
	assert.Equal(t, true, *got.Conditions.Query.IsEmpty)
}

func Test_DeleteSearchRule(t *testing.T) {
	sv := setup(t, "")
	t.Cleanup(cleanupSearchRules(sv))

	resp, err := sv.ExperimentalFeatures().SetDynamicSearchRules(true).Update()
	require.NoError(t, err)
	require.NotNil(t, resp)
	require.True(t, resp.DynamicSearchRules)

	uid := "black-friday"
	updateSearchRuleAndWait(t, sv, uid, &meilisearch.SearchRulesRequest{
		Conditions: &meilisearch.SearchRuleConditions{
			Query: &meilisearch.QueryCondition{IsEmpty: boolPtr(true)},
		},
		Actions: []meilisearch.Action{
			{
				Selector: meilisearch.Selector{IndexUid: "products", ID: "123"},
				Action:   meilisearch.ActionDef{Type: "pin", Position: 1},
			},
		},
	})

	deleteTask, err := sv.DeleteSearchRule(uid)
	require.NoError(t, err)
	require.NotNil(t, deleteTask)
	assert.Equal(t, meilisearch.TaskStatusEnqueued, deleteTask.Status)
	_, err = sv.WaitForTask(deleteTask.TaskUID, 0)
	require.NoError(t, err)

	got, err := sv.GetSearchRule(uid)
	require.Error(t, err)
	assert.Nil(t, got)

	missingTask, err := sv.DeleteSearchRule("missing-rule")
	require.NoError(t, err)
	require.NotNil(t, missingTask)
	_, err = sv.WaitForTask(missingTask.TaskUID, 0)
	require.NoError(t, err)
}

func Test_DeleteAllSearchRules(t *testing.T) {
	sv := setup(t, "")
	t.Cleanup(cleanupSearchRules(sv))

	resp, err := sv.ExperimentalFeatures().SetDynamicSearchRules(true).Update()
	require.NoError(t, err)
	require.NotNil(t, resp)
	require.True(t, resp.DynamicSearchRules)

	for _, uid := range []string{"rule-one", "rule-two"} {
		updateSearchRuleAndWait(t, sv, uid, &meilisearch.SearchRulesRequest{
			Actions: []meilisearch.Action{
				{
					Selector: meilisearch.Selector{ID: uid},
					Action:   meilisearch.ActionDef{Type: "pin", Position: 1},
				},
			},
		})
	}

	deleteTask, err := sv.DeleteAllSearchRules()
	require.NoError(t, err)
	require.NotNil(t, deleteTask)
	assert.Equal(t, meilisearch.TaskStatusEnqueued, deleteTask.Status)
	_, err = sv.WaitForTask(deleteTask.TaskUID, 0)
	require.NoError(t, err)

	results, err := sv.ListSearchRules(&meilisearch.SearchRulesParams{})
	require.NoError(t, err)
	require.NotNil(t, results)
	assert.Zero(t, results.Total)
}
