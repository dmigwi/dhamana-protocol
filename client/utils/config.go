// Copyright (c) 2023 Migwi Ndung'u
// See LICENSE for details.

package utils

import "os"

const (
	// WelcomeText is an easter egg placed in the code for cryptography enthusiasts
	// to attempt solving.
	// Should you be successful in decrypting it, please do what it says!
	WelcomeText = "Pi encrypted  715b48145c501951595355505d194050475a135a5d4251595d5d40021301705647585b5d600c7d527f615760417f63477d5a696045610348437f470d790757"

	// FullDateformat defines the full date format supported
	FullDateformat = "Mon 15:04:05 2006-01-02"

	// TLSCertFile defines the tls certificate file stored at the datadir directory.
	// The certificate is auto-generated if its missing.
	TLSCertFile = "server.crt"

	// TLSKeyFile defines the tls key file stored at the datadir directory.
	// The key file is auto-generated if its missing.
	TLSKeyFile = "server.key"

	// FilePerm defines the file permission used to manage the application data.
	FilePerm = os.FileMode(0o0700)

	// JSONRPCVersion defines the JSON version supportted for all the backend requests
	// recieved by the server.
	JSONRPCVersion = "2.0"
)
