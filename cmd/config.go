package cmd

import (
	"errors"
	"fmt"
	"os"
	"strconv"

	"github.com/manifoldco/promptui"
	"github.com/spf13/viper"
)

func writeConfig() {
	//id := "AshleshaId"
	//secret := "AshleshaAPIsecret"
	//pat := "persdoihaofhw0tw345y8l"
	apiKey, apiSecret, pat, err := ConfigSecrets()
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
	accountId, workspaceId, teamId, err := ConfigDefaults()
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
	if apiKey == "" && apiSecret == "" && pat == "" {
		fmt.Println("Empty configurations, nothing to write, exiting...")
	} else {
		vp := viper.New()
		vp.SetConfigName("bmConfig")
		vp.SetConfigType("yaml")
		vp.AddConfigPath("$HOME")
		vp.Set("id", apiKey)
		vp.Set("secret", apiSecret)
		vp.Set("pat", pat)
		vp.Set("accountId", accountId)
		vp.Set("workspaceId", workspaceId)
		vp.Set("teamId", teamId)
		if err := vp.SafeWriteConfig(); err != nil {
			if _, ok := err.(viper.ConfigFileAlreadyExistsError); ok {
				fmt.Printf("Config file already exists, please edit the file directly or remove it from $HOME\n")
			} else {
				fmt.Printf("Error writing config file: %v\n", err)
			}
		} else {
			fmt.Printf("Config file written successfully\n")
		}
	}
}

// PromptForInput prompts the user for input without validation.
func PromptForInput(label string) (string, error) {
	prompt := promptui.Prompt{
		Label:       label,
		HideEntered: true,
	}
	result, err := prompt.Run()
	if err != nil {
		return "", fmt.Errorf("prompt failed: %w", err)
	}
	return result, nil
}

// ConfigPrompt prompts the user to enter API key, API secret, and PAT.
func ConfigSecrets() (string, string, string, error) {
	apiKey, err := PromptForInput("Enter the API key")
	if err != nil {
		return "", "", "", err
	}
	apiSecret, err := PromptForInput("Enter the API secret")
	if err != nil {
		return "", "", "", err
	}
	pat, err := PromptForInput("Enter the PAT for API Monitoring")
	if err != nil {
		return "", "", "", err
	}
	return apiKey, apiSecret, pat, nil
}

func ConfigDefaults() (int, int, string, error) {
	validate := func(input string) error {
		if _, err := strconv.Atoi(input); err != nil {
			return errors.New("invalid number")
		}
		return nil
	}

	promptAccount := promptui.Prompt{
		Label:    "Enter the Default Account ID",
		Validate: validate,
	}
	promptWorkspace := promptui.Prompt{
		Label:    "Enter the Default Workspace ID",
		Validate: validate,
	}
	promptTeam := promptui.Prompt{
		Label: "Enter the Default Team ID for API monitoring",
	}

	accountIdStr, err := promptAccount.Run()
	if err != nil {
		fmt.Printf("Prompt failed %v\n", err)
		os.Exit(1)
	}
	workspaceIdStr, err := promptWorkspace.Run()
	if err != nil {
		fmt.Printf("Prompt failed %v\n", err)
		os.Exit(1)
	}
	teamId, err := promptTeam.Run()
	if err != nil {
		fmt.Printf("Prompt failed %v\n", err)
		os.Exit(1)
	}
	accountId, err := strconv.Atoi(accountIdStr)
	if err != nil {
		return 0, 0, "", err
	}
	workspaceId, err := strconv.Atoi(workspaceIdStr)
	if err != nil {
		return 0, 0, "", err
	}
	return accountId, workspaceId, teamId, nil
}
