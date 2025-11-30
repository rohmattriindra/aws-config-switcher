package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/ktr0731/go-fuzzyfinder"
)

const awsDir = ".awsconfigs"

func main() {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		fmt.Println("Error getting home directory:", err)
		return
	}

	awsConfigPath := filepath.Join(homeDir, ".aws", "config")
	awsCredentialsPath := filepath.Join(homeDir, ".aws", "credentials")
	configsPath := filepath.Join(homeDir, awsDir)

	if _, err := os.Stat(configsPath); os.IsNotExist(err) {
		fmt.Printf("Creating directory for AWS configs: %s\n", configsPath)
		if err := os.Mkdir(configsPath, 0755); err != nil {
			fmt.Println("Error creating config directory:", err)
			return
		}
	}

	configDirs, err := getConfigDirs(configsPath)
	if err != nil {
		fmt.Println("Error reading config directory:", err)
		return
	}

	selectedIndex, err := fuzzyfinder.Find(
		configDirs,
		func(i int) string {
			return configDirs[i]
		},
	)
	if err != nil {
		fmt.Println("No selection made or error:", err)
		return
	}

	selectedProfilePath := filepath.Join(configsPath, configDirs[selectedIndex])

	if err := backupAndReplaceAWSConfigs(awsConfigPath, awsCredentialsPath, selectedProfilePath); err != nil {
		fmt.Println("Error switching AWS config:", err)
		return
	}

	callerIdentity, err := getAWSCallerIdentity()
	if err != nil {
		fmt.Println("Error getting AWS caller identity:", err)
		return
	}

	fmt.Println("AWS configuration switched successfully!")
	fmt.Printf("Switched AWS configuration to profile: %s\n", filepath.Base(selectedProfilePath))
	fmt.Println("####################")

	printAWSIdentityWithOrange(callerIdentity)
}

// getConfigDirs returns a list of directories in the given AWS configs path.
func getConfigDirs(dir string) ([]string, error) {
	files, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	var configDirs []string
	for _, file := range files {
		if file.IsDir() {
			configDirs = append(configDirs, file.Name())
		}
	}
	return configDirs, nil
}

// backupAndReplaceAWSConfigs backs up the current AWS config and credentials, then sets the selected profile as default.
func backupAndReplaceAWSConfigs(configPath, credentialsPath, selectedProfilePath string) error {
	// Backup current config and credentials files
	if err := backupFile(configPath); err != nil {
		return fmt.Errorf("failed to backup config: %w", err)
	}
	if err := backupFile(credentialsPath); err != nil {
		return fmt.Errorf("failed to backup credentials: %w", err)
	}

	// Replace the selected profile's config and credentials
	if err := copyFile(filepath.Join(selectedProfilePath, "config"), configPath); err != nil {
		return fmt.Errorf("failed to replace config: %w", err)
	}

	// Replace the credentials with the selected profile's credentials under the [default] section
	if err := rewriteCredentialsAsDefault(credentialsPath, filepath.Join(selectedProfilePath, "credentials")); err != nil {
		return fmt.Errorf("failed to rewrite credentials as default: %w", err)
	}

	return nil
}

// backupFile creates a backup of the specified file.
func backupFile(filePath string) error {
	backupPath := filePath + ".backup"
	input, err := os.ReadFile(filePath)
	if err != nil {
		return err
	}
	return os.WriteFile(backupPath, input, 0644)
}

// copyFile copies a file from src to dst.
func copyFile(src, dst string) error {
	input, err := os.ReadFile(src)
	if err != nil {
		return err
	}

	if err := os.WriteFile(dst, input, 0644); err != nil {
		return err
	}
	return nil
}

// rewriteCredentialsAsDefault replaces the credentials in the file, placing the selected profile's credentials under [default].
func rewriteCredentialsAsDefault(credentialsPath, selectedCredentialsPath string) error {
	// Read the selected profile's credentials
	selectedCreds, err := os.ReadFile(selectedCredentialsPath)
	if err != nil {
		return err
	}

	// Write the selected credentials under the [default] section
	credContent := string(selectedCreds)
	credContent = strings.Replace(credContent, fmt.Sprintf("[%s]", filepath.Base(selectedCredentialsPath)), "[default]", 1)

	// Overwrite the credentials file with the modified content
	if err := os.WriteFile(credentialsPath, []byte(credContent), 0644); err != nil {
		return err
	}

	return nil
}

// getAWSCallerIdentity executes the `aws sts get-caller-identity` command and returns the output.
func getAWSCallerIdentity() (string, error) {
	cmd := exec.Command("aws", "sts", "get-caller-identity")
	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("failed to get AWS caller identity: %w", err)
	}
	return string(output), nil
}

// printAWSIdentityWithOrange prints the AWS identity with ARN in orange color.
func printAWSIdentityWithOrange(identity string) {
	// Define the ANSI color codes
	const (
		orange = "\033[38;5;214m" // Use brown as a workaround for orange
		reset  = "\033[0m"
	)

	// Print the standard part without color
	fmt.Println("Current AWS Identity:")

	// Print the identity with ARN in orange
	lines := strings.Split(identity, "\n")
	for _, line := range lines {
		if strings.Contains(line, "Arn") {
			// Print the ARN part in orange
			fmt.Printf("%s%s%s\n", orange, line, reset)
		} else {
			// Print other lines in default color
			fmt.Println(line)
		}
	}
}
