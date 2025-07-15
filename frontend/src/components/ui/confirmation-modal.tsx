import { Card, CardContent } from '@/components/ui/card'
import { Button } from '@/components/ui/button'
import { type ReactNode } from 'react'

interface ConfirmationModalProps {
  isOpen: boolean
  title: string
  message: string | ReactNode
  confirmText?: string
  cancelText?: string
  confirmVariant?: 'default' | 'destructive'
  onConfirm: () => void
  onCancel: () => void
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
            <Button onClick={onConfirm} variant={confirmVariant} disabled={isLoading}>
              {isLoading ? 'Loading...' : confirmText}
            </Button>
          </div>
        </CardContent>
      </Card>
    </div>
  )
}
