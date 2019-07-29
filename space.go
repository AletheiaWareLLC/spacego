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
	SPACE_HOST_TEST      = "test-space.aletheiaware.com"
	SPACE_PREFIX         = "Space-"
	SPACE_PREFIX_FILE    = "Space-File-"
	SPACE_PREFIX_META    = "Space-Meta-"
	SPACE_PREFIX_PREVIEW = "Space-Preview-"
	SPACE_PREFIX_SHARE   = "Space-Share-"
	SPACE_PREFIX_TAG     = "Space-Tag-"
)

func GetSpaceHost() string {
	if bcgo.IsDebug() {
		return SPACE_HOST_TEST
	}
	return SPACE_HOST
}

func GetSpaceWebsite() string {
	return "https://" + GetSpaceHost()
}

func OpenFileChannel(alias string) *bcgo.PoWChannel {
	// TODO consider lowering threshold once Periodic Validation Blockchains are implemented
	return bcgo.OpenPoWChannel(SPACE_PREFIX_FILE+alias, bcgo.THRESHOLD_STANDARD)
}

func OpenMetaChannel(alias string) *bcgo.PoWChannel {
	return bcgo.OpenPoWChannel(SPACE_PREFIX_META+alias, bcgo.THRESHOLD_STANDARD)
}

func OpenShareChannel(alias string) *bcgo.PoWChannel {
	return bcgo.OpenPoWChannel(SPACE_PREFIX_SHARE+alias, bcgo.THRESHOLD_STANDARD)
}

func OpenPreviewChannel(metaId string) *bcgo.PoWChannel {
	return bcgo.OpenPoWChannel(SPACE_PREFIX_PREVIEW+metaId, bcgo.THRESHOLD_STANDARD)
}

func OpenTagChannel(metaId string) *bcgo.PoWChannel {
	return bcgo.OpenPoWChannel(SPACE_PREFIX_TAG+metaId, bcgo.THRESHOLD_STANDARD)
}

func GetFile(files bcgo.Channel, cache bcgo.Cache, network bcgo.Network, alias string, key *rsa.PrivateKey, recordHash []byte, callback func(*bcgo.BlockEntry, []byte, []byte) error) error {
	return bcgo.Read(files.GetName(), files.GetHead(), nil, cache, network, alias, key, recordHash, callback)
}

func GetMeta(metas bcgo.Channel, cache bcgo.Cache, network bcgo.Network, alias string, key *rsa.PrivateKey, recordHash []byte, callback func(*bcgo.BlockEntry, []byte, *Meta) error) error {
	return bcgo.Read(metas.GetName(), metas.GetHead(), nil, cache, network, alias, key, recordHash, func(entry *bcgo.BlockEntry, key, data []byte) error {
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

func GetShare(shares bcgo.Channel, cache bcgo.Cache, network bcgo.Network, alias string, key *rsa.PrivateKey, recordHash []byte, callback func(*bcgo.BlockEntry, []byte, *Share) error) error {
	return bcgo.Read(shares.GetName(), shares.GetHead(), nil, cache, network, alias, key, recordHash, func(entry *bcgo.BlockEntry, key, data []byte) error {
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

func GetSharedMeta(cache bcgo.Cache, network bcgo.Network, owner string, recordHash, key []byte, callback func(*bcgo.BlockEntry, *Meta) error) error {
	metas := OpenMetaChannel(owner)
	if err := bcgo.LoadHead(metas, cache, network); err != nil {
		return err
	}
	if err := bcgo.Pull(metas, cache, network); err != nil {
		return err
	}
	block, err := bcgo.GetBlockContainingRecord(metas.GetName(), cache, network, recordHash)
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

func GetSharedFile(cache bcgo.Cache, network bcgo.Network, owner string, recordHash, key []byte, callback func(*bcgo.BlockEntry, []byte) error) error {
	files := OpenFileChannel(owner)
	if err := bcgo.LoadHead(files, cache, network); err != nil {
		return err
	}
	block, err := bcgo.GetBlockContainingRecord(files.GetName(), cache, network, recordHash)
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

func GetPreview(previews bcgo.Channel, cache bcgo.Cache, network bcgo.Network, alias string, key *rsa.PrivateKey, recordHash []byte, callback func(*bcgo.BlockEntry, []byte, *Preview) error) error {
	return bcgo.Read(previews.GetName(), previews.GetHead(), nil, cache, network, alias, key, recordHash, func(entry *bcgo.BlockEntry, key, data []byte) error {
		// Unmarshal as Preview
		preview := &Preview{}
		if err := proto.Unmarshal(data, preview); err != nil {
			return err
		} else if err := callback(entry, key, preview); err != nil {
			return err
		}
		return nil
	})
}

func GetTag(tags bcgo.Channel, cache bcgo.Cache, network bcgo.Network, alias string, key *rsa.PrivateKey, recordHash []byte, callback func(*bcgo.BlockEntry, []byte, *Tag) error) error {
	return bcgo.Read(tags.GetName(), tags.GetHead(), nil, cache, network, alias, key, recordHash, func(entry *bcgo.BlockEntry, key, data []byte) error {
		// Unmarshal as Tag
		tag := &Tag{}
		if err := proto.Unmarshal(data, tag); err != nil {
			return err
		} else if err := callback(entry, key, tag); err != nil {
			return err
		}
		return nil
	})
}

func CreateRemoteMiningRequest(host, feature string, record *bcgo.Record) (*http.Request, error) {
	var buffer bytes.Buffer
	writer := bufio.NewWriter(&buffer)
	if err := bcgo.WriteDelimitedProtobuf(writer, record); err != nil {
		return nil, err
	}
	return http.NewRequest("POST", host+"/mining/"+feature, &buffer)
}

func ParseRemoteMiningResponse(response *http.Response) (*bcgo.Reference, error) {
	switch response.StatusCode {
	case http.StatusOK:
		defer response.Body.Close()
		reference := &bcgo.Reference{}
		if err := bcgo.ReadDelimitedProtobuf(bufio.NewReader(response.Body), reference); err != nil {
			return nil, err
		}
		return reference, nil
	default:
		return nil, errors.New("Response status: " + response.Status)
	}
}
