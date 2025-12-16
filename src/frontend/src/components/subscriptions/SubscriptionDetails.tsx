import React, { useState, useEffect } from 'react';
import { useParams, useNavigate } from 'react-router-dom';
import { subscriptionsApi } from '@/services/api/subscriptions';
import { licensesApi } from '@/services/api/licenses';
import { Subscription, SubscriptionStatistics, License } from '@/types';
import { Button, Card, Badge, Spinner } from '@/components/ui';
import { SubscriptionStatusBadge } from './SubscriptionStatusBadge';

type TabType = 'overview' | 'licenses' | 'statistics';

export const SubscriptionDetails: React.FC = () => {
  const { customerId, subscriptionId } = useParams<{ customerId: string; subscriptionId: string }>();
  const navigate = useNavigate();
  const [subscription, setSubscription] = useState<Subscription | null>(null);
  const [statistics, setStatistics] = useState<SubscriptionStatistics | null>(null);
  const [licenses, setLicenses] = useState<License[]>([]);
  const [loading, setLoading] = useState(true);
  const [loadingStats, setLoadingStats] = useState(false);
  const [loadingLicenses, setLoadingLicenses] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [activeTab, setActiveTab] = useState<TabType>('overview');

  useEffect(() => {
    if (customerId && subscriptionId) {
      loadSubscription();
      loadStatistics();
      loadLicenses();
    }
  }, [customerId, subscriptionId]);

  const loadSubscription = async () => {
    if (!customerId || !subscriptionId) return;
    try {
      setLoading(true);
      setError(null);
      const data = await subscriptionsApi.getById(customerId, subscriptionId);
      setSubscription(data);
    } catch (err: any) {
      setError(err.response?.data?.error?.message || 'Failed to load subscription');
      console.error('Error loading subscription:', err);
    } finally {
      setLoading(false);
    }
  };

  const loadStatistics = async () => {
    if (!customerId || !subscriptionId) return;
    try {
      setLoadingStats(true);
      const stats = await subscriptionsApi.getStatistics(customerId, subscriptionId);
      setStatistics(stats);
    } catch (err: any) {
      console.error('Error loading statistics:', err);
    } finally {
      setLoadingStats(false);
    }
  };

  const loadLicenses = async () => {
    if (!customerId || !subscriptionId) return;
    try {
      setLoadingLicenses(true);
      const response = await licensesApi.getAll(customerId, subscriptionId);
      setLicenses(response?.data || []);
    } catch (err: any) {
      console.error('Error loading licenses:', err);
    } finally {
      setLoadingLicenses(false);
    }
  };

  const formatDate = (dateString: string) => {
    return new Date(dateString).toLocaleDateString('en-US', {
      year: 'numeric',
      month: 'long',
      day: 'numeric',
    });
  };

  if (loading) {
    return (
      <div className="flex items-center justify-center h-64">
        <Spinner size="lg" />
      </div>
    );
  }

  if (error || !subscription) {
    return (
      <div className="bg-red-50 border border-red-200 text-red-700 px-4 py-3 rounded">
        {error || 'Subscription not found'}
      </div>
    );
  }

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-3xl font-bold text-gray-900">
            {subscription.name || subscription.subscription_id}
          </h1>
          <p className="text-gray-500 mt-1">Subscription Details</p>
        </div>
        <div className="flex gap-2">
          <Button
            variant="secondary"
            onClick={() => navigate(`/customers/${customerId}/subscriptions/${subscriptionId}/edit`)}
          >
            Edit
          </Button>
          <Button
            variant="primary"
            onClick={() => navigate(`/customers/${customerId}/subscriptions/${subscriptionId}/licenses/new`)}
          >
            Assign License
          </Button>
        </div>
      </div>

      {/* Tabs */}
      <div className="border-b border-gray-200">
        <nav className="-mb-px flex space-x-8">
          <button
            onClick={() => setActiveTab('overview')}
            className={`py-4 px-1 border-b-2 font-medium text-sm ${
              activeTab === 'overview'
                ? 'border-blue-500 text-blue-600'
                : 'border-transparent text-gray-500 hover:text-gray-700 hover:border-gray-300'
            }`}
          >
            Overview
          </button>
          <button
            onClick={() => setActiveTab('licenses')}
            className={`py-4 px-1 border-b-2 font-medium text-sm ${
              activeTab === 'licenses'
                ? 'border-blue-500 text-blue-600'
                : 'border-transparent text-gray-500 hover:text-gray-700 hover:border-gray-300'
            }`}
          >
            Licenses ({licenses.length})
          </button>
          <button
            onClick={() => setActiveTab('statistics')}
            className={`py-4 px-1 border-b-2 font-medium text-sm ${
              activeTab === 'statistics'
                ? 'border-blue-500 text-blue-600'
                : 'border-transparent text-gray-500 hover:text-gray-700 hover:border-gray-300'
            }`}
          >
            Statistics
          </button>
        </nav>
      </div>

      {/* Tab Content */}
      {activeTab === 'overview' && (
        <Card title="Subscription Information">
          <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
            <div>
              <h3 className="text-sm font-medium text-gray-500 mb-1">Subscription ID</h3>
              <p className="text-gray-900">{subscription.subscription_id}</p>
            </div>
            <div>
              <h3 className="text-sm font-medium text-gray-500 mb-1">Status</h3>
              <SubscriptionStatusBadge status={subscription.status} />
            </div>
            {subscription.name && (
              <div>
                <h3 className="text-sm font-medium text-gray-500 mb-1">Name</h3>
                <p className="text-gray-900">{subscription.name}</p>
              </div>
            )}
            <div>
              <h3 className="text-sm font-medium text-gray-500 mb-1">Start Date</h3>
              <p className="text-gray-900">{formatDate(subscription.start_date)}</p>
            </div>
            {subscription.end_date && (
              <div>
                <h3 className="text-sm font-medium text-gray-500 mb-1">End Date</h3>
                <p className="text-gray-900">{formatDate(subscription.end_date)}</p>
              </div>
            )}
            {subscription.description && (
              <div className="md:col-span-2">
                <h3 className="text-sm font-medium text-gray-500 mb-1">Description</h3>
                <p className="text-gray-900">{subscription.description}</p>
              </div>
            )}
            {subscription.notes && (
              <div className="md:col-span-2">
                <h3 className="text-sm font-medium text-gray-500 mb-1">Notes</h3>
                <p className="text-gray-900">{subscription.notes}</p>
              </div>
            )}
          </div>
        </Card>
      )}

      {activeTab === 'licenses' && (
        <div className="space-y-4">
          {loadingLicenses ? (
            <div className="flex items-center justify-center h-64">
              <Spinner size="lg" />
            </div>
          ) : licenses.length === 0 ? (
            <Card>
              <div className="text-center py-8 text-gray-500">
                No licenses assigned to this subscription
              </div>
            </Card>
          ) : (
            <Card>
              <div className="overflow-x-auto">
                <table className="min-w-full divide-y divide-gray-200">
                  <thead className="bg-gray-50">
                    <tr>
                      <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">
                        License ID
                      </th>
                      <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">
                        Product
                      </th>
                      <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">
                        Type
                      </th>
                      <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">
                        Seats
                      </th>
                      <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">
                        Status
                      </th>
                      <th className="px-6 py-3 text-right text-xs font-medium text-gray-500 uppercase">
                        Actions
                      </th>
                    </tr>
                  </thead>
                  <tbody className="bg-white divide-y divide-gray-200">
                    {licenses.map((license) => (
                      <tr key={license.id} className="hover:bg-gray-50">
                        <td className="px-6 py-4 whitespace-nowrap text-sm font-medium text-gray-900">
                          {license.license_id}
                        </td>
                        <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-900">
                          {license.product_id}
                        </td>
                        <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-500">
                          {license.license_type}
                        </td>
                        <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-500">
                          {license.number_of_seats}
                        </td>
                        <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-500">
                          {license.status}
                        </td>
                        <td className="px-6 py-4 whitespace-nowrap text-right text-sm font-medium">
                          <Button
                            variant="secondary"
                            size="sm"
                            onClick={() =>
                              navigate(
                                `/customers/${customerId}/subscriptions/${subscriptionId}/licenses/${license.license_id}`
                              )
                            }
                          >
                            View
                          </Button>
                        </td>
                      </tr>
                    ))}
                  </tbody>
                </table>
              </div>
            </Card>
          )}
        </div>
      )}

      {activeTab === 'statistics' && (
        <Card title="Subscription Statistics">
          {loadingStats ? (
            <div className="flex items-center justify-center h-64">
              <Spinner size="lg" />
            </div>
          ) : statistics ? (
            <div className="grid grid-cols-1 md:grid-cols-3 gap-6">
              <div className="bg-blue-50 p-4 rounded-lg">
                <h3 className="text-sm font-medium text-blue-600 mb-1">Total Licenses</h3>
                <p className="text-2xl font-bold text-blue-900">{statistics.total_licenses}</p>
              </div>
              <div className="bg-green-50 p-4 rounded-lg">
                <h3 className="text-sm font-medium text-green-600 mb-1">Active Licenses</h3>
                <p className="text-2xl font-bold text-green-900">{statistics.active_licenses}</p>
              </div>
              <div className="bg-gray-50 p-4 rounded-lg">
                <h3 className="text-sm font-medium text-gray-600 mb-1">Total Seats</h3>
                <p className="text-2xl font-bold text-gray-900">{statistics.total_seats}</p>
              </div>
              <div className="bg-purple-50 p-4 rounded-lg">
                <h3 className="text-sm font-medium text-purple-600 mb-1">Perpetual Licenses</h3>
                <p className="text-2xl font-bold text-purple-900">{statistics.perpetual_licenses}</p>
              </div>
              <div className="bg-orange-50 p-4 rounded-lg">
                <h3 className="text-sm font-medium text-orange-600 mb-1">Time-based Licenses</h3>
                <p className="text-2xl font-bold text-orange-900">{statistics.time_based_licenses}</p>
              </div>
              <div className="bg-red-50 p-4 rounded-lg">
                <h3 className="text-sm font-medium text-red-600 mb-1">Expired Licenses</h3>
                <p className="text-2xl font-bold text-red-900">{statistics.expired_licenses}</p>
              </div>
            </div>
          ) : (
            <div className="text-center py-8 text-gray-500">No statistics available</div>
          )}
        </Card>
      )}
    </div>
  );
};

