package meilisearch

import (
	"crypto/tls"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGet_ExperimentalFeatures(t *testing.T) {
	sv := setup(t, "")
	customSv := setup(t, "", WithCustomClientWithTLS(&tls.Config{
		InsecureSkipVerify: true,
	}))

	tests := []struct {
		name   string
		client ServiceManager
	}{
		{
			name:   "TestGetStats",
			client: sv,
		},
		{
			name:   "TestGetStatsWithCustomClient",
			client: customSv,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ef := tt.client.ExperimentalFeatures()
			gotResp, err := ef.Get()
			require.NoError(t, err)
			require.NotNil(t, gotResp, "ExperimentalFeatures.Get() should not return nil value")
		})
	}
}

func TestUpdate_ExperimentalFeatures(t *testing.T) {
	sv := setup(t, "")
	customSv := setup(t, "", WithCustomClientWithTLS(&tls.Config{
		InsecureSkipVerify: true,
	}))

	tests := []struct {
		name   string
		client ServiceManager
	}{
		{
			name:   "TestUpdateStats",
			client: sv,
		},
		{
			name:   "TestUpdateStatsWithCustomClient",
			client: customSv,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ef := tt.client.ExperimentalFeatures()
			ef.SetVectorStore(true)
			ef.SetLogsRoute(true)
			ef.SetMetrics(true)
			ef.SetEditDocumentsByFunction(true)
			ef.SetContainsFilter(true)
			gotResp, err := ef.Update()
			require.NoError(t, err)
			require.Equal(t, true, gotResp.VectorStore, "ExperimentalFeatures.Update() should return vectorStore as true")
			require.Equal(t, true, gotResp.LogsRoute, "ExperimentalFeatures.Update() should return logsRoute as true")
			require.Equal(t, true, gotResp.Metrics, "ExperimentalFeatures.Update() should return metrics as true")
			require.Equal(t, true, gotResp.EditDocumentsByFunction, "ExperimentalFeatures.Update() should return editDocumentsByFunction as true")
			require.Equal(t, true, gotResp.ContainsFilter, "ExperimentalFeatures.Update() should return containsFilter as true")
		})
	}
}
