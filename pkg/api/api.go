package api

import "net/http"

var (
	// APIUrl represents Online's endpoint
	APIUrl = "https://api.online.net/api/v1"
)

// OnlineAPI is used to communicate with Online API
type OnlineAPI struct {
	client http.Client // use getClient() to copy the structure and avoid race conditions
}

// NewC14API returns a new API
func NewC14API(client http.Client) (api *OnlineAPI) {
	api = &OnlineAPI{
		client: client,
	}
	return
}
