import { useCallback, useEffect, useState } from 'react'

const STORAGE_KEY = 'featherweight:history'
const MAX_ENTRIES = 20

function isToday(isoString) {
  const date = new Date(isoString)
  const now = new Date()
  return (
    date.getFullYear() === now.getFullYear() &&
    date.getMonth() === now.getMonth() &&
    date.getDate() === now.getDate()
  )
}

function load() {
  try {
    const raw = window.localStorage.getItem(STORAGE_KEY)
    const parsed = raw ? JSON.parse(raw) : []
    return Array.isArray(parsed) ? parsed : []
  } catch {
    return []
  }
}

/** Tracks a running list of this session's uploads (name, outcome, savings) in localStorage. */
export function useUploadHistory() {
  const [entries, setEntries] = useState(load)

  useEffect(() => {
    window.localStorage.setItem(STORAGE_KEY, JSON.stringify(entries))
  }, [entries])

  const addEntry = useCallback((entry) => {
    setEntries((current) => [
      { id: crypto.randomUUID(), createdAt: new Date().toISOString(), ...entry },
      ...current,
    ].slice(0, MAX_ENTRIES))
  }, [])

  const clearHistory = useCallback(() => setEntries([]), [])

  const todaysEntries = entries.filter((entry) => isToday(entry.createdAt))

  return { entries: todaysEntries, addEntry, clearHistory }
}
