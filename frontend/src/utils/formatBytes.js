const UNITS = ['B', 'KB', 'MB', 'GB']

/** Format a byte count as a short human-readable string, e.g. 482331 -> "471 KB". */
export function formatBytes(bytes) {
  if (bytes === null || bytes === undefined || Number.isNaN(bytes)) return '—'
  if (bytes === 0) return '0 B'

  let value = bytes
  let unitIndex = 0
  while (value >= 1024 && unitIndex < UNITS.length - 1) {
    value /= 1024
    unitIndex += 1
  }

  const decimals = unitIndex === 0 ? 0 : value < 10 ? 1 : 0
  return `${value.toFixed(decimals)} ${UNITS[unitIndex]}`
}

/** Percentage reduction from original to compressed size, clamped to [0, 100]. */
export function percentSmaller(originalBytes, compressedBytes) {
  if (!originalBytes) return 0
  const reduction = ((originalBytes - compressedBytes) / originalBytes) * 100
  return Math.max(0, Math.min(100, Math.round(reduction)))
}
