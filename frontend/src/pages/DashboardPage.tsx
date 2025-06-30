import React from 'react'
import {
  PortfolioSummaryCard,
  AssetAllocationCard,
  TopPerformersCard,
  PositionCard,
} from '@/components/dashboard'
import { Title } from '@/components/ui/title'

const DashboardPage: React.FC = () => {
  return (
    <div className="min-h-screen bg-background text-foreground">
      <div className="container mx-auto px-4 py-6">
        {/* Page Header */}
        <div className="flex justify-between items-start mb-8">
          <Title as="h1">Portfolio Overview</Title>
        </div>

        {/* Dashboard Grid */}
        <div className="space-y-6">
          {/* Portfolio Summary - Full Width */}
          <PortfolioSummaryCard />

          {/* Positions - Full Width */}
          <PositionCard />

          {/* Main Content Area */}
          <div className="grid grid-cols-1 lg:grid-cols-3 gap-6">
            {/* Main Area - Asset Allocation */}
            <div className="lg:col-span-2">
              <AssetAllocationCard />
            </div>

            {/* Aside Area - Top Performers & Recent Activity */}
            <div className="space-y-6">
              <TopPerformersCard />
            </div>
          </div>
        </div>
      </div>
    </div>
  )
}

export default DashboardPage
