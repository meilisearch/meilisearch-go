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
	tests := []struct {
		name              string
		initializeNetwork func(t *testing.T, sv meilisearch.ServiceManager)
		update            *meilisearch.UpdateNetworkRequest
		want              *meilisearch.Network
	}{
		{
			name: "SetNetworkLeader",
			update: &meilisearch.UpdateNetworkRequest{
				Self:   meilisearch.String("ms-00"),
				Leader: meilisearch.String("ms-00"),
				Remotes: meilisearch.NewOpt(map[string]meilisearch.Opt[meilisearch.UpdateRemote]{
					"ms-00": meilisearch.NewOpt(meilisearch.UpdateRemote{
						URL:          meilisearch.String(getDefaultHost()),
						SearchAPIKey: meilisearch.NewOpt(masterKey),
						WriteAPIKey:  meilisearch.String(masterKey),
					}),
				}),
			},
			want: &meilisearch.Network{
				Self:   "ms-00",
				Leader: "ms-00",
				Remotes: map[string]meilisearch.Remote{
					"ms-00": {
						URL:          getDefaultHost(),
						SearchAPIKey: masterKey,
						WriteAPIKey:  masterKey,
					},
				},
			},
		},
		{
			name: "RemoveNetworkLeader",
			initializeNetwork: func(t *testing.T, sv meilisearch.ServiceManager) {
				value, err := sv.UpdateNetwork(&meilisearch.UpdateNetworkRequest{
					Self:   meilisearch.String("ms-00"),
					Leader: meilisearch.String("ms-00"),
					Remotes: meilisearch.NewOpt(map[string]meilisearch.Opt[meilisearch.UpdateRemote]{
						"ms-00": meilisearch.NewOpt(meilisearch.UpdateRemote{
							URL:          meilisearch.String(getDefaultHost()),
							SearchAPIKey: meilisearch.NewOpt(masterKey),
							WriteAPIKey:  meilisearch.String(masterKey),
						}),
					}),
				})
				require.NoError(t, err)
				task, ok := value.(*meilisearch.Task)
				if !ok {
					require.Fail(t, "expected task to be returned but got %T", value)
				}
				testWaitForTask(t, sv, task)
			},
			update: &meilisearch.UpdateNetworkRequest{
				Self:    meilisearch.String("new-self"),
				Leader:  meilisearch.Null[string](),
				Remotes: meilisearch.Null[map[string]meilisearch.Opt[meilisearch.UpdateRemote]](),
			},
			want: &meilisearch.Network{
				Self:    "new-self",
				Leader:  "",
				Remotes: map[string]meilisearch.Remote{},
			},
		},
		{
			name: "UpdateRemoteAndSelf",
			update: &meilisearch.UpdateNetworkRequest{
				Self: meilisearch.String("SELF"),
				Remotes: meilisearch.NewOpt(map[string]meilisearch.Opt[meilisearch.UpdateRemote]{
					"ms-00": meilisearch.NewOpt(meilisearch.UpdateRemote{
						URL:          meilisearch.String("https://updated.com"),
						SearchAPIKey: meilisearch.String("UPDATED_API_KEY"),
						WriteAPIKey:  meilisearch.Null[string](),
					}),
				}),
			},
			want: &meilisearch.Network{
				Self: "SELF",
				Remotes: map[string]meilisearch.Remote{
					"ms-00": {
						URL:          "https://updated.com",
						SearchAPIKey: "UPDATED_API_KEY",
						WriteAPIKey:  "",
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sv := setup(t, "")
			t.Cleanup(cleanupNetwork(sv))

			experimentalFeatures, err := sv.ExperimentalFeatures().SetNetwork(true).Update()
			require.NoError(t, err)
			require.True(t, experimentalFeatures.Network)

			if tt.initializeNetwork != nil {
				tt.initializeNetwork(t, sv)
			}

			value, err := sv.UpdateNetwork(tt.update)
			require.NoError(t, err)

			if updatedNetwork, ok := value.(*meilisearch.Network); ok {
				network, err := sv.GetNetwork()
				require.NoError(t, err)
				require.Equal(t, tt.want.Self, updatedNetwork.Self)
				require.Equal(t, tt.want.Leader, updatedNetwork.Leader)
				require.Equal(t, tt.want.Remotes, updatedNetwork.Remotes)
				require.Equal(t, updatedNetwork, network)
			} else if task, ok := value.(*meilisearch.Task); ok {
				testWaitForTask(t, sv, task)
				network, err := sv.GetNetwork()
				require.NoError(t, err)

				require.Equal(t, tt.want.Self, network.Self)
				require.Equal(t, tt.want.Leader, network.Leader)
				require.Equal(t, tt.want.Remotes, network.Remotes)
			} else {
				require.Fail(t, "unexpected type returned: %T", value)
			}
		})
	}
}
