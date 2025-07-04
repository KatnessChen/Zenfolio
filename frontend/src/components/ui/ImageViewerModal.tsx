import React, { useState, useEffect, useCallback } from 'react'
import { Button } from '@/components/ui/button'
import {
  XIcon,
  ZoomInIcon,
  ZoomOutIcon,
  RotateLeftIcon,
  RotateRightIcon,
  ResetIcon,
} from '@/components/icons'

interface ImageViewerModalProps {
  isOpen: boolean
  imageUrl: string
  imageName?: string
  onClose: () => void
}

export const ImageViewerModal: React.FC<ImageViewerModalProps> = ({
  isOpen,
  imageUrl,
  imageName = 'Image',
  onClose,
}) => {
  const [scale, setScale] = useState(1)
  const [rotation, setRotation] = useState(0)
  const [position, setPosition] = useState({ x: 0, y: 0 })
  const [isDragging, setIsDragging] = useState(false)
  const [dragStart, setDragStart] = useState({ x: 0, y: 0 })

  // Reset all transforms when modal opens
  useEffect(() => {
    if (isOpen) {
      setScale(1)
      setRotation(0)
      setPosition({ x: 0, y: 0 })
    }
  }, [isOpen])

  // Handle keyboard shortcuts
  useEffect(() => {
    if (!isOpen) return

    const handleKeyPress = (e: KeyboardEvent) => {
      switch (e.key) {
        case 'Escape':
          onClose()
          break
        case '+':
        case '=':
          handleZoomIn()
          break
        case '-':
          handleZoomOut()
          break
        case 'r':
        case 'R':
          if (e.shiftKey) {
            handleRotateLeft()
          } else {
            handleRotateRight()
          }
          break
        case '0':
          handleReset()
          break
      }
    }

    window.addEventListener('keydown', handleKeyPress)
    return () => window.removeEventListener('keydown', handleKeyPress)
  }, [isOpen, onClose])

  // Handle wheel zoom
  const handleWheel = useCallback((e: React.WheelEvent) => {
    e.preventDefault()
    const delta = e.deltaY > 0 ? -0.1 : 0.1
    setScale((prev) => Math.max(0.1, Math.min(5, prev + delta)))
  }, [])

  // Handle mouse drag
  const handleMouseDown = useCallback(
    (e: React.MouseEvent) => {
      if (e.button === 0) {
        // Left mouse button
        setIsDragging(true)
        setDragStart({ x: e.clientX - position.x, y: e.clientY - position.y })
      }
    },
    [position]
  )

  const handleMouseMove = useCallback(
    (e: React.MouseEvent) => {
      if (isDragging) {
        setPosition({
          x: e.clientX - dragStart.x,
          y: e.clientY - dragStart.y,
        })
      }
    },
    [isDragging, dragStart]
  )

  const handleMouseUp = useCallback(() => {
    setIsDragging(false)
  }, [])

  // Control functions
  const handleZoomIn = () => setScale((prev) => Math.min(5, prev + 0.25))
  const handleZoomOut = () => setScale((prev) => Math.max(0.1, prev - 0.25))
  const handleRotateLeft = () => setRotation((prev) => prev - 90)
  const handleRotateRight = () => setRotation((prev) => prev + 90)
  const handleReset = () => {
    setScale(1)
    setRotation(0)
    setPosition({ x: 0, y: 0 })
  }

  if (!isOpen) return null

  return (
    <div className="fixed inset-0 z-50 bg-black/90 flex items-center justify-center">
      {/* Header with title and close button */}
      <div className="absolute top-4 left-4 right-4 flex justify-between items-center z-10">
        <h3 className="text-white text-lg font-medium truncate">{imageName}</h3>
        <Button
          onClick={onClose}
          variant="ghost"
          size="icon"
          className="text-white hover:bg-white/20"
        >
          <XIcon size={24} />
        </Button>
      </div>

      {/* Control panel */}
      <div className="absolute bottom-4 left-1/2 transform -translate-x-1/2 z-10">
        <div className="bg-black/60 backdrop-blur-sm rounded-lg p-3 flex items-center gap-2">
          <Button
            onClick={handleZoomOut}
            variant="ghost"
            size="icon"
            className="text-white hover:bg-white/20"
            title="Zoom Out (-)"
          >
            <ZoomOutIcon size={20} />
          </Button>

          <span className="text-white text-sm px-2 min-w-[60px] text-center">
            {Math.round(scale * 100)}%
          </span>

          <Button
            onClick={handleZoomIn}
            variant="ghost"
            size="icon"
            className="text-white hover:bg-white/20"
            title="Zoom In (+)"
          >
            <ZoomInIcon size={20} />
          </Button>

          <div className="w-px h-6 bg-white/30 mx-1" />

          <Button
            onClick={handleRotateLeft}
            variant="ghost"
            size="icon"
            className="text-white hover:bg-white/20"
            title="Rotate Left (Shift+R)"
          >
            <RotateLeftIcon size={20} />
          </Button>

          <Button
            onClick={handleRotateRight}
            variant="ghost"
            size="icon"
            className="text-white hover:bg-white/20"
            title="Rotate Right (R)"
          >
            <RotateRightIcon size={20} />
          </Button>

          <div className="w-px h-6 bg-white/30 mx-1" />

          <Button
            onClick={handleReset}
            variant="ghost"
            size="icon"
            className="text-white hover:bg-white/20"
            title="Reset (0)"
          >
            <ResetIcon size={20} />
          </Button>
        </div>
      </div>

      {/* Image container */}
      <div
        className="w-full h-full flex items-center justify-center cursor-move overflow-hidden"
        onWheel={handleWheel}
        onMouseDown={handleMouseDown}
        onMouseMove={handleMouseMove}
        onMouseUp={handleMouseUp}
        onMouseLeave={handleMouseUp}
      >
        <img
          src={imageUrl}
          alt={imageName}
          className="max-w-none max-h-none select-none pointer-events-none"
          style={{
            transform: `translate(${position.x}px, ${position.y}px) scale(${scale}) rotate(${rotation}deg)`,
            transition: isDragging ? 'none' : 'transform 0.2s ease-out',
          }}
          draggable={false}
        />
      </div>

      {/* Help text */}
      <div className="absolute top-4 left-1/2 transform -translate-x-1/2 z-10">
        <div className="bg-black/60 backdrop-blur-sm rounded-lg px-3 py-1">
          <p className="text-white/80 text-xs text-center">
            Scroll to zoom • Drag to pan • ESC to close
          </p>
        </div>
      </div>
    </div>
  )
}
