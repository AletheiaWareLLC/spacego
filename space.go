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

func GetShare(shares *bcgo.Channel, alias string, key *rsa.PrivateKey, recordHash []byte, callback func(*bcgo.BlockEntry, []byte, *Share) error) error {
	return shares.Read(alias, key, recordHash, func(entry *bcgo.BlockEntry, key, data []byte) error {
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

func GetSharedMeta(owner string, recordHash, key []byte, callback func(*bcgo.BlockEntry, *Meta) error) error {
	metas, err := bcgo.OpenChannel(SPACE_PREFIX_META + owner)
	if err != nil {
		return err
	}
	block, err := metas.GetRemoteBlock(&bcgo.Reference{
		ChannelName: metas.Name,
		RecordHash:  recordHash,
	})
	if err != nil {
		return err
	}
	for _, entry := range block.Entry {
		if bytes.Equal(recordHash, entry.RecordHash) {
			data, err := bcgo.DecryptPayload(entry, key)
			if err != nil {
				return err
			}
			// Unmarshal as Meta
			meta := &Meta{}
			if err := proto.Unmarshal(data, meta); err != nil {
				return err
			} else if err := callback(entry, meta); err != nil {
				return err
			}
		}
	}
	return nil
}

func GetSharedFile(owner string, recordHash, key []byte, callback func(*bcgo.BlockEntry, []byte) error) error {
	files, err := bcgo.OpenChannel(SPACE_PREFIX_FILE + owner)
	if err != nil {
		return err
	}
	block, err := files.GetRemoteBlock(&bcgo.Reference{
		ChannelName: files.Name,
		RecordHash:  recordHash,
	})
	if err != nil {
		return err
	}
	for _, entry := range block.Entry {
		if bytes.Equal(recordHash, entry.RecordHash) {
			data, err := bcgo.DecryptPayload(entry, key)
			if err != nil {
				return err
			}
			return callback(entry, data)
		}
	}
	return nil
}

func PostRecord(feature string, record *bcgo.Record) (*bcgo.Reference, error) {
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
		return nil, errors.New("Post status: " + response.Status)
	}
}
