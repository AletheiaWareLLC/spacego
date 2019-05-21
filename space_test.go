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

package spacego_test

import (
	"github.com/AletheiaWareLLC/bcgo"
	"github.com/AletheiaWareLLC/spacego"
	"github.com/AletheiaWareLLC/testinggo"
	"testing"
)

func TestCreateRemoteMiningRequest(t *testing.T) {
	t.Run("URL", func(t *testing.T) {
		record := &bcgo.Record{
			Timestamp: 1234,
		}
		request, err := spacego.CreateRemoteMiningRequest("foo.bar", "baz", record)
		testinggo.AssertNoError(t, err)
		expected := "foo.bar/mining/baz"
		if request.URL.String() != expected {
			t.Fatalf("Incorrect URL; expected '%s', got '%s'", expected, request.URL.String())
		}
	})
}

func TestParseRemoteMiningResponse(t *testing.T) {

}
