import React, { useState } from 'react'
import { Dropdown, DropdownItem } from '@/components/ui/dropdown'
import DropdownTrigger from '@/components/ui/dropdown-trigger'
import { Card, CardHeader, CardTitle, CardContent } from '@/components/ui/card'
import { Label } from '@/components/ui/label'

const currencies = ['USD', 'TWD', 'CAD', 'EUR', 'GBP', 'JPY', 'AUD']
const languages = ['English', '繁體中文']

const currencyFlags: Record<string, string> = {
  USD: '🇺🇸',
  TWD: '🇹🇼',
  CAD: '🇨🇦',
  EUR: '🇪🇺',
  GBP: '🇬🇧',
  JPY: '🇯🇵',
  AUD: '🇦🇺',
}

const PreferenceCard: React.FC = () => {
  const [currency, setCurrency] = useState('USD')
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
                  <span className="mt-1">{currencyFlags[currency] || ''}</span>
                  <span className="ml-2">{currency}</span>
                </span>
              </DropdownTrigger>
            }
          >
            {currencies.map((cur) => (
              <DropdownItem key={cur} onClick={() => handleCurrencyChange(cur)}>
                <span className="inline-flex items-center">
                  <span className="mt-1">{currencyFlags[cur] || ''}</span>
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
        <button
          className={`px-6 py-2 rounded font-semibold transition text-gray-100 ${
            changed ? 'bg-[#9DC0B2] hover:bg-[#7BAA97]' : 'bg-[#3A4A43] cursor-not-allowed'
          }`}
          disabled={!changed}
        >
          Save Changes
        </button>
      </CardContent>
    </Card>
  )
}

export default PreferenceCard
