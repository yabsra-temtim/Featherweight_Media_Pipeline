const FORMATS = [
  { label: 'JPEG', note: 'Best all-round compatibility' },
  { label: 'PNG', note: 'Lossless, supports transparency' },
  { label: 'WebP', note: 'Smaller than JPEG at similar quality' },
  { label: 'AVIF', note: 'The smallest files, newest browser support' },
]

export default function About() {
  return (
    <div className="mx-auto flex max-w-xl flex-col gap-8 px-4 py-16">
      <div>
        <h1 className="font-display text-2xl italic text-parchment-100">About Featherweight</h1>
        <p className="mt-3 text-sm leading-relaxed text-slate-400">
          Featherweight takes a single uploaded image and, in the background,
          resizes it to a target width, compresses it, and re-encodes it into
          several modern formats — so you get web-ready assets without
          running your own build pipeline.
        </p>
      </div>

      <div className="rounded-2xl border border-ink-600 bg-ink-800/40 p-5">
        <p className="text-xs font-medium text-slate-400">Output formats</p>
        <ul className="mt-3 space-y-2">
          {FORMATS.map((format) => (
            <li key={format.label} className="flex items-baseline justify-between gap-3 text-sm">
              <span className="font-mono text-quill-soft">{format.label}</span>
              <span className="text-right text-xs text-slate-400">{format.note}</span>
            </li>
          ))}
        </ul>
      </div>

      <div className="rounded-2xl border border-ink-600 bg-ink-800/40 p-5">
        <p className="text-xs font-medium text-slate-400">How it works</p>
        <ol className="mt-3 space-y-2 text-sm text-slate-400">
          <li>1. Your image uploads directly to the Featherweight backend.</li>
          <li>2. A background worker resizes, compresses, and converts it.</li>
          <li>3. You poll for status and download each result once it's ready.</li>
          <li>4. Files are automatically removed after their retention window.</li>
        </ol>
      </div>

      <p className="text-xs text-slate-500">
        Built with a Go backend and a React + Tailwind frontend.
      </p>
    </div>
  )
}
