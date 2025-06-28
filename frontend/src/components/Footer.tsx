export default function Footer() {
  return (
    <footer className="border-t border-border/30 py-6">
      <div className="container mx-auto px-4 text-center">
        <div className="space-y-2">
          <p className="text-sm text-muted-foreground">© 2025 Zenfolio. All rights reserved.</p>
          <a
            href="/contact"
            className="text-sm text-primary hover:text-primary/90 transition-colors"
          >
            Contact Us
          </a>
        </div>
      </div>
    </footer>
  )
}
