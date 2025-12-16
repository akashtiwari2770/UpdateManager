import apiClient from './client';
import {
  Subscription,
  CreateSubscriptionRequest,
  UpdateSubscriptionRequest,
  PaginatedResponse,
  ListSubscriptionsQuery,
  SubscriptionStatistics,
} from '@/types';

export const subscriptionsApi = {
  getAll: async (
    customerId: string,
    query?: ListSubscriptionsQuery
  ): Promise<PaginatedResponse<Subscription>> => {
    try {
      const response = await apiClient.get<any>(
        `/customers/${customerId}/subscriptions`,
        { params: query }
      );
      if (response.data && response.data.subscriptions && response.data.pagination) {
        return {
          data: response.data.subscriptions,
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
      console.warn('Unexpected response format for getAll subscriptions:', response.data);
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
      console.error('Error fetching subscriptions:', error);
      throw error;
    }
  },

  getById: async (customerId: string, subscriptionId: string): Promise<Subscription> => {
    try {
      const response = await apiClient.get<any>(
        `/customers/${customerId}/subscriptions/${subscriptionId}`
      );
      if (response.data && response.data.data) {
        return response.data.data;
      }
      return response.data;
    } catch (error) {
      console.error('Error fetching subscription:', error);
      throw error;
    }
  },

  create: async (
    customerId: string,
    data: CreateSubscriptionRequest
  ): Promise<Subscription> => {
    try {
      const response = await apiClient.post<any>(
        `/customers/${customerId}/subscriptions`,
        data
      );
      if (response.data && response.data.data) {
        return response.data.data;
      }
      return response.data;
    } catch (error) {
      console.error('Error creating subscription:', error);
      throw error;
    }
  },

  update: async (
    customerId: string,
    subscriptionId: string,
    data: UpdateSubscriptionRequest
  ): Promise<Subscription> => {
    try {
      const response = await apiClient.put<any>(
        `/customers/${customerId}/subscriptions/${subscriptionId}`,
        data
      );
      if (response.data && response.data.data) {
        return response.data.data;
      }
      return response.data;
    } catch (error) {
      console.error('Error updating subscription:', error);
      throw error;
    }
  },

  delete: async (customerId: string, subscriptionId: string): Promise<void> => {
    try {
      await apiClient.delete(`/customers/${customerId}/subscriptions/${subscriptionId}`);
    } catch (error) {
      console.error('Error deleting subscription:', error);
      throw error;
    }
  },

  getStatistics: async (
    customerId: string,
    subscriptionId: string
  ): Promise<SubscriptionStatistics> => {
    try {
      const response = await apiClient.get<any>(
        `/customers/${customerId}/subscriptions/${subscriptionId}/statistics`
      );
      if (response.data && response.data.data) {
        return response.data.data;
      }
      return response.data;
    } catch (error) {
      console.error('Error fetching subscription statistics:', error);
      throw error;
    }
  },

  renew: async (
    customerId: string,
    subscriptionId: string,
    endDate: string
  ): Promise<Subscription> => {
    try {
      const response = await apiClient.post<any>(
        `/customers/${customerId}/subscriptions/${subscriptionId}/renew`,
        { end_date: endDate }
      );
      if (response.data && response.data.data) {
        return response.data.data;
      }
      return response.data;
    } catch (error) {
      console.error('Error renewing subscription:', error);
      throw error;
    }
  },
};

