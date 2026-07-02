package meilisearch

import "encoding/json"

// Network represents the Meilisearch network configuration.
type Network struct {
	Self    string            `json:"self,omitempty"`
	Leader  string            `json:"leader,omitempty"`
	Remotes map[string]Remote `json:"remotes,omitempty"`
	Shards  map[string]Shard  `json:"shards,omitempty"`
	Version string            `json:"version,omitempty"`
}

// Remote describes a single remote Meilisearch node.
type Remote struct {
	URL          string `json:"url"`
	SearchAPIKey string `json:"searchApiKey"`
	WriteAPIKey  string `json:"writeApiKey"`
}

// Shard represents a shard in the network
type Shard struct {
	Remotes []string `json:"remotes"`
}

// UpdateNetworkRequest represents the Meilisearch network configuration update without leader.
// Each field is wrapped in an Opt so it can be explicitly included,
// set to JSON null, or omitted entirely.
type UpdateNetworkRequest struct {
	Self    Opt[string]                            `json:"self,omitempty"`
	Leader  Opt[string]                            `json:"leader,omitempty"`
	Remotes Opt[map[string]Opt[UpdateRemote]]      `json:"remotes,omitempty"`
	Shards  Opt[map[string]Opt[UpdateRemoteShard]] `json:"shards,omitempty"` // FIXED
	Version Opt[string]                            `json:"version,omitempty"`
}

func (n UpdateNetworkRequest) MarshalJSON() ([]byte, error) {
	m := make(map[string]any)

	if n.Self.Valid() {
		m["self"] = n.Self.Value
	} else if n.Self.Null() {
		m["self"] = nil
	}

	if n.Remotes.Valid() {
		m["remotes"] = n.Remotes.Value
	} else if n.Remotes.Null() {
		m["remotes"] = nil
	}

	if n.Leader.Valid() {
		m["leader"] = n.Leader.Value
	} else if n.Leader.Null() {
		m["leader"] = nil
	}

	if n.Shards.Valid() {
		m["shards"] = n.Shards.Value
	} else if n.Shards.Null() {
		m["shards"] = nil
	}

	if n.Version.Valid() {
		m["version"] = n.Version.Value
	} else if n.Version.Null() {
		m["version"] = nil
	}

	return json.Marshal(m)
}

// UpdateRemote describes a single remote Meilisearch node.
// Each field is wrapped in an Opt so it can be explicitly included,
// set to JSON null, or omitted entirely.
type UpdateRemote struct {
	URL          Opt[string] `json:"url"`
	SearchAPIKey Opt[string] `json:"searchApiKey,omitempty"`
	WriteAPIKey  Opt[string] `json:"writeApiKey,omitempty"`
}

func (r UpdateRemote) MarshalJSON() ([]byte, error) {
	m := make(map[string]any)

	if r.URL.Valid() {
		m["url"] = r.URL.Value
	} else if r.URL.Null() {
		m["url"] = nil
	}

	if r.SearchAPIKey.Valid() {
		m["searchApiKey"] = r.SearchAPIKey.Value
	} else if r.SearchAPIKey.Null() {
		m["searchApiKey"] = nil
	}

	if r.WriteAPIKey.Valid() {
		m["writeApiKey"] = r.WriteAPIKey.Value
	} else if r.WriteAPIKey.Null() {
		m["writeApiKey"] = nil
	}

	return json.Marshal(m)
}

// UpdateRemoteShard represents a shard update request
type UpdateRemoteShard struct {
	Remotes       Opt[[]string] `json:"remotes,omitempty"`
	AddRemotes    Opt[[]string] `json:"addRemotes,omitempty"`
	RemoveRemotes Opt[[]string] `json:"removeRemotes,omitempty"`
}

func (u UpdateRemoteShard) MarshalJSON() ([]byte, error) {
	m := make(map[string]any)

	if u.Remotes.Valid() {
		m["remotes"] = u.Remotes.Value
	} else if u.Remotes.Null() {
		m["remotes"] = nil
	}

	if u.AddRemotes.Valid() {
		m["addRemotes"] = u.AddRemotes.Value
	} else if u.AddRemotes.Null() {
		m["addRemotes"] = nil
	}

	if u.RemoveRemotes.Valid() {
		m["removeRemotes"] = u.RemoveRemotes.Value
	} else if u.RemoveRemotes.Null() {
		m["removeRemotes"] = nil
	}

	return json.Marshal(m)
}
