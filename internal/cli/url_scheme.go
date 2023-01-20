package cli

import (
	"fmt"
	"net/url"
)

type UrlScheme string

const (
	httpsScheme UrlScheme = "https"
	httpScheme  UrlScheme = "http"
)

func addUrlScheme(urlString string, scheme UrlScheme) (string, error) {
	parsedUrl, err := url.Parse(urlString)
	// The url.Parse function may return an error when there's no scheme due to parser ambiguities,
	// so we need to verify if the error is because the url lacks a scheme
	if err != nil && parsedUrl.Scheme != "" {
		return "", err
	}

	if parsedUrl.Scheme == string(scheme) || urlString == "" {
		return urlString, nil
	}

	// Avoids adding a scheme if it already has a different one
	if parsedUrl.Scheme != "" && parsedUrl.Scheme != string(scheme) {
		err := fmt.Errorf("The URL \"%s\" already has a different scheme: \"%s\"", urlString, scheme)
		return "", err
	}

	parsedUrl.Scheme = string(scheme)
	urlWithScheme := parsedUrl.String()
	return urlWithScheme, nil
}
