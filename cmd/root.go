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
	"time"
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
		if len(args) < 1 {
			log.Fatal("Please indicate the extraction process interval")
		}

		interval, err := time.ParseDuration(args[0])
		if err != nil {
			log.Fatalf("invalid time interval format indicated %s", err)
		}

		extractor.Start(cmd.Context(), interval)
	},
}

func Execute() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	rootCmd.AddCommand(runExtractorCmd)
	rootCmd.SetContext(ctx)

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
