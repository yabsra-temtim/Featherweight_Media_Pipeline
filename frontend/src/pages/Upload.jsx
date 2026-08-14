import { useEffect, useMemo, useRef, useState } from 'react'
import Dropzone from '../components/Dropzone'
import OptionsPanel from '../components/OptionsPanel'
import UploadProgress from '../components/UploadProgress'
import JobStatus from '../components/JobStatus'
import OptimizationStats from '../components/OptimizationStats'
import ResultsPanel from '../components/ResultsPanel'
import HistoryList from '../components/HistoryList'
import FeatherMark from '../components/FeatherMark'
import { uploadImage } from '../api/client'
import { useJobPolling } from '../hooks/useJobPolling'
import { useUploadHistory } from '../hooks/useUploadHistory'
import { toFriendlyError, validateFile } from '../utils/errorMessages'
import { percentSmaller } from '../utils/formatBytes'

const DEFAULT_FORMATS = ['jpeg', 'png', 'webp', 'avif']

function pickSmallestOutput(outputs) {
  return outputs.reduce((best, output) => (output.size_bytes < best.size_bytes ? output : best), outputs[0])
}

export default function Upload() {
  const [file, setFile] = useState(null)
  const [previewUrl, setPreviewUrl] = useState(null)
  const [width, setWidth] = useState(1200)
  const [quality, setQuality] = useState(80)
  const [formats, setFormats] = useState(DEFAULT_FORMATS)

  const [jobId, setJobId] = useState(null)
  const [isUploading, setIsUploading] = useState(false)
  const [uploadPercent, setUploadPercent] = useState(0)
  const [submitError, setSubmitError] = useState(null)

  const { job, error: pollingError } = useJobPolling(jobId)
  const { entries, addEntry, clearHistory } = useUploadHistory()
  const recordedJobId = useRef(null)

  const isJobActive = job && (job.status === 'pending' || job.status === 'processing')
  const isLocked = isUploading || Boolean(isJobActive)
  const showResults = job?.status === 'completed'

  useEffect(() => {
    if (!file) {
      setPreviewUrl(null)
      return undefined
    }
    const url = URL.createObjectURL(file)
    setPreviewUrl(url)
    return () => URL.revokeObjectURL(url)
  }, [file])

  // Log each finished job to the "today's uploads" history, exactly once.
  useEffect(() => {
    if (!job || !jobId || recordedJobId.current === jobId) return

    if (job.status === 'completed') {
      const smallest = pickSmallestOutput(job.outputs)
      addEntry({
        fileName: job.file_name,
        status: 'success',
        savedPercent: percentSmaller(job.original_size_bytes, smallest.size_bytes),
        optimizedBytes: smallest.size_bytes,
      })
      recordedJobId.current = jobId
    } else if (job.status === 'failed') {
      addEntry({ fileName: job.file_name, status: 'failed' })
      recordedJobId.current = jobId
    }
  }, [job, jobId, addEntry])

  const handleSelectFile = (picked) => {
    setSubmitError(null)
    setFile(picked)
  }

  const handleClearFile = () => {
    setFile(null)
    setJobId(null)
    setSubmitError(null)
    recordedJobId.current = null
  }

  const handleToggleFormat = (value) => {
    setFormats((current) =>
      current.includes(value) ? current.filter((item) => item !== value) : [...current, value],
    )
  }

  const handleReset = () => {
    setFile(null)
    setJobId(null)
    setSubmitError(null)
    setUploadPercent(0)
    recordedJobId.current = null
  }

  const canSubmit = file && formats.length > 0 && !isLocked

  const handleSubmit = async (event) => {
    event.preventDefault()
    if (!canSubmit) return

    const validationError = validateFile(file)
    if (validationError) {
      setSubmitError(validationError)
      return
    }

    setIsUploading(true)
    setUploadPercent(0)
    setSubmitError(null)

    try {
      const result = await uploadImage({
        file,
        width,
        quality,
        formats,
        onProgress: setUploadPercent,
      })
      setJobId(result.job_id)
    } catch (err) {
      setSubmitError(toFriendlyError(err))
      addEntry({ fileName: file.name, status: 'failed' })
    } finally {
      setIsUploading(false)
    }
  }

  const submitLabel = useMemo(() => {
    if (isUploading) return 'Uploading…'
    if (isJobActive) return 'Processing…'
    return 'Optimize image'
  }, [isUploading, isJobActive])

  return (
    <div className="mx-auto flex max-w-xl flex-col gap-8 px-4 py-10 sm:py-16">
      {!showResults && (
        <form onSubmit={handleSubmit} className="flex flex-col gap-5">
          <Dropzone
            file={file}
            previewUrl={previewUrl}
            onSelect={handleSelectFile}
            onClear={handleClearFile}
            disabled={isLocked}
          />

          <OptionsPanel
            width={width}
            quality={quality}
            formats={formats}
            onWidthChange={setWidth}
            onQualityChange={setQuality}
            onToggleFormat={handleToggleFormat}
            disabled={isLocked}
          />

          {submitError && (
            <p className="rounded-xl border border-clay/40 bg-clay/10 px-4 py-3 text-xs text-clay">
              {submitError}
            </p>
          )}

          <button
            type="submit"
            disabled={!canSubmit}
            className="flex items-center justify-center gap-2 rounded-full bg-quill px-5 py-3 text-sm font-semibold text-ink-950 transition hover:bg-quill-soft disabled:cursor-not-allowed disabled:bg-ink-700 disabled:text-slate-500"
          >
            <FeatherMark className="h-4 w-4" animated={isLocked} />
            {submitLabel}
          </button>
        </form>
      )}

      {isUploading && <UploadProgress percent={uploadPercent} />}

      {!isUploading && !showResults && (
        <JobStatus job={job} pollingError={pollingError} />
      )}

      {showResults && (
        <>
          <OptimizationStats originalBytes={job.original_size_bytes} outputs={job.outputs} />
          <ResultsPanel job={job} onReset={handleReset} />
        </>
      )}

      <HistoryList entries={entries} onClear={clearHistory} />

      <p className="text-center text-[11px] text-slate-500">
        Files are processed in the background and automatically removed after their retention window.
      </p>
    </div>
  )
}
