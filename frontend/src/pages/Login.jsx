
import { useState } from 'react'
import { Link, useNavigate } from 'react-router-dom'
import {
  LockKeyhole,
  Mail,
  ArrowRight,
  Sparkles,
  CheckCircle2,
  Eye,
  EyeOff,
} from 'lucide-react'
import { useAuth } from '../context/AuthContext'
import './Login.css'

function AnimeCompanion({ mood }) {
  const messages = {
    welcome: {
      title: 'Welcome back!',
      subtitle: 'Your workspace is waiting for you.',
    },
    focus: {
      title: 'Take your time.',
      subtitle: 'I am paying attention. No pressure.',
    },
    loading: {
      title: 'One moment...',
      subtitle: 'Checking your credentials.',
    },
    error: {
      title: 'Something went wrong!',
      subtitle: 'Let us check those details again.',
    },
    success: {
      title: 'Welcome aboard!',
      subtitle: 'Your next chapter starts now.',
    },
  }

  const mascotImages = {
    welcome: '/images/mascot/mascot-welcome.png',
    focus: '/images/mascot/mascot-thinking.png',
    loading: '/images/mascot/mascot-focus.png',
    error: '/images/mascot/mascot-sad.png',
    success: '/images/mascot/mascot-happy.png',
  }

  const current = messages[mood] || messages.welcome
  const currentImage = mascotImages[mood] || mascotImages.welcome

  return (
    <div className={`companion-area companion-${mood}`}>
      <div className="companion-orbit orbit-one" />
      <div className="companion-orbit orbit-two" />

      <div className="companion-sparkle sparkle-one">✦</div>
      <div className="companion-sparkle sparkle-two">✧</div>
      <div className="companion-sparkle sparkle-three">✦</div>

      <div className="companion-card">
        <div className="companion-status">
          <span className="status-indicator" />
          TASKFLOW COMPANION
        </div>

        <div className="anime-character">
          <div className="character-halo" />

          <img
            src={currentImage}
            alt="TaskFlow companion"
            className="companion-image"
          />
        </div>

        <div className="companion-dialog">
          <span className="dialog-decoration">✧</span>
          <h3>{current.title}</h3>
          <p>{current.subtitle}</p>
        </div>
      </div>

      <div className="floating-note note-one">
        <Sparkles size={15} />
        <span>Stay productive</span>
      </div>

      <div className="floating-note note-two">
        <CheckCircle2 size={15} />
        <span>One step at a time</span>
      </div>
    </div>
  )
}

