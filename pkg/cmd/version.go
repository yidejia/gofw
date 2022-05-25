package cmd

import (
    "github.com/yidejia/gofw/pkg/console"

    "github.com/spf13/cobra"
)

// CmdVersion 输出框架版本命令
var CmdVersion = &cobra.Command{
    Use:   "version",
    Short:  "Show Gofw Framework version",
    Run: runVersion,
    Args:  cobra.ExactArgs(0),
}

func runVersion(cmd *cobra.Command, args []string) {
    console.Success("Gofw Framework 1.0.0")
}