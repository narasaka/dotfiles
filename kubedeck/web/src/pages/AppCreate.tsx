import { useState } from 'react'
import { useNavigate } from 'react-router-dom'
import { ChevronDown, ChevronUp } from 'lucide-react'
import { api } from '../api/client'
import { toast } from 'sonner'
import EnvEditor from '../components/EnvEditor'

export default function AppCreate() {
  const navigate = useNavigate()
  const [loading, setLoading] = useState(false)
  const [showAdvanced, setShowAdvanced] = useState(false)

  const [form, setForm] = useState({
    name: '',
    display_name: '',
    git_url: '',
    git_branch: 'main',
    dockerfile_path: 'Dockerfile',
    registry_image: '',
    port: 8080,
    namespace: 'default',
    replicas: 1,
    env_vars: '{}',
    auto_deploy: true,
    ingress_host: '',
    ingress_tls: false,
  })

  const update = (field: string, value: any) =>
    setForm((prev) => ({ ...prev, [field]: value }))

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    setLoading(true)

    try {
      const app = await api.post('apps', { json: form }).json<{ id: string }>()
      toast.success('App created')
      navigate(`/apps/${app.id}`)
    } catch {
      toast.error('Failed to create app')
    } finally {
      setLoading(false)
    }
  }

  return (
    <div className="max-w-2xl">
      <h1 className="text-2xl font-bold mb-2">Create Application</h1>
      <p className="text-text-secondary text-sm mb-8">
        Configure a new application to build and deploy.
      </p>

      <form onSubmit={handleSubmit} className="space-y-6">
        <div className="grid grid-cols-2 gap-4">
          <div>
            <label className="block text-sm font-medium mb-1.5">App Name</label>
            <input
              type="text"
              value={form.name}
              onChange={(e) => update('name', e.target.value)}
              className="w-full bg-surface border border-border rounded-md px-3 py-2 text-sm font-mono focus:outline-none focus:border-accent"
              placeholder="my-api"
              required
            />
          </div>
          <div>
            <label className="block text-sm font-medium mb-1.5">Display Name</label>
            <input
              type="text"
              value={form.display_name}
              onChange={(e) => update('display_name', e.target.value)}
              className="w-full bg-surface border border-border rounded-md px-3 py-2 text-sm focus:outline-none focus:border-accent"
              placeholder="My API"
            />
          </div>
        </div>

        <div>
          <label className="block text-sm font-medium mb-1.5">Git Repository URL</label>
          <input
            type="text"
            value={form.git_url}
            onChange={(e) => update('git_url', e.target.value)}
            className="w-full bg-surface border border-border rounded-md px-3 py-2 text-sm font-mono focus:outline-none focus:border-accent"
            placeholder="https://github.com/org/repo.git"
            required
          />
        </div>

        <div className="grid grid-cols-2 gap-4">
          <div>
            <label className="block text-sm font-medium mb-1.5">Branch</label>
            <input
              type="text"
              value={form.git_branch}
              onChange={(e) => update('git_branch', e.target.value)}
              className="w-full bg-surface border border-border rounded-md px-3 py-2 text-sm font-mono focus:outline-none focus:border-accent"
            />
          </div>
          <div>
            <label className="block text-sm font-medium mb-1.5">Dockerfile Path</label>
            <input
              type="text"
              value={form.dockerfile_path}
              onChange={(e) => update('dockerfile_path', e.target.value)}
              className="w-full bg-surface border border-border rounded-md px-3 py-2 text-sm font-mono focus:outline-none focus:border-accent"
            />
          </div>
        </div>

        <div className="grid grid-cols-2 gap-4">
          <div>
            <label className="block text-sm font-medium mb-1.5">Registry Image</label>
            <input
              type="text"
              value={form.registry_image}
              onChange={(e) => update('registry_image', e.target.value)}
              className="w-full bg-surface border border-border rounded-md px-3 py-2 text-sm font-mono focus:outline-none focus:border-accent"
              placeholder="ghcr.io/org/my-api"
            />
          </div>
          <div>
            <label className="block text-sm font-medium mb-1.5">Container Port</label>
            <input
              type="number"
              value={form.port}
              onChange={(e) => update('port', parseInt(e.target.value) || 8080)}
              className="w-full bg-surface border border-border rounded-md px-3 py-2 text-sm font-mono focus:outline-none focus:border-accent"
            />
          </div>
        </div>

        <button
          type="button"
          onClick={() => setShowAdvanced(!showAdvanced)}
          className="flex items-center gap-2 text-sm text-text-secondary hover:text-text-primary transition-colors"
        >
          {showAdvanced ? <ChevronUp className="w-4 h-4" /> : <ChevronDown className="w-4 h-4" />}
          Advanced Options
        </button>

        {showAdvanced && (
          <div className="space-y-6 border border-border rounded-lg p-4">
            <div className="grid grid-cols-2 gap-4">
              <div>
                <label className="block text-sm font-medium mb-1.5">Namespace</label>
                <input
                  type="text"
                  value={form.namespace}
                  onChange={(e) => update('namespace', e.target.value)}
                  className="w-full bg-[#0A0A0A] border border-border rounded-md px-3 py-2 text-sm font-mono focus:outline-none focus:border-accent"
                />
              </div>
              <div>
                <label className="block text-sm font-medium mb-1.5">Replicas</label>
                <input
                  type="number"
                  value={form.replicas}
                  onChange={(e) => update('replicas', parseInt(e.target.value) || 1)}
                  min={1}
                  className="w-full bg-[#0A0A0A] border border-border rounded-md px-3 py-2 text-sm font-mono focus:outline-none focus:border-accent"
                />
              </div>
            </div>

            <div className="grid grid-cols-2 gap-4">
              <div>
                <label className="block text-sm font-medium mb-1.5">Ingress Host</label>
                <input
                  type="text"
                  value={form.ingress_host}
                  onChange={(e) => update('ingress_host', e.target.value)}
                  className="w-full bg-[#0A0A0A] border border-border rounded-md px-3 py-2 text-sm font-mono focus:outline-none focus:border-accent"
                  placeholder="api.example.com"
                />
              </div>
              <div className="flex items-end gap-6">
                <label className="flex items-center gap-2 text-sm">
                  <input
                    type="checkbox"
                    checked={form.ingress_tls}
                    onChange={(e) => update('ingress_tls', e.target.checked)}
                    className="rounded border-border"
                  />
                  TLS
                </label>
                <label className="flex items-center gap-2 text-sm">
                  <input
                    type="checkbox"
                    checked={form.auto_deploy}
                    onChange={(e) => update('auto_deploy', e.target.checked)}
                    className="rounded border-border"
                  />
                  Auto Deploy
                </label>
              </div>
            </div>

            <div>
              <label className="block text-sm font-medium mb-1.5">Environment Variables</label>
              <EnvEditor
                value={form.env_vars}
                onChange={(v) => update('env_vars', v)}
              />
            </div>
          </div>
        )}

        <div className="flex gap-3">
          <button
            type="submit"
            disabled={loading}
            className="bg-accent text-black font-medium px-6 py-2 rounded-md hover:bg-accent/90 transition-colors disabled:opacity-50"
          >
            {loading ? 'Creating...' : 'Create Application'}
          </button>
          <button
            type="button"
            onClick={() => navigate('/')}
            className="px-6 py-2 rounded-md border border-border text-text-secondary hover:text-text-primary transition-colors"
          >
            Cancel
          </button>
        </div>
      </form>
    </div>
  )
}
