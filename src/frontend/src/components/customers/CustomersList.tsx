import React, { useState, useEffect } from 'react';
import { useNavigate } from 'react-router-dom';
import { customersApi } from '@/services/api/customers';
import { Customer, CustomerStatus, ListCustomersQuery } from '@/types';
import { Button, Card, Badge, Spinner, Input, Select } from '@/components/ui';

export const CustomersList: React.FC = () => {
  const navigate = useNavigate();
  const [customers, setCustomers] = useState<Customer[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  
  const [filters, setFilters] = useState<ListCustomersQuery>({
    page: 1,
    limit: 20,
  });
  const [searchTerm, setSearchTerm] = useState('');
  const [totalPages, setTotalPages] = useState(1);
  const [total, setTotal] = useState(0);

  useEffect(() => {
    loadCustomers();
  }, [filters]);

  const loadCustomers = async () => {
    try {
      setLoading(true);
      setError(null);
      
      const query: ListCustomersQuery = {
        ...filters,
        search: searchTerm || undefined,
        page: filters.page || 1,
        limit: filters.limit || 20,
      };

      const response = await customersApi.getAll(query);
      setCustomers(response?.data || []);
      setTotalPages(response?.pagination?.total_pages || 1);
      setTotal(response?.pagination?.total || 0);
    } catch (err: any) {
      setError(err.response?.data?.error?.message || 'Failed to load customers');
      console.error('Error loading customers:', err);
    } finally {
      setLoading(false);
    }
  };

  const handleSearch = (value: string) => {
    setSearchTerm(value);
    setFilters((prev) => ({ ...prev, search: value || undefined, page: 1 }));
  };

  const handleFilterChange = (key: keyof ListCustomersQuery, value: any) => {
    setFilters((prev) => ({
      ...prev,
      [key]: value,
      page: 1,
    }));
  };

  const handlePageChange = (newPage: number) => {
    setFilters((prev) => ({ ...prev, page: newPage }));
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

  if (loading && customers.length === 0) {
    return (
      <div className="flex items-center justify-center h-64">
        <Spinner size="lg" />
      </div>
    );
  }

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <h1 className="text-3xl font-bold text-gray-900">Customers</h1>
        <Button onClick={() => navigate('/customers/new')} variant="primary">
          Create Customer
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
          <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
            <div className="md:col-span-2">
              <Input
                label="Search"
                placeholder="Search by name, email, or organization..."
                value={searchTerm}
                onChange={(e) => handleSearch(e.target.value)}
              />
            </div>
            <Select
              label="Status"
              value={filters.status || ''}
              onChange={(e) =>
                handleFilterChange('status', e.target.value || undefined)
              }
              options={[
                { value: '', label: 'All Statuses' },
                { value: CustomerStatus.ACTIVE, label: 'Active' },
                { value: CustomerStatus.INACTIVE, label: 'Inactive' },
                { value: CustomerStatus.SUSPENDED, label: 'Suspended' },
              ]}
            />
          </div>

          {/* Table */}
          {customers.length === 0 ? (
            <div className="text-center py-12">
              <p className="text-gray-500">No customers found</p>
              <Button
                onClick={() => navigate('/customers/new')}
                variant="primary"
                className="mt-4"
              >
                Create First Customer
              </Button>
            </div>
          ) : (
            <>
              <div className="overflow-x-auto">
                <table className="min-w-full divide-y divide-gray-200">
                  <thead className="bg-gray-50">
                    <tr>
                      <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                        Customer ID
                      </th>
                      <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                        Name
                      </th>
                      <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                        Organization
                      </th>
                      <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                        Email
                      </th>
                      <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                        Status
                      </th>
                      <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                        Actions
                      </th>
                    </tr>
                  </thead>
                  <tbody className="bg-white divide-y divide-gray-200">
                    {customers.map((customer) => (
                      <tr key={customer.id} className="hover:bg-gray-50">
                        <td className="px-6 py-4 whitespace-nowrap text-sm font-medium text-gray-900">
                          {customer.customer_id}
                        </td>
                        <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-900">
                          {customer.name}
                        </td>
                        <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-500">
                          {customer.organization_name || '-'}
                        </td>
                        <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-500">
                          {customer.email}
                        </td>
                        <td className="px-6 py-4 whitespace-nowrap">
                          {getStatusBadge(customer.account_status)}
                        </td>
                        <td className="px-6 py-4 whitespace-nowrap text-sm font-medium">
                          <div className="flex space-x-2">
                            <Button
                              variant="ghost"
                              size="sm"
                              onClick={() => navigate(`/customers/${customer.id}`)}
                            >
                              View
                            </Button>
                            <Button
                              variant="ghost"
                              size="sm"
                              onClick={() => navigate(`/customers/${customer.id}/edit`)}
                            >
                              Edit
                            </Button>
                          </div>
                        </td>
                      </tr>
                    ))}
                  </tbody>
                </table>
              </div>

              {/* Pagination */}
              {totalPages > 1 && (
                <div className="flex items-center justify-between border-t border-gray-200 bg-white px-4 py-3 sm:px-6">
                  <div className="flex flex-1 justify-between sm:hidden">
                    <Button
                      variant="ghost"
                      size="sm"
                      onClick={() => handlePageChange(Math.max(1, (filters.page || 1) - 1))}
                      disabled={filters.page === 1}
                    >
                      Previous
                    </Button>
                    <Button
                      variant="ghost"
                      size="sm"
                      onClick={() => handlePageChange(Math.min(totalPages, (filters.page || 1) + 1))}
                      disabled={filters.page === totalPages}
                    >
                      Next
                    </Button>
                  </div>
                  <div className="hidden sm:flex sm:flex-1 sm:items-center sm:justify-between">
                    <div>
                      <p className="text-sm text-gray-700">
                        Showing{' '}
                        <span className="font-medium">
                          {((filters.page || 1) - 1) * (filters.limit || 20) + 1}
                        </span>{' '}
                        to{' '}
                        <span className="font-medium">
                          {Math.min((filters.page || 1) * (filters.limit || 20), total)}
                        </span>{' '}
                        of <span className="font-medium">{total}</span> results
                      </p>
                    </div>
                    <div>
                      <nav className="isolate inline-flex -space-x-px rounded-md shadow-sm">
                        <Button
                          variant="ghost"
                          size="sm"
                          onClick={() => handlePageChange(Math.max(1, (filters.page || 1) - 1))}
                          disabled={filters.page === 1}
                        >
                          Previous
                        </Button>
                        {Array.from({ length: totalPages }, (_, i) => i + 1)
                          .filter(
                            (page) =>
                              page === 1 ||
                              page === totalPages ||
                              (page >= (filters.page || 1) - 2 && page <= (filters.page || 1) + 2)
                          )
                          .map((page, idx, arr) => (
                            <React.Fragment key={page}>
                              {idx > 0 && arr[idx - 1] !== page - 1 && (
                                <span className="px-4 py-2 text-sm text-gray-700">...</span>
                              )}
                              <Button
                                variant={filters.page === page ? 'primary' : 'ghost'}
                                size="sm"
                                onClick={() => handlePageChange(page)}
                              >
                                {page}
                              </Button>
                            </React.Fragment>
                          ))}
                        <Button
                          variant="ghost"
                          size="sm"
                          onClick={() => handlePageChange(Math.min(totalPages, (filters.page || 1) + 1))}
                          disabled={filters.page === totalPages}
                        >
                          Next
                        </Button>
                      </nav>
                    </div>
                  </div>
                </div>
              )}
            </>
          )}
        </div>
      </Card>
    </div>
  );
};

