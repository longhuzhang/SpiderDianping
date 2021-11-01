package main

import "embed"

//go:embed tmp/1aef4652.json
//go:embed tmp/4ac7e624.json
//go:embed tmp/5a5dc163.json
//go:embed tmp/91a5dbad.json
//go:embed tmp/8119f472.json
//go:embed tmp/44761ea4.json
//go:embed tmp/b80115ef.json
//go:embed tmp/ecba0360.json
//go:embed tmp/ed9f41cd.json
//go:embed tmp/ef6fa92a.json
//go:embed tmp/f7fec93e.json
//go:embed tmp/fd3e5bf5.json
//go:embed tmp/febad456.json
//go:embed tmp/0e5c4c69.json
//go:embed tmp/61478e36.json
//go:embed tmp/c75d3538.json
//go:embed tmp/ec1ec5f6.json
var jsonFile embed.FS

var FileStore = map[string]bool{
	"tmp/1aef4652.json": true,
	"tmp/4ac7e624.json": true,
	"tmp/5a5dc163.json": true,
	"tmp/91a5dbad.json": true,
	"tmp/8119f472.json": true,
	"tmp/44761ea4.json": true,
	"tmp/b80115ef.json": true,
	"tmp/ecba0360.json": true,
	"tmp/ed9f41cd.json": true,
	"tmp/ef6fa92a.json": true,
	"tmp/f7fec93e.json": true,
	"tmp/fd3e5bf5.json": true,
	"tmp/febad456.json": true,
	"tmp/0e5c4c69.json": true,
	"tmp/61478e36.json": true,
	"tmp/c75d3538.json": true,
	"tmp/ec1ec5f6.json": true,
}
