import { cn } from '../lib/utils'

const statusStyles: Record<string, string> = {
  running: 'bg-emerald-500/10 text-emerald-400 border-emerald-500/20',
  building: 'bg-cyan-500/10 text-cyan-400 border-cyan-500/20 animate-pulse',
  deploying: 'bg-cyan-500/10 text-cyan-400 border-cyan-500/20 animate-pulse',
  rolling_out: 'bg-cyan-500/10 text-cyan-400 border-cyan-500/20 animate-pulse',
  success: 'bg-emerald-500/10 text-emerald-400 border-emerald-500/20',
  failed: 'bg-red-500/10 text-red-400 border-red-500/20',
  inactive: 'bg-zinc-500/10 text-zinc-400 border-zinc-500/20',
  pending: 'bg-amber-500/10 text-amber-400 border-amber-500/20',
  cancelled: 'bg-zinc-500/10 text-zinc-400 border-zinc-500/20',
  rolled_back: 'bg-amber-500/10 text-amber-400 border-amber-500/20',
}

export default function StatusBadge({ status }: { status: string }) {
  return (
    <span
      className={cn(
        'inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium border',
        statusStyles[status] || statusStyles.inactive
      )}
    >
      {status.replace('_', ' ')}
    </span>
  )
}
