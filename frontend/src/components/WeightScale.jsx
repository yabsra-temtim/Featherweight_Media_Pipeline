import { formatBytes, percentSmaller } from '../utils/formatBytes'

/**
 * Signature visual: a set of receding horizontal bars comparing the original
 * file weight against each generated output, like feather barbs falling
 * away from the quill.
 */
export default function WeightScale({ originalSize, outputs }) {
  if (!outputs?.length) return null

  const maxSize = Math.max(originalSize, ...outputs.map((o) => o.size_bytes))

  const rows = [
    { label: 'original', size_bytes: originalSize, isOriginal: true },
    ...outputs,
  ]

  return (
    <div className="rounded-2xl border border-ink-600 bg-ink-800/40 p-5">
      <p className="text-xs font-medium text-slate-400">Weight comparison</p>
      <div className="mt-4 space-y-2.5">
        {rows.map((row) => {
          const widthPct = Math.max(2, (row.size_bytes / maxSize) * 100)
          const reduction = row.isOriginal ? null : percentSmaller(originalSize, row.size_bytes)
          return (
            <div key={row.isOriginal ? 'original' : row.format} className="flex items-center gap-3">
              <span className="w-14 shrink-0 font-mono text-[11px] uppercase text-slate-400">
                {row.isOriginal ? 'orig' : row.format}
              </span>
              <div className="h-2.5 flex-1 overflow-hidden rounded-full bg-ink-700">
                <div
                  className={`h-full rounded-full ${row.isOriginal ? 'bg-slate-500' : 'bg-quill'}`}
                  style={{ width: `${widthPct}%` }}
                />
              </div>
              <span className="w-16 shrink-0 text-right font-mono text-[11px] text-slate-400">
                {formatBytes(row.size_bytes)}
              </span>
              <span className="w-12 shrink-0 text-right font-mono text-[11px] text-moss">
                {reduction !== null ? `-${reduction}%` : ''}
              </span>
            </div>
          )
        })}
      </div>
    </div>
  )
}
