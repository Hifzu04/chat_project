import React, { useState } from 'react'
import { Eye, EyeOff, Loader2 } from 'lucide-react'
import { Link } from 'react-router-dom'
import toast, { Toaster } from 'react-hot-toast'
import { useAuthStore } from '../Store/useAuthStore'

export default function Signup() {
  const [showPassword, setShowPassword] = useState(false)
  const [formData, setFormData] = useState({
    fullname: '',
    email: '',
    password: '',
  })
  const { signup, isSigningup } = useAuthStore()

  const validateForm = () => {
    if (!formData.fullname.trim())
      return toast.error('Full name is required')
    if (!formData.email.trim()) return toast.error('Email is required')
    if (!/\S+@\S+\.\S+/.test(formData.email))
      return toast.error('Invalid email address')
    if (!formData.password.trim())
      return toast.error('Password is required')
    if (formData.password.length < 6)
      return toast.error('Password must be at least 6 characters')
    return true
  }

  const handleSubmit = e => {
    e.preventDefault()
    if (validateForm()) signup(formData)
  }

  return (
    <div className="min-h-screen flex flex-col-reverse md:flex-row">
      {/* Form Panel */}
      <div className="w-full md:w-1/2 flex items-center justify-center py-12 px-4 sm:px-6 md:px-8 bg-base-200">
        <div className="w-full max-w-md space-y-6">
          {/* Logo */}
          <div className="flex justify-center">
            <img
              src="/Meetme.png"
              alt="App Logo"
              className="h-10 w-auto sm:h-12"
            />
          </div>

          {/* Heading */}
          <h2 className="text-center text-2xl sm:text-3xl md:text-4xl font-extrabold">
            Create your account
          </h2>

          {/* Form */}
          <form className="mt-8 space-y-4" onSubmit={handleSubmit}>
            {/* Full Name */}
            <div>
              <label htmlFor="fullname" className="block text-sm font-medium">
                Full Name
              </label>
              <input
                id="fullname"
                name="fullname"
                type="text"
                placeholder="Enter your name"
                required
                value={formData.fullname}
                onChange={e =>
                  setFormData({ ...formData, fullname: e.target.value })
                }
                className="mt-1 block w-full px-3 py-2 border rounded-lg shadow-sm placeholder-gray-400 focus:outline-none focus:ring-indigo-500 focus:border-indigo-500"
              />
            </div>

            {/* Email */}
            <div>
              <label htmlFor="email" className="block text-sm font-medium">
                Email address
              </label>
              <input
                id="email"
                name="email"
                type="email"
                placeholder="you@gmail.com"
                required
                value={formData.email}
                onChange={e =>
                  setFormData({ ...formData, email: e.target.value })
                }
                className="mt-1 block w-full px-3 py-2 border rounded-lg shadow-sm placeholder-gray-400 focus:outline-none focus:ring-indigo-500 focus:border-indigo-500"
              />
            </div>

            {/* Password */}
            <div className="relative">
              <label htmlFor="password" className="block text-sm font-medium">
                Password
              </label>
              <input
                id="password"
                name="password"
                type={showPassword ? 'text' : 'password'}
                placeholder="Enter your password"
                required
                value={formData.password}
                onChange={e =>
                  setFormData({ ...formData, password: e.target.value })
                }
                className="mt-1 block w-full px-3 py-2 border rounded-lg shadow-sm placeholder-gray-400 focus:outline-none focus:ring-indigo-500 focus:border-indigo-500"
              />
              <button
                type="button"
                onClick={() => setShowPassword(!showPassword)}
                className="absolute inset-y-0 translate-y-2  right-0 pr-3 flex items-center"
                aria-label="Toggle password visibility"
              >
                {showPassword ? (
                  <EyeOff className="h-5 w-5  text-base-content/40" />
                ) : (
                  <Eye className="h-5 w-5 text-base-content/40" />
                )}
              </button>
            </div>

            {/* Submit */}
            <button
              type="submit"
              disabled={isSigningup}
              className="w-full flex items-center justify-center py-2 px-4 border border-transparent rounded-lg shadow-sm bg-indigo-600 text-base-300 hover:bg-indigo-700 focus:outline-none focus:ring-2 focus:ring-indigo-500"
            >
              {isSigningup ? (
                <>
                  <Loader2 className="h-5 w-5 animate-spin" />
                  <span className="ml-2 hidden sm:inline">Loading...</span>
                </>
              ) : (
                'Create Account'
              )}
            </button>
          </form>

          {/* Link to login */}
          <p className="mt-4 text-center text-sm">
            Already registered?{' '}
            <Link
              to="/login"
              className="font-medium text-indigo-600 hover:text-indigo-500"
            >
              Go to Login
            </Link>
          </p>
        </div>
      </div>

      {/* Promo Panel */}
      <div className="hidden md:flex w-1/2 items-center justify-center bg-gradient-to-br from-indigo-500 to-purple-600 text-white p-8">
        <div className="max-w-lg text-center space-y-6">
          <h1 className="text-5xl font-bold">ChatterNest</h1>
          <p className="text-lg italic">
            “Pop in, say hi, and keep the convo rolling!” ☕💬
          </p>
        </div>
      </div>

      {/* Toast Notifications */}
      <Toaster position="top-center" reverseOrder={false} />
    </div>
  )
}
