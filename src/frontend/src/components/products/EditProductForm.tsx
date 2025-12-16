import React, { useState, useEffect } from 'react';
import { useParams, useNavigate } from 'react-router-dom';
import { productsApi } from '@/services/api/products';
import { Product, CreateProductRequest, ProductType } from '@/types';
import { Button, Card, Input, Select, Alert, Spinner } from '@/components/ui';

export const EditProductForm: React.FC = () => {
  const { id } = useParams<{ id: string }>();
  const navigate = useNavigate();
  const [product, setProduct] = useState<Product | null>(null);
  const [loading, setLoading] = useState(true);
  const [saving, setSaving] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [formData, setFormData] = useState<Partial<CreateProductRequest>>({});
  const [validationErrors, setValidationErrors] = useState<
    Record<string, string>
  >({});

  useEffect(() => {
    if (id) {
      loadProduct();
    }
  }, [id]);

  const loadProduct = async () => {
    try {
      setLoading(true);
      setError(null);
      const data = await productsApi.getById(id!);
      setProduct(data);
      setFormData({
        name: data.name,
        type: data.type,
        description: data.description || '',
        vendor: data.vendor || '',
      });
    } catch (err: any) {
      setError(err.response?.data?.error?.message || 'Failed to load product');
      console.error('Error loading product:', err);
    } finally {
      setLoading(false);
    }
  };

  const validate = (): boolean => {
    const errors: Record<string, string> = {};

    if (!formData.name?.trim()) {
      errors.name = 'Product Name is required';
    }

    if (!formData.type) {
      errors.type = 'Product Type is required';
    }

    setValidationErrors(errors);
    return Object.keys(errors).length === 0;
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();

    if (!validate() || !id) {
      return;
    }

    try {
      setSaving(true);
      setError(null);
      const updatedProduct = await productsApi.update(id, formData);
      navigate(`/products/${updatedProduct.id}`);
    } catch (err: any) {
      const errorMessage =
        err.response?.data?.error?.message || 'Failed to update product';
      setError(errorMessage);
    } finally {
      setSaving(false);
    }
  };

  const handleChange = (
    field: keyof CreateProductRequest,
    value: string | ProductType
  ) => {
    setFormData((prev) => ({ ...prev, [field]: value }));
    // Clear validation error for this field
    if (validationErrors[field]) {
      setValidationErrors((prev) => {
        const newErrors = { ...prev };
        delete newErrors[field];
        return newErrors;
      });
    }
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
        <h1 className="text-3xl font-bold text-gray-900">Edit Product</h1>
        <Button variant="ghost" onClick={() => navigate(`/products/${id}`)}>
          Cancel
        </Button>
      </div>

      <Card title="Product Details">
        <form onSubmit={handleSubmit} className="space-y-6">
          {error && (
            <Alert variant="error" title="Error" onClose={() => setError(null)}>
              {error}
            </Alert>
          )}

          <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
            <div>
              <label className="block text-sm font-medium text-gray-700 mb-2">
                Product ID
              </label>
              <input
                type="text"
                className="w-full px-3 py-2 border border-gray-300 rounded-lg bg-gray-50 text-gray-500 cursor-not-allowed"
                value={product.product_id}
                disabled
              />
              <p className="mt-1 text-xs text-gray-500">
                Product ID cannot be changed after creation
              </p>
            </div>

            <Input
              label="Product Name"
              required
              value={formData.name || ''}
              onChange={(e) => handleChange('name', e.target.value)}
              error={validationErrors.name}
              placeholder="e.g., My Awesome Product"
            />
          </div>

          <div>
            <Select
              label="Product Type"
              required
              value={formData.type || ProductType.SERVER}
              onChange={(e) => handleChange('type', e.target.value as ProductType)}
              error={validationErrors.type}
              options={[
                { value: ProductType.SERVER, label: 'Server' },
                { value: ProductType.CLIENT, label: 'Client' },
              ]}
            />
          </div>

          <div>
            <label className="block text-sm font-medium text-gray-700 mb-2">
              Description
            </label>
            <textarea
              className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-blue-500"
              rows={4}
              value={formData.description || ''}
              onChange={(e) => handleChange('description', e.target.value)}
              placeholder="Optional description of the product..."
            />
          </div>

          <div className="flex items-center justify-end gap-4 pt-4 border-t">
            <Button
              type="button"
              variant="secondary"
              onClick={() => navigate(`/products/${id}`)}
              disabled={saving}
            >
              Cancel
            </Button>
            <Button type="submit" variant="primary" isLoading={saving}>
              Save Changes
            </Button>
          </div>
        </form>
      </Card>
    </div>
  );
};

