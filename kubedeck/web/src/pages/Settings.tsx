import { useState, useEffect } from 'react'
import { api } from '../api/client'
import { toast } from 'sonner'
import { Save } from 'lucide-react'

export default function Settings() {
  const [settings, setSettings] = useState<Record<string, string>>({})
  const [loading, setLoading] = useState(true)

  useEffect(() => {
    const fetch = async () => {
      try {
        const data = await api.get('settings').json<Record<string, string>>()
        setSettings(data)
      } catch {
        toast.error('Failed to load settings')
      } finally {
        setLoading(false)
      }
    }
    fetch()
  }, [])

  const handleSave = async () => {
    try {
      await api.put('settings', { json: settings })
      toast.success('Settings saved')
    } catch {
      toast.error('Failed to save settings')
    }
  }

  const update = (key: string, value: string) =>
    setSettings((prev) => ({ ...prev, [key]: value }))

  if (loading) {
    return <div className="text-text-secondary">Loading...</div>
  }

  return (
    <div className="max-w-2xl">
      <div className="flex items-center justify-between mb-8">
        <div>
          <h1 className="text-2xl font-bold">Settings</h1>
          <p className="text-text-secondary text-sm mt-1">
            Configure global Kubedeck settings
          </p>
        </div>
        <button
          onClick={handleSave}
          className="flex items-center gap-2 bg-accent text-black font-medium px-4 py-2 rounded-md hover:bg-accent/90 transition-colors"
        >
          <Save className="w-4 h-4" />
          Save
        </button>
      </div>

      <div className="space-y-8">
        <section>
          <h2 className="text-lg font-semibold mb-4">Container Registry</h2>
          <div className="space-y-4">
            <div>
              <label className="block text-sm font-medium mb-1.5">Registry URL</label>
              <input
                type="text"
                value={settings.registry_url || ''}
                onChange={(e) => update('registry_url', e.target.value)}
                className="w-full bg-surface border border-border rounded-md px-3 py-2 text-sm font-mono focus:outline-none focus:border-accent"
                placeholder="ghcr.io"
              />
            </div>
            <div className="grid grid-cols-2 gap-4">
              <div>
                <label className="block text-sm font-medium mb-1.5">Username</label>
                <input
                  type="text"
                  value={settings.registry_username || ''}
                  onChange={(e) => update('registry_username', e.target.value)}
                  className="w-full bg-surface border border-border rounded-md px-3 py-2 text-sm focus:outline-none focus:border-accent"
                />
              </div>
              <div>
                <label className="block text-sm font-medium mb-1.5">Password</label>
                <input
                  type="password"
                  value={settings.registry_password || ''}
                  onChange={(e) => update('registry_password', e.target.value)}
                  className="w-full bg-surface border border-border rounded-md px-3 py-2 text-sm focus:outline-none focus:border-accent"
                />
              </div>
            </div>
          </div>
        </section>

        <section>
          <h2 className="text-lg font-semibold mb-4">Defaults</h2>
          <div className="grid grid-cols-2 gap-4">
            <div>
              <label className="block text-sm font-medium mb-1.5">Default Namespace</label>
              <input
                type="text"
                value={settings.default_namespace || ''}
                onChange={(e) => update('default_namespace', e.target.value)}
                className="w-full bg-surface border border-border rounded-md px-3 py-2 text-sm font-mono focus:outline-none focus:border-accent"
              />
            </div>
            <div>
              <label className="block text-sm font-medium mb-1.5">Default Domain</label>
              <input
                type="text"
                value={settings.default_domain || ''}
                onChange={(e) => update('default_domain', e.target.value)}
                className="w-full bg-surface border border-border rounded-md px-3 py-2 text-sm font-mono focus:outline-none focus:border-accent"
                placeholder="example.com"
              />
            </div>
          </div>
        </section>

        <section>
          <h2 className="text-lg font-semibold mb-4">Build</h2>
          <div>
            <label className="block text-sm font-medium mb-1.5">Kaniko Image</label>
            <input
              type="text"
              value={settings.kaniko_image || ''}
              onChange={(e) => update('kaniko_image', e.target.value)}
              className="w-full bg-surface border border-border rounded-md px-3 py-2 text-sm font-mono focus:outline-none focus:border-accent"
            />
          </div>
        </section>
      </div>
    </div>
  )
}
