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
	require.NotEmpty(t, network.Version)
	require.Empty(t, network.Remotes)
	require.Empty(t, network.Self)
	require.Empty(t, network.Leader)
}

func Test_UpdateNetwork(t *testing.T) {
	sv := setup(t, "")
	t.Cleanup(cleanupNetwork(sv))

	experimentalFeatures, err := sv.ExperimentalFeatures().SetNetwork(true).Update()
	require.NoError(t, err)
	require.True(t, experimentalFeatures.Network)

	initialNetwork := &meilisearch.UpdateNetworkRequest{
		Self: meilisearch.String("TEST"),
		Remotes: meilisearch.NewOpt(map[string]meilisearch.Opt[meilisearch.UpdateRemote]{
			"ms-00": meilisearch.NewOpt(meilisearch.UpdateRemote{
				URL:          meilisearch.String("https://example.com"),
				SearchAPIKey: meilisearch.String("TEST"),
			}),
		}),
	}
	network, err := sv.UpdateNetwork(initialNetwork)
	require.NoError(t, err)
	require.Equal(t, "TEST", network.Self)
	require.Empty(t, network.Leader)

	tests := []struct {
		name   string
		update *meilisearch.UpdateNetworkRequest
		want   *meilisearch.Network
	}{
		{
			name: "update self only",
			update: &meilisearch.UpdateNetworkRequest{
				Self: meilisearch.NewOpt("NEW-SELF"),
			},
			want: &meilisearch.Network{
				Self: "NEW-SELF",
				Remotes: map[string]meilisearch.Remote{
					"ms-00": {
						URL:          "https://example.com",
						SearchAPIKey: "TEST",
						WriteAPIKey:  "",
					},
				},
			},
		},
		{
			name: "update remote url and search API key",
			update: &meilisearch.UpdateNetworkRequest{
				Remotes: meilisearch.NewOpt(map[string]meilisearch.Opt[meilisearch.UpdateRemote]{
					"ms-00": meilisearch.NewOpt(meilisearch.UpdateRemote{
						URL:          meilisearch.String("https://updated.com"),
						SearchAPIKey: meilisearch.String("UPDATED_API_KEY"),
					}),
				}),
			},
			want: &meilisearch.Network{
				Self: "NEW-SELF",
				Remotes: map[string]meilisearch.Remote{
					"ms-00": {
						URL:          "https://updated.com",
						SearchAPIKey: "UPDATED_API_KEY",
						WriteAPIKey:  "",
					},
				},
			},
		},
		{
			name: "remove remote",
			update: &meilisearch.UpdateNetworkRequest{
				Remotes: meilisearch.NewOpt(map[string]meilisearch.Opt[meilisearch.UpdateRemote]{
					"ms-00": meilisearch.Null[meilisearch.UpdateRemote](),
				}),
			},
			want: &meilisearch.Network{
				Self:    "NEW-SELF",
				Remotes: map[string]meilisearch.Remote{},
			},
		},
		{
			name: "add new remote with keys",
			update: &meilisearch.UpdateNetworkRequest{
				Remotes: meilisearch.NewOpt(map[string]meilisearch.Opt[meilisearch.UpdateRemote]{
					"ms-01": meilisearch.NewOpt(meilisearch.UpdateRemote{
						URL:          meilisearch.String("https://new-remote.com"),
						SearchAPIKey: meilisearch.NewOpt("NEW-REMOTE-KEY"),
						WriteAPIKey:  meilisearch.String("WRITE-API-KEY"),
					}),
				}),
			},
			want: &meilisearch.Network{
				Self: "NEW-SELF",
				Remotes: map[string]meilisearch.Remote{
					"ms-01": {
						URL:          "https://new-remote.com",
						SearchAPIKey: "NEW-REMOTE-KEY",
						WriteAPIKey:  "WRITE-API-KEY",
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := sv.UpdateNetwork(tt.update)
			require.NoError(t, err)

			require.Equal(t, tt.want.Self, got.Self)
			require.Equal(t, tt.want.Leader, got.Leader)
			require.Equal(t, tt.want.Remotes, got.Remotes)
		})
	}
}
