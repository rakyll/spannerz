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
