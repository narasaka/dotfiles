import { useState } from 'react'
import { useNavigate } from 'react-router-dom'
import { Ship } from 'lucide-react'
import { api } from '../api/client'
import { useAuthStore } from '../stores/authStore'
import { toast } from 'sonner'

interface Props {
  needsSetup: boolean
  onSetupComplete: () => void
}

export default function Login({ needsSetup, onSetupComplete }: Props) {
  const navigate = useNavigate()
  const { setUser } = useAuthStore()
  const [isSetup, setIsSetup] = useState(needsSetup)
  const [email, setEmail] = useState('')
  const [password, setPassword] = useState('')
  const [name, setName] = useState('')
  const [loading, setLoading] = useState(false)

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    setLoading(true)

    try {
      if (isSetup) {
        const user = await api.post('auth/setup', {
          json: { email, password, name },
        }).json()
        setUser(user as any)
        onSetupComplete()
        toast.success('Admin account created')
      } else {
        const user = await api.post('auth/login', {
          json: { email, password },
        }).json()
        setUser(user as any)
        toast.success('Logged in')
      }
      navigate('/')
    } catch (err: any) {
      toast.error(isSetup ? 'Failed to create account' : 'Invalid credentials')
    } finally {
      setLoading(false)
    }
  }

  return (
    <div className="min-h-screen flex items-center justify-center bg-[#0A0A0A]">
      <div className="w-full max-w-sm">
        <div className="flex items-center justify-center gap-3 mb-8">
          <Ship className="w-10 h-10 text-accent" />
          <h1 className="text-2xl font-bold">Kubeploy</h1>
        </div>

        <div className="bg-surface border border-border rounded-lg p-6">
          <h2 className="text-lg font-semibold mb-1">
            {isSetup ? 'Create Admin Account' : 'Sign In'}
          </h2>
          <p className="text-sm text-text-secondary mb-6">
            {isSetup
              ? 'Set up your first admin account to get started.'
              : 'Sign in to your Kubeploy instance.'}
          </p>

          <form onSubmit={handleSubmit} className="space-y-4">
            {isSetup && (
              <div>
                <label className="block text-sm font-medium mb-1.5">Name</label>
                <input
                  type="text"
                  value={name}
                  onChange={(e) => setName(e.target.value)}
                  className="w-full bg-[#0A0A0A] border border-border rounded-md px-3 py-2 text-sm focus:outline-none focus:border-accent"
                  placeholder="Admin"
                />
              </div>
            )}
            <div>
              <label className="block text-sm font-medium mb-1.5">Email</label>
              <input
                type="email"
                value={email}
                onChange={(e) => setEmail(e.target.value)}
                className="w-full bg-[#0A0A0A] border border-border rounded-md px-3 py-2 text-sm focus:outline-none focus:border-accent"
                placeholder="admin@example.com"
                required
              />
            </div>
            <div>
              <label className="block text-sm font-medium mb-1.5">Password</label>
              <input
                type="password"
                value={password}
                onChange={(e) => setPassword(e.target.value)}
                className="w-full bg-[#0A0A0A] border border-border rounded-md px-3 py-2 text-sm focus:outline-none focus:border-accent"
                placeholder="••••••••"
                required
              />
            </div>
            <button
              type="submit"
              disabled={loading}
              className="w-full bg-accent text-black font-medium py-2 rounded-md hover:bg-accent/90 transition-colors disabled:opacity-50"
            >
              {loading ? 'Loading...' : isSetup ? 'Create Account' : 'Sign In'}
            </button>
          </form>
        </div>
      </div>
    </div>
  )
}