export default function Login() {
  const { login } = useAuth()
  const navigate = useNavigate()

  const [form, setForm] = useState({
    email: '',
    password: '',
  })

  const [error, setError] = useState('')
  const [loading, setLoading] = useState(false)
  const [showPassword, setShowPassword] = useState(false)
  const [focusedField, setFocusedField] = useState('')
  const [success, setSuccess] = useState(false)

  const mood = success
    ? 'success'
    : loading
      ? 'loading'
      : error
        ? 'error'
        : focusedField
          ? 'focus'
          : 'welcome'

  function handleChange(event) {
    const { name, value } = event.target

    setForm((current) => ({
      ...current,
      [name]: value,
    }))

    if (error) setError('')
  }

  async function handleSubmit(event) {
    event.preventDefault()
    setError('')
    setLoading(true)

    try {
      const user = await login(form.email, form.password)

      setSuccess(true)

      await new Promise((resolve) => setTimeout(resolve, 900))

      navigate(
        user.role === 'ADMIN' || user.role === 'SUPER_ADMIN'
          ? '/admin'
          : '/dashboard'
      )
    } catch (err) {
      setError(err.message || 'Login gagal. Periksa kembali data kamu.')
    } finally {
      setLoading(false)
    }
  }

  return (
    <main className="auth-page">
      <div className="auth-background">
        <div className="auth-glow glow-left" />
        <div className="auth-glow glow-right" />
        <div className="auth-grid-pattern" />

        <span className="background-star star-a">✦</span>
        <span className="background-star star-b">✧</span>
        <span className="background-star star-c">✦</span>
        <span className="background-star star-d">✧</span>
      </div>

      <section className="auth-brand">
        <Link to="/login" className="auth-brand-logo">
          <div className="brand-symbol">
            <span>T</span>
            <i>✦</i>
          </div>

          <div>
            <strong>TaskFlow</strong>
            <small>YOUR PRODUCTIVITY SPACE</small>
          </div>
        </Link>

        <div className="auth-brand-copy">
          <span className="auth-section-label">
            <span /> THE MIDNIGHT ACADEMY
          </span>

          <h1>
            Make your
            <br />
            own <span>progress.</span>
          </h1>

          <p>
            A thoughtful space to organize your work,
            track your progress, and turn small steps
            into meaningful achievements.
          </p>

          <div className="auth-feature-list">
            <div>
              <CheckCircle2 size={16} />
              <span>Organized workspace</span>
            </div>
            <div>
              <CheckCircle2 size={16} />
              <span>Progress tracking</span>
            </div>
            <div>
              <CheckCircle2 size={16} />
              <span>Focus on what matters</span>
            </div>
          </div>
        </div>

        <AnimeCompanion mood={mood} />

        <div className="auth-brand-footer">
          <span>DESIGNED FOR BETTER WORKDAYS</span>
          <span>EST. TASKFLOW</span>
        </div>
      </section>

      <section className="auth-form-section">
        <div className="auth-form-container">
          <div className="mobile-brand">
            <div className="brand-symbol">
              <span>T</span>
              <i>✦</i>
            </div>
            <strong>TaskFlow</strong>
          </div>

          <div className="form-top-decoration">
            <span>✦</span>
            <span>✧</span>
            <span>✦</span>
          </div>

          <span className="eyebrow">MEMBER ACCESS</span>

          <h2>Welcome back.</h2>

          <p className="form-description">
            Sign in to continue your journey.
          </p>

          {error && (
            <div className="alert-error">
              <span>!</span>
              {error}
            </div>
          )}

          {success && (
            <div className="alert-success">
              <CheckCircle2 size={17} />
              Login successful. Preparing your workspace...
            </div>
          )}

          <form onSubmit={handleSubmit}>
            <label htmlFor="email">Email address</label>

            <div
              className={`input-wrapper ${
                focusedField === 'email' ? 'input-focused' : ''
              }`}
            >
              <Mail size={18} />

              <input
                id="email"
                name="email"
                type="email"
                placeholder="name@example.com"
                value={form.email}
                onChange={handleChange}
                onFocus={() => setFocusedField('email')}
                onBlur={() => setFocusedField('')}
                autoComplete="email"
                required
              />

              {form.email && (
                <CheckCircle2
                  className="input-check"
                  size={16}
                />
              )}
            </div>

            <label htmlFor="password">Password</label>

            <div
              className={`input-wrapper ${
                focusedField === 'password' ? 'input-focused' : ''
              }`}
            >
              <LockKeyhole size={18} />

              <input
                id="password"
                name="password"
                type={showPassword ? 'text' : 'password'}
                placeholder="Enter your password"
                value={form.password}
                onChange={handleChange}
                onFocus={() => setFocusedField('password')}
                onBlur={() => setFocusedField('')}
                autoComplete="current-password"
                required
              />

              <button
                type="button"
                className="password-toggle"
                onClick={() => setShowPassword((value) => !value)}
                aria-label={showPassword ? 'Hide password' : 'Show password'}
              >
                {showPassword ? (
                  <EyeOff size={17} />
                ) : (
                  <Eye size={17} />
                )}
              </button>
            </div>

            <button
              className="auth-submit"
              type="submit"
              disabled={loading || success}
            >
              <span>
                {loading
                  ? 'Signing in...'
                  : success
                    ? 'Welcome!'
                    : 'Sign in'}
              </span>

              {loading ? (
                <span className="button-loader" />
              ) : (
                <ArrowRight size={18} />
              )}
            </button>
          </form>

          <div className="auth-divider">
            <span />
            <i>✧</i>
            <span />
          </div>

          <p className="auth-switch">
            Don't have an account?{' '}
            <Link to="/register">Create account</Link>
          </p>

          <p className="auth-security">
            <LockKeyhole size={13} />
            Your workspace is protected and secure.
          </p>
        </div>
      </section>
    </main>
  )
}
