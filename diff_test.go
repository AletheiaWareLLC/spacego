/*
 * Copyright 2021 Aletheia Ware LLC
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
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCompact(t *testing.T) {
	for name, tt := range map[string]struct {
		deltas, expected []*spacego.Delta
	}{
		"empty": {},
		"single": {
			deltas: []*spacego.Delta{
				&spacego.Delta{},
			},
			expected: []*spacego.Delta{
				&spacego.Delta{},
			},
		},
		"consecutive": {
			deltas: []*spacego.Delta{
				&spacego.Delta{},
				&spacego.Delta{
					Offset: 1,
				},
			},
			expected: []*spacego.Delta{
				&spacego.Delta{},
				&spacego.Delta{
					Offset: 1,
				},
			},
		},
		"delete_delete": {
			deltas: []*spacego.Delta{
				&spacego.Delta{
					Delete: 1,
				},
				&spacego.Delta{
					Offset: 1,
					Delete: 1,
				},
			},
			expected: []*spacego.Delta{
				&spacego.Delta{
					Delete: 2,
				},
			},
		},
		"insert_insert": {
			deltas: []*spacego.Delta{
				&spacego.Delta{
					Insert: []byte("a"),
				},
				&spacego.Delta{
					Insert: []byte("b"),
				},
			},
			expected: []*spacego.Delta{
				&spacego.Delta{
					Insert: []byte("ab"),
				},
			},
		},
		"delete_insert": {
			deltas: []*spacego.Delta{
				&spacego.Delta{
					Delete: 1,
				},
				&spacego.Delta{
					Insert: []byte("a"),
				},
			},
			expected: []*spacego.Delta{
				&spacego.Delta{
					Delete: 1,
					Insert: []byte("a"),
				},
			},
		},
		"insert_delete": {
			deltas: []*spacego.Delta{
				&spacego.Delta{
					Insert: []byte("a"),
				},
				&spacego.Delta{
					Delete: 1,
				},
			},
			expected: []*spacego.Delta{
				&spacego.Delta{
					Delete: 1,
					Insert: []byte("a"),
				},
			},
		},
	} {
		t.Run(name, func(t *testing.T) {
			assert.Equal(t, tt.expected, spacego.Compact(tt.deltas))
		})
	}
}

func TestDifference(t *testing.T) {
	for name, tt := range map[string]struct {
		a, b     string
		expected []*spacego.Delta
	}{
		"empty": {},
		"equal": {
			a: "foobar",
			b: "foobar",
		},
		"greeting": {
			a: "Hello World",
			b: "Hi Earth",
			expected: []*spacego.Delta{
				&spacego.Delta{
					Offset: 1,
					Delete: 4,
					Insert: []byte("i"),
				},
				&spacego.Delta{
					Offset: 3,
					Delete: 2,
					Insert: []byte("Ea"),
				},
				&spacego.Delta{
					Offset: 6,
					Delete: 2,
					Insert: []byte("th"),
				},
			},
		},
		"insert_prefix": {
			a: "bar",
			b: "foobar",
			expected: []*spacego.Delta{
				&spacego.Delta{
					Insert: []byte("foo"),
				},
			},
		},
		"insert_infix": {
			a: "foar",
			b: "foobar",
			expected: []*spacego.Delta{
				&spacego.Delta{
					Offset: 2,
					Insert: []byte("ob"),
				},
			},
		},
		"insert_suffix": {
			a: "foo",
			b: "foobar",
			expected: []*spacego.Delta{
				&spacego.Delta{
					Offset: 3,
					Insert: []byte("bar"),
				},
			},
		},
		"delete_prefix": {
			a: "foobar",
			b: "bar",
			expected: []*spacego.Delta{
				&spacego.Delta{
					Delete: 3,
				},
			},
		},
		"delete_infix": {
			a: "foobar",
			b: "foar",
			expected: []*spacego.Delta{
				&spacego.Delta{
					Offset: 2,
					Delete: 2,
				},
			},
		},
		"delete_suffix": {
			a: "foobar",
			b: "foo",
			expected: []*spacego.Delta{
				&spacego.Delta{
					Offset: 3,
					Delete: 3,
				},
			},
		},
		"swap": {
			a: "foobar",
			b: "barfoo",
			expected: []*spacego.Delta{
				&spacego.Delta{
					Delete: 3,
				},
				&spacego.Delta{
					Offset: 3,
					Insert: []byte("foo"),
				},
			},
		},
		"delete_vowels": {
			a: "foobar",
			b: "fbr",
			expected: []*spacego.Delta{
				&spacego.Delta{
					Offset: 1,
					Delete: 2,
				},
				&spacego.Delta{
					Offset: 2,
					Delete: 1,
				},
			},
		},
		"delete_consonants": {
			a: "foobar",
			b: "ooa",
			expected: []*spacego.Delta{
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
		},
		"insert_vowels": {
			a: "fbr",
			b: "foobar",
			expected: []*spacego.Delta{
				&spacego.Delta{
					Offset: 1,
					Insert: []byte("oo"),
				},
				&spacego.Delta{
					Offset: 4,
					Insert: []byte("a"),
				},
			},
		},
		"insert_consonants": {
			a: "ooa",
			b: "foobar",
			expected: []*spacego.Delta{
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
		},
		"replace": {
			a: "foo",
			b: "bar",
			expected: []*spacego.Delta{
				&spacego.Delta{
					Delete: 3,
					Insert: []byte("bar"),
				},
			},
		},
		"reverse": {
			a: "foobar",
			b: "raboof",
			expected: []*spacego.Delta{
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
		},
	} {
		t.Run(name, func(t *testing.T) {
			assert.Equal(t, tt.expected, spacego.Difference([]byte(tt.a), []byte(tt.b)))
		})
	}
}
