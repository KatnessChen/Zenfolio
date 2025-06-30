import React, { useState } from 'react'
import { Input } from '@/components/ui/input'
import { Card, CardHeader, CardTitle, CardContent } from '@/components/ui/card'
import { Label } from '@/components/ui/label'
import { Button } from '@/components/ui/button'

const SecurityCard: React.FC = () => {
  const [currentPassword, setCurrentPassword] = useState('')
  const [currentValid, setCurrentValid] = useState(false)
  const [newPassword, setNewPassword] = useState('')
  const [confirmPassword, setConfirmPassword] = useState('')
  const [resetEnabled, setResetEnabled] = useState(false)
  const [showResetInfo, setShowResetInfo] = useState(false)

  // Simulate password validation
  const validateCurrent = () => {
    if (!resetEnabled) {
      return
    }
    // Replace with real validation
    setCurrentValid(currentPassword.length > 0)
  }

  const handleReset = () => {
    // Replace with real reset logic
    setResetEnabled(false)
    setCurrentPassword('')
    setNewPassword('')
    setConfirmPassword('')
    setCurrentValid(false)
  }

  const handleSendResetLink = () => {
    setShowResetInfo(true)
  }

  return (
    <Card className="max-w-xl">
      <CardHeader>
        <CardTitle className="text-xl font-semibold text-foreground">Password Management</CardTitle>
      </CardHeader>
      <CardContent>
        {/* Reset Password Section */}
        <div className="mb-6">
          <div className="mb-2 text-muted-foreground text-sm">
            Change your password if you know your current one.
          </div>
          <Label htmlFor="current-password-input">Current Password</Label>
          <Input
            id="current-password-input"
            type="password"
            value={currentPassword}
            onChange={(e) => setCurrentPassword(e.target.value)}
          />
          <Button
            className="mt-3"
            type="button"
            onClick={validateCurrent}
            disabled={!currentPassword}
            variant="default"
          >
            Check Password
          </Button>
          {currentValid && (
            <div className="mt-6">
              <Label htmlFor="new-password-input">New Password</Label>
              <Input
                id="new-password-input"
                type="password"
                value={newPassword}
                onChange={(e) => setNewPassword(e.target.value)}
              />
              <div className="text-muted-foreground text-xs mb-2">
                Must be at least 8 characters, with letters, numbers, and symbols.
              </div>
              <Label htmlFor="confirm-password-input">Confirm New Password</Label>
              <Input
                id="confirm-password-input"
                type="password"
                value={confirmPassword}
                onChange={(e) => setConfirmPassword(e.target.value)}
                className="mb-4"
              />
              <Button
                variant="default"
                className={`px-6 py-2 rounded font-semibold transition text-gray-100`}
                disabled={
                  !(
                    newPassword &&
                    confirmPassword &&
                    newPassword === confirmPassword &&
                    newPassword.length >= 8
                  )
                }
                onClick={handleReset}
              >
                Reset Password
              </Button>
            </div>
          )}
        </div>
        <div className="flex items-center my-4">
          <div className="flex-1 border-t border-border" />
          <span className="mx-4 text-muted-foreground text-sm">OR</span>
          <div className="flex-1 border-t border-border" />
        </div>
        {/* Forgot Password Section */}
        <div>
          <div className="mb-2 text-muted-foreground text-sm">
            If you've forgotten your password, we can send you a reset link.
          </div>
          <Button
            variant="outline"
            className="px-6 py-2 rounded font-semibold hover:bg-card/80 transition"
            onClick={handleSendResetLink}
          >
            Send Reset Link
          </Button>
          {showResetInfo && (
            <div className="mt-2 text-muted-foreground text-xs">
              A password reset link will be sent to your registered email address, valid for 30
              minutes.
            </div>
          )}
        </div>
      </CardContent>
    </Card>
  )
}

export default SecurityCard
