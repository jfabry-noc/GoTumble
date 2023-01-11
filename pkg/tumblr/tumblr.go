package tumblr

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/url"

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

// VerifyBlog verifies that the blogID is valid for the current user.
func (c *TumblrClient) VerifyBlog(blogID string) bool {
	result := false
	userData, err := c.Client.GetUser()
	if err != nil {
		return result
	}

	for _, blog := range userData.Blogs {
		if blogID == blog.Name {
			return true
		} else if blogID == blog.Url {
			return true
		} else if cleanUrl(blogID) == cleanUrl(blog.Url) {
			return true
		}
	}

	return false
}

// cleanUrl returns an FQDN with nothing before or after.
func cleanUrl(potentialUrl string) string {
	processed, err := url.Parse(potentialUrl)

	if err != nil {
		// Not a URL, return the raw data.
		return potentialUrl
	}

	return processed.Host
}

func (c *TumblrClient) addPost(formData url.Values, format string) error {
	postPath := fmt.Sprintf("blog/%v/post", c.Blog)
	resp, err := c.Client.PostWithParams(postPath, formData)
	if err != nil {
		return err
	}

	var postResponse TumblrPostResp
	postErr := json.Unmarshal(resp.GetBody(), &postResponse)
	if postErr != nil {
		return err
	}

	if postResponse.Meta.Status != 201 {
		errorMessage := fmt.Sprintf("Tumblr response code was: %v", postResponse.Meta.Status)
		return errors.New(errorMessage)
	}
	return nil
}

// AddTextPost adds a new post to the Tumblr account.
func (c *TumblrClient) AddTextPost(content string, format string, tags string) error {
	formData := url.Values{}
	formData.Add("type", "text")
	formData.Add("state", "published")
	formData.Add("body", content)
	formData.Add("format", format)

	if tags != "" {
		formData.Add("tags", tags)
	}

	return c.addPost(formData, format)
}

func (c *TumblrClient) AddLinkPost(description string, link string, format string, tags string) error {
	formData := url.Values{}
	formData.Add("type", "link")
	formData.Add("state", "published")
	formData.Add("url", link)

	if description != "" {
		formData.Add("description", description)
	}

	if tags != "" {
		formData.Add("tags", tags)
	}

	return c.addPost(formData, format)
}
