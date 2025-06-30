import { useState, useCallback } from 'react'
import { useNavigate } from 'react-router-dom'
import { ROUTES } from '@/constants'
import { Card, CardContent } from '@/components/ui/card'
import { Button } from '@/components/ui/button'
import FolderPlusIcon from '@/components/icons/FolderPlusIcon'
import { Title } from '@/components/ui/title'

export default function ImageUploadPage() {
  const navigate = useNavigate()
  const [isDragOver, setIsDragOver] = useState(false)
  const [selectedFiles, setSelectedFiles] = useState<File[]>([])

  const handleDragOver = useCallback((e: React.DragEvent) => {
    e.preventDefault()
    setIsDragOver(true)
  }, [])

  const handleDragLeave = useCallback((e: React.DragEvent) => {
    e.preventDefault()
    setIsDragOver(false)
  }, [])

  const handleUpload = useCallback(async () => {
    if (selectedFiles.length === 0) return

    // Navigate to processing page
    navigate(ROUTES.TRANSACTIONS_UPLOAD_PROCESSING, {
      state: { files: selectedFiles.map((f) => f.name) },
    })

    // TODO: Implement actual file upload to backend
    console.log('Uploading files:', selectedFiles)
  }, [navigate, selectedFiles])

  const handleDrop = useCallback((e: React.DragEvent) => {
    e.preventDefault()
    setIsDragOver(false)

    const files = Array.from(e.dataTransfer.files).filter((file) =>
      ['image/jpeg', 'image/png', 'application/pdf'].includes(file.type)
    )

    setSelectedFiles((prev) => {
      const newFiles = [...prev, ...files]
      if (newFiles.length > 10) {
        alert('Maximum 10 files allowed per upload')
        return prev
      }
      return newFiles
    })
  }, [])

  const handleFileSelect = useCallback((e: React.ChangeEvent<HTMLInputElement>) => {
    const files = Array.from(e.target.files || []).filter((file) =>
      ['image/jpeg', 'image/png', 'application/pdf'].includes(file.type)
    )

    setSelectedFiles((prev) => {
      const newFiles = [...prev, ...files]
      if (newFiles.length > 10) {
        alert('Maximum 10 files allowed per upload')
        return prev
      }
      return newFiles
    })
  }, [])

  const handleRemoveFile = useCallback((index: number) => {
    setSelectedFiles((prev) => prev.filter((_, i) => i !== index))
  }, [])

  const handleBrowseClick = () => {
    const input = document.createElement('input')
    input.type = 'file'
    input.multiple = true
    input.accept = 'image/jpeg,image/png,application/pdf'
    input.onchange = (e) => {
      const target = e.target as HTMLInputElement
      const changeEvent = {
        target,
        currentTarget: target,
      } as React.ChangeEvent<HTMLInputElement>
      handleFileSelect(changeEvent)
    }
    input.click()
  }

  const handleManualAdd = useCallback(() => {
    // Navigate to manual transaction entry page
    navigate(ROUTES.TRANSACTIONS_MANUAL_ADD)
  }, [navigate])

  return (
    <div className="min-h-screen bg-background">
      {/* Main Content */}
      <main className="container mx-auto py-8 px-4 max-w-4xl">
        {/* Page Title */}
        <div className="text-center mb-8">
          <Title as="h1" className="mb-4">
            Upload Broker Screenshots
          </Title>
          <p className="text-muted-foreground">
            Our AI will automatically extract data from your trade history.
          </p>
        </div>

        {/* SECTION 1: IMAGE UPLOAD */}
        <div className="mb-12">
          {/* File Upload Area */}
          <Card className="mb-4">
            <CardContent className="p-8">
              <div
                className={`border-2 border-dashed rounded-lg p-12 text-center transition-all duration-200 ${
                  isDragOver
                    ? 'border-primary bg-primary/5'
                    : 'border-border hover:border-primary/50'
                }`}
                onDragOver={handleDragOver}
                onDragLeave={handleDragLeave}
                onDrop={handleDrop}
              >
                <div className="space-y-4">
                  {/* Upload Icon */}
                  <div className="flex justify-center">
                    <div className="w-16 h-16 bg-muted rounded-full flex items-center justify-center">
                      <FolderPlusIcon size={32} className="text-foreground" />
                    </div>
                  </div>

                  {/* Drag & Drop Text */}
                  <div>
                    <p className="text-lg text-muted-foreground mb-2">
                      Drag & drop your screenshots here
                    </p>
                    <p className="text-muted-foreground mb-4">OR</p>

                    {/* Browse Button */}
                    <Button variant="secondary" onClick={handleBrowseClick} className="mb-4">
                      Browse Files
                    </Button>

                    {/* Upload Limits */}
                    <p className="text-sm text-muted-foreground">
                      Accepts JPG, PNG, PDF files. Max 10 files per upload.
                    </p>
                  </div>
                </div>
              </div>

              {/* Selected Files Display */}
              {selectedFiles.length > 0 && (
                <div className="mt-6">
                  <h6 className="text-base font-medium text-foreground mb-3 text-center">
                    Selected Files
                  </h6>
                  <div className="space-y-2">
                    {selectedFiles.map((file, index) => (
                      <div
                        key={index}
                        className="flex items-center justify-between text-foreground p-2 rounded-lg hover:bg-muted/50 transition-colors"
                      >
                        <span className="text-sm">• {file.name}</span>
                        <button
                          onClick={() => handleRemoveFile(index)}
                          className="text-muted-foreground hover:text-destructive transition-colors ml-2"
                          title="Remove file"
                        >
                          ✕
                        </button>
                      </div>
                    ))}
                  </div>
                </div>
              )}

              {/* Upload Button */}
              {selectedFiles.length > 0 && (
                <div className="mt-6 flex justify-center">
                  <Button variant="default" onClick={handleUpload}>
                    Upload
                  </Button>
                </div>
              )}
            </CardContent>
          </Card>

          {/* Tips Section */}
          <div className="px-6">
            <p className="text-muted-foreground">
              Use clear, high-resolution screenshots. Capture the full transaction table. Avoid
              cropping or watermarks.
            </p>
          </div>
        </div>

        {/* SECTION 2: MANUAL INPUT */}
        <div className="border-t border-border pt-8">
          <div className="text-center">
            <p className="text-muted-foreground mb-4">OR</p>
            <Button variant="default" onClick={handleManualAdd} className="px-8">
              Manually Add Transaction
            </Button>
          </div>
        </div>
      </main>
    </div>
  )
}
