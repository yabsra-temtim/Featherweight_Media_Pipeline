import { NavLink } from 'react-router-dom'
import FeatherMark from './FeatherMark'
import HealthDot from './HealthDot'
import ThemeToggle from './ThemeToggle'

const LINKS = [
    { to: '/', label: 'Home', end: true },
    { to: '/upload', label: 'Upload' },
    { to: '/about', label: 'About' },
]

export default function NavBar() {
    return (
        <header className="border-b border-ink-600">
            <div className="mx-auto flex max-w-3xl items-center justify-between gap-4 px-4 py-4">
                <NavLink to="/" className="flex items-center gap-2.5">
                    <span className="flex h-8 w-8 items-center justify-center rounded-lg border border-ink-600 bg-ink-800 text-quill">
                        <FeatherMark className="h-4 w-4" />
                    </span>
                    <span className="font-display text-lg italic leading-none text-parchment-100">
                        Featherweight
                    </span>
                </NavLink>

                <nav className="flex items-center gap-1 rounded-full border border-ink-600 bg-ink-800/60 p-1">
                    {LINKS.map(({ to, label, end }) => (
                        <NavLink
                            key={to}
                            to={to}
                            end={end}
                            className={({ isActive }) =>
                                `rounded-full px-3.5 py-1.5 text-xs font-medium transition ${isActive
                                    ? 'bg-quill text-ink-950'
                                    : 'text-slate-400 hover:text-parchment-100'
                                }`
                            }
                        >
                            {label}
                        </NavLink>
                    ))}
                </nav>

                <div className="flex items-center gap-3">
                    <HealthDot />
                    <ThemeToggle />
                </div>
            </div>
        </header>
    )
}
