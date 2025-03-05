/*
Copyright (c) 2022 Cisco Systems, Inc. and others.  All rights reserved.
*/
package cmd

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
	"path"
	"runtime"
	rtpl "sqbu-github.cisco.com/tievens/render-tpl/pkg/render-tpl"
	"strconv"
	"strings"
	"time"
)

var cfgFile string
var cfgDebug bool

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "render-tpl",
	Short: "Render golang template",
	Long:  `Render golang template`,

	Run: func(cmd *cobra.Command, args []string) {
		if cfgDebug {
			log.SetLevel(log.DebugLevel)
		} else {
			log.SetLevel(log.InfoLevel)
		}

		log.SetOutput(os.Stdout)

		// Validate args
		if len(rtpl.TemplateFile) <= 0 {
			log.Fatalf("Missing template filename")
		}

		if len(rtpl.ValueFiles) <= 0 {
			log.Fatalf("Missing value file")
		}

		// Entrypoint
		rtpl.Run()
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.render-tpl.yaml)")
	rootCmd.PersistentFlags().BoolVar(&cfgDebug, "debug", false,
		"Debug logging (default is false)")

	rootCmd.PersistentFlags().StringSliceVarP(&rtpl.ValueFiles, "values", "v",
		nil, "REQUIRED: Value filename. Can repeat arg or define using a comma delimited list.")

	rootCmd.PersistentFlags().StringVarP(&rtpl.TemplateFile, "template", "t",
		"", "REQUIRED: Template filename.")

	rootCmd.PersistentFlags().BoolVarP(&rtpl.UseStdout, "stdout", "s", false,
		"Write to STDOUT instead of template filename sans .tpl (default is false)")

	rootCmd.MarkPersistentFlagRequired("values")
	rootCmd.MarkPersistentFlagRequired("template")

	log.SetReportCaller(true)
	log.SetFormatter(&log.TextFormatter{
		CallerPrettyfier:       logCbPrettyfier,
		ForceColors:            false,
		FullTimestamp:          true,
		DisableLevelTruncation: false,
		TimestampFormat:        time.RFC3339,
	})
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		// Search config in home directory with name ".render-tpl" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigType("yaml")
		viper.SetConfigName(".render-tpl")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
	}

}

func logCbPrettyfier(frame *runtime.Frame) (function string, file string) {
	var sb strings.Builder
	sb.WriteString(path.Base(frame.Function))
	sb.WriteByte('[')
	sb.WriteString(strconv.Itoa(frame.Line))
	sb.WriteByte(']')

	return sb.String(), path.Base(frame.File)
}
