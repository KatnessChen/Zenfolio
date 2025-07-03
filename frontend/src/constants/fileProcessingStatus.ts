// File processing status constants
export const FILE_STATUS_PENDING = 'pending'
export const FILE_STATUS_PROCESSING = 'processing'
export const FILE_STATUS_COMPLETED = 'completed'
export const FILE_STATUS_ERROR = 'error'

export const FILE_PROCESSING_STATUSES = [
  FILE_STATUS_PENDING,
  FILE_STATUS_PROCESSING,
  FILE_STATUS_COMPLETED,
  FILE_STATUS_ERROR,
] as const

export type FileProcessingStatus = (typeof FILE_PROCESSING_STATUSES)[number]
