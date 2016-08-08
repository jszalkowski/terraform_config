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

type TerraformConfig struct {
	// s3 bucket name
	S3_bucket string
	// s3 key
	S3_key string
	// aws region
	Aws_region string
	// config dir path
	Config_path string
	// state dir path
	State_path string
	// tf vars file
	Tf_vars string
}

var cyan = color.New(color.FgCyan).SprintFunc()
var red = color.New(color.FgRed).SprintFunc()
var bold = color.New(color.FgWhite, color.Bold).SprintFunc()

func LoadConfig(config string, environment string) *TerraformConfig {
	derivedConfigPath := configLocation(config)
	tfVars := TerraformVars(configLocation(config), environment)
	terraformConfig, err := LoadTerraformConfig(tfVars)
	terraformConfig.Tf_vars = tfVars
	terraformConfig.State_path = terraformState(environment)
	terraformConfig.Config_path = derivedConfigPath

	if err != nil {
		command.Error("Error Loading Terraform Vars", err)
	}
	fmt.Printf("Using terraform config: %s\n", cyan(tfVars))
	fmt.Println()
	fmt.Println("AWS credentials")
	fmt.Println("s3 bucket: ", bold(terraformConfig.S3_bucket))
	fmt.Println("s3 key:    ", bold(terraformConfig.S3_key))
	fmt.Println("aws region:", bold(terraformConfig.Aws_region))
	fmt.Println()
	return terraformConfig
}

func LoadTerraformConfig(path string) (*TerraformConfig, error) {
	var value TerraformConfig

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

// has side effects
func terraformState(environment string) string {
	directory := "." + string(filepath.Separator) + "tfstate" + string(filepath.Separator) + environment
	src, err := os.Stat(directory)
	terraformState := fmt.Sprintf("%s/terraform.tfstate", directory)
	if err != nil {
		os.MkdirAll(directory, 0777)
	} else if !src.IsDir() {
		command.Error("tfstate is not a directory", err)
	}
	return terraformState
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
