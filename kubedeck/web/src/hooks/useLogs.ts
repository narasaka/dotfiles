import { useEffect, useRef, useCallback } from 'react'

export function useLogs(url: string | null) {
  const wsRef = useRef<WebSocket | null>(null)
  const callbackRef = useRef<((data: string) => void) | null>(null)

  const connect = useCallback((onMessage: (data: string) => void) => {
    if (!url) return

    callbackRef.current = onMessage

    const protocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:'
    const wsUrl = `${protocol}//${window.location.host}${url}`

    const ws = new WebSocket(wsUrl)
    wsRef.current = ws

    ws.onmessage = (event) => {
      if (callbackRef.current) {
        callbackRef.current(event.data)
      }
    }

    ws.onerror = (err) => {
      console.error('WebSocket error:', err)
    }

    ws.onclose = () => {
      wsRef.current = null
    }
  }, [url])

  const disconnect = useCallback(() => {
    if (wsRef.current) {
      wsRef.current.close()
      wsRef.current = null
    }
  }, [])

  useEffect(() => {
    return () => disconnect()
  }, [disconnect])

  return { connect, disconnect }
}
