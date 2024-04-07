package analyzer

import (
	"fmt"
	"go/token"
	"path/filepath"
	"strings"
)

const (
	ErrTypeWarning ErrorType = "warning"
	ErrTypeIssue   ErrorType = "issue"
)

type (
	ErrorType string

	AnalyzeError struct {
		Err     string
		RelPath string
		Type    ErrorType
	}

	AnalyzeErrors []*AnalyzeError
)

func newAnalyzeError(basePath string, err error, pos token.Position, typ ErrorType) *AnalyzeError {
	relPath, fErr := filepath.Rel(basePath, pos.Filename)
	if fErr != nil {
		relPath = basePath
	}
	return &AnalyzeError{
		RelPath: fmt.Sprintf("%s:%d:%d", relPath, pos.Line, pos.Column),
		Err:     fmt.Sprintf("%s: %v", typ, err),
		Type:    typ,
	}
}

func (a *AnalyzeError) Error() string {
	return fmt.Sprintf("%s: %s", a.RelPath, a.Err)
}

func (a AnalyzeErrors) Error() string {
	var sb strings.Builder

	for _, e := range a {
		sb.WriteString(e.Error())
		sb.WriteString("\n")
	}

	return sb.String()
}
