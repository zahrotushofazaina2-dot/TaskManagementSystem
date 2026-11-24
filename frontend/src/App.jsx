
import { BrowserRouter, Navigate, Route, Routes } from 'react-router-dom'
import { AuthProvider, useAuth } from './context/AuthContext'
import Login from './pages/Login'
import Register from './pages/Register'
import Dashboard from './pages/Dashboard'
import PageTransition from './components/PageTransition'
import './App.css'

function ProtectedRoute({ children, role }) {
  const { user, loading } = useAuth()

  if (loading) {
    return <div className="loading-screen">Memuat aplikasi...</div>
  }

  if (!user) {
    return <Navigate to="/login" replace />
  }

  if (role && user.role !== role) {
    return (
      <Navigate
        to={user.role === 'ADMIN' ? '/admin' : '/dashboard'}
        replace
      />
    )
  }

  return children
}

function AppRoutes() {
  const { user, loading } = useAuth()

  if (loading) {
    return <div className="loading-screen">Memuat aplikasi...</div>
  }

  return (
    <PageTransition>
      <Routes>
        <Route
          path="/login"
          element={
            user
              ? <Navigate to={user.role === 'ADMIN' ? '/admin' : '/dashboard'} replace />
              : <Login />
          }
        />

        <Route
          path="/register"
          element={
            user
              ? <Navigate to={user.role === 'ADMIN' ? '/admin' : '/dashboard'} replace />
              : <Register />
          }
        />

        <Route
          path="/dashboard"
          element={
            <ProtectedRoute role="USER">
              <Dashboard />
            </ProtectedRoute>
          }
        />

        <Route
          path="/admin"
          element={
            <ProtectedRoute role="ADMIN">
              <Dashboard />
            </ProtectedRoute>
          }
        />

        <Route
          path="*"
          element={
            <Navigate
              to={
                user
                  ? user.role === 'ADMIN' ? '/admin' : '/dashboard'
                  : '/login'
              }
              replace
            />
          }
        />
      </Routes>
    </PageTransition>
  )
}

export default function App() {
  return (
    <BrowserRouter>
      <AuthProvider>
        <AppRoutes />
      </AuthProvider>
    </BrowserRouter>
  )
}
