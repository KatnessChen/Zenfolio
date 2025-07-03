import type { Meta, StoryObj } from '@storybook/react'
import {
  EditIcon,
  DeleteIcon,
  ChevronDownIcon,
  PlusIcon,
  GoogleIcon,
  ClockIcon,
  SpinnerIcon,
  CheckIcon,
  XIcon,
} from '@/components/icons'

const IconShowcase = () => {
  return (
    <div className="p-6 space-y-8">
      <div>
        <h2 className="text-2xl font-medium mb-4 text-primary">Icons</h2>
        <p className="text-muted-foreground mb-6">
          Icon components used throughout the Zenfolio Transaction Tracker application.
        </p>
      </div>
      <div className="grid grid-cols-2 md:grid-cols-4 lg:grid-cols-6 gap-6">
        <div className="flex flex-col items-center space-y-2">
          <div className="p-3 bg-card border border-border rounded-md">
            <EditIcon size={24} />
          </div>
          <span className="text-xs text-muted-foreground">EditIcon</span>
        </div>
        <div className="flex flex-col items-center space-y-2">
          <div className="p-3 bg-card border border-border rounded-md">
            <DeleteIcon size={24} />
          </div>
          <span className="text-xs text-muted-foreground">DeleteIcon</span>
        </div>
        <div className="flex flex-col items-center space-y-2">
          <div className="p-3 bg-card border border-border rounded-md">
            <PlusIcon size={24} />
          </div>
          <span className="text-xs text-muted-foreground">PlusIcon</span>
        </div>
        <div className="flex flex-col items-center space-y-2">
          <div className="p-3 bg-card border border-border rounded-md">
            <ChevronDownIcon size={24} />
          </div>
          <span className="text-xs text-muted-foreground">ChevronDownIcon</span>
        </div>
        <div className="flex flex-col items-center space-y-2">
          <div className="p-3 bg-card border border-border rounded-md">
            <ClockIcon size={24} />
          </div>
          <span className="text-xs text-muted-foreground">ClockIcon</span>
        </div>
        <div className="flex flex-col items-center space-y-2">
          <div className="p-3 bg-card border border-border rounded-md">
            <SpinnerIcon size={24} className="animate-spin" />
          </div>
          <span className="text-xs text-muted-foreground">SpinnerIcon</span>
        </div>
        <div className="flex flex-col items-center space-y-2">
          <div className="p-3 bg-card border border-border rounded-md">
            <CheckIcon size={24} />
          </div>
          <span className="text-xs text-muted-foreground">CheckIcon</span>
        </div>
        <div className="flex flex-col items-center space-y-2">
          <div className="p-3 bg-card border border-border rounded-md">
            <XIcon size={24} />
          </div>
          <span className="text-xs text-muted-foreground">XIcon</span>
        </div>
      </div>

      {/* Brand icons section */}
      <section className="space-y-4">
        <h3 className="text-xl font-semibold text-primary border-b border-border pb-2">
          Brand icons
        </h3>
        <div className="grid grid-cols-2 md:grid-cols-4 lg:grid-cols-6 gap-6">
          <div className="flex flex-col items-center space-y-2">
            <div className="p-3 bg-card border border-border rounded-md">
              <GoogleIcon size={24} />
            </div>
            <span className="text-xs text-muted-foreground">GoogleIcon</span>
          </div>
        </div>
      </section>

      <section className="space-y-4">
        <h3 className="text-xl font-semibold text-primary border-b border-border pb-2">
          Size Variations
        </h3>
        <div className="grid grid-cols-4 gap-6">
          <div className="flex flex-col items-center space-y-2">
            <EditIcon size={12} className="text-muted-foreground" />
            <span className="text-xs text-muted-foreground">12px</span>
          </div>
          <div className="flex flex-col items-center space-y-2">
            <EditIcon size={16} className="text-muted-foreground" />
            <span className="text-xs text-muted-foreground">16px (default)</span>
          </div>
          <div className="flex flex-col items-center space-y-2">
            <EditIcon size={20} className="text-muted-foreground" />
            <span className="text-xs text-muted-foreground">20px</span>
          </div>
          <div className="flex flex-col items-center space-y-2">
            <EditIcon size={24} className="text-muted-foreground" />
            <span className="text-xs text-muted-foreground">24px</span>
          </div>
        </div>
      </section>

      <section className="space-y-4">
        <h3 className="text-xl font-semibold text-primary border-b border-border pb-2">
          Interactive States
        </h3>
        <div className="grid grid-cols-2 md:grid-cols-3 gap-6">
          <div className="flex items-center gap-2 p-3 bg-card border border-border rounded-md">
            <EditIcon
              size={16}
              className="text-muted-foreground hover:text-primary transition-colors cursor-pointer"
            />
            <span className="text-sm">Hover for primary</span>
          </div>
          <div className="flex items-center gap-2 p-3 bg-card border border-border rounded-md">
            <DeleteIcon
              size={16}
              className="text-muted-foreground hover:text-destructive transition-colors cursor-pointer"
            />
            <span className="text-sm">Hover for destructive</span>
          </div>
          <div className="flex items-center gap-2 p-3 bg-muted rounded-md">
            <PlusIcon size={16} className="text-muted-foreground opacity-50 cursor-not-allowed" />
            <span className="text-sm opacity-50">Disabled state</span>
          </div>
        </div>
      </section>
    </div>
  )
}

const meta: Meta<typeof IconShowcase> = {
  title: 'Components/Icons',
  component: IconShowcase,
  parameters: {
    layout: 'fullscreen',
  },
  tags: ['autodocs'],
}

export default meta
type Story = StoryObj<typeof meta>

export const AllIcons: Story = {}
