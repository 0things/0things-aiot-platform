import { type TFunction } from 'i18next'
import {
  Construction,
  LayoutDashboard,
  Monitor,
  Bug,
  ListTodo,
  FileX,
  HelpCircle,
  Lock,
  Bell,
  Package,
  Palette,
  ServerOff,
  Settings,
  Wrench,
  UserCog,
  UserX,
  Users,
  MessagesSquare,
  ShieldCheck,
  AudioWaveform,
  Command,
  GalleryVerticalEnd,
  Smartphone,
  Activity,
  Upload,
  List,
  Zap,
  Cpu,
  Workflow,
} from 'lucide-react'
import { ClerkLogo } from '@/assets/clerk-logo'
import { type SidebarData } from '../types'

export const getSidebarData = (t: TFunction): SidebarData => ({
  user: {
    name: 'satnaing',
    email: 'satnaingdev@gmail.com',
    avatar: '/avatars/shadcn.jpg',
  },
  teams: [
    {
      name: '0things Admin',
      logo: Command,
      plan: 'Vite + ShadcnUI',
    },
    {
      name: 'Acme Inc',
      logo: GalleryVerticalEnd,
      plan: 'Enterprise',
    },
    {
      name: 'Acme Corp.',
      logo: AudioWaveform,
      plan: 'Startup',
    },
  ],
  navGroups: [
    {
      title: t('navigation:sidebar.general'),
      items: [
        {
          title: t('navigation:sidebar.dashboard'),
          url: '/',
          icon: LayoutDashboard,
        },
        {
          title: 'IoT Dashboard',
          url: '/iot-dashboard',
          icon: Cpu,
        },
        {
          title: t('navigation:sidebar.tasks'),
          url: '/tasks',
          icon: ListTodo,
        },
        {
          title: t('navigation:sidebar.apps'),
          url: '/apps',
          icon: Package,
        },
        {
          title: t('navigation:sidebar.chats'),
          url: '/chats',
          badge: '3',
          icon: MessagesSquare,
        },
        {
          title: t('navigation:sidebar.users'),
          url: '/users',
          icon: Users,
        },
        {
          title: t('navigation:sidebar.deviceManagement'),
          icon: Smartphone,
          items: [
            {
              title: t('navigation:sidebar.products'),
              url: '/device-management/products',
              icon: Package,
            },
            {
              title: t('navigation:sidebar.devices'),
              url: '/device-management/devices',
              icon: Smartphone,
            },
          ],
        },
        {
          title: t('operationsMonitoring:title'),
          icon: Activity,
          items: [
            {
              title: t('ota:title'),
              url: '/operations-monitoring/ota/packages',
              icon: Upload,
            },
            {
              title: t('operationsMonitoring:events.title'),
              url: '/operations-monitoring/events',
              icon: List,
            },
          ],
        },
        {
          title: t('navigation:sidebar.ruleEngine'),
          icon: Zap,
          items: [
            {
              title: t('navigation:sidebar.sceneLinkage'),
              url: '/rule-engine/scene-linkage',
              icon: Workflow,
            },
          ],
        },
        {
          title: 'Secured by Clerk',
          icon: ClerkLogo,
          items: [
            {
              title: 'Sign In',
              url: '/clerk/sign-in',
            },
            {
              title: 'Sign Up',
              url: '/clerk/sign-up',
            },
            {
              title: 'User Management',
              url: '/clerk/user-management',
            },
          ],
        },
      ],
    },
    {
      title: t('navigation:sidebar.pages'),
      items: [
        {
          title: 'Auth',
          icon: ShieldCheck,
          items: [
            {
              title: t('auth:signIn.title'),
              url: '/sign-in',
            },
            {
              title: `${t('auth:signIn.title')} (2 Col)`,
              url: '/sign-in-2',
            },
            {
              title: t('auth:signUp.title'),
              url: '/sign-up',
            },
            {
              title: t('auth:forgotPassword.title'),
              url: '/forgot-password',
            },
            {
              title: t('auth:otp.title'),
              url: '/otp',
            },
          ],
        },
        {
          title: 'Errors',
          icon: Bug,
          items: [
            {
              title: 'Unauthorized',
              url: '/errors/unauthorized',
              icon: Lock,
            },
            {
              title: 'Forbidden',
              url: '/errors/forbidden',
              icon: UserX,
            },
            {
              title: 'Not Found',
              url: '/errors/not-found',
              icon: FileX,
            },
            {
              title: 'Internal Server Error',
              url: '/errors/internal-server-error',
              icon: ServerOff,
            },
            {
              title: 'Maintenance Error',
              url: '/errors/maintenance-error',
              icon: Construction,
            },
          ],
        },
      ],
    },
    {
      title: 'Other',
      items: [
        {
          title: t('navigation:sidebar.settings'),
          icon: Settings,
          items: [
            {
              title: 'Profile',
              url: '/settings',
              icon: UserCog,
            },
            {
              title: t('settings:tabs.account'),
              url: '/settings/account',
              icon: Wrench,
            },
            {
              title: t('settings:tabs.appearance'),
              url: '/settings/appearance',
              icon: Palette,
            },
            {
              title: t('settings:tabs.notifications'),
              url: '/settings/notifications',
              icon: Bell,
            },
            {
              title: t('settings:tabs.display'),
              url: '/settings/display',
              icon: Monitor,
            },
          ],
        },
        {
          title: t('navigation:sidebar.helpCenter'),
          url: '/help-center',
          icon: HelpCircle,
        },
      ],
    },
  ],
})
