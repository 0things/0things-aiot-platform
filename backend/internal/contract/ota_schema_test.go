package contract_test

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/xeipuuv/gojsonschema"
)

func TestOTAContractExamples(t *testing.T) {
	root := filepath.Join("..", "..", "..")
	cases := []struct {
		name    string
		schema  string
		example string
	}{
		{"command", filepath.Join(root, "openspec", "changes", "ota-batch-upgrade-pipeline", "schemas", "ota-upgrade-command.v1.json"), filepath.Join(root, "openspec", "changes", "ota-batch-upgrade-pipeline", "schemas", "examples", "ota-upgrade-command.v1.json")},
		{"progress report", filepath.Join(root, "openspec", "changes", "ota-batch-upgrade-pipeline", "schemas", "ota-upgrade-report.v1.json"), filepath.Join(root, "openspec", "changes", "ota-batch-upgrade-pipeline", "schemas", "examples", "ota-upgrade-progress.v1.json")},
		{"inform report", filepath.Join(root, "openspec", "changes", "ota-batch-upgrade-pipeline", "schemas", "ota-upgrade-report.v1.json"), filepath.Join(root, "openspec", "changes", "ota-batch-upgrade-pipeline", "schemas", "examples", "ota-upgrade-inform.v1.json")},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			result, err := gojsonschema.Validate(gojsonschema.NewReferenceLoader("file://"+tc.schema), gojsonschema.NewReferenceLoader("file://"+tc.example))
			require.NoError(t, err)
			require.True(t, result.Valid(), result.Errors())
		})
	}
}
