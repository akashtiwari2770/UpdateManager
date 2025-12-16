import React, { useState, useEffect } from 'react';
import { useParams, useNavigate } from 'react-router-dom';
import { productsApi } from '@/services/api/products';
import { versionsApi } from '@/services/api/versions';
import { Product, ProductType, Version, VersionState } from '@/types';
import { Button, Card, Badge, Spinner, Alert, Modal } from '@/components/ui';

export const ProductDetails: React.FC = () => {
  const { id } = useParams<{ id: string }>();
  const navigate = useNavigate();
  const [product, setProduct] = useState<Product | null>(null);
  const [versions, setVersions] = useState<Version[]>([]);
  const [loadingVersions, setLoadingVersions] = useState(false);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [showDeleteDialog, setShowDeleteDialog] = useState(false);
  const [deleting, setDeleting] = useState(false);

  useEffect(() => {
    if (id) {
      loadProduct();
      loadVersions();
    }
  }, [id]);

  const loadProduct = async () => {
    try {
      setLoading(true);
      setError(null);
      const data = await productsApi.getById(id!);
      setProduct(data);
    } catch (err: any) {
      setError(err.response?.data?.error?.message || 'Failed to load product');
      console.error('Error loading product:', err);
    } finally {
      setLoading(false);
    }
  };

  const loadVersions = async () => {
    try {
      setLoadingVersions(true);
      const productData = await productsApi.getById(id!);
      const versionsData = await versionsApi.getByProduct(productData.product_id);
      // Ensure versions is always an array
      setVersions(Array.isArray(versionsData) ? versionsData : []);
    } catch (err: any) {
      console.error('Error loading versions:', err);
      setVersions([]); // Set empty array on error
    } finally {
      setLoadingVersions(false);
    }
  };

  const handleDelete = async () => {
    if (!product) return;

    try {
      setDeleting(true);
      setError(null);
      await productsApi.delete(product.id);
      // Navigate to products list after successful delete
      // The list will automatically reload and filter out inactive products
      navigate('/products');
    } catch (err: any) {
      console.error('Error deleting product:', err);
      let errorMessage = 'Failed to delete product';
      
      if (err.response?.data) {
        if (err.response.data.error) {
          errorMessage = err.response.data.error.message || err.response.data.error;
        } else if (typeof err.response.data === 'string') {
          errorMessage = err.response.data;
        } else if (err.response.data.message) {
          errorMessage = err.response.data.message;
        }
      } else if (err.message) {
        errorMessage = err.message;
      }
      
      setError(errorMessage);
      setShowDeleteDialog(false); // Close dialog on error so user can see the error message
    } finally {
      setDeleting(false);
    }
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

  if (error && !product) {
    return (
      <div className="space-y-4">
        <Button variant="ghost" onClick={() => navigate('/products')}>
          ← Back to Products
        </Button>
        <Alert variant="error" title="Error">
          {error}
        </Alert>
      </div>
    );
  }

  if (!product) {
    return (
      <div className="space-y-4">
        <Button variant="ghost" onClick={() => navigate('/products')}>
          ← Back to Products
        </Button>
        <Alert variant="info" title="Not Found">
          Product not found
        </Alert>
      </div>
    );
  }

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <div className="flex items-center gap-4">
          <Button variant="ghost" onClick={() => navigate('/products')}>
            ← Back
          </Button>
          <h1 className="text-3xl font-bold text-gray-900">{product.name}</h1>
        </div>
        <div className="flex items-center gap-2">
          <Button
            variant="secondary"
            onClick={() => navigate(`/products/${product.id}/edit`)}
          >
            Edit Product
          </Button>
          <Button
            variant="danger"
            onClick={() => setShowDeleteDialog(true)}
          >
            Delete Product
          </Button>
        </div>
      </div>

      {error && (
        <Alert variant="error" title="Error" onClose={() => setError(null)}>
          {error}
        </Alert>
      )}

      <div className="grid grid-cols-1 lg:grid-cols-3 gap-6">
        {/* Main Product Information */}
        <div className="lg:col-span-2 space-y-6">
          <Card title="Product Information">
            <dl className="grid grid-cols-1 gap-4 sm:grid-cols-2">
              <div>
                <dt className="text-sm font-medium text-gray-500">Product ID</dt>
                <dd className="mt-1 text-sm text-gray-900 font-mono">
                  {product.product_id}
                </dd>
              </div>
              <div>
                <dt className="text-sm font-medium text-gray-500">Product Name</dt>
                <dd className="mt-1 text-sm text-gray-900">{product.name}</dd>
              </div>
              <div>
                <dt className="text-sm font-medium text-gray-500">Product Type</dt>
                <dd className="mt-1">
                  <Badge
                    variant={
                      product.type === ProductType.SERVER ? 'info' : 'default'
                    }
                  >
                    {product.type}
                  </Badge>
                </dd>
              </div>
              <div>
                <dt className="text-sm font-medium text-gray-500">Status</dt>
                <dd className="mt-1">
                  <Badge variant={product.is_active ? 'success' : 'default'}>
                    {product.is_active ? 'Active' : 'Inactive'}
                  </Badge>
                </dd>
              </div>
              <div>
                <dt className="text-sm font-medium text-gray-500">Created</dt>
                <dd className="mt-1 text-sm text-gray-900">
                  {formatDate(product.created_at)}
                </dd>
              </div>
              <div>
                <dt className="text-sm font-medium text-gray-500">Last Updated</dt>
                <dd className="mt-1 text-sm text-gray-900">
                  {formatDate(product.updated_at)}
                </dd>
              </div>
            </dl>
            {product.description && (
              <div className="mt-4">
                <dt className="text-sm font-medium text-gray-500">Description</dt>
                <dd className="mt-1 text-sm text-gray-900 whitespace-pre-wrap">
                  {product.description}
                </dd>
              </div>
            )}
          </Card>

          <Card title="Versions">
            <div className="flex items-center justify-between mb-4">
              <h3 className="text-lg font-medium text-gray-900">
                Versions ({versions.length})
              </h3>
              <Button
                variant="primary"
                onClick={() => navigate(`/products/${product.product_id}/versions/new`)}
              >
                Create Version
              </Button>
            </div>

            {loadingVersions ? (
              <div className="flex justify-center py-8">
                <Spinner />
              </div>
            ) : versions.length === 0 ? (
              <div className="text-center py-8">
                <p className="text-gray-500 mb-4">No versions found for this product.</p>
                <Button
                  variant="secondary"
                  onClick={() => navigate(`/products/${product.product_id}/versions/new`)}
                >
                  Create First Version
                </Button>
              </div>
            ) : (
              <div className="overflow-x-auto">
                <table className="min-w-full divide-y divide-gray-200">
                  <thead className="bg-gray-50">
                    <tr>
                      <th className="px-4 py-3 text-left text-xs font-medium text-gray-500 uppercase">Version</th>
                      <th className="px-4 py-3 text-left text-xs font-medium text-gray-500 uppercase">Type</th>
                      <th className="px-4 py-3 text-left text-xs font-medium text-gray-500 uppercase">State</th>
                      <th className="px-4 py-3 text-left text-xs font-medium text-gray-500 uppercase">Release Date</th>
                      <th className="px-4 py-3 text-right text-xs font-medium text-gray-500 uppercase">Actions</th>
                    </tr>
                  </thead>
                  <tbody className="bg-white divide-y divide-gray-200">
                    {versions.map((version) => (
                      <tr key={version.id} className="hover:bg-gray-50">
                        <td className="px-4 py-3 whitespace-nowrap">
                          <button
                            onClick={() => navigate(`/versions/${version.id}`)}
                            className="text-sm font-medium text-blue-600 hover:text-blue-900"
                          >
                            {version.version_number}
                          </button>
                        </td>
                        <td className="px-4 py-3 whitespace-nowrap text-sm text-gray-500">
                          {version.release_type}
                        </td>
                        <td className="px-4 py-3 whitespace-nowrap">
                          <Badge
                            color={
                              version.state === VersionState.RELEASED
                                ? 'green'
                                : version.state === VersionState.APPROVED
                                ? 'blue'
                                : version.state === VersionState.PENDING_REVIEW
                                ? 'yellow'
                                : 'gray'
                            }
                          >
                            {version.state.replace('_', ' ').replace(/\b\w/g, (l) => l.toUpperCase())}
                          </Badge>
                        </td>
                        <td className="px-4 py-3 whitespace-nowrap text-sm text-gray-500">
                          {new Date(version.release_date).toLocaleDateString()}
                        </td>
                        <td className="px-4 py-3 whitespace-nowrap text-right text-sm font-medium">
                          <Button
                            variant="ghost"
                            size="sm"
                            onClick={() => navigate(`/versions/${version.id}`)}
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

            <div className="mt-4">
              <Button
                variant="secondary"
                onClick={() => navigate('/versions', { state: { productId: product.product_id } })}
              >
                View All Versions
              </Button>
            </div>
          </Card>
        </div>

        {/* Sidebar Actions */}
        <div className="space-y-6">
          <Card title="Quick Actions">
            <div className="space-y-2">
              <Button
                variant="primary"
                className="w-full"
                onClick={() => navigate(`/products/${product.product_id}/versions/new`)}
              >
                Create Version
              </Button>
              <Button
                variant="secondary"
                className="w-full"
                onClick={() => navigate('/versions', { state: { productId: product.product_id } })}
              >
                View All Versions
              </Button>
            </div>
          </Card>
        </div>
      </div>

      {/* Delete Confirmation Dialog */}
      <Modal
        isOpen={showDeleteDialog}
        onClose={() => setShowDeleteDialog(false)}
        title="Delete Product"
        size="md"
        footer={
          <div className="flex items-center justify-end gap-2">
            <Button
              variant="secondary"
              onClick={() => setShowDeleteDialog(false)}
              disabled={deleting}
            >
              Cancel
            </Button>
            <Button
              variant="danger"
              onClick={handleDelete}
              isLoading={deleting}
            >
              Delete
            </Button>
          </div>
        }
      >
        <div className="space-y-4">
          <p className="text-gray-700">
            Are you sure you want to delete <strong>{product.name}</strong>?
          </p>
          <Alert variant="warning" title="Warning">
            This action cannot be undone. All associated versions and data will be
            affected.
          </Alert>
        </div>
      </Modal>
    </div>
  );
};

