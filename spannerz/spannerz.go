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

package spannerz

import (
	"context"
	"fmt"
	"html/template"
	"net/http"
	"os/exec"
	"time"

	"cloud.google.com/go/spanner"
	structpb "github.com/golang/protobuf/ptypes/struct"
	"google.golang.org/api/iterator"
	sppb "google.golang.org/genproto/googleapis/spanner/v1"
)

type Handler struct {
	Client *spanner.Client
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Check if dot is installed.
	_, err := exec.LookPath("dot")
	if err != nil {
		h.serveError(err, w, r)
		return
	}

	var query string
	var image string
	var queryStats map[string]string

	if r.Method == http.MethodPost {
		if err := r.ParseForm(); err != nil {
			h.serveError(err, w, r)
			return
		}
		query = r.FormValue("q")

		var err error
		stats, g, err := h.queryPlan(r.Context(), query)
		if err != nil {
			h.serveError(err, w, r)
			return
		}
		queryStats = stats
		image, err = g.SVG()
		if err != nil {
			h.serveError(err, w, r)
			return
		}
	}
	if err := indexTmpl.Execute(w, &IndexData{
		Query: query,
		Stats: queryStats,
		Image: template.HTML(image),
	}); err != nil {
		fmt.Fprintf(w, "Failed to render the page: %v", err)
	}
}

func (h *Handler) serveError(err error, w http.ResponseWriter, r *http.Request) {
	if err := indexTmpl.Execute(w, &IndexData{
		Error: err,
	}); err != nil {
		fmt.Fprintf(w, "Failed to render the page: %v", err)
	}
}

func (h *Handler) queryPlan(ctx context.Context, query string) (map[string]string, *Graph, error) {
	stmt := spanner.NewStatement(query)
	it := h.Client.Single().QueryWithStats(ctx, stmt)
	defer it.Stop()
	for {
		_, err := it.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, nil, err
		}
	}

	plan := it.QueryPlan
	stats := make(map[string]string, len(it.QueryStats))
	for k, v := range it.QueryStats {
		stats[k] = fmt.Sprintf("%v", v)
	}

	graph := NewGraph()
	for i, node := range plan.PlanNodes {
		if node.Kind == sppb.PlanNode_SCALAR {
			// Scalars are to pretty print references.
			continue
		}

		var latency, cpuTime time.Duration
		for k, v := range node.ExecutionStats.GetFields() {
			switch k {
			case "latency":
				latency = getDuration(v)
			case "cpu_time":
				cpuTime = getDuration(v)
			}
		}
		n := &Node{
			ID:      i,
			Name:    node.DisplayName,
			Attrs:   getAttrs(node),
			Latency: latency,
			CPUTime: cpuTime,
		}
		graph.AddNode(n)
	}
	return stats, graph, nil
}

func getAttrs(n *sppb.PlanNode) map[string]string {
	attrs := make(map[string]string)
	for k, v := range n.Metadata.GetFields() {
		// TODO(jbd): Set the actual value.
		var str string
		switch v.Kind.(type) {
		case *structpb.Value_StringValue:
			str = v.GetStringValue()
		case *structpb.Value_NullValue:
			str = "NULL"
		case *structpb.Value_BoolValue:
			str = fmt.Sprintf("%v", v)
		case *structpb.Value_NumberValue:
			str = fmt.Sprintf("%v", v)
		default:
			str = "unknown"
		}
		attrs[k] = str
	}
	return attrs
}

func getDuration(v *structpb.Value) time.Duration {
	var dur string
	switch st := v.GetKind().(type) {
	case *structpb.Value_StructValue:
		for k, v := range st.StructValue.GetFields() {
			switch k {
			case "total":
				dur = v.GetStringValue()
			case "unit":
				// TODO(jbd): Not sure if the unit is
				// always ms, handle and convert to Go unit.
			}
		}
	}
	tdur, _ := time.ParseDuration(dur + "ms")
	return tdur
}
