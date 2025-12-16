import React, { useState } from 'react';
import { Button, Modal, Input, Alert } from '@/components/ui';
import { upgradePathsApi } from '@/services/api/upgrade-paths';

export interface BlockUpgradePathDialogProps {
  isOpen: boolean;
  productId: string;
  fromVersion: string;
  toVersion: string;
  onClose: () => void;
  onSuccess: () => void;
}

export const BlockUpgradePathDialog: React.FC<BlockUpgradePathDialogProps> = ({
  isOpen,
  productId,
  fromVersion,
  toVersion,
  onClose,
  onSuccess,
}) => {
  const [reason, setReason] = useState('');
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [validationError, setValidationError] = useState<string | null>(null);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();

    if (!reason.trim()) {
      setValidationError('Reason is required');
      return;
    }

    try {
      setLoading(true);
      setError(null);
      setValidationError(null);

      await upgradePathsApi.block(productId, fromVersion, toVersion, reason);
      onSuccess();
      setReason('');
      onClose();
    } catch (err: any) {
      setError(err.response?.data?.error?.message || 'Failed to block upgrade path');
    } finally {
      setLoading(false);
    }
  };

  const handleClose = () => {
    setReason('');
    setError(null);
    setValidationError(null);
    onClose();
  };

  return (
    <Modal
      isOpen={isOpen}
      onClose={handleClose}
      title="Block Upgrade Path"
    >
      <form onSubmit={handleSubmit} className="space-y-4">
        {error && (
          <Alert variant="error" title="Error" onClose={() => setError(null)}>
            {error}
          </Alert>
        )}

        <div>
          <p className="text-sm text-gray-600 mb-4">
            Are you sure you want to block the upgrade path from <strong>{fromVersion}</strong> to <strong>{toVersion}</strong>?
          </p>
          <p className="text-sm text-gray-600 mb-4">
            This action will prevent users from upgrading along this path.
          </p>
        </div>

        <div>
          <label className="block text-sm font-medium text-gray-700 mb-2">
            Reason <span className="text-red-500">*</span>
          </label>
          <textarea
            className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
            rows={4}
            value={reason}
            onChange={(e) => {
              setReason(e.target.value);
              setValidationError(null);
            }}
            placeholder="Enter the reason for blocking this upgrade path..."
            required
          />
          {validationError && (
            <p className="mt-1 text-sm text-red-600">{validationError}</p>
          )}
        </div>

        <div className="flex justify-end gap-2 pt-4 border-t">
          <Button variant="secondary" onClick={handleClose} disabled={loading}>
            Cancel
          </Button>
          <Button variant="primary" type="submit" isLoading={loading}>
            Block Path
          </Button>
        </div>
      </form>
    </Modal>
  );
};

