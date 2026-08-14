 import { Link } from 'react-router-dom'
import FeatherMark from '../components/FeatherMark'

const STEPS = [
  { title: 'Upload', body: 'Drop in a PNG, JPEG, or WebP — any size.' },
  { title: 'Optimize', body: 'We compress, resize, and convert it in the background.' },
  { title: 'Download', body: 'Grab web-ready JPEG, PNG, WebP, and AVIF versions.' },
]

export default function Home() {
  return (
    <div className="mx-auto flex max-w-xl flex-col gap-10 px-4 py-16 text-center sm:py-24">
      <div className="flex flex-col items-center gap-5">
        <span className="flex h-14 w-14 items-center justify-center rounded-2xl border border-ink-600 bg-ink-800 text-quill">
          <FeatherMark className="h-6 w-6" />
        </span>
        <h1 className="font-display text-4xl italic leading-tight text-parchment-100">
          Make your images fly light.
        </h1>
        <p className="max-w-sm text-sm text-slate-400">
          Featherweight compresses, resizes, and converts your images into
          modern web formats automatically — no config, no waiting around.
        </p>
        <Link
          to="/upload"
          className="mt-2 rounded-full bg-quill px-6 py-3 text-sm font-semibold text-ink-950 transition hover:bg-quill-soft"
        >
          Optimize an image
        </Link>
      </div>

      <div className="grid gap-3 sm:grid-cols-3">
        {STEPS.map((step, index) => (
          <div key={step.title} className="rounded-2xl border border-ink-600 bg-ink-800/40 p-5 text-left">
            <span className="font-mono text-xs text-quill">{String(index + 1).padStart(2, '0')}</span>
            <p className="mt-2 text-sm font-medium text-parchment-100">{step.title}</p>
            <p className="mt-1 text-xs text-slate-400">{step.body}</p>
          </div>
        ))}
      </div>
    </div>
  )
}
