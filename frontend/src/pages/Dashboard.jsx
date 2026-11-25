
import { useCallback, useEffect, useRef, useState } from 'react'
import {
  Activity,
  AlertCircle,
  CheckCircle2,
  ClipboardList,
  Clock3,
  Folder,
  LogOut,
  Plus,
  RefreshCw,
  Search,
  ShieldCheck,
  Trash2,
  Users,
  X,
  Pencil,
  Tag,
  Save,
} from 'lucide-react'

import { useAuth } from '../context/AuthContext'
import { api } from '../lib/api'
import './Dashboard.css'
import AdminUsers from './AdminUsers'

const PAGE_LIMIT = 10

const emptyForm = {
  category_id: '',
  title: '',
  description: '',
  priority: 'MEDIUM',
  status: 'TODO',
  due_date: '',
}

const mascotMoods = {
  idle: {
    emoji: '✦',
    title: 'Standing by',
    message: "I'm ready to accompany you!",
  },
  focus: {
    emoji: '⌨',
    title: 'Focus mode',
    message: 'Working on something?',
  },
  working: {
    emoji: '⚙',
    title: 'Processing',
    message: 'Just a moment, processing...',
  },
  happy: {
    emoji: '✧',
    title: 'Mission complete',
    message: 'Nice! Successfully completed.',
  },
  excited: {
    emoji: '★',
    title: 'Amazing!',
    message: 'You just completed a task!',
  },
  confused: {
    emoji: '?',
    title: 'Something is wrong',
    message: 'Hmm, please check your input.',
  },
  sad: {
    emoji: '!',
    title: 'System error',
    message: "Something went wrong. Let's try again.",
  },
}

function unwrapData(response) {
  return response?.data ?? response ?? {}
}

function extractCategories(response) {
  const data = unwrapData(response)

  if (Array.isArray(data)) return data
  if (Array.isArray(data.categories)) return data.categories
  if (Array.isArray(data.data)) return data.data
  if (Array.isArray(data.data?.categories)) return data.data.categories

  return []
}

function formatDate(value) {
  if (!value) return '-'

  const date = new Date(value)

  if (Number.isNaN(date.getTime())) return '-'

  return date.toLocaleDateString('id-ID', {
    day: '2-digit',
    month: 'short',
    year: 'numeric',
  })
}

function StatCard({
  icon: Icon,
  label,
  value,
  detail,
  className = '',
}) {
  return (
    <article className={`stat-card ${className}`}>
      <div className="stat-card-top">
        <span>{label}</span>
        <div className="stat-icon">
          <Icon size={19} />
        </div>
      </div>

      <strong>{value ?? 0}</strong>
      {detail && <small>{detail}</small>}
    </article>
  )
}

function StatusBadge({ status }) {
  const normalized = String(status || 'TODO').toUpperCase()

  const labels = {
    TODO: 'To Do',
    IN_PROGRESS: 'In Progress',
    DONE: 'Done',
  }

  return (
    <span className={`status-badge status-${normalized.toLowerCase()}`}>
      {labels[normalized] || normalized}
    </span>
  )
}

function PriorityBadge({ priority }) {
  const normalized = String(priority || 'MEDIUM').toUpperCase()

  return (
    <span className={`priority-badge priority-${normalized.toLowerCase()}`}>
      {normalized}
    </span>
  )
}

function Pagination({ page, totalPages, onPageChange }) {
  const pages = Math.max(1, totalPages || 1)

  return (
    <div className="pagination">
      <span>
        page {page} of {pages}
      </span>

      <div>
        <button
          className="pagination-button"
          disabled={page <= 1}
          onClick={() => onPageChange(page - 1)}
        >
          Prev
        </button>

        <button
          className="pagination-button"
          disabled={page >= pages}
          onClick={() => onPageChange(page + 1)}
        >
          Next
        </button>
      </div>
    </div>
  )
}

/* =========================================
   DASHBOARD MASCOT
========================================= */

