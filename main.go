package main

import (
	"github.com/spf13/cobra"

	"github.com/korableg/mysql-lint/cmd"
)

func main() {
	cobra.CheckErr(cmd.NewCommand().Execute())
}
