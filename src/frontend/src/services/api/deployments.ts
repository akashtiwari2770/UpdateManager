import apiClient from './client';
import {
  Deployment,
  CreateDeploymentRequest,
  UpdateDeploymentRequest,
  Version,
} from '@/types';

export const deploymentsApi = {
  create: async (
    customerId: string,
    tenantId: string,
    data: CreateDeploymentRequest
  ): Promise<Deployment> => {
    try {
      const response = await apiClient.post<any>(
        `/customers/${customerId}/tenants/${tenantId}/deployments`,
        data
      );
      if (response.data && response.data.data) {
        return response.data.data;
      }
      return response.data;
    } catch (error) {
      console.error('Error creating deployment:', error);
      throw error;
    }
  },

  getById: async (
    customerId: string,
    tenantId: string,
    deploymentId: string
  ): Promise<Deployment> => {
    try {
      const response = await apiClient.get<any>(
        `/customers/${customerId}/tenants/${tenantId}/deployments/${deploymentId}`
      );
      if (response.data && response.data.data) {
        return response.data.data;
      }
      return response.data;
    } catch (error) {
      console.error('Error fetching deployment:', error);
      throw error;
    }
  },

  update: async (
    customerId: string,
    tenantId: string,
    deploymentId: string,
    data: UpdateDeploymentRequest
  ): Promise<Deployment> => {
    try {
      const response = await apiClient.put<any>(
        `/customers/${customerId}/tenants/${tenantId}/deployments/${deploymentId}`,
        data
      );
      if (response.data && response.data.data) {
        return response.data.data;
      }
      return response.data;
    } catch (error) {
      console.error('Error updating deployment:', error);
      throw error;
    }
  },

  delete: async (
    customerId: string,
    tenantId: string,
    deploymentId: string
  ): Promise<void> => {
    try {
      await apiClient.delete(
        `/customers/${customerId}/tenants/${tenantId}/deployments/${deploymentId}`
      );
    } catch (error) {
      console.error('Error deleting deployment:', error);
      throw error;
    }
  },

  getAvailableUpdates: async (
    customerId: string,
    tenantId: string,
    deploymentId: string
  ): Promise<Version[]> => {
    try {
      const response = await apiClient.get<any>(
        `/customers/${customerId}/tenants/${tenantId}/deployments/${deploymentId}/updates`
      );
      if (response.data && response.data.data) {
        return response.data.data;
      }
      if (Array.isArray(response.data)) {
        return response.data;
      }
      return [];
    } catch (error) {
      console.error('Error fetching available updates:', error);
      throw error;
    }
  },
};

