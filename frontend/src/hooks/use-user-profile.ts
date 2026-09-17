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
  const { fetchUserInfo, getIdTokenClaims } = useLogto()
  const [profile, setProfile] = useState<{
    name?: string
    email?: string
    picture?: string
  }>({})

  useEffect(() => {
    let active = true
    void Promise.all([
      fetchUserInfo().catch(() => null),
      getIdTokenClaims().catch(() => null),
    ]).then(([userInfo, claims]) => {
      if (!active) return
      setProfile({
        name: userInfo?.name ?? userInfo?.username ?? claims?.name ?? undefined,
        email: userInfo?.email ?? claims?.email ?? undefined,
        picture: userInfo?.picture ?? claims?.picture ?? undefined,
      })
    })
    return () => {
      active = false
    }
  }, [fetchUserInfo, getIdTokenClaims])

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
