import * as React from 'react'
import { cn } from '@/lib/utils'

interface AvatarProps {
  children?: React.ReactNode
  fallback?: string
  className?: string
  size?: 'sm' | 'md' | 'lg'
}

const Avatar = React.forwardRef<HTMLDivElement, AvatarProps>(
  ({ children, fallback, className, size = 'md' }, ref) => {
    const sizeClasses = {
      sm: 'h-8 w-8 text-xs',
      md: 'h-10 w-10 text-sm',
      lg: 'h-12 w-12 text-base',
    }

    return (
      <div
        ref={ref}
        className={cn(
          'relative flex shrink-0 overflow-hidden rounded-full',
          'bg-primary/10 border border-border/20',
          'flex items-center justify-center',
          'text-primary font-medium',
          'hover:bg-primary/20 transition-colors cursor-pointer',
          sizeClasses[size],
          className
        )}
      >
        {children || fallback}
      </div>
    )
  }
)
Avatar.displayName = 'Avatar'

export { Avatar }
