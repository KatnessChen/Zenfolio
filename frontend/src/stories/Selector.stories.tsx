import type { Meta, StoryObj } from '@storybook/react'
import { useState } from 'react'
import { Selector } from '@/components/ui/selector'

const meta: Meta<typeof Selector> = {
  title: 'Components/Selector',
  component: Selector,
  parameters: {
    layout: 'centered',
    docs: {
      description: {
        component:
          'A reusable selector component used throughout the dashboard for time periods, asset types, and performance views.',
      },
    },
  },
  tags: ['autodocs'],
  argTypes: {
    options: {
      description: 'Array of options with value and label properties',
    },
    value: {
      description: 'Currently selected value',
    },
    onChange: {
      description: 'Callback function called when selection changes',
      action: 'changed',
    },
    className: {
      description: 'Additional CSS classes to apply',
    },
  },
}

export default meta
type Story = StoryObj<typeof meta>

export const TimePeriodSelector: Story = {
  args: {
    options: [
      { value: '1D', label: '1D' },
      { value: '1W', label: '1W' },
      { value: '1M', label: '1M' },
      { value: '3M', label: '3M' },
      { value: '6M', label: '6M' },
      { value: 'YTD', label: 'YTD' },
      { value: '1Y', label: '1Y' },
      { value: 'ALL', label: 'ALL' },
    ],
    value: 'ALL',
    onChange: () => {},
  },
}

export const AssetAllocationSelector: Story = {
  args: {
    options: [
      { value: 'Asset Type', label: 'By Asset Type' },
      { value: 'Industry', label: 'By Industry' },
      { value: 'Broker', label: 'By Broker' },
    ],
    value: 'Asset Type',
    onChange: () => {},
  },
}

export const PerformanceSelector: Story = {
  args: {
    options: [
      { value: 'Top 5', label: 'Top 5' },
      { value: 'Worst 5', label: 'Worst 5' },
    ],
    value: 'Top 5',
    onChange: () => {},
  },
}

export const CustomSelector: Story = {
  args: {
    options: [
      { value: 'option1', label: 'Option 1' },
      { value: 'option2', label: 'Option 2' },
      { value: 'option3', label: 'Option 3' },
    ],
    value: 'option1',
    onChange: () => {},
    className: 'w-full',
  },
}

export const LongLabels: Story = {
  args: {
    options: [
      { value: 'very-long-option-1', label: 'Very Long Option Name 1' },
      { value: 'very-long-option-2', label: 'Very Long Option Name 2' },
      { value: 'very-long-option-3', label: 'Very Long Option Name 3' },
    ],
    value: 'very-long-option-1',
    onChange: () => {},
  },
}

export const SingleOption: Story = {
  args: {
    options: [{ value: 'only-option', label: 'Only Option' }],
    value: 'only-option',
    onChange: () => {},
  },
}

// Interactive example with state
const InteractiveSelector = () => {
  const [selectedValue, setSelectedValue] = useState('ALL')

  const options = [
    { value: '1D', label: '1D' },
    { value: '1W', label: '1W' },
    { value: '1M', label: '1M' },
    { value: '3M', label: '3M' },
    { value: '6M', label: '6M' },
    { value: 'YTD', label: 'YTD' },
    { value: '1Y', label: '1Y' },
    { value: 'ALL', label: 'ALL' },
  ]

  return (
    <div className="space-y-4">
      <Selector options={options} value={selectedValue} onChange={setSelectedValue} />
      <p className="text-sm text-muted-foreground">
        Selected: <span className="font-medium">{selectedValue}</span>
      </p>
    </div>
  )
}

export const Interactive: Story = {
  render: () => <InteractiveSelector />,
}

const VerticalSelector = () => {
  const [selected, setSelected] = useState('profile')
  return (
    <div className="w-56">
      <Selector
        options={[
          { value: 'profile', label: 'Profile' },
          { value: 'preference', label: 'Preference' },
          { value: 'security', label: 'Security' },
        ]}
        value={selected}
        onChange={setSelected}
        type="vertical"
      />
    </div>
  )
}

export const Vertical: Story = {
  render: () => <VerticalSelector />,
}
