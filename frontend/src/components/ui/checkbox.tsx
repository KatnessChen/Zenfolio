import React, { forwardRef } from 'react'
import { cn } from '@/lib/utils'

export interface CheckboxProps extends Omit<React.InputHTMLAttributes<HTMLInputElement>, 'type'> {
  label?: string
  description?: string
  indeterminate?: boolean
}

const Checkbox = forwardRef<HTMLInputElement, CheckboxProps>(
  ({ className, label, description, indeterminate, ...props }, ref) => {
    const inputRef = React.useRef<HTMLInputElement>(null)

    React.useImperativeHandle(ref, () => inputRef.current as HTMLInputElement)

    React.useEffect(() => {
      if (inputRef.current) {
        inputRef.current.indeterminate = indeterminate || false
      }
    }, [indeterminate])

    return (
      <div className="flex items-center space-x-2">
        <input
          ref={inputRef}
          type="checkbox"
          className={cn(
            'peer h-4 w-4 shrink-0 rounded-sm border border-primary ring-offset-background',
            'focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2',
            'disabled:cursor-not-allowed disabled:opacity-50',
            'data-[state=checked]:bg-primary data-[state=checked]:text-primary-foreground',
            'accent-primary',
            className
          )}
          {...props}
        />
        {label && (
          <div className="grid gap-1.5 leading-none">
            <label
              htmlFor={props.id}
              className="text-sm font-medium leading-none peer-disabled:cursor-not-allowed peer-disabled:opacity-70"
            >
              {label}
            </label>
            {description && <p className="text-xs text-muted-foreground">{description}</p>}
          </div>
        )}
      </div>
    )
  }
)

Checkbox.displayName = 'Checkbox'

export { Checkbox }
