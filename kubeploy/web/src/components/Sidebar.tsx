import { NavLink } from 'react-router-dom'
import { LayoutDashboard, Settings, Ship } from 'lucide-react'
import { cn } from '../lib/utils'

const navItems = [
  { to: '/', icon: LayoutDashboard, label: 'Dashboard' },
  { to: '/settings', icon: Settings, label: 'Settings' },
]

export default function Sidebar() {
  return (
    <aside className="fixed left-0 top-0 h-screen w-60 bg-surface border-r border-border flex flex-col">
      <div className="flex items-center gap-3 px-6 py-5 border-b border-border">
        <Ship className="w-7 h-7 text-accent" />
        <span className="text-lg font-semibold tracking-tight">Kubeploy</span>
      </div>

      <nav className="flex-1 px-3 py-4 space-y-1">
        {navItems.map(({ to, icon: Icon, label }) => (
          <NavLink
            key={to}
            to={to}
            end={to === '/'}
            className={({ isActive }) =>
              cn(
                'flex items-center gap-3 px-3 py-2 rounded-md text-sm transition-colors',
                isActive
                  ? 'bg-accent/10 text-accent'
                  : 'text-text-secondary hover:text-text-primary hover:bg-white/5'
              )
            }
          >
            <Icon className="w-4 h-4" />
            {label}
          </NavLink>
        ))}
      </nav>

      <div className="px-4 py-3 border-t border-border text-xs text-text-secondary">
        Kubeploy v0.1.0
      </div>
    </aside>
  )
}
