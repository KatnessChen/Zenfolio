import { ChevronDownIcon } from '@/components/icons'
import { cn } from '@/lib/utils'

interface DropdownTriggerProps {
  children: React.ReactNode
  className?: string
}

export default function DropdownTrigger({ children, className }: DropdownTriggerProps) {
  return (
    <div
      className={cn(
        // Same base styling as Input component
        'bg-input text-foreground placeholder:text-muted-foreground',
        // Focus state: Brighter Sage Green border and ring
        'focus-visible:border-primary focus-visible:ring-primary/50 focus-visible:ring-2',
        // Base styling with subtle borders
        'border-border/40 flex h-10 w-full min-w-0 rounded-md border px-3 py-2 text-sm',
        // Transitions and states
        'transition-[color,border-color,box-shadow] outline-none',
        // Cursor and layout
        'cursor-pointer justify-between items-center',
        // Selection styling with primary colors
        'selection:bg-primary selection:text-primary-foreground',
        // Disabled states
        'disabled:pointer-events-none disabled:cursor-not-allowed disabled:opacity-50',
        // Error states with Alert Red
        'aria-invalid:ring-destructive/30 aria-invalid:border-destructive',
        className
      )}
    >
      {children}
      <ChevronDownIcon className="ml-2 h-4 w-4 shrink-0" />
    </div>
  )
}
