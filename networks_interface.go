package meilisearch

import "context"

type NetworksManager interface {
	NetworksReader

	// UpdateNetwork updates the network object.
	// Updates are partial; only the provided fields are updated.
	UpdateNetwork(params *Network) (*Network, error)

	// UpdateNetworkWithContext updates the network object with a context.
	// Updates are partial; only the provided fields are updated.
	UpdateNetworkWithContext(ctx context.Context, params *Network) (*Network, error)
}

type NetworksReader interface {
	// GetNetwork gets the current value of the instance’s network object.
	GetNetwork() (*Network, error)

	// GetNetworkWithContext gets the current value of the instance’s network object with a context.
	GetNetworkWithContext(ctx context.Context) (*Network, error)
}
