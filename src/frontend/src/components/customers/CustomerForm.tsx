import React, { useState, useEffect } from 'react';
import { useNavigate, useParams } from 'react-router-dom';
import { customersApi } from '@/services/api/customers';
import {
  CreateCustomerRequest,
  UpdateCustomerRequest,
  CustomerStatus,
  NotificationPreferences,
} from '@/types';
import { Button, Card, Input, Select, Spinner } from '@/components/ui';

interface CustomerFormProps {
  customerId?: string;
}

export const CustomerForm: React.FC<CustomerFormProps> = ({ customerId }) => {
  const navigate = useNavigate();
  const params = useParams();
  const id = customerId || params.id;
  const isEdit = !!id;

  const [loading, setLoading] = useState(false);
  const [loadingData, setLoadingData] = useState(isEdit);
  const [error, setError] = useState<string | null>(null);
  const [formData, setFormData] = useState<CreateCustomerRequest>({
    name: '',
    email: '',
    account_status: CustomerStatus.ACTIVE,
    notification_preferences: {
      email_enabled: true,
      in_app_enabled: true,
      uat_notifications: true,
      production_notifications: true,
    },
  });
  const [validationErrors, setValidationErrors] = useState<Record<string, string>>({});

  useEffect(() => {
    if (isEdit && id) {
      loadCustomer();
    }
  }, [isEdit, id]);

  const loadCustomer = async () => {
    try {
      setLoadingData(true);
      const customer = await customersApi.getById(id!);
      setFormData({
        name: customer.name,
        email: customer.email,
        organization_name: customer.organization_name,
        phone: customer.phone,
        address: customer.address,
        account_status: customer.account_status,
        notification_preferences: customer.notification_preferences,
      });
    } catch (err: any) {
      setError(err.response?.data?.error?.message || 'Failed to load customer');
    } finally {
      setLoadingData(false);
    }
  };

  const validate = (): boolean => {
    const errors: Record<string, string> = {};

    if (!formData.name.trim()) {
      errors.name = 'Name is required';
    }

    if (!formData.email.trim()) {
      errors.email = 'Email is required';
    } else if (!/^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(formData.email)) {
      errors.email = 'Invalid email format';
    }

    if (!formData.account_status) {
      errors.account_status = 'Account status is required';
    }

    setValidationErrors(errors);
    return Object.keys(errors).length === 0;
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();

    if (!validate()) {
      return;
    }

    try {
      setLoading(true);
      setError(null);

      if (isEdit && id) {
        const updateData: UpdateCustomerRequest = {
          name: formData.name,
          email: formData.email,
          organization_name: formData.organization_name,
          phone: formData.phone,
          address: formData.address,
          account_status: formData.account_status,
          notification_preferences: formData.notification_preferences,
        };
        await customersApi.update(id, updateData);
      } else {
        await customersApi.create(formData);
      }

      navigate('/customers');
    } catch (err: any) {
      console.error('Customer form error:', err);
      let errorMessage = isEdit ? 'Failed to update customer' : 'Failed to create customer';

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

      if (errorMessage.toLowerCase().includes('already exists') || 
          errorMessage.toLowerCase().includes('duplicate')) {
        setValidationErrors({
          email: 'A customer with this email already exists',
        });
      }
    } finally {
      setLoading(false);
    }
  };

  const handleChange = (field: keyof CreateCustomerRequest, value: any) => {
    setFormData((prev) => ({ ...prev, [field]: value }));
    if (validationErrors[field]) {
      setValidationErrors((prev) => {
        const newErrors = { ...prev };
        delete newErrors[field];
        return newErrors;
      });
    }
  };

  const handleNotificationPreferenceChange = (
    field: keyof NotificationPreferences,
    value: boolean
  ) => {
    setFormData((prev) => ({
      ...prev,
      notification_preferences: {
        ...prev.notification_preferences,
        [field]: value,
      },
    }));
  };

  if (loadingData) {
    return (
      <div className="flex items-center justify-center h-64">
        <Spinner size="lg" />
      </div>
    );
  }

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <h1 className="text-3xl font-bold text-gray-900">
          {isEdit ? 'Edit Customer' : 'Create Customer'}
        </h1>
        <Button variant="ghost" onClick={() => navigate('/customers')}>
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

            <Input
              label="Email *"
              type="email"
              value={formData.email}
              onChange={(e) => handleChange('email', e.target.value)}
              error={validationErrors.email}
              required
            />

            <Input
              label="Organization Name"
              value={formData.organization_name || ''}
              onChange={(e) => handleChange('organization_name', e.target.value)}
            />

            <Input
              label="Phone"
              value={formData.phone || ''}
              onChange={(e) => handleChange('phone', e.target.value)}
            />

            <div className="md:col-span-2">
              <Input
                label="Address"
                value={formData.address || ''}
                onChange={(e) => handleChange('address', e.target.value)}
              />
            </div>

            <Select
              label="Account Status *"
              value={formData.account_status}
              onChange={(e) => handleChange('account_status', e.target.value as CustomerStatus)}
              error={validationErrors.account_status}
              required
              options={[
                { value: CustomerStatus.ACTIVE, label: 'Active' },
                { value: CustomerStatus.INACTIVE, label: 'Inactive' },
                { value: CustomerStatus.SUSPENDED, label: 'Suspended' },
              ]}
            />
          </div>

          {/* Notification Preferences */}
          <div className="border-t pt-6">
            <h3 className="text-lg font-medium text-gray-900 mb-4">Notification Preferences</h3>
            <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
              <label className="flex items-center">
                <input
                  type="checkbox"
                  checked={formData.notification_preferences.email_enabled}
                  onChange={(e) =>
                    handleNotificationPreferenceChange('email_enabled', e.target.checked)
                  }
                  className="mr-2"
                />
                <span className="text-sm text-gray-700">Email Notifications</span>
              </label>
              <label className="flex items-center">
                <input
                  type="checkbox"
                  checked={formData.notification_preferences.in_app_enabled}
                  onChange={(e) =>
                    handleNotificationPreferenceChange('in_app_enabled', e.target.checked)
                  }
                  className="mr-2"
                />
                <span className="text-sm text-gray-700">In-App Notifications</span>
              </label>
              <label className="flex items-center">
                <input
                  type="checkbox"
                  checked={formData.notification_preferences.uat_notifications}
                  onChange={(e) =>
                    handleNotificationPreferenceChange('uat_notifications', e.target.checked)
                  }
                  className="mr-2"
                />
                <span className="text-sm text-gray-700">UAT Notifications</span>
              </label>
              <label className="flex items-center">
                <input
                  type="checkbox"
                  checked={formData.notification_preferences.production_notifications}
                  onChange={(e) =>
                    handleNotificationPreferenceChange('production_notifications', e.target.checked)
                  }
                  className="mr-2"
                />
                <span className="text-sm text-gray-700">Production Notifications</span>
              </label>
            </div>
          </div>

          <div className="flex justify-end space-x-4 pt-4 border-t">
            <Button
              type="button"
              variant="ghost"
              onClick={() => navigate('/customers')}
              disabled={loading}
            >
              Cancel
            </Button>
            <Button type="submit" variant="primary" disabled={loading} isLoading={loading}>
              {isEdit ? 'Update Customer' : 'Create Customer'}
            </Button>
          </div>
        </form>
      </Card>
    </div>
  );
};

