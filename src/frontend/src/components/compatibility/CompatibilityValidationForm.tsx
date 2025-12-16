import React, { useState } from 'react';
import { Button, Card, Input, Alert } from '@/components/ui';
import { compatibilityApi } from '@/services/api/compatibility';
import { CompatibilityMatrix } from '@/types';

export interface CompatibilityValidationFormProps {
  productId: string;
  versionNumber: string;
  onValidate: (result: CompatibilityMatrix) => void;
  onCancel?: () => void;
}

export const CompatibilityValidationForm: React.FC<CompatibilityValidationFormProps> = ({
  productId,
  versionNumber,
  onValidate,
  onCancel,
}) => {
  const [minServerVersion, setMinServerVersion] = useState('');
  const [maxServerVersion, setMaxServerVersion] = useState('');
  const [recommendedServerVersion, setRecommendedServerVersion] = useState('');
  const [incompatibleVersion, setIncompatibleVersion] = useState('');
  const [incompatibleVersions, setIncompatibleVersions] = useState<string[]>([]);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const handleAddIncompatibleVersion = () => {
    if (incompatibleVersion.trim() && !incompatibleVersions.includes(incompatibleVersion.trim())) {
      setIncompatibleVersions([...incompatibleVersions, incompatibleVersion.trim()]);
      setIncompatibleVersion('');
    }
  };

  const handleRemoveIncompatibleVersion = (version: string) => {
    setIncompatibleVersions(incompatibleVersions.filter(v => v !== version));
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();

    try {
      setLoading(true);
      setError(null);

      const result = await compatibilityApi.validate(
        productId,
        versionNumber,
        {
          min_server_version: minServerVersion.trim() || undefined,
          max_server_version: maxServerVersion.trim() || undefined,
          recommended_server_version: recommendedServerVersion.trim() || undefined,
          incompatible_versions: incompatibleVersions.length > 0 ? incompatibleVersions : undefined,
        }
      );

      onValidate(result);
    } catch (err: any) {
      setError(err.response?.data?.error?.message || 'Failed to validate compatibility');
    } finally {
      setLoading(false);
    }
  };

  return (
    <Card title="Validate Compatibility">
      <form onSubmit={handleSubmit} className="space-y-4">
        {error && (
          <Alert variant="error" title="Error" onClose={() => setError(null)}>
            {error}
          </Alert>
        )}

        <div>
          <label className="block text-sm font-medium text-gray-700 mb-1">
            Version <span className="text-red-500">*</span>
          </label>
          <Input
            value={versionNumber}
            disabled
            className="bg-gray-50"
          />
          <p className="mt-1 text-xs text-gray-500">Version being validated</p>
        </div>

        <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
          <Input
            label="Min Server Version (Optional)"
            value={minServerVersion}
            onChange={(e) => setMinServerVersion(e.target.value)}
            placeholder="e.g., 1.0.0"
          />

          <Input
            label="Max Server Version (Optional)"
            value={maxServerVersion}
            onChange={(e) => setMaxServerVersion(e.target.value)}
            placeholder="e.g., 2.0.0"
          />

          <Input
            label="Recommended Server Version (Optional)"
            value={recommendedServerVersion}
            onChange={(e) => setRecommendedServerVersion(e.target.value)}
            placeholder="e.g., 1.5.0"
          />
        </div>

        <div>
          <label className="block text-sm font-medium text-gray-700 mb-2">
            Incompatible Versions (Optional)
          </label>
          <div className="flex gap-2 mb-2">
            <Input
              value={incompatibleVersion}
              onChange={(e) => setIncompatibleVersion(e.target.value)}
              placeholder="e.g., 1.0.0"
              onKeyPress={(e) => {
                if (e.key === 'Enter') {
                  e.preventDefault();
                  handleAddIncompatibleVersion();
                }
              }}
            />
            <Button
              type="button"
              variant="secondary"
              onClick={handleAddIncompatibleVersion}
              disabled={!incompatibleVersion.trim()}
            >
              Add
            </Button>
          </div>
          {incompatibleVersions.length > 0 && (
            <div className="flex flex-wrap gap-2 mt-2">
              {incompatibleVersions.map((version) => (
                <span
                  key={version}
                  className="inline-flex items-center gap-1 px-3 py-1 rounded-full text-sm bg-red-100 text-red-800"
                >
                  {version}
                  <button
                    type="button"
                    onClick={() => handleRemoveIncompatibleVersion(version)}
                    className="hover:text-red-900"
                  >
                    ×
                  </button>
                </span>
              ))}
            </div>
          )}
        </div>

        <div className="flex items-center justify-end gap-4 pt-4 border-t">
          {onCancel && (
            <Button variant="secondary" onClick={onCancel} disabled={loading}>
              Cancel
            </Button>
          )}
          <Button variant="primary" type="submit" isLoading={loading}>
            Validate Compatibility
          </Button>
        </div>
      </form>
    </Card>
  );
};

