import apiClient from './client';
import {
  License,
  CreateLicenseRequest,
  UpdateLicenseRequest,
  PaginatedResponse,
  ListLicensesQuery,
  LicenseStatistics,
} from '@/types';

export const licensesApi = {
  getAll: async (
    customerId: string,
    subscriptionId: string,
    query?: ListLicensesQuery
  ): Promise<PaginatedResponse<License>> => {
    try {
      const response = await apiClient.get<any>(
        `/customers/${customerId}/subscriptions/${subscriptionId}/licenses`,
        { params: query }
      );
      if (response.data && response.data.licenses && response.data.pagination) {
        return {
          data: response.data.licenses,
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
      console.warn('Unexpected response format for getAll licenses:', response.data);
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
      console.error('Error fetching licenses:', error);
      throw error;
    }
  },

  getById: async (
    customerId: string,
    subscriptionId: string,
    licenseId: string
  ): Promise<License> => {
    try {
      const response = await apiClient.get<any>(
        `/customers/${customerId}/subscriptions/${subscriptionId}/licenses/${licenseId}`
      );
      if (response.data && response.data.data) {
        return response.data.data;
      }
      return response.data;
    } catch (error) {
      console.error('Error fetching license:', error);
      throw error;
    }
  },

  assign: async (
    customerId: string,
    subscriptionId: string,
    data: CreateLicenseRequest
  ): Promise<License> => {
    try {
      const response = await apiClient.post<any>(
        `/customers/${customerId}/subscriptions/${subscriptionId}/licenses`,
        data
      );
      if (response.data && response.data.data) {
        return response.data.data;
      }
      return response.data;
    } catch (error) {
      console.error('Error assigning license:', error);
      throw error;
    }
  },

  update: async (
    customerId: string,
    subscriptionId: string,
    licenseId: string,
    data: UpdateLicenseRequest
  ): Promise<License> => {
    try {
      const response = await apiClient.put<any>(
        `/customers/${customerId}/subscriptions/${subscriptionId}/licenses/${licenseId}`,
        data
      );
      if (response.data && response.data.data) {
        return response.data.data;
      }
      return response.data;
    } catch (error) {
      console.error('Error updating license:', error);
      throw error;
    }
  },

  revoke: async (
    customerId: string,
    subscriptionId: string,
    licenseId: string
  ): Promise<void> => {
    try {
      await apiClient.delete(
        `/customers/${customerId}/subscriptions/${subscriptionId}/licenses/${licenseId}`
      );
    } catch (error) {
      console.error('Error revoking license:', error);
      throw error;
    }
  },

  getStatistics: async (
    customerId: string,
    subscriptionId: string,
    licenseId: string
  ): Promise<LicenseStatistics> => {
    try {
      const response = await apiClient.get<any>(
        `/customers/${customerId}/subscriptions/${subscriptionId}/licenses/${licenseId}/statistics`
      );
      if (response.data && response.data.data) {
        return response.data.data;
      }
      return response.data;
    } catch (error) {
      console.error('Error fetching license statistics:', error);
      throw error;
    }
  },

  renew: async (
    customerId: string,
    subscriptionId: string,
    licenseId: string,
    endDate: string
  ): Promise<License> => {
    try {
      const response = await apiClient.post<any>(
        `/customers/${customerId}/subscriptions/${subscriptionId}/licenses/${licenseId}/renew`,
        { end_date: endDate }
      );
      if (response.data && response.data.data) {
        return response.data.data;
      }
      return response.data;
    } catch (error) {
      console.error('Error renewing license:', error);
      throw error;
    }
  },
};

