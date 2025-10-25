package config

import "flag"

var (
	Channel   = flag.String("c", "", "target channel")
	DirName   = flag.String("d", ".", "target dir")
	GetFile   = flag.String("g", "", "get target file")
	IsListing = flag.Bool("l", false, "list all files")
)
