package main

import (
    "encoding/json"
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

type TumblrDetails struct {
    ConsumerKey    string
    ConsumerSecret string
    Token          string
    TokenSecret    string
    TargetBlog     string
    Client         *tumblrclient.Client
}

func (t *TumblrDetails) createClient() {
    t.Client = tumblrclient.NewClientWithToken(
        t.ConsumerKey,
        t.ConsumerSecret,
        t.Token,
        t.TokenSecret,
    )
}

