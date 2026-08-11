import { useCallback, useRef, useState } from 'react'
import { formatBytes } from '../utils/formatBytes'
import { ACCEPTED_TYPES, validateFile } from '../utils/errorMessages'
import FeatherMark from './FeatherMark'

/** Drag-and-drop / click-to-browse target for choosing the source image. */
export default function Dropzone({ file, previewUrl, onSelect, onClear, disabled }) {
    const [isDragging, setIsDragging] = useState(false)
    const [localError, setLocalError] = useState(null)
    const inputRef = useRef(null)

    const handleFiles = useCallback(
        (fileList) => {
            const picked = fileList?.[0]
            if (!picked) return

            const validationError = validateFile(picked)
            if (validationError) {
                setLocalError(validationError)
                return
            }

            setLocalError(null)
            onSelect(picked)
        },
        [onSelect],
    )

    const onDrop = (event) => {
        event.preventDefault()
        setIsDragging(false)
        if (disabled) return
        handleFiles(event.dataTransfer.files)
    }

    if (file) {
        return (
            <div className="rounded-2xl border border-ink-600 bg-ink-800/60 p-4 flex items-center gap-4">
                <div className="h-20 w-20 shrink-0 overflow-hidden rounded-xl border border-ink-600 bg-ink-900">
                    {previewUrl && (
                        <img src={previewUrl} alt="" className="h-full w-full object-cover" />
                    )}
                </div>
                <div className="min-w-0 flex-1">
                    <p className="truncate text-sm font-medium text-parchment-100">{file.name}</p>
                    <p className="mt-0.5 font-mono text-xs text-slate-400">{formatBytes(file.size)}</p>
                </div>
                {!disabled && (
                    <button
                        type="button"
                        onClick={onClear}
                        className="shrink-0 rounded-full border border-ink-600 px-3 py-1.5 text-xs text-slate-400 transition hover:border-clay hover:text-clay"
                    >
                        Remove
                    </button>
                )}
            </div>
        )
    }

    return (
        <div>
            <div
                onDragOver={(event) => {
                    event.preventDefault()
                    setIsDragging(true)
                }}
                onDragLeave={() => setIsDragging(false)}
                onDrop={onDrop}
                onClick={() => inputRef.current?.click()}
                role="button"
                tabIndex={0}
                onKeyDown={(event) => {
                    if (event.key === 'Enter' || event.key === ' ') inputRef.current?.click()
                }}
                className={`group cursor-pointer rounded-2xl border-2 border-dashed p-10 text-center transition-colors ${isDragging
                        ? 'border-quill bg-quill/5'
                        : localError
                            ? 'border-clay/50'
                            : 'border-ink-600 hover:border-slate-500 bg-ink-800/40'
                    }`}
            >
                <input
                    ref={inputRef}
                    type="file"
                    accept={ACCEPTED_TYPES.join(',')}
                    className="hidden"
                    onChange={(event) => handleFiles(event.target.files)}
                />
                <FeatherMark
                    className={`mx-auto h-9 w-9 transition-all ${isDragging ? 'scale-110 text-quill' : 'text-slate-500 group-hover:text-quill'
                        }`}
                    animated={isDragging}
                />
                {isDragging ? (
                    <p className="mt-4 text-sm font-semibold text-quill">Drop your image here</p>
                ) : (
                    <>
                        <p className="mt-4 text-sm font-medium text-parchment-100">
                            Drop an image here, or click to browse
                        </p>
                        <p className="mt-1 text-xs text-slate-400">PNG, JPEG, or WebP</p>
                    </>
                )}
            </div>
            {localError && <p className="mt-2 text-xs text-clay">{localError}</p>}
        </div>
    )
}
