package integration

import (
	"fmt"
	"testing"

	"github.com/meilisearch/meilisearch-go"
	"github.com/stretchr/testify/require"
)

func Test_AddWebhook(t *testing.T) {
	sv := setup(t, "")
	t.Cleanup(cleanupWebhook(sv))

	webhook := &meilisearch.AddWebhookRequest{
		URL:     "http://example.com",
		Headers: map[string]string{"FOO": "BAR", "BAR": ""},
	}
	result, err := sv.AddWebhook(webhook)
	require.NoError(t, err)
	require.Equal(t, webhook.URL, result.URL)
	require.Equal(t, webhook.Headers, result.Headers)
	require.True(t, result.IsEditable)
	require.NotZero(t, result.UUID)
}

func Test_GetWebhook(t *testing.T) {
	sv := setup(t, "")
	t.Cleanup(cleanupWebhook(sv))

	wb, err := sv.AddWebhook(&meilisearch.AddWebhookRequest{
		URL:     "http://example.com",
		Headers: map[string]string{"FOO": "BAR"},
	})
	require.NoError(t, err)
	got, err := sv.GetWebhook(wb.UUID)
	require.NoError(t, err)
	require.Equal(t, wb, got)
}

func Test_UpdateWebhook(t *testing.T) {
	sv := setup(t, "")
	t.Cleanup(cleanupWebhook(sv))

	wb, err := sv.AddWebhook(&meilisearch.AddWebhookRequest{
		URL:     "http://example.com",
		Headers: map[string]string{"FOO": "BAR"},
	})
	require.NoError(t, err)

	updatedWb, err := sv.UpdateWebhook(wb.UUID, &meilisearch.UpdateWebhookRequest{
		URL:     "http://update.com",
		Headers: map[string]string{"FOO": "UPDATED", "BAR": ""},
	})
	require.NoError(t, err)

	got, err := sv.GetWebhook(wb.UUID)
	require.NoError(t, err)
	require.Equal(t, updatedWb, got)
}

func Test_ListWebhooks(t *testing.T) {
	sv := setup(t, "")
	t.Cleanup(cleanupWebhook(sv))

	n := 5
	for i := 0; i < n; i++ {
		_, err := sv.AddWebhook(&meilisearch.AddWebhookRequest{
			URL:     fmt.Sprintf("http://example_%d.com", i),
			Headers: map[string]string{"FOO": "BAR"},
		})
		require.NoError(t, err)
	}
	result, err := sv.ListWebhooks()
	require.NoError(t, err)
	require.Len(t, result.Result, n)
}

func Test_DeleteWebhook(t *testing.T) {
	sv := setup(t, "")
	t.Cleanup(cleanupWebhook(sv))

	wb, err := sv.AddWebhook(&meilisearch.AddWebhookRequest{
		URL:     "http://example.com",
		Headers: map[string]string{"FOO": "BAR"},
	})
	require.NoError(t, err)

	err = sv.DeleteWebhook(wb.UUID)
	require.NoError(t, err)
}
