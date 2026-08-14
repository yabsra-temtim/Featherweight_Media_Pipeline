import { formatBytes, percentSmaller } from '../utils/formatBytes'
import { resolveDownloadUrl } from '../api/client'

/** One generated output: a compressed file in a specific format, with a preview and download link. */
export default function ResultCard({ output, originalSize }) {
  const url = resolveDownloadUrl(output.url)
  const reduction = percentSmaller(originalSize, output.size_bytes)

  return (
    <div className="group overflow-hidden rounded-2xl border border-ink-600 bg-ink-800/40 transition hover:border-ink-500">
      <div className="aspect-[4/3] w-full overflow-hidden bg-ink-900">
        <img
          src={url}
          alt={`${output.format} preview`}
          className="h-full w-full object-cover transition duration-300 group-hover:scale-[1.03]"
          loading="lazy"
        />
      </div>
      <div className="flex items-center justify-between gap-3 p-4">
        <div>
          <p className="font-mono text-xs uppercase tracking-wide text-quill-soft">{output.format}</p>
          <p className="mt-0.5 font-mono text-xs text-slate-400">
            {formatBytes(output.size_bytes)}
            {reduction > 0 && <span className="text-moss"> · -{reduction}%</span>}
          </p>
        </div>
        <a
          href={url}
          download
          className="shrink-0 rounded-full border border-ink-600 px-3.5 py-1.5 text-xs font-medium text-parchment-100 transition hover:border-quill hover:text-quill-soft"
        >
          Download
        </a>
      </div>
    </div>
  )
}
