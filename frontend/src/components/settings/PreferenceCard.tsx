import React, { useState } from 'react'
import { Dropdown, DropdownItem } from '@/components/ui/dropdown'
import DropdownTrigger from '@/components/ui/dropdown-trigger'
import { Card, CardHeader, CardTitle, CardContent } from '@/components/ui/card'
import { Label } from '@/components/ui/label'
import { CURRENCY, CURRENCY_FLAGS } from '@/constants'
import { Button } from '@/components/ui/button'

const currencies = Object.values(CURRENCY)
const languages = ['English', '繁體中文']

const PreferenceCard: React.FC = () => {
  const [currency, setCurrency] = useState(CURRENCY.USD)
  const [language, setLanguage] = useState('English')
  const [changed, setChanged] = useState(false)

  const handleCurrencyChange = (cur: string) => {
    setCurrency(cur)
    setChanged(true)
  }
  const handleLanguageChange = (lang: string) => {
    setLanguage(lang)
    setChanged(true)
  }

  return (
    <Card className="max-w-xl">
      <CardHeader>
        <CardTitle className="text-xl font-semibold text-foreground">Display Preferences</CardTitle>
      </CardHeader>
      <CardContent>
        <div className="mb-4">
          <Label>Currency</Label>
          <Dropdown
            trigger={
              <DropdownTrigger>
                <span className="inline-flex items-center">
                  <span className="mt-1">{CURRENCY_FLAGS[currency] || ''}</span>
                  <span className="ml-2">{currency}</span>
                </span>
              </DropdownTrigger>
            }
          >
            {currencies.map((cur) => (
              <DropdownItem key={cur} onClick={() => handleCurrencyChange(cur)}>
                <span className="inline-flex items-center">
                  <span className="mt-1">{CURRENCY_FLAGS[cur] || ''}</span>
                  <span className="ml-2">{cur}</span>
                </span>
              </DropdownItem>
            ))}
          </Dropdown>
        </div>
        <div className="mb-6">
          <Label htmlFor="language-dropdown">Language</Label>
          <Dropdown trigger={<DropdownTrigger>{language}</DropdownTrigger>}>
            {languages.map((lang) => (
              <DropdownItem key={lang} onClick={() => handleLanguageChange(lang)}>
                {lang}
              </DropdownItem>
            ))}
          </Dropdown>
        </div>
        <Button disabled={!changed}>Save Changes</Button>
      </CardContent>
    </Card>
  )
}

export default PreferenceCard
