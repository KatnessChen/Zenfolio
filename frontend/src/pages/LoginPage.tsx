import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import { Card, CardContent } from '@/components/ui/card'
import { GoogleIcon } from '@/components/icons'
import { Title } from '@/components/ui/title'
import { ROUTES } from '@/constants'
import { Link } from '@/components/ui/link'

export default function LoginPage() {
  return (
    <div className="min-h-screen flex flex-col items-center justify-center py-16">
      <div className="container mx-auto px-4 max-w-md">
        {/* Page Title */}
        <div className="text-center mb-8 space-y-4">
          <Title as="h1">Login to Your Zenfolio Account</Title>
        </div>

        {/* Login Form Card - Dark Grey-Green Background */}
        <Card className="bg-card border-border">
          <CardContent className="space-y-6 p-8">
            {/* Form Fields */}
            <div className="space-y-4">
              <div className="space-y-2">
                <Label htmlFor="email" className="text-foreground">
                  Email Address
                </Label>
                <Input
                  id="email"
                  type="email"
                  placeholder="Enter your email"
                  required
                  className="bg-card text-foreground border-border placeholder:text-muted-foreground"
                />
              </div>

              <div className="space-y-2">
                <Label htmlFor="password" className="text-foreground">
                  Password
                </Label>
                <Input
                  id="password"
                  type="password"
                  placeholder="Enter your password"
                  required
                  className="bg-card text-foreground border-border placeholder:text-muted-foreground"
                />
              </div>
            </div>

            {/* Forgot Password Link */}
            <div className="text-right">
              <Link
                to="/forgot-password"
                className="text-sm text-muted-foreground hover:text-primary transition-colors"
              >
                Forgot Password?
              </Link>
            </div>

            {/* Primary Login Button */}
            <Button size="lg" className="w-full">
              Login
            </Button>

            {/* Divider */}
            <div className="relative">
              <div className="absolute inset-0 flex items-center">
                <span className="w-full border-t border-border" />
              </div>
              <div className="relative flex justify-center text-xs uppercase">
                <span className="bg-card px-2 text-muted-foreground">OR</span>
              </div>
            </div>

            {/* Google Login Button */}
            <Button variant="outline" size="lg" className="w-full space-x-3">
              <GoogleIcon />
              <span>Login with Google</span>
            </Button>

            {/* Sign Up Link */}
            <div className="text-center">
              <p className="text-sm text-muted-foreground">
                Don't have an account? <Link to={ROUTES.SIGN_UP}>Sign Up now</Link>
              </p>
            </div>
          </CardContent>
        </Card>
      </div>
    </div>
  )
}
