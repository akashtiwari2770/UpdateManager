import React, { useState, useEffect } from 'react';
import { useParams, useNavigate } from 'react-router-dom';
import { customersApi } from '@/services/api/customers';
import {
  Customer,
  CustomerStatus,
  CustomerStatistics,
} from '@/types';
import { Button, Card, Badge, Spinner } from '@/components/ui';
import { TenantsList } from '@/components/tenants';
import { SubscriptionsList } from '@/components/subscriptions';
import { pendingUpdatesApi } from '@/services/api/pending-updates';
import { CustomerPendingUpdatesSummary } from '@/types';

type TabType = 'overview' | 'tenants' | 'subscriptions' | 'activity';

export const CustomerDetails: React.FC = () => {
  const { id } = useParams<{ id: string }>();
  const navigate = useNavigate();
  const [customer, setCustomer] = useState<Customer | null>(null);
  const [statistics, setStatistics] = useState<CustomerStatistics | null>(null);
  const [pendingUpdatesSummary, setPendingUpdatesSummary] = useState<CustomerPendingUpdatesSummary | null>(null);
  const [loading, setLoading] = useState(true);
  const [loadingStats, setLoadingStats] = useState(false);
  const [loadingPendingUpdates, setLoadingPendingUpdates] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [activeTab, setActiveTab] = useState<TabType>('overview');

  useEffect(() => {
    if (id) {
      loadCustomer();
      loadStatistics();
      loadPendingUpdates();
    }
  }, [id]);

  const loadCustomer = async () => {
    try {
      setLoading(true);
      setError(null);
      const data = await customersApi.getById(id!);
      setCustomer(data);
    } catch (err: any) {
      setError(err.response?.data?.error?.message || 'Failed to load customer');
      console.error('Error loading customer:', err);
    } finally {
      setLoading(false);
    }
  };

  const loadStatistics = async () => {
    if (!id) return;
    try {
      setLoadingStats(true);
      const stats = await customersApi.getStatistics(id);
      setStatistics(stats);
    } catch (err: any) {
      console.error('Error loading statistics:', err);
    } finally {
      setLoadingStats(false);
    }
  };

  const loadPendingUpdates = async () => {
    if (!id) return;
    try {
      setLoadingPendingUpdates(true);
      const summary = await pendingUpdatesApi.getCustomerPendingUpdates(id);
      setPendingUpdatesSummary(summary);
    } catch (err: any) {
      console.error('Error loading pending updates:', err);
    } finally {
      setLoadingPendingUpdates(false);
    }
  };

  const getStatusBadge = (status: CustomerStatus) => {
    const statusConfig = {
      [CustomerStatus.ACTIVE]: { label: 'Active', className: 'bg-green-100 text-green-800' },
      [CustomerStatus.INACTIVE]: { label: 'Inactive', className: 'bg-gray-100 text-gray-800' },
      [CustomerStatus.SUSPENDED]: { label: 'Suspended', className: 'bg-red-100 text-red-800' },
    };
    const config = statusConfig[status] || statusConfig[CustomerStatus.INACTIVE];
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

  if (error || !customer) {
    return (
      <div className="space-y-6">
        <div className="flex items-center justify-between">
          <h1 className="text-3xl font-bold text-gray-900">Customer Details</h1>
          <Button variant="ghost" onClick={() => navigate('/customers')}>
            Back to Customers
          </Button>
        </div>
        <Card>
          <div className="text-center py-12">
            <p className="text-red-600 mb-4">{error || 'Customer not found'}</p>
            <Button variant="primary" onClick={() => navigate('/customers')}>
              Back to Customers
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
          <h1 className="text-3xl font-bold text-gray-900">{customer.name}</h1>
          <p className="text-gray-500 mt-1">Customer ID: {customer.customer_id}</p>
        </div>
        <div className="flex space-x-2">
          <Button variant="ghost" onClick={() => navigate('/customers')}>
            Back
          </Button>
          <Button variant="primary" onClick={() => navigate(`/customers/${customer.id}/edit`)}>
            Edit
          </Button>
        </div>
      </div>

      {/* Tabs */}
      <Card>
        <div className="border-b border-gray-200">
          <nav className="-mb-px flex space-x-8">
            {(['overview', 'tenants', 'subscriptions', 'activity'] as TabType[]).map((tab) => (
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
              {/* Customer Information */}
              <div>
                <h3 className="text-lg font-medium text-gray-900 mb-4">Customer Information</h3>
                <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
                  <div>
                    <label className="text-sm font-medium text-gray-500">Name</label>
                    <p className="mt-1 text-sm text-gray-900">{customer.name}</p>
                  </div>
                  <div>
                    <label className="text-sm font-medium text-gray-500">Customer ID</label>
                    <p className="mt-1 text-sm text-gray-900">{customer.customer_id}</p>
                  </div>
                  <div>
                    <label className="text-sm font-medium text-gray-500">Organization</label>
                    <p className="mt-1 text-sm text-gray-900">
                      {customer.organization_name || '-'}
                    </p>
                  </div>
                  <div>
                    <label className="text-sm font-medium text-gray-500">Email</label>
                    <p className="mt-1 text-sm text-gray-900">{customer.email}</p>
                  </div>
                  <div>
                    <label className="text-sm font-medium text-gray-500">Phone</label>
                    <p className="mt-1 text-sm text-gray-900">{customer.phone || '-'}</p>
                  </div>
                  <div>
                    <label className="text-sm font-medium text-gray-500">Status</label>
                    <div className="mt-1">{getStatusBadge(customer.account_status)}</div>
                  </div>
                  {customer.address && (
                    <div className="md:col-span-2">
                      <label className="text-sm font-medium text-gray-500">Address</label>
                      <p className="mt-1 text-sm text-gray-900">{customer.address}</p>
                    </div>
                  )}
                  <div>
                    <label className="text-sm font-medium text-gray-500">Created</label>
                    <p className="mt-1 text-sm text-gray-900">{formatDate(customer.created_at)}</p>
                  </div>
                  <div>
                    <label className="text-sm font-medium text-gray-500">Last Updated</label>
                    <p className="mt-1 text-sm text-gray-900">{formatDate(customer.updated_at)}</p>
                  </div>
                </div>
              </div>

              {/* Notification Preferences */}
              <div>
                <h3 className="text-lg font-medium text-gray-900 mb-4">Notification Preferences</h3>
                <div className="grid grid-cols-2 gap-4">
                  <div className="flex items-center">
                    <span className="text-sm text-gray-700">
                      Email Notifications:{' '}
                      {customer.notification_preferences.email_enabled ? 'Enabled' : 'Disabled'}
                    </span>
                  </div>
                  <div className="flex items-center">
                    <span className="text-sm text-gray-700">
                      In-App Notifications:{' '}
                      {customer.notification_preferences.in_app_enabled ? 'Enabled' : 'Disabled'}
                    </span>
                  </div>
                  <div className="flex items-center">
                    <span className="text-sm text-gray-700">
                      UAT Notifications:{' '}
                      {customer.notification_preferences.uat_notifications ? 'Enabled' : 'Disabled'}
                    </span>
                  </div>
                  <div className="flex items-center">
                    <span className="text-sm text-gray-700">
                      Production Notifications:{' '}
                      {customer.notification_preferences.production_notifications
                        ? 'Enabled'
                        : 'Disabled'}
                    </span>
                  </div>
                </div>
              </div>

              {/* Statistics */}
              {statistics && (
                <div>
                  <h3 className="text-lg font-medium text-gray-900 mb-4">Statistics</h3>
                  {loadingStats ? (
                    <Spinner size="md" />
                  ) : (
                    <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
                      <Card className="p-4">
                        <div className="text-sm font-medium text-gray-500">Total Tenants</div>
                        <div className="text-2xl font-bold text-gray-900 mt-2">
                          {statistics.total_tenants}
                        </div>
                      </Card>
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
                    <div className="grid grid-cols-1 md:grid-cols-4 gap-4">
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
                        <div className="text-sm font-medium text-gray-500">Critical</div>
                        <div className="text-2xl font-bold text-red-600 mt-2">
                          {pendingUpdatesSummary.by_priority?.critical || 0}
                        </div>
                      </Card>
                      <Card className="p-4">
                        <div className="text-sm font-medium text-gray-500">High Priority</div>
                        <div className="text-2xl font-bold text-orange-600 mt-2">
                          {pendingUpdatesSummary.by_priority?.high || 0}
                        </div>
                      </Card>
                    </div>
                  )}
                </div>
              )}
            </div>
          )}

          {/* Tenants Tab */}
          {activeTab === 'tenants' && (
            <div>
              <div className="flex items-center justify-between mb-4">
                <h3 className="text-lg font-medium text-gray-900">Tenants</h3>
                <Button
                  variant="primary"
                  onClick={() => navigate(`/customers/${customer.id}/tenants/new`)}
                >
                  Add Tenant
                </Button>
              </div>
              <TenantsList customerId={customer.customer_id} />
            </div>
          )}

          {/* Subscriptions Tab */}
          {activeTab === 'subscriptions' && (
            <div>
              <SubscriptionsList customerId={customer.customer_id} />
            </div>
          )}

          {/* Activity Tab */}
          {activeTab === 'activity' && (
            <div>
              <h3 className="text-lg font-medium text-gray-900 mb-4">Recent Activity</h3>
              <p className="text-gray-500">Activity timeline will be displayed here</p>
              {/* TODO: Implement activity timeline with notifications and audit logs */}
            </div>
          )}
        </div>
      </Card>
    </div>
  );
};

