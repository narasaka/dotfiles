import { useState, useEffect } from 'react'
import { useParams, useNavigate } from 'react-router-dom'
import { Play, RotateCcw, Trash2, GitBranch, Globe } from 'lucide-react'
import { api } from '../api/client'
import { toast } from 'sonner'
import { App, Build, Deployment } from '../stores/appStore'
import StatusBadge from '../components/StatusBadge'
import BuildList from '../components/BuildList'
import DeploymentStatus from '../components/DeploymentStatus'
import LogViewer from '../components/LogViewer'
import EnvEditor from '../components/EnvEditor'

type Tab = 'overview' | 'builds' | 'logs' | 'settings'

export default function AppDetail() {
  const { id } = useParams<{ id: string }>()
  const navigate = useNavigate()
  const [tab, setTab] = useState<Tab>('overview')
  const [app, setApp] = useState<App | null>(null)
  const [builds, setBuilds] = useState<Build[]>([])
  const [deployments, setDeployments] = useState<Deployment[]>([])
  const [editForm, setEditForm] = useState<any>(null)

  const fetchData = async () => {
    try {
      const [appData, buildsData, depsData] = await Promise.all([
        api.get(`apps/${id}`).json<App>(),
        api.get(`apps/${id}/builds`).json<Build[]>(),
        api.get(`apps/${id}/deployments`).json<Deployment[]>(),
      ])
      setApp(appData)
      setBuilds(buildsData)
      setDeployments(depsData)
    } catch {
      toast.error('Failed to load app')
      navigate('/')
    }
  }

  useEffect(() => {
    fetchData()
    const interval = setInterval(fetchData, 5000)
    return () => clearInterval(interval)
  }, [id])

  const triggerBuild = async () => {
    try {
      await api.post(`apps/${id}/builds`).json()
      toast.success('Build triggered')
      fetchData()
    } catch {
      toast.error('Failed to trigger build')
    }
  }

  const deleteApp = async () => {
    if (!confirm('Are you sure you want to delete this application?')) return
    try {
      await api.delete(`apps/${id}`)
      toast.success('App deleted')
      navigate('/')
    } catch {
      toast.error('Failed to delete app')
    }
  }

  const saveSettings = async () => {
    if (!editForm) return
    try {
      await api.put(`apps/${id}`, { json: editForm }).json()
      toast.success('Settings saved')
      fetchData()
    } catch {
      toast.error('Failed to save settings')
    }
  }

  const rollback = async (depId: string) => {
    try {
      await api.post(`deployments/${depId}/rollback`).json()
      toast.success('Rollback triggered')
      fetchData()
    } catch {
      toast.error('Failed to rollback')
    }
  }

  if (!app) {
    return <div className="text-text-secondary">Loading...</div>
  }

  const latestDep = deployments[0] || null

  return (
    <div>
      <div className="flex items-start justify-between mb-8">
        <div>
          <div className="flex items-center gap-3">
            <h1 className="text-2xl font-bold">{app.display_name || app.name}</h1>
            <StatusBadge status={app.status} />
          </div>
          <div className="flex items-center gap-4 mt-2 text-sm text-text-secondary">
            <span className="flex items-center gap-1 font-mono">
              <GitBranch className="w-3.5 h-3.5" />
              {app.git_branch}
            </span>
            {app.ingress_host && (
              <span className="flex items-center gap-1">
                <Globe className="w-3.5 h-3.5" />
                {app.ingress_host}
              </span>
            )}
          </div>
        </div>
        <div className="flex items-center gap-2">
          <button
            onClick={triggerBuild}
            className="flex items-center gap-2 bg-accent text-black font-medium px-4 py-2 rounded-md hover:bg-accent/90 transition-colors"
          >
            <Play className="w-4 h-4" />
            Build
          </button>
        </div>
      </div>

      {/* Tabs */}
      <div className="flex gap-1 border-b border-border mb-6">
        {(['overview', 'builds', 'logs', 'settings'] as Tab[]).map((t) => (
          <button
            key={t}
            onClick={() => {
              setTab(t)
              if (t === 'settings' && app) setEditForm({ ...app })
            }}
            className={`px-4 py-2.5 text-sm font-medium transition-colors border-b-2 -mb-px ${
              tab === t
                ? 'border-accent text-accent'
                : 'border-transparent text-text-secondary hover:text-text-primary'
            }`}
          >
            {t.charAt(0).toUpperCase() + t.slice(1)}
          </button>
        ))}
      </div>

      {/* Tab content */}
      {tab === 'overview' && (
        <div className="space-y-6">
          <DeploymentStatus deployment={latestDep} />

          {deployments.length > 1 && (
            <div>
              <h3 className="font-medium mb-3">Deployment History</h3>
              <div className="space-y-2">
                {deployments.slice(1).map((dep) => (
                  <div
                    key={dep.id}
                    className="flex items-center justify-between p-3 bg-surface border border-border rounded-lg"
                  >
                    <div className="flex items-center gap-3">
                      <StatusBadge status={dep.status} />
                      <span className="text-sm font-mono">{dep.k8s_deployment_name}</span>
                      <span className="text-xs text-text-secondary">
                        {dep.replicas_ready}/{dep.replicas_desired} replicas
                      </span>
                    </div>
                    <button
                      onClick={() => rollback(dep.id)}
                      className="flex items-center gap-1 text-xs text-text-secondary hover:text-accent transition-colors"
                    >
                      <RotateCcw className="w-3 h-3" />
                      Rollback
                    </button>
                  </div>
                ))}
              </div>
            </div>
          )}
        </div>
      )}

      {tab === 'builds' && (
        <BuildList builds={builds} appId={app.id} />
      )}

      {tab === 'logs' && (
        <LogViewer wsUrl={`/api/v1/apps/${app.id}/logs/ws`} />
      )}

      {tab === 'settings' && editForm && (
        <div className="max-w-2xl space-y-6">
          <div className="grid grid-cols-2 gap-4">
            <div>
              <label className="block text-sm font-medium mb-1.5">App Name</label>
              <input
                type="text"
                value={editForm.name}
                onChange={(e) => setEditForm({ ...editForm, name: e.target.value })}
                className="w-full bg-surface border border-border rounded-md px-3 py-2 text-sm font-mono focus:outline-none focus:border-accent"
              />
            </div>
            <div>
              <label className="block text-sm font-medium mb-1.5">Display Name</label>
              <input
                type="text"
                value={editForm.display_name}
                onChange={(e) => setEditForm({ ...editForm, display_name: e.target.value })}
                className="w-full bg-surface border border-border rounded-md px-3 py-2 text-sm focus:outline-none focus:border-accent"
              />
            </div>
          </div>

          <div>
            <label className="block text-sm font-medium mb-1.5">Git URL</label>
            <input
              type="text"
              value={editForm.git_url}
              onChange={(e) => setEditForm({ ...editForm, git_url: e.target.value })}
              className="w-full bg-surface border border-border rounded-md px-3 py-2 text-sm font-mono focus:outline-none focus:border-accent"
            />
          </div>

          <div className="grid grid-cols-3 gap-4">
            <div>
              <label className="block text-sm font-medium mb-1.5">Branch</label>
              <input
                type="text"
                value={editForm.git_branch}
                onChange={(e) => setEditForm({ ...editForm, git_branch: e.target.value })}
                className="w-full bg-surface border border-border rounded-md px-3 py-2 text-sm font-mono focus:outline-none focus:border-accent"
              />
            </div>
            <div>
              <label className="block text-sm font-medium mb-1.5">Dockerfile</label>
              <input
                type="text"
                value={editForm.dockerfile_path}
                onChange={(e) => setEditForm({ ...editForm, dockerfile_path: e.target.value })}
                className="w-full bg-surface border border-border rounded-md px-3 py-2 text-sm font-mono focus:outline-none focus:border-accent"
              />
            </div>
            <div>
              <label className="block text-sm font-medium mb-1.5">Port</label>
              <input
                type="number"
                value={editForm.port}
                onChange={(e) => setEditForm({ ...editForm, port: parseInt(e.target.value) || 8080 })}
                className="w-full bg-surface border border-border rounded-md px-3 py-2 text-sm font-mono focus:outline-none focus:border-accent"
              />
            </div>
          </div>

          <div className="grid grid-cols-2 gap-4">
            <div>
              <label className="block text-sm font-medium mb-1.5">Registry Image</label>
              <input
                type="text"
                value={editForm.registry_image}
                onChange={(e) => setEditForm({ ...editForm, registry_image: e.target.value })}
                className="w-full bg-surface border border-border rounded-md px-3 py-2 text-sm font-mono focus:outline-none focus:border-accent"
              />
            </div>
            <div>
              <label className="block text-sm font-medium mb-1.5">Namespace</label>
              <input
                type="text"
                value={editForm.namespace}
                onChange={(e) => setEditForm({ ...editForm, namespace: e.target.value })}
                className="w-full bg-surface border border-border rounded-md px-3 py-2 text-sm font-mono focus:outline-none focus:border-accent"
              />
            </div>
          </div>

          <div className="grid grid-cols-2 gap-4">
            <div>
              <label className="block text-sm font-medium mb-1.5">Ingress Host</label>
              <input
                type="text"
                value={editForm.ingress_host}
                onChange={(e) => setEditForm({ ...editForm, ingress_host: e.target.value })}
                className="w-full bg-surface border border-border rounded-md px-3 py-2 text-sm font-mono focus:outline-none focus:border-accent"
                placeholder="api.example.com"
              />
            </div>
            <div>
              <label className="block text-sm font-medium mb-1.5">Replicas</label>
              <input
                type="number"
                value={editForm.replicas}
                onChange={(e) => setEditForm({ ...editForm, replicas: parseInt(e.target.value) || 1 })}
                min={1}
                className="w-full bg-surface border border-border rounded-md px-3 py-2 text-sm font-mono focus:outline-none focus:border-accent"
              />
            </div>
          </div>

          <div className="flex items-center gap-6">
            <label className="flex items-center gap-2 text-sm">
              <input
                type="checkbox"
                checked={editForm.auto_deploy}
                onChange={(e) => setEditForm({ ...editForm, auto_deploy: e.target.checked })}
              />
              Auto Deploy on Push
            </label>
            <label className="flex items-center gap-2 text-sm">
              <input
                type="checkbox"
                checked={editForm.ingress_tls}
                onChange={(e) => setEditForm({ ...editForm, ingress_tls: e.target.checked })}
              />
              Enable TLS
            </label>
          </div>

          <div>
            <label className="block text-sm font-medium mb-1.5">Environment Variables</label>
            <EnvEditor
              value={editForm.env_vars}
              onChange={(v) => setEditForm({ ...editForm, env_vars: v })}
            />
          </div>

          <div className="flex items-center justify-between pt-4 border-t border-border">
            <button
              onClick={saveSettings}
              className="bg-accent text-black font-medium px-6 py-2 rounded-md hover:bg-accent/90 transition-colors"
            >
              Save Changes
            </button>
            <button
              onClick={deleteApp}
              className="flex items-center gap-2 text-red-400 hover:text-red-300 transition-colors text-sm"
            >
              <Trash2 className="w-4 h-4" />
              Delete Application
            </button>
          </div>
        </div>
      )}
    </div>
  )
}
