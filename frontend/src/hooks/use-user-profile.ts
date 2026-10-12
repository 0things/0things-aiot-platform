import { useEffect, useState } from 'react'
import { useLogto } from '@logto/react'

export interface UserProfileInfo {
  displayName: string
  email?: string
  initials: string
  avatarSrc?: string
}

export function useUserProfile(fallback?: {
  name?: string
  email?: string
  avatar?: string
}): UserProfileInfo {
  const { fetchUserInfo } = useLogto()
  const [profile, setProfile] = useState<{
    name?: string
    email?: string
    picture?: string
  }>({})

  useEffect(() => {
    let active = true
    void fetchUserInfo()
      .catch(() => null)
      .then((userInfo) => {
        if (!active) return
        setProfile({
          name: userInfo?.name ?? userInfo?.username ?? undefined,
          email: userInfo?.email ?? undefined,
          picture: userInfo?.picture ?? undefined,
        })
      })
    return () => {
      active = false
    }
  }, [fetchUserInfo])

  const name = profile.name || fallback?.name
  const email = profile.email || fallback?.email
  const displayName = name || email || 'User'
  const initials = (name || email || 'U').slice(0, 2).toUpperCase()
  const avatarSrc = profile.picture || fallback?.avatar

  return {
    displayName,
    email,
    initials,
    avatarSrc,
  }
}
