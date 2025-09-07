import * as RadixTooltip from '@radix-ui/react-tooltip'
import React from 'react'
import type { ReactNode } from 'react'

interface TooltipProps {
  content: ReactNode
  children: ReactNode
}

export const Tooltip: React.FC<TooltipProps> = ({ content, children }) => (
  <RadixTooltip.Provider>
    <RadixTooltip.Root delayDuration={200}>
      <RadixTooltip.Trigger asChild>
        <span style={{ display: 'inline-flex', alignItems: 'center', gap: 4 }}>
          {children}
          <span
            style={{
              display: 'inline-flex',
              alignItems: 'center',
              justifyContent: 'center',
              width: 14,
              height: 14,
              borderRadius: '50%',
              border: '1px solid #fefefe',
              color: '#fefefe',
              fontSize: 11,
              fontWeight: 600,
              cursor: 'pointer',
              background: 'none',
            }}
            aria-label="Tooltip info"
            tabIndex={0}
          >
            i
          </span>
        </span>
      </RadixTooltip.Trigger>
      <RadixTooltip.Portal>
        <RadixTooltip.Content
          side="top"
          align="center"
          style={{
            background: 'rgba(0,0,0,0.85)',
            color: '#fff',
            padding: '6px 12px',
            borderRadius: 4,
            fontSize: 14,
            boxShadow: '0 2px 8px rgba(0,0,0,0.15)',
            maxWidth: '90vw',
            minWidth: '120px',
            zIndex: 100,
            whiteSpace: 'pre-line',
          }}
        >
          {content}
        </RadixTooltip.Content>
      </RadixTooltip.Portal>
    </RadixTooltip.Root>
  </RadixTooltip.Provider>
)
