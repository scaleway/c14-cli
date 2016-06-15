package api

// OnlineGetSafe represents the response of a GET /safe/UUID
type OnlineGetSafe struct {
	// _ref         string `json:"$ref"`
	Description string `json:"description"`
	Name        string `json:"name"`
	Status      string `json:"status"`
	UUIDRef     string `json:"uuid_ref"`
}
