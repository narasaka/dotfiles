import { Link } from 'react-router-dom'
import { Build } from '../stores/appStore'
import StatusBadge from './StatusBadge'
import { shortSha, timeAgo, formatDate } from '../lib/utils'

interface Props {
  builds: Build[]
  appId: string
}

export default function BuildList({ builds, appId }: Props) {
  if (builds.length === 0) {
    return (
      <div className="text-center py-12 text-text-secondary">
        No builds yet. Trigger a build or push to your repository.
      </div>
    )
  }

  return (
    <div className="space-y-2">
      {builds.map((build) => (
        <Link
          key={build.id}
          to={`/apps/${appId}/builds/${build.id}`}
          className="flex items-center justify-between p-4 bg-surface border border-border rounded-lg hover:border-accent/30 transition-colors"
        >
          <div className="flex items-center gap-4">
            <StatusBadge status={build.status} />
            <div>
              <div className="flex items-center gap-2">
                <span className="font-mono text-sm text-accent">
                  {shortSha(build.commit_sha)}
                </span>
                <span className="text-sm text-text-primary truncate max-w-[300px]">
                  {build.commit_message || 'Manual build'}
                </span>
              </div>
              <div className="text-xs text-text-secondary mt-1">
                {build.commit_author && `${build.commit_author} · `}
                {timeAgo(build.created_at)}
              </div>
            </div>
          </div>
          <div className="text-xs text-text-secondary">
            {build.started_at && build.finished_at
              ? `${Math.round((new Date(build.finished_at).getTime() - new Date(build.started_at).getTime()) / 1000)}s`
              : build.started_at
                ? 'Running...'
                : 'Queued'}
          </div>
        </Link>
      ))}
    </div>
  )
}