function MascotCompanion({ mood = 'idle' }) {
  const currentMood = mascotMoods[mood] || mascotMoods.idle

  const mascotImages = {
    idle: '/images/mascot/mascot-idle.png',
    focus: '/images/mascot/mascot-focus.png',
    working: '/images/mascot/mascot-focus.png',
    happy: '/images/mascot/mascot-happy.png',
    excited: '/images/mascot/mascot-excited.png',
    confused: '/images/mascot/mascot-confused.png',
    sad: '/images/mascot/mascot-sad.png',
  }

  const mascotSrc = mascotImages[mood] || mascotImages.idle

  return (
    <aside
      className={`welcome-mascot mascot-mood-${mood}`}
      aria-label="TaskFlow mascot companion"
    >
      <div className="mascot-glow" />

      <div className="mascot-image-frame">
        <img
          key={mascotSrc}
          className="mascot-image"
          src={mascotSrc}
          alt={`TaskFlow mascot ${mood}`}
          onError={(event) => {
            event.currentTarget.style.display = 'none'
            event.currentTarget.nextElementSibling.style.display = 'flex'
          }}
        />

        <div className="mascot-fallback">
          <div className="mascot-fallback-status">
            <span />
            SYSTEM COMPANION
          </div>

          <div className="dashboard-anime-character">
            <div className="dashboard-character-halo" />
            <div className="dashboard-character-hair-back" />
            <div className="dashboard-character-neck" />

            <div className="dashboard-character-face">
              <span className="dashboard-character-ear left" />
              <span className="dashboard-character-ear right" />

              <div className="dashboard-character-eyebrow left" />
              <div className="dashboard-character-eyebrow right" />

              <div className="dashboard-character-eye left">
                <span />
              </div>

              <div className="dashboard-character-eye right">
                <span />
              </div>

              <div className="dashboard-character-blush left" />
              <div className="dashboard-character-blush right" />
              <div className="dashboard-character-mouth" />
            </div>

            <div className="dashboard-character-hair">
              <span className="dashboard-hair-strand strand-one" />
              <span className="dashboard-hair-strand strand-two" />
              <span className="dashboard-hair-strand strand-three" />
            </div>

            <div className="dashboard-character-accessory">✦</div>

            <div className="dashboard-character-body">
              <span className="dashboard-character-collar left" />
              <span className="dashboard-character-collar right" />
              <span className="dashboard-character-tie" />
            </div>

            <div className="dashboard-character-shadow" />
          </div>

          <div className="mascot-fallback-caption">
            <strong>TaskFlow Companion</strong>
            <small>READY TO HELP YOU</small>
          </div>
        </div>

        <div className="mascot-emotion-indicator">
          {currentMood.emoji}
        </div>
      </div>

      <div className="mascot-dialog">
        <span className="mascot-dialog-status">
          {currentMood.title}
        </span>

        <p>{currentMood.message}</p>
      </div>

      <span className="mascot-sparkle sparkle-one">✦</span>
      <span className="mascot-sparkle sparkle-two">✧</span>
      <span className="mascot-sparkle sparkle-three">✦</span>
    </aside>
  )
}

/* =========================================
   MAIN DASHBOARD
========================================= */

