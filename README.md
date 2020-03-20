# spannerz

spannerz adds an HTTP handler to your binary
to report query plans from [Google Cloud Spanner](http://cloud.google.com/spanner) clients.
You can use spannerz as a standalone binary too.

**NOTE**: You need [Graphviz](https://www.graphviz.org/) installed for visualization features.

## Goals

* Allow users to investigate their client setup in production without having to redeploy new versions.
* Allow users to run a visualizer outside of Google Cloud Console.

## Usage

Standalone binary:

```
$ go get -u github.com/rakyll/spannerz
$ spannerz -db projects/PROJECT/instances/SPANNER_INSTANCE/databases/SPANNER_DB
```

HTTP handler:

``` go
import (
    "cloud.google.com/go/spanner"
    "github.com/rakyll/spannerz/spannerz"
)

client, err := spanner.NewClient(ctx, "projects/PROJECT/instances/SPANNER_INSTANCE/databases/SPANNER_DB")
if err != nil {
    log.Fatalf("Cannot create Spanner client: %v", err)
}
http.Handle("/spannerz", &spannerz.Handler{
    Client: client,
})
log.Fatal(http.ListenAndServe(":9090", nil))
```

![Screenshot](https://i.imgur.com/06XVWjh.png)


## Roadmap
* Allow running only the query planner without executing the query.
* Support different optimizer versions.
* Support partition queries.
* Support read/write transactions, we currently support read-only ones.