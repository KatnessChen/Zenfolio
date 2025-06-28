import React from 'react'
import { Selector } from '@/components/ui/selector'

const navItems = [
  { key: 'profile', label: 'Profile' },
  { key: 'preference', label: 'Preference' },
  { key: 'security', label: 'Security' },
]

interface SettingsSidebarProps {
  activeTab: string
  setActiveTab: (tab: 'profile' | 'preference' | 'security') => void
}

const SettingsSidebar: React.FC<SettingsSidebarProps> = ({ activeTab, setActiveTab }) => (
  <Selector
    options={navItems.map((item) => ({ value: item.key, label: item.label }))}
    value={activeTab}
    onChange={(val) => setActiveTab(val as 'profile' | 'preference' | 'security')}
    type="vertical"
    className="w-56"
  />
)

export default SettingsSidebar
