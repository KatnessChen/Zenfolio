import { useState, useEffect } from 'react'
import { useNavigate, useLocation, Navigate } from 'react-router-dom'
import { useDispatch, useSelector } from 'react-redux'
import { Button, Input, Label, Card, CardContent, Title, Link } from '@/components/ui'
import { GoogleIcon } from '@/components/icons'
import { ROUTES } from '@/constants'
import type { RootState, AppDispatch } from '@/store'
import { loginUser, clearError } from '@/store/authSlice'

interface LocationState {
  from?: {
    pathname: string
  }
}

export default function LoginPage() {
  const [email, setEmail] = useState('')
  const [password, setPassword] = useState('')
  const [formErrors, setFormErrors] = useState<{ email?: string; password?: string }>({})

  const dispatch = useDispatch<AppDispatch>()
  const navigate = useNavigate()
  const location = useLocation()

  const { isAuthenticated, isLoading, error } = useSelector((state: RootState) => state.auth)

  // Get the intended destination from location state
  const from = (location.state as LocationState)?.from?.pathname || ROUTES.DASHBOARD

  useEffect(() => {
    // Clear any previous errors when component mounts
    dispatch(clearError())
  }, [dispatch])

  // Redirect if already authenticated
  if (isAuthenticated) {
    return <Navigate to={from} replace />
  }

  const validateForm = () => {
    const errors: { email?: string; password?: string } = {}

    if (!email) {
      errors.email = 'Email is required'
    } else if (!/\S+@\S+\.\S+/.test(email)) {
      errors.email = 'Please enter a valid email address'
    }

    if (!password) {
      errors.password = 'Password is required'
    }

    setFormErrors(errors)
    return Object.keys(errors).length === 0
  }

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()

    if (!validateForm()) {
      return
    }

    try {
      const result = await dispatch(loginUser({ email, password }))

      if (loginUser.fulfilled.match(result)) {
        // Login successful, navigate to intended page
        navigate(from, { replace: true })
      }
    } catch (err) {
      // Error is handled by the reducer
      console.error('Login error:', err)
    }
  }
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
            {/* Error Display */}
            {error && (
              <div className="bg-destructive/10 border border-destructive/20 text-destructive px-4 py-3 rounded-md text-sm">
                {error}
              </div>
            )}

            {/* Form Fields */}
            <form onSubmit={handleSubmit} className="space-y-4">
              <div className="space-y-2">
                <Label htmlFor="email" className="text-foreground">
                  Email Address
                </Label>
                <Input
                  id="email"
                  type="email"
                  placeholder="Enter your email"
                  value={email}
                  onChange={(e) => setEmail(e.target.value)}
                  required
                  className="bg-card text-foreground border-border placeholder:text-muted-foreground"
                />
                {formErrors.email && <p className="text-destructive text-sm">{formErrors.email}</p>}
              </div>

              <div className="space-y-2">
                <Label htmlFor="password" className="text-foreground">
                  Password
                </Label>
                <Input
                  id="password"
                  type="password"
                  placeholder="Enter your password"
                  value={password}
                  onChange={(e) => setPassword(e.target.value)}
                  required
                  className="bg-card text-foreground border-border placeholder:text-muted-foreground"
                />
                {formErrors.password && (
                  <p className="text-destructive text-sm">{formErrors.password}</p>
                )}
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
              <Button type="submit" size="lg" className="w-full" disabled={isLoading}>
                {isLoading ? 'Logging in...' : 'Login'}
              </Button>
            </form>

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
