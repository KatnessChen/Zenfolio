import * as React from 'react'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import { Card, CardContent } from '@/components/ui/card'
import { GoogleIcon } from '@/components/icons'
import { Title } from '@/components/ui/title'
import { ROUTES } from '@/constants'
import { Link } from '@/components/ui/link'

export default function SignUpPage() {
  const [email, setEmail] = React.useState('')
  const [firstName, setFirstName] = React.useState('')
  const [lastName, setLastName] = React.useState('')
  const [password, setPassword] = React.useState('')
  const [confirmPassword, setConfirmPassword] = React.useState('')
  const [passwordTooShort, setPasswordTooShort] = React.useState(false)
  const [passwordMismatch, setPasswordMismatch] = React.useState(false)
  const isFormValid =
    email.trim() && firstName.trim() && password.length >= 8 && password === confirmPassword

  function handleSignUp(e: React.FormEvent) {
    e.preventDefault()
    setPasswordTooShort(password.length > 0 && password.length < 8)
    setPasswordMismatch(confirmPassword.length > 0 && password !== confirmPassword)
    if (!isFormValid) return
    // TODO: submit form logic here
  }

  return (
    <div className="min-h-screen flex flex-col items-center justify-center py-16">
      <div className="container mx-auto px-4 max-w-md">
        {/* Page Title */}
        <div className="text-center mb-8 space-y-4">
          <Title as="h1">Create Your Zenfolio Account</Title>
        </div>

        {/* Sign Up Form Card */}
        <Card className="bg-card border-border">
          <CardContent className="space-y-6 p-8">
            <form onSubmit={handleSignUp} className="space-y-6">
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
                    value={email}
                    onChange={(e) => setEmail(e.target.value)}
                  />
                </div>

                <div className="space-y-2">
                  <Label htmlFor="firstName">First Name</Label>
                  <Input
                    id="firstName"
                    type="text"
                    placeholder="Enter your first name"
                    required
                    value={firstName}
                    onChange={(e) => setFirstName(e.target.value)}
                  />
                </div>

                <div className="space-y-2">
                  <Label htmlFor="lastName">Last Name</Label>
                  <Input
                    id="lastName"
                    type="text"
                    placeholder="Enter your last name (optional)"
                    value={lastName}
                    onChange={(e) => setLastName(e.target.value)}
                  />
                </div>

                <div className="space-y-2">
                  <Label htmlFor="password">Password</Label>
                  <Input
                    id="password"
                    type="password"
                    placeholder="Create a password at least 8 characters long"
                    required
                    value={password}
                    onChange={(e) => setPassword(e.target.value)}
                  />
                  {passwordTooShort && (
                    <div className="text-xs text-destructive mt-1">
                      Password must be at least 8 characters.
                    </div>
                  )}
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
                    value={confirmPassword}
                    onChange={(e) => setConfirmPassword(e.target.value)}
                  />
                  {passwordMismatch && (
                    <div className="text-xs text-destructive mt-1">Passwords do not match.</div>
                  )}
                </div>
              </div>

              {/* Primary Sign Up Button */}
              <Button size="lg" className="w-full" type="submit">
                Sign Up
              </Button>

              {/* Reminder Message */}
              <div className="text-xs text-muted-foreground mt-2">
                Password must be at least 8 characters. Password and confirm password must match.
              </div>

              {/* Divider */}
              <div className="relative">
                <div className="absolute inset-0 flex items-center">
                  <span className="w-full border-t border-border" />
                </div>
                <div className="relative flex justify-center text-xs uppercase">
                  <span className="bg-card px-2 text-muted-foreground">OR</span>
                </div>
              </div>

              {/* Google Sign Up Button */}
              <Button variant="outline" size="lg" className="w-full space-x-3">
                <GoogleIcon />
                <span>Sign Up with Google</span>
              </Button>

              {/* Login Link */}
              <div className="text-center">
                <div className="text-sm text-muted-foreground">
                  Already have an account? <Link to={ROUTES.LOGIN}>Login here</Link>
                </div>
              </div>
            </form>
          </CardContent>
        </Card>
      </div>
    </div>
  )
}
