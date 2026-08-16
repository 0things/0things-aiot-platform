import { createFileRoute, redirect } from '@tanstack/react-router'

export const Route = createFileRoute('/_authenticated/rule-engine/')({
  beforeLoad: () => {
    throw redirect({ to: '/rule-engine/scene-linkage' })
  },
})
