// Copyright 2020 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package spannerz_test

import (
	"context"
	"log"
	"net/http"

	"cloud.google.com/go/spanner"
	"github.com/rakyll/spannerz/spannerz"
)

var ctx = context.Background()

func Example() {
	client, err := spanner.NewClient(ctx, "projects/PROJECT/instances/SPANNER_INSTANCE/databases/SPANNER_DB")
	if err != nil {
		log.Fatalf("Cannot create Spanner client: %v", err)
	}
	http.Handle("/spannerz", &spannerz.Handler{
		Client: client,
	})
	log.Fatal(http.ListenAndServe(":9090", nil))
}
