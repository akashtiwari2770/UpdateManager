import apiClient from './client';
import { Product, CreateProductRequest, PaginatedResponse, ListProductsQuery } from '@/types';

export const productsApi = {
  getAll: async (query?: ListProductsQuery): Promise<PaginatedResponse<Product>> => {
    try {
      const response = await apiClient.get<any>('/products', { params: query });
      // Backend returns {success: true, data: [...], meta: {...}}
      // Interceptor transforms it to {data: [...], meta: {...}}
      // We need to transform meta to pagination
      if (response.data && response.data.data && response.data.meta) {
        return {
          data: response.data.data || [],
          pagination: {
            page: response.data.meta.page || 1,
            limit: response.data.meta.limit || 25,
            total: response.data.meta.total || 0,
            total_pages: response.data.meta.total_pages || 1,
          },
        };
      }
      // Fallback if structure is different
      if (Array.isArray(response.data)) {
        return {
          data: response.data,
          pagination: {
            page: query?.page || 1,
            limit: query?.limit || 25,
            total: response.data.length,
            total_pages: 1,
          },
        };
      }
      // If response.data is already in the correct format
      if (response.data && response.data.data && response.data.pagination) {
        return response.data;
      }
      // Last resort - return empty
      console.warn('Unexpected response format for getAll:', response.data);
      return {
        data: [],
        pagination: {
          page: query?.page || 1,
          limit: query?.limit || 25,
          total: 0,
          total_pages: 1,
        },
      };
    } catch (error) {
      console.error('Error fetching products:', error);
      throw error;
    }
  },

  getById: async (id: string): Promise<Product> => {
    const response = await apiClient.get<Product>(`/products/${id}`);
    return response.data;
  },

  create: async (data: CreateProductRequest): Promise<Product> => {
    try {
      const response = await apiClient.post<any>('/products', data);
      // Response interceptor should have unwrapped {success: true, data: {...}} to {...}
      // But handle both cases: unwrapped and wrapped
      let product: Product;
      if (response.data && response.data.data && response.data.success) {
        // Still wrapped - interceptor didn't work
        product = response.data.data;
      } else if (response.data && response.data.id) {
        // Already unwrapped
        product = response.data;
      } else {
        console.error('Unexpected response format:', response.data);
        throw new Error('Invalid response format from server');
      }
      
      if (!product.id) {
        console.error('Product missing ID:', product);
        throw new Error('Product ID is missing from server response');
      }
      
      return product;
    } catch (error) {
      console.error('Error creating product:', error);
      throw error;
    }
  },

  update: async (id: string, data: Partial<CreateProductRequest>): Promise<Product> => {
    const response = await apiClient.put<Product>(`/products/${id}`, data);
    return response.data;
  },

  delete: async (id: string): Promise<void> => {
    try {
      const response = await apiClient.delete(`/products/${id}`);
      // Backend returns {success: true, data: {message: "..."}}
      // Response interceptor unwraps it, so response.data should be {message: "..."}
      // For DELETE, we just need to verify the request succeeded (status 200)
      // The response.data might be the unwrapped message object or undefined
      // Either way, if we get here without an error, the delete was successful
      return;
    } catch (error) {
      console.error('Error deleting product:', error);
      throw error;
    }
  },

  getActive: async (): Promise<Product[]> => {
    const response = await apiClient.get<Product[]>('/products/active');
    return response.data;
  },
};

