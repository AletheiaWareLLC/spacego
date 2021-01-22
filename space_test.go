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

package spacego_test

import (
	"aletheiaware.com/spacego"
	"aletheiaware.com/testinggo"
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

func TestApplyDelta(t *testing.T) {
	for name, tt := range map[string]struct {
		delta   *spacego.Delta
		initial string
		want    string
	}{
		"Empty": {
			delta:   &spacego.Delta{},
			initial: "foobar",
			want:    "foobar",
		},
		"Delete": {
			delta: &spacego.Delta{
				Delete: uint64(len([]byte("foo"))),
			},
			initial: "foobar",
			want:    "bar",
		},
		"DeleteAll": {
			delta: &spacego.Delta{
				Delete: uint64(len([]byte("foobar"))),
			},
			initial: "foobar",
			want:    "",
		},
		"Append": {
			delta: &spacego.Delta{
				Offset: 6,
				Insert: []byte("blah"),
			},
			initial: "foobar",
			want:    "foobarblah",
		},
		"Insert": {
			delta: &spacego.Delta{
				Offset: 3,
				Insert: []byte("blah"),
			},
			initial: "foobar",
			want:    "fooblahbar",
		},
		"Replace": {
			delta: &spacego.Delta{
				Offset: 3,
				Delete: uint64(len([]byte("bar"))),
				Insert: []byte("blah"),
			},
			initial: "foobar",
			want:    "fooblah",
		},
	} {
		t.Run(name, func(t *testing.T) {
			got := string(spacego.ApplyDelta(tt.delta, []byte(tt.initial)))
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestCreateDeltas(t *testing.T) {
	for name, tt := range map[string]struct {
		initial string
		want    []*spacego.Delta
	}{
		"Empty": {
			initial: "",
		},
		"Single": {
			initial: "foobar",
			want: []*spacego.Delta{
				&spacego.Delta{
					Insert: []byte("foobar"),
				},
			},
		},
		"Double": {
			initial: "foobarfoobar",
			want: []*spacego.Delta{
				&spacego.Delta{
					Insert: []byte("foobarfoob"),
				},
				&spacego.Delta{
					Offset: 10,
					Insert: []byte("ar"),
				},
			},
		},
	} {
		t.Run(name, func(t *testing.T) {
			var got []*spacego.Delta
			testinggo.AssertNoError(t, spacego.CreateDeltas(strings.NewReader(tt.initial), 10, func(d *spacego.Delta) error {
				got = append(got, d)
				return nil
			}))
			assert.Equal(t, tt.want, got)
		})
	}
}
