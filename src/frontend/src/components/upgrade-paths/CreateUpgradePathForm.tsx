import React, { useState } from 'react';
import { Button, Card, Input, Select, Alert } from '@/components/ui';
import { upgradePathsApi } from '@/services/api/upgrade-paths';
import { UpgradePathType } from '@/types';

export interface CreateUpgradePathFormProps {
  productId: string;
  onSuccess: (path: any) => void;
  onCancel?: () => void;
}

export const CreateUpgradePathForm: React.FC<CreateUpgradePathFormProps> = ({
  productId,
  onSuccess,
  onCancel,
}) => {
  const [fromVersion, setFromVersion] = useState('');
  const [toVersion, setToVersion] = useState('');
  const [pathType, setPathType] = useState<UpgradePathType>(UpgradePathType.DIRECT);
  const [intermediateVersions, setIntermediateVersions] = useState<string[]>([]);
  const [currentIntermediate, setCurrentIntermediate] = useState('');
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [validationErrors, setValidationErrors] = useState<Record<string, string>>({});

  const addIntermediateVersion = () => {
    if (currentIntermediate.trim()) {
      setIntermediateVersions([...intermediateVersions, currentIntermediate.trim()]);
      setCurrentIntermediate('');
    }
  };

  const removeIntermediateVersion = (index: number) => {
    setIntermediateVersions(intermediateVersions.filter((_, i) => i !== index));
  };

  const validate = (): boolean => {
    const errors: Record<string, string> = {};

    if (!fromVersion.trim()) {
      errors.fromVersion = 'From version is required';
    }
    if (!toVersion.trim()) {
      errors.toVersion = 'To version is required';
    }
    if (fromVersion === toVersion) {
      errors.toVersion = 'To version must be different from from version';
    }
    if (pathType === UpgradePathType.MULTI_STEP && intermediateVersions.length === 0) {
      errors.intermediateVersions = 'At least one intermediate version is required for multi-step paths';
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

      const pathData: any = {
        from_version: fromVersion,
        to_version: toVersion,
        path_type: pathType,
      };

      if (pathType === UpgradePathType.MULTI_STEP && intermediateVersions.length > 0) {
        pathData.intermediate_versions = intermediateVersions;
      }

      const path = await upgradePathsApi.create(productId, pathData);
      onSuccess(path);
    } catch (err: any) {
      setError(err.response?.data?.error?.message || 'Failed to create upgrade path');
    } finally {
      setLoading(false);
    }
  };

  return (
    <Card title="Create Upgrade Path">
      <form onSubmit={handleSubmit} className="space-y-4">
        {error && (
          <Alert variant="error" title="Error" onClose={() => setError(null)}>
            {error}
          </Alert>
        )}

        <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
          <Input
            label="From Version"
            value={fromVersion}
            onChange={(e) => setFromVersion(e.target.value)}
            placeholder="e.g., 1.0.0"
            error={validationErrors.fromVersion}
            required
          />

          <Input
            label="To Version"
            value={toVersion}
            onChange={(e) => setToVersion(e.target.value)}
            placeholder="e.g., 2.0.0"
            error={validationErrors.toVersion}
            required
          />
        </div>

        <Select
          label="Path Type"
          value={pathType}
          onChange={(e) => setPathType(e.target.value as UpgradePathType)}
          options={[
            { value: UpgradePathType.DIRECT, label: 'Direct' },
            { value: UpgradePathType.MULTI_STEP, label: 'Multi-Step' },
            { value: UpgradePathType.BLOCKED, label: 'Blocked' },
          ]}
        />

        {pathType === UpgradePathType.MULTI_STEP && (
          <div>
            <label className="block text-sm font-medium text-gray-700 mb-2">
              Intermediate Versions
            </label>
            <div className="space-y-2">
              {intermediateVersions.map((version, index) => (
                <div key={index} className="flex items-center gap-2">
                  <Input
                    value={version}
                    disabled
                    className="flex-1 bg-gray-50"
                  />
                  <Button
                    variant="ghost"
                    size="sm"
                    onClick={() => removeIntermediateVersion(index)}
                  >
                    Remove
                  </Button>
                </div>
              ))}
              <div className="flex items-center gap-2">
                <Input
                  value={currentIntermediate}
                  onChange={(e) => setCurrentIntermediate(e.target.value)}
                  placeholder="e.g., 1.5.0"
                  onKeyPress={(e) => {
                    if (e.key === 'Enter') {
                      e.preventDefault();
                      addIntermediateVersion();
                    }
                  }}
                />
                <Button
                  type="button"
                  variant="secondary"
                  onClick={addIntermediateVersion}
                >
                  Add
                </Button>
              </div>
            </div>
            {validationErrors.intermediateVersions && (
              <p className="mt-1 text-sm text-red-600">
                {validationErrors.intermediateVersions}
              </p>
            )}
          </div>
        )}

        <div className="flex items-center justify-end gap-4 pt-4 border-t">
          {onCancel && (
            <Button variant="secondary" onClick={onCancel} disabled={loading}>
              Cancel
            </Button>
          )}
          <Button variant="primary" type="submit" isLoading={loading}>
            Create Upgrade Path
          </Button>
        </div>
      </form>
    </Card>
  );
};

