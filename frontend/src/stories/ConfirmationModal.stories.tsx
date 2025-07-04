import type { Meta, StoryObj } from '@storybook/react'
import { fn } from '@storybook/test'
import { ConfirmationModal } from '@/components/ui/confirmation-modal'

const meta: Meta<typeof ConfirmationModal> = {
  title: 'UI/ConfirmationModal',
  component: ConfirmationModal,
  parameters: {
    layout: 'centered',
    docs: {
      description: {
        component:
          'A reusable confirmation modal for user confirmations with customizable text and actions.',
      },
    },
  },
  tags: ['autodocs'],
  argTypes: {
    isOpen: {
      control: 'boolean',
      description: 'Whether the modal is visible',
    },
    title: {
      control: 'text',
      description: 'Modal title',
    },
    message: {
      control: 'text',
      description: 'Modal description/message',
    },
    confirmText: {
      control: 'text',
      description: 'Text for the confirm button',
    },
    cancelText: {
      control: 'text',
      description: 'Text for the cancel button',
    },
    confirmVariant: {
      control: 'select',
      options: ['default'],
      description: 'Variant for the confirm button',
    },
    isLoading: {
      control: 'boolean',
      description: 'Whether the confirm action is loading',
    },
  },
}

export default meta
type Story = StoryObj<typeof meta>

export const Default: Story = {
  args: {
    isOpen: true,
    title: 'Confirm Action',
    message: 'Are you sure you want to perform this action?',
    confirmText: 'Confirm',
    cancelText: 'Cancel',
    confirmVariant: 'default',
    isLoading: false,
    onConfirm: fn(),
    onCancel: fn(),
  },
}
