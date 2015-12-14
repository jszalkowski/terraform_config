package terraform_config

import (
	"os"
	"path/filepath"

	"testing"
	"github.com/stretchr/testify/assert"
)

func TestTerraformStateCreatesDir(t *testing.T) {

	stateDir := "." + string(filepath.Separator) + "tfstate"
	os.Remove(stateDir)

	tfState := TerraformState("foo")
	src, err := os.Stat(stateDir)

	assert.Nil(t, err)
	assert.Equal(t, "./tfstate/foo/terraform.tfstate", tfState)
	assert.True(t, src.IsDir())
	os.Remove(stateDir)
}

