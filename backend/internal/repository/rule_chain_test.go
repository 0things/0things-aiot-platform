package repository

import (
	"context"
	"encoding/json"
	"testing"

	"aiot-backend/internal/model"
	"aiot-backend/internal/tenant"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

func TestRuleChainRepositoryPersistsGraphWithinTenant(t *testing.T) {
	db, err := gorm.Open(sqlite.Open("file:rule-chain-repository?mode=memory&cache=shared"), &gorm.Config{})
	if err != nil {
		t.Fatal(err)
	}
	if err := db.AutoMigrate(&model.RuleChain{}); err != nil {
		t.Fatal(err)
	}

	repo := NewRuleChainRepository(db)
	ctx := tenant.WithTenant(context.Background(), 42)
	chain := &model.RuleChain{UUID: "chain-1", Name: "Telemetry flow", Status: "draft", Version: 1, Graph: json.RawMessage(`{"nodes":[{"id":"entry"}],"edges":[]}`)}
	if err := repo.Create(ctx, chain); err != nil {
		t.Fatal(err)
	}

	loaded, err := repo.Find(ctx, "chain-1")
	if err != nil {
		t.Fatal(err)
	}
	if loaded.OrganizationID != 42 || string(loaded.Graph) != string(chain.Graph) {
		t.Fatalf("unexpected persisted chain: %+v", loaded)
	}

	otherTenant := tenant.WithTenant(context.Background(), 43)
	if _, err := repo.Find(otherTenant, "chain-1"); err != ErrNotFound {
		t.Fatalf("cross-tenant find error = %v, want ErrNotFound", err)
	}
}
