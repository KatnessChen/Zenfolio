import type { Meta, StoryObj } from '@storybook/react'
import { Logo } from '@/components/ui/logo'
import { MemoryRouter } from 'react-router-dom'

const meta: Meta<typeof Logo> = {
  title: 'Components/Logo',
  component: Logo,
  parameters: {
    layout: 'centered',
  },
  tags: ['autodocs'],
  argTypes: {
    size: {
      control: 'select',
      options: ['sm', 'default', 'lg'],
    },
  },
  decorators: [
    (Story) => (
      <MemoryRouter>
        <Story />
      </MemoryRouter>
    ),
  ],
}

export default meta
type Story = StoryObj<typeof meta>

export const AllSizes: Story = {
  render: () => (
    <div className="flex flex-wrap gap-6 items-center">
      <div className="text-center">
        <Logo size="sm" />
      </div>
      <div className="text-center">
        <Logo size="default" />
      </div>
      <div className="text-center">
        <Logo size="lg" />
      </div>
    </div>
  ),
}
