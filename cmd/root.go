package cmd

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/kdisneur/changelog/pkg/changelog"
	"github.com/kdisneur/changelog/pkg/configuration"
)

var overrideConfigPath string
var configurationFile configuration.File
var configurationCommands configuration.Command

var rootCmd = &cobra.Command{
	Use:   "changelog [flags] <commit-reference> <new-version-name>",
	Short: "Generate a Changelog based on a Git history",
	Long:  "Read every commit, and fetch the bug tracker (e.g. GitHub pull request) description for every commits in the Git History",
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) == 2 {
			return nil
		}

		return fmt.Errorf("please check the arguments. expected 2, received %d\nArguments: %s", len(args), strings.Join(args, ", "))
	},
	Run: func(cmd *cobra.Command, args []string) {
		configurationCommands.From = args[0]
		configurationCommands.VersionName = args[1]
		configurationCommands.Date = time.Now()

		conf, err := configuration.Validate(configurationFile, configurationCommands)
		if err != nil {
			Exit(err.Error())
		}

		formattedChangelog, err := changelog.BuildChangelog(conf)
		if err != nil {
			Exit(err.Error())
		}

		fmt.Println(formattedChangelog)
	},
}

func Exit(message string) {
	fmt.Fprintln(os.Stderr, message)
	os.Exit(1)
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(loadConfigurationFile)

	defaultConfigurationPath, err := configuration.DefaultFilePath()
	if err != nil {
		Exit(err.Error())
	}

	rootCmd.PersistentFlags().StringVar(&overrideConfigPath, "config", "", fmt.Sprintf("config file (default is %s)", defaultConfigurationPath))
	rootCmd.Flags().StringVarP(&configurationCommands.RepositoryName, "repository", "r", "", "name of the GitHub repository (e.g. kdisneur/changelog)")
	rootCmd.Flags().StringVarP(&configurationCommands.RepositoryLocalPath, "change-dir", "C", ".", "path to the local repository path (e.g. ~/Workspace/kdisneur/changelog)")
	rootCmd.Flags().StringVarP(&configurationCommands.To, "branch", "b", "", `name of the base branch (default "master")`)
	rootCmd.Flags().StringVarP(&configurationCommands.MergeStrategy, "strategy", "", "", `commit history followed merge strategy (one of "squash" or "merge") (default "squash")`)
}

func loadConfigurationFile() {
	viper.SetConfigType("toml")

	if overrideConfigPath != "" {
		viper.SetConfigFile(overrideConfigPath)
	} else {
		defaultConfigurationFolder, err := configuration.DefaultFolderPath()
		if err != nil {
			Exit(err.Error())
		}

		viper.AddConfigPath(defaultConfigurationFolder)
		viper.SetConfigName(configuration.DefaultFileName())
	}

	viper.AutomaticEnv()

	err := viper.ReadInConfig()
	if err != nil {
		Exit(err.Error())
	}

	if viper.Unmarshal(&configurationFile) != nil {
		Exit("can't parse configuration file")
	}
}
