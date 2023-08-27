### Point of Access (PoA) Applications

This describes a list of applications that offer users already vetted by an approved trust organisation, access to the dhamana-protocol.
Communication between PoA Apps and the dhamana-protocol server backend is double encrypted to reduced chances of a middle man attack.

In the [Dhamana Protocol](../client/doc.go) overview the PoA app are at the top layer.

## List of PoA Applications Supported
1. [Lotus](/poa/lotus/README.md) - Implemented with [Golang](https://go.dev/) and [GIO](https://gioui.org/) library.
2. ..

## Simple Golang PoA Application Implemented.

```go
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

const clientCert = `-----BEGIN CERTIFICATE-----
MIIDwTCCAqmgAwIBAgIUUsPJTjeAN0vn8+KsMyMoUjTIfnYwDQYJKoZIhvcNAQEL
BQAwNzEQMA4GA1UEAwwHMC4wLjAuMDELMAkGA1UEBhMCVVMxFjAUBgNVBAcMDVNh
biBGcmFuc2lzY28wHhcNMjMwODExMDkxNjUyWhcNMjQwODEwMDkxNjUyWjCBjTEL
MAkGA1UEBhMCVVMxEzARBgNVBAgMCkNhbGlmb3JuaWExFjAUBgNVBAcMDVNhbiBG
cmFuc2lzY28xGTAXBgNVBAoMEERoYW1hbmEtcHJvdG9jb2wxJDAiBgNVBAsMG0Ro
YW1hbmEtcHJvdG9jb2wgY2xpZW50LWRldjEQMA4GA1UEAwwHMC4wLjAuMDCCASIw
DQYJKoZIhvcNAQEBBQADggEPADCCAQoCggEBAOkVZR6ZUqQ1T9leO0g9ltX5o6bx
VX/rBDdhDQ9bGhw57zQxxJAe0XfmnVPyB+NLQ9q8WiS0asTYKzY65PBETfxjkmsm
Ot99Lyigewu/5MGclZsGDdh0g/tqzp9qSlQEOrO2JE9YlnL8DrmvtaIqyKC13u+5
bbon39MvyXNbUFnAnhRiDNWIlziURR0ydKRgV9MDMSK0P1J3i75bwfiflqdIB/S9
o0ettzDn8lq14TycIXbKidgwYo8DlPrdZP404c0W3qlEUP10kwZkg/CamU8FOL7I
urejbSPJ2d2ZM1DTkVCpmCfZ+3KwWpTe0yAgWW4qWf9GmkXFHS2asCJgUnMCAwEA
AaNuMGwwHwYDVR0jBBgwFoAUCChoVu/oFWiNEQtdrm5lnt9584swCQYDVR0TBAIw
ADALBgNVHQ8EBAMCBPAwEgYDVR0RBAswCYIHMC4wLjAuMDAdBgNVHQ4EFgQUvdjy
LiA9wDxXvYFDyitDdGWWP/MwDQYJKoZIhvcNAQELBQADggEBAC6QNho7D+gryQhY
gzjfj9oTbuj8Hknq4oL3cfYWYMKoRY4wMzyZPoLJccy19GWlm/Gc8MohxuwAmUlP
P9LXxWscwhbNWM27eovEfXxAmrLWvKxdOBXgnSw8MFYiGnO4YLg12ojwbfhhTXir
8pvnVG5BMMT85lUmC78uzPIfLTYW1f0GilpXtxNdSflRga3KVyXrAPE/AKTuDPEw
HO0+tCaprNcYaNRw+z7ObtSsb1tCtKWOSHaFMvMDf6DW+wIqhETQh6r5808g+y7X
NUzFRTqfzpKKefxbBtqi31Ep7XLa/ILnrrXCN4kzZ3NfZRTIQ6YgaMTxLncOqYcy
RBDIGGA=
-----END CERTIFICATE-----`

const clientKey = `-----BEGIN PRIVATE KEY-----
MIIEvAIBADANBgkqhkiG9w0BAQEFAASCBKYwggSiAgEAAoIBAQDpFWUemVKkNU/Z
XjtIPZbV+aOm8VV/6wQ3YQ0PWxocOe80McSQHtF35p1T8gfjS0PavFoktGrE2Cs2
OuTwRE38Y5JrJjrffS8ooHsLv+TBnJWbBg3YdIP7as6fakpUBDqztiRPWJZy/A65
r7WiKsigtd7vuW26J9/TL8lzW1BZwJ4UYgzViJc4lEUdMnSkYFfTAzEitD9Sd4u+
W8H4n5anSAf0vaNHrbcw5/JateE8nCF2yonYMGKPA5T63WT+NOHNFt6pRFD9dJMG
ZIPwmplPBTi+yLq3o20jydndmTNQ05FQqZgn2ftysFqU3tMgIFluKln/RppFxR0t
mrAiYFJzAgMBAAECggEAFgUEFTdINrnP81RoUVN63ngZB8/LfgGGgnIpY7FS7tEM
/LypjzFPef5QD/rQ1/OSU3sby6qSXkpLEB0m57PsPNMYpVq/5mvTJV6+gRuofB9T
hYkK4c9lcJO54Bbl3WNmOvgSx2mHvREm6vqN3uYf0gmcfWqNJ9sAgC5cNQkEDEEX
VbpTNirZzL5b8/efxR1MBErKg6bgsET/SXFTSw/10N2FwtiOAcbiFgC4X1/s7Eib
zuXLsRxGNJAajC2PEMDWQO+T6/gGX9pP3pil5sfqrVQCaRVFzZsWu4oG+tEGgnvu
PDCqh13o931JRy5Xpa5XZlWc5UyTvjVlhBASRwvfgQKBgQDsexs2Wgkoh0Hiabj+
/MU2yV/C/mzkWTJcKVpMhhZXQkYL1sC1ATgrpr3jX75aqGpU+4oLB0tGuWlfsxFj
lvXHzVHbpSJ5M18dgEgckumk2910oLg9adWsXWCSt0xkASrVWPLmaPNQqCZYqvix
tSaymetttsIHMPmCy5lI06hMYQKBgQD8UoC2wwcHLHo+cGXQ+dg47SFCtTca8mth
3ukruLGDD6MwZLxBLaQb1eiFHyzvZsk9RaBCxMR/48kpRyGEYZmAhjaLu1WKnH1k
ZVB0SkruHQpIKm1D/JYx2ZfJcanmBI1/WKbJc1qJ0fimoBEgO0maEOAaK2csnZxI
3d+7ayHvUwKBgHE6BG6Cr43jLS2ON1CHkJnJ03sWvOacupscBatMLFg9WDKE8aH4
4n8sCBFdH1Ri/P6RrafYJzfGwOhcYcAQYL/40+/Z4marrSf/6wcbZJlV2HPmHDDz
gqZT01CMSRw83thmDW864v1EdY/Q1OCpfszXG71dWwt8bIsulsci6JshAoGAKXsO
ufz57M39EsK7mk1YpJMnQqYz9vQffyl8P7nPRPKPK2eEI1rzfbf+z9O+OWU2dCI8
JH3gp/20llqhQfghmiV2ViZn+6+aVaTqQxPrmZWgmRiQefrOXkedUnqjKbNZ57OV
R8z0929TZ2EtL5VPlkpr7SFxhr4qcTg4jcEhBQ0CgYBrKrd2Ca7ZXA2sK3K+FJyu
YDGvWDkQ4vF0xk7JrXWhyGvwkmkmUTqZ8zhk73RKBA0OLhBVH97lF1SM1+qu677F
rWaEiOGn7czixcTa0ZiQhogpjRgUXf193jt6pURcOv9P/hWtS+vAmvKP2x41DasT
haVNL+FzEkSbVr5kvhlF6g==
-----END PRIVATE KEY-----`

func main() {
	cert, err := tls.X509KeyPair([]byte(clientCert), []byte(clientKey))
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

```