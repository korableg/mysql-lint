package analyzer

import (
	"context"
	"errors"
	"fmt"
	"go/constant"
	"go/token"
	"os"
	"path/filepath"

	"github.com/pingcap/tidb/pkg/parser"
	_ "github.com/pingcap/tidb/pkg/parser/test_driver"
	"golang.org/x/tools/go/callgraph/cha"
	"golang.org/x/tools/go/packages"
	"golang.org/x/tools/go/ssa"
	"golang.org/x/tools/go/ssa/ssautil"
)

const (
	pkgMode = packages.NeedName | packages.NeedFiles | packages.NeedImports | packages.NeedDeps |
		packages.NeedTypes | packages.NeedSyntax | packages.NeedTypesInfo
)

type (
	Query struct {
		Query string
		Pos   token.Position
	}

	Analyzer struct {
		fSet     *token.FileSet
		noLint   *NoLint
		sqlFuncs SQLFuncCollection
	}
)

func NewAnalyzer() *Analyzer {
	var fSet = token.NewFileSet()
	return &Analyzer{
		fSet:     fSet,
		noLint:   NewNolint(fSet),
		sqlFuncs: NewSQLFuncsCollection(),
	}
}

func (a *Analyzer) Analyze(ctx context.Context, root string) (err error) {
	root, err = filepath.Abs(root)
	if err != nil {
		return fmt.Errorf("calculate file absolute path: %w", err)
	}

	info, err := os.Lstat(root)
	if err != nil {
		return fmt.Errorf("get file info: %w", err)
	}

	if !info.IsDir() {
		return errors.New("file is not a dir")
	}

	queries, err := a.getSQLQueries(ctx, root)
	if err != nil {
		return fmt.Errorf("get sql queries: %w", err)
	}

	var errs AnalyzeErrors
	for _, q := range queries {
		_, pWarns, pErr := parser.New().ParseSQL(q.Query)
		for _, w := range pWarns {
			errs = append(errs, newAnalyzeError(root, w, q.Pos, ErrTypeWarning))
		}

		if pErr != nil {
			errs = append(errs, newAnalyzeError(root, pErr, q.Pos, ErrTypeIssue))
		}
	}

	if len(errs) == 0 {
		return nil
	}

	return errs
}

func (a *Analyzer) getSQLQueries(ctx context.Context, root string) ([]Query, error) {
	pkgs, err := packages.Load(
		&packages.Config{
			Mode:    pkgMode,
			Context: ctx,
			Dir:     root,
			Fset:    a.fSet,
		},
		filepath.Join(root, "..."))
	if err != nil {
		return nil, fmt.Errorf("load packages: %w", err)
	}

	prog, _ := ssautil.Packages(pkgs, 0)
	prog.Build()

	var (
		cg      = cha.CallGraph(prog)
		queries []Query
	)

	for _, n := range cg.Nodes {
		paramNum, ok := a.sqlFuncs.Get(n.Func)
		if !ok {
			continue
		}

		for _, edge := range n.In {
			queryString := queryFromSSAValue(edge.Site.Common().Args[paramNum])
			if queryString == "" {
				continue
			}

			var (
				pos           = a.fSet.Position(edge.Site.Pos())
				ignored, iErr = a.noLint.Ignored(pos.Filename, pos.Line)
			)
			if iErr != nil {
				return nil, err
			}

			if ignored {
				continue
			}

			queries = append(queries, Query{
				Query: queryString,
				Pos:   pos,
			})
		}
	}

	return queries, nil
}

func queryFromSSAValue(v ssa.Value) string {
	switch typedV := v.(type) {
	case *ssa.Const:
		return constant.StringVal(typedV.Value)
	case *ssa.BinOp:
		if typedV.Op != token.ADD {
			return ""
		}
		return queryFromSSAValue(typedV.X) + queryFromSSAValue(typedV.Y)
	}

	return ""
}
