### Point of Access (PoA) Applications

This describes a list of applications that offer users already vetted by an approved trust organisation, access to the dhamana-protocol.
Communication between PoA Apps and the dhamana-protocol server backend is double encrypted to reduced chances of a middle man attack.

In the [Dhamana Protocol](../client/doc.go) overview the PoA app are at the top layer.

## List of PoA Applications Supported
1. [Lotus](/poa/lotus/README.md) - Implemented with [Golang](https://go.dev/) and [GIO](https://gioui.org/) library.
2. ..

## Simple Golang PoA Application Implemented.

Certicates used below are fake and are just required to test for server connection.

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

const clientKey = `-----BEGIN PRIVATE KEY-----
MIIEvgIBADANBgkqhkiG9w0BAQEFAASCBKgwggSkAgEAAoIBAQC1tovRQXY4Ni4l
Gqk0dm4uNSAimpP/Px7ABlw23nTkGvixnfpZZtYKrKYPLQvBPvQso7xeftuCgrly
H46BzbR8JVUfarP1VpTMlZIQVnJw0++B3u0/ddgHJucaU/03OkuMpCOqfqabFz1M
GS95wdSCScuZQtJP/MW25KwM+zdU52e/gQSN2dO0JyNkDvn26XbqbyJJ6TYVJWVB
H7W4sBR1xjQU2YNTnUKhxsld4lumFuOAFWTOPuhac+H9ydSHOWa+EhQrGccvEHjU
Pmh2xFkxljo0eJEdsbi6aDzmHGLpsnjZwjbkywd0V15tq70XkYSbtnKaf8ph4r+9
i5NXikO/AgMBAAECggEALpeT2h/S1C7wPgtL/2uubKKqjrTb5cKle88lrdv0VNil
k8VR/8Jid3I1UgbW0MH3kcqZ7hDQ7/Hc4uo8fAPmlz5rSRbu9aUxmhNv1EiWJ1/w
NXiXlIH+1jafYxzN/G8yF+muS3UV+wZGbVC335xXhOCvF6kOi/vgJjkT9HAli9sr
FmgKtit+A5sotFXmqCVB0EROYPmzj8j90SEeaMHdMGqJ8Vj/1q67BLAE/f2awB5U
tCR4eFwE4zSq2D5xBb1cL+L8IBors+Tp1jzPyhpK7GpOPWu+xjt8LGKdUcqx7mdF
Bi0bazBXcWEv8TbsLczLP7kHZIoLc9ufXkTp8YqxoQKBgQDdXGn7EqmLE/xfRErD
LutdcWe6CCVx6I9pbTpC3TnpVr4TtCyora89Pd09ShTdGvLQEVEdgxN2YpeKYL3C
jcST+qAaLt82QVs+mHmSG8Gbn6bZ/LWy6qZZEvVVbooa7Ua+nGcsCgJgsC7GdXrK
mKsXpRMP1UaROZxUqMIpywzlFQKBgQDSJdxWOq7JtK9sUgAhmscqJYMCiDUp21jm
SMqVA8xdqjfODWOmTrBD3Ed834oV8qhrog76nwNEpGm/cA4vCf7B3Vs/GkwO6qbR
3B31Y7FvCnOSEbSqnaa6g7JgHizigfpz5TT+LM9Qg0+ieaP6Cn5zHDuPqGDFpqyf
dX2pzQNigwKBgGHApriORDq7p642TUGmXZ/VLbY0VLzZs0MeTiUq5qEJgkTXQwV2
NbW4tROUvGPru6BwlT6QHK8h2MPt9r6MtmuWuM73NfESqYWZ9c203imoNhl7hI2v
G6ioO5jviKNddulDzjffb69c/jr7tC71flChwCo0x8XoCAZGw/+KwHYFAoGBAJ2e
PhQC8cRiHE0vd9+8mnNXLVtB1DYvyg73O9LmxWrfV/nZewtq67QKTSgw9f4eQgpw
w7FggPAELTikEE9hvM2lfGHpFHD/uN2grmu2OYgim6pMU2jA1CQC0VBccaf2e2Zf
3Q5jh59Izfr8J2xMYKlv3JCUZvj4WXNEiVtJZKeHAoGBAKR3KCny3fAnmX2Lg0EX
3ekcZGdJLuqfPetAa30hIwovqFcvTrdHW9PZ7PbttDhjh2uj+RpiTiJwXwu+xKdb
SsiDIvpkcLH8vreBwZa1ipSdq4+fcqKtkvkvqeCzt2BOqCxNpKBlaEO6dkQA36H1
Ui0Tnga6gdiZCvVFzft7f9n8
-----END PRIVATE KEY-----`

const clientCert = `-----BEGIN CERTIFICATE-----
MIIDrDCCApSgAwIBAgIEXz3xDTANBgkqhkiG9w0BAQsFADBvMQswCQYDVQQGEwJL
RTEZMBcGA1UEAwwQZGhhbWFuYS1wcm90b2NvbDEQMA4GA1UECAwHTmFpcm9iaTEV
MBMGA1UEBwwMRGVmYXVsdCBDaXR5MRwwGgYDVQQKDBNEZWZhdWx0IENvbXBhbnkg
THRkMB4XDTIzMDgyODIxMTUwNVoXDTI1MDgyNzIxMTUwNVowbzELMAkGA1UEBhMC
S0UxGTAXBgNVBAMMEGRoYW1hbmEtcHJvdG9jb2wxEDAOBgNVBAgMB05haXJvYmkx
FTATBgNVBAcMDERlZmF1bHQgQ2l0eTEcMBoGA1UECgwTRGVmYXVsdCBDb21wYW55
IEx0ZDCCASIwDQYJKoZIhvcNAQEBBQADggEPADCCAQoCggEBALW2i9FBdjg2LiUa
qTR2bi41ICKak/8/HsAGXDbedOQa+LGd+llm1gqspg8tC8E+9CyjvF5+24KCuXIf
joHNtHwlVR9qs/VWlMyVkhBWcnDT74He7T912Acm5xpT/Tc6S4ykI6p+ppsXPUwZ
L3nB1IJJy5lC0k/8xbbkrAz7N1TnZ7+BBI3Z07QnI2QO+fbpdupvIknpNhUlZUEf
tbiwFHXGNBTZg1OdQqHGyV3iW6YW44AVZM4+6Fpz4f3J1Ic5Zr4SFCsZxy8QeNQ+
aHbEWTGWOjR4kR2xuLpoPOYcYumyeNnCNuTLB3RXXm2rvReRhJu2cpp/ymHiv72L
k1eKQ78CAwEAAaNQME4wHQYDVR0OBBYEFFXlWHBbIY40ZhGpoKRs6dwxNZwfMB8G
A1UdIwQYMBaAFFXlWHBbIY40ZhGpoKRs6dwxNZwfMAwGA1UdEwQFMAMBAf8wDQYJ
KoZIhvcNAQELBQADggEBADXgEbxGbWIty9cdD2m72f2Fq0WizyEWSeZm/czpYBfi
rv6LhhC3t883FNaK8vp2OBTtjPgU7DYJ629bvyauOu19W4ig+QvF0p1NWUH63mnm
IBKlcD5Lu1GcOfZiCNu8rzcxaHpwfJN/7edkW9vJ7SVaRiPGxk+0ZcCnKALTPDS+
uZ22cB0f3btGMIVKfPcqxcqNvG7bmm4MiIxX+W/UUMNdSkGMqQta2E7uVSioTg3p
X0sOFym1yIr8hB9KaRL0CW+avHH8Pl+m3e8K3ndVk7Fm1bnkKonBvdzO+wbaWMPI
FUQhkms1UNbkyjEqwiCefrXpsFmlm1B9QfuGLO7sVL0=
-----END CERTIFICATE-----
`

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

On running this with a server on **https://127.0.0.1:30443** you should get the following response to confirm connection to the dhamana-protocol server.

```
2023/08/29 00:23:20 client: "Pi encrypted  715b48145c501951595355505d194050475a135a5d4251595d5d40021301705647585b5d600c7d527f615760417f63477d5a696045610348437f470d790757"
```

<h4 align="center">HAPPY BUILDING ON DHAMANA PROTOCOL!</h4>