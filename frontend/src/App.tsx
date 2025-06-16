import { BrowserRouter as Router, Routes, Route } from 'react-router-dom'
import { Provider } from 'react-redux'
import { store } from '@/store'
import Navigation from '@/components/Navigation'
import HomePage from '@/pages/HomePage'
import TransactionExtractPage from '@/pages/TransactionExtractPage'
import TransactionHistoryPage from '@/pages/TransactionHistoryPage'

function App() {
  return (
    <Provider store={store}>
      <Router>
        <div className="min-h-screen bg-background">
          <Navigation />
          <main>
            <Routes>
              <Route path="/" element={<HomePage />} />
              <Route path="/extract" element={<TransactionExtractPage />} />
              <Route path="/history" element={<TransactionHistoryPage />} />
            </Routes>
          </main>
        </div>
      </Router>
    </Provider>
  )
}

export default App
