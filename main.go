package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/spf13/cobra"
)

var rootCmd = cobra.Command{
	Use:   "awsthingy",
	Short: "A handy CLI for managing AWS resources.",
}

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	rootCmd.AddCommand(ec2List())

	_ = rootCmd.ExecuteContext(ctx)
}
