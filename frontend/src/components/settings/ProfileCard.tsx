import React, { useState } from 'react'
import { Input } from '@/components/ui/input'
import { Card, CardHeader, CardTitle, CardContent } from '@/components/ui/card'
import { Label } from '@/components/ui/label'

const ProfileCard: React.FC = () => {
  const [firstName, setFirstName] = useState('John')
  const [lastName, setLastName] = useState('Doe')
  const [changed, setChanged] = useState(false)

  const handleChange = (setter: (v: string) => void, value: string) => {
    setter(value)
    setChanged(true)
  }

  return (
    <Card className="max-w-xl">
      <CardHeader>
        <CardTitle className="text-xl font-semibold text-foreground">
          Personal Information
        </CardTitle>
      </CardHeader>
      <CardContent>
        <div className="mb-4">
          <Label htmlFor="first-name-input">First Name</Label>
          <Input
            id="first-name-input"
            value={firstName}
            onChange={(e) => handleChange(setFirstName, e.target.value)}
            placeholder="John"
          />
        </div>
        <div className="mb-6">
          <Label htmlFor="last-name-input">Last Name</Label>
          <Input
            id="last-name-input"
            value={lastName}
            onChange={(e) => handleChange(setLastName, e.target.value)}
            placeholder="Doe"
          />
        </div>
        <button
          className={`px-6 py-2 rounded font-semibold transition text-gray-100 ${
            changed ? 'bg-primary/80 hover:bg-primary' : 'bg-muted cursor-not-allowed'
          }`}
          disabled={!changed}
        >
          Save Changes
        </button>
      </CardContent>
    </Card>
  )
}

export default ProfileCard
