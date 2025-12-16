import React, { useState, useEffect } from 'react';
import { useNavigate, useParams } from 'react-router-dom';
import { subscriptionsApi } from '@/services/api/subscriptions';
import {
  CreateSubscriptionRequest,
  UpdateSubscriptionRequest,
  SubscriptionStatus,
} from '@/types';
import { Button, Card, Input, Select, Spinner } from '@/components/ui';

export const SubscriptionForm: React.FC = () => {
  const navigate = useNavigate();
  const { customerId, subscriptionId } = useParams<{ customerId: string; subscriptionId: string }>();
  const isEdit = !!subscriptionId;

  const [loading, setLoading] = useState(false);
  const [loadingData, setLoadingData] = useState(isEdit);
  const [error, setError] = useState<string | null>(null);
  const [formData, setFormData] = useState<CreateSubscriptionRequest>({
    subscription_id: '',
    name: '',
    description: '',
    start_date: new Date().toISOString().split('T')[0],
    end_date: undefined,
    status: SubscriptionStatus.ACTIVE,
    notes: '',
  });
  const [validationErrors, setValidationErrors] = useState<Record<string, string>>({});

  useEffect(() => {
    if (isEdit && subscriptionId && customerId) {
      loadSubscription();
    }
  }, [isEdit, subscriptionId, customerId]);

  const loadSubscription = async () => {
    if (!customerId || !subscriptionId) return;
    try {
      setLoadingData(true);
      const subscription = await subscriptionsApi.getById(customerId, subscriptionId);
      setFormData({
        subscription_id: subscription.subscription_id,
        name: subscription.name || '',
        description: subscription.description || '',
        start_date: subscription.start_date.split('T')[0],
        end_date: subscription.end_date ? subscription.end_date.split('T')[0] : undefined,
        status: subscription.status,
        notes: subscription.notes || '',
      });
    } catch (err: any) {
      setError(err.response?.data?.error?.message || 'Failed to load subscription');
    } finally {
      setLoadingData(false);
    }
  };

  const validate = (): boolean => {
    const errors: Record<string, string> = {};

    if (!formData.subscription_id.trim()) {
      errors.subscription_id = 'Subscription ID is required';
    }

    if (!formData.start_date) {
      errors.start_date = 'Start date is required';
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

    if (!validate() || !customerId) {
      return;
    }

    try {
      setLoading(true);
      setError(null);

      if (isEdit && subscriptionId) {
        const updateData: UpdateSubscriptionRequest = {
          name: formData.name || undefined,
          description: formData.description || undefined,
          start_date: convertDateToISO(formData.start_date),
          end_date: convertDateToISO(formData.end_date),
          status: formData.status,
          notes: formData.notes || undefined,
        };
        await subscriptionsApi.update(customerId, subscriptionId, updateData);
      } else {
        const createData: CreateSubscriptionRequest = {
          ...formData,
          start_date: convertDateToISO(formData.start_date)!,
          end_date: convertDateToISO(formData.end_date),
        };
        await subscriptionsApi.create(customerId, createData);
      }

      navigate(`/customers/${customerId}`);
    } catch (err: any) {
      setError(err.response?.data?.error?.message || 'Failed to save subscription');
      console.error('Error saving subscription:', err);
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
        <h1 className="text-3xl font-bold text-gray-900">
          {isEdit ? 'Edit Subscription' : 'Create Subscription'}
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
                Subscription ID <span className="text-red-500">*</span>
              </label>
              <Input
                value={formData.subscription_id}
                onChange={(e) =>
                  setFormData({ ...formData, subscription_id: e.target.value })
                }
                disabled={isEdit}
                required
              />
              {validationErrors.subscription_id && (
                <p className="mt-1 text-sm text-red-600">{validationErrors.subscription_id}</p>
              )}
            </div>

            <div>
              <label className="block text-sm font-medium text-gray-700 mb-1">Name</label>
              <Input
                value={formData.name}
                onChange={(e) => setFormData({ ...formData, name: e.target.value })}
              />
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
              <label className="block text-sm font-medium text-gray-700 mb-1">End Date</label>
              <Input
                type="date"
                value={formData.end_date || ''}
                onChange={(e) =>
                  setFormData({ ...formData, end_date: e.target.value || undefined })
                }
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
                  setFormData({ ...formData, status: e.target.value as SubscriptionStatus })
                }
                required
              >
                <option value={SubscriptionStatus.ACTIVE}>Active</option>
                <option value={SubscriptionStatus.INACTIVE}>Inactive</option>
                <option value={SubscriptionStatus.EXPIRED}>Expired</option>
                <option value={SubscriptionStatus.SUSPENDED}>Suspended</option>
              </Select>
            </div>
          </div>

          <div>
            <label className="block text-sm font-medium text-gray-700 mb-1">Description</label>
            <textarea
              className="w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-blue-500 focus:border-blue-500"
              rows={3}
              value={formData.description}
              onChange={(e) => setFormData({ ...formData, description: e.target.value })}
            />
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
              onClick={() => navigate(`/customers/${customerId}`)}
            >
              Cancel
            </Button>
            <Button type="submit" variant="primary" disabled={loading}>
              {loading ? 'Saving...' : isEdit ? 'Update' : 'Create'}
            </Button>
          </div>
        </form>
      </Card>
    </div>
  );
};

