package meilisearch

import (
	"bytes"
	"context"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"time"
)

func IsValidUUID(uuid string) bool {
	r := regexp.MustCompile("^[a-fA-F0-9]{8}-[a-fA-F0-9]{4}-4[a-fA-F0-9]{3}-[8|9|aA|bB][a-fA-F0-9]{3}-[a-fA-F0-9]{12}$")
	return r.MatchString(uuid)
}

// This function allows the user to create a Key with an ExpiresAt in time.Time
// and transform the Key structure into a KeyParsed structure to send the time format
// managed by meilisearch
func convertKeyToParsedKey(key Key) (resp KeyParsed) {
	resp = KeyParsed{
		Name:        key.Name,
		Description: key.Description,
		UID:         key.UID,
		Actions:     key.Actions,
		Indexes:     key.Indexes,
	}

	// Convert time.Time to *string to feat the exact ISO-8601
	// format of meilisearch
	if !key.ExpiresAt.IsZero() {
		resp.ExpiresAt = formatDate(key.ExpiresAt, true)
	}
	return resp
}

func encodeTasksQuery(param *TasksQuery, req *internalRequest) {
	if param.Limit != 0 {
		req.withQueryParams["limit"] = strconv.FormatInt(param.Limit, 10)
	}
	if param.From != 0 {
		req.withQueryParams["from"] = strconv.FormatInt(param.From, 10)
	}
	if len(param.Statuses) != 0 {
		var statuses []string
		for _, status := range param.Statuses {
			statuses = append(statuses, string(status))
		}
		req.withQueryParams["statuses"] = strings.Join(statuses, ",")
	}
	if len(param.Types) != 0 {
		var types []string
		for _, t := range param.Types {
			types = append(types, string(t))
		}
		req.withQueryParams["types"] = strings.Join(types, ",")
	}
	if len(param.IndexUIDS) != 0 {
		req.withQueryParams["indexUids"] = strings.Join(param.IndexUIDS, ",")
	}
	if len(param.UIDS) != 0 {
		req.withQueryParams["uids"] = strings.Trim(strings.Join(strings.Fields(fmt.Sprint(param.UIDS)), ","), "[]")
	}
	if len(param.CanceledBy) != 0 {
		req.withQueryParams["canceledBy"] = strings.Trim(strings.Join(strings.Fields(fmt.Sprint(param.CanceledBy)), ","), "[]")
	}
	if !param.BeforeEnqueuedAt.IsZero() {
		req.withQueryParams["beforeEnqueuedAt"] = *formatDate(param.BeforeEnqueuedAt, false)
	}
	if !param.AfterEnqueuedAt.IsZero() {
		req.withQueryParams["afterEnqueuedAt"] = *formatDate(param.AfterEnqueuedAt, false)
	}
	if !param.BeforeStartedAt.IsZero() {
		req.withQueryParams["beforeStartedAt"] = *formatDate(param.BeforeStartedAt, false)
	}
	if !param.AfterStartedAt.IsZero() {
		req.withQueryParams["afterStartedAt"] = *formatDate(param.AfterStartedAt, false)
	}
	if !param.BeforeFinishedAt.IsZero() {
		req.withQueryParams["beforeFinishedAt"] = *formatDate(param.BeforeFinishedAt, false)
	}
	if !param.AfterFinishedAt.IsZero() {
		req.withQueryParams["afterFinishedAt"] = *formatDate(param.AfterFinishedAt, false)
	}
}

func formatDate(date time.Time, _ bool) *string {
	const format = "2006-01-02T15:04:05Z"
	timeParsedToString := date.Format(format)
	return &timeParsedToString
}

func transformStringToMap(primaryKey *string) (options map[string]string) {
	if primaryKey != nil {
		return map[string]string{
			"primaryKey": *primaryKey,
		}
	}
	return nil
}

func transformCsvDocumentsQueryToMap(options *CsvDocumentsQuery) map[string]string {
	var optionsMap map[string]string
	data, _ := json.Marshal(options)
	_ = json.Unmarshal(data, &optionsMap)
	return optionsMap
}

func generateQueryForOptions(options map[string]string) (urlQuery string) {
	q := url.Values{}
	for key, val := range options {
		q.Add(key, val)
	}
	return q.Encode()
}

func sendCsvRecords(ctx context.Context, documentsCsvFunc func(ctx context.Context, recs []byte, op *CsvDocumentsQuery) (resp *TaskInfo, err error), records [][]string, options *CsvDocumentsQuery) (*TaskInfo, error) {
	b := new(bytes.Buffer)
	w := csv.NewWriter(b)
	w.UseCRLF = true

	err := w.WriteAll(records)
	if err != nil {
		return nil, fmt.Errorf("could not write CSV records: %w", err)
	}

	resp, err := documentsCsvFunc(ctx, b.Bytes(), options)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func joinInt64(vals []int64) string {
	if len(vals) == 0 {
		return ""
	}
	result := make([]string, len(vals))
	for i, v := range vals {
		result[i] = strconv.FormatInt(v, 10)
	}
	return joinString(result)
}

func joinString(vals []string) string {
	if len(vals) == 0 {
		return ""
	}

	return strings.Join(vals, ",")
}
