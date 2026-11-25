
import { useState } from 'react'
import { Link, useNavigate } from 'react-router-dom'
import {
  UserRound,
  Mail,
  LockKeyhole,
  ArrowRight,
  Eye,
  EyeOff,
  ShieldCheck,
  Terminal,
  Check,
  X,
} from 'lucide-react'
import { useAuth } from '../context/AuthContext'
import './Register.css'

export default function Register() {
  const { register } = useAuth()
  const navigate = useNavigate()

  const [form, setForm] = useState({
    name: '',
    email: '',
    password: '',
  })

  const [error, setError] = useState('')
  const [success, setSuccess] = useState('')
  const [loading, setLoading] = useState(false)
  const [showPassword, setShowPassword] = useState(false)
  const [focusedField, setFocusedField] = useState('')
  const [characterMood, setCharacterMood] = useState('idle')

  const passwordChecks = [
    {
      label: 'Minimal 8 karakter',
      valid: form.password.length >= 8,
    },
    {
      label: 'Memiliki huruf dan angka',
      valid: /[A-Za-z]/.test(form.password) && /\d/.test(form.password),
    },
    {
      label: 'Memiliki huruf kapital',
      valid: /[A-Z]/.test(form.password),
    },
  ]

  const passwordStrength = passwordChecks.filter(
    (item) => item.valid
  ).length

  const mascotImages = {
    idle: '/images/mascot/mascot-welcome.png',
    thinking: '/images/mascot/mascot-thinking.png',
    working: '/images/mascot/mascot-focus.png',
    happy: '/images/mascot/mascot-happy.png',
    error: '/images/mascot/mascot-sad.png',
  }

  const mascotSrc = mascotImages[characterMood] || mascotImages.idle

  // Dipanggil ketika salah satu textbox mulai aktif
  function handleFieldFocus(fieldName) {
    setFocusedField(fieldName)

    if (!loading && !success) {
      setCharacterMood('thinking')
    }
  }

  // Dipanggil ketika pengguna meninggalkan textbox
  function handleFieldBlur() {
    setFocusedField('')

    if (!loading && !success && characterMood === 'thinking') {
      setCharacterMood('idle')
    }
  }

  function handleChange(event) {
    const { name, value } = event.target

    setForm((current) => ({
      ...current,
      [name]: value,
    }))

    setError('')
    setSuccess('')
  }

  async function handleSubmit(event) {
    event.preventDefault()

    setError('')
    setSuccess('')

    if (form.name.trim().length < 2) {
      setError('Nama minimal 2 karakter.')
      setCharacterMood('error')
      return
    }

    if (form.password.length < 8) {
      setError('Password minimal 8 karakter.')
      setCharacterMood('error')
      return
    }

    setLoading(true)
    setCharacterMood('working')

    try {
      await register(
        form.name.trim(),
        form.email.trim(),
        form.password
      )

      setSuccess('ACCESS GRANTED. Account successfully created.')
      setCharacterMood('happy')

      setTimeout(() => {
        navigate('/login')
      }, 1600)
    } catch (err) {
      setError(err.message || 'Registrasi gagal.')
      setCharacterMood('error')
    } finally {
      setLoading(false)
    }
  }

  return (
    <main className="cyber-auth-page register-page">
      <div className="cyber-grid-background" />
      <div className="cyber-scanline" />
      <div className="cyber-noise" />

      <div className="cyber-topbar">
        <Link to="/login" className="cyber-brand">
          <span className="cyber-brand-icon">T</span>
          <span>
            TASK<span>FLOW</span>
          </span>
        </Link>

        <div className="cyber-system-status">
          <span className="system-live-dot" />
          SYSTEM ONLINE
          <span className="system-version">SYS.02.26</span>
        </div>
      </div>

      <div className="register-layout">
        <section className="register-visual">
          <div className="visual-top-label">
            <Terminal size={14} />
            <span>NEW USER PROTOCOL</span>
          </div>

          <div className="visual-main">
            <div className="character-orbit orbit-one" />
            <div className="character-orbit orbit-two" />

            <div className={`register-character mood-${characterMood}`}>
              <div className="character-image-placeholder">
                <img
                  key={mascotSrc}
                  src={mascotSrc}
                  alt="TaskFlow cyberpunk mascot"
                  onError={(event) => {
                    event.currentTarget.style.display = 'none'
                  }}
                />
              </div>
            </div>

            <div className="visual-floating-tag tag-top">
              <span>01</span>
              IDENTITY
            </div>

            <div className="visual-floating-tag tag-bottom">
              <span>02</span>
              AUTHENTICATION
            </div>
          </div>

          <div className="visual-copy">
            <span className="cyber-eyebrow">WELCOME TO THE GRID</span>

            <h1>
              CREATE
              <br />
              YOUR <span>LEGACY.</span>
            </h1>

            <p>
              Your next chapter starts here.
              Organize your missions, track your progress,
              and take control of your workflow.
            </p>
          </div>

          <div className="visual-footer">
            <span>NEO-TOKYO NETWORK</span>
            <span>LATENCY: 012 MS</span>
            <span>ENCRYPTION: ACTIVE</span>
          </div>
        </section>

        <section className="register-form-section">
          <div className="register-terminal">
            <div className="terminal-heading">
              <div className="terminal-indicator">
                <span />
                <span />
                <span />
              </div>

              <span>IDENTITY_REGISTRATION.EXE</span>

              <span className="terminal-code">[ 001 ]</span>
            </div>

            <div className="register-form-content">
              <div className="form-title-area">
                <span className="cyber-eyebrow">
                  <span className="title-line" />
                  ACCOUNT INITIALIZATION
                </span>

                <h2>
                  JOIN THE
                  <br />
                  <span>NETWORK.</span>
                </h2>

                <p>
                  Create your credentials to access TaskFlow.
                </p>
              </div>

              {error && (
                <div className="cyber-alert cyber-alert-error">
                  <X size={17} />
                  <span>{error}</span>
                </div>
              )}

              {success && (
                <div className="cyber-alert cyber-alert-success">
                  <Check size={17} />
                  <span>{success}</span>
                </div>
              )}

              <form
                onSubmit={handleSubmit}
                className="cyber-register-form"
              >
                {/* FULL NAME */}
                <div
                  className={`cyber-field ${
                    focusedField === 'name' ? 'field-active' : ''
                  }`}
                >
                  <label htmlFor="register-name">
                    <span>01 / FULL NAME</span>
                    <span className="field-required">REQUIRED</span>
                  </label>

                  <div className="cyber-input-wrap">
                    <UserRound size={18} />

                    <input
                      id="register-name"
                      name="name"
                      type="text"
                      placeholder="Enter your identity"
                      value={form.name}
                      onChange={handleChange}
                      onFocus={() => handleFieldFocus('name')}
                      onBlur={handleFieldBlur}
                      minLength={2}
                      maxLength={100}
                      autoComplete="name"
                      required
                    />

                    <span className="input-index">01</span>
                  </div>
                </div>

                {/* EMAIL */}
                <div
                  className={`cyber-field ${
                    focusedField === 'email' ? 'field-active' : ''
                  }`}
                >
                  <label htmlFor="register-email">
                    <span>02 / EMAIL ADDRESS</span>
                    <span className="field-required">REQUIRED</span>
                  </label>

                  <div className="cyber-input-wrap">
                    <Mail size={18} />

                    <input
                      id="register-email"
                      name="email"
                      type="email"
                      placeholder="user@network.com"
                      value={form.email}
                      onChange={handleChange}
                      onFocus={() => handleFieldFocus('email')}
                      onBlur={handleFieldBlur}
                      autoComplete="email"
                      required
                    />

                    <span className="input-index">02</span>
                  </div>
                </div>

                {/* PASSWORD */}
                <div
                  className={`cyber-field ${
                    focusedField === 'password' ? 'field-active' : ''
                  }`}
                >
                  <label htmlFor="register-password">
                    <span>03 / PASSWORD</span>
                    <span className="field-required">REQUIRED</span>
                  </label>

                  <div className="cyber-input-wrap">
                    <LockKeyhole size={18} />

                    <input
                      id="register-password"
                      name="password"
                      type={showPassword ? 'text' : 'password'}
                      placeholder="Create secure password"
                      value={form.password}
                      onChange={handleChange}
                      onFocus={() => handleFieldFocus('password')}
                      onBlur={handleFieldBlur}
                      minLength={8}
                      maxLength={72}
                      autoComplete="new-password"
                      required
                    />

                    <button
                      type="button"
                      className="password-toggle"
                      onClick={() => setShowPassword(!showPassword)}
                      aria-label={
                        showPassword ? 'Hide password' : 'Show password'
                      }
                    >
                      {showPassword ? (
                        <EyeOff size={17} />
                      ) : (
                        <Eye size={17} />
                      )}
                    </button>
                  </div>

                  <div className="password-strength">
                    <div className="strength-bars">
                      {[1, 2, 3].map((item) => (
                        <span
                          key={item}
                          className={
                            passwordStrength >= item
                              ? 'strength-filled'
                              : ''
                          }
                        />
                      ))}
                    </div>

                    <span>
                      {!form.password
                        ? 'SECURITY LEVEL: WAITING'
                        : passwordStrength === 1
                          ? 'SECURITY LEVEL: LOW'
                          : passwordStrength === 2
                            ? 'SECURITY LEVEL: MEDIUM'
                            : 'SECURITY LEVEL: HIGH'}
                    </span>
                  </div>

                  {form.password && (
                    <div className="password-checks">
                      {passwordChecks.map((item) => (
                        <div
                          key={item.label}
                          className={item.valid ? 'check-valid' : ''}
                        >
                          {item.valid ? (
                            <Check size={12} />
                          ) : (
                            <span className="check-empty" />
                          )}
                          {item.label}
                        </div>
                      ))}
                    </div>
                  )}
                </div>

                {/* SUBMIT */}
                <button
                  className="cyber-submit-button"
                  type="submit"
                  disabled={loading || Boolean(success)}
                >
                  <span className="submit-glitch" />

                  <span className="submit-button-content">
                    <ShieldCheck size={18} />
                    {loading
                      ? 'INITIALIZING...'
                      : success
                        ? 'ACCESS GRANTED'
                        : 'CREATE ACCOUNT'}
                  </span>

                  {!loading && !success && (
                    <ArrowRight size={19} />
                  )}
                </button>
              </form>

              <div className="register-divider">
                <span />
                <small>EXISTING NETWORK MEMBER?</small>
                <span />
              </div>

              <Link to="/login" className="cyber-login-link">
                SIGN IN TO YOUR ACCOUNT
                <ArrowRight size={15} />
              </Link>

              <div className="register-security-note">
                <ShieldCheck size={15} />
                <span>Your credentials are transmitted securely.</span>
              </div>
            </div>

            <div className="terminal-bottom">
              <span>TF_AUTH_MODULE</span>
              <span>BUILD 2026.11</span>
              <span>END OF TRANSMISSION</span>
            </div>
          </div>
        </section>
      </div>

      <div className="cyber-bottom-bar">
        <span>© TASKFLOW SYSTEMS</span>
        <span>DESIGNED FOR THE NEXT GENERATION</span>
        <span>ALL SYSTEMS OPERATIONAL</span>
      </div>
    </main>
  )
}
