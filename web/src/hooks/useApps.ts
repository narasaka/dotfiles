import { useEffect, useCallback } from 'react'
import { api } from '../api/client'
import { useAppStore, App } from '../stores/appStore'

export function useApps() {
  const { apps, setApps } = useAppStore()

  const fetchApps = useCallback(async () => {
    try {
      const data = await api.get('apps').json<App[]>()
      setApps(data)
    } catch (err) {
      console.error('Failed to fetch apps:', err)
    }
  }, [setApps])

  useEffect(() => {
    fetchApps()
    const interval = setInterval(fetchApps, 10000)
    return () => clearInterval(interval)
  }, [fetchApps])

  return { apps, refetch: fetchApps }
}
