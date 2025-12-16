import apiClient from './client';
import { UpdateDetection } from '@/types';

export interface ListUpdateDetectionsQuery {
  page?: number;
  limit?: number;
  endpoint_id?: string;
  product_id?: string;
}

export interface PaginatedUpdateDetectionResponse {
  data: UpdateDetection[];
  meta?: {
    page: number;
    limit: number;
    total: number;
    total_pages: number;
  };
}

export const updateDetectionsApi = {
  list: async (query?: ListUpdateDetectionsQuery): Promise<PaginatedUpdateDetectionResponse> => {
    const params = new URLSearchParams();
    if (query?.page) params.append('page', query.page.toString());
    if (query?.limit) params.append('limit', query.limit.toString());
    if (query?.endpoint_id) params.append('endpoint_id', query.endpoint_id);
    if (query?.product_id) params.append('product_id', query.product_id);

    const response = await apiClient.get<PaginatedUpdateDetectionResponse>(
      `/update-detections${params.toString() ? `?${params.toString()}` : ''}`
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

  detect: async (data: Partial<UpdateDetection>): Promise<UpdateDetection> => {
    const response = await apiClient.post<UpdateDetection>('/update-detections', data);
    return response.data;
  },

  updateAvailableVersion: async (
    endpointId: string,
    productId: string,
    availableVersion: string
  ): Promise<UpdateDetection> => {
    const response = await apiClient.put<UpdateDetection>(
      `/update-detections/${endpointId}/${productId}/available-version`,
      { available_version: availableVersion }
    );
    return response.data;
  },
};

