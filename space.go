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

package spacego

import (
	"bufio"
	"bytes"
	"crypto/rsa"
	"errors"
	"github.com/AletheiaWareLLC/bcgo"
	"github.com/AletheiaWareLLC/cryptogo"
	"github.com/golang/protobuf/proto"
	"net/http"
	"strings"
)

const (
	SPACE              = "S P A C E"
	SPACE_HOUR         = "Space-Hour"
	SPACE_DAY          = "Space-Day"
	SPACE_YEAR         = "Space-Year"
	SPACE_CHARGE       = "Space-Charge"
	SPACE_INVOICE      = "Space-Invoice"
	SPACE_MINER        = "Space-Miner"
	SPACE_REGISTRAR    = "Space-Registrar"
	SPACE_REGISTRATION = "Space-Registration"
	SPACE_SUBSCRIPTION = "Space-Subscription"
	SPACE_USAGE_RECORD = "Space-Usage-Record"

	SPACE_PREFIX         = "Space-"
	SPACE_PREFIX_FILE    = "Space-File-"
	SPACE_PREFIX_META    = "Space-Meta-"
	SPACE_PREFIX_PREVIEW = "Space-Preview-"
	SPACE_PREFIX_SHARE   = "Space-Share-"
	SPACE_PREFIX_TAG     = "Space-Tag-"

	MIME_TYPE_UNKNOWN    = "?/?"
	MIME_TYPE_IMAGE_JPEG = "image/jpeg"
	MIME_TYPE_IMAGE_JPG  = "image/jpg"
	MIME_TYPE_IMAGE_GIF  = "image/gif"
	MIME_TYPE_IMAGE_PNG  = "image/png"
	MIME_TYPE_IMAGE_WEBP = "image/webp"
	MIME_TYPE_TEXT_PLAIN = "text/plain"
	MIME_TYPE_PDF        = "application/pdf"
	MIME_TYPE_PROTOBUF   = "application/x-protobuf"
	MIME_TYPE_VIDEO_MPEG = "video/mpeg"
	MIME_TYPE_AUDIO_MPEG = "audio/mpeg"

	MIME_TYPE_IMAGE_DEFAULT = "image/jpeg"
	MIME_TYPE_VIDEO_DEFAULT = "video/mpeg"
	MIME_TYPE_AUDIO_DEFAULT = "audio/mpeg"

	PREVIEW_IMAGE_SIZE  = 128
	PREVIEW_TEXT_LENGTH = 64
)

func GetSpaceHosts() []string {
	if bcgo.IsLive() {
		return []string{
			"space-nyc.aletheiaware.com",
			"space-sfo.aletheiaware.com",
		}
	}
	return []string{
		"test-space.aletheiaware.com",
	}
}

func OpenHourChannel() *bcgo.Channel {
	return bcgo.OpenPoWChannel(SPACE_HOUR, bcgo.THRESHOLD_PERIOD_HOUR)
}

func OpenDayChannel() *bcgo.Channel {
	return bcgo.OpenPoWChannel(SPACE_DAY, bcgo.THRESHOLD_PERIOD_DAY)
}

func OpenYearChannel() *bcgo.Channel {
	return bcgo.OpenPoWChannel(SPACE_YEAR, bcgo.THRESHOLD_PERIOD_YEAR)
}

func OpenChargeChannel() *bcgo.Channel {
	return bcgo.OpenPoWChannel(SPACE_CHARGE, bcgo.THRESHOLD_G)
}

func OpenInvoiceChannel() *bcgo.Channel {
	return bcgo.OpenPoWChannel(SPACE_INVOICE, bcgo.THRESHOLD_G)
}

func OpenMinerChannel() *bcgo.Channel {
	return bcgo.OpenPoWChannel(SPACE_MINER, bcgo.THRESHOLD_G)
}

func OpenRegistrarChannel() *bcgo.Channel {
	return bcgo.OpenPoWChannel(SPACE_REGISTRAR, bcgo.THRESHOLD_G)
}

func OpenRegistrationChannel() *bcgo.Channel {
	return bcgo.OpenPoWChannel(SPACE_REGISTRATION, bcgo.THRESHOLD_G)
}

