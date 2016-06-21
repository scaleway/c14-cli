package api

/*
 * Get
 */

// OnlineGetSafe represents the response of a GET /safe/UUID
type OnlineGetSafe struct {
	// _ref         string `json:"$ref"`
	Description string `json:"description"`
	Name        string `json:"name"`
	Status      string `json:"status"`
	UUIDRef     string `json:"uuid_ref"`
}

// OnlineGetPlatform represents the reponse of a GET /platform/UUID
type OnlineGetPlatform struct {
	// _ref       string `json:"$ref"`
	Datacenter string `json:"datacenter"`
	ID         int    `json:"id"`
	Name       string `json:"name"`
}

type OnlineGetSSHKey struct {
	// _ref string `json:"$ref"`
	Desc        string `json:"description"`
	Fingerprint string `json:"fingerprint"`
	UUIDRef     string `json:"uuid_ref"`
}

/*
 * Post
 */

type OnlinePostSafe struct {
	Description string `json:"description"`
	Name        string `json:"name"`
}

type OnlinePostArchive struct {
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Parity      string   `json:"parity,omitempty"`
	Protocols   []string `json:"protocols"`
	SSHKeys     []string `json:"ssh_keys"`
	Platforms   []string `json:"platforms"`
	Days        int      `json:"days"`
}
