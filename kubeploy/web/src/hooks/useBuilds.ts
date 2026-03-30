import { useState, useEffect, useCallback } from 'react'
import { api } from '../api/client'
import { Build } from '../stores/appStore'

export function useBuilds(appId: string) {
  const [builds, setBuilds] = useState<Build[]>([])
  const [loading, setLoading] = useState(true)

  const fetchBuilds = useCallback(async () => {
    try {
      const data = await api.get(`apps/${appId}/builds`).json<Build[]>()
      setBuilds(data)
    } catch (err) {
      console.error('Failed to fetch builds:', err)
    } finally {
      setLoading(false)
    }
  }, [appId])

  useEffect(() => {
    fetchBuilds()
    const interval = setInterval(fetchBuilds, 5000)
    return () => clearInterval(interval)
  }, [fetchBuilds])

  return { builds, loading, refetch: fetchBuilds }
}
