import apiClient from './client';
import { UpdateRollout, RolloutStatus } from '@/types';

export interface ListUpdateRolloutsQuery {
  page?: number;
  limit?: number;
  endpoint_id?: string;
  product_id?: string;
  status?: RolloutStatus;
}

export interface PaginatedUpdateRolloutResponse {
  data: UpdateRollout[];
  meta?: {
    page: number;
    limit: number;
    total: number;
    total_pages: number;
  };
}

export interface InitiateRolloutRequest {
  endpoint_id: string;
  product_id: string;
  from_version: string;
  to_version: string;
}

export const updateRolloutsApi = {
  list: async (query?: ListUpdateRolloutsQuery): Promise<PaginatedUpdateRolloutResponse> => {
    const params = new URLSearchParams();
    if (query?.page) params.append('page', query.page.toString());
    if (query?.limit) params.append('limit', query.limit.toString());
    if (query?.endpoint_id) params.append('endpoint_id', query.endpoint_id);
    if (query?.product_id) params.append('product_id', query.product_id);
    if (query?.status) params.append('status', query.status);

    const response = await apiClient.get<PaginatedUpdateRolloutResponse>(
      `/update-rollouts${params.toString() ? `?${params.toString()}` : ''}`
    );
    
    const responseData = response.data as any;
    if (responseData && responseData.data && responseData.meta) {
      return {
        data: responseData.data,
        meta: responseData.meta,
      };
    }
    
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

  getById: async (id: string): Promise<UpdateRollout> => {
    const response = await apiClient.get<UpdateRollout>(`/update-rollouts/${id}`);
    return response.data;
  },

  initiate: async (data: InitiateRolloutRequest): Promise<UpdateRollout> => {
    const response = await apiClient.post<UpdateRollout>('/update-rollouts', data);
    return response.data;
  },

  updateProgress: async (id: string, progress: number): Promise<UpdateRollout> => {
    const response = await apiClient.put<UpdateRollout>(`/update-rollouts/${id}/progress`, {
      progress,
    });
    return response.data;
  },

  updateStatus: async (id: string, status: RolloutStatus, errorMessage?: string): Promise<UpdateRollout> => {
    const response = await apiClient.put<UpdateRollout>(`/update-rollouts/${id}/status`, {
      status,
      error_message: errorMessage,
    });
    return response.data;
  },
};

