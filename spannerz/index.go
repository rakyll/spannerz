// Copyright 2020 Google Inc. All Rights Reserved.
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

import "html/template"

var indexTmpl, _ = template.New("index").Parse(`<!DOCTYPE html>
<html lang="en">
  <head>
    <meta charset="utf-8">
    <title>Cloud Spanner Query Analysis Plan</title>
  </head>
  <style>
	.container {
		width: 770px;
		margin: 50px auto 50px auto;
		font-family: Arial, Helvetica, sans-serif;
	}
	a {
		color: #1a73e8;
	}
	form { margin-bottom: 30px; }
	input, textarea {font-size: 20px;}
	#q { width: 740px; padding: 15px; }
	.button {
		background-color: #1a73e8;
		font-weight: bold;
		border: none;
		color: white;
		padding: 8px 25px;
		text-align: center;
		display: inline-block;
		font-size: 16px;
	}
	svg { 
		display:block;
		margin:auto; 
	}
	.docs {
		font-size: 14px;
		color: #666;
		margin-left: 10px;
	}
	.red {
		color: red;
	}
	.stats {
		width: 300px;
		float: left;
		line-height: 140%;
	}
  </style>
  <body>
	<div class="container">
	<h1>Google Cloud Spanner Query Analyzer</h1>  
	<form action="/" method="post">
		<textarea id="q" name="q" type="text" placeholder="Enter your query..." rows="3">{{.Query}}</textarea><br>
		<input class="button" type="submit" value="Analyze">
		<span class="docs">
		See Spanner's <a href="https://cloud.google.com/spanner/docs/query-execution-plans">Query Execution Plans</a> 
		and <a href="https://github.com/rakyll/spannerz">GitHub</a> for documentation.
		</span>
	</form>
	{{if .Error}}
		<span class="red">{{ .Error }}</span>
	{{end}}
	{{if .Stats}}
		<div class="stats">
		{{ range $key, $value := .Stats }}
   			<strong>{{ $key }}</strong>: {{ $value }}<br>
		{{ end }}
		</div>
	{{end}}
	{{if .Image}}
		{{ .Image }}
	{{end}}
	</div>
  </body>
</html>
`)

type IndexData struct {
	Query string
	Stats map[string]string
	Image template.HTML
	Error error
}
