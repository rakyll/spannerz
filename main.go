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

// Program spannerz starts a web server to visualize
// the query plan for a Google Cloud Spanner query.
package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"google.golang.org/api/option"

	"cloud.google.com/go/spanner"
	"github.com/rakyll/spannerz/spannerz"
)

var (
	database string
	httpv    string
)

func main() {
	ctx := context.Background()
	flag.StringVar(&database, "db", "", "")
	flag.StringVar(&httpv, "http", "localhost:9090", "") // starts and HTTP server at the target.
	flag.Usage = usage
	flag.Parse()

	if database == "" {
		usage()
		os.Exit(1)
	}

	client, err := spanner.NewClient(ctx, database, option.WithUserAgent("spannerz/0.1"))
	if err != nil {
		log.Fatalf("Cannot create Spanner client: %v", err)
	}
	fmt.Printf("Starting the server at http://%v\n", httpv)
	log.Fatal(http.ListenAndServe(httpv, &spannerz.Handler{
		Client: client,
	}))
}

func usage() {
	fmt.Println(`spannerz [options...]

Options:
-db    Google Cloud Spanner database string, e.g.
       projects/PROJECT/instances/SPANNER_INSTANCE/databases/SPANNER_DB
-http  Host and port to start the visualization server, 
       the default is "localhost:9090".`)
}
