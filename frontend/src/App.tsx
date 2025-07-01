import { BrowserRouter as Router, Routes, Route } from 'react-router-dom'
import { Provider } from 'react-redux'
import { store } from '@/store'
import { ROUTES } from '@/constants'
import Navigation from '@/components/Navigation'
import Footer from '@/components/Footer'
import ScrollToTop from '@/components/ScrollToTop'
import { ProtectedRoute } from '@/components/ProtectedRoute'
import HomePage from '@/pages/HomePage'
import DashboardPage from '@/pages/DashboardPage'
import SignUpPage from '@/pages/SignUpPage'
import LoginPage from '@/pages/LoginPage'
import TransactionHistoryPage from '@/pages/TransactionHistoryPage'
import ImageUploadPage from '@/pages/ImageUploadPage'
import ProcessingPage from '@/pages/ProcessingPage'
import DataReviewPage from '@/pages/DataReviewPage'
import SettingsPage from '@/pages/SettingsPage'
import ManualTransactionPage from '@/pages/ManualTransactionPage'
import BatchReviewPage from '@/pages/BatchReviewPage'

function App() {
  return (
    <Provider store={store}>
      <Router>
        <ScrollToTop />
        {/* Apply dark class for Zenfolio dark theme */}
        <div className="min-h-screen bg-background dark flex flex-col">
          <div className="fixed top-0 left-0 right-0 z-50 border-b border-border/50">
            <Navigation />
          </div>
          <main className="flex-1 pt-24">
            <Routes>
              {/* Public routes */}
              <Route path={ROUTES.LOGIN} element={<LoginPage />} />
              <Route path={ROUTES.SIGN_UP} element={<SignUpPage />} />

              {/* Protected routes */}
              <Route element={<ProtectedRoute />}>
                <Route path={ROUTES.HOME} element={<HomePage />} />
                <Route path={ROUTES.DASHBOARD} element={<DashboardPage />} />
                <Route path={ROUTES.TRANSACTIONS} element={<TransactionHistoryPage />} />
                <Route path={ROUTES.TRANSACTIONS_UPLOAD} element={<ImageUploadPage />} />
                <Route path={ROUTES.TRANSACTIONS_UPLOAD_PROCESSING} element={<ProcessingPage />} />
                <Route path={ROUTES.TRANSACTIONS_UPLOAD_REVIEW} element={<DataReviewPage />} />
                <Route path={ROUTES.TRANSACTIONS_MANUAL_ADD} element={<ManualTransactionPage />} />
                <Route path={ROUTES.TRANSACTIONS_MANUAL_REVIEW} element={<BatchReviewPage />} />
                <Route path={ROUTES.SETTINGS} element={<SettingsPage />} />
              </Route>
            </Routes>
          </main>
          <Footer />
        </div>
      </Router>
    </Provider>
  )
}

export default App
