import { Badge } from '@/components/ui/badge'
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from '@/components/ui/card'

export function InstanceInfo() {
  return (
    <div className='space-y-4'>
      <div className='grid gap-4 md:grid-cols-2'>
        <Card>
          <CardHeader className='pb-3'>
            <CardTitle className='text-base'>Instance Details</CardTitle>
          </CardHeader>
          <CardContent className='space-y-3'>
            <div className='flex items-center justify-between'>
              <span className='text-sm text-muted-foreground'>Instance ID</span>
              <Badge variant='secondary'>IOT-PROD-001</Badge>
            </div>
            <div className='flex items-center justify-between'>
              <span className='text-sm text-muted-foreground'>Region</span>
              <span className='text-sm font-medium'>us-east-1</span>
            </div>
            <div className='flex items-center justify-between'>
              <span className='text-sm text-muted-foreground'>Version</span>
              <span className='text-sm font-medium'>v3.2.1</span>
            </div>
            <div className='flex items-center justify-between'>
              <span className='text-sm text-muted-foreground'>Created</span>
              <span className='text-sm font-medium'>2024-01-15</span>
            </div>
          </CardContent>
        </Card>

        <Card>
          <CardHeader className='pb-3'>
            <CardTitle className='text-base'>Performance Metrics</CardTitle>
          </CardHeader>
          <CardContent className='space-y-3'>
            <div className='flex items-center justify-between'>
              <span className='text-sm text-muted-foreground'>CPU Usage</span>
              <span className='text-sm font-medium'>45.2%</span>
            </div>
            <div className='flex items-center justify-between'>
              <span className='text-sm text-muted-foreground'>Memory</span>
              <span className='text-sm font-medium'>2.1 GB / 4 GB</span>
            </div>
            <div className='flex items-center justify-between'>
              <span className='text-sm text-muted-foreground'>Storage</span>
              <span className='text-sm font-medium'>78 GB / 100 GB</span>
            </div>
            <div className='flex items-center justify-between'>
              <span className='text-sm text-muted-foreground'>Network I/O</span>
              <span className='text-sm font-medium'>124 MB/s</span>
            </div>
          </CardContent>
        </Card>
      </div>

      <Card>
        <CardHeader className='pb-3'>
          <CardTitle className='text-base'>Connection Status</CardTitle>
          <CardDescription>Real-time connection information</CardDescription>
        </CardHeader>
        <CardContent>
          <div className='grid gap-4 md:grid-cols-3'>
            <div className='text-center'>
              <div className='text-2xl font-bold text-green-600'>1,234</div>
              <p className='text-sm text-muted-foreground'>Connected Devices</p>
            </div>
            <div className='text-center'>
              <div className='text-2xl font-bold text-blue-600'>89.2 MB/s</div>
              <p className='text-sm text-muted-foreground'>Data Throughput</p>
            </div>
            <div className='text-center'>
              <div className='text-2xl font-bold text-purple-600'>12.4 ms</div>
              <p className='text-sm text-muted-foreground'>Avg Latency</p>
            </div>
          </div>
        </CardContent>
      </Card>

      <Card>
        <CardHeader className='pb-3'>
          <CardTitle className='text-base'>Service Health</CardTitle>
        </CardHeader>
        <CardContent>
          <div className='space-y-3'>
            <div className='flex items-center justify-between'>
              <div className='flex items-center space-x-2'>
                <div className='h-2 w-2 rounded-full bg-green-500'></div>
                <span className='text-sm'>Message Broker</span>
              </div>
              <Badge variant='default' className='bg-green-100 text-green-800'>
                Healthy
              </Badge>
            </div>
            <div className='flex items-center justify-between'>
              <div className='flex items-center space-x-2'>
                <div className='h-2 w-2 rounded-full bg-green-500'></div>
                <span className='text-sm'>Data Storage</span>
              </div>
              <Badge variant='default' className='bg-green-100 text-green-800'>
                Healthy
              </Badge>
            </div>
            <div className='flex items-center justify-between'>
              <div className='flex items-center space-x-2'>
                <div className='h-2 w-2 rounded-full bg-yellow-500'></div>
                <span className='text-sm'>API Gateway</span>
              </div>
              <Badge
                variant='default'
                className='bg-yellow-100 text-yellow-800'
              >
                Warning
              </Badge>
            </div>
            <div className='flex items-center justify-between'>
              <div className='flex items-center space-x-2'>
                <div className='h-2 w-2 rounded-full bg-green-500'></div>
                <span className='text-sm'>Authentication</span>
              </div>
              <Badge variant='default' className='bg-green-100 text-green-800'>
                Healthy
              </Badge>
            </div>
          </div>
        </CardContent>
      </Card>
    </div>
  )
}
