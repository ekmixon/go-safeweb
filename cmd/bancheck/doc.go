// Copyright 2021 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// 	https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Package main contains the CLI used for detecting risky APIs.
// See https://pkg.go.dev/github.com/google/go-safeweb/safehttp#hdr-Restricting_Risky_APIs
// for a high level overview.
//
// Overview
//
// Bancheck is a program that implements a static analysis module and lets you run it
// on go packages of your choice. Under the hood it uses the go/analysis package
// https://pkg.go.dev/golang.org/x/tools/go/analysis which provides all the tools that
// are needed for static code analysis. The tool resolves fully qualified function
// and import names and checks them against a config file that defines banned tokens.
//
// Usage
//
// Apart from the standard https://pkg.go.dev/golang.org/x/tools/go/analysis#Analyzer flags
// the command requires a config flag where a list of config files should be provided.
//
//  $ ./bancheck
//  bannedAPI: Checks for usage of banned APIs
//
//  Usage: bannedAPI [-flag] [package]
//
//  Flags:
//  -V    print version and exit
//  -all
//  		no effect (deprecated)
//  -c int
//  		display offending line with this many lines of context (default -1)
//  -configs string
// 		Config files with banned APIs separated by a comma
//  -cpuprofile string
// 		write CPU profile to this file
//  -debug string
// 		debug flags, any subset of "fpstv"
//  -fix
// 		apply all suggested fixes
//  -flags
//  		print analyzer flags in JSON
//  -json
//  		emit JSON output
//  -memprofile string
//  		write memory profile to this file
//  -source
//  		no effect (deprecated)
//  -tags string
// 	 	no effect (deprecated)
//  -trace string
//  		write trace log to this file
//  -v    no effect (deprecated)
//
// Config
//
// Config lets you specify which APIs should be banned, explain why they are risky to use
// and allow a list of packages for which the check should be skipped.
// The structure of a config can be found in go-safeweb/cmd/bancheck/config/config.go.
//
// Note: It is possible to have colliding config files e.g. one config file bans an API
// but another one exempts it. The tool applies checks from each config file separately
// i.e. one warning will still be returned.
//
// Example config:
//  {
// 	"functions": [
// 		{
// 			"name": "fmt.Printf",
// 			"msg": "Banned by team A"
// 		}
// 	],
// 	"imports": [
// 		{
// 			"name": "fmt",
// 			"msg": "Banned by team A",
//			"exemptions": [
//				{
//					"justification": "#yolo",
//					"allowedPkg": "main"
//				}
//			]
// 		}
// 	]
//  }
//
// Example
//
// The example below shows a simple use case where "fmt" package and "fmt.Printf" function were banned
// by two separate teams.
//
// main.go
//  package main
//
//  import "fmt"
//
//  func main() {
//   fmt.Printf("Hello")
//  }
//
// config.json
//  {
//   "functions": [
// 	  {
// 	   "name": "fmt.Printf",
// 	   "msg": "Banned by team A"
// 	  }
// 	 ],
//   "imports": [
// 	  {
// 	   "name": "fmt",
// 	   "msg": "Banned by team B"
// 	  }
// 	 ],
//  }
//
// CLI usage
//  $ ./bancheck -configs config.json main.go
//  /go-safeweb/cmd/bancheck/test/main.go:3:8: Banned API found "fmt". Additional info: Banned by team B
//  /go-safeweb/cmd/bancheck/test/main.go:6:6: Banned API found "fmt.Printf". Additional info: Banned by team A
package main
