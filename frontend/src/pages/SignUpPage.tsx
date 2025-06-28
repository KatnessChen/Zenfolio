import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import { Card, CardContent } from '@/components/ui/card'
import { GoogleIcon } from '@/components/icons'
import { Link } from 'react-router-dom'
import { Title } from '@/components/ui/title'

export default function SignUpPage() {
  return (
    <div className="min-h-screen flex flex-col items-center justify-center py-16">
      <div className="container mx-auto px-4 max-w-md">
        {/* Page Title */}
        <div className="text-center mb-8 space-y-4">
          <Title as="h1">Create Your Zenfolio Account</Title>
        </div>

        {/* Sign Up Form Card - Dark Grey-Green Background */}
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
                <Label htmlFor="firstName" className="text-foreground">
                  First Name
                </Label>
                <Input
                  id="firstName"
                  type="text"
                  placeholder="Enter your first name"
                  required
                  className="bg-card text-foreground border-border placeholder:text-muted-foreground"
                />
              </div>

              <div className="space-y-2">
                <Label htmlFor="lastName" className="text-foreground">
                  Last Name
                </Label>
                <Input
                  id="lastName"
                  type="text"
                  placeholder="Enter your last name (optional)"
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
                  placeholder="Create a password"
                  required
                  className="bg-card text-foreground border-border placeholder:text-muted-foreground"
                />
              </div>

              <div className="space-y-2">
                <Label htmlFor="confirmPassword" className="text-foreground">
                  Confirm Password
                </Label>
                <Input
                  id="confirmPassword"
                  type="password"
                  placeholder="Confirm your password"
                  required
                  className="bg-card text-foreground border-border placeholder:text-muted-foreground"
                />
              </div>
            </div>

            {/* Primary Sign Up Button - Brighter Sage Green background, Pure Light Gray text */}
            <Button
              size="lg"
              className="w-full bg-primary text-primary-foreground hover:bg-primary/90"
            >
              Sign Up
            </Button>

            {/* Divider - Medium Grey-Green text */}
            <div className="relative">
              <div className="absolute inset-0 flex items-center">
                <span className="w-full border-t border-border" />
              </div>
              <div className="relative flex justify-center text-xs uppercase">
                <span className="bg-card px-2 text-muted-foreground">OR</span>
              </div>
            </div>

            {/* Google Sign Up Button - Pure Light Gray text */}
            <Button
              variant="outline"
              size="lg"
              className="w-full space-x-3 text-foreground border-border hover:bg-card/50"
            >
              <GoogleIcon />
              <span>Sign Up with Google</span>
            </Button>

            {/* Login Link - Medium Grey-Green text */}
            <div className="text-center">
              <p className="text-sm text-muted-foreground">
                Already have an account?{' '}
                <Link
                  to="/login"
                  className="text-primary hover:text-primary/90 transition-colors font-medium"
                >
                  Login here
                </Link>
              </p>
            </div>
          </CardContent>
        </Card>
      </div>
    </div>
  )
}
