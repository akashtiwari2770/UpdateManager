import React, { useState, useEffect } from 'react';
import {
  StatisticsCard,
  RecentUpdates,
  PendingApprovals,
  ActivityTimeline,
} from '@/components/dashboard';
import { dashboardApi } from '@/services/api/dashboard';
import { Version } from '@/types';
import { Spinner, Alert } from '@/components/ui';

export const Dashboard: React.FC = () => {
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [stats, setStats] = useState({
    totalProducts: 0,
    activeVersions: 0,
    pendingUpdates: 0,
    activeRollouts: 0,
    totalCustomers: 0,
    totalTenants: 0,
    totalDeployments: 0,
    totalUsers: 0,
  });
  const [recentUpdates, setRecentUpdates] = useState<Version[]>([]);
  const [pendingApprovals, setPendingApprovals] = useState<Version[]>([]);
  const [recentActivity, setRecentActivity] = useState<any[]>([]);

  useEffect(() => {
    loadDashboardData();
  }, []);

  const loadDashboardData = async () => {
    setLoading(true);
    setError(null);
    try {
      const data = await dashboardApi.getAll();
      setStats(data.stats);
      setRecentUpdates(data.recentUpdates);
      setPendingApprovals(data.pendingApprovals);
      setRecentActivity(data.recentActivity);
    } catch (err: any) {
      setError(err.message || 'Failed to load dashboard data');
      console.error('Error loading dashboard:', err);
    } finally {
      setLoading(false);
    }
  };

  const handleApprove = (versionId: string) => {
    // Remove approved version from pending list
    setPendingApprovals((prev) => prev.filter((v) => v.id !== versionId));
    // Update stats
    setStats((prev) => ({
      ...prev,
      pendingUpdates: Math.max(0, prev.pendingUpdates - 1),
    }));
  };

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <h1 className="text-3xl font-bold text-gray-900">Dashboard</h1>
        <button
          onClick={loadDashboardData}
          className="text-sm text-blue-600 hover:text-blue-700 font-medium"
          disabled={loading}
        >
          Refresh
        </button>
      </div>

      {error && <Alert variant="error">{error}</Alert>}

      {/* Statistics Cards */}
      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6">
        <StatisticsCard
          title="Total Products"
          value={loading ? '...' : stats.totalProducts}
          icon={
            <svg className="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path
                strokeLinecap="round"
                strokeLinejoin="round"
                strokeWidth={2}
                d="M20 7l-8-4-8 4m16 0l-8 4m8-4v10l-8 4m0-10L4 7m8 4v10M4 7v10l8 4"
              />
            </svg>
          }
          linkTo="/products"
          loading={loading}
        />
        <StatisticsCard
          title="Active Versions"
          value={loading ? '...' : stats.activeVersions}
          icon={
            <svg className="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path
                strokeLinecap="round"
                strokeLinejoin="round"
                strokeWidth={2}
                d="M7 16a4 4 0 01-.88-7.903A5 5 0 1115.9 6L16 6a5 5 0 011 9.9M15 13l-3-3m0 0l-3 3m3-3v12"
              />
            </svg>
          }
          linkTo="/versions"
          loading={loading}
        />
        <StatisticsCard
          title="Pending Updates"
          value={loading ? '...' : stats.pendingUpdates}
          icon={
            <svg className="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path
                strokeLinecap="round"
                strokeLinejoin="round"
                strokeWidth={2}
                d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z"
              />
            </svg>
          }
          linkTo="/versions?state=pending_review"
          loading={loading}
        />
        <StatisticsCard
          title="Active Rollouts"
          value={loading ? '...' : stats.activeRollouts}
          icon={
            <svg className="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path
                strokeLinecap="round"
                strokeLinejoin="round"
                strokeWidth={2}
                d="M13 10V3L4 14h7v7l9-11h-7z"
              />
            </svg>
          }
          linkTo="/updates"
          loading={loading}
        />
      </div>

      {/* Customer Management Statistics */}
      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6">
        <StatisticsCard
          title="Total Customers"
          value={loading ? '...' : stats.totalCustomers}
          icon={
            <svg className="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path
                strokeLinecap="round"
                strokeLinejoin="round"
                strokeWidth={2}
                d="M17 20h5v-2a3 3 0 00-5.356-1.857M17 20H7m10 0v-2c0-.656-.126-1.283-.356-1.857M7 20H2v-2a3 3 0 015.356-1.857M7 20v-2c0-.656.126-1.283.356-1.857m0 0a5.002 5.002 0 019.288 0M15 7a3 3 0 11-6 0 3 3 0 016 0zm6 3a2 2 0 11-4 0 2 2 0 014 0zM7 10a2 2 0 11-4 0 2 2 0 014 0z"
              />
            </svg>
          }
          linkTo="/customers"
          loading={loading}
        />
        <StatisticsCard
          title="Total Tenants"
          value={loading ? '...' : stats.totalTenants}
          icon={
            <svg className="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path
                strokeLinecap="round"
                strokeLinejoin="round"
                strokeWidth={2}
                d="M19 21V5a2 2 0 00-2-2H7a2 2 0 00-2 2v16m14 0h2m-2 0h-5m-9 0H3m2 0h5M9 7h1m-1 4h1m4-4h1m-1 4h1m-5 10v-5a1 1 0 011-1h2a1 1 0 011 1v5m-4 0h4"
              />
            </svg>
          }
          loading={loading}
        />
        <StatisticsCard
          title="Total Deployments"
          value={loading ? '...' : stats.totalDeployments}
          icon={
            <svg className="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path
                strokeLinecap="round"
                strokeLinejoin="round"
                strokeWidth={2}
                d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z"
              />
            </svg>
          }
          loading={loading}
        />
        <StatisticsCard
          title="Total Users"
          value={loading ? '...' : stats.totalUsers}
          icon={
            <svg className="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path
                strokeLinecap="round"
                strokeLinejoin="round"
                strokeWidth={2}
                d="M12 4.354a4 4 0 110 5.292M15 21H3v-1a6 6 0 0112 0v1zm0 0h6v-1a6 6 0 00-9-5.197M13 7a4 4 0 11-8 0 4 4 0 018 0z"
              />
            </svg>
          }
          loading={loading}
        />
      </div>

      {/* Main Content Grid */}
      <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
        {/* Recent Updates */}
        <div className="lg:col-span-1">
          <RecentUpdates updates={recentUpdates} loading={loading} />
        </div>

        {/* Pending Approvals */}
        <div className="lg:col-span-1">
          <PendingApprovals
            approvals={pendingApprovals}
            loading={loading}
            onApprove={handleApprove}
          />
        </div>
      </div>

      {/* Activity Timeline */}
      <div>
        <ActivityTimeline activities={recentActivity} loading={loading} />
      </div>
    </div>
  );
};
