import { createMCPClient } from "@ai-sdk/mcp";

export const ALLOWED_MCP_TOOLS = new Set([
  "iot_query_devices",
  "iot_get_device_detail",
  "iot_query_telemetry_history",
]);

export interface MCPContext {
  tools: Record<string, any>;
  close: () => Promise<void>;
}

export async function createAuthenticatedMCPClient(
  rawToken: string,
  mcpUrl: string,
  timeoutMs: number = 8000
): Promise<MCPContext> {
  const timeoutPromise = new Promise<never>((_, reject) =>
    setTimeout(() => reject(new Error("MCP connection timed out")), timeoutMs)
  );

  const clientPromise = (async () => {
    const client = await createMCPClient({
      transport: {
        type: "http",
        url: mcpUrl,
        headers: {
          Authorization: `Bearer ${rawToken}`,
        },
      },
    });

    const allTools = await client.tools();
    const whitelistedTools: Record<string, any> = {};

    for (const [toolName, toolDef] of Object.entries(allTools)) {
      if (ALLOWED_MCP_TOOLS.has(toolName)) {
        whitelistedTools[toolName] = toolDef;
      }
    }

    return {
      tools: whitelistedTools,
      close: async () => {
        try {
          await client.close();
        } catch {
          // ignore cleanup errors
        }
      },
    };
  })();

  return Promise.race([clientPromise, timeoutPromise]);
}
