import { cva } from 'class-variance-authority'

export const buttonVariants = cva(
  // Base styles with minimal interactive effects and subtle click feedback
  "inline-flex items-center justify-center gap-2 whitespace-nowrap rounded-md text-sm font-medium transition-all duration-200 ease-out disabled:pointer-events-none disabled:opacity-50 [&_svg]:pointer-events-none [&_svg:not([class*='size-'])]:size-4 shrink-0 [&_svg]:shrink-0 outline-none focus-visible:border-ring focus-visible:ring-ring/50 focus-visible:ring-[3px] aria-invalid:ring-destructive/20 dark:aria-invalid:ring-destructive/40 aria-invalid:border-destructive transform-gpu active:scale-[0.99] active:duration-75",
  {
    variants: {
      variant: {
        default:
          'bg-primary text-primary-foreground hover:bg-primary/90 dark:bg-primary dark:text-primary-foreground dark:hover:bg-primary/85', // Minimal sage green
        destructive:
          'bg-destructive text-destructive-foreground hover:bg-destructive/90 focus-visible:ring-destructive/20 dark:focus-visible:ring-destructive/40', // Minimal alert red
        outline:
          'border border-primary text-primary bg-transparent hover:bg-primary/10 dark:border-primary dark:text-primary dark:hover:bg-primary/15', // Minimal outline
        secondary:
          'bg-secondary text-secondary-foreground hover:bg-secondary/90 dark:bg-secondary dark:text-secondary-foreground dark:hover:bg-secondary/80', // Minimal secondary
        ghost: 'text-primary hover:underline dark:text-primary', // Minimal ghost - now uses primary color with underline
        link: 'text-foreground hover:bg-accent hover:text-accent-foreground dark:text-foreground dark:hover:bg-accent/20', // Minimal link - now uses foreground color with background
        profit: 'bg-profit text-profit-foreground hover:bg-profit/90', // Minimal profit
        loss: 'bg-loss text-loss-foreground hover:bg-loss/90', // Minimal loss
      },
      size: {
        default: 'h-10 px-5 py-2 has-[>svg]:px-4', // Slightly larger for better touch targets
        sm: 'h-8 rounded-md gap-1.5 px-3 has-[>svg]:px-2.5',
        lg: 'h-12 rounded-md px-6 has-[>svg]:px-5',
        icon: 'size-10', // Slightly larger for better touch targets
      },
    },
    defaultVariants: {
      variant: 'default',
      size: 'default',
    },
  }
)
