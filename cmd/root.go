package cmd

import (
	"fmt"
	"os"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	_ "github.com/go-sql-driver/mysql"

	"github.com/s-capybara/db-schema-markdown/lib"
)

var cfgFile string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "db-schema-markdown",
	Short: "Generates database definition for Markdown from existing table",
	Long: `"db-schema-markdown" is a CLI tool to generate a database definition document for Markdown
from an existing database table.

Positional arguments specify columns to show.`,
	Run: Root,
}

func Root(cmd *cobra.Command, args []string) {
	dbsm.Root(Auth(), Table(), Columns(cmd, args))
}

func Auth() dbsm.Auth {
	return dbsm.Auth{
		Username: viper.GetString("username"),
		Password: viper.GetString("password"),
		Database: viper.GetString("database")}
}

func Table() string {
	return viper.GetString("table")
}

func Columns(cmd *cobra.Command, args []string) []string {
	// TODO: It's better to get columns by flags.
	if viper.GetBool("full") {
		// Full columns.
		return []string{"Field", "Type", "Collation", "Null", "Key", "Default", "Extra", "Privileges", "Comment"}
	} else if len(args) == 0 {
		// Standard columns.
		return []string{"Field", "Type", "Null", "Default", "Comment"}
	} else {
		// Specified columns.
		return args
	}
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default $HOME/.db-schema-markdown.yml)")
	rootCmd.PersistentFlags().BoolP("full", "f", false, "shows all columns if true")

	rootCmd.Flags().StringP("username", "u", "root", "username for database")
	rootCmd.Flags().StringP("password", "p", "", "password for database")
	rootCmd.Flags().StringP("database", "D", "", "database name")
	rootCmd.Flags().StringP("table", "t", "", "table name")

	rootCmd.MarkFlagRequired("database")
	rootCmd.MarkFlagRequired("table")

	viper.BindPFlag("username", rootCmd.Flags().Lookup("username"))
	viper.BindPFlag("password", rootCmd.Flags().Lookup("password"))
	viper.BindPFlag("database", rootCmd.Flags().Lookup("database"))
	viper.BindPFlag("table", rootCmd.Flags().Lookup("table"))
	viper.BindPFlag("full", rootCmd.PersistentFlags().Lookup("full"))
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// Search config in home directory with name ".db-schema-markdown" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(".db-schema-markdown")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		// fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}
