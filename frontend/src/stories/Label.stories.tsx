import type { Meta, StoryObj } from '@storybook/react'
import { Label } from '@/components/ui/label'
import { Input } from '@/components/ui/input'

const meta: Meta<typeof Label> = {
  title: 'Components/Label',
  component: Label,
  parameters: {
    layout: 'centered',
  },
  tags: ['autodocs'],
}

export default meta
type Story = StoryObj<typeof meta>

export const Default: Story = {
  args: {
    children: 'Stock Symbol',
    htmlFor: 'symbol',
  },
}

export const Required: Story = {
  render: () => (
    <Label htmlFor="required-field">
      Email Address <span className="text-destructive">*</span>
    </Label>
  ),
}

export const WithInput: Story = {
  render: () => (
    <div className="space-y-2 w-64">
      <Label htmlFor="example">Transaction Amount</Label>
      <Input id="example" type="number" placeholder="0.00" />
    </div>
  ),
}

export const FormExample: Story = {
  render: () => (
    <div className="grid grid-cols-1 md:grid-cols-2 gap-6 w-96">
      <div className="space-y-2">
        <Label htmlFor="symbol">Stock Symbol</Label>
        <Input id="symbol" placeholder="e.g. AAPL" className="font-mono" />
      </div>
      <div className="space-y-2">
        <Label htmlFor="shares">Number of Shares</Label>
        <Input id="shares" type="number" placeholder="0" />
      </div>
      <div className="space-y-2">
        <Label htmlFor="price">Price per Share</Label>
        <Input id="price" type="number" step="0.01" placeholder="0.00" />
      </div>
      <div className="space-y-2">
        <Label htmlFor="date">Transaction Date</Label>
        <Input id="date" type="date" />
      </div>
    </div>
  ),
}
