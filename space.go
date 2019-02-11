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
	"bytes"
	"crypto/rsa"
	"errors"
	"github.com/AletheiaWareLLC/bcgo"
	"github.com/golang/protobuf/proto"
	"net/http"
)

const (
	SPACE_HOST           = "space.aletheiaware.com"
	SPACE_WEBSITE        = "https://space.aletheiaware.com"
	SPACE_PREFIX         = "Space-"
	SPACE_PREFIX_FILE    = "Space-File-"
	SPACE_PREFIX_META    = "Space-Meta-"
	SPACE_PREFIX_PREVIEW = "Space-Preview-"
	SPACE_PREFIX_SHARE   = "Space-Share-"
	SPACE_PREFIX_TAG     = "Space-Tag-"
)

func GetFile(files *bcgo.Channel, alias string, key *rsa.PrivateKey, recordHash []byte, callback func(*bcgo.BlockEntry, []byte, []byte) error) error {
	return files.Read(alias, key, recordHash, callback)
}

func GetMeta(metas *bcgo.Channel, alias string, key *rsa.PrivateKey, recordHash []byte, callback func(*bcgo.BlockEntry, []byte, *Meta) error) error {
	return metas.Read(alias, key, recordHash, func(entry *bcgo.BlockEntry, key, data []byte) error {
		// Unmarshal as Meta
		meta := &Meta{}
		if err := proto.Unmarshal(data, meta); err != nil {
			return err
		} else if err := callback(entry, key, meta); err != nil {
			return err
		}
		return nil
	})
}

func GetSharedMeta(metas *bcgo.Channel, recordHash, key []byte, callback func(*bcgo.BlockEntry, *Meta) error) error {
	return metas.Iterate(func(h []byte, b *bcgo.Block) error {
		for _, e := range b.Entry {
			if recordHash == nil || bytes.Equal(recordHash, e.RecordHash) {
				data, err := bcgo.DecryptPayload(e, key)
				if err != nil {
					return err
				}
				// Unmarshal as Meta
				meta := &Meta{}
				if err := proto.Unmarshal(data, meta); err != nil {
					return err
				} else if err := callback(e, meta); err != nil {
					return err
				}
				return nil
			}
		}
		return nil
	})
}

func GetSharedFile(shares *bcgo.Channel, recordHash, key []byte, callback func(*bcgo.BlockEntry, []byte) error) error {
	return shares.Iterate(func(h []byte, b *bcgo.Block) error {
		for _, e := range b.Entry {
			if recordHash == nil || bytes.Equal(recordHash, e.RecordHash) {
				data, err := bcgo.DecryptPayload(e, key)
				if err != nil {
					return err
				}
				return callback(e, data)
			}
		}
		return nil
	})
}

func GetTag(tags *bcgo.Channel, alias string, key *rsa.PrivateKey, recordHash []byte, callback func(*bcgo.BlockEntry, []byte, *Tag) error) error {
	return tags.Read(alias, key, recordHash, func(entry *bcgo.BlockEntry, key, data []byte) error {
		// Unmarshal as Tag
		meta := &Tag{}
		if err := proto.Unmarshal(data, meta); err != nil {
			return err
		} else if err := callback(entry, key, meta); err != nil {
			return err
		}
		return nil
	})
}

func GetShare(metas *bcgo.Channel, alias string, key *rsa.PrivateKey, recordHash []byte, callback func(*bcgo.BlockEntry, []byte, *Share) error) error {
	return metas.Read(alias, key, recordHash, func(entry *bcgo.BlockEntry, key, data []byte) error {
		// Unmarshal as Share
		share := &Share{}
		if err := proto.Unmarshal(data, share); err != nil {
			return err
		} else if err := callback(entry, key, share); err != nil {
			return err
		}
		return nil
	})
}

func UploadRecord(feature string, record *bcgo.Record) (*bcgo.Reference, error) {
	var buffer bytes.Buffer
	writer := bufio.NewWriter(&buffer)
	err := bcgo.WriteRecord(writer, record)
	if err != nil {
		return nil, err
	}
	client := http.Client{}
	request, err := http.NewRequest("POST", SPACE_WEBSITE+"/mining/"+feature, &buffer)
	if err != nil {
		return nil, err
	}
	response, err := client.Do(request)
	if err != nil {
		return nil, err
	}
	switch response.StatusCode {
	case http.StatusOK:
		defer response.Body.Close()
		return bcgo.ReadReference(bufio.NewReader(response.Body))
	default:
		return nil, errors.New("Upload status: " + response.Status)
	}
}
