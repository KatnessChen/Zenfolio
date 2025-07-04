import {
  FILE_STATUS_PENDING,
  FILE_STATUS_PROCESSING,
  FILE_STATUS_COMPLETED,
  FILE_STATUS_ERROR,
} from '@/constants/fileProcessingStatus'
import type { FileProcessingState } from '@/types'

interface SerializableFile {
  name: string
  size: number
  type: string
  lastModified: number
  dataUrl: string // Base64 encoded file data
}

/**
 * Convert a File object to a serializable format for Redux storage
 */
export const fileToSerializable = async (file: File): Promise<SerializableFile> => {
  return new Promise((resolve, reject) => {
    const reader = new FileReader()

    reader.onload = () => {
      resolve({
        name: file.name,
        size: file.size,
        type: file.type,
        lastModified: file.lastModified,
        dataUrl: reader.result as string,
      })
    }

    reader.onerror = () => {
      reject(new Error(`Failed to read file: ${file.name}`))
    }

    reader.readAsDataURL(file)
  })
}

/**
 * Convert multiple File objects to serializable format
 */
export const filesToSerializable = async (files: File[]): Promise<SerializableFile[]> => {
  const promises = files.map((file) => fileToSerializable(file))
  return Promise.all(promises)
}

/**
 * Convert a serializable file back to a File-like object for display purposes
 */
export const serializableToFile = (serializableFile: SerializableFile): File => {
  // Convert dataUrl back to File object
  const byteString = atob(serializableFile.dataUrl.split(',')[1])
  const mimeString = serializableFile.dataUrl.split(',')[0].split(':')[1].split(';')[0]

  const ab = new ArrayBuffer(byteString.length)
  const ia = new Uint8Array(ab)

  for (let i = 0; i < byteString.length; i++) {
    ia[i] = byteString.charCodeAt(i)
  }

  const blob = new Blob([ab], { type: mimeString })

  // Create a File object from the blob
  const file = new File([blob], serializableFile.name, {
    type: serializableFile.type,
    lastModified: serializableFile.lastModified,
  })

  return file
}

/**
 * Get a blob URL from a serializable file for image display
 */
export const getSerializableFileUrl = (serializableFile: SerializableFile): string => {
  return serializableFile.dataUrl
}

/**
 * Get the status of a file by index from fileStates array
 */
export const getFileStatus = (
  fileStates: FileProcessingState[],
  fileIndex: number
): FileProcessingState['status'] => {
  if (!fileStates[fileIndex]) return FILE_STATUS_COMPLETED
  return fileStates[fileIndex].status
}

/**
 * Get the appropriate status icon component name for a given status string
 * (Consumer should render the icon component by name)
 */
export const getStatusIconName = (
  status: string
): 'ClockIcon' | 'SpinnerIcon' | 'CheckIcon' | 'XIcon' => {
  switch (status) {
    case FILE_STATUS_PENDING:
      return 'ClockIcon'
    case FILE_STATUS_PROCESSING:
      return 'SpinnerIcon'
    case FILE_STATUS_COMPLETED:
      return 'CheckIcon'
    case FILE_STATUS_ERROR:
      return 'XIcon'
    default:
      return 'CheckIcon'
  }
}
