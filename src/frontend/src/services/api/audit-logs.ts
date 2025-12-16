import apiClient from './client';
import { AuditLog, PaginatedResponse, ListAuditLogsQuery } from '@/types';

export const auditLogsApi = {
  getAll: async (query?: ListAuditLogsQuery): Promise<PaginatedResponse<AuditLog>> => {
    const response = await apiClient.get<PaginatedResponse<AuditLog>>('/audit-logs', {
      params: query,
    });
    return response.data;
  },

  getById: async (id: string): Promise<AuditLog> => {
    const response = await apiClient.get<AuditLog>(`/audit-logs/${id}`);
    return response.data;
  },
};

