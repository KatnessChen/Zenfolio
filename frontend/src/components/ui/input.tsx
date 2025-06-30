import * as React from 'react'
import { cn } from '@/lib/utils'
import { EyeIcon, EyeOffIcon } from '@/components/icons/EyeIcon'

const Input = React.forwardRef<HTMLInputElement, React.ComponentProps<'input'>>(
  ({ className, type, ...props }, ref) => {
    const [show, setShow] = React.useState(false)
    const isPassword = type === 'password'
    return (
      <div className={cn('relative', className)}>
        <input
          ref={ref}
          type={isPassword ? (show ? 'text' : 'password') : type}
          data-slot="input"
          className={cn(
            // Zenfolio Design: Dark Grey-Green background (#3A4B40) with Light Off-White text
            'bg-input text-foreground placeholder:text-muted-foreground',
            // Focus state: Brighter Sage Green border and ring
            'focus-visible:border-primary focus-visible:ring-primary/50 focus-visible:ring-2',
            // Base styling with subtle borders
            'border-border/40 flex h-10 w-full min-w-0 rounded-md border px-3 py-2 text-sm',
            // Transitions and states
            'transition-[color,border-color,box-shadow] outline-none',
            // File input styling
            'file:text-foreground file:inline-flex file:h-7 file:border-0 file:bg-transparent file:text-sm file:font-medium',
            // Selection styling with primary colors
            'selection:bg-primary selection:text-primary-foreground',
            // Disabled states
            'disabled:pointer-events-none disabled:cursor-not-allowed disabled:opacity-50',
            // Error states with Alert Red
            'aria-invalid:ring-destructive/30 aria-invalid:border-destructive',
            isPassword ? 'pr-10' : ''
          )}
          {...props}
        />
        {isPassword && (
          <button
            type="button"
            tabIndex={-1}
            aria-label={show ? 'Hide password' : 'Show password'}
            className="absolute right-2 top-1/2 -translate-y-1/2 p-1 text-muted-foreground hover:text-foreground focus:outline-none"
            onClick={() => setShow((v) => !v)}
          >
            {show ? <EyeOffIcon width={20} height={20} /> : <EyeIcon width={20} height={20} />}
          </button>
        )}
      </div>
    )
  }
)
Input.displayName = 'Input'

export { Input }
