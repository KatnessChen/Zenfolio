import { useEffect } from 'react'
import { useSelector } from 'react-redux'
import { useNavigate } from 'react-router-dom'
import { ROUTES } from '@/constants'
import ExtractedDataReviewPage from '@/pages/ExtractedDataReviewPage'
import type { RootState } from '@/store'

export default function UploadReviewGuard() {
  const navigate = useNavigate()
  const files = useSelector((state: RootState) => state.fileProcessing.files)

  useEffect(() => {
    if (files.length === 0) {
      navigate(ROUTES.TRANSACTIONS_UPLOAD, { replace: true })
    }
    // 如果有檔案就留在本頁
    else {
      // do nothing, stay on this page
    }
  }, [files, navigate])

  // 只有有檔案時才渲染內容
  if (files.length === 0) return null
  return <ExtractedDataReviewPage />
}
