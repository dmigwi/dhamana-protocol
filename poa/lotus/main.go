// Copyright (c) 2023 Migwi Ndung'u
// See LICENSE for details.

package main

import (
	"crypto/tls"
	"io"
	"log"
	"net/http"
	"time"
)

const (
	clientCert    = "certs/client.crt"
	clientKeyFile = "certs/client.key"
)

func main() {
	cert, err := tls.LoadX509KeyPair(clientCert, clientKeyFile)
	if err != nil {
		log.Fatalf("server: loadkeys: %s", err)
	}

	tr := &http.Transport{
		MaxIdleConns:       10,
		IdleConnTimeout:    30 * time.Second,
		DisableCompression: true,
		TLSClientConfig: &tls.Config{
			Certificates:       []tls.Certificate{cert},
			InsecureSkipVerify: true,
		},
	}

	client := &http.Client{Transport: tr}
	resp, err := client.Get("https://127.0.0.1:30443")
	if err != nil {
		log.Fatalf("request err : %v", err)
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("reading the request body failed: %v", err)
	}

	log.Printf("client: %s", string(data))
}
