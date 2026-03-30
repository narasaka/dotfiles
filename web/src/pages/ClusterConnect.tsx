import { useState } from 'react'
import { useNavigate } from 'react-router-dom'
import { api } from '../api/client'
import { toast } from 'sonner'
import { ArrowLeft, Check, Cloud, Loader2 } from 'lucide-react'

interface DiscoveredCluster {
  id: string
  name: string
  location: string
  status: string
  k8s_version: string
  node_count: number
}

export default function ClusterConnect() {
  const navigate = useNavigate()
  const [step, setStep] = useState<'credentials' | 'select' | 'confirm'>('credentials')
  const [provider] = useState('gke')
  const [projectId, setProjectId] = useState('')
  const [serviceAccountJson, setServiceAccountJson] = useState('')
  const [validating, setValidating] = useState(false)
  const [discovering, setDiscovering] = useState(false)
  const [clusters, setClusters] = useState<DiscoveredCluster[]>([])
  const [selected, setSelected] = useState<DiscoveredCluster | null>(null)
  const [displayName, setDisplayName] = useState('')
  const [connecting, setConnecting] = useState(false)

  const creds = {
    project_id: projectId,
    service_account_json: serviceAccountJson,
  }

  const handleValidate = async () => {
    setValidating(true)
    try {
      await api.post(`providers/${provider}/validate`, { json: creds }).json()
      toast.success('Credentials validated')

      // Discover clusters
      setDiscovering(true)
      const data = await api.post(`providers/${provider}/discover`, { json: creds }).json<DiscoveredCluster[]>()
      setClusters(data || [])
      setStep('select')
    } catch (err: any) {
      toast.error('Invalid credentials')
    } finally {
      setValidating(false)
      setDiscovering(false)
    }
  }

  const handleSelect = (cluster: DiscoveredCluster) => {
    setSelected(cluster)
    setDisplayName(cluster.name)
    setStep('confirm')
  }

  const handleConnect = async () => {
    if (!selected) return
    setConnecting(true)
    try {
      const result = await api.post('clusters', {
        json: {
          name: selected.name,
          display_name: displayName,
          provider,
          provider_cluster_id: selected.id,
          project_id: projectId,
          credentials: creds,
        },
      }).json<{ id: string }>()
      toast.success('Cluster connected')
      navigate(`/clusters/${result.id}`)
    } catch (err: any) {
      toast.error('Failed to connect cluster')
    } finally {
      setConnecting(false)
    }
  }

  return (
    <div className="max-w-2xl">
      <button
        onClick={() => navigate('/clusters')}
        className="flex items-center gap-1.5 text-sm text-text-secondary hover:text-text-primary mb-6 transition-colors"
      >
        <ArrowLeft className="w-4 h-4" />
        Back to clusters
      </button>

      <h1 className="text-2xl font-bold mb-2">Connect Cluster</h1>
      <p className="text-text-secondary text-sm mb-8">
        Connect an existing GKE cluster to manage it through Kubeploy.
      </p>

      {/* Step indicator */}
      <div className="flex items-center gap-2 mb-8">
        {['Credentials', 'Select Cluster', 'Confirm'].map((label, i) => {
          const stepIndex = ['credentials', 'select', 'confirm'].indexOf(step)
          const isActive = i === stepIndex
          const isDone = i < stepIndex
          return (
            <div key={label} className="flex items-center gap-2">
              {i > 0 && <div className={`w-8 h-px ${isDone ? 'bg-accent' : 'bg-border'}`} />}
              <div className={`flex items-center gap-2 text-sm ${isActive ? 'text-accent' : isDone ? 'text-emerald-400' : 'text-text-secondary'}`}>
                <div className={`w-6 h-6 rounded-full flex items-center justify-center text-xs border ${isActive ? 'border-accent bg-accent/10' : isDone ? 'border-emerald-400 bg-emerald-400/10' : 'border-border'}`}>
                  {isDone ? <Check className="w-3 h-3" /> : i + 1}
                </div>
                {label}
              </div>
            </div>
          )
        })}
      </div>

      {/* Step 1: Credentials */}
      {step === 'credentials' && (
        <div className="space-y-6">
          <div className="bg-surface border border-border rounded-lg p-5">
            <div className="flex items-center gap-3 mb-4">
              <Cloud className="w-5 h-5 text-accent" />
              <h3 className="font-medium">Google Kubernetes Engine (GKE)</h3>
            </div>

            <div className="space-y-4">
              <div>
                <label className="block text-sm font-medium mb-1.5">GCP Project ID</label>
                <input
                  type="text"
                  value={projectId}
                  onChange={(e) => setProjectId(e.target.value)}
                  className="w-full bg-[#0A0A0A] border border-border rounded-md px-3 py-2 text-sm font-mono focus:outline-none focus:border-accent"
                  placeholder="my-gcp-project"
                />
              </div>
              <div>
                <label className="block text-sm font-medium mb-1.5">Service Account JSON Key</label>
                <textarea
                  value={serviceAccountJson}
                  onChange={(e) => setServiceAccountJson(e.target.value)}
                  rows={8}
                  className="w-full bg-[#0A0A0A] border border-border rounded-md px-3 py-2 text-sm font-mono focus:outline-none focus:border-accent resize-none"
                  placeholder='{"type": "service_account", ...}'
                />
                <p className="text-xs text-text-secondary mt-1">
                  Requires roles: Kubernetes Engine Admin, Compute Viewer
                </p>
              </div>
            </div>
          </div>

          <button
            onClick={handleValidate}
            disabled={validating || !projectId || !serviceAccountJson}
            className="flex items-center gap-2 bg-accent text-black font-medium px-6 py-2 rounded-md hover:bg-accent/90 transition-colors disabled:opacity-50"
          >
            {validating ? <Loader2 className="w-4 h-4 animate-spin" /> : null}
            {validating ? 'Validating...' : discovering ? 'Discovering clusters...' : 'Validate & Discover'}
          </button>
        </div>
      )}

      {/* Step 2: Select cluster */}
      {step === 'select' && (
        <div className="space-y-4">
          {clusters.length === 0 ? (
            <div className="text-center py-12 text-text-secondary">
              No clusters found in project {projectId}
            </div>
          ) : (
            clusters.map((cluster) => (
              <button
                key={cluster.id}
                onClick={() => handleSelect(cluster)}
                className="w-full text-left bg-surface border border-border rounded-lg p-4 hover:border-accent/30 transition-colors"
              >
                <div className="flex items-center justify-between">
                  <div>
                    <h3 className="font-medium">{cluster.name}</h3>
                    <div className="flex items-center gap-4 mt-1 text-xs text-text-secondary">
                      <span>{cluster.location}</span>
                      <span>v{cluster.k8s_version}</span>
                      <span>{cluster.node_count} nodes</span>
                    </div>
                  </div>
                  <StatusBadge status={cluster.status.toLowerCase()} />
                </div>
              </button>
            ))
          )}

          <button
            onClick={() => setStep('credentials')}
            className="text-sm text-text-secondary hover:text-text-primary transition-colors"
          >
            Back to credentials
          </button>
        </div>
      )}

      {/* Step 3: Confirm */}
      {step === 'confirm' && selected && (
        <div className="space-y-6">
          <div className="bg-surface border border-border rounded-lg p-5">
            <h3 className="font-medium mb-4">Connection Summary</h3>
            <div className="grid grid-cols-2 gap-4 text-sm">
              <div>
                <span className="text-text-secondary">Cluster</span>
                <p className="font-mono mt-1">{selected.name}</p>
              </div>
              <div>
                <span className="text-text-secondary">Location</span>
                <p className="mt-1">{selected.location}</p>
              </div>
              <div>
                <span className="text-text-secondary">K8s Version</span>
                <p className="font-mono mt-1">{selected.k8s_version}</p>
              </div>
              <div>
                <span className="text-text-secondary">Nodes</span>
                <p className="font-mono mt-1">{selected.node_count}</p>
              </div>
            </div>
          </div>

          <div>
            <label className="block text-sm font-medium mb-1.5">Display Name</label>
            <input
              type="text"
              value={displayName}
              onChange={(e) => setDisplayName(e.target.value)}
              className="w-full bg-surface border border-border rounded-md px-3 py-2 text-sm focus:outline-none focus:border-accent"
            />
          </div>

          <div className="flex gap-3">
            <button
              onClick={handleConnect}
              disabled={connecting}
              className="flex items-center gap-2 bg-accent text-black font-medium px-6 py-2 rounded-md hover:bg-accent/90 transition-colors disabled:opacity-50"
            >
              {connecting ? <Loader2 className="w-4 h-4 animate-spin" /> : null}
              {connecting ? 'Connecting...' : 'Connect Cluster'}
            </button>
            <button
              onClick={() => setStep('select')}
              className="px-6 py-2 rounded-md border border-border text-text-secondary hover:text-text-primary transition-colors"
            >
              Back
            </button>
          </div>
        </div>
      )}
    </div>
  )
}
