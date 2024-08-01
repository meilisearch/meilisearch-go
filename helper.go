package meilisearch

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// This function allows the user to create a Key with an ExpiresAt in time.Time
// and transform the Key structure into a KeyParsed structure to send the time format
// managed by Meilisearch
func convertKeyToParsedKey(key Key) (resp KeyParsed) {
	resp = KeyParsed{Name: key.Name, Description: key.Description, UID: key.UID, Actions: key.Actions, Indexes: key.Indexes}

	// Convert time.Time to *string to feat the exact ISO-8601
	// format of Meilisearch
	if !key.ExpiresAt.IsZero() {
		resp.ExpiresAt = formatDate(key.ExpiresAt, true)
	}
	return resp
}

func IsValidUUID(uuid string) bool {
	r := regexp.MustCompile("^[a-fA-F0-9]{8}-[a-fA-F0-9]{4}-4[a-fA-F0-9]{3}-[8|9|aA|bB][a-fA-F0-9]{3}-[a-fA-F0-9]{12}$")
	return r.MatchString(uuid)
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

func formatDate(date time.Time, key bool) *string {
	const format = "2006-01-02T15:04:05Z"
	timeParsedToString := date.Format(format)
	return &timeParsedToString
}
