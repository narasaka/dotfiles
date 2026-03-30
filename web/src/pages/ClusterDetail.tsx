import { useState, useEffect } from 'react'
import { useParams, useNavigate } from 'react-router-dom'
import { RefreshCw, Plus, Trash2, Server, Cpu, HardDrive, ArrowUpDown, Clock } from 'lucide-react'
import { api } from '../api/client'
import { toast } from 'sonner'
import StatusBadge from '../components/StatusBadge'
import { formatDate, timeAgo } from '../lib/utils'

interface Cluster {
  id: string
  name: string
  display_name: string
  provider: string
  location: string
  status: string
  k8s_version: string
  endpoint: string
  node_count: number
  last_synced_at: string | null
  created_at: string
}

interface NodePool {
  id: string
  name: string
  machine_type: string
  disk_size_gb: number
  disk_type: string
  node_count: number
  min_nodes: number
  max_nodes: number
  autoscaling: boolean
  preemptible: boolean
  spot_instance: boolean
  status: string
  k8s_version: string
}

interface ClusterEvent {
  id: string
  event_type: string
  message: string
  created_at: string
}

interface ClusterMetrics {
  total_nodes: number
  ready_nodes: number
  total_pods: number
  running_pods: number
}

type Tab = 'overview' | 'node-pools' | 'events'

