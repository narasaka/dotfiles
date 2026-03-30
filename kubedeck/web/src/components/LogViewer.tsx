import { useEffect, useRef } from 'react'
import { Terminal } from 'xterm'
import { FitAddon } from 'xterm-addon-fit'
import 'xterm/css/xterm.css'

interface Props {
  wsUrl?: string | null
  staticLogs?: string
}

export default function LogViewer({ wsUrl, staticLogs }: Props) {
  const containerRef = useRef<HTMLDivElement>(null)
  const termRef = useRef<Terminal | null>(null)

  useEffect(() => {
    if (!containerRef.current) return

    const term = new Terminal({
      theme: {
        background: '#0A0A0A',
        foreground: '#FAFAFA',
        cursor: '#22D3EE',
        selectionBackground: '#22D3EE33',
        black: '#0A0A0A',
        green: '#10B981',
        red: '#EF4444',
        yellow: '#F59E0B',
        cyan: '#22D3EE',
      },
      fontFamily: '"JetBrains Mono", monospace',
      fontSize: 13,
      lineHeight: 1.4,
      cursorBlink: false,
      disableStdin: true,
      convertEol: true,
    })

    const fitAddon = new FitAddon()
    term.loadAddon(fitAddon)
    term.open(containerRef.current)
    fitAddon.fit()
    termRef.current = term

    const handleResize = () => fitAddon.fit()
    window.addEventListener('resize', handleResize)

    // Write static logs if provided
    if (staticLogs) {
      term.write(staticLogs)
    }

    // Connect WebSocket if URL provided
    let ws: WebSocket | null = null
    if (wsUrl) {
      const protocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:'
      ws = new WebSocket(`${protocol}//${window.location.host}${wsUrl}`)
      ws.onmessage = (event) => {
        term.write(event.data)
      }
      ws.onerror = () => {
        term.write('\r\n\x1b[31mWebSocket connection error\x1b[0m\r\n')
      }
      ws.onclose = () => {
        term.write('\r\n\x1b[33m--- Stream ended ---\x1b[0m\r\n')
      }
    }

    return () => {
      window.removeEventListener('resize', handleResize)
      ws?.close()
      term.dispose()
    }
  }, [wsUrl, staticLogs])

  return (
    <div
      ref={containerRef}
      className="h-[500px] rounded-lg border border-border overflow-hidden"
    />
  )
}
