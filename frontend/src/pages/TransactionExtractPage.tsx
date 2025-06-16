import { useState } from 'react'
import { Button } from '@/components/ui/button'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'

export default function TransactionExtractPage() {
  const [selectedFiles, setSelectedFiles] = useState<File[]>([])
  const [isLoading, setIsLoading] = useState(false)

  const handleFileChange = (event: React.ChangeEvent<HTMLInputElement>) => {
    const files = Array.from(event.target.files || [])
    setSelectedFiles(files)
  }

  const handleExtract = async () => {
    if (selectedFiles.length === 0) return

    setIsLoading(true)
    try {
      // TODO: Implement API call to extract transactions
      console.log('Extracting transactions from files:', selectedFiles)
      await new Promise((resolve) => setTimeout(resolve, 2000)) // Simulate API call
    } catch (error) {
      console.error('Error extracting transactions:', error)
    } finally {
      setIsLoading(false)
    }
  }

  return (
    <div className="container mx-auto py-10">
      <div className="max-w-2xl mx-auto space-y-8">
        <div className="text-center space-y-4">
          <h1 className="text-3xl font-bold tracking-tight">Upload Transactions</h1>
          <p className="text-lg text-muted-foreground">
            Upload trading screenshots to automatically extract transaction data using AI
          </p>
        </div>

        <Card>
          <CardHeader>
            <CardTitle>Upload Transaction Screenshots</CardTitle>
            <CardDescription>
              Select one or more screenshots of your trading transactions
            </CardDescription>
          </CardHeader>
          <CardContent className="space-y-6">
            <div className="space-y-2">
              <Label htmlFor="screenshots">Transaction Screenshots</Label>
              <Input
                id="screenshots"
                type="file"
                multiple
                accept="image/*"
                onChange={handleFileChange}
                className="cursor-pointer"
              />
              {selectedFiles.length > 0 && (
                <p className="text-sm text-muted-foreground">
                  {selectedFiles.length} file(s) selected
                </p>
              )}
            </div>

            <Button
              onClick={handleExtract}
              disabled={selectedFiles.length === 0 || isLoading}
              className="w-full"
            >
              {isLoading ? 'Extracting...' : 'Extract Transactions'}
            </Button>
          </CardContent>
        </Card>
      </div>
    </div>
  )
}
