package utils

import (
	"bytes"
	"fmt"
	"io"
	"net/http"

	"golang.org/x/net/html"
)

func ParsePage(uri string) ([]string, error) {
	var images []string
	client := &http.Client{}
	req, err := http.NewRequest("GET", uri, nil)
	resp, err := client.Do(req)
	if err != nil {
		return images, fmt.Errorf("Unable to fetch image %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return images, fmt.Errorf("Unable to read http response body %v", err)
	}

	doc, err := html.Parse(bytes.NewReader(body))
	if err != nil {
		return images, fmt.Errorf("Unable to parse http response body %v as html", err)
	}
	images = searchForImages(doc)

	return images, nil
}

func searchForImages(n *html.Node) []string {
	var images []string

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		if n.Type == html.ElementNode && n.Data == "body" {
			images = append(images, searchInBody(c)...)
		} else if n.Type == html.ElementNode && n.Data == "head" {
			images = append(images, searchInHeader(c)...)
		} else {
			images = append(images, searchForImages(c)...)
		}

	}
	return images
}

func searchInHeader(n *html.Node) []string {
	var images []string
	if n.Type == html.ElementNode && n.Data == "meta" {
		imageTag := false
		image := ""
		for _, a := range n.Attr {
			if a.Key == "content" {
				image = a.Val
			}
			if (a.Key == "name" || a.Key == "property" || a.Key == "itemprop") &&
				(a.Val == "og:image" || a.Val == "twitter:image" || a.Val == "image") {
				imageTag = true
			}
		}
		if imageTag {
			images = append(images, image)
		}
	} else {
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			images = append(images, searchInHeader(c)...)
		}
	}
	return images
}

func searchInBody(n *html.Node) []string {
	var images []string
	if n.Type == html.ElementNode && n.Data == "img" {
		for _, a := range n.Attr {
			if a.Key == "src" {
				images = append(images, a.Val)
				break
			}
		}
	} else {
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			images = append(images, searchInBody(c)...)
		}
	}
	return images
}
