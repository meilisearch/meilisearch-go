package integration

import (
	"testing"

	"github.com/meilisearch/meilisearch-go"
	"github.com/stretchr/testify/require"
)

func Test_GetNetwork(t *testing.T) {
	sv := setup(t, "")
	t.Cleanup(cleanupNetwork(sv))

	experimentalFeatures, err := sv.ExperimentalFeatures().SetNetwork(true).Update()
	require.NoError(t, err)
	require.True(t, experimentalFeatures.Network)

	network, err := sv.GetNetwork()
	require.NoError(t, err)
	require.Nil(t, network.Self)
	require.Len(t, network.Remotes, 0)
}

func Test_UpdateNetwork(t *testing.T) {
	sv := setup(t, "")
	t.Cleanup(cleanupNetwork(sv))

	experimentalFeatures, err := sv.ExperimentalFeatures().SetNetwork(true).Update()
	require.NoError(t, err)
	require.True(t, experimentalFeatures.Network)

	tests := []struct {
		name        string
		input       *meilisearch.Network
		wantSelf    string
		wantRemotes map[string]*meilisearch.Remote
	}{
		{
			name: "set initial network",
			input: &meilisearch.Network{
				Self: meilisearch.String("TEST"),
				Remotes: map[string]*meilisearch.Remote{
					"ms-00": {
						URL:          meilisearch.String("https://example.com"),
						SearchApiKey: meilisearch.String("TEST"),
					},
				},
			},
			wantSelf: "TEST",
			wantRemotes: map[string]*meilisearch.Remote{
				"ms-00": {
					URL:          meilisearch.String("https://example.com"),
					SearchApiKey: meilisearch.String("TEST"),
				},
			},
		},
		{
			name: "update self only",
			input: &meilisearch.Network{
				Self: meilisearch.String("NEW-SELF"),
				Remotes: map[string]*meilisearch.Remote{
					"ms-00": {
						URL:          meilisearch.String("https://example.com"),
						SearchApiKey: meilisearch.String("TEST"),
					},
				},
			},
			wantSelf: "NEW-SELF",
			wantRemotes: map[string]*meilisearch.Remote{
				"ms-00": {
					URL:          meilisearch.String("https://example.com"),
					SearchApiKey: meilisearch.String("TEST"),
				},
			},
		},
		{
			name: "update remote url only",
			input: &meilisearch.Network{
				Self: meilisearch.String("NEW-SELF"),
				Remotes: map[string]*meilisearch.Remote{
					"ms-00": {
						URL:          meilisearch.String("https://updated.com"),
						SearchApiKey: meilisearch.String("TEST"),
					},
				},
			},
			wantSelf: "NEW-SELF",
			wantRemotes: map[string]*meilisearch.Remote{
				"ms-00": {
					URL:          meilisearch.String("https://updated.com"),
					SearchApiKey: meilisearch.String("TEST"),
				},
			},
		},
		{
			name: "update remote searchApiKey only",
			input: &meilisearch.Network{
				Self: meilisearch.String("NEW-SELF"),
				Remotes: map[string]*meilisearch.Remote{
					"ms-00": {
						URL:          meilisearch.String("https://updated.com"),
						SearchApiKey: meilisearch.String("NEW-KEY"),
					},
				},
			},
			wantSelf: "NEW-SELF",
			wantRemotes: map[string]*meilisearch.Remote{
				"ms-00": {
					URL:          meilisearch.String("https://updated.com"),
					SearchApiKey: meilisearch.String("NEW-KEY"),
				},
			},
		},
		{
			name: "remove remote",
			input: &meilisearch.Network{
				Self: meilisearch.String("NEW-SELF"),
				Remotes: map[string]*meilisearch.Remote{
					"ms-00": nil, // remove by setting nil
				},
			},
			wantSelf:    "NEW-SELF",
			wantRemotes: map[string]*meilisearch.Remote{},
		},
		{
			name: "add new remote",
			input: &meilisearch.Network{
				Self: meilisearch.String("NEW-SELF"),
				Remotes: map[string]*meilisearch.Remote{
					"ms-01": {
						URL:          meilisearch.String("https://new-remote.com"),
						SearchApiKey: meilisearch.String("NEW-REMOTE-KEY"),
					},
				},
			},
			wantSelf: "NEW-SELF",
			wantRemotes: map[string]*meilisearch.Remote{
				"ms-01": {
					URL:          meilisearch.String("https://new-remote.com"),
					SearchApiKey: meilisearch.String("NEW-REMOTE-KEY"),
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			network, err := sv.UpdateNetwork(tt.input)
			require.NoError(t, err)

			require.Equal(t, tt.wantSelf, *network.Self)
			require.Equal(t, tt.wantRemotes, network.Remotes)
		})
	}
}
