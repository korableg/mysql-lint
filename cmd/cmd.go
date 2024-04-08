package cmd

import (
	"context"
	"errors"
	"fmt"
	"os"
	"runtime/debug"

	"github.com/spf13/cobra"

	"github.com/korableg/mysql-lint/analyzer"
)

const (
	AppName = "mysql-lint"

	FlagDir = "dir"
)

func NewCommand() *cobra.Command {
	var (
		version = "(devel)"
		bi, ok  = debug.ReadBuildInfo()
	)

	if ok {
		version = bi.Main.Version
	}

	c := &cobra.Command{
		Use:     AppName,
		Version: version,
		Short:   "MySQL query linter",
		Run: func(cmd *cobra.Command, args []string) {
			var dir, _ = cmd.Flags().GetString(FlagDir)
			analyze(cmd.Context(), dir)
		},
	}

	c.Flags().String(FlagDir, ".", "Path to the GO code directory")

	return c
}

func analyze(ctx context.Context, dir string) {
	fmt.Println("MySQL query linter: Start checking...")

	a := analyzer.NewAnalyzer()
	err := a.Analyze(ctx, dir)
	if err == nil {
		fmt.Println("Checking finished successfully! 🎉")
		os.Exit(0)
	}

	var aErrs analyzer.AnalyzeErrors
	if !errors.As(err, &aErrs) {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Printf("Found %d errors:\n", len(aErrs))

	for _, e := range aErrs {
		fmt.Println(e)
	}

	os.Exit(1)
}
