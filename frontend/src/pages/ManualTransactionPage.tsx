import { MultiTransactionEditor } from '@/components/MultiTransactionEditor'
import { Title } from '@/components/ui/title'
import { useLocation } from 'react-router-dom'
import type { TransactionData } from '@/types'

export default function ManualTransactionPage() {
  const location = useLocation()
  const isEditMode = Boolean(location.state?.initial as TransactionData | undefined)

  return (
    <div className="min-h-screen bg-background">
      <main className="container mx-auto py-8 px-4 max-w-6xl">
        {/* Page Title */}
        <div className="text-center mb-8">
          <Title as="h1" className="mb-4">
            {isEditMode ? 'Edit Transaction' : 'Add Transactions Manually'}
          </Title>
          <p className="text-muted-foreground">
            {isEditMode
              ? 'Edit transaction details with enhanced search and validation'
              : 'Add one or multiple transactions with enhanced search and validation'}
          </p>
        </div>

        {/* Multi-Transaction Editor */}
        <div className="w-full">
          <MultiTransactionEditor />
        </div>
      </main>
    </div>
  )
}
