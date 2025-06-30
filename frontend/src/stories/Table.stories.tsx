import type { Meta, StoryObj } from '@storybook/react'
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from '@/components/ui/table'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'
import { EditIcon, DeleteIcon } from '@/components/icons'
import { useState } from 'react'

const meta: Meta<typeof Table> = {
  title: 'Components/Table',
  component: Table,
  parameters: {
    layout: 'padded',
  },
  tags: ['autodocs'],
}

export default meta
type Story = StoryObj<typeof meta>

export const Basic: Story = {
  render: () => (
    <Table>
      <TableHeader>
        <TableRow>
          <TableHead>Symbol</TableHead>
          <TableHead>Type</TableHead>
          <TableHead>Price</TableHead>
          <TableHead className="text-right">Amount</TableHead>
        </TableRow>
      </TableHeader>
      <TableBody>
        <TableRow>
          <TableCell className="font-mono font-medium">AAPL</TableCell>
          <TableCell>Buy</TableCell>
          <TableCell>$150.25</TableCell>
          <TableCell className="text-right">$15,025.00</TableCell>
        </TableRow>
        <TableRow>
          <TableCell className="font-mono font-medium">GOOGL</TableCell>
          <TableCell>Sell</TableCell>
          <TableCell>$2,750.80</TableCell>
          <TableCell className="text-right">$137,540.00</TableCell>
        </TableRow>
      </TableBody>
    </Table>
  ),
}

export const WithActions: Story = {
  render: () => (
    <Table>
      <TableHeader>
        <TableRow>
          <TableHead>Symbol</TableHead>
          <TableHead>Type</TableHead>
          <TableHead>Price</TableHead>
          <TableHead className="text-right">Amount</TableHead>
          <TableHead>Notes</TableHead>
          <TableHead>Actions</TableHead>
        </TableRow>
      </TableHeader>
      <TableBody>
        <TableRow>
          <TableCell className="font-mono font-medium">AAPL</TableCell>
          <TableCell>Buy</TableCell>
          <TableCell>$150.25</TableCell>
          <TableCell className="text-right text-muted-foreground">$15,025.00</TableCell>
          <TableCell>Long-term investment</TableCell>
          <TableCell>
            <div className="flex items-center gap-2">
              <button className="p-1 hover:bg-muted rounded transition-colors">
                <EditIcon size={16} className="text-muted-foreground hover:text-primary" />
              </button>
              <button className="p-1 hover:bg-muted rounded transition-colors">
                <DeleteIcon size={16} className="text-muted-foreground hover:text-destructive" />
              </button>
            </div>
          </TableCell>
        </TableRow>
        <TableRow>
          <TableCell className="font-mono font-medium">MSFT</TableCell>
          <TableCell>Dividend</TableCell>
          <TableCell>$2.75</TableCell>
          <TableCell className="text-right text-muted-foreground">$550.00</TableCell>
          <TableCell>Quarterly dividend</TableCell>
          <TableCell>
            <div className="flex items-center gap-2">
              <button className="p-1 hover:bg-muted rounded transition-colors">
                <EditIcon size={16} className="text-muted-foreground hover:text-primary" />
              </button>
              <button className="p-1 hover:bg-muted rounded transition-colors">
                <DeleteIcon size={16} className="text-muted-foreground hover:text-destructive" />
              </button>
            </div>
          </TableCell>
        </TableRow>
      </TableBody>
    </Table>
  ),
}

export const Sortable: Story = {
  render: () => <SortableTableExample />,
}

const SortableTableExample = () => {
  const [sortDirection, setSortDirection] = useState<'asc' | 'desc' | null>('desc')

  return (
    <Table>
      <TableHeader>
        <TableRow>
          <TableHead
            sortable
            sortDirection={sortDirection}
            onSort={() => setSortDirection(sortDirection === 'asc' ? 'desc' : 'asc')}
          >
            Symbol
          </TableHead>
          <TableHead sortable>Type</TableHead>
          <TableHead sortable>Price</TableHead>
          <TableHead sortable className="text-right">
            Amount
          </TableHead>
          <TableHead>Notes</TableHead>
        </TableRow>
      </TableHeader>
      <TableBody>
        <TableRow>
          <TableCell className="font-mono font-medium">AAPL</TableCell>
          <TableCell>Buy</TableCell>
          <TableCell>$150.25</TableCell>
          <TableCell className="text-right text-muted-foreground">$15,025.00</TableCell>
          <TableCell>Long-term investment</TableCell>
        </TableRow>
        <TableRow>
          <TableCell className="font-mono font-medium">GOOGL</TableCell>
          <TableCell>Sell</TableCell>
          <TableCell>$2,750.80</TableCell>
          <TableCell className="text-right text-muted-foreground">$137,540.00</TableCell>
          <TableCell>Profit taking</TableCell>
        </TableRow>
      </TableBody>
    </Table>
  )
}

export const InCard: Story = {
  render: () => (
    <Card className="w-full max-w-4xl">
      <CardHeader>
        <CardTitle>Transaction Data Table</CardTitle>
        <CardDescription>Example table with sortable headers and financial styling</CardDescription>
      </CardHeader>
      <CardContent className="p-0">
        <Table>
          <TableHeader>
            <TableRow>
              <TableHead>Symbol</TableHead>
              <TableHead>Type</TableHead>
              <TableHead>Price</TableHead>
              <TableHead className="text-right">Amount</TableHead>
              <TableHead>Notes</TableHead>
              <TableHead>Actions</TableHead>
            </TableRow>
          </TableHeader>
          <TableBody>
            <TableRow>
              <TableCell className="font-mono font-medium">AAPL</TableCell>
              <TableCell>Buy</TableCell>
              <TableCell>$150.25</TableCell>
              <TableCell className="text-right text-muted-foreground">$15,025.00</TableCell>
              <TableCell>Long-term investment</TableCell>
              <TableCell>
                <div className="flex items-center gap-2">
                  <button className="p-1 hover:bg-muted rounded transition-colors">
                    <EditIcon size={16} className="text-muted-foreground hover:text-primary" />
                  </button>
                  <button className="p-1 hover:bg-muted rounded transition-colors">
                    <DeleteIcon
                      size={16}
                      className="text-muted-foreground hover:text-destructive"
                    />
                  </button>
                </div>
              </TableCell>
            </TableRow>
            <TableRow>
              <TableCell className="font-mono font-medium">GOOGL</TableCell>
              <TableCell>Sell</TableCell>
              <TableCell>$2,750.80</TableCell>
              <TableCell className="text-right text-muted-foreground">$137,540.00</TableCell>
              <TableCell>Profit taking</TableCell>
              <TableCell>
                <div className="flex items-center gap-2">
                  <button className="p-1 hover:bg-muted rounded transition-colors">
                    <EditIcon size={16} className="text-muted-foreground hover:text-primary" />
                  </button>
                  <button className="p-1 hover:bg-muted rounded transition-colors">
                    <DeleteIcon
                      size={16}
                      className="text-muted-foreground hover:text-destructive"
                    />
                  </button>
                </div>
              </TableCell>
            </TableRow>
          </TableBody>
        </Table>
      </CardContent>
    </Card>
  ),
}
