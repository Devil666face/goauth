package assets

import (
	"embed"
)

//go:embed templates/*
var Viewfs embed.FS

//go:embed static/*
var Staticfs embed.FS

var (
	Mediadir  string = "media"
	Staticdir string = "static"
)
