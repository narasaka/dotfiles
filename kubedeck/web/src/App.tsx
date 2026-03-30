import { Routes, Route, Navigate } from 'react-router-dom'
import { useAuthStore } from './stores/authStore'
import { useEffect, useState } from 'react'
import Layout from './components/Layout'
import Dashboard from './pages/Dashboard'
import AppDetail from './pages/AppDetail'
import AppCreate from './pages/AppCreate'
import BuildDetail from './pages/BuildDetail'
import Settings from './pages/Settings'
import Login from './pages/Login'
import { api } from './api/client'

function ProtectedRoute({ children }: { children: React.ReactNode }) {
  const { user, isLoading } = useAuthStore()
  if (isLoading) return <div className="flex items-center justify-center h-screen">Loading...</div>
  if (!user) return <Navigate to="/login" />
  return <>{children}</>
}

export default function App() {
  const { setUser, setLoading } = useAuthStore()
  const [needsSetup, setNeedsSetup] = useState<boolean | null>(null)

  useEffect(() => {
    const checkAuth = async () => {
      try {
        const res = await api.get('auth/check').json<{ needs_setup: boolean }>()
        setNeedsSetup(res.needs_setup)

        if (!res.needs_setup) {
          try {
            const user = await api.get('auth/me').json()
            setUser(user as any)
          } catch {
            setUser(null)
          }
        }
      } catch {
        setUser(null)
      } finally {
        setLoading(false)
      }
    }
    checkAuth()
  }, [setUser, setLoading])

  if (needsSetup === null) {
    return <div className="flex items-center justify-center h-screen text-text-secondary">Loading...</div>
  }

  return (
    <Routes>
      <Route path="/login" element={<Login needsSetup={needsSetup} onSetupComplete={() => setNeedsSetup(false)} />} />
      <Route path="/" element={<ProtectedRoute><Layout /></ProtectedRoute>}>
        <Route index element={<Dashboard />} />
        <Route path="apps/new" element={<AppCreate />} />
        <Route path="apps/:id" element={<AppDetail />} />
        <Route path="apps/:appId/builds/:buildId" element={<BuildDetail />} />
        <Route path="settings" element={<Settings />} />
      </Route>
    </Routes>
  )
}
