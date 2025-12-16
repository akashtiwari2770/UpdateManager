import React, { useState, useEffect } from 'react';
import { UpdateRollout, RolloutStatus as RolloutStatusType } from '@/types';
import { Card, Badge, Button, Spinner, Alert, Modal, Input, Select } from '@/components/ui';
import { updateRolloutsApi } from '@/services/api/update-rollouts';
import { useNavigate } from 'react-router-dom';

export interface RolloutStatusProps {
  rolloutId: string;
  autoRefresh?: boolean;
  refreshInterval?: number;
}

export const RolloutStatus: React.FC<RolloutStatusProps> = ({
  rolloutId,
  autoRefresh = true,
  refreshInterval = 5000,
}) => {
  const navigate = useNavigate();
  const [rollout, setRollout] = useState<UpdateRollout | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [showProgressModal, setShowProgressModal] = useState(false);
  const [showStatusModal, setShowStatusModal] = useState(false);
  const [progressValue, setProgressValue] = useState(0);
  const [statusValue, setStatusValue] = useState<RolloutStatusType>(RolloutStatusType.PENDING);
  const [errorMessage, setErrorMessage] = useState('');
  const [updating, setUpdating] = useState(false);

  useEffect(() => {
    loadRollout();
  }, [rolloutId]);

  useEffect(() => {
    if (autoRefresh && rollout && rollout.status === RolloutStatusType.IN_PROGRESS) {
      const interval = setInterval(() => {
        loadRollout();
      }, refreshInterval);
      return () => clearInterval(interval);
    }
  }, [autoRefresh, rollout, refreshInterval]);

  const loadRollout = async () => {
    try {
      setLoading(true);
      setError(null);
      const data = await updateRolloutsApi.getById(rolloutId);
      setRollout(data);
    } catch (err: any) {
      setError(err.response?.data?.error?.message || 'Failed to load rollout status');
    } finally {
      setLoading(false);
    }
  };

  const handleUpdateProgress = async () => {
    if (!rollout) return;

    try {
      setUpdating(true);
      await updateRolloutsApi.updateProgress(rollout.id, progressValue);
      await loadRollout();
      setShowProgressModal(false);
      setProgressValue(0);
    } catch (err: any) {
      setError(err.response?.data?.error?.message || 'Failed to update progress');
    } finally {
      setUpdating(false);
    }
  };

  const handleUpdateStatus = async () => {
    if (!rollout) return;

    try {
      setUpdating(true);
      await updateRolloutsApi.updateStatus(rollout.id, statusValue, errorMessage || undefined);
      await loadRollout();
      setShowStatusModal(false);
      setStatusValue(RolloutStatusType.PENDING);
      setErrorMessage('');
    } catch (err: any) {
      setError(err.response?.data?.error?.message || 'Failed to update status');
    } finally {
      setUpdating(false);
    }
  };

  const getStatusColor = (status: RolloutStatusType): string => {
    switch (status) {
      case RolloutStatusType.PENDING:
        return 'yellow';
      case RolloutStatusType.IN_PROGRESS:
        return 'blue';
      case RolloutStatusType.COMPLETED:
        return 'green';
      case RolloutStatusType.FAILED:
        return 'red';
      case RolloutStatusType.CANCELLED:
        return 'gray';
      default:
        return 'gray';
    }
  };

  const getStatusLabel = (status: RolloutStatusType): string => {
    return status.replace('_', ' ').replace(/\b\w/g, (l) => l.toUpperCase());
  };

  const formatDate = (dateStr?: string): string => {
    if (!dateStr) return 'N/A';
    return new Date(dateStr).toLocaleString();
  };

  if (loading && !rollout) {
    return (
      <div className="flex justify-center items-center h-64">
        <Spinner />
      </div>
    );
  }

  if (error && !rollout) {
    return (
      <Alert variant="error" title="Error" onClose={() => setError(null)}>
        {error}
      </Alert>
    );
  }

  if (!rollout) {
    return (
      <div className="text-center py-12">
        <p className="text-gray-500">Rollout not found.</p>
      </div>
    );
  }

  return (
    <div className="space-y-6">
      {error && (
        <Alert variant="error" title="Error" onClose={() => setError(null)}>
          {error}
        </Alert>
      )}

      {/* Rollout Overview */}
      <Card title="Rollout Information">
        <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
          <div>
            <h3 className="text-sm font-medium text-gray-500 mb-2">Status</h3>
            <Badge color={getStatusColor(rollout.status)} size="lg">
              {getStatusLabel(rollout.status)}
            </Badge>
          </div>

          <div>
            <h3 className="text-sm font-medium text-gray-500 mb-2">Progress</h3>
            <div className="space-y-2">
              <div className="flex items-center justify-between text-sm">
                <span className="text-gray-700">{rollout.progress}%</span>
                <span className="text-gray-500">Complete</span>
              </div>
              <div className="w-full bg-gray-200 rounded-full h-2.5">
                <div
                  className="bg-blue-600 h-2.5 rounded-full transition-all duration-300"
                  style={{ width: `${rollout.progress}%` }}
                />
              </div>
            </div>
          </div>

          <div>
            <h3 className="text-sm font-medium text-gray-500 mb-1">Product ID</h3>
            <p className="text-gray-900">{rollout.product_id}</p>
          </div>

          <div>
            <h3 className="text-sm font-medium text-gray-500 mb-1">Endpoint ID</h3>
            <p className="text-gray-900">{rollout.endpoint_id}</p>
          </div>

          <div>
            <h3 className="text-sm font-medium text-gray-500 mb-1">From Version</h3>
            <p className="text-gray-900">{rollout.from_version}</p>
          </div>

          <div>
            <h3 className="text-sm font-medium text-gray-500 mb-1">To Version</h3>
            <p className="text-gray-900 font-semibold text-blue-600">{rollout.to_version}</p>
          </div>

          <div>
            <h3 className="text-sm font-medium text-gray-500 mb-1">Initiated By</h3>
            <p className="text-gray-900">{rollout.initiated_by}</p>
          </div>

          <div>
            <h3 className="text-sm font-medium text-gray-500 mb-1">Initiated At</h3>
            <p className="text-gray-900">{formatDate(rollout.initiated_at)}</p>
          </div>

          {rollout.started_at && (
            <div>
              <h3 className="text-sm font-medium text-gray-500 mb-1">Started At</h3>
              <p className="text-gray-900">{formatDate(rollout.started_at)}</p>
            </div>
          )}

          {rollout.completed_at && (
            <div>
              <h3 className="text-sm font-medium text-gray-500 mb-1">Completed At</h3>
              <p className="text-gray-900">{formatDate(rollout.completed_at)}</p>
            </div>
          )}

          {rollout.failed_at && (
            <div>
              <h3 className="text-sm font-medium text-gray-500 mb-1">Failed At</h3>
              <p className="text-gray-900 text-red-600">{formatDate(rollout.failed_at)}</p>
            </div>
          )}
        </div>

        {rollout.error_message && (
          <div className="mt-4 p-4 bg-red-50 border border-red-200 rounded-lg">
            <h4 className="text-sm font-medium text-red-900 mb-1">Error Message</h4>
            <p className="text-sm text-red-700">{rollout.error_message}</p>
          </div>
        )}
      </Card>

      {/* Action Buttons */}
      <div className="flex gap-4">
        <Button
          variant="secondary"
          onClick={() => {
            setProgressValue(rollout.progress);
            setShowProgressModal(true);
          }}
        >
          Update Progress
        </Button>
        <Button
          variant="secondary"
          onClick={() => {
            setStatusValue(rollout.status);
            setShowStatusModal(true);
          }}
        >
          Update Status
        </Button>
        <Button variant="secondary" onClick={() => navigate('/updates/rollouts')}>
          Back to Rollouts
        </Button>
      </div>

      {/* Update Progress Modal */}
      <Modal
        isOpen={showProgressModal}
        onClose={() => setShowProgressModal(false)}
        title="Update Progress"
      >
        <div className="space-y-4">
          <Input
            label="Progress (0-100%)"
            type="number"
            min="0"
            max="100"
            value={progressValue.toString()}
            onChange={(e) => setProgressValue(parseInt(e.target.value) || 0)}
          />
          <div className="flex justify-end gap-2">
            <Button variant="secondary" onClick={() => setShowProgressModal(false)} disabled={updating}>
              Cancel
            </Button>
            <Button variant="primary" onClick={handleUpdateProgress} isLoading={updating}>
              Update
            </Button>
          </div>
        </div>
      </Modal>

      {/* Update Status Modal */}
      <Modal
        isOpen={showStatusModal}
        onClose={() => setShowStatusModal(false)}
        title="Update Status"
      >
        <div className="space-y-4">
          <Select
            label="Status"
            value={statusValue}
            onChange={(e) => setStatusValue(e.target.value as RolloutStatusType)}
            options={[
              { value: RolloutStatusType.PENDING, label: 'Pending' },
              { value: RolloutStatusType.IN_PROGRESS, label: 'In Progress' },
              { value: RolloutStatusType.COMPLETED, label: 'Completed' },
              { value: RolloutStatusType.FAILED, label: 'Failed' },
              { value: RolloutStatusType.CANCELLED, label: 'Cancelled' },
            ]}
          />
          {statusValue === RolloutStatusType.FAILED && (
            <Input
              label="Error Message (Optional)"
              value={errorMessage}
              onChange={(e) => setErrorMessage(e.target.value)}
              placeholder="Describe the error..."
            />
          )}
          <div className="flex justify-end gap-2">
            <Button variant="secondary" onClick={() => setShowStatusModal(false)} disabled={updating}>
              Cancel
            </Button>
            <Button variant="primary" onClick={handleUpdateStatus} isLoading={updating}>
              Update
            </Button>
          </div>
        </div>
      </Modal>
    </div>
  );
};

