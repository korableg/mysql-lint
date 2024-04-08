package analyzer

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUnit_Analyzer(t *testing.T) {
	var (
		err   = NewAnalyzer().Analyze(context.Background(), "../test/database_sql")
		aErrs AnalyzeErrors
	)
	require.ErrorAs(t,
		err,
		&aErrs)
	assert.Len(t, aErrs, 2)

	var wantErrs = map[string]struct{}{
		`sql.go:79:32: issue: line 1 column 68 near "%s GROUP BY level ORDER BY level DESC LIMIT 15 AND K = 4" `: {},
		`sql.go:58:21: issue: line 1 column 49 near "%s)" `:                                                      {},
	}

	for _, e := range aErrs {
		_, ok := wantErrs[e.Error()]
		assert.True(t, ok)
	}

	assert.ErrorContains(t, err, `sql.go:79:32`)
	assert.ErrorContains(t, err, `sql.go:58:21`)
}
