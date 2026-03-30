import { useState, useEffect } from 'react'
import { Link } from 'react-router-dom'
import { Plus, RefreshCw, Server, Cloud } from 'lucide-react'
import { api } from '../api/client'
import StatusBadge from '../components/StatusBadge'
import { timeAgo } from '../lib/utils'

interface Cluster {
  id: string
  name: string
  display_name: string
  provider: string
  location: string
  status: string
  k8s_version: string
  node_count: number
  created_at: string
  last_synced_at: string | null
}

const providerIcons: Record<string, string> = {
  gke: 'Google Cloud',
  eks: 'AWS',
  aks: 'Azure',
}

export default function Clusters() {
  const [clusters, setClusters] = useState<Cluster[]>([])
  const [loading, setLoading] = useState(true)

  const fetchClusters = async () => {
    try {
      const data = await api.get('clusters').json<Cluster[]>()
      setClusters(data)
    } catch (err) {
      console.error('Failed to fetch clusters:', err)
    } finally {
      setLoading(false)
    }
  }

  useEffect(() => {
    fetchClusters()
    const interval = setInterval(fetchClusters, 15000)
    return () => clearInterval(interval)
  }, [])

  if (loading) {
    return <div className="text-text-secondary">Loading clusters...</div>
  }

  return (
    <div>
      <div className="flex items-center justify-between mb-8">
        <div>
          <h1 className="text-2xl font-bold">Clusters</h1>
          <p className="text-text-secondary text-sm mt-1">
            Manage your Kubernetes clusters and node pools
          </p>
        </div>
        <Link
          to="/clusters/connect"
          className="flex items-center gap-2 bg-accent text-black font-medium px-4 py-2 rounded-md hover:bg-accent/90 transition-colors"
        >
          <Plus className="w-4 h-4" />
          Connect Cluster
        </Link>
      </div>

      {clusters.length === 0 ? (
        <div className="text-center py-20">
          <Cloud className="w-12 h-12 text-text-secondary mx-auto mb-4" />
          <p className="text-text-secondary text-lg mb-2">No clusters connected</p>
          <p className="text-text-secondary text-sm mb-6">
            Connect your GKE cluster to manage nodes and deployments
          </p>
          <Link
            to="/clusters/connect"
            className="inline-flex items-center gap-2 bg-accent text-black font-medium px-4 py-2 rounded-md hover:bg-accent/90 transition-colors"
          >
            <Plus className="w-4 h-4" />
            Connect your first cluster
          </Link>
        </div>
      ) : (
        <div className="grid grid-cols-1 lg:grid-cols-2 gap-4">
          {clusters.map((cluster) => (
            <Link
              key={cluster.id}
              to={`/clusters/${cluster.id}`}
              className="block bg-surface border border-border rounded-lg p-5 hover:border-accent/30 transition-colors"
            >
              <div className="flex items-start justify-between mb-3">
                <div className="flex items-center gap-3">
                  <Server className="w-5 h-5 text-accent" />
                  <div>
                    <h3 className="font-medium text-text-primary">
                      {cluster.display_name || cluster.name}
                    </h3>
                    <p className="text-xs text-text-secondary font-mono">{cluster.name}</p>
                  </div>
                </div>
                <StatusBadge status={cluster.status} />
              </div>

              <div className="grid grid-cols-2 gap-3 text-xs text-text-secondary mt-4">
                <div>
                  <span className="block text-text-secondary/60">Provider</span>
                  <span className="text-text-primary">{providerIcons[cluster.provider] || cluster.provider.toUpperCase()}</span>
                </div>
                <div>
                  <span className="block text-text-secondary/60">Location</span>
                  <span className="text-text-primary">{cluster.location}</span>
                </div>
                <div>
                  <span className="block text-text-secondary/60">Nodes</span>
                  <span className="text-text-primary font-mono">{cluster.node_count}</span>
                </div>
                <div>
                  <span className="block text-text-secondary/60">K8s Version</span>
                  <span className="text-text-primary font-mono">{cluster.k8s_version}</span>
                </div>
              </div>

              {cluster.last_synced_at && (
                <div className="flex items-center gap-1 mt-3 text-xs text-text-secondary">
                  <RefreshCw className="w-3 h-3" />
                  Synced {timeAgo(cluster.last_synced_at)}
                </div>
              )}
            </Link>
          ))}
        </div>
      )}
    </div>
  )
}
