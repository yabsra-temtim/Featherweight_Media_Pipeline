import { useEffect, useRef, useState } from 'react'
import { getJob } from '../api/client'

const TERMINAL_STATUSES = new Set(['completed', 'failed'])
const POLL_INTERVAL_MS = 900

/**
 * Polls GET /api/v1/jobs/:id while a job is pending/processing, stopping
 * automatically once the backend reports "completed" or "failed".
 */
export function useJobPolling(jobId) {
  const [job, setJob] = useState(null)
  const [error, setError] = useState(null)
  const timeoutId = useRef(null)

  useEffect(() => {
    setJob(null)
    setError(null)

    if (!jobId) return undefined

    const controller = new AbortController()
    let cancelled = false

    async function tick() {
      try {
        const data = await getJob(jobId, controller.signal)
        if (cancelled) return
        setJob(data)
        setError(null)
        if (!TERMINAL_STATUSES.has(data.status)) {
          timeoutId.current = setTimeout(tick, POLL_INTERVAL_MS)
        }
      } catch (err) {
        if (cancelled || err.name === 'AbortError') return
        setError(err.message)
        timeoutId.current = setTimeout(tick, POLL_INTERVAL_MS)
      }
    }

    tick()

    return () => {
      cancelled = true
      controller.abort()
      clearTimeout(timeoutId.current)
    }
  }, [jobId])

  return { job, error }
}
