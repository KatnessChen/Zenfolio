import { useEffect, useState, useCallback, useRef } from 'react'
import { useNavigate, useLocation } from 'react-router-dom'
import { useDispatch, useSelector } from 'react-redux'
import { ROUTES } from '@/constants'
import { ExtractService } from '@/services/extract.service'
import { ClockIcon, SpinnerIcon, CheckIcon, XIcon } from '@/components/icons'
import { initializeFiles, updateFileStatus } from '@/store/fileProcessingSlice'
import { filesToSerializable } from '@/utils/fileUtils'
import type { RootState, AppDispatch } from '@/store'
import type { ExtractResponse } from '@/types'
import {
  FILE_STATUS_PENDING,
  FILE_STATUS_PROCESSING,
  FILE_STATUS_COMPLETED,
  FILE_STATUS_ERROR,
} from '@/constants/fileProcessingStatus'

interface LocationState {
  files: File[]
}

export default function ProcessingPage() {
  const navigate = useNavigate()
  const location = useLocation()
  const state = location.state as LocationState
  const dispatch = useDispatch<AppDispatch>()
  const { fileStates } = useSelector((state: RootState) => state.fileProcessing)

  const [hasNavigated, setHasNavigated] = useState(false)
  const hasStartedExtractionRef = useRef(false)

  const startParallelExtraction = useCallback(
    async (files: File[]) => {
      try {
        // Update all files to processing status
        files.forEach((_, index) => {
          dispatch(updateFileStatus({ fileIndex: index, status: 'processing' }))
        })

        // Process files in parallel with progress updates
        await ExtractService.extractTransactionsParallel(
          files,
          (fileIndex: number, result: ExtractResponse, error?: string) => {
            const status = result.success ? 'completed' : 'error'
            dispatch(
              updateFileStatus({
                fileIndex,
                status,
                result: result.data,
                error: error || (!result.success ? result.message : undefined),
              })
            )

            // Navigate as soon as ANY file is successfully processed OR all files are done (including all errors)
            const currentStates = [...fileStates]
            currentStates[fileIndex] = {
              ...currentStates[fileIndex],
              status,
              result: result.data,
              error: error || (!result.success ? result.message : undefined),
            }

            const hasAnyCompleted = currentStates.some((fs) => fs.status === 'completed')
            const allFilesProcessed = currentStates.every(
              (fs) => fs.status === 'completed' || fs.status === 'error'
            )

            if ((hasAnyCompleted || allFilesProcessed) && !hasNavigated) {
              setHasNavigated(true)

              // Navigate to review page after a short delay
              setTimeout(() => {
                navigate(ROUTES.TRANSACTIONS_UPLOAD_REVIEW)
              }, 1000)
            }
          }
        )
      } catch (error) {
        console.error('Extraction failed:', error)
        files.forEach((_, index) => {
          dispatch(
            updateFileStatus({
              fileIndex: index,
              status: 'error',
              error: 'Extraction failed',
            })
          )
        })
      }
    },
    [navigate, hasNavigated, dispatch, fileStates]
  )

  useEffect(() => {
    if (!state?.files || state.files.length === 0) {
      navigate(ROUTES.TRANSACTIONS_UPLOAD)
      return
    }

    // Prevent multiple executions
    if (hasStartedExtractionRef.current) return

    // Convert files to serializable format and initialize in store
    const initializeAndStart = async () => {
      try {
        const serializableFiles = await filesToSerializable(state.files)
        dispatch(initializeFiles({ files: serializableFiles }))

        // Start parallel extraction with original File objects
        startParallelExtraction(state.files)
      } catch (error) {
        console.error('Failed to convert files to serializable format:', error)
        navigate(ROUTES.TRANSACTIONS_UPLOAD)
      }
    }

    initializeAndStart()

    // Mark extraction as started
    hasStartedExtractionRef.current = true
  }, [navigate, state, startParallelExtraction, dispatch])

  const completedCount = fileStates.filter((fs) => fs.status === 'completed').length
  const errorCount = fileStates.filter((fs) => fs.status === 'error').length
  const totalCount = fileStates.length

  const getStatusIcon = (status: (typeof fileStates)[0]['status']) => {
    switch (status) {
      case FILE_STATUS_PENDING:
        return <ClockIcon size={24} />
      case FILE_STATUS_PROCESSING:
        return <SpinnerIcon size={24} />
      case FILE_STATUS_COMPLETED:
        return <CheckIcon size={24} />
      case FILE_STATUS_ERROR:
        return <XIcon size={24} />
      default:
        return <ClockIcon size={24} />
    }
  }

  return (
    <div className="min-h-screen bg-background flex items-center justify-center">
      <div className="text-center space-y-6 max-w-md mx-auto px-4">
        {/* Loading Animation */}
        <div className="flex justify-center mb-8">
          <div className="relative">
            {/* Zenfolio Logo with Spin Animation */}
            <div className="w-16 h-16 bg-primary rounded-full flex items-center justify-center animate-spin">
              <span className="text-primary-foreground font-bold text-xl">Z</span>
            </div>

            {/* Outer Ring Animation */}
            <div className="absolute inset-0 border-4 border-primary/20 rounded-full animate-pulse"></div>
          </div>
        </div>

        {/* Loading Text */}
        <div className="space-y-3">
          <h2 className="text-2xl font-semibold text-foreground">Processing your screenshots...</h2>
          <p className="text-muted-foreground">
            {totalCount > 1
              ? `Processing ${totalCount} files in parallel. You'll be redirected when the first file is ready.`
              : 'This may take a moment, please do not close this page.'}
          </p>
          <p className="text-sm text-muted-foreground">
            The review page will appear once the first file is processed.
          </p>
          {totalCount > 1 && (
            <p className="text-sm text-muted-foreground">
              <span className="inline-flex items-center gap-1">
                <CheckIcon className="w-3 h-3" />
                {completedCount} completed
              </span>
              {' • '}
              <span className="inline-flex items-center gap-1">
                <XIcon className="w-3 h-3" />
                {errorCount} failed
              </span>
              {' • '}
              <span className="inline-flex items-center gap-1">
                <ClockIcon className="w-3 h-3" />
                {totalCount - completedCount - errorCount} remaining
              </span>
            </p>
          )}
        </div>

        {/* File Processing Status */}
        {fileStates.length > 0 && (
          <div className="mt-8 p-4 bg-muted rounded-lg max-w-md">
            <div className="space-y-2 max-h-64 overflow-y-auto">
              {fileStates.map((fileState, index) => (
                <div key={index} className="flex items-center justify-between text-xs">
                  <span className="text-foreground font-mono truncate flex-1 mr-2">
                    {fileState.file.name}
                  </span>
                  <div className="flex items-center space-x-2">
                    <span>{getStatusIcon(fileState.status)}</span>
                    {fileState.status === 'completed' && fileState.result && (
                      <span className="text-green-600 text-xs">
                        {fileState.result.transaction_count} found
                      </span>
                    )}
                    {fileState.status === FILE_STATUS_ERROR && (
                      <span className="text-red-600 text-xs" title={fileState.error}>
                        Error
                      </span>
                    )}
                  </div>
                </div>
              ))}
            </div>
          </div>
        )}
      </div>
    </div>
  )
}
