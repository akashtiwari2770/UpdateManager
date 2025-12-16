import React, { useState, useEffect } from 'react';
import { useNavigate, useParams } from 'react-router-dom';
import { licenseAllocationsApi } from '@/services/api/license-allocations';
import { tenantsApi } from '@/services/api/tenants';
import { deploymentsApi } from '@/services/api/deployments';
import { AllocateLicenseRequest, CustomerTenant, Deployment } from '@/types';
import { Button, Card, Input, Select, Spinner } from '@/components/ui';

export const AllocateLicenseForm: React.FC = () => {
  const navigate = useNavigate();
  const { customerId, subscriptionId, licenseId } = useParams<{
    customerId: string;
    subscriptionId: string;
    licenseId: string;
  }>();

  const [loading, setLoading] = useState(false);
  const [loadingData, setLoadingData] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [tenants, setTenants] = useState<CustomerTenant[]>([]);
  const [deployments, setDeployments] = useState<Deployment[]>([]);
  const [selectedTenantId, setSelectedTenantId] = useState<string>('');
  const [formData, setFormData] = useState<AllocateLicenseRequest>({
    tenant_id: undefined,
    deployment_id: undefined,
    number_of_seats_allocated: 1,
    notes: '',
  });
  const [validationErrors, setValidationErrors] = useState<Record<string, string>>({});

  useEffect(() => {
    if (customerId) {
      loadTenants();
    }
  }, [customerId]);

  useEffect(() => {
    if (selectedTenantId && customerId) {
      loadDeployments(selectedTenantId);
    } else {
      setDeployments([]);
      setFormData((prev) => ({ ...prev, deployment_id: undefined }));
    }
  }, [selectedTenantId, customerId]);

  const loadTenants = async () => {
    if (!customerId) return;
    try {
      setLoadingData(true);
      const response = await tenantsApi.getAll(customerId);
      setTenants(response?.data || []);
    } catch (err: any) {
      console.error('Error loading tenants:', err);
    } finally {
      setLoadingData(false);
    }
  };

  const loadDeployments = async (tenantId: string) => {
    if (!customerId) return;
    try {
      const response = await deploymentsApi.getAll(customerId, tenantId);
      setDeployments(response?.data || []);
    } catch (err: any) {
      console.error('Error loading deployments:', err);
    }
  };

  const validate = (): boolean => {
    const errors: Record<string, string> = {};

    if (!formData.tenant_id && !formData.deployment_id) {
      errors.allocation_target = 'Either tenant or deployment must be selected';
    }

    if (formData.number_of_seats_allocated < 1) {
      errors.number_of_seats_allocated = 'Number of seats must be at least 1';
    }

    setValidationErrors(errors);
    return Object.keys(errors).length === 0;
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();

    if (!validate() || !customerId || !subscriptionId || !licenseId) {
      return;
    }

    try {
      setLoading(true);
      setError(null);

      const allocationData: AllocateLicenseRequest = {
        tenant_id: formData.tenant_id || undefined,
        deployment_id: formData.deployment_id || undefined,
        number_of_seats_allocated: formData.number_of_seats_allocated,
        notes: formData.notes || undefined,
      };

      await licenseAllocationsApi.allocate(customerId, subscriptionId, licenseId, allocationData);

      navigate(
        `/customers/${customerId}/subscriptions/${subscriptionId}/licenses/${licenseId}`
      );
    } catch (err: any) {
      setError(err.response?.data?.error?.message || 'Failed to allocate license');
      console.error('Error allocating license:', err);
    } finally {
      setLoading(false);
    }
  };

  if (loadingData) {
    return (
      <div className="flex items-center justify-center h-64">
        <Spinner size="lg" />
      </div>
    );
  }

  return (
    <div className="max-w-4xl mx-auto">
      <div className="mb-6">
        <h1 className="text-3xl font-bold text-gray-900">Allocate License</h1>
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
              <label className="block text-sm font-medium text-gray-700 mb-1">Tenant</label>
              <Select
                value={selectedTenantId}
                onChange={(e) => {
                  setSelectedTenantId(e.target.value);
                  setFormData({
                    ...formData,
                    tenant_id: e.target.value || undefined,
                    deployment_id: undefined,
                  });
                }}
              >
                <option value="">Select Tenant (Optional)</option>
                {tenants.map((tenant) => (
                  <option key={tenant.id} value={tenant.tenant_id}>
                    {tenant.name} ({tenant.tenant_id})
                  </option>
                ))}
              </Select>
            </div>

            <div>
              <label className="block text-sm font-medium text-gray-700 mb-1">Deployment</label>
              <Select
                value={formData.deployment_id || ''}
                onChange={(e) =>
                  setFormData({
                    ...formData,
                    deployment_id: e.target.value || undefined,
                  })
                }
                disabled={!selectedTenantId}
              >
                <option value="">Select Deployment (Optional)</option>
                {deployments.map((deployment) => (
                  <option key={deployment.id} value={deployment.deployment_id}>
                    {deployment.deployment_id} ({deployment.product_id})
                  </option>
                ))}
              </Select>
              {!selectedTenantId && (
                <p className="mt-1 text-sm text-gray-500">Select a tenant first</p>
              )}
            </div>

            <div>
              <label className="block text-sm font-medium text-gray-700 mb-1">
                Number of Seats <span className="text-red-500">*</span>
              </label>
              <Input
                type="number"
                min="1"
                value={formData.number_of_seats_allocated}
                onChange={(e) =>
                  setFormData({
                    ...formData,
                    number_of_seats_allocated: parseInt(e.target.value) || 1,
                  })
                }
                required
              />
              {validationErrors.number_of_seats_allocated && (
                <p className="mt-1 text-sm text-red-600">
                  {validationErrors.number_of_seats_allocated}
                </p>
              )}
            </div>
          </div>

          {validationErrors.allocation_target && (
            <div className="bg-yellow-50 border border-yellow-200 text-yellow-700 px-4 py-3 rounded">
              {validationErrors.allocation_target}
            </div>
          )}

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
              onClick={() =>
                navigate(
                  `/customers/${customerId}/subscriptions/${subscriptionId}/licenses/${licenseId}`
                )
              }
            >
              Cancel
            </Button>
            <Button type="submit" variant="primary" disabled={loading}>
              {loading ? 'Allocating...' : 'Allocate License'}
            </Button>
          </div>
        </form>
      </Card>
    </div>
  );
};