func OpenSubscriptionChannel() *bcgo.Channel {
	return bcgo.OpenPoWChannel(SPACE_SUBSCRIPTION, bcgo.THRESHOLD_G)
}

func OpenUsageRecordChannel() *bcgo.Channel {
	return bcgo.OpenPoWChannel(SPACE_USAGE_RECORD, bcgo.THRESHOLD_G)
}

func GetFileChannelName(alias string) string {
	return SPACE_PREFIX_FILE + alias
}

func GetMetaChannelName(alias string) string {
	return SPACE_PREFIX_META + alias
}

func GetShareChannelName(alias string) string {
	return SPACE_PREFIX_SHARE + alias
}

func GetPreviewChannelName(metaId string) string {
	return SPACE_PREFIX_PREVIEW + metaId
}

func GetTagChannelName(metaId string) string {
	return SPACE_PREFIX_TAG + metaId
}

func OpenFileChannel(alias string) *bcgo.Channel {
	return bcgo.OpenPoWChannel(GetFileChannelName(alias), bcgo.THRESHOLD_I)
}

func OpenMetaChannel(alias string) *bcgo.Channel {
	return bcgo.OpenPoWChannel(GetMetaChannelName(alias), bcgo.THRESHOLD_G)
}

func OpenShareChannel(alias string) *bcgo.Channel {
	return bcgo.OpenPoWChannel(GetShareChannelName(alias), bcgo.THRESHOLD_G)
}

func OpenPreviewChannel(metaId string) *bcgo.Channel {
	return bcgo.OpenPoWChannel(GetPreviewChannelName(metaId), bcgo.THRESHOLD_G)
}

func OpenTagChannel(metaId string) *bcgo.Channel {
	return bcgo.OpenPoWChannel(GetTagChannelName(metaId), bcgo.THRESHOLD_G)
}

func GetThreshold(channel string) uint64 {
	if strings.HasPrefix(channel, SPACE_PREFIX_FILE) {
		return bcgo.THRESHOLD_I
	}
	switch channel {
	case SPACE_HOUR:
		return bcgo.THRESHOLD_PERIOD_HOUR
	case SPACE_DAY:
		return bcgo.THRESHOLD_PERIOD_DAY
	case SPACE_YEAR:
		return bcgo.THRESHOLD_PERIOD_YEAR
	default:
		return bcgo.THRESHOLD_G
	}
}

func GetMiner(miners *bcgo.Channel, cache bcgo.Cache, network bcgo.Network, alias string) (*Miner, error) {
	var miner *Miner
	if err := bcgo.Read(miners.Name, miners.Head, nil, cache, network, "", nil, nil, func(entry *bcgo.BlockEntry, key, data []byte) error {
		// Unmarshal as Miner
		m := &Miner{}
		err := proto.Unmarshal(data, m)
		if err != nil {
			return err
		}
		if m.Merchant.Alias == alias {
			miner = m
			return bcgo.StopIterationError{}
		}
		return nil
	}); err != nil {
		switch err.(type) {
		case bcgo.StopIterationError:
			// Do nothing
			break
		default:
			return nil, err
		}
	}
	return miner, nil
}

func GetMiners(miners *bcgo.Channel, cache bcgo.Cache, network bcgo.Network) (ms []*Miner) {
	bcgo.Read(miners.Name, miners.Head, nil, cache, network, "", nil, nil, func(entry *bcgo.BlockEntry, key, data []byte) error {
		// Unmarshal as Miner
		m := &Miner{}
		err := proto.Unmarshal(data, m)
		if err != nil {
			return err
		}
		ms = append(ms, m)
		return nil
	})
	return ms
}

func GetRegistrar(registrars *bcgo.Channel, cache bcgo.Cache, network bcgo.Network, alias string) (*Registrar, error) {
	var registrar *Registrar
	if err := bcgo.Read(registrars.Name, registrars.Head, nil, cache, network, "", nil, nil, func(entry *bcgo.BlockEntry, key, data []byte) error {
		// Unmarshal as Registrar
		r := &Registrar{}
		err := proto.Unmarshal(data, r)
		if err != nil {
			return err
		}
		if r.Merchant.Alias == alias {
			registrar = r
			return bcgo.StopIterationError{}
		}
		return nil
	}); err != nil {
		switch err.(type) {
		case bcgo.StopIterationError:
			// Do nothing
			break
		default:
			return nil, err
		}
	}
	return registrar, nil
}

