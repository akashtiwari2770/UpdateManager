import React, { useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { notificationsApi } from '@/services/api/notifications';
import { NotificationType, NotificationPriority } from '@/types';
import { Button, Input, Select, Alert, Card } from '@/components/ui';

export const CreateNotificationForm: React.FC = () => {
  const navigate = useNavigate();
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [formData, setFormData] = useState({
    recipient_id: '',
    type: NotificationType.UPDATE_AVAILABLE,
    title: '',
    message: '',
    priority: NotificationPriority.NORMAL,
    product_id: '',
    version_id: '',
  });

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setError(null);
    setLoading(true);

    try {
      // Note: The API endpoint structure may need adjustment based on backend
      await notificationsApi.create({
        recipient_id: formData.recipient_id,
        type: formData.type,
        title: formData.title,
        message: formData.message,
        priority: formData.priority,
        product_id: formData.product_id || undefined,
        version_id: formData.version_id || undefined,
      });

      navigate('/notifications');
    } catch (err: any) {
      setError(err.response?.data?.message || 'Failed to create notification');
    } finally {
      setLoading(false);
    }
  };

  const handleChange = (
    e: React.ChangeEvent<HTMLInputElement | HTMLSelectElement | HTMLTextAreaElement>
  ) => {
    const { name, value } = e.target;
    setFormData((prev) => ({ ...prev, [name]: value }));
  };

  return (
    <div className="max-w-2xl mx-auto">
      <h1 className="text-3xl font-bold text-gray-900 mb-6">Create Notification</h1>

      <Card>
        <form onSubmit={handleSubmit} className="space-y-6">
          {error && <Alert variant="error">{error}</Alert>}

          <div>
            <label htmlFor="recipient_id" className="block text-sm font-medium text-gray-700 mb-1">
              Recipient ID <span className="text-red-500">*</span>
            </label>
            <Input
              id="recipient_id"
              name="recipient_id"
              value={formData.recipient_id}
              onChange={handleChange}
              required
              placeholder="Enter recipient user ID"
            />
          </div>

          <div>
            <label htmlFor="type" className="block text-sm font-medium text-gray-700 mb-1">
              Type <span className="text-red-500">*</span>
            </label>
            <Select
              id="type"
              name="type"
              value={formData.type}
              onChange={handleChange}
              required
            >
              <option value={NotificationType.UPDATE_AVAILABLE}>Update Available</option>
              <option value={NotificationType.NEW_VERSION}>New Version</option>
              <option value={NotificationType.SECURITY_RELEASE}>Security Release</option>
              <option value={NotificationType.EOL_WARNING}>EOL Warning</option>
            </Select>
          </div>

          <div>
            <label htmlFor="title" className="block text-sm font-medium text-gray-700 mb-1">
              Title <span className="text-red-500">*</span>
            </label>
            <Input
              id="title"
              name="title"
              value={formData.title}
              onChange={handleChange}
              required
              placeholder="Enter notification title"
            />
          </div>

          <div>
            <label htmlFor="message" className="block text-sm font-medium text-gray-700 mb-1">
              Message <span className="text-red-500">*</span>
            </label>
            <textarea
              id="message"
              name="message"
              value={formData.message}
              onChange={handleChange}
              required
              rows={4}
              className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
              placeholder="Enter notification message"
            />
          </div>

          <div>
            <label htmlFor="priority" className="block text-sm font-medium text-gray-700 mb-1">
              Priority <span className="text-red-500">*</span>
            </label>
            <Select
              id="priority"
              name="priority"
              value={formData.priority}
              onChange={handleChange}
              required
            >
              <option value={NotificationPriority.LOW}>Low</option>
              <option value={NotificationPriority.NORMAL}>Normal</option>
              <option value={NotificationPriority.HIGH}>High</option>
              <option value={NotificationPriority.CRITICAL}>Critical</option>
            </Select>
          </div>

          <div>
            <label htmlFor="product_id" className="block text-sm font-medium text-gray-700 mb-1">
              Product ID (Optional)
            </label>
            <Input
              id="product_id"
              name="product_id"
              value={formData.product_id}
              onChange={handleChange}
              placeholder="Enter product ID"
            />
          </div>

          <div>
            <label htmlFor="version_id" className="block text-sm font-medium text-gray-700 mb-1">
              Version ID (Optional)
            </label>
            <Input
              id="version_id"
              name="version_id"
              value={formData.version_id}
              onChange={handleChange}
              placeholder="Enter version ID"
            />
          </div>

          <div className="flex gap-4 pt-4">
            <Button type="submit" isLoading={loading}>
              Create Notification
            </Button>
            <Button
              type="button"
              variant="secondary"
              onClick={() => navigate('/notifications')}
            >
              Cancel
            </Button>
          </div>
        </form>
      </Card>
    </div>
  );
};

