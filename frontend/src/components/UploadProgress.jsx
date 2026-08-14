import FeatherMark from './FeatherMark'

const BAR_LENGTH = 24

function blockBar(percent) {
  const filled = Math.round((percent / 100) * BAR_LENGTH)
  return '█'.repeat(filled) + '░'.repeat(Math.max(0, BAR_LENGTH - filled))
}

/** Loading indicator shown while the file is actively being sent to the server. */
export default function UploadProgress({ percent }) {
  return (
    <div className="animate-fade-up rounded-2xl border border-ink-600 bg-ink-800/40 p-5">
      <div className="flex items-center gap-2.5">
        <FeatherMark className="h-4 w-4 text-quill" animated />
        <span className="text-sm font-medium text-parchment-100">Uploading…</span>
        <span className="ml-auto font-mono text-xs text-slate-400">{percent}%</span>
      </div>

      <div className="mt-3 h-1.5 w-full overflow-hidden rounded-full bg-ink-700">
        <div
          className="h-full rounded-full bg-quill transition-all duration-200 ease-out"
          style={{ width: `${Math.max(4, percent)}%` }}
        />
      </div>

      <p aria-hidden="true" className="mt-2 select-none font-mono text-[11px] tracking-tight text-quill-soft/70">
        {blockBar(percent)}
      </p>
    </div>
  )
}
