import React, { useState } from 'react'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'

interface ExpandableCardProps {
  title: string
  description: string
  defaultExpanded?: boolean
  children: React.ReactNode
  className?: string
}

export const ExpandableCard: React.FC<ExpandableCardProps> = ({
  title,
  description,
  defaultExpanded = false,
  children,
  className,
}) => {
  const [isExpanded, setIsExpanded] = useState(defaultExpanded)

  return (
    <Card className={className}>
      <CardHeader className="cursor-pointer" onClick={() => setIsExpanded(!isExpanded)}>
        <div className="flex items-center justify-between">
          <div>
            <CardTitle>{title}</CardTitle>
            <CardDescription>{description}</CardDescription>
          </div>
          <div
            className="text-muted-foreground transition-transform duration-200"
            style={{ transform: isExpanded ? 'rotate(180deg)' : 'rotate(0deg)' }}
          >
            ▼
          </div>
        </div>
      </CardHeader>
      {isExpanded && <CardContent className="pt-0">{children}</CardContent>}
    </Card>
  )
}

export default ExpandableCard
