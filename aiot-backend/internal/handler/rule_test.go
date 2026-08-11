package handler

import (
	"encoding/json"
	"testing"
	"time"

	"0things-backend/internal/model"
	"github.com/stretchr/testify/require"
)

func TestRuleResponseMapping(t *testing.T) {
	now := time.Now().UTC().Truncate(time.Second)
	rule := ruleJSON(model.Rule{
		ID: 7, Name: "temperature", TriggerConfig: json.RawMessage(`{"source":"device"}`),
		Tags: json.RawMessage(`["production"]`), CreatedAt: now, UpdatedAt: now,
	})
	require.EqualValues(t, 7, rule.ID)
	require.Equal(t, `{"source":"device"}`, rule.TriggerConfig)
	require.Equal(t, []string{"production"}, rule.Tags)

	execution := ruleExecutionJSON(model.RuleExecution{ID: 8, RuleID: 7, TriggerData: json.RawMessage(`{"value":42}`), TriggeredAt: now, CreatedAt: now})
	require.EqualValues(t, 7, execution.RuleID)
	require.Equal(t, `{"value":42}`, execution.TriggerData)
}
