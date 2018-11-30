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
	"errors"
	"fmt"
	//bc "github.com/AletheiaWareLLC/bcgo"
	bcutils "github.com/AletheiaWareLLC/bcgo/utils"
	"github.com/golang/protobuf/proto"
	"golang.org/x/crypto/ssh/terminal"
	"io"
	"log"
	"os"
	"os/user"
	"path"
	"syscall"
)

const (
	SPACE_PREFIX         = "Space "
	SPACE_REGISTRATION   = "Space Registration"
	SPACE_FILE_PREFIX    = "Space File "
	SPACE_META_PREFIX    = "Space Meta "
	SPACE_PREVIEW_PREFIX = "Space Preview "
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

	if bcutils.HasRSAPrivateKey(keystore) {
		fmt.Println("Found keystore under " + keystore)
		var password []byte
		pwd, ok := os.LookupEnv("PASSWORD")
		if ok {
			password = []byte(pwd)
		} else {
			fmt.Print("Enter keystore password: ")
			var err error
			password, err = terminal.ReadPassword(int(syscall.Stdin))
			if err != nil {
				return nil, err
			}
			fmt.Println()
		}
		key, err := bcutils.GetRSAPrivateKey(keystore, password)
		if err != nil {
			return nil, err
		}

		publicKeyBase64, err := bcutils.RSAPublicKeyToBase64(&key.PublicKey)
		if err != nil {
			return nil, err
		}
		fmt.Println(publicKeyBase64)
		return key, nil
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

		key, err := bcutils.CreateRSAPrivateKey(keystore, password)
		if err != nil {
			return nil, err
		}

		publicKeyBase64, err := bcutils.RSAPublicKeyToBase64(&key.PublicKey)
		if err != nil {
			return nil, err
		}
		fmt.Printf("Successfully created key pair, visit https://space.aletheiaware.com/register and register your public key;\n\n%s\n", publicKeyBase64)
		return key, nil
	}
}

func ReadStorageRequest(reader io.Reader) (*StorageRequest, error) {
	var data [1024]byte
	n, err := reader.Read(data[:])
	if err != nil {
		return nil, err
	}
	if n <= 0 {
		return nil, errors.New("Could not read data")
	}
	size, s := proto.DecodeVarint(data[:])
	if s <= 0 {
		return nil, errors.New("Could not read size")
	}

	// Create new larger buffer
	buffer := make([]byte, size)
	// Calculate data received
	count := uint64(n - s)
	// Copy data into new buffer
	copy(buffer[:count], data[s:n])
	// Read addition bytes
	for count < size {
		n, err := reader.Read(buffer[count:])
		if err != nil {
			return nil, err
		}
		if n <= 0 {
			return nil, errors.New("Could not read data")
		}
		count = count + uint64(n)
	}

	// Unmarshal as StorageRequest
	request := &StorageRequest{}
	if err = proto.Unmarshal(buffer, request); err != nil {
		return nil, err
	}
	return request, nil
}

func WriteStorageResponse(writer io.Writer, response *StorageResponse) error {
	data, err := proto.Marshal(response)
	if err != nil {
		return err
	}
	size := uint64(len(data))
	// Write response size varint
	if _, err := writer.Write(proto.EncodeVarint(size)); err != nil {
		return err
	}
	// Write response data
	if _, err := writer.Write(data); err != nil {
		return err
	}
	return nil
}
