package utils

import (
	"context"
	"encoding/base64"
	"fmt"

	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/security/keyvault/azsecrets"
)

func getAzureClient() (*azidentity.DefaultAzureCredential, error) {
	creds, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		return nil, err
	}

	return creds, nil
}

func GetKeyvaultSecret(resource_group, vault_name, secret_name, secret_version string) (string, error) {
	// if secret_version == "" {
	// 	secret_version = "latest"
	// }
	// Create a KeyVault client with the Azure credentials
	creds, err := getAzureClient()
	if err != nil {
		return "", err
	}
	vault_url := fmt.Sprintf("https://%s.vault.azure.net", vault_name)

	// Fetch the secret from the key vault
	secretClient, err := azsecrets.NewClient(vault_url, creds, nil)
	if err != nil {
		return "", err
	}

	secretResponse, err := secretClient.GetSecret(context.TODO(), secret_name, secret_version, nil)
	if err != nil {
		return "", err
	}

	// Return the secret value as a string
	fmt.Println(*secretResponse.Value)
	return *secretResponse.Value, nil
}

func Base64Encode(data []byte) string {
	encodedSecret := base64.StdEncoding.EncodeToString(data)
	return encodedSecret
}
