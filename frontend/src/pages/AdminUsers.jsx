import { useCallback, useEffect, useState } from 'react'
import {
  ArrowLeft,
  Search,
  RefreshCw,
  Users,
  ShieldCheck,
  UserCheck,
  UserX,
  Eye,
  X,
  AlertCircle,
} from 'lucide-react'
import { api } from '../lib/api'
import './AdminUsers.css'

const LIMIT = 10

function getData(response) {
  return response?.data ?? response ?? {}
}

function formatDate(value) {
  if (!value) return '-'
  const date = new Date(value)
  if (Number.isNaN(date.getTime())) return '-'

  return date.toLocaleDateString('id-ID', {
    day: '2-digit',
    month: 'long',
    year: 'numeric',
  })
}

export default function AdminUsers({ currentUser, onBack }) {
  const [users, setUsers] = useState([])
  const [pagination, setPagination] = useState({
    total: 0,
    page: 1,
    total_pages: 1,
  })

  const [page, setPage] = useState(1)
  const [searchInput, setSearchInput] = useState('')
  const [search, setSearch] = useState('')
  const [roleFilter, setRoleFilter] = useState('')
  const [statusFilter, setStatusFilter] = useState('')

  const [loading, setLoading] = useState(true)
  const [refreshing, setRefreshing] = useState(false)
  const [error, setError] = useState('')
  const [notice, setNotice] = useState('')

  const [selectedUser, setSelectedUser] = useState(null)
  const [detailLoading, setDetailLoading] = useState(false)
  const [detailError, setDetailError] = useState('')

  const [updatingID, setUpdatingID] = useState(null)

  const loadUsers = useCallback(async (refresh = false) => {
    if (refresh) setRefreshing(true)
    else setLoading(true)

    setError('')
    setNotice('')

    try {
      const params = new URLSearchParams({
        page: String(page),
        limit: String(LIMIT),
      })

      if (search) params.set('search', search)
      if (roleFilter) params.set('role', roleFilter)
      if (statusFilter !== '') params.set('is_active', statusFilter)

      const response = await api.get(
        `/api/v1/admin/users?${params.toString()}`
      )

      const data = getData(response)
      const resultUsers = data.users || data.items || []

      setUsers(Array.isArray(resultUsers) ? resultUsers : [])

      const paginationData = data.pagination || data

      setPagination({
        total: paginationData.total ?? 0,
        page: paginationData.page ?? page,
        total_pages: paginationData.total_pages ?? 1,
      })
    } catch (err) {
      setError(err.message || 'Gagal mengambil daftar pengguna.')
    } finally {
      setLoading(false)
      setRefreshing(false)
    }
  }, [page, search, roleFilter, statusFilter])

  useEffect(() => {
    loadUsers()
  }, [loadUsers])

  function submitSearch(event) {
    event.preventDefault()
    setPage(1)
    setSearch(searchInput.trim())
  }

  function resetFilters() {
    setSearchInput('')
    setSearch('')
    setRoleFilter('')
    setStatusFilter('')
    setPage(1)
  }

  async function openUserDetail(userID) {
    setSelectedUser(null)
    setDetailError('')
    setDetailLoading(true)

    try {
      const response = await api.get(`/api/v1/admin/users/${userID}`)
      setSelectedUser(getData(response))
    } catch (err) {
      setDetailError(err.message || 'Gagal mengambil detail pengguna.')
    } finally {
      setDetailLoading(false)
    }
  }

  async function updateRole(user, newRole) {
    if (newRole === user.role) return

    const confirmed = window.confirm(
      `Ubah role ${user.name} dari ${user.role} menjadi ${newRole}?`
    )

    if (!confirmed) return

    setUpdatingID(user.id)
    setError('')
    setNotice('')

    try {
      await api.put(`/api/v1/admin/users/${user.id}/role`, {
        role: newRole,
      })

      setNotice(`Role ${user.name} berhasil diubah menjadi ${newRole}.`)
      await loadUsers(true)

      if (selectedUser?.id === user.id) {
        await openUserDetail(user.id)
      }
    } catch (err) {
      setError(err.message || 'Gagal mengubah role pengguna.')
    } finally {
      setUpdatingID(null)
    }
  }

  async function updateStatus(user) {
    const nextStatus = !user.is_active
    const action = nextStatus ? 'mengaktifkan' : 'menonaktifkan'

    const confirmed = window.confirm(
      `Yakin ingin ${action} akun ${user.name}?`
    )

    if (!confirmed) return

    setUpdatingID(user.id)
    setError('')
    setNotice('')

    try {
      await api.put(`/api/v1/admin/users/${user.id}/status`, {
        is_active: nextStatus,
      })

      setNotice(
        `Akun ${user.name} berhasil ${nextStatus ? 'diaktifkan' : 'dinonaktifkan'}.`
      )

      await loadUsers(true)

      if (selectedUser?.id === user.id) {
        await openUserDetail(user.id)
      }
    } catch (err) {
      setError(err.message || 'Gagal mengubah status akun.')
    } finally {
      setUpdatingID(null)
    }
  }

  const totalPages = Math.max(1, pagination.total_pages || 1)
  const activeCount = users.filter((user) => user.is_active).length
  const inactiveCount = users.filter((user) => !user.is_active).length
  const adminCount = users.filter((user) => user.role === 'ADMIN').length

  const detailUser = selectedUser?.user || selectedUser

  return (
    <section className="admin-users-page">
      <div className="admin-users-heading">
        <div>
          <button className="admin-back-button" onClick={onBack}>
            <ArrowLeft size={17} />
            Kembali ke Dashboard
          </button>

          <span className="admin-eyebrow">ADMINISTRATION</span>
          <h2>Kelola Pengguna</h2>
          <p>Kelola role dan status akun seluruh pengguna sistem.</p>
        </div>

        <button
          className="admin-refresh-button"
          onClick={() => loadUsers(true)}
          disabled={loading || refreshing}
        >
          <RefreshCw size={16} className={refreshing ? 'admin-spinning' : ''} />
          Refresh
        </button>
      </div>

      <div className="admin-user-summary">
        <article>
          <div className="admin-summary-icon"><Users size={19} /></div>
          <span>Total pengguna</span>
          <strong>{pagination.total}</strong>
        </article>

        <article>
          <div className="admin-summary-icon"><UserCheck size={19} /></div>
          <span>Aktif di halaman ini</span>
          <strong>{activeCount}</strong>
        </article>

        <article>
          <div className="admin-summary-icon"><UserX size={19} /></div>
          <span>Nonaktif di halaman ini</span>
          <strong>{inactiveCount}</strong>
        </article>

        <article>
          <div className="admin-summary-icon"><ShieldCheck size={19} /></div>
          <span>Admin di halaman ini</span>
          <strong>{adminCount}</strong>
        </article>
      </div>

      {error && (
        <div className="admin-message admin-message-error">
          <AlertCircle size={17} />
          {error}
          <button onClick={() => setError('')}><X size={16} /></button>
        </div>
      )}

      {notice && (
        <div className="admin-message admin-message-success">
          {notice}
          <button onClick={() => setNotice('')}><X size={16} /></button>
        </div>
      )}

      <section className="admin-users-panel">
        <div className="admin-users-panel-heading">
          <div>
            <h3>Daftar Pengguna</h3>
            <p>Cari akun, filter role, dan atur status pengguna.</p>
          </div>
          <span className="admin-total-label">
            {pagination.total} pengguna
          </span>
        </div>

        <div className="admin-user-filters">
          <form className="admin-user-search" onSubmit={submitSearch}>
            <Search size={17} />
            <input
              value={searchInput}
              onChange={(event) => setSearchInput(event.target.value)}
              placeholder="Cari nama atau email..."
            />
            <button type="submit">Cari</button>
          </form>

          <select
            value={roleFilter}
            onChange={(event) => {
              setPage(1)
              setRoleFilter(event.target.value)
            }}
          >
            <option value="">Semua Role</option>
            <option value="USER">USER</option>
            <option value="ADMIN">ADMIN</option>
          </select>

          <select
            value={statusFilter}
            onChange={(event) => {
              setPage(1)
              setStatusFilter(event.target.value)
            }}
          >
            <option value="">Semua Status</option>
            <option value="true">Aktif</option>
            <option value="false">Nonaktif</option>
          </select>

          <button className="admin-reset-button" onClick={resetFilters}>
            Reset
          </button>
        </div>

        <div className="admin-users-table-wrapper">
          <table className="admin-users-table">
            <thead>
              <tr>
                <th>Pengguna</th>
                <th>ID</th>
                <th>Role</th>
                <th>Status Akun</th>
                <th>Bergabung</th>
                <th>Kelola</th>
              </tr>
            </thead>

            <tbody>
              {loading ? (
                <tr>
                  <td colSpan="6" className="admin-users-empty">
                    Memuat daftar pengguna...
                  </td>
                </tr>
              ) : users.length === 0 ? (
                <tr>
                  <td colSpan="6" className="admin-users-empty">
                    Tidak ada pengguna yang cocok dengan filter.
                  </td>
                </tr>
              ) : (
                users.map((user) => {
                  const isSelf = Number(user.id) === Number(currentUser?.id)
                  const isUpdating = updatingID === user.id

                  return (
                    <tr key={user.id}>
                      <td>
                        <div className="admin-user-identity">
                          <div className="admin-user-avatar">
                            {user.name?.charAt(0)?.toUpperCase() || 'U'}
                          </div>
                          <div>
                            <strong>{user.name}</strong>
                            <span>{user.email}</span>
                            {isSelf && <small>Akun kamu</small>}
                          </div>
                        </div>
                      </td>

                      <td>#{user.id}</td>

                      <td>
                        <select
                          className="admin-role-select"
                          value={user.role}
                          disabled={isSelf || isUpdating}
                          onChange={(event) => updateRole(user, event.target.value)}
                        >
                          <option value="USER">USER</option>
                          <option value="ADMIN">ADMIN</option>
                        </select>
                      </td>

                      <td>
                        <span className={
                          user.is_active
                            ? 'admin-account-status active'
                            : 'admin-account-status inactive'
                        }>
                          {user.is_active ? 'Aktif' : 'Nonaktif'}
                        </span>
                      </td>

                      <td>{formatDate(user.created_at)}</td>

                      <td>
                        <div className="admin-user-actions">
                          <button
                            className="admin-detail-button"
                            title="Lihat detail"
                            onClick={() => openUserDetail(user.id)}
                          >
                            <Eye size={16} />
                          </button>

                          <button
                            className={
                              user.is_active
                                ? 'admin-status-button deactivate'
                                : 'admin-status-button activate'
                            }
                            disabled={isSelf || isUpdating}
                            onClick={() => updateStatus(user)}
                          >
                            {isUpdating
                              ? '...'
                              : user.is_active
                                ? 'Nonaktifkan'
                                : 'Aktifkan'}
                          </button>
                        </div>
                      </td>
                    </tr>
                  )
                })
              )}
            </tbody>
          </table>
        </div>

        <div className="admin-users-pagination">
          <span>
            Halaman {page} dari {totalPages}
          </span>

          <div>
            <button
              disabled={page <= 1 || loading}
              onClick={() => setPage((value) => value - 1)}
            >
              Sebelumnya
            </button>
            <button
              disabled={page >= totalPages || loading}
              onClick={() => setPage((value) => value + 1)}
            >
              Berikutnya
            </button>
          </div>
        </div>
      </section>

      {(detailLoading || selectedUser || detailError) && (
        <div
          className="admin-detail-backdrop"
          onMouseDown={(event) => {
            if (event.target === event.currentTarget) {
              setSelectedUser(null)
              setDetailError('')
            }
          }}
        >
          <section className="admin-detail-modal">
            <div className="admin-detail-heading">
              <div>
                <span className="admin-eyebrow">USER INFORMATION</span>
                <h3>Detail Pengguna</h3>
              </div>
              <button
                className="admin-detail-close"
                onClick={() => {
                  setSelectedUser(null)
                  setDetailError('')
                }}
              >
                <X size={19} />
              </button>
            </div>

            {detailLoading ? (
              <p className="admin-detail-loading">Memuat detail pengguna...</p>
            ) : detailError ? (
              <div className="admin-message admin-message-error">
                {detailError}
              </div>
            ) : detailUser ? (
              <div className="admin-detail-content">
                <div className="admin-detail-avatar">
                  {detailUser.name?.charAt(0)?.toUpperCase() || 'U'}
                </div>

                <h4>{detailUser.name || '-'}</h4>
                <p>{detailUser.email || '-'}</p>

                <div className="admin-detail-grid">
                  <div>
                    <span>User ID</span>
                    <strong>#{detailUser.id}</strong>
                  </div>
                  <div>
                    <span>Role</span>
                    <strong>{detailUser.role || '-'}</strong>
                  </div>
                  <div>
                    <span>Status</span>
                    <strong>{detailUser.is_active ? 'Aktif' : 'Nonaktif'}</strong>
                  </div>
                  <div>
                    <span>Bergabung</span>
                    <strong>{formatDate(detailUser.created_at)}</strong>
                  </div>
                  <div>
                    <span>Total Task</span>
                    <strong>
                      {detailUser.total_tasks ??
                        detailUser.task_count ??
                        detailUser.tasks_count ??
                        '-'}
                    </strong>
                  </div>
                </div>
              </div>
            ) : null}
          </section>
        </div>
      )}
    </section>
  )
}