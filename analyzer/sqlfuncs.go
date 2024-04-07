package analyzer

import "golang.org/x/tools/go/ssa"

type SQLFuncCollection map[string]map[string]int

func NewSQLFuncsCollection() SQLFuncCollection {
	return SQLFuncCollection{
		"database/sql": {
			"ExecContext":     2,
			"QueryContext":    2,
			"QueryRowContext": 2,
			"PrepareContext":  2,
			"Prepare":         1,
			"Exec":            1,
			"Query":           1,
			"QueryRow":        1,
		},
	}
}

func (s SQLFuncCollection) Get(f *ssa.Function) (int, bool) {
	if f == nil || f.Pkg == nil || f.Pkg.Pkg == nil {
		return 0, false
	}

	pkgCol, ok := s[f.Pkg.Pkg.Path()]
	if !ok {
		return 0, false
	}

	num, ok := pkgCol[f.Name()]
	return num, ok
}
