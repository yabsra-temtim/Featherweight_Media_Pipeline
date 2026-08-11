import { formatBytes } from '../utils/formatBytes'

/** "Today's uploads" — a running log of this session's results. */
export default function HistoryList({ entries, onClear }) {
    if (entries.length === 0) return null

    return (
        <div className="rounded-2xl border border-ink-600 bg-ink-800/40 p-5">
            <div className="flex items-center justify-between">
                <p className="text-xs font-medium text-slate-400">Today&rsquo;s uploads</p>
                <button
                    type="button"
                    onClick={onClear}
                    className="text-[11px] text-slate-500 transition hover:text-clay"
                >
                    Clear
                </button>
            </div>

            <ul className="mt-3 max-h-56 space-y-1 overflow-y-auto scrollbar-thin pr-1">
                {entries.map((entry) => (
                    <li
                        key={entry.id}
                        className="flex items-center justify-between gap-3 rounded-lg px-2 py-1.5 text-xs transition hover:bg-ink-700/50"
                    >
                        <span className="flex min-w-0 items-center gap-2">
                            <span aria-hidden="true" className={entry.status === 'success' ? 'text-moss' : 'text-clay'}>
                                {entry.status === 'success' ? '✔' : '✖'}
                            </span>
                            <span className="truncate text-parchment-100">{entry.fileName}</span>
                        </span>
                        <span className="shrink-0 font-mono text-[11px] text-slate-400">
                            {entry.status === 'success' ? `-${entry.savedPercent}%` : 'failed'}
                            {entry.status === 'success' && (
                                <span className="ml-2 text-slate-500">{formatBytes(entry.optimizedBytes)}</span>
                            )}
                        </span>
                    </li>
                ))}
            </ul>
        </div>
    )
}
