import React, { useState, useEffect } from 'react';
import { useNavigate, useParams } from 'react-router-dom';
import { versionsApi } from '@/services/api/versions';
import { productsApi } from '@/services/api/products';
import { CreateVersionRequest, Product, ReleaseType } from '@/types';
import { Button, Card, Input, Select, Alert, Spinner } from '@/components/ui';

export const CreateVersionForm: React.FC = () => {
  const navigate = useNavigate();
  const { productId } = useParams<{ productId?: string }>();
  const [products, setProducts] = useState<Product[]>([]);
  const [loading, setLoading] = useState(false);
  const [saving, setSaving] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [validationErrors, setValidationErrors] = useState<Record<string, string>>({});
  
  const [formData, setFormData] = useState<CreateVersionRequest>({
    version_number: '',
    release_date: new Date().toISOString().split('T')[0],
    release_type: ReleaseType.FEATURE,
    eol_date: '',
    min_server_version: '',
    max_server_version: '',
    recommended_server_version: '',
  });

  const [selectedProductId, setSelectedProductId] = useState<string>(productId || '');

  useEffect(() => {
    loadProducts();
    if (productId) {
      setSelectedProductId(productId);
    }
  }, [productId]);

  const loadProducts = async () => {
    try {
      const response = await productsApi.getAll({ page: 1, limit: 1000, is_active: true });
      setProducts(response.data || []);
    } catch (err) {
      console.error('Error loading products:', err);
    }
  };

  const validateVersionNumber = (version: string): boolean => {
    // Semantic versioning: major.minor.patch
    const semverRegex = /^\d+\.\d+\.\d+(-[a-zA-Z0-9]+)?(\+[a-zA-Z0-9]+)?$/;
    return semverRegex.test(version);
  };

  const validate = (): boolean => {
    const errors: Record<string, string> = {};

    if (!selectedProductId) {
      errors.product_id = 'Product is required';
    }

    if (!formData.version_number.trim()) {
      errors.version_number = 'Version number is required';
    } else if (!validateVersionNumber(formData.version_number)) {
      errors.version_number = 'Version number must follow semantic versioning (e.g., 1.0.0)';
    }

    if (!formData.release_date) {
      errors.release_date = 'Release date is required';
    }

    if (!formData.release_type) {
      errors.release_type = 'Release type is required';
    }

    setValidationErrors(errors);
    return Object.keys(errors).length === 0;
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    
    if (!validate()) {
      return;
    }

    if (!selectedProductId) {
      setError('Please select a product');
      return;
    }

    try {
      setSaving(true);
      setError(null);
      
      // Convert date strings (YYYY-MM-DD) to ISO datetime strings
      // Date input returns YYYY-MM-DD, we need to convert to ISO datetime
      const convertDateToISO = (dateStr: string): string => {
        if (!dateStr) return new Date().toISOString();
        // Create date at midnight UTC
        const date = new Date(dateStr + 'T00:00:00.000Z');
        return date.toISOString();
      };
      
      const payload: CreateVersionRequest = {
        ...formData,
        release_date: convertDateToISO(formData.release_date),
        eol_date: formData.eol_date ? convertDateToISO(formData.eol_date) : undefined,
      };
      
      const version = await versionsApi.create(selectedProductId, payload);
      navigate(`/versions/${version.id}`);
    } catch (err: any) {
      const errorMessage = err.response?.data?.error?.message || 'Failed to create version';
      setError(errorMessage);
      
      if (errorMessage.includes('already exists') || errorMessage.includes('duplicate')) {
        setValidationErrors({ version_number: 'This version number already exists for this product' });
      }
    } finally {
      setSaving(false);
    }
  };

  const handleChange = (field: keyof CreateVersionRequest, value: any) => {
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

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <h1 className="text-3xl font-bold text-gray-900">Create Version</h1>
        <Button variant="ghost" onClick={() => navigate('/versions')}>
          ← Back to Versions
        </Button>
      </div>

      {error && (
        <Alert variant="error" title="Error" onClose={() => setError(null)}>
          {error}
        </Alert>
      )}

      <Card>
        <form onSubmit={handleSubmit} className="space-y-6">
          {/* Product Selection */}
          <div>
            <label className="block text-sm font-medium text-gray-700 mb-2">
              Product <span className="text-red-500">*</span>
            </label>
            <Select
              value={selectedProductId}
              onChange={(e) => {
                setSelectedProductId(e.target.value);
                if (validationErrors.product_id) {
                  setValidationErrors((prev) => {
                    const newErrors = { ...prev };
                    delete newErrors.product_id;
                    return newErrors;
                  });
                }
              }}
              disabled={!!productId}
            >
              <option value="">Select a product</option>
              {products.map((product) => (
                <option key={product.id} value={product.product_id}>
                  {product.name} ({product.product_id})
                </option>
              ))}
            </Select>
            {validationErrors.product_id && (
              <p className="mt-1 text-sm text-red-600">{validationErrors.product_id}</p>
            )}
          </div>

          {/* Version Number */}
          <div>
            <label className="block text-sm font-medium text-gray-700 mb-2">
              Version Number <span className="text-red-500">*</span>
            </label>
            <Input
              type="text"
              placeholder="e.g., 1.0.0"
              value={formData.version_number}
              onChange={(e) => handleChange('version_number', e.target.value)}
            />
            {validationErrors.version_number && (
              <p className="mt-1 text-sm text-red-600">{validationErrors.version_number}</p>
            )}
            <p className="mt-1 text-sm text-gray-500">Follow semantic versioning (major.minor.patch)</p>
          </div>

          {/* Release Type */}
          <div>
            <label className="block text-sm font-medium text-gray-700 mb-2">
              Release Type <span className="text-red-500">*</span>
            </label>
            <Select
              value={formData.release_type}
              onChange={(e) => handleChange('release_type', e.target.value as ReleaseType)}
            >
              <option value={ReleaseType.SECURITY}>Security</option>
              <option value={ReleaseType.FEATURE}>Feature</option>
              <option value={ReleaseType.MAINTENANCE}>Maintenance</option>
              <option value={ReleaseType.MAJOR}>Major</option>
            </Select>
            {validationErrors.release_type && (
              <p className="mt-1 text-sm text-red-600">{validationErrors.release_type}</p>
            )}
          </div>

          {/* Release Date */}
          <div>
            <label className="block text-sm font-medium text-gray-700 mb-2">
              Release Date <span className="text-red-500">*</span>
            </label>
            <Input
              type="date"
              value={formData.release_date}
              onChange={(e) => handleChange('release_date', e.target.value)}
            />
            {validationErrors.release_date && (
              <p className="mt-1 text-sm text-red-600">{validationErrors.release_date}</p>
            )}
          </div>

          {/* EOL Date (Optional) */}
          <div>
            <label className="block text-sm font-medium text-gray-700 mb-2">
              End of Life Date (Optional)
            </label>
            <Input
              type="date"
              value={formData.eol_date || ''}
              onChange={(e) => handleChange('eol_date', e.target.value || undefined)}
            />
          </div>

          {/* Server Version Requirements */}
          <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
            <div>
              <label className="block text-sm font-medium text-gray-700 mb-2">
                Min Server Version
              </label>
              <Input
                type="text"
                placeholder="e.g., 1.0.0"
                value={formData.min_server_version || ''}
                onChange={(e) => handleChange('min_server_version', e.target.value || undefined)}
              />
            </div>
            <div>
              <label className="block text-sm font-medium text-gray-700 mb-2">
                Max Server Version
              </label>
              <Input
                type="text"
                placeholder="e.g., 2.0.0"
                value={formData.max_server_version || ''}
                onChange={(e) => handleChange('max_server_version', e.target.value || undefined)}
              />
            </div>
            <div>
              <label className="block text-sm font-medium text-gray-700 mb-2">
                Recommended Server Version
              </label>
              <Input
                type="text"
                placeholder="e.g., 1.5.0"
                value={formData.recommended_server_version || ''}
                onChange={(e) => handleChange('recommended_server_version', e.target.value || undefined)}
              />
            </div>
          </div>

          {/* Actions */}
          <div className="flex items-center justify-end gap-4 pt-4 border-t">
            <Button
              type="button"
              variant="secondary"
              onClick={() => navigate('/versions')}
              disabled={saving}
            >
              Cancel
            </Button>
            <Button type="submit" disabled={saving}>
              {saving ? <Spinner size="sm" /> : 'Create Version'}
            </Button>
          </div>
        </form>
      </Card>
    </div>
  );
};

