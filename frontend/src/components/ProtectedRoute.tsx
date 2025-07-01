import { useEffect } from 'react'
import { Navigate, useLocation, Outlet } from 'react-router-dom'
import { useSelector, useDispatch } from 'react-redux'
import type { RootState, AppDispatch } from '@/store'
import { getCurrentUser } from '@/store/authSlice'
import { ROUTES } from '@/constants'

export function ProtectedRoute() {
  const dispatch = useDispatch<AppDispatch>()
  const location = useLocation()
  const { isAuthenticated, token, user, isLoading } = useSelector((state: RootState) => state.auth)

  useEffect(() => {
    // If we have a token but no user data, fetch user info
    if (token && !user && !isLoading) {
      dispatch(getCurrentUser())
    }
  }, [token, user, isLoading, dispatch])

  // If not authenticated, redirect to login with return URL
  if (!isAuthenticated) {
    return <Navigate to={ROUTES.LOGIN} state={{ from: location }} replace />
  }

  // If authenticated but still loading user data, show loading
  if (isAuthenticated && token && !user && isLoading) {
    return (
      <div className="min-h-screen flex items-center justify-center">
        <div className="text-center">
          <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-primary mx-auto mb-4"></div>
          <p className="text-muted-foreground">Loading...</p>
        </div>
      </div>
    )
  }

  return <Outlet />
}
