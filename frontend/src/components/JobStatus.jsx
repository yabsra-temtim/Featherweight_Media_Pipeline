import FeatherMark from './FeatherMark'

const STATUS_COPY = {
    pending: 'Queued for processing…',
    processing: 'Compressing, resizing, and converting…',
    completed: 'Done — your images are ready.',
    failed: 'Something went wrong.',
}

/** Progress readout for an in-flight (or finished) optimization job. */
export default function JobStatus({ job, pollingError }) {
    if (!job) return null

    const isActive = job.status === 'pending' || job.status === 'processing'
    const isFailed = job.status === 'failed'

    return (
        <div className="animate-fade-up rounded-2xl border border-ink-600 bg-ink-800/40 p-5">
            <div className="flex items-center justify-between gap-3">
                <div className="flex items-center gap-2.5">
                    <FeatherMark
                        className={`h-4 w-4 ${isFailed ? 'text-clay' : 'text-quill'}`}
                        animated={isActive}
                    />
                    <span className={`text-sm font-medium ${isFailed ? 'text-clay' : 'text-parchment-100'}`}>
                        {STATUS_COPY[job.status] || job.status}
                    </span>
                </div>
                <span className="font-mono text-xs text-slate-400">{job.progress}%</span>
            </div>

            <div className="mt-3 h-1.5 w-full overflow-hidden rounded-full bg-ink-700">
                <div
                    className={`h-full rounded-full transition-all duration-500 ease-out ${isFailed ? 'bg-clay' : 'shimmer-bg animate-shimmer'
                        }`}
                    style={{ width: `${Math.max(4, job.progress)}%` }}
                />
            </div>

            {job.error && (
                <p className="mt-3 text-xs text-clay">{job.error}</p>
            )}

            {pollingError && !isFailed && (
                <p className="mt-3 text-xs text-slate-400">
                    Having trouble reaching the server — retrying…
                </p>
            )}

            {job.notes?.length > 0 && (
                <ul className="mt-3 space-y-1">
                    {job.notes.map((note) => (
                        <li key={note} className="text-xs text-quill-soft/80">
                            {note}
                        </li>
                    ))}
                </ul>
            )}
        </div>
    )
}
