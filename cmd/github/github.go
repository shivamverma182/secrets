/*
Copyright © 2024 NAME HERE shivamverma182@gmail.com
*/
package github

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/shivamverma182/gh-secrets/utils"
	"github.com/spf13/cobra"
)

var (
	githubUrl   string
	githubOwner string
	githubRepo  string
	secretName  string
	secretValue string = "1234"
	githubEnv   string
)

// githubCmd represents the github command
var GithubCmd = &cobra.Command{
	Use:   "github",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(githubEnv) == 0 {
			_, err := utils.SetRepoSecret(githubUrl, githubOwner, githubRepo, secretName, secretValue, os.Getenv("GITHUB_TOKEN"))
			if err != nil {
				log.Fatal(err)
			}
		} else {
			_, err := utils.SetEnvSecret(githubUrl, githubOwner, githubRepo, secretName, secretValue, githubEnv, os.Getenv("GITHUB_TOKEN"))
			if err != nil {
				log.Fatal(err)
			}
		}
		fmt.Printf("Secret %s has been created.\n", secretName)
	},
}

func init() {
	GithubCmd.PersistentFlags().StringVar(&githubUrl, "url", "", "github server url")
	GithubCmd.PersistentFlags().StringVar(&githubOwner, "owner", "shivamverma182", "Github Repository Owner")
	GithubCmd.PersistentFlags().StringVar(&githubRepo, "repository", "secrets", "Github Repository name")
	GithubCmd.PersistentFlags().StringVar(&secretName, "secret-name", "api_key", "Github Secret Name")
	GithubCmd.PersistentFlags().StringVar(&githubEnv, "github-env", "", "Github Environment Name")
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}
}
