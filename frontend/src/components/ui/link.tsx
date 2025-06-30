import * as React from 'react'
import { Link as RouterLink, type LinkProps } from 'react-router-dom'
import { cn } from '@/lib/utils'

export interface LinkPropsCustom extends LinkProps {
  className?: string
}

export const Link = React.forwardRef<HTMLAnchorElement, LinkPropsCustom>(
  ({ className, ...props }, ref) => (
    <RouterLink
      ref={ref}
      className={cn('text-primary hover:text-primary/90 transition-colors font-medium', className)}
      {...props}
    />
  )
)
Link.displayName = 'Link'
