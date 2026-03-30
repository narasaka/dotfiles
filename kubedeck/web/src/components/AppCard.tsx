import { Link } from 'react-router-dom'
import { GitBranch, Clock } from 'lucide-react'
import { App } from '../stores/appStore'
import StatusBadge from './StatusBadge'
import { timeAgo } from '../lib/utils'

export default function AppCard({ app }: { app: App }) {
  return (
    <Link
      to={`/apps/${app.id}`}
      className="block bg-surface border border-border rounded-lg p-5 hover:border-accent/30 transition-colors"
    >
      <div className="flex items-start justify-between mb-3">
        <div>
          <h3 className="font-medium text-text-primary">
            {app.display_name || app.name}
          </h3>
          <p className="text-sm text-text-secondary font-mono">{app.name}</p>
        </div>
        <StatusBadge status={app.status} />
      </div>

      <div className="flex items-center gap-4 text-xs text-text-secondary">
        <span className="flex items-center gap-1">
          <GitBranch className="w-3 h-3" />
          {app.git_branch}
        </span>
        <span className="flex items-center gap-1">
          <Clock className="w-3 h-3" />
          {timeAgo(app.updated_at)}
        </span>
      </div>
    </Link>
  )
}
