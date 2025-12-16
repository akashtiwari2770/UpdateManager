import React, { useState, useEffect } from 'react';
import { useNavigate } from 'react-router-dom';
import { subscriptionsApi } from '@/services/api/subscriptions';
import { licensesApi } from '@/services/api/licenses';
import { customersApi } from '@/services/api/customers';
import { Subscription, License, Customer, LicenseType, LicenseStatus } from '@/types';
import { Button, Card, Badge, Spinner, Select, Input } from '@/components/ui';
import { LicenseTypeBadge, LicenseStatusBadge } from '@/components/licenses';

export const Licenses: React.FC = () => {
  const navigate = useNavigate();
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [allLicenses, setAllLicenses] = useState<Array<License & { customer: Customer; subscription: Subscription }>>([]);
  const [filteredLicenses, setFilteredLicenses] = useState<typeof allLicenses>([]);
  const [customers, setCustomers] = useState<Customer[]>([]);
  const [filters, setFilters] = useState({
    licenseType: '' as LicenseType | '',
    status: '' as LicenseStatus | '',
    customerId: '',
    search: '',
  });

  useEffect(() => {
    loadData();
  }, []);

  useEffect(() => {
    applyFilters();
  }, [filters, allLicenses]);

  const loadData = async () => {
    try {
      setLoading(true);
      setError(null);

      // Load all customers
      const customersResponse = await customersApi.getAll({ limit: 1000 });
      const customersList = customersResponse?.data || [];
      setCustomers(customersList);

      // Load all subscriptions and licenses for each customer
      const licensesData: typeof allLicenses = [];
      
      for (const customer of customersList) {
        try {
          const subscriptionsResponse = await subscriptionsApi.getAll(customer.customer_id, { limit: 1000 });
          const subscriptions = subscriptionsResponse?.data || [];

          for (const subscription of subscriptions) {
            try {
              const licensesResponse = await licensesApi.getAll(
                customer.customer_id,
                subscription.subscription_id,
                { limit: 1000 }
              );
              const licenses = licensesResponse?.data || [];

              for (const license of licenses) {
                licensesData.push({
                  ...license,
                  customer,
                  subscription,
                });
              }
            } catch (err) {
              console.error(`Error loading licenses for subscription ${subscription.subscription_id}:`, err);
            }
          }
        } catch (err) {
          console.error(`Error loading subscriptions for customer ${customer.customer_id}:`, err);
        }
      }

      setAllLicenses(licensesData);
    } catch (err: any) {
      setError(err.response?.data?.error?.message || 'Failed to load licenses');
      console.error('Error loading licenses:', err);
    } finally {
      setLoading(false);
    }
  };

  const applyFilters = () => {
    let filtered = [...allLicenses];

    if (filters.licenseType) {
      filtered = filtered.filter((l) => l.license_type === filters.licenseType);
    }

    if (filters.status) {
      filtered = filtered.filter((l) => l.status === filters.status);
    }

    if (filters.customerId) {
      filtered = filtered.filter((l) => l.customer.id === filters.customerId);
    }

    if (filters.search) {
      const searchLower = filters.search.toLowerCase();
      filtered = filtered.filter(
        (l) =>
          l.license_id.toLowerCase().includes(searchLower) ||
          l.product_id.toLowerCase().includes(searchLower) ||
          l.customer.name.toLowerCase().includes(searchLower) ||
          l.subscription.subscription_id.toLowerCase().includes(searchLower)
      );
    }

    setFilteredLicenses(filtered);
  };

  const formatDate = (dateString: string) => {
    return new Date(dateString).toLocaleDateString('en-US', {
      year: 'numeric',
      month: 'short',
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

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <h1 className="text-3xl font-bold text-gray-900">License Management</h1>
        <Button variant="primary" onClick={() => navigate('/customers')}>
          Manage via Customers
        </Button>
      </div>

      {error && (
        <div className="bg-red-50 border border-red-200 text-red-700 px-4 py-3 rounded">
          {error}
        </div>
      )}

      <Card>
        <div className="space-y-4">
          {/* Filters */}
          <div className="grid grid-cols-1 md:grid-cols-4 gap-4">
            <div>
              <label className="block text-sm font-medium text-gray-700 mb-1">Search</label>
              <Input
                placeholder="Search licenses..."
                value={filters.search}
                onChange={(e) => setFilters({ ...filters, search: e.target.value })}
              />
            </div>
            <div>
              <label className="block text-sm font-medium text-gray-700 mb-1">Customer</label>
              <Select
                value={filters.customerId}
                onChange={(e) => setFilters({ ...filters, customerId: e.target.value })}
              >
                <option value="">All Customers</option>
                {customers.map((customer) => (
                  <option key={customer.id} value={customer.id}>
                    {customer.name}
                  </option>
                ))}
              </Select>
            </div>
            <div>
              <label className="block text-sm font-medium text-gray-700 mb-1">License Type</label>
              <Select
                value={filters.licenseType}
                onChange={(e) =>
                  setFilters({ ...filters, licenseType: e.target.value as LicenseType | '' })
                }
              >
                <option value="">All Types</option>
                <option value={LicenseType.PERPETUAL}>Perpetual</option>
                <option value={LicenseType.TIME_BASED}>Time-based</option>
              </Select>
            </div>
            <div>
              <label className="block text-sm font-medium text-gray-700 mb-1">Status</label>
              <Select
                value={filters.status}
                onChange={(e) =>
                  setFilters({ ...filters, status: e.target.value as LicenseStatus | '' })
                }
              >
                <option value="">All Statuses</option>
                <option value={LicenseStatus.ACTIVE}>Active</option>
                <option value={LicenseStatus.INACTIVE}>Inactive</option>
                <option value={LicenseStatus.EXPIRED}>Expired</option>
                <option value={LicenseStatus.REVOKED}>Revoked</option>
              </Select>
            </div>
          </div>

          {/* Statistics */}
          <div className="grid grid-cols-1 md:grid-cols-4 gap-4 pt-4 border-t">
            <div className="bg-blue-50 p-4 rounded-lg">
              <h3 className="text-sm font-medium text-blue-600 mb-1">Total Licenses</h3>
              <p className="text-2xl font-bold text-blue-900">{allLicenses.length}</p>
            </div>
            <div className="bg-green-50 p-4 rounded-lg">
              <h3 className="text-sm font-medium text-green-600 mb-1">Active</h3>
              <p className="text-2xl font-bold text-green-900">
                {allLicenses.filter((l) => l.status === LicenseStatus.ACTIVE).length}
              </p>
            </div>
            <div className="bg-purple-50 p-4 rounded-lg">
              <h3 className="text-sm font-medium text-purple-600 mb-1">Perpetual</h3>
              <p className="text-2xl font-bold text-purple-900">
                {allLicenses.filter((l) => l.license_type === LicenseType.PERPETUAL).length}
              </p>
            </div>
            <div className="bg-orange-50 p-4 rounded-lg">
              <h3 className="text-sm font-medium text-orange-600 mb-1">Time-based</h3>
              <p className="text-2xl font-bold text-orange-900">
                {allLicenses.filter((l) => l.license_type === LicenseType.TIME_BASED).length}
              </p>
            </div>
          </div>

          {/* Licenses Table */}
          {filteredLicenses.length === 0 ? (
            <div className="text-center py-8 text-gray-500">
              {allLicenses.length === 0
                ? 'No licenses found'
                : 'No licenses match the current filters'}
            </div>
          ) : (
            <div className="overflow-x-auto">
              <table className="min-w-full divide-y divide-gray-200">
                <thead className="bg-gray-50">
                  <tr>
                    <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">
                      License ID
                    </th>
                    <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">
                      Customer
                    </th>
                    <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">
                      Subscription
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
                    <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">
                      End Date
                    </th>
                    <th className="px-6 py-3 text-right text-xs font-medium text-gray-500 uppercase">
                      Actions
                    </th>
                  </tr>
                </thead>
                <tbody className="bg-white divide-y divide-gray-200">
                  {filteredLicenses.map((license) => (
                    <tr key={license.id} className="hover:bg-gray-50">
                      <td className="px-6 py-4 whitespace-nowrap text-sm font-medium text-gray-900">
                        {license.license_id}
                      </td>
                      <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-900">
                        <button
                          onClick={() => navigate(`/customers/${license.customer.id}`)}
                          className="text-blue-600 hover:text-blue-800"
                        >
                          {license.customer.name}
                        </button>
                      </td>
                      <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-500">
                        {license.subscription.subscription_id}
                      </td>
                      <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-500">
                        {license.product_id}
                      </td>
                      <td className="px-6 py-4 whitespace-nowrap">
                        <LicenseTypeBadge type={license.license_type} />
                      </td>
                      <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-500">
                        {license.number_of_seats}
                      </td>
                      <td className="px-6 py-4 whitespace-nowrap">
                        <LicenseStatusBadge status={license.status} />
                      </td>
                      <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-500">
                        {license.end_date ? formatDate(license.end_date) : '-'}
                      </td>
                      <td className="px-6 py-4 whitespace-nowrap text-right text-sm font-medium">
                        <Button
                          variant="secondary"
                          size="sm"
                          onClick={() =>
                            navigate(
                              `/customers/${license.customer.customer_id}/subscriptions/${license.subscription.subscription_id}/licenses/${license.license_id}`
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
          )}
        </div>
      </Card>
    </div>
  );
};

