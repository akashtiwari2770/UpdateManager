import React, { useState, useEffect } from 'react';
import { useNavigate } from 'react-router-dom';
import { productsApi } from '@/services/api/products';
import { Product, ProductType, ListProductsQuery } from '@/types';
import { Button, Card, Badge, Spinner, Input, Select } from '@/components/ui';

export const ProductsList: React.FC = () => {
  const navigate = useNavigate();
  const [products, setProducts] = useState<Product[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  
  // Filters and pagination - default to showing only active products
  const [filters, setFilters] = useState<ListProductsQuery>({
    page: 1,
    limit: 25,
    is_active: true, // Show only active products by default
  });
  const [searchTerm, setSearchTerm] = useState('');
  const [totalPages, setTotalPages] = useState(1);
  const [total, setTotal] = useState(0);

  // Sorting
  const [sortField, setSortField] = useState<string>('created_at');
  const [sortOrder, setSortOrder] = useState<'asc' | 'desc'>('desc');

  useEffect(() => {
    loadProducts();
  }, [filters, sortField, sortOrder]);

  const loadProducts = async () => {
    try {
      setLoading(true);
      setError(null);
      
      const query: ListProductsQuery = {
        ...filters,
        page: filters.page || 1,
        limit: filters.limit || 25,
      };

      const response = await productsApi.getAll(query);
      const productsData = response?.data || [];
      // Ensure all products have required fields
      const normalizedProducts = productsData.map((p: Product) => ({
        ...p,
        product_id: p.product_id || '',
        name: p.name || '',
      }));
      setProducts(normalizedProducts);
      setTotalPages(response?.pagination?.total_pages || 1);
      setTotal(response?.pagination?.total || 0);
    } catch (err: any) {
      setError(err.response?.data?.error?.message || 'Failed to load products');
      console.error('Error loading products:', err);
    } finally {
      setLoading(false);
    }
  };

  const handleSearch = (value: string) => {
    setSearchTerm(value);
    // TODO: Implement search API call when backend supports it
    // For now, filter client-side
    if (value.trim()) {
      const filtered = products.filter(
        (p) =>
          p.name.toLowerCase().includes(value.toLowerCase()) ||
          p.product_id.toLowerCase().includes(value.toLowerCase())
      );
      setProducts(filtered);
    } else {
      loadProducts();
    }
  };

  const handleFilterChange = (key: keyof ListProductsQuery, value: any) => {
    setFilters((prev) => ({
      ...prev,
      [key]: value,
      page: 1, // Reset to first page on filter change
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
    return new Date(dateString).toLocaleDateString('en-US', {
      year: 'numeric',
      month: 'short',
      day: 'numeric',
    });
  };

  if (loading && products.length === 0) {
    return (
      <div className="flex items-center justify-center h-64">
        <Spinner size="lg" />
      </div>
    );
  }

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <h1 className="text-3xl font-bold text-gray-900">Products</h1>
        <Button onClick={() => navigate('/products/new')} variant="primary">
          Create Product
        </Button>
      </div>

      <Card>
        <div className="space-y-4">
          {/* Filters */}
          <div className="grid grid-cols-1 md:grid-cols-4 gap-4">
            <div className="md:col-span-2">
              <Input
                label="Search"
                placeholder="Search by name or ID..."
                value={searchTerm}
                onChange={(e) => handleSearch(e.target.value)}
              />
            </div>
            <Select
              label="Type"
              value={filters.type || ''}
              onChange={(e) =>
                handleFilterChange('type', e.target.value || undefined)
              }
              options={[
                { value: '', label: 'All Types' },
                { value: ProductType.SERVER, label: 'Server' },
                { value: ProductType.CLIENT, label: 'Client' },
              ]}
            />
            <Select
              label="Status"
              value={filters.is_active !== undefined ? String(filters.is_active) : ''}
              onChange={(e) =>
                handleFilterChange(
                  'is_active',
                  e.target.value === '' ? undefined : e.target.value === 'true'
                )
              }
              options={[
                { value: '', label: 'All Status' },
                { value: 'true', label: 'Active' },
                { value: 'false', label: 'Inactive' },
              ]}
            />
          </div>

          {error && (
            <div className="bg-red-50 border border-red-200 text-red-800 px-4 py-3 rounded-lg">
              {error}
            </div>
          )}

          {/* Products Table */}
          {products.length === 0 ? (
            <div className="text-center py-12">
              <p className="text-gray-500 text-lg">No products found</p>
              <Button
                onClick={() => navigate('/products/new')}
                variant="primary"
                className="mt-4"
              >
                Create First Product
              </Button>
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
                        <div className="flex items-center gap-2">
                          ID {getSortIcon('product_id')}
                        </div>
                      </th>
                      <th
                        className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider cursor-pointer hover:bg-gray-100"
                        onClick={() => handleSort('name')}
                      >
                        <div className="flex items-center gap-2">
                          Name {getSortIcon('name')}
                        </div>
                      </th>
                      <th
                        className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider cursor-pointer hover:bg-gray-100"
                        onClick={() => handleSort('type')}
                      >
                        <div className="flex items-center gap-2">
                          Type {getSortIcon('type')}
                        </div>
                      </th>
                      <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                        Status
                      </th>
                      <th
                        className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider cursor-pointer hover:bg-gray-100"
                        onClick={() => handleSort('updated_at')}
                      >
                        <div className="flex items-center gap-2">
                          Last Updated {getSortIcon('updated_at')}
                        </div>
                      </th>
                      <th className="px-6 py-3 text-right text-xs font-medium text-gray-500 uppercase tracking-wider">
                        Actions
                      </th>
                    </tr>
                  </thead>
                  <tbody className="bg-white divide-y divide-gray-200">
                    {products.map((product) => (
                      <tr
                        key={product.id}
                        className="hover:bg-gray-50 cursor-pointer"
                        onClick={() => {
                          if (product.id) {
                            navigate(`/products/${product.id}`);
                          } else {
                            console.error('Product missing ID:', product);
                          }
                        }}
                      >
                        <td className="px-6 py-4 whitespace-nowrap text-sm font-mono text-gray-900">
                          {product.product_id || <span className="text-gray-400 italic">N/A</span>}
                        </td>
                        <td className="px-6 py-4 whitespace-nowrap text-sm font-medium text-gray-900">
                          {product.name || <span className="text-gray-400 italic">Unnamed</span>}
                        </td>
                        <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-500">
                          <Badge
                            variant={
                              product.type === ProductType.SERVER
                                ? 'info'
                                : 'default'
                            }
                          >
                            {product.type}
                          </Badge>
                        </td>
                        <td className="px-6 py-4 whitespace-nowrap text-sm">
                          <Badge variant={product.is_active ? 'success' : 'default'}>
                            {product.is_active ? 'Active' : 'Inactive'}
                          </Badge>
                        </td>
                        <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-500">
                          {formatDate(product.updated_at)}
                        </td>
                        <td className="px-6 py-4 whitespace-nowrap text-right text-sm font-medium">
                          <div
                            className="flex items-center justify-end gap-2"
                            onClick={(e) => e.stopPropagation()}
                          >
                            <Button
                              variant="ghost"
                              size="sm"
                              onClick={(e) => {
                                e.stopPropagation();
                                if (product.id) {
                                  navigate(`/products/${product.id}`);
                                } else {
                                  console.error('Product missing ID:', product);
                                }
                              }}
                            >
                              View
                            </Button>
                            <Button
                              variant="ghost"
                              size="sm"
                              onClick={(e) => {
                                e.stopPropagation();
                                if (product.id) {
                                  navigate(`/products/${product.id}/edit`);
                                } else {
                                  console.error('Product missing ID:', product);
                                }
                              }}
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
              <div className="flex items-center justify-between border-t border-gray-200 px-4 py-3">
                <div className="flex items-center gap-4">
                  <span className="text-sm text-gray-700">
                    Showing {((filters.page || 1) - 1) * (filters.limit || 25) + 1} to{' '}
                    {Math.min((filters.page || 1) * (filters.limit || 25), total)} of {total}{' '}
                    products
                  </span>
                  <div className="w-32">
                    <Select
                      value={String(filters.limit || 25)}
                      onChange={(e) => handlePageSizeChange(Number(e.target.value))}
                      options={[
                        { value: '10', label: '10 per page' },
                        { value: '25', label: '25 per page' },
                        { value: '50', label: '50 per page' },
                        { value: '100', label: '100 per page' },
                      ]}
                    />
                  </div>
                </div>
                <div className="flex items-center gap-2">
                  <Button
                    variant="ghost"
                    size="sm"
                    disabled={filters.page === 1}
                    onClick={() => handlePageChange((filters.page || 1) - 1)}
                  >
                    Previous
                  </Button>
                  <span className="text-sm text-gray-700">
                    Page {filters.page || 1} of {totalPages}
                  </span>
                  <Button
                    variant="ghost"
                    size="sm"
                    disabled={(filters.page || 1) >= totalPages}
                    onClick={() => handlePageChange((filters.page || 1) + 1)}
                  >
                    Next
                  </Button>
                </div>
              </div>
            </>
          )}
        </div>
      </Card>
    </div>
  );
};

