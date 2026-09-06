package handler

import (
	"net/http"

	v1 "aiot-backend/api/v1"
	"aiot-backend/internal/model"
	"aiot-backend/internal/service"

	"github.com/gin-gonic/gin"
)

type RuleNodeDefinitionHandler struct {
	*Handler
	svc service.RuleNodeDefinitionServiceInterface
}

func NewRuleNodeDefinitionHandler(h *Handler, svc service.RuleNodeDefinitionServiceInterface) *RuleNodeDefinitionHandler {
	return &RuleNodeDefinitionHandler{Handler: h, svc: svc}
}

func ruleNodeDefinitionJSON(definition model.RuleNodeDefinition) v1.RuleNodeDefinition {
	return v1.RuleNodeDefinition{
		UUID:          definition.UUID,
		Key:           definition.Key,
		Category:      definition.Category,
		Name:          definition.Name,
		Description:   definition.Description,
		Icon:          definition.Icon,
		ConfigSchema:  definition.ConfigSchema,
		DefaultConfig: definition.DefaultConfig,
		InputPorts:    definition.InputPorts,
		OutputPorts:   definition.OutputPorts,
		ExecutorKey:   definition.ExecutorKey,
		IsSystem:      definition.IsSystem,
	}
}

// ListRuleNodeDefinitions godoc
// @Summary List rule node definitions
// @Description Lists all enabled system rule node definitions.
// @Tags Rule node definitions
// @Produce json
// @Security Bearer
// @Success 200 {object} v1.ApiResponse[v1.ListRuleNodeDefinitionsResponse]
// @Router /rule-node-definitions [get]
func (h *RuleNodeDefinitionHandler) ListRuleNodeDefinitions(c *gin.Context) {
	definitions, err := h.svc.ListEnabled(c)
	if err != nil {
		v1.HandleError(c, http.StatusInternalServerError, err, nil)
		return
	}

	items := make([]v1.RuleNodeDefinition, len(definitions))
	for i, definition := range definitions {
		items[i] = ruleNodeDefinitionJSON(definition)
	}
	v1.HandleSuccess(c, v1.ListRuleNodeDefinitionsResponse{Items: items})
}
