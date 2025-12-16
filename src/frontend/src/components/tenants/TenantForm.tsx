import React, { useState, useEffect } from 'react';
import { useNavigate, useParams } from 'react-router-dom';
import { tenantsApi } from '@/services/api/tenants';
import {
  CreateTenantRequest,
  UpdateTenantRequest,
  TenantStatus,
} from '@/types';
import { Button, Card, Input, Select, Spinner } from '@/components/ui';

export const TenantForm: React.FC = () => {
  const navigate = useNavigate();
  const params = useParams<{ customerId: string; tenantId?: string }>();
  const { customerId, tenantId } = params;
  const isEdit = !!tenantId;

  const [loading, setLoading] = useState(false);
  const [loadingData, setLoadingData] = useState(isEdit);
  const [error, setError] = useState<string | null>(null);
  const [formData, setFormData] = useState<CreateTenantRequest>({
    name: '',
    description: '',
    status: TenantStatus.ACTIVE,
  });
  const [validationErrors, setValidationErrors] = useState<Record<string, string>>({});

  useEffect(() => {
    if (isEdit && tenantId && customerId) {
      loadTenant();
    }
  }, [isEdit, tenantId, customerId]);

  const loadTenant = async () => {
    if (!customerId || !tenantId) return;
    try {
      setLoadingData(true);
      const tenant = await tenantsApi.getById(customerId, tenantId);
      setFormData({
        name: tenant.name,
        description: tenant.description || '',
        status: tenant.status,
      });
    } catch (err: any) {
      setError(err.response?.data?.error?.message || 'Failed to load tenant');
    } finally {
      setLoadingData(false);
    }
  };

  const validate = (): boolean => {
    const errors: Record<string, string> = {};

    if (!formData.name.trim()) {
      errors.name = 'Name is required';
    }

    setValidationErrors(errors);
    return Object.keys(errors).length === 0;
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();

    if (!validate() || !customerId) {
      return;
    }

    try {
      setLoading(true);
      setError(null);

      if (isEdit && tenantId) {
        const updateData: UpdateTenantRequest = {
          name: formData.name,
          description: formData.description,
          status: formData.status,
        };
        await tenantsApi.update(customerId, tenantId, updateData);
      } else {
        await tenantsApi.create(customerId, formData);
      }

      navigate(`/customers/${customerId}`);
    } catch (err: any) {
      console.error('Tenant form error:', err);
      let errorMessage = isEdit ? 'Failed to update tenant' : 'Failed to create tenant';

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
    } finally {
      setLoading(false);
    }
  };

  const handleChange = (field: keyof CreateTenantRequest, value: any) => {
    setFormData((prev) => ({ ...prev, [field]: value }));
    if (validationErrors[field]) {
      setValidationErrors((prev) => {
        const newErrors = { ...prev };
        delete newErrors[field];
        return newErrors;
      });
    }
  };

  if (loadingData) {
    return (
      <div className="flex items-center justify-center h-64">
        <Spinner size="lg" />
      </div>
    );
  }

  if (!customerId) {
    return (
      <div className="space-y-6">
        <Card>
          <div className="text-center py-12">
            <p className="text-red-600 mb-4">Customer ID is required</p>
            <Button variant="primary" onClick={() => navigate('/customers')}>
              Back to Customers
            </Button>
          </div>
        </Card>
      </div>
    );
  }

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <h1 className="text-3xl font-bold text-gray-900">
          {isEdit ? 'Edit Tenant' : 'Create Tenant'}
        </h1>
        <Button variant="ghost" onClick={() => navigate(`/customers/${customerId}`)}>
          Cancel
        </Button>
      </div>

      {error && (
        <div className="bg-red-50 border border-red-200 text-red-700 px-4 py-3 rounded">
          {error}
        </div>
      )}

      <Card>
        <form onSubmit={handleSubmit} className="space-y-6">
          <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
            <Input
              label="Name *"
              value={formData.name}
              onChange={(e) => handleChange('name', e.target.value)}
              error={validationErrors.name}
              required
            />

            <Select
              label="Status *"
              value={formData.status}
              onChange={(e) => handleChange('status', e.target.value as TenantStatus)}
              required
              options={[
                { value: TenantStatus.ACTIVE, label: 'Active' },
                { value: TenantStatus.INACTIVE, label: 'Inactive' },
              ]}
            />

            <div className="md:col-span-2">
              <label className="block text-sm font-medium text-gray-700 mb-1">
                Description
              </label>
              <textarea
                value={formData.description || ''}
                onChange={(e) => handleChange('description', e.target.value)}
                rows={4}
                className="w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-blue-500 focus:border-blue-500"
                placeholder="Optional description for this tenant"
              />
            </div>
          </div>

          <div className="flex justify-end space-x-4 pt-4 border-t">
            <Button
              type="button"
              variant="ghost"
              onClick={() => navigate(`/customers/${customerId}`)}
              disabled={loading}
            >
              Cancel
            </Button>
            <Button type="submit" variant="primary" disabled={loading} isLoading={loading}>
              {isEdit ? 'Update Tenant' : 'Create Tenant'}
            </Button>
          </div>
        </form>
      </Card>
    </div>
  );
};

