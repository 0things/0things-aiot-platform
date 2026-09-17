package service

import (
	"context"
	"encoding/json"
	"testing"

	"aiot-backend/internal/model"
	"aiot-backend/internal/repository"
	"aiot-backend/internal/tenant"

	"github.com/glebarez/sqlite"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

func TestValidateRuleChainGraph(t *testing.T) {
	tests := []struct {
		name  string
		graph string
		valid bool
	}{
		{name: "entry node without edges is valid", graph: `{"nodes":[{"id":"entry"}],"edges":[]}`, valid: true},
		{name: "connected graph is valid", graph: `{"nodes":[{"id":"entry"},{"id":"filter"}],"edges":[{"source":"entry","target":"filter"}]}`, valid: true},
		{name: "missing nodes is invalid", graph: `{"nodes":[],"edges":[]}`, valid: false},
		{name: "malformed JSON is invalid", graph: `{`, valid: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			valid := validateRuleChainGraph(json.RawMessage(tt.graph)) == nil
			if valid != tt.valid {
				t.Fatalf("validateRuleChainGraph() = %v, want %v", valid, tt.valid)
			}
		})
	}
}

func TestRuleChainServiceRejectsDuplicateName(t *testing.T) {
	db, err := gorm.Open(sqlite.Open("file:rule-chain-service-duplicate-name?mode=memory&cache=shared"), &gorm.Config{})
	require.NoError(t, err)
	require.NoError(t, db.AutoMigrate(&model.RuleChain{}))

	svc := NewRuleChainService(repository.NewRuleChainRepository(db))
	ctx := tenant.WithOrganization(context.Background(), "org-42")
	graph := json.RawMessage(`{"nodes":[{"id":"entry"}],"edges":[]}`)
	require.NoError(t, svc.Create(ctx, &model.RuleChain{Name: "Same name", Graph: graph}))
	err = svc.Create(ctx, &model.RuleChain{Name: "Same name", Graph: graph})
	require.ErrorIs(t, err, ErrRuleChainNameExists)
}
