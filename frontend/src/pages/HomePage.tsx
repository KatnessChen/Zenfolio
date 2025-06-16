import { Button } from '@/components/ui/button'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'

export default function HomePage() {
  return (
    <div className="container mx-auto py-10">
      <div className="flex flex-col items-center justify-center min-h-[600px] space-y-8">
        <div className="text-center space-y-4">
          <h1 className="text-5xl font-extrabold tracking-tighter bg-gradient-to-r from-blue-500 to-purple-600 text-transparent bg-clip-text">
            TradeVault
          </h1>
          <p className="text-xl text-muted-foreground max-w-2xl">
            Trading on multiple exchanges? Upload screenshots and let AI do the work of reading and
            organizing your transactions for you
          </p>
        </div>

        <div className="grid grid-cols-1 md:grid-cols-2 gap-6 w-full max-w-4xl">
          <Card>
            <CardHeader>
              <CardTitle>Upload Screenshots</CardTitle>
              <CardDescription>
                Upload trading screenshots to automatically extract transaction data
              </CardDescription>
            </CardHeader>
            <CardContent>
              <Button className="w-full">Start Extracting Transactions</Button>
            </CardContent>
          </Card>

          <Card>
            <CardHeader>
              <CardTitle>View Transactions</CardTitle>
              <CardDescription>Browse and manage your uploaded transaction history</CardDescription>
            </CardHeader>
            <CardContent>
              <Button variant="outline" className="w-full">
                View Transaction History
              </Button>
            </CardContent>
          </Card>
        </div>
      </div>
    </div>
  )
}
