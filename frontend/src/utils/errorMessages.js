export const ACCEPTED_TYPES = ['image/jpeg', 'image/png', 'image/webp']
export const MAX_UPLOAD_MB = Number(import.meta.env.VITE_MAX_UPLOAD_MB) || 10

/**
 * Validate a file client-side before it's sent, so people get an immediate,
 * specific reason instead of waiting on a round trip to find out.
 * Returns an error string, or null if the file looks fine.
 */
export function validateFile(file) {
  if (!file) return 'Please choose an image to upload.'

  if (!ACCEPTED_TYPES.includes(file.type)) {
    return 'Image must be PNG, JPEG, or WebP.'
  }

  const maxBytes = MAX_UPLOAD_MB * 1024 * 1024
  if (file.size > maxBytes) {
    return `Maximum file size is ${MAX_UPLOAD_MB} MB.`
  }

  return null
}

/**
 * Turn a raw error (from the network layer or the backend's JSON error
 * field) into something a person can act on.
 */
export function toFriendlyError(error) {
  if (!error) return 'Something went wrong. Please try again.'

  // Network-level failures (backend down, CORS blocked, offline, etc.)
  if (error.name === 'TypeError' || /failed to fetch|networkerror/i.test(error.message || '')) {
    return 'Server is currently unavailable. Please try again in a moment.'
  }

  const message = typeof error === 'string' ? error : error.message || ''

  const knownPatterns = [
    [/too large/i, `That image is too large. Maximum file size is ${MAX_UPLOAD_MB} MB.`],
    [/only image files|unsupported output format|not a valid supported image/i, 'Image must be PNG, JPEG, or WebP.'],
    [/width/i, 'The selected width is outside the allowed range.'],
    [/quality/i, 'Quality must be a number between 1 and 100.'],
    [/at least one output format/i, 'Choose at least one output format.'],
    [/server is busy/i, 'Server is busy right now. Please try again shortly.'],
    [/job not found/i, "We couldn't find that job — it may have expired."],
  ]

  for (const [pattern, friendly] of knownPatterns) {
    if (pattern.test(message)) return friendly
  }

  return message || 'Upload failed. Please try again.'
}