export default function Dashboard() {
  const { user, logout } = useAuth()

  // ADMIN dan SUPER_ADMIN menggunakan dashboard admin.
  const isAdmin =
    user?.role === 'ADMIN' || user?.role === 'SUPER_ADMIN'

  const [adminView, setAdminView] = useState('overview')

  const [stats, setStats] = useState(null)
  const [tasks, setTasks] = useState([])
  const [categories, setCategories] = useState([])

  const [pagination, setPagination] = useState({
    total: 0,
    page: 1,
    limit: PAGE_LIMIT,
    total_pages: 1,
  })

  const [page, setPage] = useState(1)
  const [searchInput, setSearchInput] = useState('')
  const [search, setSearch] = useState('')
  const [statusFilter, setStatusFilter] = useState('')
  const [priorityFilter, setPriorityFilter] = useState('')

  const [loading, setLoading] = useState(true)
  const [error, setError] = useState('')
  const [refreshing, setRefreshing] = useState(false)

  const [modalOpen, setModalOpen] = useState(false)
  const [editingTask, setEditingTask] = useState(null)
  const [form, setForm] = useState(emptyForm)
  const [saving, setSaving] = useState(false)
  const [formError, setFormError] = useState('')

  const [categoryModalOpen, setCategoryModalOpen] = useState(false)
  const [editingCategory, setEditingCategory] = useState(null)
  const [categoryName, setCategoryName] = useState('')
  const [categorySaving, setCategorySaving] = useState(false)
  const [categoryError, setCategoryError] = useState('')

  /* MASCOT MOOD */

  const [mascotMood, setMascotMood] = useState('idle')
  const mascotTimer = useRef(null)

  const triggerMascotMood = useCallback((mood, duration = 2600) => {
    setMascotMood(mood)

    if (mascotTimer.current) {
      clearTimeout(mascotTimer.current)
    }

    if (mood !== 'idle') {
      mascotTimer.current = setTimeout(() => {
        setMascotMood('idle')
        mascotTimer.current = null
      }, duration)
    }
  }, [])

  useEffect(() => {
    return () => {
      if (mascotTimer.current) {
        clearTimeout(mascotTimer.current)
      }
    }
  }, [])

  /* =========================================
     LOAD DASHBOARD
  ========================================= */

  const loadDashboard = useCallback(async (showRefresh = false) => {
    if (showRefresh) {
      setRefreshing(true)
    } else {
      setLoading(true)
    }

    setError('')
    triggerMascotMood('working', 10000)

    try {
      if (isAdmin) {
        const [dashboardResponse, taskResponse] = await Promise.all([
          api.get('/api/v1/admin/dashboard'),
          api.get(`/api/v1/admin/tasks?page=${page}&limit=${PAGE_LIMIT}`),
        ])

        const dashboardData = unwrapData(dashboardResponse)
        const taskData = unwrapData(taskResponse)

        setStats(dashboardData)
        setTasks(taskData.tasks || [])

        setPagination({
          total: taskData.total || 0,
          page: taskData.page || page,
          limit: taskData.limit || PAGE_LIMIT,
          total_pages: taskData.total_pages || 1,
        })
      } else {
        const params = new URLSearchParams({
          page: String(page),
          limit: String(PAGE_LIMIT),
          sort_by: 'created_at',
          sort_order: 'DESC',
        })

        if (search) params.set('search', search)
        if (statusFilter) params.set('status', statusFilter)
        if (priorityFilter) params.set('priority', priorityFilter)

        const [
          dashboardResponse,
          taskResponse,
          categoryResponse,
        ] = await Promise.all([
          api.get('/api/v1/dashboard'),
          api.get(`/api/v1/tasks?${params.toString()}`),
          api.get('/api/v1/categories'),
        ])

        const dashboardData = unwrapData(dashboardResponse)
        const taskData = unwrapData(taskResponse)

        setStats(dashboardData)
        setTasks(taskData.tasks || [])

        setPagination(
          taskData.pagination || {
            total: 0,
            page,
            limit: PAGE_LIMIT,
            total_pages: 1,
          }
        )

        setCategories(extractCategories(categoryResponse))
      }

      setMascotMood('idle')
    } catch (err) {
      setError(err.message || 'Gagal memuat dashboard.')
      triggerMascotMood('sad', 4000)
    } finally {
      setLoading(false)
      setRefreshing(false)
    }
  }, [
    isAdmin,
    page,
    search,
    statusFilter,
    priorityFilter,
    triggerMascotMood,
  ])

  useEffect(() => {
    loadDashboard()
  }, [loadDashboard])

  /* =========================================
     TASK HANDLERS
  ========================================= */

  function openCreateModal() {
    setEditingTask(null)
    setForm(emptyForm)
    setFormError('')
    setModalOpen(true)
    triggerMascotMood('focus')
  }

  function openEditModal(task) {
    setEditingTask(task)

    setForm({
      category_id: task.category_id
        ? String(task.category_id)
        : '',
      title: task.title || '',
      description: task.description || '',
      priority: String(task.priority || 'MEDIUM').toUpperCase(),
      status: String(task.status || 'TODO').toUpperCase(),
      due_date: task.due_date
        ? String(task.due_date).slice(0, 10)
        : '',
    })

    setFormError('')
    setModalOpen(true)
    triggerMascotMood('focus')
  }

  function updateForm(field, value) {
    setForm((current) => ({
      ...current,
      [field]: value,
    }))

    triggerMascotMood('focus', 1800)
  }

  async function handleSubmit(event) {
    event.preventDefault()
    setFormError('')

    if (form.title.trim().length < 2) {
      setFormError('Judul task minimal 2 karakter.')
      triggerMascotMood('confused', 3500)
      return
    }

    if (form.title.trim().length > 200) {
      setFormError('Judul task maksimal 200 karakter.')
      triggerMascotMood('confused', 3500)
      return
    }

    if (form.description.length > 5000) {
      setFormError('Deskripsi maksimal 5000 karakter.')
      triggerMascotMood('confused', 3500)
      return
    }

    const payload = {
      category_id: form.category_id
        ? Number(form.category_id)
        : null,
      title: form.title.trim(),
      description: form.description.trim() || null,
      priority: form.priority,
      status: form.status,
      due_date: form.due_date || null,
    }

    setSaving(true)
    triggerMascotMood('working', 10000)

    try {
      if (editingTask) {
        await api.put(
          `/api/v1/tasks/${editingTask.id}`,
          payload
        )
      } else {
        await api.post('/api/v1/tasks', payload)
      }

      setModalOpen(false)
      setForm(emptyForm)

      await loadDashboard(true)

      triggerMascotMood(
        payload.status === 'DONE' ? 'excited' : 'happy',
        4000
      )
    } catch (err) {
      setFormError(err.message || 'Gagal menyimpan task.')
      triggerMascotMood('sad', 4000)
    } finally {
      setSaving(false)
    }
  }

  async function handleDelete(task) {
    const confirmed = window.confirm(
      `Yakin ingin menghapus task "${task.title}"?`
    )

    if (!confirmed) return

    triggerMascotMood('working', 10000)

    try {
      await api.delete(`/api/v1/tasks/${task.id}`)
      await loadDashboard(true)
      triggerMascotMood('happy', 3500)
    } catch (err) {
      window.alert(err.message || 'Gagal menghapus task.')
      triggerMascotMood('sad', 4000)
    }
  }

  /* =========================================
     CATEGORY HANDLERS
  ========================================= */

  function openCreateCategoryModal() {
    setEditingCategory(null)
    setCategoryName('')
    setCategoryError('')
    setCategoryModalOpen(true)
    triggerMascotMood('focus')
  }

  function openEditCategoryModal(category) {
    setEditingCategory(category)
    setCategoryName(category.name || '')
    setCategoryError('')
    setCategoryModalOpen(true)
    triggerMascotMood('focus')
  }

  async function handleCategorySubmit(event) {
    event.preventDefault()
    setCategoryError('')

    const trimmedName = categoryName.trim()

    if (!trimmedName) {
      setCategoryError('Category name is required.')
      triggerMascotMood('confused', 3500)
      return
    }

    if (trimmedName.length > 100) {
      setCategoryError('Category name maximum 100 characters.')
      triggerMascotMood('confused', 3500)
      return
    }

    setCategorySaving(true)
    triggerMascotMood('working', 10000)

    try {
      if (editingCategory) {
        await api.put(
          `/api/v1/categories/${editingCategory.id}`,
          { name: trimmedName }
        )
      } else {
        await api.post('/api/v1/categories', {
          name: trimmedName,
        })
      }

      setCategoryModalOpen(false)
      setEditingCategory(null)
      setCategoryName('')

      await loadDashboard(true)
      triggerMascotMood('happy', 3500)
    } catch (err) {
      setCategoryError(err.message || 'Failed to save category.')
      triggerMascotMood('sad', 4000)
    } finally {
      setCategorySaving(false)
    }
  }

  async function handleDeleteCategory(category) {
    const confirmed = window.confirm(
      `Are you sure you want to delete the category? "${category.name}"? Tasks that use this category will not be deleted.`
    )

    if (!confirmed) return

    triggerMascotMood('working', 10000)

    try {
      await api.delete(`/api/v1/categories/${category.id}`)
      await loadDashboard(true)
      triggerMascotMood('happy', 3500)
    } catch (err) {
      window.alert(err.message || 'Failed to delete category.')
      triggerMascotMood('sad', 4000)
    }
  }

  /* =========================================
     FILTER HANDLERS
  ========================================= */

  function handleSearch(event) {
    event.preventDefault()
    setPage(1)
    setSearch(searchInput.trim())
    triggerMascotMood('focus')
  }

  function clearFilters() {
    setSearchInput('')
    setSearch('')
    setStatusFilter('')
    setPriorityFilter('')
    setPage(1)
    triggerMascotMood('happy')
  }

  const totalPages = pagination.total_pages || 1

  return (
    <main className="dashboard-page">
      <header className="dashboard-header">
        <div className="dashboard-brand">
          <div className="brand-symbol">
            <ClipboardList size={23} />
          </div>

          <div>
            <h1>TaskFlow</h1>
            <p>Task Management System</p>
          </div>
        </div>

        <div className="header-user">
          <div className="user-avatar">
            {user?.name?.charAt(0)?.toUpperCase() || 'U'}
          </div>

          <div className="user-details">
            <strong>{user?.name || 'User'}</strong>
            <span>{user?.role || 'USER'}</span>
          </div>

          <button
            className="logout-button"
            onClick={logout}
            title="Logout"
          >
            <LogOut size={18} />
            <span>Logout</span>
          </button>
        </div>
      </header>

      <MascotCompanion mood={mascotMood} />

      <section className="dashboard-content">
        <div className="welcome-section">
          <div className="welcome-copy">
            <span className="eyebrow">
              {isAdmin ? 'ADMINISTRATION' : 'PERSONAL WORKSPACE'}
            </span>

            <h2>
              {isAdmin
                ? 'Welcome back, Administrator!'
                : `Halo, ${user?.name || 'User'}!`}
            </h2>

            <p>
              {isAdmin
                ? 'Monitor system activity and manage TaskFlow users.'
                : 'Every small step brings you closer to your goals.'}
            </p>

            <div className="welcome-meta">
              <span className="welcome-status-dot" />
              <span>Workspace is ready</span>
            </div>
          </div>

          <button
            className="refresh-button"
            onClick={() => loadDashboard(true)}
            disabled={loading || refreshing}
          >
            <RefreshCw
              size={17}
              className={refreshing ? 'spinning' : ''}
            />
            Refresh
          </button>
        </div>

        {error && (
          <div className="error-banner">
            <AlertCircle size={18} />
            <span>{error}</span>

            <button onClick={() => loadDashboard(true)}>
              Coba lagi
            </button>
          </div>
        )}

        <div
          key={isAdmin ? adminView : 'user-dashboard'}
          className="dashboard-view-transition"
        >
          {loading ? (
            <div className="dashboard-loading">
              <div className="loading-spinner" />
              <p>Memuat data dashboard...</p>
            </div>
          ) : isAdmin && adminView === 'users' ? (
            <AdminUsers
              currentUser={user}
              onBack={() => setAdminView('overview')}
              onActionSuccess={() => triggerMascotMood('happy')}
            />
          ) : isAdmin ? (
            <AdminDashboard
              stats={stats}
              tasks={tasks}
              pagination={pagination}
              page={page}
              setPage={setPage}
              totalPages={totalPages}
              onManageUsers={() => setAdminView('users')}
            />
          ) : (
            <UserDashboard
              stats={stats}
              tasks={tasks}
              categories={categories}
              pagination={pagination}
              page={page}
              setPage={setPage}
              totalPages={totalPages}
              searchInput={searchInput}
              setSearchInput={setSearchInput}
              statusFilter={statusFilter}
              setStatusFilter={setStatusFilter}
              priorityFilter={priorityFilter}
              setPriorityFilter={setPriorityFilter}
              handleSearch={handleSearch}
              clearFilters={clearFilters}
              openCreateModal={openCreateModal}
              openEditModal={openEditModal}
              handleDelete={handleDelete}
              openCreateCategoryModal={openCreateCategoryModal}
              openEditCategoryModal={openEditCategoryModal}
              handleDeleteCategory={handleDeleteCategory}
              triggerMascotMood={triggerMascotMood}
            />
          )}
        </div>

        {/* TASK MODAL */}

        {!isAdmin && modalOpen && (
          <div
            className="modal-backdrop"
            onMouseDown={(event) => {
              if (
                event.target === event.currentTarget &&
                !saving
              ) {
                setModalOpen(false)
                triggerMascotMood('idle')
              }
            }}
          >
            <section className="task-modal">
              <div className="modal-header">
                <div>
                  <span className="eyebrow">TASK MANAGEMENT</span>
                  <h2>
                    {editingTask ? 'Edit Task' : 'Create New Task'}
                  </h2>
                </div>

                <button
                  className="icon-button"
                  onClick={() => {
                    setModalOpen(false)
                    triggerMascotMood('idle')
                  }}
                  disabled={saving}
                >
                  <X size={20} />
                </button>
              </div>

              <form onSubmit={handleSubmit} className="task-form">
                <label>
                  Judul Task <span>*</span>

                  <input
                    value={form.title}
                    onFocus={() => triggerMascotMood('focus')}
                    onChange={(event) =>
                      updateForm('title', event.target.value)
                    }
                    placeholder="Contoh: Membuat halaman login"
                    minLength={2}
                    maxLength={200}
                    required
                  />
                </label>

                <label>
                  Deskripsi

                  <textarea
                    value={form.description}
                    onFocus={() => triggerMascotMood('focus')}
                    onChange={(event) =>
                      updateForm('description', event.target.value)
                    }
                    placeholder="additional details about the task..."
                    maxLength={5000}
                    rows={4}
                  />

                  <small>
                    {form.description.length}/5000 karakter
                  </small>
                </label>

                <label>
                  Category

                  <select
                    value={form.category_id}
                    onFocus={() => triggerMascotMood('focus')}
                    onChange={(event) =>
                      updateForm('category_id', event.target.value)
                    }
                  >
                    <option value="">No category</option>

                    {categories.map((category) => (
                      <option
                        key={category.id}
                        value={category.id}
                      >
                        {category.name}
                      </option>
                    ))}
                  </select>
                </label>

                <div className="form-row">
                  <label>
                    Prioritas

                    <select
                      value={form.priority}
                      onFocus={() => triggerMascotMood('focus')}
                      onChange={(event) =>
                        updateForm('priority', event.target.value)
                      }
                    >
                      <option value="LOW">Low</option>
                      <option value="MEDIUM">Medium</option>
                      <option value="HIGH">High</option>
                    </select>
                  </label>

                  <label>
                    Status

                    <select
                      value={form.status}
                      onFocus={() => triggerMascotMood('focus')}
                      onChange={(event) =>
                        updateForm('status', event.target.value)
                      }
                    >
                      <option value="TODO">To Do</option>
                      <option value="IN_PROGRESS">In Progress</option>
                      <option value="DONE">Done</option>
                    </select>
                  </label>
                </div>

                <label>
                  Tenggat Waktu

                  <input
                    type="date"
                    value={form.due_date}
                    onFocus={() => triggerMascotMood('focus')}
                    onChange={(event) =>
                      updateForm('due_date', event.target.value)
                    }
                  />
                </label>

                {formError && (
                  <div className="form-error">
                    <AlertCircle size={16} />
                    {formError}
                  </div>
                )}

                <div className="modal-actions">
                  <button
                    type="button"
                    className="secondary-button"
                    onClick={() => {
                      setModalOpen(false)
                      triggerMascotMood('idle')
                    }}
                    disabled={saving}
                  >
                    Batal
                  </button>

                  <button
                    type="submit"
                    className="primary-button"
                    disabled={saving}
                  >
                    {saving
                      ? 'Saving...'
                      : editingTask
                        ? 'Save Changes'
                        : 'Create Task'}
                  </button>
                </div>
              </form>
            </section>
          </div>
        )}

        {/* CATEGORY MODAL */}

        {!isAdmin && categoryModalOpen && (
          <div
            className="modal-backdrop"
            onMouseDown={(event) => {
              if (
                event.target === event.currentTarget &&
                !categorySaving
              ) {
                setCategoryModalOpen(false)
                triggerMascotMood('idle')
              }
            }}
          >
            <section className="task-modal category-modal">
              <div className="modal-header">
                <div>
                  <span className="eyebrow">CATEGORY MANAGEMENT</span>

                  <h2>
                    {editingCategory
                      ? 'Edit Category'
                      : 'New Category'}
                  </h2>
                </div>

                <button
                  className="icon-button"
                  onClick={() => {
                    setCategoryModalOpen(false)
                    triggerMascotMood('idle')
                  }}
                  disabled={categorySaving}
                >
                  <X size={20} />
                </button>
              </div>

              <form
                onSubmit={handleCategorySubmit}
                className="task-form"
              >
                <label>
                  Category Name <span>*</span>

                  <input
                    value={categoryName}
                    onFocus={() => triggerMascotMood('focus')}
                    onChange={(event) => {
                      setCategoryName(event.target.value)
                      triggerMascotMood('focus', 1800)
                    }}
                    placeholder="Ex: Homework, Work, Personal"
                    maxLength={100}
                    required
                    autoFocus
                  />

                  <small>{categoryName.length}/100 karakter</small>
                </label>

                {categoryError && (
                  <div className="form-error">
                    <AlertCircle size={16} />
                    {categoryError}
                  </div>
                )}

                <div className="modal-actions">
                  <button
                    type="button"
                    className="secondary-button"
                    onClick={() => {
                      setCategoryModalOpen(false)
                      triggerMascotMood('idle')
                    }}
                    disabled={categorySaving}
                  >
                    Batal
                  </button>

                  <button
                    type="submit"
                    className="primary-button"
                    disabled={categorySaving}
                  >
                    <Save size={16} />
                    {categorySaving
                      ? 'Saving...'
                      : editingCategory
                        ? 'Save Changes'
                        : 'Create Category'}
                  </button>
                </div>
              </form>
            </section>
          </div>
        )}
      </section>
    </main>
  )
}

