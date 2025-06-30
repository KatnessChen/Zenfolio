import { useState } from 'react'
import ProfileCard from '@/components/settings/ProfileCard'
import PreferenceCard from '@/components/settings/PreferenceCard'
import SecurityCard from '@/components/settings/SecurityCard'
import SettingsSidebar from '@/components/settings/SettingsSidebar'
import Title from '@/components/ui/title'

const TABS = [
  { key: 'profile', component: ProfileCard },
  { key: 'preference', component: PreferenceCard },
  { key: 'security', component: SecurityCard },
] as const

type TabKey = (typeof TABS)[number]['key']

export default function SettingsPage() {
  const [activeTab, setActiveTab] = useState<TabKey>(TABS[0][`key`])

  return (
    <div className="min-h-screen">
      {/* Page Title - Top Level */}
      <div className="container mx-auto px-4 py-6">
        <Title as="h1">Settings</Title>
      </div>
      {/* Main Content */}
      <main className="container mx-auto px-4 flex gap-6">
        {/* Sidebar */}
        <aside className="hidden md:block flex-shrink-0">
          <SettingsSidebar activeTab={activeTab} setActiveTab={setActiveTab} />
        </aside>
        {/* Main Area */}
        <section className="flex-1 min-w-0">
          <div>{TABS.map((tab) => activeTab === tab.key && <tab.component key={tab.key} />)}</div>
        </section>
      </main>
    </div>
  )
}
