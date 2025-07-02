import { useSelector } from 'react-redux'
import ExtractedDataReviewPage from '@/pages/ExtractedDataReviewPage'
import type { RootState } from '@/store'
import ImageUploadPage from '@/pages/ImageUploadPage'

export default function UploadImageGuard() {
  const files = useSelector((state: RootState) => state.fileProcessing.files)

  if (files.length === 0) return <ImageUploadPage />
  return <ExtractedDataReviewPage />
}
