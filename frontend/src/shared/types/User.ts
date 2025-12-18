/**
 * Типы для работы с пользователями
 */

export interface User {
  id: number;
  uuid: string;
  login: string;
  email: string;
  name: string;
  phone?: string;
  role: 'buyer' | 'manager' | 'admin';
  created_at: string;
  updated_at: string;
}

export interface LoginRequest {
  login: string;
  password: string;
}

export interface RegisterRequest {
  login: string;
  email: string;
  name: string;
  password: string;
  phone?: string;
  role?: string;
}

export interface AuthResponse {
  access_token: string;
  refresh_token: string;
  expires_at: string;
  user: User;
}

export interface UpdateProfileRequest {
  name?: string;
  phone?: string;
  email?: string;
}

