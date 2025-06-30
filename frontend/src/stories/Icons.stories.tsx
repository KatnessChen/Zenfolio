import type { Meta, StoryObj } from '@storybook/react'
import { EditIcon, DeleteIcon, ChevronDownIcon, PlusIcon, GoogleIcon } from '@/components/icons'

const IconShowcase = () => {
  return (
    <div className="p-6 space-y-8">
      <div>
        <h2 className="text-2xl font-medium mb-4 text-primary">Icons</h2>
        <p className="text-muted-foreground mb-6">
          Icon components used throughout the Zenfolio Transaction Tracker application.
        </p>
      </div>

      <section className="space-y-4">
        <h3 className="text-xl font-semibold text-primary border-b border-border pb-2">
          Action Icons
        </h3>
        <div className="grid grid-cols-2 md:grid-cols-4 lg:grid-cols-6 gap-6">
          <div className="flex flex-col items-center space-y-2">
            <div className="p-3 bg-card border border-border rounded-md">
              <EditIcon
                size={24}
                className="text-muted-foreground hover:text-primary transition-colors"
              />
            </div>
            <span className="text-xs text-muted-foreground">EditIcon</span>
          </div>
          <div className="flex flex-col items-center space-y-2">
            <div className="p-3 bg-card border border-border rounded-md">
              <DeleteIcon
                size={24}
                className="text-muted-foreground hover:text-destructive transition-colors"
              />
            </div>
            <span className="text-xs text-muted-foreground">DeleteIcon</span>
          </div>
          <div className="flex flex-col items-center space-y-2">
            <div className="p-3 bg-card border border-border rounded-md">
              <PlusIcon
                size={24}
                className="text-muted-foreground hover:text-primary transition-colors"
              />
            </div>
            <span className="text-xs text-muted-foreground">PlusIcon</span>
          </div>
        </div>
      </section>

      <section className="space-y-4">
        <h3 className="text-xl font-semibold text-primary border-b border-border pb-2">UI Icons</h3>
        <div className="grid grid-cols-2 md:grid-cols-4 lg:grid-cols-6 gap-6">
          <div className="flex flex-col items-center space-y-2">
            <div className="p-3 bg-card border border-border rounded-md">
              <ChevronDownIcon className="text-muted-foreground w-6 h-6" />
            </div>
            <span className="text-xs text-muted-foreground">ChevronDownIcon</span>
          </div>
          <div className="flex flex-col items-center space-y-2">
            <div className="p-3 bg-card border border-border rounded-md">
              <GoogleIcon />
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

export const EditIcon_: Story = {
  render: () => (
    <div className="flex gap-4 items-center">
      <EditIcon size={16} />
      <EditIcon size={20} />
      <EditIcon size={24} />
    </div>
  ),
}

export const DeleteIcon_: Story = {
  render: () => (
    <div className="flex gap-4 items-center">
      <DeleteIcon size={16} />
      <DeleteIcon size={20} />
      <DeleteIcon size={24} />
    </div>
  ),
}

export const InteractiveButtons: Story = {
  render: () => (
    <div className="flex gap-4">
      <button className="p-2 hover:bg-muted rounded transition-colors">
        <EditIcon size={16} className="text-muted-foreground hover:text-primary" />
      </button>
      <button className="p-2 hover:bg-muted rounded transition-colors">
        <DeleteIcon size={16} className="text-muted-foreground hover:text-destructive" />
      </button>
      <button className="p-2 hover:bg-muted rounded transition-colors">
        <PlusIcon size={16} className="text-muted-foreground hover:text-primary" />
      </button>
    </div>
  ),
}
