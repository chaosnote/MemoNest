package service

import (
	"flag"

	"go.uber.org/zap"
)

type FlagImpl struct {
	AllowDemo bool
}

func NewFlagImpl(logger *zap.Logger) (*FlagImpl, error) {
	allow_demo := flag.Bool("allow_demo", false, "是否開啟 Demo Router")

	flag.Parse()

	return &FlagImpl{
		AllowDemo: *allow_demo,
	}, nil
}
