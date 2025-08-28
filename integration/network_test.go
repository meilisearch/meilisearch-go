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

	tests := []struct {
		name         string
		input        *meilisearch.Network
		wantSelf     string
		wantSharding bool
		wantRemotes  map[string]meilisearch.Opt[meilisearch.Remote]
	}{
		{
			name: "set initial network",
			input: &meilisearch.Network{
				Self: meilisearch.String("TEST"),
				Remotes: meilisearch.NewOpt(map[string]meilisearch.Opt[meilisearch.Remote]{
					"ms-00": meilisearch.NewOpt(meilisearch.Remote{
						URL:          meilisearch.String("https://example.com"),
						SearchAPIKey: meilisearch.String("TEST"),
					}),
				}),
			},
			wantSelf:     "TEST",
			wantSharding: false,
			wantRemotes: map[string]meilisearch.Opt[meilisearch.Remote]{
				"ms-00": meilisearch.NewOpt(meilisearch.Remote{
					URL:          meilisearch.String("https://example.com"),
					SearchAPIKey: meilisearch.String("TEST"),
					WriteAPIKey:  meilisearch.Null[string](),
				}),
			},
		},
		{
			name: "update self only (remotes omitted)",
			input: &meilisearch.Network{
				Self: meilisearch.NewOpt("NEW-SELF"),
			},
			wantSelf:     "NEW-SELF",
			wantSharding: false,
			wantRemotes: map[string]meilisearch.Opt[meilisearch.Remote]{
				"ms-00": meilisearch.NewOpt(meilisearch.Remote{
					URL:          meilisearch.String("https://example.com"),
					SearchAPIKey: meilisearch.String("TEST"),
					WriteAPIKey:  meilisearch.Null[string](),
				}),
			},
		},
		{
			name: "update remote url only",
			input: &meilisearch.Network{
				Remotes: meilisearch.NewOpt(map[string]meilisearch.Opt[meilisearch.Remote]{
					"ms-00": meilisearch.NewOpt(meilisearch.Remote{
						URL:         meilisearch.String("https://updated.com"),
						WriteAPIKey: meilisearch.Null[string](),
					}),
				}),
			},
			wantSelf:     "NEW-SELF",
			wantSharding: false,
			wantRemotes: map[string]meilisearch.Opt[meilisearch.Remote]{
				"ms-00": meilisearch.NewOpt(meilisearch.Remote{
					URL:          meilisearch.String("https://updated.com"),
					SearchAPIKey: meilisearch.String("TEST"), // unchanged
					WriteAPIKey:  meilisearch.Null[string](),
				}),
			},
		},
		{
			name: "update remote searchApiKey only",
			input: &meilisearch.Network{
				Remotes: meilisearch.NewOpt(map[string]meilisearch.Opt[meilisearch.Remote]{
					"ms-00": meilisearch.NewOpt(meilisearch.Remote{
						SearchAPIKey: meilisearch.String("UPDATED_API_KEY"),
						WriteAPIKey:  meilisearch.Null[string](),
					}),
				}),
			},
			wantSelf:     "NEW-SELF",
			wantSharding: false,
			wantRemotes: map[string]meilisearch.Opt[meilisearch.Remote]{
				"ms-00": meilisearch.NewOpt(meilisearch.Remote{
					URL:          meilisearch.String("https://updated.com"),
					SearchAPIKey: meilisearch.String("UPDATED_API_KEY"), // unchanged
					WriteAPIKey:  meilisearch.Null[string](),
				}),
			},
		},
		{
			name: "remove searchApiKey",
			input: &meilisearch.Network{
				Remotes: meilisearch.NewOpt(map[string]meilisearch.Opt[meilisearch.Remote]{
					"ms-00": meilisearch.NewOpt(meilisearch.Remote{
						SearchAPIKey: meilisearch.Null[string](),
						WriteAPIKey:  meilisearch.Null[string](),
					}),
				}),
			},
			wantSelf:     "NEW-SELF",
			wantSharding: false,
			wantRemotes: map[string]meilisearch.Opt[meilisearch.Remote]{
				"ms-00": meilisearch.NewOpt(meilisearch.Remote{
					URL:          meilisearch.String("https://updated.com"),
					SearchAPIKey: meilisearch.Null[string](),
					WriteAPIKey:  meilisearch.Null[string](),
				}),
			},
		},
		{
			name: "remove remote",
			input: &meilisearch.Network{
				Remotes: meilisearch.NewOpt(map[string]meilisearch.Opt[meilisearch.Remote]{
					"ms-00": meilisearch.Null[meilisearch.Remote](),
				}),
			},
			wantSharding: false,
			wantSelf:     "NEW-SELF",
			wantRemotes:  map[string]meilisearch.Opt[meilisearch.Remote]{},
		},
		{
			name: "add new remote",
			input: &meilisearch.Network{
				Remotes: meilisearch.NewOpt(map[string]meilisearch.Opt[meilisearch.Remote]{
					"ms-01": meilisearch.NewOpt(meilisearch.Remote{
						URL:          meilisearch.String("https://new-remote.com"),
						SearchAPIKey: meilisearch.NewOpt("NEW-REMOTE-KEY"),
					}),
				}),
			},
			wantSelf:     "NEW-SELF",
			wantSharding: false,
			wantRemotes: map[string]meilisearch.Opt[meilisearch.Remote]{
				"ms-01": meilisearch.NewOpt(meilisearch.Remote{
					URL:          meilisearch.String("https://new-remote.com"),
					SearchAPIKey: meilisearch.NewOpt("NEW-REMOTE-KEY"),
					WriteAPIKey:  meilisearch.Null[string](),
				}),
			},
		},
		{
			name: "enable sharding",
			input: &meilisearch.Network{
				Sharding: meilisearch.Bool(true),
			},
			wantSelf:     "NEW-SELF",
			wantSharding: true,
			wantRemotes: map[string]meilisearch.Opt[meilisearch.Remote]{
				"ms-01": meilisearch.NewOpt(meilisearch.Remote{
					URL:          meilisearch.String("https://new-remote.com"),
					SearchAPIKey: meilisearch.NewOpt("NEW-REMOTE-KEY"),
					WriteAPIKey:  meilisearch.Null[string](),
				}),
			},
		},
		{
			name: "disable sharding",
			input: &meilisearch.Network{
				Sharding: meilisearch.Bool(false),
			},
			wantSharding: false,
			wantSelf:     "NEW-SELF",
			wantRemotes: map[string]meilisearch.Opt[meilisearch.Remote]{
				"ms-01": meilisearch.NewOpt(meilisearch.Remote{
					URL:          meilisearch.String("https://new-remote.com"),
					SearchAPIKey: meilisearch.NewOpt("NEW-REMOTE-KEY"),
					WriteAPIKey:  meilisearch.Null[string](),
				}),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			network, err := sv.UpdateNetwork(tt.input)
			require.NoError(t, err)

			require.Equal(t, tt.wantSelf, network.Self.Value)
			require.Equal(t, tt.wantRemotes, network.Remotes.Value)
			require.Equal(t, tt.wantSharding, network.Sharding.Value)
		})
	}
}
