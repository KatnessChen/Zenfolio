import * as React from 'react'
import { Button, Input, Label, Card, CardContent, Title, Link } from '@/components/ui'
import { GoogleIcon } from '@/components/icons'
import { ROUTES } from '@/constants'
import { AuthService } from '@/services/auth.service'
import { useNavigate } from 'react-router-dom'

export default function SignUpPage() {
  const [email, setEmail] = React.useState('')
  const [firstName, setFirstName] = React.useState('')
  const [lastName, setLastName] = React.useState('')
  const [password, setPassword] = React.useState('')
  const [confirmPassword, setConfirmPassword] = React.useState('')
  const [passwordTooShort, setPasswordTooShort] = React.useState(false)
  const [passwordMismatch, setPasswordMismatch] = React.useState(false)
  const [isLoading, setIsLoading] = React.useState(false)
  const [error, setError] = React.useState<string | null>(null)
  const [success, setSuccess] = React.useState(false)
  const navigate = useNavigate()
  const isFormValid =
    email.trim() && firstName.trim() && password.length >= 8 && password === confirmPassword

  async function handleSignUp(e: React.FormEvent) {
    e.preventDefault()
    setPasswordTooShort(password.length > 0 && password.length < 8)
    setPasswordMismatch(confirmPassword.length > 0 && password !== confirmPassword)
    setError(null)
    setSuccess(false)
    if (!isFormValid) return
    setIsLoading(true)
    try {
      const result = await AuthService.signUp({
        email,
        firstName,
        lastName: lastName || undefined,
        password,
        confirmPassword,
      })
      if (result.success) {
        setSuccess(true)
        setError(null)
        setTimeout(() => navigate(ROUTES.LOGIN), 1200)
      } else {
        setError(result.message || 'Sign up failed')
      }
    } catch (err: unknown) {
      const errorMessage = err instanceof Error ? err.message : 'Sign up failed'
      setError(errorMessage)
    } finally {
      setIsLoading(false)
    }
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
              <Button size="lg" className="w-full" type="submit" disabled={isLoading}>
                {isLoading ? 'Signing Up...' : 'Sign Up'}
              </Button>
              {error && <div className="text-xs text-destructive mt-2 text-center">{error}</div>}
              {success && (
                <div className="text-xs text-success mt-2 text-center">
                  Sign up successful! Please log in.
                </div>
              )}

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
