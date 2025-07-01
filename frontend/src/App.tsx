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
              <Route
                path={ROUTES.HOME}
                element={
                  <ProtectedRoute>
                    <HomePage />
                  </ProtectedRoute>
                }
              />
              <Route
                path={ROUTES.DASHBOARD}
                element={
                  <ProtectedRoute>
                    <DashboardPage />
                  </ProtectedRoute>
                }
              />
              <Route
                path={ROUTES.TRANSACTIONS}
                element={
                  <ProtectedRoute>
                    <TransactionHistoryPage />
                  </ProtectedRoute>
                }
              />
              <Route
                path={ROUTES.TRANSACTIONS_UPLOAD}
                element={
                  <ProtectedRoute>
                    <ImageUploadPage />
                  </ProtectedRoute>
                }
              />
              <Route
                path={ROUTES.TRANSACTIONS_UPLOAD_PROCESSING}
                element={
                  <ProtectedRoute>
                    <ProcessingPage />
                  </ProtectedRoute>
                }
              />
              <Route
                path={ROUTES.TRANSACTIONS_UPLOAD_REVIEW}
                element={
                  <ProtectedRoute>
                    <DataReviewPage />
                  </ProtectedRoute>
                }
              />
              <Route
                path={ROUTES.TRANSACTIONS_MANUAL_ADD}
                element={
                  <ProtectedRoute>
                    <ManualTransactionPage />
                  </ProtectedRoute>
                }
              />
              <Route
                path={ROUTES.TRANSACTIONS_MANUAL_REVIEW}
                element={
                  <ProtectedRoute>
                    <BatchReviewPage />
                  </ProtectedRoute>
                }
              />
              <Route
                path={ROUTES.SETTINGS}
                element={
                  <ProtectedRoute>
                    <SettingsPage />
                  </ProtectedRoute>
                }
              />
            </Routes>
          </main>
          <Footer />
        </div>
      </Router>
    </Provider>
  )
}

export default App
