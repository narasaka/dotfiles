import { Deployment } from '../stores/appStore'
import StatusBadge from './StatusBadge'
import { formatDate } from '../lib/utils'

interface Props {
  deployment: Deployment | null
}

export default function DeploymentStatus({ deployment }: Props) {
  if (!deployment) {
    return (
      <div className="text-center py-12 text-text-secondary">
        No deployments yet.
      </div>
    )
  }

  return (
    <div className="bg-surface border border-border rounded-lg p-6">
      <div className="flex items-center justify-between mb-4">
        <h3 className="font-medium">Current Deployment</h3>
        <StatusBadge status={deployment.status} />
      </div>

      <div className="grid grid-cols-2 gap-4 text-sm">
        <div>
          <span className="text-text-secondary">Replicas</span>
          <p className="font-mono mt-1">
            {deployment.replicas_ready}/{deployment.replicas_desired}
          </p>
        </div>
        <div>
          <span className="text-text-secondary">Deployed</span>
          <p className="mt-1">{formatDate(deployment.created_at)}</p>
        </div>
        <div>
          <span className="text-text-secondary">K8s Deployment</span>
          <p className="font-mono mt-1">{deployment.k8s_deployment_name}</p>
        </div>
      </div>
    </div>
  )
}
