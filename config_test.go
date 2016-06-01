package terraform_config

import (
	"os"
	"path/filepath"

	"testing"
	"github.com/stretchr/testify/assert"
)

func TestTerraformStateCreatesTfStateDir(t *testing.T) {

	stateDir := "." + string(filepath.Separator) + "tfstate"
	os.RemoveAll(stateDir)

	tfState := TerraformState("foo")
	src, err := os.Stat(stateDir)

	assert.Nil(t, err)
	assert.Equal(t, "./tfstate/foo/terraform.tfstate", tfState)
	assert.True(t, src.IsDir())
	os.RemoveAll(stateDir)
}

func TestTerraformStateCreatesEnvironmentDir(t *testing.T) {

	stateDir := "." + string(filepath.Separator) + "tfstate"
	os.RemoveAll(stateDir)
	os.MkdirAll("." + string(filepath.Separator) + "tfstate", 0777)

	tfState := TerraformState("foo")
	src, err := os.Stat(stateDir)

	assert.Nil(t, err)
	assert.Equal(t, "./tfstate/foo/terraform.tfstate", tfState)
	assert.True(t, src.IsDir())
	os.RemoveAll(stateDir)
}

