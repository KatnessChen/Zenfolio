import type { Meta, StoryObj } from '@storybook/react'
import { Dropdown, DropdownItem, DropdownSeparator } from '@/components/ui/dropdown'
import DropdownTrigger from '@/components/ui/dropdown-trigger'

const meta: Meta<typeof Dropdown> = {
  title: 'Components/Dropdown',
  component: Dropdown,
  parameters: {
    layout: 'centered',
  },
  tags: ['autodocs'],
  argTypes: {
    trigger: {
      control: false,
    },
  },
}

export default meta
type Story = StoryObj<typeof meta>

export const WithTrigger: Story = {
  args: {
    trigger: <DropdownTrigger>Stock Selection</DropdownTrigger>,
    children: (
      <>
        <DropdownItem>AAPL - Apple Inc.</DropdownItem>
        <DropdownItem>GOOGL - Alphabet Inc.</DropdownItem>
        <DropdownItem>MSFT - Microsoft Corp.</DropdownItem>
        <DropdownSeparator />
        <DropdownItem disabled>TSLA - Tesla Inc. (Unavailable)</DropdownItem>
      </>
    ),
  },
}