export default function ClusterDetail() {
  const { id } = useParams<{ id: string }>()
  const navigate = useNavigate()
  const [tab, setTab] = useState<Tab>('overview')
  const [cluster, setCluster] = useState<Cluster | null>(null)
  const [nodePools, setNodePools] = useState<NodePool[]>([])
  const [events, setEvents] = useState<ClusterEvent[]>([])
  const [metrics, setMetrics] = useState<ClusterMetrics | null>(null)
  const [syncing, setSyncing] = useState(false)

  // Scale modal state
  const [scalePool, setScalePool] = useState<NodePool | null>(null)
  const [scaleCount, setScaleCount] = useState(0)
  const [scaleMin, setScaleMin] = useState(0)
  const [scaleMax, setScaleMax] = useState(0)

  // Create pool modal
  const [showCreate, setShowCreate] = useState(false)
  const [newPool, setNewPool] = useState({
    name: '',
    machine_type: 'e2-medium',
    disk_size_gb: 100,
    initial_node_count: 1,
    min_nodes: 1,
    max_nodes: 3,
    autoscaling: true,
    preemptible: false,
    spot_instance: false,
  })

  const fetchData = async () => {
    try {
      const [c, pools, evts] = await Promise.all([
        api.get(`clusters/${id}`).json<Cluster>(),
        api.get(`clusters/${id}/node-pools`).json<NodePool[]>(),
        api.get(`clusters/${id}/events?limit=20`).json<ClusterEvent[]>(),
      ])
      setCluster(c)
      setNodePools(pools)
      setEvents(evts)

      try {
        const m = await api.get(`clusters/${id}/metrics`).json<ClusterMetrics>()
        setMetrics(m)
      } catch {}
    } catch {
      toast.error('Failed to load cluster')
      navigate('/clusters')
    }
  }

  useEffect(() => {
    fetchData()
    const interval = setInterval(fetchData, 10000)
    return () => clearInterval(interval)
  }, [id])

  const handleSync = async () => {
    setSyncing(true)
    try {
      await api.post(`clusters/${id}/sync`).json()
      toast.success('Cluster synced')
      fetchData()
    } catch {
      toast.error('Sync failed')
    } finally {
      setSyncing(false)
    }
  }

  const handleDelete = async () => {
    if (!confirm('Disconnect this cluster? This will not delete the actual cluster.')) return
    try {
      await api.delete(`clusters/${id}`)
      toast.success('Cluster disconnected')
      navigate('/clusters')
    } catch {
      toast.error('Failed to disconnect')
    }
  }

  const handleScale = async () => {
    if (!scalePool) return
    try {
      await api.put(`clusters/${id}/node-pools/${scalePool.id}`, {
        json: {
          node_count: scaleCount,
          min_nodes: scaleMin,
          max_nodes: scaleMax,
          autoscaling: scalePool.autoscaling,
        },
      }).json()
      toast.success(`Node pool "${scalePool.name}" scaling to ${scaleCount} nodes`)
      setScalePool(null)
      fetchData()
    } catch {
      toast.error('Failed to scale node pool')
    }
  }

  const handleCreatePool = async () => {
    try {
      await api.post(`clusters/${id}/node-pools`, { json: newPool }).json()
      toast.success('Node pool created')
      setShowCreate(false)
      setNewPool({ name: '', machine_type: 'e2-medium', disk_size_gb: 100, initial_node_count: 1, min_nodes: 1, max_nodes: 3, autoscaling: true, preemptible: false, spot_instance: false })
      fetchData()
    } catch {
      toast.error('Failed to create node pool')
    }
  }

  const handleDeletePool = async (pool: NodePool) => {
    if (!confirm(`Delete node pool "${pool.name}"? This will drain and remove all nodes.`)) return
    try {
      await api.delete(`clusters/${id}/node-pools/${pool.id}`)
      toast.success('Node pool deleted')
      fetchData()
    } catch {
      toast.error('Failed to delete node pool')
    }
  }

  if (!cluster) return <div className="text-text-secondary">Loading...</div>

  return (
    <div>
      {/* Header */}
      <div className="flex items-start justify-between mb-8">
        <div>
          <div className="flex items-center gap-3">
            <Server className="w-6 h-6 text-accent" />
            <h1 className="text-2xl font-bold">{cluster.display_name || cluster.name}</h1>
            <StatusBadge status={cluster.status} />
          </div>
          <div className="flex items-center gap-4 mt-2 text-sm text-text-secondary">
            <span>{cluster.location}</span>
            <span className="font-mono">v{cluster.k8s_version}</span>
            <span>{cluster.node_count} nodes</span>
          </div>
        </div>
        <div className="flex items-center gap-2">
          <button
            onClick={handleSync}
            disabled={syncing}
            className="flex items-center gap-2 border border-border px-3 py-2 rounded-md text-sm hover:bg-white/5 transition-colors disabled:opacity-50"
          >
            <RefreshCw className={`w-4 h-4 ${syncing ? 'animate-spin' : ''}`} />
            Sync
          </button>
          <button
            onClick={handleDelete}
            className="flex items-center gap-2 border border-red-400/30 text-red-400 px-3 py-2 rounded-md text-sm hover:bg-red-400/10 transition-colors"
          >
            <Trash2 className="w-4 h-4" />
            Disconnect
          </button>
        </div>
      </div>

      {/* Tabs */}
      <div className="flex gap-1 border-b border-border mb-6">
        {(['overview', 'node-pools', 'events'] as Tab[]).map((t) => (
          <button
            key={t}
            onClick={() => setTab(t)}
            className={`px-4 py-2.5 text-sm font-medium transition-colors border-b-2 -mb-px ${
              tab === t ? 'border-accent text-accent' : 'border-transparent text-text-secondary hover:text-text-primary'
            }`}
          >
            {t === 'node-pools' ? 'Node Pools' : t.charAt(0).toUpperCase() + t.slice(1)}
          </button>
        ))}
      </div>

      {/* Overview */}
      {tab === 'overview' && (
        <div className="space-y-6">
          {/* Metrics cards */}
          <div className="grid grid-cols-2 lg:grid-cols-4 gap-4">
            <div className="bg-surface border border-border rounded-lg p-4">
              <div className="flex items-center gap-2 text-text-secondary text-sm mb-2">
                <Server className="w-4 h-4" />
                Nodes
              </div>
              <p className="text-2xl font-bold font-mono">
                {metrics?.ready_nodes ?? cluster.node_count}/{metrics?.total_nodes ?? cluster.node_count}
              </p>
            </div>
            <div className="bg-surface border border-border rounded-lg p-4">
              <div className="flex items-center gap-2 text-text-secondary text-sm mb-2">
                <HardDrive className="w-4 h-4" />
                Node Pools
              </div>
              <p className="text-2xl font-bold font-mono">{nodePools.length}</p>
            </div>
            <div className="bg-surface border border-border rounded-lg p-4">
              <div className="flex items-center gap-2 text-text-secondary text-sm mb-2">
                <Cpu className="w-4 h-4" />
                Pods
              </div>
              <p className="text-2xl font-bold font-mono">
                {metrics?.running_pods ?? '-'}/{metrics?.total_pods ?? '-'}
              </p>
            </div>
            <div className="bg-surface border border-border rounded-lg p-4">
              <div className="flex items-center gap-2 text-text-secondary text-sm mb-2">
                <Clock className="w-4 h-4" />
                Last Synced
              </div>
              <p className="text-sm font-mono mt-1">
                {cluster.last_synced_at ? timeAgo(cluster.last_synced_at) : 'Never'}
              </p>
            </div>
          </div>

          {/* Cluster details */}
          <div className="bg-surface border border-border rounded-lg p-6">
            <h3 className="font-medium mb-4">Cluster Details</h3>
            <div className="grid grid-cols-2 gap-4 text-sm">
              <div>
                <span className="text-text-secondary">Endpoint</span>
                <p className="font-mono mt-1 text-xs break-all">{cluster.endpoint}</p>
              </div>
              <div>
                <span className="text-text-secondary">Provider</span>
                <p className="mt-1">{cluster.provider.toUpperCase()}</p>
              </div>
              <div>
                <span className="text-text-secondary">Created</span>
                <p className="mt-1">{formatDate(cluster.created_at)}</p>
              </div>
              <div>
                <span className="text-text-secondary">K8s Version</span>
                <p className="font-mono mt-1">{cluster.k8s_version}</p>
              </div>
            </div>
          </div>
        </div>
      )}

      {/* Node Pools */}
      {tab === 'node-pools' && (
        <div className="space-y-4">
          <div className="flex justify-end">
            <button
              onClick={() => setShowCreate(true)}
              className="flex items-center gap-2 bg-accent text-black font-medium px-4 py-2 rounded-md hover:bg-accent/90 transition-colors"
            >
              <Plus className="w-4 h-4" />
              Add Node Pool
            </button>
          </div>

          {nodePools.length === 0 ? (
            <div className="text-center py-12 text-text-secondary">
              No node pools found. Sync the cluster or add a new pool.
            </div>
          ) : (
            nodePools.map((pool) => (
              <div key={pool.id} className="bg-surface border border-border rounded-lg p-5">
                <div className="flex items-start justify-between mb-4">
                  <div>
                    <div className="flex items-center gap-3">
                      <h3 className="font-medium">{pool.name}</h3>
                      <StatusBadge status={pool.status.toLowerCase()} />
                    </div>
                    <p className="text-xs text-text-secondary font-mono mt-1">{pool.machine_type}</p>
                  </div>
                  <div className="flex items-center gap-2">
                    <button
                      onClick={() => {
                        setScalePool(pool)
                        setScaleCount(pool.node_count)
                        setScaleMin(pool.min_nodes)
                        setScaleMax(pool.max_nodes)
                      }}
                      className="flex items-center gap-1 text-xs border border-border px-2.5 py-1.5 rounded-md hover:bg-white/5 transition-colors"
                    >
                      <ArrowUpDown className="w-3 h-3" />
                      Scale
                    </button>
                    <button
                      onClick={() => handleDeletePool(pool)}
                      className="flex items-center gap-1 text-xs border border-red-400/30 text-red-400 px-2.5 py-1.5 rounded-md hover:bg-red-400/10 transition-colors"
                    >
                      <Trash2 className="w-3 h-3" />
                    </button>
                  </div>
                </div>

                <div className="grid grid-cols-4 gap-4 text-xs">
                  <div>
                    <span className="text-text-secondary">Nodes</span>
                    <p className="font-mono mt-1 text-text-primary">{pool.node_count}</p>
                  </div>
                  <div>
                    <span className="text-text-secondary">Autoscaling</span>
                    <p className="font-mono mt-1 text-text-primary">
                      {pool.autoscaling ? `${pool.min_nodes}-${pool.max_nodes}` : 'Off'}
                    </p>
                  </div>
                  <div>
                    <span className="text-text-secondary">Disk</span>
                    <p className="font-mono mt-1 text-text-primary">{pool.disk_size_gb}GB {pool.disk_type}</p>
                  </div>
                  <div>
                    <span className="text-text-secondary">Type</span>
                    <p className="font-mono mt-1 text-text-primary">
                      {pool.spot_instance ? 'Spot' : pool.preemptible ? 'Preemptible' : 'On-demand'}
                    </p>
                  </div>
                </div>
              </div>
            ))
          )}

          {/* Scale Modal */}
          {scalePool && (
            <div className="fixed inset-0 bg-black/50 flex items-center justify-center z-50" onClick={() => setScalePool(null)}>
              <div className="bg-surface border border-border rounded-lg p-6 w-96" onClick={(e) => e.stopPropagation()}>
                <h3 className="font-medium mb-4">Scale "{scalePool.name}"</h3>
                <div className="space-y-4">
                  <div>
                    <label className="block text-sm font-medium mb-1.5">Node Count</label>
                    <input
                      type="number"
                      value={scaleCount}
                      onChange={(e) => setScaleCount(parseInt(e.target.value) || 0)}
                      min={0}
                      className="w-full bg-[#0A0A0A] border border-border rounded-md px-3 py-2 text-sm font-mono focus:outline-none focus:border-accent"
                    />
                  </div>
                  {scalePool.autoscaling && (
                    <div className="grid grid-cols-2 gap-4">
                      <div>
                        <label className="block text-sm font-medium mb-1.5">Min Nodes</label>
                        <input
                          type="number"
                          value={scaleMin}
                          onChange={(e) => setScaleMin(parseInt(e.target.value) || 0)}
                          min={0}
                          className="w-full bg-[#0A0A0A] border border-border rounded-md px-3 py-2 text-sm font-mono focus:outline-none focus:border-accent"
                        />
                      </div>
                      <div>
                        <label className="block text-sm font-medium mb-1.5">Max Nodes</label>
                        <input
                          type="number"
                          value={scaleMax}
                          onChange={(e) => setScaleMax(parseInt(e.target.value) || 0)}
                          min={1}
                          className="w-full bg-[#0A0A0A] border border-border rounded-md px-3 py-2 text-sm font-mono focus:outline-none focus:border-accent"
                        />
                      </div>
                    </div>
                  )}
                  <div className="flex gap-3">
                    <button onClick={handleScale} className="bg-accent text-black font-medium px-4 py-2 rounded-md hover:bg-accent/90 transition-colors">
                      Apply
                    </button>
                    <button onClick={() => setScalePool(null)} className="px-4 py-2 rounded-md border border-border text-text-secondary hover:text-text-primary transition-colors">
                      Cancel
                    </button>
                  </div>
                </div>
              </div>
            </div>
          )}

          {/* Create Pool Modal */}
          {showCreate && (
            <div className="fixed inset-0 bg-black/50 flex items-center justify-center z-50" onClick={() => setShowCreate(false)}>
              <div className="bg-surface border border-border rounded-lg p-6 w-[500px] max-h-[80vh] overflow-y-auto" onClick={(e) => e.stopPropagation()}>
                <h3 className="font-medium mb-4">Create Node Pool</h3>
                <div className="space-y-4">
                  <div>
                    <label className="block text-sm font-medium mb-1.5">Name</label>
                    <input
                      type="text"
                      value={newPool.name}
                      onChange={(e) => setNewPool({ ...newPool, name: e.target.value })}
                      className="w-full bg-[#0A0A0A] border border-border rounded-md px-3 py-2 text-sm font-mono focus:outline-none focus:border-accent"
                      placeholder="worker-pool"
                    />
                  </div>
                  <div className="grid grid-cols-2 gap-4">
                    <div>
                      <label className="block text-sm font-medium mb-1.5">Machine Type</label>
                      <select
                        value={newPool.machine_type}
                        onChange={(e) => setNewPool({ ...newPool, machine_type: e.target.value })}
                        className="w-full bg-[#0A0A0A] border border-border rounded-md px-3 py-2 text-sm focus:outline-none focus:border-accent"
                      >
                        <option value="e2-micro">e2-micro (2 vCPU, 1GB)</option>
                        <option value="e2-small">e2-small (2 vCPU, 2GB)</option>
                        <option value="e2-medium">e2-medium (2 vCPU, 4GB)</option>
                        <option value="e2-standard-2">e2-standard-2 (2 vCPU, 8GB)</option>
                        <option value="e2-standard-4">e2-standard-4 (4 vCPU, 16GB)</option>
                        <option value="e2-standard-8">e2-standard-8 (8 vCPU, 32GB)</option>
                        <option value="e2-standard-16">e2-standard-16 (16 vCPU, 64GB)</option>
                        <option value="n2-standard-2">n2-standard-2 (2 vCPU, 8GB)</option>
                        <option value="n2-standard-4">n2-standard-4 (4 vCPU, 16GB)</option>
                        <option value="n2-standard-8">n2-standard-8 (8 vCPU, 32GB)</option>
                      </select>
                    </div>
                    <div>
                      <label className="block text-sm font-medium mb-1.5">Disk Size (GB)</label>
                      <input
                        type="number"
                        value={newPool.disk_size_gb}
                        onChange={(e) => setNewPool({ ...newPool, disk_size_gb: parseInt(e.target.value) || 100 })}
                        className="w-full bg-[#0A0A0A] border border-border rounded-md px-3 py-2 text-sm font-mono focus:outline-none focus:border-accent"
                      />
                    </div>
                  </div>
                  <div className="grid grid-cols-3 gap-4">
                    <div>
                      <label className="block text-sm font-medium mb-1.5">Initial Nodes</label>
                      <input
                        type="number"
                        value={newPool.initial_node_count}
                        onChange={(e) => setNewPool({ ...newPool, initial_node_count: parseInt(e.target.value) || 1 })}
                        min={1}
                        className="w-full bg-[#0A0A0A] border border-border rounded-md px-3 py-2 text-sm font-mono focus:outline-none focus:border-accent"
                      />
                    </div>
                    <div>
                      <label className="block text-sm font-medium mb-1.5">Min Nodes</label>
                      <input
                        type="number"
                        value={newPool.min_nodes}
                        onChange={(e) => setNewPool({ ...newPool, min_nodes: parseInt(e.target.value) || 0 })}
                        min={0}
                        className="w-full bg-[#0A0A0A] border border-border rounded-md px-3 py-2 text-sm font-mono focus:outline-none focus:border-accent"
                      />
                    </div>
                    <div>
                      <label className="block text-sm font-medium mb-1.5">Max Nodes</label>
                      <input
                        type="number"
                        value={newPool.max_nodes}
                        onChange={(e) => setNewPool({ ...newPool, max_nodes: parseInt(e.target.value) || 3 })}
                        min={1}
                        className="w-full bg-[#0A0A0A] border border-border rounded-md px-3 py-2 text-sm font-mono focus:outline-none focus:border-accent"
                      />
                    </div>
                  </div>
                  <div className="flex items-center gap-6">
                    <label className="flex items-center gap-2 text-sm">
                      <input type="checkbox" checked={newPool.autoscaling} onChange={(e) => setNewPool({ ...newPool, autoscaling: e.target.checked })} />
                      Autoscaling
                    </label>
                    <label className="flex items-center gap-2 text-sm">
                      <input type="checkbox" checked={newPool.spot_instance} onChange={(e) => setNewPool({ ...newPool, spot_instance: e.target.checked, preemptible: false })} />
                      Spot Instances
                    </label>
                    <label className="flex items-center gap-2 text-sm">
                      <input type="checkbox" checked={newPool.preemptible} onChange={(e) => setNewPool({ ...newPool, preemptible: e.target.checked, spot_instance: false })} />
                      Preemptible
                    </label>
                  </div>
                  <div className="flex gap-3 pt-2">
                    <button onClick={handleCreatePool} disabled={!newPool.name} className="bg-accent text-black font-medium px-4 py-2 rounded-md hover:bg-accent/90 transition-colors disabled:opacity-50">
                      Create
                    </button>
                    <button onClick={() => setShowCreate(false)} className="px-4 py-2 rounded-md border border-border text-text-secondary hover:text-text-primary transition-colors">
                      Cancel
                    </button>
                  </div>
                </div>
              </div>
            </div>
          )}
        </div>
      )}

      {/* Events */}
      {tab === 'events' && (
        <div className="space-y-2">
          {events.length === 0 ? (
            <div className="text-center py-12 text-text-secondary">No events yet.</div>
          ) : (
            events.map((event) => (
              <div key={event.id} className="flex items-start gap-3 p-3 bg-surface border border-border rounded-lg">
                <div className="w-2 h-2 rounded-full mt-1.5 bg-accent flex-shrink-0" />
                <div>
                  <p className="text-sm">{event.message}</p>
                  <div className="flex items-center gap-2 mt-1 text-xs text-text-secondary">
                    <span className="font-mono">{event.event_type}</span>
                    <span>{timeAgo(event.created_at)}</span>
                  </div>
                </div>
              </div>
            ))
          )}
        </div>
      )}
    </div>
  )
}
