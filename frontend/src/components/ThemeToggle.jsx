import { useTheme } from '../context/ThemeContext'

/** Sun/moon switch for light and dark mode. */
export default function ThemeToggle() {
  const { theme, toggleTheme } = useTheme()
  const isDark = theme === 'dark'

  return (
    <button
      type="button"
      onClick={toggleTheme}
      aria-label={isDark ? 'Switch to light mode' : 'Switch to dark mode'}
      title={isDark ? 'Switch to light mode' : 'Switch to dark mode'}
      className="flex h-8 w-8 shrink-0 items-center justify-center rounded-full border border-ink-600 text-sm transition hover:border-quill"
    >
      <span aria-hidden="true">{isDark ? '🌙' : '☀️'}</span>
    </button>
  )
}
