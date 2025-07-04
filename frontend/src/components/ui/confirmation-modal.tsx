import { Card, CardContent } from '@/components/ui/card'
import { Button } from '@/components/ui/button'

interface ConfirmationModalProps {
  /** Whether the modal is visible */
  isOpen: boolean
  /** Modal title */
  title: string
  /** Modal description/message */
  message: string
  /** Text for the confirm button */
  confirmText?: string
  /** Text for the cancel button */
  cancelText?: string
  /** Variant for the confirm button */
  confirmVariant?: 'default' | 'destructive'
  /** Callback when confirm button is clicked */
  onConfirm: () => void
  /** Callback when cancel button is clicked */
  onCancel: () => void
  /** Whether the confirm action is loading */
  isLoading?: boolean
}

export function ConfirmationModal({
  isOpen,
  title,
  message,
  confirmText = 'Confirm',
  cancelText = 'Cancel',
  confirmVariant = 'default',
  onConfirm,
  onCancel,
  isLoading = false,
}: ConfirmationModalProps) {
  if (!isOpen) return null

  const confirmButtonClass =
    confirmVariant === 'destructive'
      ? 'bg-red-600 hover:bg-red-700 text-white'
      : 'bg-blue-600 hover:bg-blue-700 text-white'

  return (
    <div className="fixed inset-0 bg-black/50 flex items-center justify-center z-50">
      <Card className="w-[400px]">
        <CardContent className="p-6">
          <h3 className="text-lg font-semibold mb-4">{title}</h3>
          <p className="text-muted-foreground mb-6">{message}</p>
          <div className="flex justify-end gap-3">
            <Button onClick={onCancel} variant="outline" disabled={isLoading}>
              {cancelText}
            </Button>
            <Button onClick={onConfirm} className={confirmButtonClass} disabled={isLoading}>
              {isLoading ? 'Loading...' : confirmText}
            </Button>
          </div>
        </CardContent>
      </Card>
    </div>
  )
}