/* =========================================
   USER DASHBOARD
========================================= */

function UserDashboard({
  stats,
  tasks,
  categories,
  pagination,
  page,
  setPage,
  totalPages,
  searchInput,
  setSearchInput,
  statusFilter,
  setStatusFilter,
  priorityFilter,
  setPriorityFilter,
  handleSearch,
  clearFilters,
  openCreateModal,
  openEditModal,
  handleDelete,
  openCreateCategoryModal,
  openEditCategoryModal,
  handleDeleteCategory,
  triggerMascotMood,
}) {
  const categoryMap = new Map(
    categories.map((category) => [
      String(category.id),
      category.name,
    ])
  )

  return (
    <>
      <div className="stats-grid">
        <StatCard
          icon={ClipboardList}
          label="Total Task"
          value={stats?.total_tasks}
          detail="All your tasks"
        />

        <StatCard
          icon={Clock3}
          label="To Do"
          value={stats?.todo}
          detail="Not yet completed"
        />

        <StatCard
          icon={Activity}
          label="In Progress"
          value={stats?.in_progress}
          detail="Sedang dikerjakan"
        />

        <StatCard
          icon={CheckCircle2}
          label="Done"
          value={stats?.done}
          detail="Successfully completed"
        />
      </div>

      {/* CATEGORY MANAGEMENT */}

      <section className="panel category-panel">
        <div className="panel-heading">
          <div>
            <h3>
              <Folder size={19} />
              My Categories
            </h3>

            <p>
              Group your tasks by creating custom categories.
            </p>
          </div>

          <button
            className="primary-button"
            onClick={openCreateCategoryModal}
          >
            <Plus size={18} />
            Add Category
          </button>
        </div>

        {categories.length === 0 ? (
          <div className="empty-state category-empty-state">
            <Tag size={32} />
            <strong>No categories available</strong>
            <span>
              Create your first category to start grouping tasks.
            </span>
          </div>
        ) : (
          <div className="category-list">
            {categories.map((category) => (
              <article className="category-item" key={category.id}>
                <div className="category-item-info">
                  <div className="category-icon">
                    <Folder size={18} />
                  </div>

                  <div>
                    <strong>{category.name}</strong>
                    <small>Category #{category.id}</small>
                  </div>
                </div>

                <div className="category-actions">
                  <button
                    className="table-icon-button edit"
                    title="Edit Category"
                    onClick={() => openEditCategoryModal(category)}
                  >
                    <Pencil size={16} />
                  </button>

                  <button
                    className="table-icon-button delete"
                    title="Delete Category"
                    onClick={() => handleDeleteCategory(category)}
                  >
                    <Trash2 size={16} />
                  </button>
                </div>
              </article>
            ))}
          </div>
        )}

        <p className="table-footnote">
          Total {categories.length} My Categories.
        </p>
      </section>

      {/* TASK LIST */}

      <section className="panel task-panel">
        <div className="panel-heading">
          <div>
            <h3>Task List</h3>
            <p>Manage and track all your tasks.</p>
          </div>

          <button
            className="primary-button"
            onClick={openCreateModal}
          >
            <Plus size={18} />
            Add Task
          </button>
        </div>

        <div className="filter-toolbar">
          <form
            className="search-form"
            onSubmit={handleSearch}
          >
            <Search size={18} />

            <input
              value={searchInput}
              onFocus={() => triggerMascotMood('focus')}
              onChange={(event) => {
                setSearchInput(event.target.value)
                triggerMascotMood('focus', 1800)
              }}
              placeholder="Cari judul task..."
            />

            <button type="submit">Cari</button>
          </form>

          <select
            value={statusFilter}
            onChange={(event) => {
              setStatusFilter(event.target.value)
              setPage(1)
              triggerMascotMood('focus')
            }}
          >
            <option value="">All Status</option>
            <option value="TODO">To Do</option>
            <option value="IN_PROGRESS">In Progress</option>
            <option value="DONE">Done</option>
          </select>

          <select
            value={priorityFilter}
            onChange={(event) => {
              setPriorityFilter(event.target.value)
              setPage(1)
              triggerMascotMood('focus')
            }}
          >
            <option value="">All Priorities</option>
            <option value="LOW">Low</option>
            <option value="MEDIUM">Medium</option>
            <option value="HIGH">High</option>
          </select>

          <button
            className="clear-filter-button"
            onClick={clearFilters}
          >
            Reset
          </button>
        </div>

        <div className="table-wrapper">
          <table className="task-table">
            <thead>
              <tr>
                <th>Task</th>
                <th>Category</th>
                <th>Priority</th>
                <th>Status</th>
                <th>Deadline</th>
                <th>Actions</th>
              </tr>
            </thead>

            <tbody>
              {tasks.length === 0 ? (
                <tr>
                  <td colSpan="6">
                    <div className="empty-state">
                      <ClipboardList size={35} />
                      <strong>Belum ada task</strong>
                      <span>
                        Buat task pertamamu untuk mulai mengatur pekerjaan.
                      </span>
                    </div>
                  </td>
                </tr>
              ) : (
                tasks.map((task) => (
                  <tr key={task.id}>
                    <td>
                      <div className="task-title-cell">
                        <strong>{task.title}</strong>
                        <span>
                          {task.description || 'No description available'}
                        </span>
                      </div>
                    </td>

                    <td>
                      {task.category_id
                        ? categoryMap.get(String(task.category_id)) || 'Category not available'
                        : '-'}
                    </td>

                    <td>
                      <PriorityBadge priority={task.priority} />
                    </td>

                    <td>
                      <StatusBadge status={task.status} />
                    </td>

                    <td>{formatDate(task.due_date)}</td>

                    <td>
                      <div className="action-buttons">
                        <button
                          className="table-icon-button edit"
                          title="Edit task"
                          onClick={() => openEditModal(task)}
                        >
                          <Pencil size={16} />
                        </button>

                        <button
                          className="table-icon-button delete"
                          title="Hapus task"
                          onClick={() => handleDelete(task)}
                        >
                          <Trash2 size={16} />
                        </button>
                      </div>
                    </td>
                  </tr>
                ))
              )}
            </tbody>
          </table>
        </div>

        <Pagination
          page={page}
          totalPages={totalPages}
          onPageChange={setPage}
        />

        <p className="table-footnote">
          Menampilkan {tasks.length} dari {pagination.total || 0} task
          {categories.length > 0
            ? ` · ${categories.length} Category Available`
            : ''}
        </p>
      </section>
    </>
  )
}

