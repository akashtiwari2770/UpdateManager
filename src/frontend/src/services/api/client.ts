import axios, { AxiosInstance, AxiosError, InternalAxiosRequestConfig, AxiosResponse } from 'axios';
import { ErrorResponse } from '@/types';

// Use Vite proxy in development to avoid CORS issues
// In production, use the full API URL
const isDevelopment = import.meta.env.DEV;
const API_BASE_URL = isDevelopment 
  ? '' // Use relative path to leverage Vite proxy
  : (import.meta.env.VITE_API_BASE_URL || 'http://localhost:8080');
const API_VERSION = import.meta.env.VITE_API_VERSION || 'v1';

class ApiClient {
  private client: AxiosInstance;

  constructor() {
    this.client = axios.create({
      baseURL: `${API_BASE_URL}/api/${API_VERSION}`,
      headers: {
        'Content-Type': 'application/json',
      },
      timeout: 30000,
    });

    this.setupInterceptors();
  }

  private setupInterceptors(): void {
    // Request interceptor
    this.client.interceptors.request.use(
      (config: InternalAxiosRequestConfig) => {
        // Add auth token if available
        const token = localStorage.getItem('auth_token');
        if (token && config.headers) {
          config.headers.Authorization = `Bearer ${token}`;
        }

        // Add user ID header for audit logging
        const userId = localStorage.getItem('user_id') || 'anonymous';
        if (config.headers) {
          config.headers['X-User-ID'] = userId;
        }

        return config;
      },
      (error: AxiosError) => {
        return Promise.reject(error);
      }
    );

    // Response interceptor
    this.client.interceptors.response.use(
      (response: AxiosResponse) => {
        // Unwrap backend response format: {success: true, data: {...}, meta: {...}}
        if (response.data && typeof response.data === 'object' && 'success' in response.data) {
          const originalData = response.data;
          if (response.data.success && 'data' in response.data) {
            // Store original response for pagination metadata
            (response as any).originalData = originalData;
            // If there's meta (pagination), preserve it
            if (originalData.meta) {
              response.data = {
                data: originalData.data === null ? [] : originalData.data,
                meta: originalData.meta,
              };
            } else {
              // Return unwrapped data for non-paginated responses
              // Convert null to empty array for list responses
              response.data = originalData.data === null ? [] : originalData.data;
            }
          } else if (!response.data.success && 'error' in response.data) {
            // This shouldn't happen in success responses, but handle it
            const errorData = response.data.error;
            response.data = {
              error: {
                code: errorData.code || 'UNKNOWN_ERROR',
                message: errorData.message || 'An error occurred',
                details: errorData.details,
              },
            };
          }
        }
        return response;
      },
      (error: AxiosError<ErrorResponse>) => {
        // Handle common errors
        if (error.response) {
          const status = error.response.status;
          let errorData = error.response.data;

          // Unwrap backend error response format: {success: false, error: {...}}
          if (errorData && typeof errorData === 'object' && 'success' in errorData && !errorData.success) {
            if (errorData.error) {
              // Transform to match ErrorResponse type
              errorData = {
                error: {
                  code: errorData.error.code || 'UNKNOWN_ERROR',
                  message: errorData.error.message || 'An error occurred',
                  details: errorData.error.details,
                },
              };
              // Update the error response data
              error.response.data = errorData;
            }
          }

          switch (status) {
            case 401:
              // Unauthorized - clear auth and redirect to login
              localStorage.removeItem('auth_token');
              localStorage.removeItem('user_id');
              // window.location.href = '/login';
              break;
            case 403:
              // Forbidden
              console.error('Access forbidden:', errorData);
              break;
            case 404:
              // Not found
              console.error('Resource not found:', errorData);
              break;
            case 500:
              // Server error
              console.error('Server error:', errorData);
              break;
            default:
              console.error('API error:', errorData);
          }
        } else if (error.request) {
          // Network error
          console.error('Network error:', error.message);
        } else {
          // Request setup error
          console.error('Request error:', error.message);
        }

        return Promise.reject(error);
      }
    );
  }

  public getInstance(): AxiosInstance {
    return this.client;
  }

  public get<T = any>(url: string, config?: any): Promise<AxiosResponse<T>> {
    return this.client.get<T>(url, config);
  }

  public post<T = any>(url: string, data?: any, config?: any): Promise<AxiosResponse<T>> {
    return this.client.post<T>(url, data, config);
  }

  public put<T = any>(url: string, data?: any, config?: any): Promise<AxiosResponse<T>> {
    return this.client.put<T>(url, data, config);
  }

  public patch<T = any>(url: string, data?: any, config?: any): Promise<AxiosResponse<T>> {
    return this.client.patch<T>(url, data, config);
  }

  public delete<T = any>(url: string, config?: any): Promise<AxiosResponse<T>> {
    return this.client.delete<T>(url, config);
  }
}

export const apiClient = new ApiClient();
export default apiClient;

