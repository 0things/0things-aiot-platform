import { AssistantRuntimeProvider } from '@assistant-ui/react'
import { AssistantModal } from '@/components/assistant-ui/elements/assistant-modal.aui'
import { useCopilotRuntime } from './runtime/use-copilot-runtime'

export function AICopilot() {
  const isEnabled = import.meta.env.VITE_AI_COPILOT_ENABLED !== 'false'
  const { runtime } = useCopilotRuntime()

  if (!isEnabled) {
    return null
  }

  return (
    <AssistantRuntimeProvider runtime={runtime}>
      <AssistantModal />
    </AssistantRuntimeProvider>
  )
}

export * from './runtime/use-copilot-runtime'
