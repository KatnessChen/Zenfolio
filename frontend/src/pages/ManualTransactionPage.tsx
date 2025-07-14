import { TransactionEditor } from '@/components/TransactionEditor'
import type { TransactionData } from '@/types'
import { useLocation } from 'react-router-dom'

export default function ManualTransactionPage() {
  const location = useLocation()

  const initialFromState = location.state?.initial as TransactionData | undefined

  return (
    <div className="min-h-screen bg-background flex items-center justify-center">
      <TransactionEditor initial={initialFromState} />
    </div>
  )
}
