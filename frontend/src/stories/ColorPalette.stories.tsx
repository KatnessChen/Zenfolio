import type { Meta, StoryObj } from '@storybook/react'

// This is a documentation-only story showing the Zenfolio color palette
const ColorPalette = () => {
  return (
    <div className="p-6 space-y-8 max-w-4xl">
      <div>
        <h2 className="text-2xl font-medium mb-4 text-primary">Zenfolio Color Palette</h2>
        <p className="text-muted-foreground mb-6">
          Demonstration of the Zenfolio Design System color palette with "Yoga Version" dark mode
          aesthetics.
        </p>
      </div>

      <section className="space-y-4">
        <h3 className="text-xl font-semibold text-primary border-b border-border pb-2">
          Core Colors
        </h3>
        <div className="grid grid-cols-2 md:grid-cols-4 gap-4">
          <div className="space-y-2">
            <div className="w-full h-16 bg-background border border-border rounded-md"></div>
            <p className="text-xs text-muted-foreground">
              Deep Ink Green
              <br />
              (Background)
            </p>
          </div>
          <div className="space-y-2">
            <div className="w-full h-16 bg-card border border-border rounded-md"></div>
            <p className="text-xs text-muted-foreground">
              Dark Grey-Green
              <br />
              (Cards)
            </p>
          </div>
          <div className="space-y-2">
            <div className="w-full h-16 bg-primary border border-border rounded-md"></div>
            <p className="text-xs text-muted-foreground">
              Brighter Sage Green
              <br />
              (Primary)
            </p>
          </div>
          <div className="space-y-2">
            <div className="w-full h-16 bg-muted-foreground border border-border rounded-md"></div>
            <p className="text-xs text-muted-foreground">
              Medium Grey-Green
              <br />
              (Secondary Text)
            </p>
          </div>
        </div>
      </section>

      <section className="space-y-4">
        <h3 className="text-xl font-semibold text-primary border-b border-border pb-2">
          Financial & Status Colors
        </h3>
        <div className="grid grid-cols-2 md:grid-cols-4 gap-4">
          <div className="space-y-2">
            <div className="w-full h-16 bg-profit border border-border rounded-md"></div>
            <p className="text-xs text-muted-foreground">
              Bright Grass Green
              <br />
              (Profit)
            </p>
          </div>
          <div className="space-y-2">
            <div className="w-full h-16 bg-loss border border-border rounded-md"></div>
            <p className="text-xs text-muted-foreground">
              Alert Red
              <br />
              (Loss/Error)
            </p>
          </div>
        </div>
      </section>

      <section className="space-y-4">
        <h3 className="text-xl font-semibold text-primary border-b border-border pb-2">
          Usage Examples
        </h3>
        <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
          <div className="space-y-4">
            <h4 className="text-lg font-medium text-primary">Text Colors</h4>
            <div className="space-y-2">
              <p className="text-foreground">Primary text (foreground)</p>
              <p className="text-muted-foreground">Secondary text (muted-foreground)</p>
              <p className="text-primary">Primary accent text</p>
              <p className="text-profit">Profit/success text</p>
              <p className="text-destructive">Error/destructive text</p>
            </div>
          </div>
          <div className="space-y-4">
            <h4 className="text-lg font-medium text-primary">Background Usage</h4>
            <div className="space-y-2">
              <div className="p-3 bg-card border border-border rounded-md">
                <p className="text-sm">Card background</p>
              </div>
              <div className="p-3 bg-muted rounded-md">
                <p className="text-sm">Muted background</p>
              </div>
              <div className="p-3 bg-primary text-primary-foreground rounded-md">
                <p className="text-sm">Primary background</p>
              </div>
            </div>
          </div>
        </div>
      </section>
    </div>
  )
}

const meta: Meta<typeof ColorPalette> = {
  title: 'Design System/Color Palette',
  component: ColorPalette,
  parameters: {
    layout: 'fullscreen',
  },
  tags: ['autodocs'],
}

export default meta
type Story = StoryObj<typeof meta>

export const ZenfolioColors: Story = {}

const chartColors = [
  { var: '--chart-1', name: 'Pale Wheat' },
  { var: '--chart-2', name: 'Deep Soil' },
  { var: '--chart-3', name: 'Soft Salmon Pink' },
  { var: '--chart-4', name: 'Earthy Pink' },
  { var: '--chart-5', name: 'Terracotta (Earthy Red)' },
  { var: '--chart-6', name: 'Wheat' },
  { var: '--chart-7', name: 'Brick Brown' },
  { var: '--chart-8', name: 'Sandstone' },
  { var: '--chart-9', name: 'Olive Green' },
  { var: '--chart-10', name: 'Dusty Blue' },
  { var: '--chart-11', name: 'Forest Green' },
  { var: '--chart-12', name: 'Clay Orange' },
  { var: '--chart-13', name: 'Moss Green' },
  { var: '--chart-14', name: 'Deep Navy' },
  { var: '--chart-15', name: 'Ochre Yellow' },
  { var: '--chart-16', name: 'Deep Night Sky Blue' },
  { var: '--chart-17', name: 'Sage Green' },
  { var: '--chart-18', name: 'Olive' },
  { var: '--chart-19', name: 'Slate Green' },
  { var: '--chart-20', name: 'Teal' },
]

export const ChartColors = () => (
  <div style={{ display: 'flex', flexWrap: 'wrap', gap: 24 }}>
    {chartColors.map((color) => (
      <div key={color.var} style={{ width: 120, textAlign: 'center' }}>
        <div
          style={{
            width: 64,
            height: 64,
            borderRadius: 8,
            margin: '0 auto 8px',
            background: `hsl(var(${color.var}))`,
            border: '1px solid #ccc',
          }}
        />
        <div style={{ fontWeight: 600 }}>{color.name}</div>
        <div style={{ fontSize: 12, color: '#888' }}>{color.var}</div>
      </div>
    ))}
  </div>
)
