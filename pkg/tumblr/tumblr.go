package tumblr

import (
	"github.com/tumblr/tumblrclient.go"
)

type TumblrPostResp struct {
	Meta struct {
		Status int    `json:"status"`
		Msg    string `json:"msg"`
	} `json:"meta"`
	Response struct {
		ID          int64  `json:"id"`
		IDString    string `json:"id_string"`
		State       string `json:"state"`
		DisplayText string `json:"display_text"`
	} `json:"response"`
}

type TumblrClient struct {
	Client *tumblrclient.Client
	Blog   string
}

// CreateClient instantiates a new Tumblr client based on the official library.
func CreateClient(ck string, cs string, t string, ts string, blog string) TumblrClient {
	var client TumblrClient
	client.Blog = blog
	client.Client = tumblrclient.NewClientWithToken(ck, cs, t, ts)

	return client
}
