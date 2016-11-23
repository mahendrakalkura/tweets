package main

import (
	"net/http"
	"net/url"
	"time"
)

func get_http_client(settings *Settings, with_proxy bool) *http.Client {
	timeout := time.Duration(30 * time.Second)

	proxy := get_proxy(settings.Proxies.Hostname, settings.Proxies.Ports)
	proxy_url, err := url.Parse(proxy)
	if err != nil {
		panic(err)
	}

	client := &http.Client{}
	client.Timeout = timeout
	if with_proxy {
		client.Transport = &http.Transport{Proxy: http.ProxyURL(proxy_url)}
	}

	return client
}
