package main

import (
	"encoding/base64"
	"fmt"
	"os"
)

func cmd() {
	var err error
	var tmpStr []byte

	fmt.Println("Birb Mockery Template Generator")
	fmt.Println("")

	fmt.Println("Generating mockery template directory...")
	err = os.MkdirAll("templates", 0755)
	assert(err)

	fmt.Println("Generating mockery template...")
	tmpStr, err = base64.StdEncoding.DecodeString(MockeryTemplate)
	assert(err)
	err = os.WriteFile("templates/mockery.go.tmpl", tmpStr, 0644)
	assert(err)

	fmt.Println("Generating mockery template schema...")
	tmpStr, err = base64.StdEncoding.DecodeString(MockerySchema)
	assert(err)
	err = os.WriteFile("templates/mockery.go.tmpl.schema.json", tmpStr, 0644)
	assert(err)
}

func assert(err error) {
	if err != nil {
		panic(err)
	}
}
