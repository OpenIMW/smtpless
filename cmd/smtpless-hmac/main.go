package main

import (
	"fmt"
	"github.com/openimw/smtpless/utils"
	"os"
)

func usage() {
	fmt.Println(`
 Simple tool to generate and verify hmac hashes.

 smtpless-hmac hash [HOST] [SECRET]
 smtpless-hmac verify [HASH] [HOST] [SECRET]
	`)
}

func hash(host string, secret string) {
	fmt.Println(utils.Hash(host, secret))
}

func verify(hash string, host string, secret string) {

	if hash == utils.Hash(host, secret) {
		fmt.Println("Valid")
		return
	}

	fmt.Println("Invalid")
	os.Exit(1)
}

func main() {

	if len(os.Args) >= 2 {
		switch os.Args[1] {
		case "hash":
			hash(os.Args[2], os.Args[3])
			return
		case "verify":
			verify(os.Args[2], os.Args[3], os.Args[4])
			return
		}
	}

	usage()
}
