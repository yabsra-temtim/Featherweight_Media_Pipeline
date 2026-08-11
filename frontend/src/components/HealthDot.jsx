import { useEffect, useState } from 'react'
import { getHealth } from '../api/client'

/** Small status readout showing whether the backend's WebP/AVIF encoders are available. */
export default function HealthDot() {
    const [health, setHealth] = useState(null)
    const [reachable, setReachable] = useState(true)

    useEffect(() => {
        const controller = new AbortController()
        getHealth(controller.signal)
            .then((data) => {
                setHealth(data)
                setReachable(true)
            })
            .catch(() => setReachable(false))
        return () => controller.abort()
    }, [])

    const dotColor = !reachable ? 'bg-clay' : health ? 'bg-moss' : 'bg-slate-500'
    const label = !reachable
        ? 'Backend unreachable'
        : health
            ? `WebP ${health.webp_enabled ? 'on' : 'off'} · AVIF ${health.avif_enabled ? 'on' : 'off'}`
            : 'Checking backend…'

    return (
        <div className="flex items-center gap-2 text-xs text-slate-400 font-mono">
            <span className={`h-1.5 w-1.5 rounded-full ${dotColor}`} />
            <span>{label}</span>
        </div>
    )
}
