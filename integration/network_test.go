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
	require.True(t, network.Remotes.Valid())
	require.True(t, network.Self.Null())
}

func Test_UpdateNetwork(t *testing.T) {
	sv := setup(t, "")
	t.Cleanup(cleanupNetwork(sv))

	experimentalFeatures, err := sv.ExperimentalFeatures().SetNetwork(true).Update()
	require.NoError(t, err)
	require.True(t, experimentalFeatures.Network)

	// Initial network setup
	initialNetwork := &meilisearch.Network{
		Self: meilisearch.String("TEST"),
		Remotes: meilisearch.NewOpt(map[string]meilisearch.Opt[meilisearch.Remote]{
			"ms-00": meilisearch.NewOpt(meilisearch.Remote{
				URL:          meilisearch.String("https://example.com"),
				SearchAPIKey: meilisearch.String("TEST"),
			}),
		}),
	}

	network, err := sv.UpdateNetwork(initialNetwork)
	require.NoError(t, err)
	require.Equal(t, "TEST", network.Self.Value)
	require.Equal(t, "", network.Leader.Value)

	tests := []struct {
		name    string
		update  *meilisearch.Network
		want    *meilisearch.Network
	}{
		{
			name: "update self only",
			update: &meilisearch.Network{
				Self: meilisearch.NewOpt("NEW-SELF"),
			},
			want: &meilisearch.Network{
				Self: meilisearch.NewOpt("NEW-SELF"),
				Remotes: meilisearch.NewOpt(map[string]meilisearch.Opt[meilisearch.Remote]{
					"ms-00": meilisearch.NewOpt(meilisearch.Remote{
						URL:          meilisearch.String("https://example.com"),
						SearchAPIKey: meilisearch.String("TEST"),
						WriteAPIKey:  meilisearch.Null[string](),
					}),
				}),
			},
		},
		{
			name: "update remote url and search API key",
			update: &meilisearch.Network{
				Remotes: meilisearch.NewOpt(map[string]meilisearch.Opt[meilisearch.Remote]{
					"ms-00": meilisearch.NewOpt(meilisearch.Remote{
						URL:          meilisearch.String("https://updated.com"),
						SearchAPIKey: meilisearch.String("UPDATED_API_KEY"),
					}),
				}),
			},
			want: &meilisearch.Network{
				Self: meilisearch.NewOpt("NEW-SELF"),
				Remotes: meilisearch.NewOpt(map[string]meilisearch.Opt[meilisearch.Remote]{
					"ms-00": meilisearch.NewOpt(meilisearch.Remote{
						URL:          meilisearch.String("https://updated.com"),
						SearchAPIKey: meilisearch.String("UPDATED_API_KEY"),
						WriteAPIKey:  meilisearch.Null[string](),
					}),
				}),
			},
		},
		{
			name: "remove remote",
			update: &meilisearch.Network{
				Remotes: meilisearch.NewOpt(map[string]meilisearch.Opt[meilisearch.Remote]{
					"ms-00": meilisearch.Null[meilisearch.Remote](),
				}),
			},
			want: &meilisearch.Network{
				Self:    meilisearch.NewOpt("NEW-SELF"),
				Remotes: meilisearch.NewOpt(map[string]meilisearch.Opt[meilisearch.Remote]{}),
			},
		},
		{
			name: "add new remote with keys",
			update: &meilisearch.Network{
				Remotes: meilisearch.NewOpt(map[string]meilisearch.Opt[meilisearch.Remote]{
					"ms-01": meilisearch.NewOpt(meilisearch.Remote{
						URL:          meilisearch.String("https://new-remote.com"),
						SearchAPIKey: meilisearch.NewOpt("NEW-REMOTE-KEY"),
						WriteAPIKey:  meilisearch.String("WRITE-API-KEY"),
					}),
				}),
			},
			want: &meilisearch.Network{
				Self: meilisearch.NewOpt("NEW-SELF"),
				Remotes: meilisearch.NewOpt(map[string]meilisearch.Opt[meilisearch.Remote]{
					"ms-01": meilisearch.NewOpt(meilisearch.Remote{
						URL:          meilisearch.String("https://new-remote.com"),
						SearchAPIKey: meilisearch.NewOpt("NEW-REMOTE-KEY"),
						WriteAPIKey:  meilisearch.String("WRITE-API-KEY"),
					}),
				}),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			network, err := sv.UpdateNetwork(tt.update)
			require.NoError(t, err)
			require.Equal(t, tt.want.Self.Value, network.Self.Value)
			require.Equal(t, tt.want.Leader.Value, network.Leader.Value)
			require.Equal(t, tt.want.Remotes.Value, network.Remotes.Value)
		})
	}
}