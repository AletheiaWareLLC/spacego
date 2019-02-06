/*
 * Copyright 2019 Aletheia Ware LLC
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
	"bufio"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha512"
	//"errors"
	"github.com/AletheiaWareLLC/bcgo"
	"github.com/AletheiaWareLLC/financego"
	"github.com/golang/protobuf/proto"
	"io"
	"log"
	"net"
	"strconv"
)

const (
	SPACE_HOST           = "space.aletheiaware.com"
	SPACE_WEBSITE        = "https://space.aletheiaware.com"
	SPACE_PREFIX         = "Space "
	SPACE_PREFIX_FILE    = "Space File "
	SPACE_PREFIX_META    = "Space Meta "
	SPACE_PREFIX_PREVIEW = "Space Preview "

	PORT_UPLOAD = 23232
)

func OpenFileChannel(alias string) (*bcgo.Channel, error) {
	return bcgo.OpenChannel(SPACE_PREFIX_FILE + alias)
}

func OpenMetaChannel(alias string) (*bcgo.Channel, error) {
	return bcgo.OpenChannel(SPACE_PREFIX_META + alias)
}

func GetFile(files *bcgo.Channel, alias string, key *rsa.PrivateKey, recordHash []byte, callback func(*bcgo.BlockEntry, []byte)) error {
	return files.Read(alias, key, recordHash, callback)
}

func GetMeta(metas *bcgo.Channel, alias string, key *rsa.PrivateKey, recordHash []byte, callback func(*bcgo.BlockEntry, *Meta)) error {
	return metas.Read(alias, key, recordHash, func(entry *bcgo.BlockEntry, data []byte) {
		// Unmarshal as Meta
		meta := &Meta{}
		err := proto.Unmarshal(data, meta)
		if err != nil {
			log.Println(err)
		} else {
			callback(entry, meta)
		}
	})
}

func GetCustomer(node *bcgo.Node) (*financego.Customer, error) {
	// Open Customer Channel
	channel, err := financego.OpenCustomerChannel()
	if err != nil {
		return nil, err
	}
	// Sync channel
	if err := channel.Sync(); err != nil {
		return nil, err
	}
	return financego.GetCustomerSync(channel, node.Alias, node.Key, node.Alias)
}

func GetSubscription(node *bcgo.Node) (*financego.Subscription, error) {
	channel, err := financego.OpenSubscriptionChannel()
	if err != nil {
		return nil, err
	}
	// Sync channel
	if err := channel.Sync(); err != nil {
		return nil, err
	}
	return financego.GetSubscriptionSync(channel, node.Alias, node.Key, node.Alias)
}

func NewBundle(node *bcgo.Node, payload []byte) (*StorageRequest_Bundle, error) {
	// TODO compress payload

	// Generate a random shared key
	key := make([]byte, bcgo.AES_KEY_SIZE_BYTES)
	if _, err := io.ReadFull(rand.Reader, key); err != nil {
		return nil, err
	}

	// Create cipher
	c, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	// Create galois counter mode
	gcm, err := cipher.NewGCM(c)
	if err != nil {
		return nil, err
	}

	// Generate nonce
	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, err
	}

	// Encrypt payload
	encryptedPayload := append(nonce, gcm.Seal(nil, nonce, payload, nil)...)

	// Hash encrypted payload
	hashed := bcgo.Hash(encryptedPayload)

	// Sign hash of encrypted payload
	signature, err := bcgo.CreateSignature(node.Key, hashed, bcgo.SignatureAlgorithm_SHA512WITHRSA_PSS)
	if err != nil {
		return nil, err
	}

	// Encrypt AES key with RSA public key
	encryptedKey, err := rsa.EncryptOAEP(sha512.New(), rand.Reader, &node.Key.PublicKey, key, nil)
	if err != nil {
		return nil, err
	}

	return &StorageRequest_Bundle{
		Key:                    encryptedKey,
		KeyEncryptionAlgorithm: bcgo.EncryptionAlgorithm_RSA_ECB_OAEPPADDING,
		Payload:                encryptedPayload,
		CompressionAlgorithm:   bcgo.CompressionAlgorithm_UNKNOWN_COMPRESSION,
		EncryptionAlgorithm:    bcgo.EncryptionAlgorithm_AES_GCM_NOPADDING,
		Signature:              signature,
		SignatureAlgorithm:     bcgo.SignatureAlgorithm_SHA512WITHRSA_PSS,
	}, nil
}

func Upload(host string, request *StorageRequest) (*StorageResponse, error) {
	address := host + ":" + strconv.Itoa(PORT_UPLOAD)
	log.Println(address)
	connection, err := net.Dial("tcp", address)
	if err != nil {
		return nil, err
	}
	defer connection.Close()
	reader := bufio.NewReader(connection)
	writer := bufio.NewWriter(connection)
	if err := WriteStorageRequest(writer, request); err != nil {
		return nil, err
	}
	return ReadStorageResponse(reader)
}

func ReadStorageRequest(reader *bufio.Reader) (*StorageRequest, error) {
	// Unmarshal as StorageRequest
	request := &StorageRequest{}
	if err := bcgo.ReadDelimitedProtobuf(reader, request); err != nil {
		return nil, err
	}
	return request, nil
}

func ReadStorageResponse(reader *bufio.Reader) (*StorageResponse, error) {
	response := &StorageResponse{}
	if err := bcgo.ReadDelimitedProtobuf(reader, response); err != nil {
		return nil, err
	}
	return response, nil
}

func WriteStorageRequest(writer *bufio.Writer, request *StorageRequest) error {
	return bcgo.WriteDelimitedProtobuf(writer, request)
}

func WriteStorageResponse(writer *bufio.Writer, response *StorageResponse) error {
	return bcgo.WriteDelimitedProtobuf(writer, response)
}
