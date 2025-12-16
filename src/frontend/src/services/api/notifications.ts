import apiClient from './client';
import { Notification, PaginatedResponse, ListNotificationsQuery } from '@/types';

export const notificationsApi = {
  getAll: async (query?: ListNotificationsQuery): Promise<PaginatedResponse<Notification>> => {
    const response = await apiClient.get<PaginatedResponse<Notification>>('/notifications', {
      params: query,
    });
    return response.data;
  },

  getById: async (id: string): Promise<Notification> => {
    const response = await apiClient.get<Notification>(`/notifications/${id}`);
    return response.data;
  },

  markAsRead: async (id: string): Promise<Notification> => {
    const response = await apiClient.put<Notification>(`/notifications/${id}/read`);
    return response.data;
  },

  markAllAsRead: async (): Promise<void> => {
    await apiClient.put('/notifications/read-all');
  },

  getUnreadCount: async (): Promise<number> => {
    const response = await apiClient.get<{ count: number }>('/notifications/unread-count');
    return response.data.count;
  },

  create: async (data: {
    recipient_id: string;
    type: string;
    title: string;
    message: string;
    priority: string;
    product_id?: string;
    version_id?: string;
  }): Promise<Notification> => {
    const response = await apiClient.post<Notification>('/notifications', data);
    return response.data;
  },
};

