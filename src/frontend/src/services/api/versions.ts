import apiClient from './client';
import {
  Version,
  CreateVersionRequest,
  UpdateVersionRequest,
  PaginatedResponse,
  ListVersionsQuery,
} from '@/types';

export const versionsApi = {
  getAll: async (query?: ListVersionsQuery): Promise<PaginatedResponse<Version>> => {
    try {
      const response = await apiClient.get<any>('/versions', { params: query });
      // Backend returns {success: true, data: [...], meta: {...}}
      // Interceptor transforms it to {data: [...], meta: {...}}
      // We need to transform meta to pagination
      // Handle case where data might be null (empty result)
      if (response.data && response.data.data !== undefined && response.data.meta) {
        const versionsData = response.data.data;
        return {
          data: Array.isArray(versionsData) ? versionsData : [],
          pagination: {
            page: response.data.meta.page || 1,
            limit: response.data.meta.limit || 25,
            total: response.data.meta.total || 0,
            total_pages: response.data.meta.total_pages || 1,
          },
        };
      }
      // Fallback if structure is different
      if (Array.isArray(response.data)) {
        return {
          data: response.data,
          pagination: {
            page: query?.page || 1,
            limit: query?.limit || 25,
            total: response.data.length,
            total_pages: 1,
          },
        };
      }
      // If response.data is already in the correct format
      if (response.data && response.data.data && response.data.pagination) {
        return response.data;
      }
      // Last resort - return empty
      console.warn('Unexpected response format for getAll versions:', response.data);
      return {
        data: [],
        pagination: {
          page: query?.page || 1,
          limit: query?.limit || 25,
          total: 0,
          total_pages: 1,
        },
      };
    } catch (error) {
      console.error('Error fetching versions:', error);
      throw error;
    }
  },

  getById: async (id: string): Promise<Version> => {
    const response = await apiClient.get<Version>(`/versions/${id}`);
    return response.data;
  },

  getByProduct: async (productId: string): Promise<Version[]> => {
    try {
      const response = await apiClient.get<any>(`/products/${productId}/versions`);
      // Backend returns paginated response: {success: true, data: [...], meta: {...}}
      // Interceptor transforms it to {data: [...], meta: {...}}
      // Handle both paginated and array responses
      if (response.data && response.data.data !== undefined && response.data.meta) {
        const versionsData = response.data.data;
        return Array.isArray(versionsData) ? versionsData : [];
      }
      // If it's already an array (unwrapped)
      if (Array.isArray(response.data)) {
        return response.data;
      }
      // If data is null or undefined, return empty array
      if (response.data === null || response.data === undefined) {
        return [];
      }
      // Last resort
      console.warn('Unexpected response format for getByProduct:', response.data);
      return [];
    } catch (error) {
      console.error('Error fetching versions by product:', error);
      return [];
    }
  },

  create: async (productId: string, data: CreateVersionRequest): Promise<Version> => {
    const response = await apiClient.post<Version>(`/products/${productId}/versions`, data);
    return response.data;
  },

  update: async (id: string, data: UpdateVersionRequest): Promise<Version> => {
    const response = await apiClient.put<Version>(`/versions/${id}`, data);
    return response.data;
  },

  submitForReview: async (id: string): Promise<Version> => {
    const response = await apiClient.post<Version>(`/versions/${id}/submit`);
    return response.data;
  },

  approve: async (id: string, approvedBy: string): Promise<Version> => {
    const response = await apiClient.post<Version>(`/versions/${id}/approve`, { approved_by: approvedBy });
    return response.data;
  },

  release: async (id: string): Promise<Version> => {
    const response = await apiClient.post<Version>(`/versions/${id}/release`);
    return response.data;
  },
};

