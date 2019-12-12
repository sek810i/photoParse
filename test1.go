package main

import (
	"net/url"
	"net/http"
	"io/ioutil"
	"fmt"
	"encoding/json"
	"golang.org/x/net/html"
	"net/http/cookiejar"
	"io"
	"errors"
	"os"
)

const (
	version = "5.103"
	clientId = "7238281"
	apiURL  = "https://api.vk.com/method/"
	authURL = "https://oauth.vk.com/authorize?" +
		"client_id=%s" +
		"&scope=%s" +
		"&redirect_uri=https://oauth.vk.com/blank.html" +
		"&display=wap" +
		"&v=%s" +
		"&response_type=token"
)

type vk struct {
	accessToken string
	version     string
}

func main()  {

}

func Auth(login, password string) vk {
	urlPath := fmt.Sprintf(authURL, clientId, "photos", version)
	jar, _ := cookiejar.New(nil)
	client := &http.Client {
		Jar: jar,
	}

	resp, err := client.Get(urlPath)
	if err != nil {
		return vk{}
	}
	defer resp.Body.Close()

	args, u := parseForm(resp.Body)

	args.Add("email", login)
	args.Add("pass", password)

	resp, err = client.PostForm(u, args)
	if err != nil {
		return vk{}
	}

	if resp.Request.URL.Path != "/blank.html" {
		args, u := parseForm(resp.Body)
		resp, err := client.PostForm(u, args)
		if err != nil {
			return vk{}
		}
		defer resp.Body.Close()

		if resp.Request.URL.Path != "/blank.html" {
			return vk{}
		}
	}

	urlArgs, err := url.ParseQuery(resp.Request.URL.Fragment)
	if err != nil {
		return vk{}
	}

	return vk{version:version, accessToken:urlArgs["access_token"][0]}
}

func check(err error)  {
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func parseForm(body io.ReadCloser) (url.Values, string) {
	tokenizer := html.NewTokenizer(body)

	u := ""
	formData := map[string]string{}

	end := false
	for !end {
		tag := tokenizer.Next()

		switch tag {
		case html.ErrorToken:
			end = true
		case html.StartTagToken:
			switch token := tokenizer.Token(); token.Data {
			case "form":
				for _, attr := range token.Attr {
					if attr.Key == "action" {
						u = attr.Val
					}
				}
			case "input":
				if token.Attr[1].Val == "_origin" {
					formData["_origin"] = token.Attr[2].Val
				}
				if token.Attr[1].Val == "to" {
					formData["to"] = token.Attr[2].Val
				}
			}
		case html.SelfClosingTagToken:
			switch token := tokenizer.Token(); token.Data {
			case "input":
				if token.Attr[1].Val == "ip_h" {
					formData["ip_h"] = token.Attr[2].Val
				}
				if token.Attr[1].Val == "lg_h" {
					formData["lg_h"] = token.Attr[2].Val
				}
			}
		}
	}

	args := url.Values{}

	for key, val := range formData {
		args.Add(key, val)
	}

	return args, u
}
