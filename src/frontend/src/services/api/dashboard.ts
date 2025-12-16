import apiClient from './client';
import { productsApi } from './products';
import { versionsApi } from './versions';
import { updateRolloutsApi } from './update-rollouts';
import { auditLogsApi } from './audit-logs';
import { customersApi } from './customers';
import { Version, VersionState } from '@/types';

export interface DashboardStats {
  totalProducts: number;
  activeVersions: number;
  pendingUpdates: number;
  activeRollouts: number;
  totalCustomers: number;
  totalTenants: number;
  totalDeployments: number;
  totalUsers: number;
}

export interface DashboardData {
  stats: DashboardStats;
  recentUpdates: Version[];
  pendingApprovals: Version[];
  recentActivity: any[];
}

export const dashboardApi = {
  getStats: async (): Promise<DashboardStats> => {
    try {
      // Fetch all data in parallel
      const [
        productsResponse,
        versionsResponse,
        rolloutsResponse,
        customersResponse,
      ] = await Promise.all([
        productsApi.getAll({ limit: 1 }),
        versionsApi.getAll({ limit: 1 }),
        updateRolloutsApi.list({ limit: 1, status: 'in_progress' }),
        customersApi.getAll({ limit: 1 }),
      ]);

      // Get total counts from pagination
      const totalProducts = productsResponse.pagination?.total || 0;
      const totalVersions = versionsResponse.pagination?.total || 0;
      const totalCustomers = customersResponse.pagination?.total || 0;

      // Get active versions (released)
      const activeVersionsResponse = await versionsApi.getAll({
        state: VersionState.RELEASED,
        limit: 1,
      });
      const activeVersions = activeVersionsResponse.pagination?.total || 0;

      // Get pending updates (versions pending review)
      const pendingVersionsResponse = await versionsApi.getAll({
        state: VersionState.PENDING_REVIEW,
        limit: 1,
      });
      const pendingUpdates = pendingVersionsResponse.pagination?.total || 0;

      // Get active rollouts
      const activeRollouts = rolloutsResponse.meta?.total || rolloutsResponse.data?.length || 0;

      // Get customer management statistics
      let totalTenants = 0;
      let totalDeployments = 0;
      let totalUsers = 0;

      try {
        // Get all customers to calculate tenant and deployment stats
        const allCustomersResponse = await customersApi.getAll({ limit: 1000 });
        const customers = allCustomersResponse.data || [];

        // For each customer, get their statistics
        const customerStatsPromises = customers.map((customer) =>
          customersApi.getStatistics(customer.id).catch(() => null)
        );
        const customerStats = await Promise.all(customerStatsPromises);

        // Aggregate statistics
        customerStats.forEach((stats) => {
          if (stats) {
            totalTenants += stats.total_tenants || 0;
            totalDeployments += stats.total_deployments || 0;
            totalUsers += stats.total_users || 0;
          }
        });
      } catch (error) {
        console.error('Error fetching customer management stats:', error);
        // Continue with 0 values if customer stats fail
      }

      return {
        totalProducts,
        activeVersions,
        pendingUpdates,
        activeRollouts,
        totalCustomers,
        totalTenants,
        totalDeployments,
        totalUsers,
      };
    } catch (error) {
      console.error('Error fetching dashboard stats:', error);
      // Return default values on error
      return {
        totalProducts: 0,
        activeVersions: 0,
        pendingUpdates: 0,
        activeRollouts: 0,
        totalCustomers: 0,
        totalTenants: 0,
        totalDeployments: 0,
        totalUsers: 0,
      };
    }
  },

  getRecentUpdates: async (limit: number = 10): Promise<Version[]> => {
    try {
      const response = await versionsApi.getAll({
        state: VersionState.RELEASED,
        limit,
        page: 1,
      });
      return response.data || [];
    } catch (error) {
      console.error('Error fetching recent updates:', error);
      return [];
    }
  },

  getPendingApprovals: async (): Promise<Version[]> => {
    try {
      const response = await versionsApi.getAll({
        state: VersionState.PENDING_REVIEW,
        limit: 20,
        page: 1,
      });
      return response.data || [];
    } catch (error) {
      console.error('Error fetching pending approvals:', error);
      return [];
    }
  },

  getRecentActivity: async (limit: number = 20): Promise<any[]> => {
    try {
      const response = await auditLogsApi.getAll({
        limit,
        page: 1,
      });
      return response.data || [];
    } catch (error) {
      console.error('Error fetching recent activity:', error);
      return [];
    }
  },

  getAll: async (): Promise<DashboardData> => {
    const [stats, recentUpdates, pendingApprovals, recentActivity] = await Promise.all([
      dashboardApi.getStats(),
      dashboardApi.getRecentUpdates(10),
      dashboardApi.getPendingApprovals(),
      dashboardApi.getRecentActivity(20),
    ]);

    return {
      stats,
      recentUpdates,
      pendingApprovals,
      recentActivity,
    };
  },
};

