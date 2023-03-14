/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package main

import (
	"bytes"
	"flag"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
)

var proxyTarget string
var listen string

func main() {
	log.SetFlags(0)

	flag.StringVar(&proxyTarget, "proxy", "https://app.fossa.com", "The host to which all requests should be proxied.")
	flag.StringVar(&listen, "listen", ":3000", "The local address on which to listen.")
	flag.Parse()

	proxyUrl, err := url.Parse(proxyTarget)
	if err != nil {
		bailf("parse proxy host url: %v", err)
	}

	log.Printf("✨ Serving on '%s', forwarding to '%s'", listen, proxyUrl)

	proxy := httputil.NewSingleHostReverseProxy(proxyUrl)
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		log.Printf("🚀 Forward '%s' to '%s%s'", r.URL.Path, proxyTarget, r.URL.Path)

		// Read the body into memory so that we can write this request multiple times.
		// Typically bodies are streaming, and can only be read once.
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			bailf("read request body: %v", err)
		}

		// Replace the body in the request with a new reader that just reads from the body
		// that was previously buffered above.
		r.Body = ioutil.NopCloser(bytes.NewBuffer(body))

		// Write the request in wire protocol format to the logger.
		// Wire protocol format means "this is the text that is actually sent over HTTP to the server".
		r.Write(log.Writer())

		// Add a blank line after the wire protocol output.
		// This is because requests that don't have a body emit a blank line, but requests with a body
		// just emit the body directly with no trailing line, which makes terminal output look messy.
		// Using this gives us extra blank space at the end of requests that don't have a body,
		// but that's okay because it's more consistent.
		log.Print("\n\n")

		// Update the request URL to match the proxy destination.
		r.URL.Host = proxyUrl.Host
		r.URL.Scheme = proxyUrl.Scheme
		r.Host = proxyUrl.Host

		// Replace the body in the request with a new reader that just reads from the body
		// that was previously buffered at the start of this function, and perform the reverse proxy.
		r.Body = ioutil.NopCloser(bytes.NewBuffer(body))
		proxy.ServeHTTP(w, r)
	})

	http.ListenAndServe(listen, nil)
}

func bailf(msg string, err error) {
	log.Fatalf("😵 error: "+msg, err)
}
