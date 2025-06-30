import { Link } from '@/components/ui/link'
import { ROUTES } from '@/constants'

export default function Footer() {
  return (
    <footer className="border-t border-border/30 py-6">
      <div className="container mx-auto px-4 text-center">
        <div className="space-y-2">
          <p className="text-sm text-muted-foreground">© 2025 Zenfolio. All rights reserved.</p>
          <Link to={ROUTES.CONTACT} className="text-sm">
            Contact Us
          </Link>
        </div>
      </div>
    </footer>
  )
}