/* =========================================
   ADMIN DASHBOARD
========================================= */

function AdminDashboard({
  stats,
  tasks,
  pagination,
  page,
  setPage,
  totalPages,
  onManageUsers,
}) {
  const users = stats?.users || {}
  const taskStats = stats?.tasks || {}
  const statusStats = taskStats.by_status || {}
  const priorityStats = taskStats.by_priority || {}
  const categoryStats = stats?.categories || {}

  return (
    <>
      <div className="stats-grid admin-stats">
        <StatCard
          icon={Users}
          label="Total Users"
          value={users.total}
          detail={`${users.active || 0} active users`}
        />

        <StatCard
          icon={CheckCircle2}
          label="Active Users"
          value={users.active}
          detail="Active accounts"
        />

        <StatCard
          icon={AlertCircle}
          label="Nonactive Users"
          value={users.inactive}
          detail="Nonactive accounts"
        />

        <StatCard
          icon={ClipboardList}
          label="Total Task"
          value={taskStats.total}
          detail={`${categoryStats.total || 0} categories`}
        />
      </div>

      <div className="stats-grid admin-status-grid">
        <StatCard
          icon={Clock3}
          label="To Do"
          value={statusStats.todo}
        />

        <StatCard
          icon={Activity}
          label="In Progress"
          value={statusStats.in_progress}
        />

        <StatCard
          icon={CheckCircle2}
          label="Done"
          value={statusStats.done}
        />

        <StatCard
          icon={Folder}
          label="Category"
          value={categoryStats.total}
        />
      </div>

      <section className="panel priority-panel">
        <div className="panel-heading">
          <div>
            <h3>Priority Distribution</h3>
            <p>Number of tasks based on priority level.</p>
          </div>

          <ShieldCheck
            size={22}
            className="admin-shield"
          />
        </div>

        <div className="priority-overview">
          <div>
            <span className="priority-dot low-dot" />
            <span>Low</span>
            <strong>{priorityStats.low || 0}</strong>
          </div>

          <div>
            <span className="priority-dot medium-dot" />
            <span>Medium</span>
            <strong>{priorityStats.medium || 0}</strong>
          </div>

          <div>
            <span className="priority-dot high-dot" />
            <span>High</span>
            <strong>{priorityStats.high || 0}</strong>
          </div>
        </div>
      </section>

      <section className="panel task-panel">
        <div className="panel-heading">
          <div>
            <h3>All User Tasks</h3>
            <p>
              List of tasks from all registered accounts.
            </p>
          </div>

          <div className="admin-panel-actions">
            <button
              className="admin-manage-users-button"
              onClick={onManageUsers}
            >
              <Users size={16} />
              Manage Users
            </button>

            <span className="readonly-label">
              <ShieldCheck size={15} />
              Read Only
            </span>
          </div>
        </div>

        <div className="table-wrapper">
          <table className="task-table admin-task-table">
            <thead>
              <tr>
                <th>Task</th>
                <th>User ID</th>
                <th>Category ID</th>
                <th>Priority</th>
                <th>Status</th>
                <th>Created</th>
              </tr>
            </thead>

            <tbody>
              {tasks.length === 0 ? (
                <tr>
                  <td colSpan="6">
                    <div className="empty-state">
                      <ClipboardList size={35} />
                      <strong>No tasks available</strong>
                      <span>
                        User tasks will appear here.
                      </span>
                    </div>
                  </td>
                </tr>
              ) : (
                tasks.map((task) => (
                  <tr key={task.id}>
                    <td>
                      <div className="task-title-cell">
                        <strong>{task.title}</strong>
                        <span>
                          {task.description || 'No description available'}
                        </span>
                      </div>
                    </td>

                    <td>#{task.user_id}</td>
                    <td>{task.category_id ?? '-'}</td>

                    <td>
                      <PriorityBadge priority={task.priority} />
                    </td>

                    <td>
                      <StatusBadge status={task.status} />
                    </td>

                    <td>{formatDate(task.created_at)}</td>
                  </tr>
                ))
              )}
            </tbody>
          </table>
        </div>

        <Pagination
          page={page}
          totalPages={totalPages}
          onPageChange={setPage}
        />

        <p className="table-footnote">
          Showing {tasks.length} of {pagination.total || 0} tasks.
          Admin has read-only access.
        </p>
      </section>
    </>
  )
}
