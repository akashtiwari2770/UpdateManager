import React, { useState, useEffect } from 'react';
import { useParams, useNavigate } from 'react-router-dom';
import { tenantsApi } from '@/services/api/tenants';
import {
  CustomerTenant,
  TenantStatus,
  TenantStatistics,
} from '@/types';
import { Button, Card, Badge, Spinner } from '@/components/ui';
import { DeploymentsList } from '@/components/deployments';
import { pendingUpdatesApi } from '@/services/api/pending-updates';
import { TenantPendingUpdatesSummary } from '@/types';
import { UpdateBadge } from '@/components/ui/UpdateBadge';
import { PriorityBadge } from '@/components/ui/PriorityBadge';

type TabType = 'overview' | 'deployments';

export const TenantDetails: React.FC = () => {
  const { customerId, tenantId } = useParams<{ customerId: string; tenantId: string }>();
  const navigate = useNavigate();
  const [tenant, setTenant] = useState<CustomerTenant | null>(null);
  const [statistics, setStatistics] = useState<TenantStatistics | null>(null);
  const [pendingUpdatesSummary, setPendingUpdatesSummary] = useState<TenantPendingUpdatesSummary | null>(null);
  const [loading, setLoading] = useState(true);
  const [loadingStats, setLoadingStats] = useState(false);
  const [loadingPendingUpdates, setLoadingPendingUpdates] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [activeTab, setActiveTab] = useState<TabType>('overview');

  useEffect(() => {
    if (customerId && tenantId) {
      loadTenant();
      loadStatistics();
      loadPendingUpdates();
    }
  }, [customerId, tenantId]);

  const loadTenant = async () => {
    if (!customerId || !tenantId) return;
    try {
      setLoading(true);
      setError(null);
      const data = await tenantsApi.getById(customerId, tenantId);
      setTenant(data);
    } catch (err: any) {
      setError(err.response?.data?.error?.message || 'Failed to load tenant');
      console.error('Error loading tenant:', err);
    } finally {
      setLoading(false);
    }
  };

  const loadStatistics = async () => {
    if (!customerId || !tenantId) return;
    try {
      setLoadingStats(true);
      const stats = await tenantsApi.getStatistics(customerId, tenantId);
      setStatistics(stats);
    } catch (err: any) {
      console.error('Error loading statistics:', err);
    } finally {
      setLoadingStats(false);
    }
  };

  const loadPendingUpdates = async () => {
    if (!customerId || !tenantId) return;
    try {
      setLoadingPendingUpdates(true);
      const summary = await pendingUpdatesApi.getTenantPendingUpdates(customerId, tenantId);
      setPendingUpdatesSummary(summary);
    } catch (err: any) {
      console.error('Error loading pending updates:', err);
    } finally {
      setLoadingPendingUpdates(false);
    }
  };

  const getStatusBadge = (status: TenantStatus) => {
    const statusConfig = {
      [TenantStatus.ACTIVE]: { label: 'Active', className: 'bg-green-100 text-green-800' },
      [TenantStatus.INACTIVE]: { label: 'Inactive', className: 'bg-gray-100 text-gray-800' },
    };
    const config = statusConfig[status] || statusConfig[TenantStatus.INACTIVE];
    return <Badge className={config.className}>{config.label}</Badge>;
  };

  const formatDate = (dateString: string) => {
    return new Date(dateString).toLocaleString('en-US', {
      year: 'numeric',
      month: 'long',
      day: 'numeric',
      hour: '2-digit',
      minute: '2-digit',
    });
  };

  if (loading) {
    return (
      <div className="flex items-center justify-center h-64">
        <Spinner size="lg" />
      </div>
    );
  }

  if (error || !tenant || !customerId) {
    return (
      <div className="space-y-6">
        <div className="flex items-center justify-between">
          <h1 className="text-3xl font-bold text-gray-900">Tenant Details</h1>
          <Button variant="ghost" onClick={() => navigate(`/customers/${customerId}`)}>
            Back to Customer
          </Button>
        </div>
        <Card>
          <div className="text-center py-12">
            <p className="text-red-600 mb-4">{error || 'Tenant not found'}</p>
            <Button variant="primary" onClick={() => navigate(`/customers/${customerId}`)}>
              Back to Customer
            </Button>
          </div>
        </Card>
      </div>
    );
  }

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-3xl font-bold text-gray-900">{tenant.name}</h1>
          <p className="text-gray-500 mt-1">Tenant ID: {tenant.tenant_id}</p>
        </div>
        <div className="flex space-x-2">
          <Button variant="ghost" onClick={() => navigate(`/customers/${customerId}`)}>
            Back
          </Button>
          <Button
            variant="primary"
            onClick={() => navigate(`/customers/${customerId}/tenants/${tenant.id}/edit`)}
          >
            Edit
          </Button>
        </div>
      </div>

      {/* Tabs */}
      <Card>
        <div className="border-b border-gray-200">
          <nav className="-mb-px flex space-x-8">
            {(['overview', 'deployments'] as TabType[]).map((tab) => (
              <button
                key={tab}
                onClick={() => setActiveTab(tab)}
                className={`
                  py-4 px-1 border-b-2 font-medium text-sm
                  ${
                    activeTab === tab
                      ? 'border-blue-500 text-blue-600'
                      : 'border-transparent text-gray-500 hover:text-gray-700 hover:border-gray-300'
                  }
                `}
              >
                {tab.charAt(0).toUpperCase() + tab.slice(1)}
              </button>
            ))}
          </nav>
        </div>

        <div className="mt-6">
          {/* Overview Tab */}
          {activeTab === 'overview' && (
            <div className="space-y-6">
              {/* Tenant Information */}
              <div>
                <h3 className="text-lg font-medium text-gray-900 mb-4">Tenant Information</h3>
                <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
                  <div>
                    <label className="text-sm font-medium text-gray-500">Name</label>
                    <p className="mt-1 text-sm text-gray-900">{tenant.name}</p>
                  </div>
                  <div>
                    <label className="text-sm font-medium text-gray-500">Tenant ID</label>
                    <p className="mt-1 text-sm text-gray-900">{tenant.tenant_id}</p>
                  </div>
                  <div>
                    <label className="text-sm font-medium text-gray-500">Status</label>
                    <div className="mt-1">{getStatusBadge(tenant.status)}</div>
                  </div>
                  <div>
                    <label className="text-sm font-medium text-gray-500">Created</label>
                    <p className="mt-1 text-sm text-gray-900">{formatDate(tenant.created_at)}</p>
                  </div>
                  {tenant.description && (
                    <div className="md:col-span-2">
                      <label className="text-sm font-medium text-gray-500">Description</label>
                      <p className="mt-1 text-sm text-gray-900">{tenant.description}</p>
                    </div>
                  )}
                </div>
              </div>

              {/* Statistics */}
              {statistics && (
                <div>
                  <h3 className="text-lg font-medium text-gray-900 mb-4">Statistics</h3>
                  {loadingStats ? (
                    <Spinner size="md" />
                  ) : (
                    <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                      <Card className="p-4">
                        <div className="text-sm font-medium text-gray-500">Total Deployments</div>
                        <div className="text-2xl font-bold text-gray-900 mt-2">
                          {statistics.total_deployments}
                        </div>
                      </Card>
                      <Card className="p-4">
                        <div className="text-sm font-medium text-gray-500">Total Users</div>
                        <div className="text-2xl font-bold text-gray-900 mt-2">
                          {statistics.total_users}
                        </div>
                      </Card>
                    </div>
                  )}
                </div>
              )}

              {/* Pending Updates Summary */}
              {pendingUpdatesSummary && (
                <div>
                  <h3 className="text-lg font-medium text-gray-900 mb-4">Pending Updates</h3>
                  {loadingPendingUpdates ? (
                    <Spinner size="md" />
                  ) : (
                    <>
                      <div className="grid grid-cols-1 md:grid-cols-3 gap-4 mb-6">
                        <Card className="p-4">
                          <div className="text-sm font-medium text-gray-500">Deployments with Updates</div>
                          <div className="text-2xl font-bold text-gray-900 mt-2">
                            {pendingUpdatesSummary.deployments_with_updates}
                          </div>
                          <div className="text-xs text-gray-500 mt-1">
                            of {pendingUpdatesSummary.total_deployments} total
                          </div>
                        </Card>
                        <Card className="p-4">
                          <div className="text-sm font-medium text-gray-500">Total Pending Updates</div>
                          <div className="text-2xl font-bold text-gray-900 mt-2">
                            {pendingUpdatesSummary.total_pending_update_count}
                          </div>
                        </Card>
                        <Card className="p-4">
                          <div className="text-sm font-medium text-gray-500">By Priority</div>
                          <div className="mt-2 space-y-1">
                            <div className="flex justify-between text-sm">
                              <span>Critical:</span>
                              <span className="font-medium text-red-600">
                                {pendingUpdatesSummary.by_priority?.critical || 0}
                              </span>
                            </div>
                            <div className="flex justify-between text-sm">
                              <span>High:</span>
                              <span className="font-medium text-orange-600">
                                {pendingUpdatesSummary.by_priority?.high || 0}
                              </span>
                            </div>
                            <div className="flex justify-between text-sm">
                              <span>Normal:</span>
                              <span className="font-medium text-blue-600">
                                {pendingUpdatesSummary.by_priority?.normal || 0}
                              </span>
                            </div>
                          </div>
                        </Card>
                      </div>

                      {/* Deployments with Pending Updates List */}
                      {pendingUpdatesSummary.deployments && pendingUpdatesSummary.deployments.length > 0 && (
                        <div>
                          <h4 className="text-md font-medium text-gray-900 mb-3">Deployments Requiring Updates</h4>
                          <Card>
                            <div className="overflow-x-auto">
                              <table className="min-w-full divide-y divide-gray-200">
                                <thead className="bg-gray-50">
                                  <tr>
                                    <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                                      Deployment
                                    </th>
                                    <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                                      Product
                                    </th>
                                    <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                                      Current
                                    </th>
                                    <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                                      Latest
                                    </th>
                                    <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                                      Updates
                                    </th>
                                    <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                                      Priority
                                    </th>
                                    <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                                      Actions
                                    </th>
                                  </tr>
                                </thead>
                                <tbody className="bg-white divide-y divide-gray-200">
                                  {pendingUpdatesSummary.deployments.map((deployment) => (
                                    <tr key={deployment.deployment_id} className="hover:bg-gray-50">
                                      <td className="px-6 py-4 whitespace-nowrap text-sm font-medium text-gray-900">
                                        {deployment.deployment_id}
                                      </td>
                                      <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-500">
                                        {deployment.product_id}
                                      </td>
                                      <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-500">
                                        {deployment.current_version}
                                      </td>
                                      <td className="px-6 py-4 whitespace-nowrap text-sm font-medium text-gray-900">
                                        {deployment.latest_version}
                                      </td>
                                      <td className="px-6 py-4 whitespace-nowrap">
                                        <UpdateBadge count={deployment.update_count} priority={deployment.priority} />
                                      </td>
                                      <td className="px-6 py-4 whitespace-nowrap">
                                        <PriorityBadge priority={deployment.priority} />
                                      </td>
                                      <td className="px-6 py-4 whitespace-nowrap text-sm font-medium">
                                        {customerId && tenantId && deployment.deployment_id ? (
                                          <Button
                                            variant="ghost"
                                            size="sm"
                                            onClick={() =>
                                              navigate(
                                                `/customers/${customerId}/tenants/${tenantId}/deployments/${deployment.deployment_id}`
                                              )
                                            }
                                          >
                                            View
                                          </Button>
                                        ) : (
                                          <span className="text-gray-400 text-xs">N/A</span>
                                        )}
                                      </td>
                                    </tr>
                                  ))}
                                </tbody>
                              </table>
                            </div>
                          </Card>
                        </div>
                      )}
                    </>
                  )}
                </div>
              )}
            </div>
          )}

          {/* Deployments Tab */}
          {activeTab === 'deployments' && (
            <div>
              <div className="flex items-center justify-between mb-4">
                <h3 className="text-lg font-medium text-gray-900">Deployments</h3>
                <Button
                  variant="primary"
                  onClick={() =>
                    navigate(`/customers/${customerId}/tenants/${tenant.id}/deployments/new`)
                  }
                >
                  Add Deployment
                </Button>
              </div>
              <DeploymentsList customerId={customerId} tenantId={tenant.id} />
            </div>
          )}
        </div>
      </Card>
    </div>
  );
};

