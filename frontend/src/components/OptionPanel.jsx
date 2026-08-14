const FORMAT_OPTIONS = [
  { value: 'jpeg', label: 'JPEG' },
  { value: 'png', label: 'PNG' },
  { value: 'webp', label: 'WebP' },
  { value: 'avif', label: 'AVIF' },
]

/** Controls for target width, output quality, and which formats to generate. */
export default function OptionsPanel({ width, quality, formats, onWidthChange, onQualityChange, onToggleFormat, disabled }) {
  const toggleFormat = (value) => {
    if (disabled) return
    onToggleFormat(value)
  }

  return (
    <div className="rounded-2xl border border-ink-600 bg-ink-800/40 p-5">
      <div className="grid gap-6 sm:grid-cols-2">
        <label className="block">
          <span className="flex items-baseline justify-between text-xs font-medium text-slate-400">
            <span>Max width</span>
            <span className="font-mono text-quill-soft">{width}px</span>
          </span>
          <input
            type="range"
            min="200"
            max="4000"
            step="50"
            value={width}
            disabled={disabled}
            onChange={(event) => onWidthChange(Number(event.target.value))}
            className="mt-3 w-full accent-quill disabled:opacity-40"
          />
        </label>

        <label className="block">
          <span className="flex items-baseline justify-between text-xs font-medium text-slate-400">
            <span>Quality</span>
            <span className="font-mono text-quill-soft">{quality}</span>
          </span>
          <input
            type="range"
            min="1"
            max="100"
            value={quality}
            disabled={disabled}
            onChange={(event) => onQualityChange(Number(event.target.value))}
            className="mt-3 w-full accent-quill disabled:opacity-40"
          />
        </label>
      </div>

      <div className="mt-6">
        <span className="text-xs font-medium text-slate-400">Output formats</span>
        <div className="mt-3 flex flex-wrap gap-2">
          {FORMAT_OPTIONS.map(({ value, label }) => {
            const active = formats.includes(value)
            return (
              <button
                key={value}
                type="button"
                disabled={disabled}
                onClick={() => toggleFormat(value)}
                className={`rounded-full border px-3.5 py-1.5 text-xs font-medium font-mono tracking-wide transition disabled:cursor-not-allowed disabled:opacity-40 ${
                  active
                    ? 'border-quill bg-quill/15 text-quill-soft'
                    : 'border-ink-600 text-slate-400 hover:border-slate-500 hover:text-parchment-100'
                }`}
              >
                {label}
              </button>
            )
          })}
        </div>
      </div>
    </div>
  )
}
