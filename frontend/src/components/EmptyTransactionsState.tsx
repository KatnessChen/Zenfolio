import { Button } from '@/components/ui/button'
import { PlusIcon } from '@/components/icons'

interface EmptyTransactionsStateProps {
  onUploadClick?: () => void
}

export default function EmptyTransactionsState({
  onUploadClick = () => console.log('Upload image flow'),
}: EmptyTransactionsStateProps) {
  return (
    <div className="flex flex-col items-center justify-center py-16 space-y-6">
      <div className="text-center space-y-4">
        <div className="w-16 h-16 mx-auto bg-muted rounded-full flex items-center justify-center">
          <PlusIcon size={24} className="text-muted-foreground" />
        </div>
        <div className="space-y-2">
          <h3 className="text-xl font-medium text-foreground">No transactions yet</h3>
          <p className="text-muted-foreground max-w-md">
            Start by uploading transaction images to automatically extract and organize your trading
            data.
          </p>
        </div>
        <Button variant="default" onClick={onUploadClick} className="mt-4">
          <PlusIcon size={16} />
          <span>Upload Transaction Image</span>
        </Button>
      </div>
    </div>
  )
}
