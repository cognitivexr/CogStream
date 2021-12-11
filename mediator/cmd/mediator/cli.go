package main

import (
	"path/filepath"
	"strings"
)

type AbsPathList []string

func (l *AbsPathList) String() string {
	return strings.Join(*l, ",")
}

func (l *AbsPathList) Set(value string) (err error) {
	value, err = filepath.Abs(value)
	if err != nil {
		return
	}
	*l = append(*l, value)
	return
}
