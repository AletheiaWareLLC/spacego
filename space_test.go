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
		"Remove": {
			delta: &spacego.Delta{
				Remove: uint64(len([]byte("foo"))),
			},
			initial: "foobar",
			want:    "bar",
		},
		"RemoveAll": {
			delta: &spacego.Delta{
				Remove: uint64(len([]byte("foobar"))),
			},
			initial: "foobar",
			want:    "",
		},
		"Append": {
			delta: &spacego.Delta{
				Offset: 6,
				Add:    []byte("blah"),
			},
			initial: "foobar",
			want:    "foobarblah",
		},
		"Insert": {
			delta: &spacego.Delta{
				Offset: 3,
				Add:    []byte("blah"),
			},
			initial: "foobar",
			want:    "fooblahbar",
		},
		"Replace": {
			delta: &spacego.Delta{
				Offset: 3,
				Remove: uint64(len([]byte("bar"))),
				Add:    []byte("blah"),
			},
			initial: "foobar",
			want:    "fooblah",
		},
	} {
		t.Run(name, func(t *testing.T) {
			got := string(spacego.ApplyDelta(tt.delta, []byte(tt.initial)))
			if got != tt.want {
				t.Fatalf("Incorrect buffer; expected '%s', got '%s'", tt.want, got)
			}
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
			want:    []*spacego.Delta{},
		},
		"Single": {
			initial: "foobar",
			want: []*spacego.Delta{
				&spacego.Delta{
					Add: []byte("foobar"),
				},
			},
		},
		"Double": {
			initial: "foobarfoobar",
			want: []*spacego.Delta{
				&spacego.Delta{
					Add: []byte("foobarfoob"),
				},
				&spacego.Delta{
					Offset: 10,
					Add:    []byte("ar"),
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
			if len(got) != len(tt.want) {
				t.Fatalf("Wrong number of deltas; expected '%d', got '%d'", len(tt.want), len(got))
			}
			for i, w := range tt.want {
				g := got[i]
				if g.Offset != w.Offset {
					t.Fatalf("Incorrect offset; expected '%d', got '%d'", w.Offset, g.Offset)
				}
				if g.Remove != w.Remove {
					t.Fatalf("Incorrect remove; expected '%d', got '%d'", w.Remove, g.Remove)
				}
				if len(g.Add) != len(w.Add) {
					t.Fatalf("Incorrect add length; expected '%d', got '%d'", len(w.Add), len(g.Add))
				}
				if string(g.Add) != string(w.Add) {
					t.Fatalf("Incorrect add; expected '%s', got '%s'", string(w.Add), string(g.Add))
				}
			}
		})
	}
}
