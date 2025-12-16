import apiClient from './client';
import { CompatibilityMatrix, ValidateCompatibilityRequest } from '@/types';

export interface ListCompatibilityQuery {
  page?: number;
  limit?: number;
  product_id?: string;
  validation_status?: string;
}

export interface PaginatedCompatibilityResponse {
  data: CompatibilityMatrix[];
  meta?: {
    page: number;
    limit: number;
    total: number;
    total_pages: number;
  };
}

export const compatibilityApi = {
  validate: async (
    productId: string,
    versionNumber: string,
    data: ValidateCompatibilityRequest
  ): Promise<CompatibilityMatrix> => {
    const response = await apiClient.post<CompatibilityMatrix>(
      `/products/${productId}/versions/${versionNumber}/compatibility`,
      data
    );
    return response.data;
  },

  get: async (productId: string, versionNumber: string): Promise<CompatibilityMatrix> => {
    const response = await apiClient.get<CompatibilityMatrix>(
      `/products/${productId}/versions/${versionNumber}/compatibility`
    );
    return response.data;
  },

  list: async (query?: ListCompatibilityQuery): Promise<PaginatedCompatibilityResponse> => {
    const params = new URLSearchParams();
    if (query?.page) params.append('page', query.page.toString());
    if (query?.limit) params.append('limit', query.limit.toString());
    if (query?.product_id) params.append('product_id', query.product_id);
    if (query?.validation_status) params.append('validation_status', query.validation_status);

      const response = await apiClient.get<PaginatedCompatibilityResponse>(
        `/compatibility${params.toString() ? `?${params.toString()}` : ''}`
      );
      
      // Handle both wrapped and unwrapped responses
      const responseData = response.data as any;
      if (responseData && responseData.data && responseData.meta) {
        return {
          data: responseData.data,
          meta: responseData.meta,
        };
      }
      
      // If response is already in the correct format
      if (Array.isArray(responseData)) {
        return {
          data: responseData,
          meta: {
            page: query?.page || 1,
            limit: query?.limit || 25,
            total: responseData.length,
            total_pages: 1,
          },
        };
      }
      
      return {
        data: responseData?.data || [],
        meta: responseData?.meta || {
          page: query?.page || 1,
          limit: query?.limit || 25,
          total: 0,
          total_pages: 0,
        },
      };
  },
};

