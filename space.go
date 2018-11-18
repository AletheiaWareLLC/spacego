/*
 * Copyright 2018 Aletheia Ware LLC
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 * http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

// Package containing utilities for Space in Go
package spacego

import (
	"bytes"
	"crypto/rsa"
	"fmt"
	//bc "github.com/AletheiaWareLLC/bcgo"
	bcutils "github.com/AletheiaWareLLC/bcgo/utils"
	"golang.org/x/crypto/ssh/terminal"
	"log"
	"os"
	"os/user"
	"path"
	"syscall"
)

func GetOrCreatePrivateKey() (*rsa.PrivateKey, error) {
	keystore, ok := os.LookupEnv("KEYSTORE")
	if !ok {
		u, err := user.Current()
		if err != nil {
			return nil, err
		}
		keystore = path.Join(u.HomeDir, "bc")
	}

	var key *rsa.PrivateKey

	if bcutils.HasRSAPrivateKey(keystore) {
		fmt.Println("Found keystore under " + keystore)
		fmt.Print("Enter keystore password: ")
		password, err := terminal.ReadPassword(int(syscall.Stdin))
		if err != nil {
			return nil, err
		}
		fmt.Println()
		key, err = bcutils.GetRSAPrivateKey(keystore, password)
		if err != nil {
			return nil, err
		}

		publicKeyBase64, err := bcutils.RSAPublicKeyToBase64(&key.PublicKey)
		if err != nil {
			return nil, err
		}
		fmt.Println(publicKeyBase64)
	} else {
		fmt.Println("Creating keystore under " + keystore)

		fmt.Print("Enter keystore password: ")
		password, err := terminal.ReadPassword(int(syscall.Stdin))
		if err != nil {
			return nil, err
		}
		fmt.Println()

		fmt.Print("Confirm keystore password: ")
		confirm, err := terminal.ReadPassword(int(syscall.Stdin))
		if err != nil {
			return nil, err
		}
		fmt.Println()

		if !bytes.Equal(password, confirm) {
			log.Fatal("Passwords don't match")
		}

		key, err = bcutils.CreateRSAPrivateKey(keystore, password)
		if err != nil {
			return nil, err
		}

		publicKeyBase64, err := bcutils.RSAPublicKeyToBase64(&key.PublicKey)
		if err != nil {
			return nil, err
		}
		fmt.Printf("Successfully created key pair, visit https://space.aletheiaware.com/register and register your public key;\n\n%s\n", publicKeyBase64)
	}
	return key, nil
}
