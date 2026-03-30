import { useState, useEffect } from 'react'
import { useParams, Link } from 'react-router-dom'
import { ArrowLeft, XCircle } from 'lucide-react'
import { api } from '../api/client'
import { toast } from 'sonner'
import { Build } from '../stores/appStore'
import StatusBadge from '../components/StatusBadge'
import LogViewer from '../components/LogViewer'
import { shortSha, formatDate } from '../lib/utils'

export default function BuildDetail() {
  const { appId, buildId } = useParams<{ appId: string; buildId: string }>()
  const [build, setBuild] = useState<Build | null>(null)

  useEffect(() => {
    const fetch = async () => {
      try {
        const data = await api.get(`builds/${buildId}`).json<Build>()
        setBuild(data)
      } catch {
        toast.error('Failed to load build')
      }
    }
    fetch()
    const interval = setInterval(fetch, 3000)
    return () => clearInterval(interval)
  }, [buildId])

  const cancelBuild = async () => {
    try {
      await api.post(`builds/${buildId}/cancel`)
      toast.success('Build cancelled')
    } catch {
      toast.error('Failed to cancel build')
    }
  }

  if (!build) {
    return <div className="text-text-secondary">Loading...</div>
  }

  const isActive = build.status === 'pending' || build.status === 'building'

  return (
    <div>
      <Link
        to={`/apps/${appId}`}
        className="flex items-center gap-1.5 text-sm text-text-secondary hover:text-text-primary mb-6 transition-colors"
      >
        <ArrowLeft className="w-4 h-4" />
        Back to app
      </Link>

      <div className="flex items-start justify-between mb-6">
        <div>
          <div className="flex items-center gap-3 mb-2">
            <h1 className="text-2xl font-bold font-mono">{shortSha(build.commit_sha)}</h1>
            <StatusBadge status={build.status} />
          </div>
          <div className="space-y-1 text-sm text-text-secondary">
            {build.commit_message && <p>{build.commit_message}</p>}
            <div className="flex items-center gap-4">
              {build.commit_author && <span>by {build.commit_author}</span>}
              <span>Started: {formatDate(build.started_at)}</span>
              {build.finished_at && <span>Finished: {formatDate(build.finished_at)}</span>}
            </div>
          </div>
        </div>

        {isActive && (
          <button
            onClick={cancelBuild}
            className="flex items-center gap-2 text-red-400 hover:text-red-300 border border-red-400/30 px-4 py-2 rounded-md transition-colors"
          >
            <XCircle className="w-4 h-4" />
            Cancel
          </button>
        )}
      </div>

      <LogViewer
        wsUrl={isActive ? `/api/v1/builds/${buildId}/logs/ws` : undefined}
        staticLogs={!isActive ? build.logs : undefined}
      />
    </div>
  )
}
