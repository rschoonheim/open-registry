import { defineEventHandler, readBody, createError, setCookie } from 'h3'
import type { LoginCredentials, LoginResponse } from '~/types/auth'


export default defineEventHandler(async (event) => {
  try {
    // Get request body
    const body = await readBody<LoginCredentials>(event)

    // Validate request
    if (!body.username || !body.password) {
      throw createError({
        statusCode: 400,
        statusMessage: 'Username and password are required'
      })
    }

    // Get API URL from runtime config
    const apiUrl = import.meta.env.VITE_API_URL || "http://localhost:3000"

    // Send login request to API
    const response = await fetch(`${apiUrl}/auth/login`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json'
      },
      body: JSON.stringify({
        username: body.username,
        password: body.password
      })
    })

    console.log('Login request sent to:', `${apiUrl}/auth/login`)

    // Check if response is OK
    if (!response.ok) {
      const errorData = await response.json()
      throw createError({
        statusCode: response.status,
        statusMessage: errorData.error || 'Authentication failed'
      })
    }

    // Parse response data
    const data = await response.json() as LoginResponse

    // Set JWT token as cookie
    setCookie(event, 'auth_token', data.token, {
      httpOnly: true,
      path: '/',
      maxAge: 60 * 60 * 24 * 7, // 7 days
      secure: process.env.NODE_ENV === 'production',
      sameSite: 'strict'
    })

    // Return user data without sensitive information
    return {
      user: {
        username: data.user.username,
        email: data.user.email,
        isAdmin: data.user.is_admin
      }
    }
  } catch (error) {
    console.error('Login error:', error)

    if (error.statusCode) {
      throw error // Re-throw H3 errors
    }

    throw createError({
      statusCode: 500,
      statusMessage: 'An unexpected error occurred during login'
    })
  }
})

