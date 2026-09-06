import { streamText, convertToModelMessages, stepCountIs } from "ai";
import { createOpenAI } from "@ai-sdk/openai";
import { GatewayConfig } from "./config.js";
import { UserContext } from "./auth.js";
import { createAuthenticatedMCPClient } from "./mcp.js";

const SYSTEM_PROMPT = `You are the 0things AIoT Copilot (0things 物联网智能助手).
Your job is to assist operators in monitoring, inspecting, and analyzing IoT devices and telemetry data within their organization.

Guidelines:
1. You have access to three read-only tools:
   - iot_query_devices: Query device list, online state, and basic info.
   - iot_get_device_detail: Retrieve comprehensive device details including tags, device shadow, and thing-model properties.
   - iot_query_telemetry_history: Query time-series data points and summary statistics (min/max/avg/latest) for a specific property.
2. Read-Only Constraint: You CANNOT perform write operations, modify device states, delete resources, change configurations, or issue control commands. If the user asks for any write/control actions, politely explain that Copilot currently operates in read-only mode for safety and security.
3. Multi-Tenant Scope: Your tools automatically operate within the authenticated organization scope. Only query and reference data returned by the tools.
4. Response Format: Provide concise, clear, and structured answers (using Markdown tables, lists, or metrics when appropriate). Default to Chinese unless the user asks in another language.`;

export async function handleChatRequest(
  req: Request,
  user: UserContext,
  config: GatewayConfig
): Promise<Response> {
  let body: any;
  try {
    body = await req.json();
  } catch {
    return new Response(JSON.stringify({ error: "Invalid JSON body" }), {
      status: 400,
      headers: { "Content-Type": "application/json" },
    });
  }

  const rawMessages = body.messages;
  if (!Array.isArray(rawMessages) || rawMessages.length === 0) {
    return new Response(JSON.stringify({ error: "messages array is required" }), {
      status: 400,
      headers: { "Content-Type": "application/json" },
    });
  }

  // Normalize messages
  const normalizedMessages = rawMessages.map((m: any) => {
    if (m.parts && Array.isArray(m.parts)) {
      return m;
    }
    if (typeof m.content === "string") {
      return {
        id: m.id || String(Date.now()),
        role: m.role || "user",
        parts: [{ type: "text", text: m.content }],
      };
    }
    return m;
  });

  let mcpContext: { tools: Record<string, any>; close: () => Promise<void> } | null = null;

  try {
    // 1. Connect to MCP using the user's verified token
    mcpContext = await createAuthenticatedMCPClient(
      user.rawToken,
      config.MCP_STREAMABLE_HTTP_URL
    );

    // 2. Initialize OpenAI-compatible model provider
    const openai = createOpenAI({
      baseURL: config.AI_MODEL_BASE_URL,
      apiKey: config.AI_MODEL_API_KEY,
    });
    const model = openai(config.AI_MODEL_ID);

    // 3. Convert messages to model format
    const modelMessages = await convertToModelMessages(normalizedMessages);

    // 4. Stream text with tools
    const result = streamText({
      model,
      system: SYSTEM_PROMPT,
      messages: modelMessages,
      tools: mcpContext.tools,
      stopWhen: stepCountIs(5),
      onFinish: async () => {
        if (mcpContext) {
          await mcpContext.close();
        }
      },
      onError: async ({ error }) => {
        console.error("AI chat streaming error:", error);
        if (mcpContext) {
          await mcpContext.close();
        }
      },
    });

    return result.toUIMessageStreamResponse({
      getErrorMessage: (error: unknown) => {
        return sanitizeStreamError(error);
      },
    });
  } catch (err: any) {
    if (mcpContext) {
      await mcpContext.close();
    }
    const isTimeout =
      err?.name === "TimeoutError" ||
      err?.message?.includes("timed out") ||
      err?.code === "ETIMEDOUT";
    const userMessage = isTimeout
      ? "Connection to IoT data service timed out. Please try again later."
      : "Assistant service is currently unavailable. Please try again later.";

    return new Response(
      JSON.stringify({
        error: userMessage,
        code: isTimeout ? "MCP_TIMEOUT" : "GATEWAY_ERROR",
      }),
      {
        status: 503,
        headers: { "Content-Type": "application/json" },
      }
    );
  }
}

function sanitizeStreamError(error: unknown): string {
  if (typeof error === "object" && error !== null) {
    const err = error as any;
    const status = err.status || err.statusCode || err.response?.status;
    if (status === 401 || status === 403) {
      return "Model provider authentication failed. Please check the API key configuration.";
    }
    if (status === 429) {
      return "Model provider rate limit exceeded or quota exhausted. Please try again later.";
    }
    if (err.name === "AbortError") {
      return "Request was aborted.";
    }
  }
  return "Assistant service is currently unavailable. Please try again later.";
}

