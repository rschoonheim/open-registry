export interface User {
  id: string
  username: string
  is_admin: boolean
}

export interface LoginCredentials {
  username: string
  password: string
}

export interface LoginResponse {
  user: User
  token: string
}
