package api

// OnlineGetSafe represents the response of a GET /safe/UUID
type OnlineGetSafe struct {
	// _ref         string `json:"$ref"`
	Description string `json:"description"`
	Name        string `json:"name"`
	Status      string `json:"status"`
	UUIDRef     string `json:"uuid_ref"`
}

type OnlinePostSafe struct {
	Description string `json:"description"`
	Name        string `json:"name"`
}

// OnlineGetPlatform represents the reponse of a GET /platform/UUID
type OnlineGetPlatform struct {
	// _ref       string `json:"$ref"`
	Datacenter string `json:"datacenter"`
	ID         int    `json:"id"`
	Name       string `json:"name"`
}
