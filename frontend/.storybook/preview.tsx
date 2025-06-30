import type { Preview } from '@storybook/react'
import React from 'react'
import '../src/index.css'

// Global CSS to modify Storybook root element
const globalStyles = `
  #storybook-root {
    width: 100% !important;
    max-width: none !important;
    margin: 0 !important;
    padding: 0 !important;
  }
`

const preview: Preview = {
  parameters: {
    controls: {
      matchers: {
        color: /(background|color)$/i,
        date: /Date$/i,
      },
    },
    backgrounds: {
      default: 'dark',
      values: [
        {
          name: 'dark',
          value: '#1a1f1e',
        },
        {
          name: 'light',
          value: '#ffffff',
        },
      ],
    },
  },
  decorators: [
    (Story) => (
      <>
        <style dangerouslySetInnerHTML={{ __html: globalStyles }} />
        <div
          className="dark bg-background text-foreground min-h-screen p-4"
          style={{ width: '100vw' }}
        >
          <Story />
        </div>
      </>
    ),
  ],
}

export default preview
