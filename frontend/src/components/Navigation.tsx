import { Link, useLocation, useNavigate } from 'react-router-dom'
import { useSelector, useDispatch } from 'react-redux'
import { Logo, Avatar, Dropdown, DropdownItem, DropdownSeparator } from '@/components/ui'
import { ROUTES } from '@/constants'
import { useEffect, useRef, useState, useCallback } from 'react'
import type { RootState, AppDispatch } from '@/store'
import { logoutUser } from '@/store/authSlice'

interface NavLinkProps {
  to: string
  children: React.ReactNode
  isActive: boolean
  onLinkRef: (element: HTMLAnchorElement | null) => void
}

const NavLink = ({ to, children, isActive, onLinkRef }: NavLinkProps) => {
  return (
    <Link
      ref={onLinkRef}
      to={to}
      className={`relative py-2 px-1 text-base font-medium transition-colors duration-200 ${
        isActive ? 'text-primary' : 'text-muted-foreground hover:text-foreground'
      }`}
    >
      {children}
    </Link>
  )
}

export default function Navigation() {
  const location = useLocation()
  const navigate = useNavigate()
  const dispatch = useDispatch<AppDispatch>()
  const { isAuthenticated, user } = useSelector((state: RootState) => state.auth)
  const [underlineStyle, setUnderlineStyle] = useState({ left: 0, width: 0 })
  const navLinksRef = useRef<{ [key: string]: HTMLAnchorElement | null }>({})

  const isActive = useCallback(
    (path: string) => {
      return location.pathname === path
    },
    [location.pathname]
  )

  const handleLogout = async () => {
    await dispatch(logoutUser())
  }

  const handleSettingsClick = () => {
    navigate(ROUTES.SETTINGS)
  }

  const updateUnderline = useCallback(() => {
    const activeLink = Object.entries(navLinksRef.current).find(([path]) => isActive(path))?.[1]
    if (activeLink) {
      const parent = activeLink.parentElement
      if (parent) {
        const parentRect = parent.getBoundingClientRect()
        const linkRect = activeLink.getBoundingClientRect()
        setUnderlineStyle({
          left: linkRect.left - parentRect.left,
          width: linkRect.width,
        })
      }
    }
  }, [isActive])

  useEffect(() => {
    updateUnderline()
    window.addEventListener('resize', updateUnderline)
    return () => window.removeEventListener('resize', updateUnderline)
  }, [updateUnderline])

  return (
    <nav className="bg-card/70 backdrop-blur-md border-b border-border/20 sticky top-0 z-50">
      <div className="container mx-auto px-4 py-4">
        <div className="flex justify-between items-center">
          <Logo size="lg" />

          {isAuthenticated && (
            <div className="flex items-center space-x-8">
              <div className="flex space-x-8 relative">
                <NavLink
                  to={ROUTES.DASHBOARD}
                  isActive={isActive(ROUTES.DASHBOARD)}
                  onLinkRef={(el) => (navLinksRef.current[ROUTES.DASHBOARD] = el)}
                >
                  Dashboard
                </NavLink>
                <NavLink
                  to={ROUTES.TRANSACTIONS_UPLOAD}
                  isActive={isActive(ROUTES.TRANSACTIONS_UPLOAD)}
                  onLinkRef={(el) => (navLinksRef.current[ROUTES.TRANSACTIONS_UPLOAD] = el)}
                >
                  Upload
                </NavLink>
                <NavLink
                  to={ROUTES.TRANSACTIONS}
                  isActive={isActive(ROUTES.TRANSACTIONS)}
                  onLinkRef={(el) => (navLinksRef.current[ROUTES.TRANSACTIONS] = el)}
                >
                  History
                </NavLink>

                {/* Animated underline */}
                <div
                  className="absolute bottom-1 h-0.5 bg-primary rounded-full transition-all duration-200 ease-out"
                  style={{
                    left: `${underlineStyle.left - 32}px`,
                    width: `${underlineStyle.width}px`,
                  }}
                />
              </div>

              <Dropdown
                trigger={
                  <Avatar fallback={user?.firstName?.charAt(0)?.toUpperCase() || 'U'} size="sm" />
                }
                className="ml-4"
                align="right"
                width="160px"
              >
                <DropdownItem onClick={handleSettingsClick}>Settings</DropdownItem>
                <DropdownSeparator />
                <DropdownItem onClick={handleLogout}>Logout</DropdownItem>
              </Dropdown>
            </div>
          )}
        </div>
      </div>
    </nav>
  )
}
