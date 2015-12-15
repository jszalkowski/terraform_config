package terraform_config

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/fatih/color"
	"github.com/hashicorp/hcl"
	"github.com/nadnerb/cli_command"
)

type AwsConfig struct {
	// s3 bucket name
	S3_bucket string
	// s3 key
	S3_key string
	// aws region
	Aws_region string
	// ssh key
	Key_path string
}

var cyan = color.New(color.FgCyan).SprintFunc()
var red = color.New(color.FgRed).SprintFunc()
var bold = color.New(color.FgWhite, color.Bold).SprintFunc()

func LoadConfig(config string, environment string) *AwsConfig {
	tfVars := TerraformVars(configLocation(config), environment)
	awsConfig, err := LoadAwsConfig(tfVars)
	if err != nil {
		command.Error("Error Loading Terraform Vars", err)
	}
	fmt.Printf("Using terraform config: %s\n", cyan(tfVars))
	fmt.Println()
	fmt.Println("AWS credentials")
	fmt.Println("s3 bucket: ", bold(awsConfig.S3_bucket))
	fmt.Println("s3 key:    ", bold(awsConfig.S3_key))
	fmt.Println("aws region:", bold(awsConfig.Aws_region))
	fmt.Println()
	return awsConfig
}

func LoadAwsConfig(path string) (*AwsConfig, error) {
	var value AwsConfig

	if _, err := os.Stat(path); err != nil {
		return nil, err
	}
	err := hcl.Decode(&value, readFile(path))
	if err != nil {
		return nil, err
	}
	return &value, nil
}

func TerraformVars(configLocation string, environment string) string {
	return filepath.Clean(fmt.Sprintf("%s/%s.tfvars", configLocation, environment))
}

func TerraformState(environment string) string {
	src, err := os.Stat("." + string(filepath.Separator) + "tfstate")
	terraformState := fmt.Sprintf("./tfstate/%s/terraform.tfstate", environment)
	if err != nil {
		tempStateFile(environment, terraformState)
	} else if !src.IsDir() {
		command.Error("tfstate is not a directory", err)
	}
	return terraformState
}

func tempStateFile(environment string, terraformState string) {
	os.MkdirAll("." + string(filepath.Separator) + "tfstate" + string(filepath.Separator) + environment, 0777)
	bytes := []byte("")
	err := ioutil.WriteFile(terraformState, bytes, 0644)
	if err != nil {
		command.Warn("create tfstate fail, please make sure you ensure you are syncing state to S3", err)
	}
}

func configLocation(config string) string {
	if len(config) > 0 {
		if _, err := os.Stat(config); os.IsNotExist(err) {
			command.Error("Directory does not exist", err)
		}
		return config
	} else {
		return defaultConfig()
	}
}

func defaultConfig() string {
	defaultConfig, _ := filepath.Abs("./config/")
	fmt.Printf("Using default config location: %s\n", cyan(defaultConfig))
	if _, err := os.Stat(defaultConfig); os.IsNotExist(err) {
		command.Error("Directory does not exist", err)
	}
	return defaultConfig
}

func readFile(path string) string {
	d, err := ioutil.ReadFile(path)
	if err != nil {
		fmt.Errorf(
			"Error parsing %s: %s", path, err)
	}

	return string(d)
}
