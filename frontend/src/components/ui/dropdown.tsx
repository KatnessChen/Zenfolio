import * as React from 'react'
import { cn } from '@/lib/utils'

interface DropdownProps {
  trigger: React.ReactNode
  children: React.ReactNode
  className?: string
}

interface DropdownItemProps {
  children: React.ReactNode
  onClick?: () => void
  className?: string
  disabled?: boolean
}

const Dropdown = React.forwardRef<HTMLDivElement, DropdownProps>(
  ({ trigger, children, className }) => {
    const [isOpen, setIsOpen] = React.useState(false)
    const dropdownRef = React.useRef<HTMLDivElement>(null)

    React.useEffect(() => {
      const handleClickOutside = (event: MouseEvent) => {
        if (dropdownRef.current && !dropdownRef.current.contains(event.target as Node)) {
          setIsOpen(false)
        }
      }

      document.addEventListener('mousedown', handleClickOutside)
      return () => document.removeEventListener('mousedown', handleClickOutside)
    }, [])

    return (
      <div ref={dropdownRef} className={cn('relative inline-block', className)}>
        <div onClick={() => setIsOpen(!isOpen)} className="cursor-pointer">
          {trigger}
        </div>

        {isOpen && (
          <div
            className={cn(
              'absolute top-full left-0 min-w-[200px] z-50',
              'bg-card border border-border/50 rounded-md shadow-lg',
              'backdrop-blur-sm bg-card/95',
              'animate-in fade-in-0 zoom-in-95 duration-200'
            )}
          >
            <div className="p-1">{children}</div>
          </div>
        )}
      </div>
    )
  }
)
Dropdown.displayName = 'Dropdown'

const DropdownItem = React.forwardRef<HTMLDivElement, DropdownItemProps>(
  ({ children, onClick, className, disabled = false }, ref) => {
    return (
      <div
        ref={ref}
        onClick={disabled ? undefined : onClick}
        className={cn(
          'px-3 py-2 text-sm rounded-sm cursor-pointer',
          'text-foreground',
          // Hover effect
          'hover:bg-primary/10 hover:text-primary',
          // Focus effects
          'focus:bg-primary/10 focus:text-primary focus:outline-none',
          // Disabled state
          disabled && 'opacity-50 cursor-not-allowed hover:bg-transparent hover:text-foreground',
          className
        )}
      >
        {children}
      </div>
    )
  }
)
DropdownItem.displayName = 'DropdownItem'

const DropdownSeparator = React.forwardRef<HTMLDivElement, React.HTMLAttributes<HTMLDivElement>>(
  ({ className, ...props }, ref) => (
    <div ref={ref} className={cn('h-px bg-border/50 my-1', className)} {...props} />
  )
)
DropdownSeparator.displayName = 'DropdownSeparator'

export { Dropdown, DropdownItem, DropdownSeparator }
