import WeightScale from './WeightScale'
import ResultCard from './ResultCard'

/** Full results view once a job has completed: weight comparison + downloadable outputs. */
export default function ResultsPanel({ job, onReset }) {
  if (job.status !== 'completed') return null

  return (
    <div className="animate-fade-up space-y-5">
      <div className="flex items-center justify-between">
        <h2 className="font-display text-lg italic text-parchment-100">Optimized &amp; ready</h2>
        <button
          type="button"
          onClick={onReset}
          className="rounded-full border border-ink-600 px-4 py-1.5 text-xs font-medium text-slate-400 transition hover:border-quill hover:text-quill-soft"
        >
          Compress another
        </button>
      </div>

      <WeightScale originalSize={job.original_size_bytes} outputs={job.outputs} />

      <div className="grid gap-4 sm:grid-cols-2">
        {job.outputs.map((output) => (
          <ResultCard key={output.format} output={output} originalSize={job.original_size_bytes} />
        ))}
      </div>
    </div>
  )
}
