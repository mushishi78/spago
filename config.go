package main

type Config struct {
	Port                 int      `json:"port,omitempty"`
	ExcludedPaths        []string `json:"excludedPaths,omitempty"`
	StaticFileExtensions []string `json:"staticFileExtensions,omitempty"`
	ReverseProxyURL      string   `json:"reverseProxyUrl,omitempty"`
	ReverseProxyRoute    string   `json:"reverseProxyRoute,omitempty"`
}
