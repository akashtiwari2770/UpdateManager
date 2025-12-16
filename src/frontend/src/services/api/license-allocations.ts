import apiClient from './client';
import {
  LicenseAllocation,
  AllocateLicenseRequest,
  PaginatedResponse,
  ListLicenseAllocationsQuery,
  LicenseUtilization,
} from '@/types';

export const licenseAllocationsApi = {
  getAll: async (
    customerId: string,
    subscriptionId: string,
    licenseId: string,
    query?: ListLicenseAllocationsQuery
  ): Promise<PaginatedResponse<LicenseAllocation>> => {
    try {
      const response = await apiClient.get<any>(
        `/customers/${customerId}/subscriptions/${subscriptionId}/licenses/${licenseId}/allocations`,
        { params: query }
      );
      if (response.data && response.data.allocations && response.data.pagination) {
        return {
          data: response.data.allocations,
          pagination: response.data.pagination,
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
      console.warn('Unexpected response format for getAll allocations:', response.data);
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
      console.error('Error fetching license allocations:', error);
      throw error;
    }
  },

  allocate: async (
    customerId: string,
    subscriptionId: string,
    licenseId: string,
    data: AllocateLicenseRequest
  ): Promise<LicenseAllocation> => {
    try {
      const response = await apiClient.post<any>(
        `/customers/${customerId}/subscriptions/${subscriptionId}/licenses/${licenseId}/allocate`,
        data
      );
      if (response.data && response.data.data) {
        return response.data.data;
      }
      return response.data;
    } catch (error) {
      console.error('Error allocating license:', error);
      throw error;
    }
  },

  release: async (
    customerId: string,
    subscriptionId: string,
    licenseId: string,
    allocationId: string
  ): Promise<void> => {
    try {
      await apiClient.post(
        `/customers/${customerId}/subscriptions/${subscriptionId}/licenses/${licenseId}/allocations/${allocationId}/release`
      );
    } catch (error) {
      console.error('Error releasing allocation:', error);
      throw error;
    }
  },

  getByTenant: async (
    customerId: string,
    tenantId: string,
    query?: ListLicenseAllocationsQuery
  ): Promise<PaginatedResponse<LicenseAllocation>> => {
    try {
      const response = await apiClient.get<any>(
        `/customers/${customerId}/tenants/${tenantId}/licenses`,
        { params: query }
      );
      if (response.data && response.data.allocations && response.data.pagination) {
        return {
          data: response.data.allocations,
          pagination: response.data.pagination,
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
      console.error('Error fetching tenant license allocations:', error);
      throw error;
    }
  },

  getByDeployment: async (
    customerId: string,
    tenantId: string,
    deploymentId: string,
    query?: ListLicenseAllocationsQuery
  ): Promise<PaginatedResponse<LicenseAllocation>> => {
    try {
      const response = await apiClient.get<any>(
        `/customers/${customerId}/tenants/${tenantId}/deployments/${deploymentId}/licenses`,
        { params: query }
      );
      if (response.data && response.data.allocations && response.data.pagination) {
        return {
          data: response.data.allocations,
          pagination: response.data.pagination,
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
      console.error('Error fetching deployment license allocations:', error);
      throw error;
    }
  },

  getUtilization: async (
    customerId: string,
    subscriptionId: string,
    licenseId: string
  ): Promise<LicenseUtilization> => {
    try {
      const response = await apiClient.get<any>(
        `/customers/${customerId}/subscriptions/${subscriptionId}/licenses/${licenseId}/utilization`
      );
      if (response.data && response.data.data) {
        return response.data.data;
      }
      return response.data;
    } catch (error) {
      console.error('Error fetching license utilization:', error);
      throw error;
    }
  },
};

