export const API_BASE = (import.meta.env.VITE_API_URL || 'http://localhost:8080').replace(/\/$/, '')

async function parseErrorMessage(response, fallback) {
    try {
        const data = await response.json()
        return data.error || fallback
    } catch {
        return fallback
    }
}
/**
 * Upload an image for background optimization.
 * Uses XMLHttpRequest (rather than fetch) so we can report real upload
 * progress via onProgress(percent) while the bytes are in flight.
 * @param {{file: File, width: number, quality: number, formats: string[], onProgress?: (percent:number)=>void, signal?: AbortSignal}} input
 * @returns {Promise<{job_id: string, status: string}>}
 */
export function uploadImage({ file, width, quality, formats, onProgress, signal }) {
    const form = new FormData()
    form.append('file', file)
    form.append('width', String(width))
    form.append('quality', String(quality))
    formats.forEach((format) => form.append('formats', format))

    return new Promise((resolve, reject) => {
        const xhr = new XMLHttpRequest()
        xhr.open('POST', `${API_BASE}/api/v1/media/upload`)

        xhr.upload.onprogress = (event) => {
            if (!event.lengthComputable || !onProgress) return
            onProgress(Math.round((event.loaded / event.total) * 100))
        }

        xhr.onload = () => {
            let body = null
            try {
                body = JSON.parse(xhr.responseText)
            } catch {
                // ignore malformed body, handled below via status check
            }

            if (xhr.status >= 200 && xhr.status < 300) {
                onProgress?.(100)
                resolve(body)
            } else {
                reject(new Error(body?.error || 'Upload failed. Please try again.'))
            }
        }

        xhr.onerror = () => reject(new TypeError('Failed to fetch'))
        xhr.onabort = () => reject(new DOMException('Upload aborted', 'AbortError'))

        if (signal) {
            if (signal.aborted) {
                xhr.abort()
                return
            }
            signal.addEventListener('abort', () => xhr.abort())
        }

        xhr.send(form)
    })
}

/**
 * Fetch the current state of a processing job.
 * @param {string} jobId
 */
export async function getJob(jobId, signal) {
    const response = await fetch(`${API_BASE}/api/v1/jobs/${jobId}`, { signal })

    if (!response.ok) {
        throw new Error(await parseErrorMessage(response, 'Could not check job status.'))
    }

    return response.json()
}

/** Check which encoders the backend has available (WebP / AVIF). */
export async function getHealth(signal) {
    const response = await fetch(`${API_BASE}/api/v1/health`, { signal })

    if (!response.ok) {
        throw new Error('Health check failed.')
    }

    return response.json()
}

/** Resolve a backend-relative download path (e.g. "/downloads/xyz/optimized.webp") to an absolute URL. */
export function resolveDownloadUrl(path) {
    if (!path) return ''
    return path.startsWith('http') ? path : `${API_BASE}${path}`
}
