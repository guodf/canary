package netutil

import (
	"crypto/tls"
	"net/http"
)

var tr = &http.Transport{
	TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
}

func IgnoreTls(client *http.Client) *http.Client {
	client.Transport = tr
	return client
}
