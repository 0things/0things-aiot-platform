import { useMemo, useState } from 'react'
import {
  AssistantChatTransport,
  useChatRuntime,
} from '@assistant-ui/react-ai-sdk'
import { getAuthToken } from '@/api/clients'

export function useCopilotRuntime() {
  const [sessionKey, setSessionKey] = useState(0)

  const gatewayUrl = import.meta.env.VITE_AI_GATEWAY_URL || ''

  const transport = useMemo(() => {
    const apiEndpoint = gatewayUrl
      ? `${gatewayUrl.replace(/\/$/, '')}/v1/ai/chat`
      : '/v1/ai/chat'

    return new AssistantChatTransport({
      api: apiEndpoint,
      headers: () => {
        const token = getAuthToken()
        const headers: Record<string, string> = {}
        if (token) {
          headers['Authorization'] = `Bearer ${token}`
        }
        return headers
      },
    })
  }, [gatewayUrl, sessionKey])

  const runtime = useChatRuntime({
    transport,
  })

  const resetSession = () => {
    setSessionKey((prev) => prev + 1)
  }

  return {
    runtime,
    resetSession,
    sessionKey,
  }
}
