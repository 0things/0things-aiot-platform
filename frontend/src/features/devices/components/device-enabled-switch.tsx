import { toast } from 'sonner'
import { Switch } from '@/components/ui/switch'
import { useSetDeviceEnabled } from '../api/queries'
import { type Device } from '../data/schema'

export function DeviceEnabledSwitch({ device }: { device: Device }) {
  const setDeviceEnabled = useSetDeviceEnabled()

  return (
    <div className='flex justify-center'>
      <Switch
        checked={device.enabled}
        disabled={setDeviceEnabled.isPending}
        onCheckedChange={async (checked) => {
          try {
            await setDeviceEnabled.mutateAsync({
              deviceKey: device.deviceKey,
              data: { enabled: checked },
            })
            toast.success(
              checked
                ? 'Device enabled successfully!'
                : 'Device disabled successfully!'
            )
          } catch {
            toast.error('Failed to update device status')
          }
        }}
        aria-label='Toggle device enabled'
      />
    </div>
  )
}
