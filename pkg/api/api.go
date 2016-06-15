package api

import "net/http"

var (
	APIUrl = "https://api.online.net/api/v1"
)

type OnlineAPI struct {
	client http.Client // use getClient() to copy the structure and avoid race conditions
}

func NewC14API(client http.Client) (api *OnlineAPI) {
	api = &OnlineAPI{
		client: client,
	}
	return
}
