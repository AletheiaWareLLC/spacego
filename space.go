/*
 * Copyright 2019-2020 Aletheia Ware LLC
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
	"crypto/rsa"
	"github.com/AletheiaWareLLC/bcgo"
	"github.com/golang/protobuf/proto"
	"io"
)

const (
	SPACE              = "S P A C E"
	SPACE_HOUR         = "Space-Hour"
	SPACE_DAY          = "Space-Day"
	SPACE_YEAR         = "Space-Year"
	SPACE_CHARGE       = "Space-Charge"
	SPACE_INVOICE      = "Space-Invoice"
	SPACE_REGISTRAR    = "Space-Registrar"
	SPACE_REGISTRATION = "Space-Registration"
	SPACE_SUBSCRIPTION = "Space-Subscription"
	SPACE_USAGE_RECORD = "Space-Usage-Record"

	SPACE_PREFIX         = "Space-"
	SPACE_PREFIX_DELTA   = "Space-Delta-"
	SPACE_PREFIX_META    = "Space-Meta-"
	SPACE_PREFIX_PREVIEW = "Space-Preview-"
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

	MAX_SIZE_BYTES = bcgo.MAX_PAYLOAD_SIZE_BYTES - 1024 // 10Mb-1Kb (for delta protobuf stuff)
)

type DeltaCallback func(entry *bcgo.BlockEntry, delta *Delta) error

type MetaCallback func(entry *bcgo.BlockEntry, meta *Meta) error

type PreviewCallback func(entry *bcgo.BlockEntry, preview *Preview) error

type TagCallback func(entry *bcgo.BlockEntry, tag *Tag) error

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
	return bcgo.OpenPoWChannel(SPACE_CHARGE, bcgo.THRESHOLD_Z)
}

func OpenInvoiceChannel() *bcgo.Channel {
	return bcgo.OpenPoWChannel(SPACE_INVOICE, bcgo.THRESHOLD_Z)
}

func OpenRegistrarChannel() *bcgo.Channel {
	return bcgo.OpenPoWChannel(SPACE_REGISTRAR, bcgo.THRESHOLD_Z)
}

func OpenRegistrationChannel() *bcgo.Channel {
	return bcgo.OpenPoWChannel(SPACE_REGISTRATION, bcgo.THRESHOLD_Z)
}

func OpenSubscriptionChannel() *bcgo.Channel {
	return bcgo.OpenPoWChannel(SPACE_SUBSCRIPTION, bcgo.THRESHOLD_Z)
}

func OpenUsageRecordChannel() *bcgo.Channel {
	return bcgo.OpenPoWChannel(SPACE_USAGE_RECORD, bcgo.THRESHOLD_Z)
}

func GetDeltaChannelName(metaId string) string {
	return SPACE_PREFIX_DELTA + metaId
}

func GetMetaChannelName(alias string) string {
	return SPACE_PREFIX_META + alias
}

func GetPreviewChannelName(metaId string) string {
	return SPACE_PREFIX_PREVIEW + metaId
}

func GetTagChannelName(metaId string) string {
	return SPACE_PREFIX_TAG + metaId
}

func OpenDeltaChannel(metaId string) *bcgo.Channel {
	return bcgo.OpenPoWChannel(GetDeltaChannelName(metaId), bcgo.THRESHOLD_Z)
}

func OpenMetaChannel(alias string) *bcgo.Channel {
	return bcgo.OpenPoWChannel(GetMetaChannelName(alias), bcgo.THRESHOLD_Z)
}

func OpenPreviewChannel(metaId string) *bcgo.Channel {
	return bcgo.OpenPoWChannel(GetPreviewChannelName(metaId), bcgo.THRESHOLD_Z)
}

func OpenTagChannel(metaId string) *bcgo.Channel {
	return bcgo.OpenPoWChannel(GetTagChannelName(metaId), bcgo.THRESHOLD_Z)
}

func ApplyDelta(delta *Delta, input []byte) []byte {
	count := 0
	length := uint64(len(input))
	output := make([]byte, length-delta.Remove+uint64(len(delta.Add)))
	if delta.Offset <= length {
		count += copy(output, input[:delta.Offset])
	}
	count += copy(output[count:], delta.Add)
	index := delta.Offset + delta.Remove
	if index < length {
		copy(output[count:], input[index:])
	}
	return output
}

func CreateDeltaRecords(reader io.Reader, max uint64, callback func(*Delta) error) error {
	buffer := make([]byte, max)
	var size uint64
	for {
		count, err := reader.Read(buffer)
		if err != nil {
			if err == io.EOF {
				// Ignore EOFs
				break
			} else {
				return err
			}
		}
		add := make([]byte, count)
		copy(add, buffer[:count])
		delta := &Delta{
			Offset: size,
			Add:    add,
		}
		if err := callback(delta); err != nil {
			return err
		}
		size += uint64(count)
	}
	return nil
}

func GetThreshold(channel string) uint64 {
	switch channel {
	case SPACE_HOUR:
		return bcgo.THRESHOLD_PERIOD_HOUR
	case SPACE_DAY:
		return bcgo.THRESHOLD_PERIOD_DAY
	case SPACE_YEAR:
		return bcgo.THRESHOLD_PERIOD_YEAR
	default:
		return bcgo.THRESHOLD_Z
	}
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
	return
}

func GetDelta(deltas *bcgo.Channel, cache bcgo.Cache, network bcgo.Network, alias string, key *rsa.PrivateKey, recordHash []byte, callback DeltaCallback) error {
	return bcgo.Read(deltas.Name, deltas.Head, nil, cache, network, alias, key, recordHash, func(entry *bcgo.BlockEntry, key, data []byte) error {
		// Unmarshal as Delta
		d := &Delta{}
		err := proto.Unmarshal(data, d)
		if err != nil {
			return err
		}
		return callback(entry, d)
	})
}

func GetMeta(metas *bcgo.Channel, cache bcgo.Cache, network bcgo.Network, alias string, key *rsa.PrivateKey, recordHash []byte, callback MetaCallback) error {
	return bcgo.Read(metas.Name, metas.Head, nil, cache, network, alias, key, recordHash, func(entry *bcgo.BlockEntry, key, data []byte) error {
		// Unmarshal as Meta
		meta := &Meta{}
		if err := proto.Unmarshal(data, meta); err != nil {
			return err
		}
		return callback(entry, meta)
	})
}

func GetPreview(previews *bcgo.Channel, cache bcgo.Cache, network bcgo.Network, alias string, key *rsa.PrivateKey, recordHash []byte, callback PreviewCallback) error {
	return bcgo.Read(previews.Name, previews.Head, nil, cache, network, alias, key, recordHash, func(entry *bcgo.BlockEntry, key, data []byte) error {
		// Unmarshal as Preview
		preview := &Preview{}
		if err := proto.Unmarshal(data, preview); err != nil {
			return err
		}
		return callback(entry, preview)
	})
}

func GetTag(tags *bcgo.Channel, cache bcgo.Cache, network bcgo.Network, alias string, key *rsa.PrivateKey, recordHash []byte, callback TagCallback) error {
	return bcgo.Read(tags.Name, tags.Head, nil, cache, network, alias, key, recordHash, func(entry *bcgo.BlockEntry, key, data []byte) error {
		// Unmarshal as Tag
		tag := &Tag{}
		if err := proto.Unmarshal(data, tag); err != nil {
			return err
		}
		return callback(entry, tag)
	})
}

func GetMinimumRegistrars() int {
	if bcgo.IsLive() {
		return 3
	}
	return 1
}

func IterateDeltas(node *bcgo.Node, deltas *bcgo.Channel, callback DeltaCallback) error {
	// Iterate through chain chronologically
	return bcgo.IterateChronologically(deltas.Name, deltas.Head, nil, node.Cache, node.Network, func(hash []byte, block *bcgo.Block) error {
		for _, entry := range block.Entry {
			for _, access := range entry.Record.Access {
				if node.Alias == access.Alias {
					if err := bcgo.DecryptRecord(entry, access, node.Key, func(entry *bcgo.BlockEntry, key []byte, data []byte) error {
						// Unmarshal as Delta
						d := &Delta{}
						if err := proto.Unmarshal(data, d); err != nil {
							return err
						}
						return callback(entry, d)
					}); err != nil {
						return err
					}
				}
			}
		}
		return nil
	})
}
