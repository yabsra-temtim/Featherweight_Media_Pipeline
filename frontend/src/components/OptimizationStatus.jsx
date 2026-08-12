import { formatBytes, percentSmaller } from '../utils/formatBytes'

/** Headline before/after summary — the number people actually came here for. */
export default function OptimizationStats({ originalBytes, outputs }) {
    if (!outputs?.length) return null

    const smallest = outputs.reduce((best, output) =>
        output.size_bytes < best.size_bytes ? output : best,
        outputs[0])

    const saved = percentSmaller(originalBytes, smallest.size_bytes)

    return (
        <div className="animate-fade-up rounded-2xl border border-quill/30 bg-quill/5 p-6 text-center">
            <div className="flex flex-col items-center gap-1">
                <span className="text-xs uppercase tracking-wide text-slate-400">Original</span>
                <span className="font-mono text-lg text-parchment-100">{formatBytes(originalBytes)}</span>
            </div>

            <div aria-hidden="true" className="my-2 text-quill">↓</div>

            <div className="flex flex-col items-center gap-1">
                <span className="text-xs uppercase tracking-wide text-slate-400">
                    Optimized <span className="normal-case text-slate-500">({smallest.format})</span>
                </span>
                <span className="font-mono text-lg text-parchment-100">{formatBytes(smallest.size_bytes)}</span>
            </div>

            <div aria-hidden="true" className="my-2 text-quill">↓</div>

            <div className="flex flex-col items-center gap-1">
                <span className="text-xs uppercase tracking-wide text-slate-400">Saved</span>
                <span className="font-display text-3xl italic text-moss">{saved}%</span>
            </div>
        </div>
    )
}
