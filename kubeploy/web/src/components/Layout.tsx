import { Outlet } from 'react-router-dom'
import Sidebar from './Sidebar'
import { useAuthStore } from '../stores/authStore'
import { api } from '../api/client'
import { LogOut, User } from 'lucide-react'

export default function Layout() {
  const { user, logout } = useAuthStore()

  const handleLogout = async () => {
    await api.post('auth/logout')
    logout()
    window.location.href = '/login'
  }

  return (
    <div className="flex min-h-screen">
      <Sidebar />
      <div className="flex-1 ml-60">
        <header className="sticky top-0 z-10 flex items-center justify-between px-8 py-4 bg-[#0A0A0A]/80 backdrop-blur-sm border-b border-border">
          <div />
          <div className="flex items-center gap-4">
            <div className="flex items-center gap-2 text-sm text-text-secondary">
              <User className="w-4 h-4" />
              {user?.name || user?.email}
            </div>
            <button
              onClick={handleLogout}
              className="flex items-center gap-1.5 text-sm text-text-secondary hover:text-text-primary transition-colors"
            >
              <LogOut className="w-4 h-4" />
              Logout
            </button>
          </div>
        </header>
        <main className="p-8">
          <Outlet />
        </main>
      </div>
    </div>
  )
}
