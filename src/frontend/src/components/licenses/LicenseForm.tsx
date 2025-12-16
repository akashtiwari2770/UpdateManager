import React, { useState, useEffect } from 'react';
import { useNavigate, useParams } from 'react-router-dom';
import { licensesApi } from '@/services/api/licenses';
import { productsApi } from '@/services/api/products';
import {
  CreateLicenseRequest,
  UpdateLicenseRequest,
  LicenseType,
  LicenseStatus,
  Product,
} from '@/types';
import { Button, Card, Input, Select, Spinner } from '@/components/ui';

export const LicenseForm: React.FC = () => {
  const navigate = useNavigate();
  const { customerId, subscriptionId, licenseId } = useParams<{
    customerId: string;
    subscriptionId: string;
    licenseId: string;
  }>();
  const isEdit = !!licenseId;

  const [loading, setLoading] = useState(false);
  const [loadingData, setLoadingData] = useState(isEdit);
  const [loadingProducts, setLoadingProducts] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [products, setProducts] = useState<Product[]>([]);
  const [formData, setFormData] = useState<CreateLicenseRequest>({
    license_id: '',
    product_id: '',
    license_type: LicenseType.PERPETUAL,
    number_of_seats: 1,
    start_date: new Date().toISOString().split('T')[0],
    end_date: undefined,
    status: LicenseStatus.ACTIVE,
    notes: '',
  });
  const [validationErrors, setValidationErrors] = useState<Record<string, string>>({});

  useEffect(() => {
    loadProducts();
    if (isEdit && licenseId && customerId && subscriptionId) {
      loadLicense();
    }
  }, [isEdit, licenseId, customerId, subscriptionId]);

  const loadProducts = async () => {
    try {
      setLoadingProducts(true);
      const response = await productsApi.getAll({ limit: 100 });
      setProducts(response?.data || []);
    } catch (err: any) {
      console.error('Error loading products:', err);
    } finally {
      setLoadingProducts(false);
    }
  };

  const loadLicense = async () => {
    if (!customerId || !subscriptionId || !licenseId) return;
    try {
      setLoadingData(true);
      const license = await licensesApi.getById(customerId, subscriptionId, licenseId);
      setFormData({
        license_id: license.license_id,
        product_id: license.product_id,
        license_type: license.license_type,
        number_of_seats: license.number_of_seats,
        start_date: license.start_date.split('T')[0],
        end_date: license.end_date ? license.end_date.split('T')[0] : undefined,
        status: license.status,
        notes: license.notes || '',
      });
    } catch (err: any) {
      setError(err.response?.data?.error?.message || 'Failed to load license');
    } finally {
      setLoadingData(false);
    }
  };

  const validate = (): boolean => {
    const errors: Record<string, string> = {};

    if (!formData.license_id.trim()) {
      errors.license_id = 'License ID is required';
    }

    if (!formData.product_id) {
      errors.product_id = 'Product is required';
    }

    if (formData.number_of_seats < 1) {
      errors.number_of_seats = 'Number of seats must be at least 1';
    }

    if (!formData.start_date) {
      errors.start_date = 'Start date is required';
    }

    if (formData.license_type === LicenseType.TIME_BASED && !formData.end_date) {
      errors.end_date = 'End date is required for time-based licenses';
    }

    if (formData.end_date && formData.end_date < formData.start_date) {
      errors.end_date = 'End date must be after start date';
    }

    setValidationErrors(errors);
    return Object.keys(errors).length === 0;
  };

  const convertDateToISO = (dateStr: string | undefined): string | undefined => {
    if (!dateStr) return undefined;
    // Create date at midnight UTC
    const date = new Date(dateStr + 'T00:00:00.000Z');
    return date.toISOString();
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();

    if (!validate() || !customerId || !subscriptionId) {
      return;
    }

    try {
      setLoading(true);
      setError(null);

      if (isEdit && licenseId) {
        const updateData: UpdateLicenseRequest = {
          license_type: formData.license_type,
          number_of_seats: formData.number_of_seats,
          start_date: convertDateToISO(formData.start_date),
          end_date: convertDateToISO(formData.end_date),
          status: formData.status,
          notes: formData.notes || undefined,
        };
        await licensesApi.update(customerId, subscriptionId, licenseId, updateData);
      } else {
        const createData: CreateLicenseRequest = {
          ...formData,
          start_date: convertDateToISO(formData.start_date)!,
          end_date: convertDateToISO(formData.end_date),
        };
        await licensesApi.assign(customerId, subscriptionId, createData);
      }

      navigate(`/customers/${customerId}/subscriptions/${subscriptionId}`);
    } catch (err: any) {
      setError(err.response?.data?.error?.message || 'Failed to save license');
      console.error('Error saving license:', err);
    } finally {
      setLoading(false);
    }
  };

  if (loadingData || loadingProducts) {
    return (
      <div className="flex items-center justify-center h-64">
        <Spinner size="lg" />
      </div>
    );
  }

  return (
    <div className="max-w-4xl mx-auto">
      <div className="mb-6">
        <h1 className="text-3xl font-bold text-gray-900">
          {isEdit ? 'Edit License' : 'Assign License'}
        </h1>
      </div>

      {error && (
        <div className="mb-4 bg-red-50 border border-red-200 text-red-700 px-4 py-3 rounded">
          {error}
        </div>
      )}

      <Card>
        <form onSubmit={handleSubmit} className="space-y-6">
          <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
            <div>
              <label className="block text-sm font-medium text-gray-700 mb-1">
                License ID <span className="text-red-500">*</span>
              </label>
              <Input
                value={formData.license_id}
                onChange={(e) =>
                  setFormData({ ...formData, license_id: e.target.value })
                }
                disabled={isEdit}
                required
              />
              {validationErrors.license_id && (
                <p className="mt-1 text-sm text-red-600">{validationErrors.license_id}</p>
              )}
            </div>

            <div>
              <label className="block text-sm font-medium text-gray-700 mb-1">
                Product <span className="text-red-500">*</span>
              </label>
              <Select
                value={formData.product_id}
                onChange={(e) => setFormData({ ...formData, product_id: e.target.value })}
                disabled={isEdit}
                required
              >
                <option value="">Select Product</option>
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

            <div>
              <label className="block text-sm font-medium text-gray-700 mb-1">
                License Type <span className="text-red-500">*</span>
              </label>
              <Select
                value={formData.license_type}
                onChange={(e) =>
                  setFormData({
                    ...formData,
                    license_type: e.target.value as LicenseType,
                    end_date: e.target.value === LicenseType.PERPETUAL ? undefined : formData.end_date,
                  })
                }
                required
              >
                <option value={LicenseType.PERPETUAL}>Perpetual</option>
                <option value={LicenseType.TIME_BASED}>Time-based</option>
              </Select>
            </div>

            <div>
              <label className="block text-sm font-medium text-gray-700 mb-1">
                Number of Seats <span className="text-red-500">*</span>
              </label>
              <Input
                type="number"
                min="1"
                value={formData.number_of_seats}
                onChange={(e) =>
                  setFormData({ ...formData, number_of_seats: parseInt(e.target.value) || 1 })
                }
                required
              />
              {validationErrors.number_of_seats && (
                <p className="mt-1 text-sm text-red-600">{validationErrors.number_of_seats}</p>
              )}
            </div>

            <div>
              <label className="block text-sm font-medium text-gray-700 mb-1">
                Start Date <span className="text-red-500">*</span>
              </label>
              <Input
                type="date"
                value={formData.start_date}
                onChange={(e) => setFormData({ ...formData, start_date: e.target.value })}
                required
              />
              {validationErrors.start_date && (
                <p className="mt-1 text-sm text-red-600">{validationErrors.start_date}</p>
              )}
            </div>

            <div>
              <label className="block text-sm font-medium text-gray-700 mb-1">
                End Date {formData.license_type === LicenseType.TIME_BASED && <span className="text-red-500">*</span>}
              </label>
              <Input
                type="date"
                value={formData.end_date || ''}
                onChange={(e) =>
                  setFormData({ ...formData, end_date: e.target.value || undefined })
                }
                disabled={formData.license_type === LicenseType.PERPETUAL}
                required={formData.license_type === LicenseType.TIME_BASED}
              />
              {validationErrors.end_date && (
                <p className="mt-1 text-sm text-red-600">{validationErrors.end_date}</p>
              )}
            </div>

            <div>
              <label className="block text-sm font-medium text-gray-700 mb-1">
                Status <span className="text-red-500">*</span>
              </label>
              <Select
                value={formData.status}
                onChange={(e) =>
                  setFormData({ ...formData, status: e.target.value as LicenseStatus })
                }
                required
              >
                <option value={LicenseStatus.ACTIVE}>Active</option>
                <option value={LicenseStatus.INACTIVE}>Inactive</option>
                <option value={LicenseStatus.EXPIRED}>Expired</option>
                <option value={LicenseStatus.REVOKED}>Revoked</option>
              </Select>
            </div>
          </div>

          <div>
            <label className="block text-sm font-medium text-gray-700 mb-1">Notes</label>
            <textarea
              className="w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-blue-500 focus:border-blue-500"
              rows={3}
              value={formData.notes}
              onChange={(e) => setFormData({ ...formData, notes: e.target.value })}
            />
          </div>

          <div className="flex items-center justify-end gap-4 pt-4 border-t">
            <Button
              type="button"
              variant="secondary"
              onClick={() => navigate(`/customers/${customerId}/subscriptions/${subscriptionId}`)}
            >
              Cancel
            </Button>
            <Button type="submit" variant="primary" disabled={loading}>
              {loading ? 'Saving...' : isEdit ? 'Update' : 'Assign'}
            </Button>
          </div>
        </form>
      </Card>
    </div>
  );
};

