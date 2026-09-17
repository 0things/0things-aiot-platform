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
  const { getIdTokenClaims } = useLogto()
  const [profile, setProfile] = useState<{
    name?: string
    email?: string
    picture?: string
  }>({})

  useEffect(() => {
    let active = true
    void getIdTokenClaims().then((claims) => {
      if (!active || !claims) return
      setProfile({
        name: claims.name ?? undefined,
        email: claims.email ?? undefined,
        picture: claims.picture ?? undefined,
      })
    })
    return () => {
      active = false
    }
  }, [getIdTokenClaims])

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
