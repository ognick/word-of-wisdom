package arch

import (
	"testing"

	archgo "github.com/arch-go/arch-go/api"
	config "github.com/arch-go/arch-go/api/configuration"
	"github.com/stretchr/testify/require"
)

func TestArchitecture(t *testing.T) {
	configuration, err := config.LoadConfig("../arch-go.yml")
	require.NoError(t, err)
	moduleInfo := config.Load("github.com/ognick/word_of_wisdom")

	result := archgo.CheckArchitecture(moduleInfo, *configuration)
	if !result.Pass {
		t.Fatal("Project doesn't pass architecture tests")
	}
}