func GetRegistrars(registrars *bcgo.Channel, cache bcgo.Cache, network bcgo.Network) (rs []*Registrar) {
	bcgo.Read(registrars.Name, registrars.Head, nil, cache, network, "", nil, nil, func(entry *bcgo.BlockEntry, key, data []byte) error {
		// Unmarshal as Registrar
		r := &Registrar{}
		err := proto.Unmarshal(data, r)
		if err != nil {
			return err
		}
		rs = append(rs, r)
		return nil
	})
	return rs
}

func GetFile(files *bcgo.Channel, cache bcgo.Cache, network bcgo.Network, alias string, key *rsa.PrivateKey, recordHash []byte, callback func(*bcgo.BlockEntry, []byte, []byte) error) error {
	return bcgo.Read(files.Name, files.Head, nil, cache, network, alias, key, recordHash, callback)
}

func GetMeta(metas *bcgo.Channel, cache bcgo.Cache, network bcgo.Network, alias string, key *rsa.PrivateKey, recordHash []byte, callback func(*bcgo.BlockEntry, []byte, *Meta) error) error {
	return bcgo.Read(metas.Name, metas.Head, nil, cache, network, alias, key, recordHash, func(entry *bcgo.BlockEntry, key, data []byte) error {
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

func GetShare(shares *bcgo.Channel, cache bcgo.Cache, network bcgo.Network, alias string, key *rsa.PrivateKey, recordHash []byte, callback func(*bcgo.BlockEntry, []byte, *Share) error) error {
	return bcgo.Read(shares.Name, shares.Head, nil, cache, network, alias, key, recordHash, func(entry *bcgo.BlockEntry, key, data []byte) error {
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

func GetSharedMeta(metas *bcgo.Channel, cache bcgo.Cache, network bcgo.Network, recordHash, key []byte, callback func(*bcgo.BlockEntry, *Meta) error) error {
	block, err := bcgo.GetBlockContainingRecord(metas.Name, cache, network, recordHash)
	if err != nil {
		return err
	}
	for _, entry := range block.Entry {
		if bytes.Equal(recordHash, entry.RecordHash) {
			data, err := cryptogo.DecryptPayload(entry.Record.EncryptionAlgorithm, entry.Record.Payload, key)
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

func GetSharedFile(files *bcgo.Channel, cache bcgo.Cache, network bcgo.Network, recordHash, key []byte, callback func(*bcgo.BlockEntry, []byte) error) error {
	block, err := bcgo.GetBlockContainingRecord(files.Name, cache, network, recordHash)
	if err != nil {
		return err
	}
	for _, entry := range block.Entry {
		if bytes.Equal(recordHash, entry.RecordHash) {
			data, err := cryptogo.DecryptPayload(entry.Record.EncryptionAlgorithm, entry.Record.Payload, key)
			if err != nil {
				return err
			}
			return callback(entry, data)
		}
	}
	return nil
}

func GetPreview(previews *bcgo.Channel, cache bcgo.Cache, network bcgo.Network, alias string, key *rsa.PrivateKey, recordHash []byte, callback func(*bcgo.BlockEntry, []byte, *Preview) error) error {
	return bcgo.Read(previews.Name, previews.Head, nil, cache, network, alias, key, recordHash, func(entry *bcgo.BlockEntry, key, data []byte) error {
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

func GetTag(tags *bcgo.Channel, cache bcgo.Cache, network bcgo.Network, alias string, key *rsa.PrivateKey, recordHash []byte, callback func(*bcgo.BlockEntry, []byte, *Tag) error) error {
	return bcgo.Read(tags.Name, tags.Head, nil, cache, network, alias, key, recordHash, func(entry *bcgo.BlockEntry, key, data []byte) error {
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

func GetMinimumRegistrars() int {
	if bcgo.IsLive() {
		return 3
	}
	return 1
}
