package cmd

import (
	"context"
	"fmt"
	"github.com/spf13/cobra"
	"go-extractor/services/extractor"
	"log"
	"os"
	"os/signal"
	"syscall"
)

var rootCmd = &cobra.Command{
	Use:   "app",
	Short: "",
	Run: func(cmd *cobra.Command, args []string) {
		//
	},
}

var runExtractorCmd = &cobra.Command{
	Use:   "extractor:run",
	Short: "Asana extractor",
	Run: func(cmd *cobra.Command, args []string) {
		log.Println("Starting extractor..")
		extractor.Start(cmd.Context())
	},
}

func Execute() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	rootCmd.SetContext(ctx)
	rootCmd.AddCommand(runExtractorCmd)

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
