import React, { useState, useEffect } from 'react';
import { useParams, useNavigate } from 'react-router-dom';
import { licensesApi } from '@/services/api/licenses';
import { licenseAllocationsApi } from '@/services/api/license-allocations';
import { License, LicenseStatistics, LicenseUtilization, LicenseAllocation } from '@/types';
import { Button, Card, Badge, Spinner } from '@/components/ui';
import { LicenseTypeBadge } from './LicenseTypeBadge';
import { LicenseStatusBadge } from './LicenseStatusBadge';

type TabType = 'overview' | 'allocations' | 'statistics';

export const LicenseDetails: React.FC = () => {
  const { customerId, subscriptionId, licenseId } = useParams<{
    customerId: string;
    subscriptionId: string;
    licenseId: string;
  }>();
  const navigate = useNavigate();
  const [license, setLicense] = useState<License | null>(null);
  const [statistics, setStatistics] = useState<LicenseStatistics | null>(null);
  const [utilization, setUtilization] = useState<LicenseUtilization | null>(null);
  const [allocations, setAllocations] = useState<LicenseAllocation[]>([]);
  const [loading, setLoading] = useState(true);
  const [loadingStats, setLoadingStats] = useState(false);
  const [loadingAllocations, setLoadingAllocations] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [activeTab, setActiveTab] = useState<TabType>('overview');

  useEffect(() => {
    if (customerId && subscriptionId && licenseId) {
      loadLicense();
      loadStatistics();
      loadUtilization();
      loadAllocations();
    }
  }, [customerId, subscriptionId, licenseId]);

  const loadLicense = async () => {
    if (!customerId || !subscriptionId || !licenseId) return;
    try {
      setLoading(true);
      setError(null);
      const data = await licensesApi.getById(customerId, subscriptionId, licenseId);
      setLicense(data);
    } catch (err: any) {
      setError(err.response?.data?.error?.message || 'Failed to load license');
      console.error('Error loading license:', err);
    } finally {
      setLoading(false);
    }
  };

  const loadStatistics = async () => {
    if (!customerId || !subscriptionId || !licenseId) return;
    try {
      setLoadingStats(true);
      const stats = await licensesApi.getStatistics(customerId, subscriptionId, licenseId);
      setStatistics(stats);
    } catch (err: any) {
      console.error('Error loading statistics:', err);
    } finally {
      setLoadingStats(false);
    }
  };

  const loadUtilization = async () => {
    if (!customerId || !subscriptionId || !licenseId) return;
    try {
      const util = await licenseAllocationsApi.getUtilization(customerId, subscriptionId, licenseId);
      setUtilization(util);
    } catch (err: any) {
      console.error('Error loading utilization:', err);
    }
  };

  const loadAllocations = async () => {
    if (!customerId || !subscriptionId || !licenseId) return;
    try {
      setLoadingAllocations(true);
      const response = await licenseAllocationsApi.getAll(customerId, subscriptionId, licenseId);
      setAllocations(response?.data || []);
    } catch (err: any) {
      console.error('Error loading allocations:', err);
    } finally {
      setLoadingAllocations(false);
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

  if (error || !license) {
    return (
      <div className="bg-red-50 border border-red-200 text-red-700 px-4 py-3 rounded">
        {error || 'License not found'}
      </div>
    );
  }

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-3xl font-bold text-gray-900">{license.license_id}</h1>
          <p className="text-gray-500 mt-1">License Details</p>
        </div>
        <div className="flex gap-2">
          <Button
            variant="secondary"
            onClick={() =>
              navigate(
                `/customers/${customerId}/subscriptions/${subscriptionId}/licenses/${licenseId}/edit`
              )
            }
          >
            Edit
          </Button>
          <Button
            variant="primary"
            onClick={() =>
              navigate(
                `/customers/${customerId}/subscriptions/${subscriptionId}/licenses/${licenseId}/allocate`
              )
            }
          >
            Allocate License
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
            onClick={() => setActiveTab('allocations')}
            className={`py-4 px-1 border-b-2 font-medium text-sm ${
              activeTab === 'allocations'
                ? 'border-blue-500 text-blue-600'
                : 'border-transparent text-gray-500 hover:text-gray-700 hover:border-gray-300'
            }`}
          >
            Allocations ({allocations.length})
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
        <Card title="License Information">
          <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
            <div>
              <h3 className="text-sm font-medium text-gray-500 mb-1">License ID</h3>
              <p className="text-gray-900">{license.license_id}</p>
            </div>
            <div>
              <h3 className="text-sm font-medium text-gray-500 mb-1">Product</h3>
              <p className="text-gray-900">{license.product_id}</p>
            </div>
            <div>
              <h3 className="text-sm font-medium text-gray-500 mb-1">License Type</h3>
              <LicenseTypeBadge type={license.license_type} />
            </div>
            <div>
              <h3 className="text-sm font-medium text-gray-500 mb-1">Status</h3>
              <LicenseStatusBadge status={license.status} />
            </div>
            <div>
              <h3 className="text-sm font-medium text-gray-500 mb-1">Number of Seats</h3>
              <p className="text-gray-900">{license.number_of_seats}</p>
            </div>
            <div>
              <h3 className="text-sm font-medium text-gray-500 mb-1">Start Date</h3>
              <p className="text-gray-900">{formatDate(license.start_date)}</p>
            </div>
            {license.end_date && (
              <div>
                <h3 className="text-sm font-medium text-gray-500 mb-1">End Date</h3>
                <p className="text-gray-900">{formatDate(license.end_date)}</p>
              </div>
            )}
            {license.notes && (
              <div className="md:col-span-2">
                <h3 className="text-sm font-medium text-gray-500 mb-1">Notes</h3>
                <p className="text-gray-900">{license.notes}</p>
              </div>
            )}
          </div>
        </Card>
      )}

      {activeTab === 'allocations' && (
        <div className="space-y-4">
          {loadingAllocations ? (
            <div className="flex items-center justify-center h-64">
              <Spinner size="lg" />
            </div>
          ) : allocations.length === 0 ? (
            <Card>
              <div className="text-center py-8 text-gray-500">
                No allocations found for this license
              </div>
            </Card>
          ) : (
            <Card>
              <div className="overflow-x-auto">
                <table className="min-w-full divide-y divide-gray-200">
                  <thead className="bg-gray-50">
                    <tr>
                      <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">
                        Allocation ID
                      </th>
                      <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">
                        Tenant/Deployment
                      </th>
                      <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">
                        Seats Allocated
                      </th>
                      <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">
                        Status
                      </th>
                      <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">
                        Allocation Date
                      </th>
                      <th className="px-6 py-3 text-right text-xs font-medium text-gray-500 uppercase">
                        Actions
                      </th>
                    </tr>
                  </thead>
                  <tbody className="bg-white divide-y divide-gray-200">
                    {allocations.map((allocation) => (
                      <tr key={allocation.id} className="hover:bg-gray-50">
                        <td className="px-6 py-4 whitespace-nowrap text-sm font-medium text-gray-900">
                          {allocation.allocation_id}
                        </td>
                        <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-500">
                          {allocation.tenant_id ? `Tenant: ${allocation.tenant_id}` : ''}
                          {allocation.deployment_id ? `Deployment: ${allocation.deployment_id}` : ''}
                        </td>
                        <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-500">
                          {allocation.number_of_seats_allocated}
                        </td>
                        <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-500">
                          {allocation.status}
                        </td>
                        <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-500">
                          {formatDate(allocation.allocation_date)}
                        </td>
                        <td className="px-6 py-4 whitespace-nowrap text-right text-sm font-medium">
                          {allocation.status === 'active' && (
                            <Button
                              variant="secondary"
                              size="sm"
                              onClick={async () => {
                                if (
                                  customerId &&
                                  subscriptionId &&
                                  licenseId &&
                                  window.confirm('Are you sure you want to release this allocation?')
                                ) {
                                  try {
                                    await licenseAllocationsApi.release(
                                      customerId,
                                      subscriptionId,
                                      licenseId,
                                      allocation.allocation_id
                                    );
                                    loadAllocations();
                                    loadUtilization();
                                  } catch (err: any) {
                                    alert(err.response?.data?.error?.message || 'Failed to release allocation');
                                  }
                                }
                              }}
                            >
                              Release
                            </Button>
                          )}
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
        <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
          <Card title="License Statistics">
            {loadingStats ? (
              <div className="flex items-center justify-center h-32">
                <Spinner size="lg" />
              </div>
            ) : statistics ? (
              <div className="space-y-4">
                <div>
                  <h3 className="text-sm font-medium text-gray-500 mb-1">Total Seats</h3>
                  <p className="text-2xl font-bold text-gray-900">{statistics.total_seats}</p>
                </div>
                <div>
                  <h3 className="text-sm font-medium text-gray-500 mb-1">Allocated Seats</h3>
                  <p className="text-2xl font-bold text-blue-900">{statistics.allocated_seats}</p>
                </div>
                <div>
                  <h3 className="text-sm font-medium text-gray-500 mb-1">Available Seats</h3>
                  <p className="text-2xl font-bold text-green-900">{statistics.available_seats}</p>
                </div>
                <div>
                  <h3 className="text-sm font-medium text-gray-500 mb-1">Utilization</h3>
                  <p className="text-2xl font-bold text-purple-900">
                    {statistics.utilization_percent.toFixed(1)}%
                  </p>
                </div>
              </div>
            ) : (
              <div className="text-center py-8 text-gray-500">No statistics available</div>
            )}
          </Card>

          <Card title="Utilization Details">
            {utilization ? (
              <div className="space-y-4">
                <div>
                  <h3 className="text-sm font-medium text-gray-500 mb-1">Active Allocations</h3>
                  <p className="text-2xl font-bold text-gray-900">{utilization.active_allocations}</p>
                </div>
                <div>
                  <h3 className="text-sm font-medium text-gray-500 mb-1">Total Utilization</h3>
                  <div className="mt-2">
                    <div className="w-full bg-gray-200 rounded-full h-4">
                      <div
                        className="bg-blue-600 h-4 rounded-full transition-all duration-300"
                        style={{ width: `${utilization.utilization_percent}%` }}
                      />
                    </div>
                    <p className="text-sm text-gray-600 mt-1">
                      {utilization.utilization_percent.toFixed(1)}% utilized
                    </p>
                  </div>
                </div>
              </div>
            ) : (
              <div className="text-center py-8 text-gray-500">No utilization data available</div>
            )}
          </Card>
        </div>
      )}
    </div>
  );
};

