import type { Meta, StoryObj } from '@storybook/react'
import { Link } from '@/components/ui/link'

const meta = {
  title: 'Components/Link',
  component: Link,
  parameters: {
    layout: 'centered',
    docs: {
      description: {
        component: 'A styled link component for navigation, matching Zenfolio design system.',
      },
    },
  },
  tags: ['autodocs'],
  argTypes: {
    to: { control: 'text', description: 'Destination route' },
    children: { control: 'text', description: 'Link text' },
  },
} satisfies Meta<typeof Link>

export default meta

type Story = StoryObj<typeof meta>

export const Default: Story = {
  args: {
    to: '/login',
    children: 'Go to Login',
  },
}

export const CustomClass: Story = {
  args: {
    to: '/dashboard',
    children: 'Dashboard',
    className: 'underline text-lg',
  },
}
