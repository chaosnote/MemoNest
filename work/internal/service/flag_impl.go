package service

import (
	"flag"
)

type FlagImpl struct {
	AllowDemo bool
}

func NewFlagImpl() (*FlagImpl, error) {
	allow_demo := flag.Bool("allow_demo", false, "是否開啟 Demo Router")

	flag.Parse()

	return &FlagImpl{
		AllowDemo: *allow_demo,
	}, nil
}
