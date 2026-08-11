import { DotsHorizontalIcon } from '@radix-ui/react-icons'
import { Link } from '@tanstack/react-router'
import { type Row } from '@tanstack/react-table'
import { Pencil, Trash2, Power, PowerOff, Eye } from 'lucide-react'
import { Button } from '@/components/ui/button'
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuSeparator,
  DropdownMenuShortcut,
  DropdownMenuTrigger,
} from '@/components/ui/dropdown-menu'
import { type Device } from '../data/schema'
import { useDevices } from './devices-provider'

type DataTableRowActionsProps = {
  row: Row<Device>
}

export function DataTableRowActions({ row }: DataTableRowActionsProps) {
  const { setOpen, setCurrentRow } = useDevices()
  const device = row.original
  const isEnabled = device.enabled

  return (
    <>
      <DropdownMenu modal={false}>
        <DropdownMenuTrigger asChild>
          <Button
            variant='ghost'
            className='flex h-8 w-8 p-0 data-[state=open]:bg-muted'
          >
            <DotsHorizontalIcon className='h-4 w-4' />
            <span className='sr-only'>Open menu</span>
          </Button>
        </DropdownMenuTrigger>
        <DropdownMenuContent align='end' className='w-[180px]'>
          <DropdownMenuItem asChild>
            <Link
              to='/device-management/devices/$deviceKey'
              params={{ deviceKey: device.deviceKey }}
              className='cursor-pointer'
            >
              Detail
              <DropdownMenuShortcut>
                <Eye size={16} />
              </DropdownMenuShortcut>
            </Link>
          </DropdownMenuItem>

          <DropdownMenuSeparator />

          <DropdownMenuItem
            onClick={() => {
              setCurrentRow(row.original)
              setOpen('edit')
            }}
          >
            Edit
            <DropdownMenuShortcut>
              <Pencil size={16} />
            </DropdownMenuShortcut>
          </DropdownMenuItem>

          <DropdownMenuItem
            onClick={() => {
              setCurrentRow(row.original)
              setOpen('enable')
            }}
          >
            {isEnabled ? 'Disable' : 'Enable'}
            <DropdownMenuShortcut>
              {isEnabled ? <PowerOff size={16} /> : <Power size={16} />}
            </DropdownMenuShortcut>
          </DropdownMenuItem>

          <DropdownMenuSeparator />
          <DropdownMenuItem
            onClick={() => {
              setCurrentRow(row.original)
              setOpen('delete')
            }}
            className='text-red-500!'
          >
            Delete
            <DropdownMenuShortcut>
              <Trash2 size={16} />
            </DropdownMenuShortcut>
          </DropdownMenuItem>
        </DropdownMenuContent>
      </DropdownMenu>
    </>
  )
}
