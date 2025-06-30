import * as React from 'react'
import { Slot } from '@radix-ui/react-slot'
import { type VariantProps } from 'class-variance-authority'
import { Link } from 'react-router-dom'

import { cn } from '@/lib/utils'
import { buttonVariants } from './button-variants'

interface ButtonProps extends React.ComponentProps<'button'>, VariantProps<typeof buttonVariants> {
  asChild?: boolean
  to?: string
}

function Button({ className, variant, size, asChild = false, to, ...props }: ButtonProps) {
  if (to) {
    return (
      <Link
        to={to}
        data-slot="button"
        className={cn(buttonVariants({ variant, size, className }))}
        {...props}
      />
    )
  }
  const Comp = asChild ? Slot : 'button'
  return (
    <Comp
      data-slot="button"
      className={cn(buttonVariants({ variant, size, className }))}
      {...props}
    />
  )
}

export { Button }
