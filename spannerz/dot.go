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
	"fmt"
	"os/exec"
	"strings"
	"time"
)

type Graph struct {
	nodes []*Node
	links map[int]int
}

type Node struct {
	ID      int
	Name    string
	Latency time.Duration
	CPUTime time.Duration
	Attrs   map[string]string
}

func (n *Node) internalName() string {
	return fmt.Sprintf("n%d", n.ID)
}

func NewGraph() *Graph {
	return &Graph{
		nodes: make([]*Node, 0),
		links: make(map[int]int),
	}
}

func (g *Graph) AddNode(n *Node) {
	g.nodes = append(g.nodes, n)
}

func (g *Graph) LinkByIndex(from, to int) {
	g.links[from] = g.links[to]
}

func (g *Graph) String() string {
	b := &strings.Builder{}
	fmt.Fprintf(b, `digraph weighted {
graph [
rankdir = "TB"
];
`)
	for _, n := range g.nodes {
		fmt.Fprint(b, generateNode(n))
	}
	for i := 0; i < len(g.nodes)-1; i++ {
		fmt.Fprintf(b, "n%v -> n%v\n", i, i+1)
	}
	fmt.Fprintf(b, `}`)
	return b.String()
}

func (g *Graph) SVG() (string, error) {
	cmd := exec.Command("dot", "-Tsvg")
	cmd.Stdin = strings.NewReader(g.String())
	out, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("graphviz failed: %s", out)
	}
	return string(out), nil
}

func generateNode(n *Node) string {
	str := n.internalName() + " [\n"
	str += `label=<` + "\n" + generateName(n) + "\n>\n"
	str += "\n" + `shape = "box"`
	str += fmt.Sprintf("\nweight=%d", n.Latency.Microseconds())
	str += "\n];\n"
	return str
}

func generateName(n *Node) string {
	str := "<b>" + n.Name + "</b>"
	for k, v := range n.Attrs {
		str += "<br/>" + escape(k) + "=" + v
	}
	str += "<br/><br/>(Latency=" + n.Latency.String() + ")"
	return str
}

func escape(v string) string {
	v = strings.ReplaceAll(v, " ", "_")
	v = strings.ToLower(v)
	return v
}
