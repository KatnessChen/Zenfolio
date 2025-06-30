import { Button } from '@/components/ui/button'
import { Logo } from '@/components/ui/logo'
import { GoogleIcon } from '@/components/icons'
import { Link } from 'react-router-dom'
import { Title } from '@/components/ui/title'

export default function HomePage() {
  return (
    <div className="min-h-screen flex flex-col">
      {/* Hero Section - Split Layout */}
      <main className="flex-1 flex items-center justify-center py-16">
        <div className="container mx-auto px-4">
          <div className="grid grid-cols-1 lg:grid-cols-2 gap-12 items-center max-w-6xl mx-auto">
            {/* Left Half - Information Area */}
            <div className="space-y-8 text-center lg:text-left">
              <div className="flex justify-center lg:justify-start">
                <Logo size="lg" className="text-4xl font-bold" />
              </div>

              <div className="space-y-6">
                <Title as="h1" className="text-4xl lg:text-5xl font-bold leading-tight">
                  Trade across multiple brokers?
                </Title>
                <p className="text-xl text-muted-foreground leading-relaxed">
                  Unify your portfolio with Zenfolio. Upload your trade history, let our AI extract
                  and summarize your performance effortlessly.
                </p>
              </div>
            </div>

            {/* Right Half - Action Area */}
            <div className="space-y-6 max-w-md mx-auto lg:mx-0">
              <div className="bg-card rounded-lg border border-border p-8 space-y-6">
                {/* Primary Login Button */}
                <Link to="/login" className="block">
                  <Button size="lg" className="w-full py-6">
                    Login
                  </Button>
                </Link>

                {/* Secondary Sign Up Button */}
                <Link to="/sign-up" className="block">
                  <Button variant="outline" size="lg" className="w-full py-6">
                    Sign Up
                  </Button>
                </Link>

                {/* Divider */}
                <div className="relative">
                  <div className="absolute inset-0 flex items-center">
                    <span className="w-full border-t border-border" />
                  </div>
                  <div className="relative flex justify-center text-xs uppercase">
                    <span className="bg-card px-2 text-muted-foreground">OR</span>
                  </div>
                </div>

                {/* Google Sign In Button */}
                <Button variant="outline" size="lg" className="w-full py-6 space-x-3">
                  <GoogleIcon />
                  <span>Continue with Google</span>
                </Button>
              </div>
            </div>
          </div>
        </div>
      </main>
    </div>
  )
}
