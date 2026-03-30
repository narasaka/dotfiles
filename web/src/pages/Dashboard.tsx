import { Link } from 'react-router-dom'
import { Plus } from 'lucide-react'
import { useApps } from '../hooks/useApps'
import AppCard from '../components/AppCard'

export default function Dashboard() {
  const { apps } = useApps()

  return (
    <div>
      <div className="flex items-center justify-between mb-8">
        <div>
          <h1 className="text-2xl font-bold">Applications</h1>
          <p className="text-text-secondary text-sm mt-1">
            Manage your deployed applications
          </p>
        </div>
        <Link
          to="/apps/new"
          className="flex items-center gap-2 bg-accent text-black font-medium px-4 py-2 rounded-md hover:bg-accent/90 transition-colors"
        >
          <Plus className="w-4 h-4" />
          New App
        </Link>
      </div>

      {apps.length === 0 ? (
        <div className="text-center py-20">
          <p className="text-text-secondary text-lg mb-4">No applications yet</p>
          <Link
            to="/apps/new"
            className="inline-flex items-center gap-2 bg-accent text-black font-medium px-4 py-2 rounded-md hover:bg-accent/90 transition-colors"
          >
            <Plus className="w-4 h-4" />
            Create your first app
          </Link>
        </div>
      ) : (
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
          {apps.map((app) => (
            <AppCard key={app.id} app={app} />
          ))}
        </div>
      )}
    </div>
  )
}
