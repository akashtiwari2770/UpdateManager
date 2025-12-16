import React, { useState, useEffect } from 'react';
import { useNavigate } from 'react-router-dom';
import { versionsApi } from '@/services/api/versions';
import { productsApi } from '@/services/api/products';
import { Version, VersionState, ReleaseType, Product, ListVersionsQuery } from '@/types';
import { Button, Card, Badge, Spinner, Input, Select, Alert } from '@/components/ui';

export const VersionsList: React.FC = () => {
  const navigate = useNavigate();
  const [versions, setVersions] = useState<Version[]>([]);
  const [products, setProducts] = useState<Product[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  
  // Filters and pagination
  const [filters, setFilters] = useState<ListVersionsQuery>({
    page: 1,
    limit: 25,
  });
  const [searchTerm, setSearchTerm] = useState('');
  const [totalPages, setTotalPages] = useState(1);
  const [total, setTotal] = useState(0);

  // Sorting
  const [sortField, setSortField] = useState<string>('release_date');
  const [sortOrder, setSortOrder] = useState<'asc' | 'desc'>('desc');

  useEffect(() => {
    loadProducts();
    loadVersions();
  }, [filters, sortField, sortOrder]);

  const loadProducts = async () => {
    try {
      const response = await productsApi.getAll({ page: 1, limit: 1000, is_active: true });
      setProducts(response.data || []);
    } catch (err) {
      console.error('Error loading products:', err);
    }
  };

  const loadVersions = async () => {
    try {
      setLoading(true);
      setError(null);
      
      const query: ListVersionsQuery = {
        ...filters,
        page: filters.page || 1,
        limit: filters.limit || 25,
      };

      const response = await versionsApi.getAll(query);
      setVersions(response?.data || []);
      setTotalPages(response?.pagination?.total_pages || 1);
      setTotal(response?.pagination?.total || 0);
    } catch (err: any) {
      setError(err.response?.data?.error?.message || 'Failed to load versions');
      console.error('Error loading versions:', err);
    } finally {
      setLoading(false);
    }
  };

  const handleSearch = (value: string) => {
    setSearchTerm(value);
    // Client-side filtering for now
    if (value.trim()) {
      const filtered = versions.filter(
        (v) =>
          v.version_number.toLowerCase().includes(value.toLowerCase()) ||
          v.product_id.toLowerCase().includes(value.toLowerCase())
      );
      setVersions(filtered);
    } else {
      loadVersions();
    }
  };

  const handleFilterChange = (key: keyof ListVersionsQuery, value: any) => {
    setFilters((prev) => ({
      ...prev,
      [key]: value,
      page: 1,
    }));
  };

  const handleSort = (field: string) => {
    if (sortField === field) {
      setSortOrder(sortOrder === 'asc' ? 'desc' : 'asc');
    } else {
      setSortField(field);
      setSortOrder('asc');
    }
  };

  const handlePageChange = (newPage: number) => {
    setFilters((prev) => ({ ...prev, page: newPage }));
  };

  const handlePageSizeChange = (newSize: number) => {
    setFilters((prev) => ({ ...prev, limit: newSize, page: 1 }));
  };

  const getSortIcon = (field: string) => {
    if (sortField !== field) return '↕️';
    return sortOrder === 'asc' ? '↑' : '↓';
  };

  const formatDate = (dateString: string) => {
    if (!dateString) return 'N/A';
    return new Date(dateString).toLocaleDateString();
  };

  const getStateBadgeColor = (state: VersionState) => {
    switch (state) {
      case VersionState.DRAFT:
        return 'gray';
      case VersionState.PENDING_REVIEW:
        return 'yellow';
      case VersionState.APPROVED:
        return 'blue';
      case VersionState.RELEASED:
        return 'green';
      case VersionState.DEPRECATED:
        return 'orange';
      case VersionState.EOL:
        return 'red';
      default:
        return 'gray';
    }
  };

  const getReleaseTypeLabel = (type: ReleaseType) => {
    switch (type) {
      case ReleaseType.SECURITY:
        return 'Security';
      case ReleaseType.FEATURE:
        return 'Feature';
      case ReleaseType.MAINTENANCE:
        return 'Maintenance';
      case ReleaseType.MAJOR:
        return 'Major';
      default:
        return type;
    }
  };

  const getProductName = (productId: string) => {
    const product = products.find((p) => p.product_id === productId);
    return product?.name || productId;
  };

  if (loading && versions.length === 0) {
    return (
      <div className="flex justify-center items-center h-64">
        <Spinner />
      </div>
    );
  }

  if (error) {
    return (
      <div className="space-y-6">
        <div className="flex items-center justify-between">
          <h1 className="text-3xl font-bold text-gray-900">Versions</h1>
          <Button onClick={() => navigate('/versions/new')}>
            Create Version
          </Button>
        </div>
        <Alert variant="error" title="Error" onClose={() => setError(null)}>
          {error}
        </Alert>
      </div>
    );
  }

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <h1 className="text-3xl font-bold text-gray-900">Versions</h1>
        <Button onClick={() => navigate('/versions/new')}>
          Create Version
        </Button>
      </div>

      {error && (
        <Alert variant="error" title="Error" onClose={() => setError(null)}>
          {error}
        </Alert>
      )}

      <Card>
        {/* Filters */}
        <div className="grid grid-cols-1 md:grid-cols-4 gap-4 mb-6">
          <Input
            placeholder="Search by version or product..."
            value={searchTerm}
            onChange={(e) => handleSearch(e.target.value)}
          />
          
          <Select
            value={filters.product_id || ''}
            onChange={(e) => handleFilterChange('product_id', e.target.value || undefined)}
          >
            <option value="">All Products</option>
            {products.map((product) => (
              <option key={product.id} value={product.product_id}>
                {product.name}
              </option>
            ))}
          </Select>

          <Select
            value={filters.state || ''}
            onChange={(e) => handleFilterChange('state', e.target.value || undefined)}
          >
            <option value="">All States</option>
            {Object.values(VersionState).map((state) => (
              <option key={state} value={state}>
                {state.replace('_', ' ').replace(/\b\w/g, (l) => l.toUpperCase())}
              </option>
            ))}
          </Select>

          <Select
            value={filters.release_type || ''}
            onChange={(e) => handleFilterChange('release_type', e.target.value || undefined)}
          >
            <option value="">All Release Types</option>
            {Object.values(ReleaseType).map((type) => (
              <option key={type} value={type}>
                {getReleaseTypeLabel(type)}
              </option>
            ))}
          </Select>
        </div>

        {/* Table */}
        {versions.length === 0 ? (
          <div className="text-center py-12">
            <p className="text-gray-500">No versions found</p>
          </div>
        ) : (
          <>
            <div className="overflow-x-auto">
              <table className="min-w-full divide-y divide-gray-200">
                <thead className="bg-gray-50">
                  <tr>
                    <th
                      className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider cursor-pointer hover:bg-gray-100"
                      onClick={() => handleSort('product_id')}
                    >
                      Product {getSortIcon('product_id')}
                    </th>
                    <th
                      className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider cursor-pointer hover:bg-gray-100"
                      onClick={() => handleSort('version_number')}
                    >
                      Version {getSortIcon('version_number')}
                    </th>
                    <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                      Release Type
                    </th>
                    <th
                      className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider cursor-pointer hover:bg-gray-100"
                      onClick={() => handleSort('state')}
                    >
                      State {getSortIcon('state')}
                    </th>
                    <th
                      className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider cursor-pointer hover:bg-gray-100"
                      onClick={() => handleSort('release_date')}
                    >
                      Release Date {getSortIcon('release_date')}
                    </th>
                    <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                      Approved By
                    </th>
                    <th className="px-6 py-3 text-right text-xs font-medium text-gray-500 uppercase tracking-wider">
                      Actions
                    </th>
                  </tr>
                </thead>
                <tbody className="bg-white divide-y divide-gray-200">
                  {versions.map((version) => (
                    <tr
                      key={version.id}
                      className="hover:bg-gray-50 cursor-pointer"
                      onClick={() => {
                        if (version.id) {
                          navigate(`/versions/${version.id}`);
                        } else {
                          console.error('Version missing ID:', version);
                        }
                      }}
                    >
                      <td className="px-6 py-4 whitespace-nowrap text-sm font-medium text-gray-900">
                        {getProductName(version.product_id)}
                      </td>
                      <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-900">
                        {version.version_number}
                      </td>
                      <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-500">
                        {getReleaseTypeLabel(version.release_type)}
                      </td>
                      <td className="px-6 py-4 whitespace-nowrap">
                        <Badge color={getStateBadgeColor(version.state)}>
                          {version.state.replace('_', ' ').replace(/\b\w/g, (l) => l.toUpperCase())}
                        </Badge>
                      </td>
                      <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-500">
                        {formatDate(version.release_date)}
                      </td>
                      <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-500">
                        {version.approved_by || 'N/A'}
                      </td>
                      <td className="px-6 py-4 whitespace-nowrap text-right text-sm font-medium">
                        <Button
                          variant="ghost"
                          size="sm"
                          onClick={(e) => {
                            e.stopPropagation();
                            if (version.id) {
                              navigate(`/versions/${version.id}`);
                            } else {
                              console.error('Version missing ID:', version);
                            }
                          }}
                        >
                          View
                        </Button>
                      </td>
                    </tr>
                  ))}
                </tbody>
              </table>
            </div>

            {/* Pagination */}
            {totalPages > 1 && (
              <div className="mt-6 flex items-center justify-between">
                <div className="flex items-center gap-2">
                  <span className="text-sm text-gray-700">
                    Showing {((filters.page || 1) - 1) * (filters.limit || 25) + 1} to{' '}
                    {Math.min((filters.page || 1) * (filters.limit || 25), total)} of {total} versions
                  </span>
                  <Select
                    value={filters.limit || 25}
                    onChange={(e) => handlePageSizeChange(Number(e.target.value))}
                    className="w-20"
                  >
                    <option value="10">10</option>
                    <option value="25">25</option>
                    <option value="50">50</option>
                    <option value="100">100</option>
                  </Select>
                </div>
                <div className="flex items-center gap-2">
                  <Button
                    variant="secondary"
                    size="sm"
                    onClick={() => handlePageChange((filters.page || 1) - 1)}
                    disabled={(filters.page || 1) === 1}
                  >
                    Previous
                  </Button>
                  <span className="text-sm text-gray-700">
                    Page {filters.page || 1} of {totalPages}
                  </span>
                  <Button
                    variant="secondary"
                    size="sm"
                    onClick={() => handlePageChange((filters.page || 1) + 1)}
                    disabled={(filters.page || 1) >= totalPages}
                  >
                    Next
                  </Button>
                </div>
              </div>
            )}
          </>
        )}
      </Card>
    </div>
  );
};

