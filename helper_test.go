package meilisearch

import (
	"net/url"
	"reflect"
	"strconv"
	"testing"
	"time"
)

func (req *internalRequest) init() {
	req.withQueryParams = make(map[string]string)
}

func formatDateForComparison(date time.Time) string {
	const format = "2006-01-02T15:04:05Z"
	return date.Format(format)
}

func TestConvertKeyToParsedKey(t *testing.T) {
	key := Key{
		Name:        "test",
		Description: "test description",
		UID:         "123",
		Actions:     []string{"read", "write"},
		Indexes:     []string{"index1", "index2"},
		ExpiresAt:   time.Now(),
	}

	expectedExpiresAt := formatDateForComparison(key.ExpiresAt)
	parsedKey := convertKeyToParsedKey(key)

	if parsedKey.Name != key.Name ||
		parsedKey.Description != key.Description ||
		parsedKey.UID != key.UID ||
		!reflect.DeepEqual(parsedKey.Actions, key.Actions) ||
		!reflect.DeepEqual(parsedKey.Indexes, key.Indexes) ||
		parsedKey.ExpiresAt == nil || *parsedKey.ExpiresAt != expectedExpiresAt {
		t.Errorf("convertKeyToParsedKey(%v) = %v; want %v", key, parsedKey, key)
	}
}

func TestEncodeTasksQuery(t *testing.T) {
	param := &TasksQuery{
		Limit:            10,
		From:             5,
		Statuses:         []TaskStatus{"queued", "running"},
		Types:            []TaskType{"type1", "type2"},
		IndexUIDS:        []string{"uid1", "uid2"},
		UIDS:             []int64{1, 2, 3},
		CanceledBy:       []int64{4, 5},
		BeforeEnqueuedAt: time.Now().Add(-10 * time.Hour),
		AfterEnqueuedAt:  time.Now().Add(-20 * time.Hour),
		BeforeStartedAt:  time.Now().Add(-30 * time.Hour),
		AfterStartedAt:   time.Now().Add(-40 * time.Hour),
		BeforeFinishedAt: time.Now().Add(-50 * time.Hour),
		AfterFinishedAt:  time.Now().Add(-60 * time.Hour),
	}
	req := &internalRequest{}
	req.init()

	encodeTasksQuery(param, req)

	expectedParams := map[string]string{
		"limit":            strconv.FormatInt(param.Limit, 10),
		"from":             strconv.FormatInt(param.From, 10),
		"statuses":         "queued,running",
		"types":            "type1,type2",
		"indexUids":        "uid1,uid2",
		"uids":             "1,2,3",
		"canceledBy":       "4,5",
		"beforeEnqueuedAt": formatDateForComparison(param.BeforeEnqueuedAt),
		"afterEnqueuedAt":  formatDateForComparison(param.AfterEnqueuedAt),
		"beforeStartedAt":  formatDateForComparison(param.BeforeStartedAt),
		"afterStartedAt":   formatDateForComparison(param.AfterStartedAt),
		"beforeFinishedAt": formatDateForComparison(param.BeforeFinishedAt),
		"afterFinishedAt":  formatDateForComparison(param.AfterFinishedAt),
	}

	for k, v := range expectedParams {
		if req.withQueryParams[k] != v {
			t.Errorf("encodeTasksQuery() param %v = %v; want %v", k, req.withQueryParams[k], v)
		}
	}
}

func TestTransformStringVariadicToMap(t *testing.T) {
	tests := []struct {
		input  []string
		expect map[string]string
	}{
		{[]string{"primaryKey1"}, map[string]string{"primaryKey": "primaryKey1"}},
		{nil, nil},
	}

	for _, test := range tests {
		result := transformStringVariadicToMap(test.input...)
		if !reflect.DeepEqual(result, test.expect) {
			t.Errorf("transformStringVariadicToMap(%v) = %v; want %v", test.input, result, test.expect)
		}
	}
}

func TestGenerateQueryForOptions(t *testing.T) {
	options := map[string]string{
		"key1": "value1",
		"key2": "value2",
	}

	expected := url.Values{}
	expected.Add("key1", "value1")
	expected.Add("key2", "value2")

	result := generateQueryForOptions(options)
	if result != expected.Encode() {
		t.Errorf("generateQueryForOptions(%v) = %v; want %v", options, result, expected.Encode())
	}
}
