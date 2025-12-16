import apiClient from './client';
import { UpgradePath, UpgradePathsQuery } from '@/types';

export const upgradePathsApi = {
  getAll: async (productId: string, query?: UpgradePathsQuery): Promise<UpgradePath[]> => {
    const response = await apiClient.get<UpgradePath[]>(
      `/products/${productId}/upgrade-paths`,
      { params: query }
    );
    return response.data;
  },

  get: async (
    productId: string,
    fromVersion: string,
    toVersion: string
  ): Promise<UpgradePath> => {
    const response = await apiClient.get<UpgradePath>(
      `/products/${productId}/upgrade-paths/${fromVersion}/${toVersion}`
    );
    return response.data;
  },

  create: async (productId: string, data: Partial<UpgradePath>): Promise<UpgradePath> => {
    const response = await apiClient.post<UpgradePath>(
      `/products/${productId}/upgrade-paths`,
      data
    );
    return response.data;
  },

  block: async (
    productId: string,
    fromVersion: string,
    toVersion: string,
    reason: string
  ): Promise<UpgradePath> => {
    const response = await apiClient.post<UpgradePath>(
      `/products/${productId}/upgrade-paths/${fromVersion}/${toVersion}/block`,
      { reason }
    );
    return response.data;
  },
};

