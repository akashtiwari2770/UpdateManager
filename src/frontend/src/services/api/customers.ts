import apiClient from './client';
import {
  Customer,
  CreateCustomerRequest,
  UpdateCustomerRequest,
  PaginatedResponse,
  ListCustomersQuery,
  CustomerTenant,
  ListTenantsQuery,
  CustomerStatistics,
} from '@/types';

export const customersApi = {
  getAll: async (query?: ListCustomersQuery): Promise<PaginatedResponse<Customer>> => {
    try {
      const response = await apiClient.get<any>('/customers', { params: query });
      if (response.data && response.data.data && response.data.meta) {
        return {
          data: response.data.data || [],
          pagination: {
            page: response.data.meta.page || 1,
            limit: response.data.meta.limit || 20,
            total: response.data.meta.total || 0,
            total_pages: response.data.meta.total_pages || 1,
          },
        };
      }
      if (Array.isArray(response.data)) {
        return {
          data: response.data,
          pagination: {
            page: query?.page || 1,
            limit: query?.limit || 20,
            total: response.data.length,
            total_pages: 1,
          },
        };
      }
      if (response.data && response.data.data && response.data.pagination) {
        return response.data;
      }
      console.warn('Unexpected response format for getAll:', response.data);
      return {
        data: [],
        pagination: {
          page: query?.page || 1,
          limit: query?.limit || 20,
          total: 0,
          total_pages: 1,
        },
      };
    } catch (error) {
      console.error('Error fetching customers:', error);
      throw error;
    }
  },

  getById: async (id: string): Promise<Customer> => {
    try {
      const response = await apiClient.get<any>(`/customers/${id}`);
      if (response.data && response.data.data) {
        return response.data.data;
      }
      return response.data;
    } catch (error) {
      console.error('Error fetching customer:', error);
      throw error;
    }
  },

  create: async (data: CreateCustomerRequest): Promise<Customer> => {
    try {
      const response = await apiClient.post<any>('/customers', data);
      if (response.data && response.data.data) {
        return response.data.data;
      }
      return response.data;
    } catch (error) {
      console.error('Error creating customer:', error);
      throw error;
    }
  },

  update: async (id: string, data: UpdateCustomerRequest): Promise<Customer> => {
    try {
      const response = await apiClient.put<any>(`/customers/${id}`, data);
      if (response.data && response.data.data) {
        return response.data.data;
      }
      return response.data;
    } catch (error) {
      console.error('Error updating customer:', error);
      throw error;
    }
  },

  delete: async (id: string): Promise<void> => {
    try {
      await apiClient.delete(`/customers/${id}`);
    } catch (error) {
      console.error('Error deleting customer:', error);
      throw error;
    }
  },

  getTenants: async (
    customerId: string,
    query?: ListTenantsQuery
  ): Promise<PaginatedResponse<CustomerTenant>> => {
    try {
      const response = await apiClient.get<any>(`/customers/${customerId}/tenants`, {
        params: query,
      });
      if (response.data && response.data.data && response.data.meta) {
        return {
          data: response.data.data || [],
          pagination: {
            page: response.data.meta.page || 1,
            limit: response.data.meta.limit || 20,
            total: response.data.meta.total || 0,
            total_pages: response.data.meta.total_pages || 1,
          },
        };
      }
      if (Array.isArray(response.data)) {
        return {
          data: response.data,
          pagination: {
            page: query?.page || 1,
            limit: query?.limit || 20,
            total: response.data.length,
            total_pages: 1,
          },
        };
      }
      if (response.data && response.data.data && response.data.pagination) {
        return response.data;
      }
      console.warn('Unexpected response format for getTenants:', response.data);
      return {
        data: [],
        pagination: {
          page: query?.page || 1,
          limit: query?.limit || 20,
          total: 0,
          total_pages: 1,
        },
      };
    } catch (error) {
      console.error('Error fetching customer tenants:', error);
      throw error;
    }
  },

  getStatistics: async (customerId: string): Promise<CustomerStatistics> => {
    try {
      const response = await apiClient.get<any>(`/customers/${customerId}/statistics`);
      if (response.data && response.data.data) {
        return response.data.data;
      }
      return response.data;
    } catch (error) {
      console.error('Error fetching customer statistics:', error);
      throw error;
    }
  },
};

