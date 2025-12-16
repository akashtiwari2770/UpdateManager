import apiClient from './client';
import {
  PendingUpdatesResponse,
  TenantPendingUpdatesSummary,
  CustomerPendingUpdatesSummary,
  PendingUpdatesQuery,
  PaginatedResponse,
} from '@/types';

export const pendingUpdatesApi = {
  getDeploymentPendingUpdates: async (
    customerId: string,
    tenantId: string,
    deploymentId: string
  ): Promise<PendingUpdatesResponse> => {
    const response = await apiClient.get<PendingUpdatesResponse>(
      `/customers/${customerId}/tenants/${tenantId}/deployments/${deploymentId}/updates`
    );
    return response.data;
  },

  getTenantPendingUpdates: async (
    customerId: string,
    tenantId: string,
    query?: PendingUpdatesQuery
  ): Promise<TenantPendingUpdatesSummary> => {
    const response = await apiClient.get<TenantPendingUpdatesSummary>(
      `/customers/${customerId}/tenants/${tenantId}/deployments/pending-updates`,
      { params: query }
    );
    return response.data;
  },

  getCustomerPendingUpdates: async (
    customerId: string,
    query?: PendingUpdatesQuery
  ): Promise<CustomerPendingUpdatesSummary> => {
    const response = await apiClient.get<CustomerPendingUpdatesSummary>(
      `/customers/${customerId}/deployments/pending-updates`,
      { params: query }
    );
    return response.data;
  },

  getAllPendingUpdates: async (
    query?: PendingUpdatesQuery
  ): Promise<PaginatedResponse<PendingUpdatesResponse>> => {
    try {
      const response = await apiClient.get<any>('/updates/pending', { params: query });
      // Backend returns {success: true, data: [...], meta: {...}}
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
      console.warn('Unexpected response format for getAllPendingUpdates:', response.data);
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
      console.error('Error fetching pending updates:', error);
      throw error;
    }
  },
};

