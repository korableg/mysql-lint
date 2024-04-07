package analyzer

import (
	"fmt"
	"go/parser"
	"go/token"
	"regexp"
	"sync"
)

type NoLint struct {
	fSet *token.FileSet
	ptrn *regexp.Regexp

	cache map[string]map[int]bool
	mu    sync.Mutex
}

func NewNolint(fSet *token.FileSet) *NoLint {
	return &NoLint{
		fSet:  fSet,
		ptrn:  regexp.MustCompile(`nolint *:[ ,a-zA-Z0-9]*mysql`),
		cache: make(map[string]map[int]bool),
	}
}

func (n *NoLint) Ignored(filename string, line int) (bool, error) {
	if ignored, ok := n.getFromCache(filename, line); ok {
		return ignored, nil
	}

	return n.getFromFile(filename, line)
}

func (n *NoLint) getFromCache(filename string, line int) (bool, bool) {
	n.mu.Lock()
	defer n.mu.Unlock()

	var cacheLine = n.cache[filename]
	if cacheLine != nil {
		return cacheLine[line], true
	}

	return false, false
}

func (n *NoLint) getFromFile(filename string, line int) (bool, error) {
	fi, err := parser.ParseFile(n.fSet, filename, nil, parser.ParseComments)
	if err != nil {
		return false, fmt.Errorf("parse comments: %w", err)
	}

	var (
		cacheLine = make(map[int]bool)
		ignored   bool
	)
	for _, c := range fi.Comments {
		for _, l := range c.List {
			if !n.ptrn.MatchString(l.Text) {
				continue
			}
			ll := n.fSet.Position(l.Pos()).Line
			cacheLine[ll] = true
			if ll == line {
				ignored = true
			}
		}
	}

	n.mu.Lock()
	n.cache[filename] = cacheLine
	n.mu.Unlock()

	return ignored, nil
}
