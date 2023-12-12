package main

import (
	"context"
	"os"

	"github.com/lixvyang/betxin.one/cmd/root"
	"github.com/lixvyang/betxin.one/internal/session"
)

var (
	version = "2.0.0"
)

func main() {
	ctx := context.Background()
	s := &session.Session{Version: version}
	ctx = session.With(ctx, s)

	expandedArgs := []string{}
	if len(os.Args) > 0 {
		expandedArgs = os.Args[1:]
	}

	rootCmd := root.NewCmdRoot(version)

	rootCmd.SetArgs(expandedArgs)
	if err := rootCmd.ExecuteContext(ctx); err != nil {
		rootCmd.PrintErrln("execute failed:", err)
		os.Exit(1)
	}
}
