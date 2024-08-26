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
	githubUrl          string
	githubOwner        string
	githubRepo         string
	secretName         string
	githubEnv          string
	resourceGroup      string
	vaultName          string
	vaultSecretName    string
	vaultSecretVersion string
	encodeSecret       bool
)

// githubCmd represents the github command
var GithubCmd = &cobra.Command{
	Use:   "github",
	Short: "A command to create github secrets from azure keyvault secrets",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		secretValue, err := utils.GetKeyvaultSecret(resourceGroup, vaultName, vaultSecretName, vaultSecretVersion)
		if err != nil {
			log.Fatal(err)
		}
		if encodeSecret {
			secretValue = utils.Base64Encode([]byte(secretValue))
			fmt.Println(secretValue)
		}
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
	GithubCmd.PersistentFlags().StringVar(&githubOwner, "owner", "", "Github Repository Owner")
	GithubCmd.PersistentFlags().StringVar(&githubRepo, "repository", "", "Github Repository name")
	GithubCmd.PersistentFlags().StringVar(&secretName, "secret-name", "", "Github Secret Name")
	GithubCmd.PersistentFlags().StringVar(&githubEnv, "github-env", "", "Github Environment Name")
	GithubCmd.PersistentFlags().StringVar(&resourceGroup, "resource-group", "", "Key Vault Resource Group Name")
	GithubCmd.PersistentFlags().StringVar(&vaultName, "keyvault-name", "", "Key Vault Name")
	GithubCmd.PersistentFlags().StringVar(&vaultSecretName, "keyvault-secret-name", "", "Key Vault Secret Name")
	GithubCmd.PersistentFlags().StringVar(&vaultSecretVersion, "keyvault-secret-version", "", "Key Vault Secret Version")
	GithubCmd.PersistentFlags().BoolVar(&encodeSecret, "encode-secret", false, "Encode Secret Value")
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}
}
