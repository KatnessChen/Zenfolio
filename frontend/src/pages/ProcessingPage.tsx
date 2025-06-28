import { useEffect } from 'react'
import { useNavigate, useLocation } from 'react-router-dom'
import { ROUTES } from '@/constants/routes'

interface LocationState {
  files: string[]
}

export default function ProcessingPage() {
  const navigate = useNavigate()
  const location = useLocation()
  const state = location.state as LocationState

  useEffect(() => {
    // Simulate processing time (3-5 seconds)
    const processingTime = 3000 + Math.random() * 2000

    const timer = setTimeout(() => {
      // Navigate to review page with processed data
      navigate(ROUTES.TRANSACTIONS_UPLOAD_REVIEW, {
        state: {
          files: state?.files || [],
          processedData: generateMockProcessedData(),
        },
      })
    }, processingTime)

    return () => clearTimeout(timer)
  }, [navigate, state])

  // Mock function to generate processed transaction data
  const generateMockProcessedData = () => {
    return [
      {
        id: '1',
        ticker: 'AAPL',
        tradeType: 'Buy',
        quantity: 100,
        price: 150.25,
        amount: 15025.0,
        tradeDate: '2024-01-15',
        broker: 'Fidelity',
        currency: 'USD',
        userNotes: '',
        confidence: 0.95,
      },
      {
        id: '2',
        ticker: 'GOOGL',
        tradeType: 'Sell',
        quantity: 50,
        price: 2750.8,
        amount: 137540.0,
        tradeDate: '2024-01-14',
        broker: 'Schwab',
        currency: 'USD',
        userNotes: '',
        confidence: 0.92,
      },
    ]
  }

  return (
    <div className="min-h-screen bg-background flex items-center justify-center">
      <div className="text-center space-y-6 max-w-md mx-auto px-4">
        {/* Loading Animation */}
        <div className="flex justify-center mb-8">
          <div className="relative">
            {/* Zenfolio Logo with Spin Animation */}
            <div className="w-16 h-16 bg-primary rounded-full flex items-center justify-center animate-spin">
              <span className="text-primary-foreground font-bold text-xl">Z</span>
            </div>

            {/* Outer Ring Animation */}
            <div className="absolute inset-0 border-4 border-primary/20 rounded-full animate-pulse"></div>
          </div>
        </div>

        {/* Loading Text */}
        <div className="space-y-3">
          <h2 className="text-2xl font-semibold text-foreground">Processing your screenshots...</h2>
          <p className="text-muted-foreground">
            This may take a moment, please do not close this page.
          </p>
        </div>

        {/* Processing Steps */}
        <div className="mt-8 space-y-3 text-left">
          <div className="flex items-center space-x-3">
            <div className="w-2 h-2 bg-primary rounded-full animate-pulse"></div>
            <span className="text-sm text-muted-foreground">Analyzing images...</span>
          </div>
          <div className="flex items-center space-x-3">
            <div className="w-2 h-2 bg-muted-foreground/40 rounded-full"></div>
            <span className="text-sm text-muted-foreground/60">Extracting transaction data...</span>
          </div>
          <div className="flex items-center space-x-3">
            <div className="w-2 h-2 bg-muted-foreground/40 rounded-full"></div>
            <span className="text-sm text-muted-foreground/60">Validating information...</span>
          </div>
        </div>

        {/* File Info */}
        {state?.files && (
          <div className="mt-8 p-4 bg-muted rounded-lg">
            <p className="text-sm text-muted-foreground mb-2">Processing files:</p>
            <div className="space-y-1">
              {state.files.map((filename, index) => (
                <p key={index} className="text-xs text-foreground font-mono">
                  {filename}
                </p>
              ))}
            </div>
          </div>
        )}
      </div>
    </div>
  )
}
