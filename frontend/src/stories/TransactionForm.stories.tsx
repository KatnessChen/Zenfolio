import type { Meta, StoryObj } from '@storybook/react'
import { Button } from '@/components/ui/button'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import { Dropdown, DropdownItem, DropdownSeparator } from '@/components/ui/dropdown'
import DropdownTrigger from '@/components/ui/dropdown-trigger'

const TransactionForm = () => {
  return (
    <Card className="w-full max-w-2xl">
      <CardHeader>
        <CardTitle>Add Manual Transaction</CardTitle>
        <CardDescription>Input fields with Zenfolio styling</CardDescription>
      </CardHeader>
      <CardContent>
        <form className="space-y-6">
          <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
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
            <div className="space-y-2">
              <Label htmlFor="type">Transaction Type</Label>
              <Dropdown trigger={<DropdownTrigger>Select Type</DropdownTrigger>}>
                <DropdownItem>Buy</DropdownItem>
                <DropdownItem>Sell</DropdownItem>
                <DropdownSeparator />
                <DropdownItem>Dividends</DropdownItem>
              </Dropdown>
            </div>
            <div className="space-y-2">
              <Label htmlFor="notes">Notes (Optional)</Label>
              <Input id="notes" placeholder="Investment strategy, notes..." />
            </div>
          </div>
          <div className="flex gap-4">
            <Button variant="default">Add Transaction</Button>
            <Button variant="ghost">Cancel</Button>
            <Button variant="outline">Save as Draft</Button>
          </div>
        </form>
      </CardContent>
    </Card>
  )
}

const meta: Meta<typeof TransactionForm> = {
  title: 'Examples/Transaction Form',
  component: TransactionForm,
  parameters: {
    layout: 'centered',
  },
  tags: ['autodocs'],
}

export default meta
type Story = StoryObj<typeof meta>

export const Complete: Story = {}

export const BasicForm: Story = {
  render: () => (
    <Card className="w-full max-w-md">
      <CardHeader>
        <CardTitle>Quick Entry</CardTitle>
        <CardDescription>Simple transaction form</CardDescription>
      </CardHeader>
      <CardContent>
        <form className="space-y-4">
          <div className="space-y-2">
            <Label htmlFor="quick-symbol">Symbol</Label>
            <Input id="quick-symbol" placeholder="AAPL" className="font-mono" />
          </div>
          <div className="space-y-2">
            <Label htmlFor="quick-amount">Amount</Label>
            <Input id="quick-amount" type="number" placeholder="1000.00" />
          </div>
          <Button className="w-full">Add Transaction</Button>
        </form>
      </CardContent>
    </Card>
  ),
}

export const ValidationStates: Story = {
  render: () => (
    <Card className="w-full max-w-md">
      <CardHeader>
        <CardTitle>Form Validation</CardTitle>
        <CardDescription>Examples of input validation states</CardDescription>
      </CardHeader>
      <CardContent>
        <form className="space-y-4">
          <div className="space-y-2">
            <Label htmlFor="valid-input">Valid Input</Label>
            <Input id="valid-input" value="AAPL" className="font-mono" />
          </div>
          <div className="space-y-2">
            <Label htmlFor="error-input">Error State</Label>
            <Input
              id="error-input"
              value="INVALID"
              className="font-mono border-destructive"
              aria-invalid="true"
            />
            <p className="text-xs text-destructive">Please enter a valid stock symbol</p>
          </div>
          <div className="space-y-2">
            <Label htmlFor="required-input">
              Required Field <span className="text-destructive">*</span>
            </Label>
            <Input id="required-input" placeholder="Required..." />
          </div>
          <div className="flex gap-2">
            <Button variant="default">Submit</Button>
            <Button variant="outline">Reset</Button>
          </div>
        </form>
      </CardContent>
    </Card>
  ),
}
