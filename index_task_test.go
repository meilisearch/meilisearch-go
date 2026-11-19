package meilisearch

import (
	"context"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

// newTasksServer returns a server answering GET /tasks with an empty result,
// recording the query parameters of the last request.
func newTasksServer(t *testing.T) (*httptest.Server, *url.Values) {
	t.Helper()
	got := new(url.Values)
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodGet, r.Method)
		require.Equal(t, "/tasks", r.URL.Path)
		*got = r.URL.Query()
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"results":[],"limit":20,"from":null,"next":null,"total":0}`))
	}))
	t.Cleanup(srv.Close)
	return srv, got
}

func TestIndex_GetTasks_QueryParams(t *testing.T) {
	t.Parallel()

	afterEnqueuedAt := time.Date(2026, 11, 1, 10, 30, 0, 0, time.UTC)

	tests := []struct {
		name  string
		query *TasksQuery
		want  url.Values
	}{
		{
			name:  "nil query is scoped to the index",
			query: nil,
			want:  url.Values{"indexUids": {"books"}},
		},
		{
			name:  "empty query is scoped to the index",
			query: &TasksQuery{},
			want:  url.Values{"indexUids": {"books"}},
		},
		{
			name:  "index uids from the query are kept",
			query: &TasksQuery{IndexUIDS: []string{"movies"}},
			want:  url.Values{"indexUids": {"movies,books"}},
		},
		{
			name: "every filter reaches the server",
			query: &TasksQuery{
				UIDS:            []int64{42, 43},
				Limit:           10,
				From:            5,
				Statuses:        []TaskStatus{TaskStatusSucceeded, TaskStatusFailed},
				Types:           []TaskType{TaskTypeDocumentAdditionOrUpdate},
				CanceledBy:      []int64{7},
				AfterEnqueuedAt: afterEnqueuedAt,
				Reverse:         true,
			},
			want: url.Values{
				"indexUids":       {"books"},
				"uids":            {"42,43"},
				"limit":           {"10"},
				"from":            {"5"},
				"statuses":        {"succeeded,failed"},
				"types":           {"documentAdditionOrUpdate"},
				"canceledBy":      {"7"},
				"afterEnqueuedAt": {"2026-11-01T10:30:00Z"},
				"reverse":         {"true"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			srv, got := newTasksServer(t)

			_, err := New(srv.URL).Index("books").GetTasks(context.Background(), tt.query)
			require.NoError(t, err)
			require.Equal(t, tt.want, *got)
		})
	}
}

func TestIndex_GetTasks_DoesNotMutateQuery(t *testing.T) {
	t.Parallel()
	srv, got := newTasksServer(t)

	// Spare capacity would let a plain append write into the caller's array.
	indexUIDs := make([]string, 1, 4)
	indexUIDs[0] = "movies"
	query := &TasksQuery{IndexUIDS: indexUIDs}

	idx := New(srv.URL).Index("books")
	for n := 0; n < 2; n++ {
		_, err := idx.GetTasks(context.Background(), query)
		require.NoError(t, err)
		require.Equal(t, "movies,books", got.Get("indexUids"))
	}

	require.Equal(t, []string{"movies"}, query.IndexUIDS)
	require.Equal(t, []string{"movies", "", "", ""}, indexUIDs[:cap(indexUIDs)])
}

func TestClient_GetTasks_Reverse(t *testing.T) {
	t.Parallel()
	srv, got := newTasksServer(t)

	_, err := New(srv.URL).GetTasks(context.Background(), &TasksQuery{Reverse: true})
	require.NoError(t, err)
	require.Equal(t, url.Values{"reverse": {"true"}}, *got)
}
