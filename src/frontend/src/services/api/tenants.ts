import apiClient from './client';
import {
  CustomerTenant,
  CreateTenantRequest,
  UpdateTenantRequest,
  Deployment,
  PaginatedResponse,
  ListDeploymentsQuery,
  TenantStatistics,
} from '@/types';

export const tenantsApi = {
  create: async (
    customerId: string,
    data: CreateTenantRequest
  ): Promise<CustomerTenant> => {
    try {
      const response = await apiClient.post<any>(`/customers/${customerId}/tenants`, data);
      if (response.data && response.data.data) {
        return response.data.data;
      }
      return response.data;
    } catch (error) {
      console.error('Error creating tenant:', error);
      throw error;
    }
  },

  getById: async (customerId: string, tenantId: string): Promise<CustomerTenant> => {
    try {
      const response = await apiClient.get<any>(
        `/customers/${customerId}/tenants/${tenantId}`
      );
      if (response.data && response.data.data) {
        return response.data.data;
      }
      return response.data;
    } catch (error) {
      console.error('Error fetching tenant:', error);
      throw error;
    }
  },

  update: async (
    customerId: string,
    tenantId: string,
    data: UpdateTenantRequest
  ): Promise<CustomerTenant> => {
    try {
      const response = await apiClient.put<any>(
        `/customers/${customerId}/tenants/${tenantId}`,
        data
      );
      if (response.data && response.data.data) {
        return response.data.data;
      }
      return response.data;
    } catch (error) {
      console.error('Error updating tenant:', error);
      throw error;
    }
  },

  delete: async (customerId: string, tenantId: string): Promise<void> => {
    try {
      await apiClient.delete(`/customers/${customerId}/tenants/${tenantId}`);
    } catch (error) {
      console.error('Error deleting tenant:', error);
      throw error;
    }
  },

  getDeployments: async (
    customerId: string,
    tenantId: string,
    query?: ListDeploymentsQuery
  ): Promise<PaginatedResponse<Deployment>> => {
    try {
      const response = await apiClient.get<any>(
        `/customers/${customerId}/tenants/${tenantId}/deployments`,
        { params: query }
      );
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
      console.warn('Unexpected response format for getDeployments:', response.data);
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
      console.error('Error fetching tenant deployments:', error);
      throw error;
    }
  },

  getStatistics: async (
    customerId: string,
    tenantId: string
  ): Promise<TenantStatistics> => {
    try {
      const response = await apiClient.get<any>(
        `/customers/${customerId}/tenants/${tenantId}/statistics`
      );
      if (response.data && response.data.data) {
        return response.data.data;
      }
      return response.data;
    } catch (error) {
      console.error('Error fetching tenant statistics:', error);
      throw error;
    }
  },
};

