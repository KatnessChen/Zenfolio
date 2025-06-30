import * as React from 'react'
import { Link } from 'react-router-dom'
import { cn } from '@/lib/utils'

interface LogoProps extends Omit<React.ComponentProps<typeof Link>, 'to'> {
  size?: 'sm' | 'default' | 'lg'
  to?: string
}

function Logo({ className, size = 'default', to = '/', ...props }: LogoProps) {
  const sizeClasses = {
    sm: 'text-lg',
    default: 'text-xl',
    lg: 'text-2xl',
  }

  const baseClasses =
    'font-bold text-primary transition-colors hover:text-primary/90 inline-flex items-center'

  return (
    <Link to={to} className={cn(baseClasses, sizeClasses[size], className)} {...props}>
      <span className="text-primary">Zen</span>
      <span className="text-muted-foreground">folio</span>
    </Link>
  )
}

export { Logo }
