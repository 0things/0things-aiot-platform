package repository

import (
	"context"
	"testing"
	"time"

	"0things-backend/internal/model"
	"github.com/stretchr/testify/require"
)

func TestRuleRepository(t *testing.T) {
	store := newRepositoryTestDB(t, &model.Rule{}, &model.RuleExecution{})
	repo := NewRuleRepository(store)
	ctx := context.Background()
	rule := &model.Rule{Name: "temperature", Type: "sql", Status: "draft"}
	require.NoError(t, repo.Create(ctx, rule))
	require.NoError(t, repo.UpdateStatus(ctx, rule, "enabled"))
	require.NoError(t, repo.CreateExecution(ctx, &model.RuleExecution{RuleID: rule.ID, RuleName: rule.Name, Status: "success", TriggeredAt: time.Now()}))
	require.NoError(t, repo.UpdateExecutionStats(ctx, rule, time.Now()))

	items, total, err := repo.List(ctx, 1, 20, "sql", "enabled", "temperature")
	require.NoError(t, err)
	require.EqualValues(t, 1, total)
	require.Len(t, items, 1)
	executions, executionTotal, err := repo.ListExecutions(ctx, rule.ID, 1, 20)
	require.NoError(t, err)
	require.EqualValues(t, 1, executionTotal)
	require.Len(t, executions, 1)
}
