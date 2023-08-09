// Copyright (c) 2023 Migwi Ndung'u
// See LICENSE for details.

package main

import "fmt"

// deploymentFilePath defines the path where npm stores the deployment log information
// on deploying a contract.
const deploymentFilePath = "build/deploy_*.log"

// FieldInfo defines the properties of the auto generated file functions using the
// the field names of the properties extracted from the truffle deployment logs.
type FieldInfo struct {
	// Identifier defines field name as used in truffle deployment log.
	Identifier string
	// FuncName defines the function name to be used in the generated file.
	FuncName string
	// ReturnType defines the data type to be returned in the generated file.
	ReturnType string
	// Comment defines description of the function in the generated file.
	Comment string
	// Value defines the actual value to be returned once regex expressions are
	// evaluated for each Identifier.
	Value interface{}
}

var fields = []FieldInfo{
	{
		Identifier: "network",
		FuncName:   "GetNetwork",
		ReturnType: "string",
		Comment:    "returns the network used to make the deployment.",
	}, {
		Identifier: "networkId",
		FuncName:   "GetChainID",
		ReturnType: "uint64",
		Comment:    "returns the chain ID of the network used to make the deployment.",
	}, {
		Identifier: "contractName",
		FuncName:   "GetContractName",
		ReturnType: "string",
		Comment:    "returns the contract name deployed.",
	}, {
		Identifier: "address",
		FuncName:   "GetContractAddress",
		ReturnType: "string",
		Comment:    "returns the address of the deployed contract.",
	}, {
		Identifier: "deployed",
		FuncName:   "GetIsDeployed",
		ReturnType: "boolean",
		Comment:    "confirms if the contract was successfully deployed. If yes, it succeeded",
	}, {
		Identifier: "transactionHash",
		FuncName:   "GetTransactionHash",
		ReturnType: "string",
		Comment:    "returns the tx hash when the contract we deployed.",
	}, {
		Identifier: "timestamp",
		FuncName:   "GetDeploymentTime",
		ReturnType: "uint64",
		Comment:    "returns the timestamp in seconds when the contract was actually deployed ",
	},
}

const deploymentConfig = `
// Deployment configuration code is generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package deployment

{{ range . }}

// {{ .FuncName }} {{ .Comment }}
func {{ .FuncName }}() ({{ .ReturnType }}) {
	return {{ .Value }}
}
{{end}}
`

func main() {
	fmt.Println("GOT HERE")
}
