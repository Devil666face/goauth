package assets

import (
	"embed"
)

//go:embed templates/*
var Viewfs embed.FS

// go:embed templates/*
// var Staticfs embed.FS
