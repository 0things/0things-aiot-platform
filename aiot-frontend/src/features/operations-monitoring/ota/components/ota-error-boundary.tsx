import { Component, type ReactNode } from 'react'
import { AlertTriangle } from 'lucide-react'
import { Button } from '@/components/ui/button'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'

interface Props {
  children: ReactNode
}

interface State {
  hasError: boolean
  error?: Error
}

export class OTAErrorBoundary extends Component<Props, State> {
  constructor(props: Props) {
    super(props)
    this.state = { hasError: false }
  }

  static getDerivedStateFromError(error: Error): State {
    return { hasError: true, error }
  }

  componentDidCatch(error: Error, errorInfo: React.ErrorInfo) {
    // eslint-disable-next-line no-console
    console.error('OTA Module Error:', error, errorInfo)
  }

  handleReset = () => {
    this.setState({ hasError: false, error: undefined })
  }

  render() {
    if (this.state.hasError) {
      return (
        <div className='flex min-h-[400px] items-center justify-center p-4'>
          <Card className='w-full max-w-md'>
            <CardHeader>
              <CardTitle className='flex items-center gap-2 text-destructive'>
                <AlertTriangle className='h-5 w-5' />
                Something went wrong
              </CardTitle>
            </CardHeader>
            <CardContent className='space-y-4'>
              <p className='text-sm text-muted-foreground'>
                An error occurred in the OTA module. Please try again or contact
                support if the problem persists.
              </p>
              {this.state.error && (
                <details className='text-xs text-muted-foreground'>
                  <summary className='cursor-pointer hover:underline'>
                    Error details
                  </summary>
                  <pre className='mt-2 rounded border p-2'>
                    {this.state.error.message}
                  </pre>
                </details>
              )}
              <div className='flex gap-2'>
                <Button onClick={this.handleReset} className='flex-1'>
                  Try Again
                </Button>
                <Button
                  variant='outline'
                  onClick={() => window.location.reload()}
                  className='flex-1'
                >
                  Reload Page
                </Button>
              </div>
            </CardContent>
          </Card>
        </div>
      )
    }

    return this.props.children
  }
}
