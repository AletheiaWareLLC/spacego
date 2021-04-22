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
		given    string
		deltas   []*spacego.Delta
		expected string
	}{
		"empty": {},
		"equal": {
			given:    "foobar",
			expected: "foobar",
		},
		"insert_prefix": {
			given: "bar",
			deltas: []*spacego.Delta{
				&spacego.Delta{
					Insert: []byte("foo"),
				},
			},
			expected: "foobar",
		},
		"insert_infix": {
			given: "foar",
			deltas: []*spacego.Delta{
				&spacego.Delta{
					Offset: 2,
					Insert: []byte("ob"),
				},
			},
			expected: "foobar",
		},
		"insert_suffix": {
			given: "foo",
			deltas: []*spacego.Delta{
				&spacego.Delta{
					Offset: 3,
					Insert: []byte("bar"),
				},
			},
			expected: "foobar",
		},
		"delete_all": {
			given: "foobar",
			deltas: []*spacego.Delta{
				&spacego.Delta{
					Delete: 6,
				},
			},
		},
		"delete_prefix": {
			given: "foobar",
			deltas: []*spacego.Delta{
				&spacego.Delta{
					Delete: 3,
				},
			},
			expected: "bar",
		},
		"delete_infix": {
			given: "foobar",
			deltas: []*spacego.Delta{
				&spacego.Delta{
					Offset: 2,
					Delete: 2,
				},
			},
			expected: "foar",
		},
		"delete_suffix": {
			given: "foobar",
			deltas: []*spacego.Delta{
				&spacego.Delta{
					Offset: 3,
					Delete: 3,
				},
			},
			expected: "foo",
		},
		"swap": {
			given: "foobar",
			deltas: []*spacego.Delta{
				&spacego.Delta{
					Insert: []byte("bar"),
				},
				&spacego.Delta{
					Offset: 6,
					Delete: 3,
				},
			},
			expected: "barfoo",
		},
		"delete_vowels": {
			given: "foobar",
			deltas: []*spacego.Delta{
				&spacego.Delta{
					Offset: 1,
					Delete: 2,
				},
				&spacego.Delta{
					Offset: 2,
					Delete: 1,
				},
			},
			expected: "fbr",
		},
		"delete_consonants": {
			given: "foobar",
			deltas: []*spacego.Delta{
				&spacego.Delta{
					Delete: 1,
				},
				&spacego.Delta{
					Offset: 2,
					Delete: 1,
				},
				&spacego.Delta{
					Offset: 3,
					Delete: 1,
				},
			},
			expected: "ooa",
		},
		"insert_vowels": {
			given: "fbr",
			deltas: []*spacego.Delta{
				&spacego.Delta{
					Offset: 1,
					Insert: []byte("oo"),
				},
				&spacego.Delta{
					Offset: 4,
					Insert: []byte("a"),
				},
			},
			expected: "foobar",
		},
		"insert_consonants": {
			given: "ooa",
			deltas: []*spacego.Delta{
				&spacego.Delta{
					Insert: []byte("f"),
				},
				&spacego.Delta{
					Offset: 3,
					Insert: []byte("b"),
				},
				&spacego.Delta{
					Offset: 5,
					Insert: []byte("r"),
				},
			},
			expected: "foobar",
		},
		"replace": {
			given: "foo",
			deltas: []*spacego.Delta{
				&spacego.Delta{
					Delete: 3,
					Insert: []byte("bar"),
				},
			},
			expected: "bar",
		},
		"reverse": {
			given: "foobar",
			deltas: []*spacego.Delta{
				&spacego.Delta{
					Delete: 1,
					Insert: []byte("rab"),
				},
				&spacego.Delta{
					Offset: 5,
					Delete: 3,
					Insert: []byte("f"),
				},
			},
			expected: "raboof",
		},
	} {
		t.Run(name, func(t *testing.T) {
			buffer := []byte(tt.given)
			for _, d := range tt.deltas {
				buffer = spacego.ApplyDelta(d, buffer)
			}
			assert.Equal(t, tt.expected, string(buffer))
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
