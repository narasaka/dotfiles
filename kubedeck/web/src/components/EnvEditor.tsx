import { useState, useEffect } from 'react'
import { Plus, Trash2 } from 'lucide-react'

interface Props {
  value: string
  onChange: (value: string) => void
}

interface EnvPair {
  key: string
  value: string
}

export default function EnvEditor({ value, onChange }: Props) {
  const [pairs, setPairs] = useState<EnvPair[]>([])

  useEffect(() => {
    try {
      const parsed = JSON.parse(value || '{}')
      const entries = Object.entries(parsed).map(([key, val]) => ({
        key,
        value: val as string,
      }))
      setPairs(entries.length > 0 ? entries : [{ key: '', value: '' }])
    } catch {
      setPairs([{ key: '', value: '' }])
    }
  }, [])

  const updateAndEmit = (newPairs: EnvPair[]) => {
    setPairs(newPairs)
    const obj: Record<string, string> = {}
    newPairs.forEach(({ key, value }) => {
      if (key.trim()) obj[key.trim()] = value
    })
    onChange(JSON.stringify(obj))
  }

  const addPair = () => updateAndEmit([...pairs, { key: '', value: '' }])

  const removePair = (index: number) => {
    const newPairs = pairs.filter((_, i) => i !== index)
    updateAndEmit(newPairs.length > 0 ? newPairs : [{ key: '', value: '' }])
  }

  const updatePair = (index: number, field: 'key' | 'value', val: string) => {
    const newPairs = pairs.map((p, i) =>
      i === index ? { ...p, [field]: val } : p
    )
    updateAndEmit(newPairs)
  }

  return (
    <div className="space-y-2">
      {pairs.map((pair, i) => (
        <div key={i} className="flex items-center gap-2">
          <input
            type="text"
            placeholder="KEY"
            value={pair.key}
            onChange={(e) => updatePair(i, 'key', e.target.value)}
            className="flex-1 bg-[#0A0A0A] border border-border rounded-md px-3 py-2 text-sm font-mono focus:outline-none focus:border-accent"
          />
          <input
            type="text"
            placeholder="value"
            value={pair.value}
            onChange={(e) => updatePair(i, 'value', e.target.value)}
            className="flex-1 bg-[#0A0A0A] border border-border rounded-md px-3 py-2 text-sm font-mono focus:outline-none focus:border-accent"
          />
          <button
            onClick={() => removePair(i)}
            className="p-2 text-text-secondary hover:text-red-400 transition-colors"
          >
            <Trash2 className="w-4 h-4" />
          </button>
        </div>
      ))}
      <button
        onClick={addPair}
        className="flex items-center gap-1 text-sm text-accent hover:text-accent/80 transition-colors"
      >
        <Plus className="w-4 h-4" />
        Add variable
      </button>
    </div>
  )
}
