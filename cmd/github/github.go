/*
Copyright © 2024 NAME HERE shivamverma182@gmail.com
*/
package github

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"log"
	"net/url"
	"os"

	"github.com/google/go-github/v63/github"
	"github.com/joho/godotenv"
	"github.com/spf13/cobra"
	"golang.org/x/crypto/nacl/box"
)

var (
	githubUrl     string
	githubOwner   string
	githubRepo    string
	secretName    string
	secretValue   string = "1234"
	githubEnv     string
	actionService github.ActionsService
)

func getClient(token, url string) (*github.Client, error) {
	if url == "" {
		return github.NewClient(nil).WithAuthToken(token), nil
	}
	client, err := github.NewClient(nil).WithAuthToken(token).WithEnterpriseURLs(url, url)
	if err != nil {
		return nil, err
	}
	return client, nil
}

func encryptSecret(pk, secret string) (string, error) {
	publicKeyBytes, err := base64.StdEncoding.DecodeString(pk)
	if err != nil {
		return "", err
	}
	var pk32 [32]byte
	copy(pk32[:], publicKeyBytes)
	secretBytes := secret

	out := make([]byte, 0, len(secretBytes)+box.Overhead+len(pk32))
	enc, err := box.SealAnonymous(out, []byte(secretBytes), &pk32, rand.Reader)
	if err != nil {
		return "", err
	}
	encEnc := base64.StdEncoding.EncodeToString(enc)
	return encEnc, nil
}

func getRepoPublicKeyDetails(owner, repo, token, url string) (keyId, keyData string, err error) {
	client, err := getClient(token, url)
	if err != nil {
		return "", "", err
	}
	key, _, err := client.Actions.GetRepoPublicKey(context.TODO(), owner, repo)
	if err != nil {
		return "", "", err
	}
	return key.GetKeyID(), key.GetKey(), nil
}

func getEnvPublicKeyDetails(token, url, env string, repoId int) (keyId, keyData string, err error) {
	client, err := getClient(token, url)
	if err != nil {
		return "", "", err
	}
	key, _, err := client.Actions.GetEnvPublicKey(context.TODO(), repoId, env)
	if err != nil {
		return "", "", err
	}
	return key.GetKeyID(), key.GetKey(), nil
}

func SetSecret(github_url, owner, repo, secret, env, token string) (int, error) {

	client, err := getClient(token, github_url)
	if err != nil {
		return 0, err
	}
	if env == "" {
		keyId, keyData, err := getRepoPublicKeyDetails(owner, repo, token, github_url)
		if err != nil {
			return 0, err
		}
		encryptedValue, err := encryptSecret(keyData, secret)
		if err != nil {
			return 0, err
		}
		resp, err := client.Actions.CreateOrUpdateRepoSecret(context.TODO(), owner, repo, &github.EncryptedSecret{
			Name:           secretName,
			KeyID:          keyId,
			EncryptedValue: encryptedValue,
		})
		if err != nil {
			return 0, err
		}
		return resp.StatusCode, nil
	}
	escapedEnv := url.PathEscape(env)
	repo_o, _, err := client.Repositories.Get(context.TODO(), owner, repo)
	if err != nil {
		return 0, err
	}
	keyId, keyData, err := getEnvPublicKeyDetails(token, github_url, env, int(repo_o.GetID()))
	encryptedValue, err := encryptSecret(keyData, secret)
	resp, err := client.Actions.CreateOrUpdateEnvSecret(context.TODO(), int(repo_o.GetID()), escapedEnv, &github.EncryptedSecret{
		Name:           secretName,
		KeyID:          keyId,
		EncryptedValue: encryptedValue,
	})
	if err != nil {
		return 0, err
	}
	return resp.StatusCode, nil
}

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
		_, err := SetSecret(githubUrl, githubOwner, githubRepo, secretValue, githubEnv, os.Getenv("GITHUB_TOKEN"))
		if err != nil {
			log.Fatal(err)
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
