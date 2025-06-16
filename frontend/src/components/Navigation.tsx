import { Link } from 'react-router-dom'
import { Button } from '@/components/ui/button'

export default function Navigation() {
  return (
    <nav className="border-b">
      <div className="container mx-auto px-4 py-4">
        <div className="flex justify-between items-center">
          <Link to="/" className="text-xl font-bold">
            Transaction Tracker
          </Link>

          <div className="flex space-x-4">
            <Button variant="ghost" asChild>
              <Link to="/">Home</Link>
            </Button>
            <Button variant="ghost" asChild>
              <Link to="/extract">Extract</Link>
            </Button>
            <Button variant="ghost" asChild>
              <Link to="/history">History</Link>
            </Button>
          </div>
        </div>
      </div>
    </nav>
  )
}
