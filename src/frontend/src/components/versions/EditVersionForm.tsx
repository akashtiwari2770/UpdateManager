import React, { useState, useEffect } from 'react';
import { useParams, useNavigate } from 'react-router-dom';
import { versionsApi } from '@/services/api/versions';
import { Version, UpdateVersionRequest, ReleaseType, VersionState } from '@/types';
import { Button, Card, Input, Select, Alert, Spinner } from '@/components/ui';

export const EditVersionForm: React.FC = () => {
  const { id } = useParams<{ id: string }>();
  const navigate = useNavigate();
  const [version, setVersion] = useState<Version | null>(null);
  const [loading, setLoading] = useState(true);
  const [saving, setSaving] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [validationErrors, setValidationErrors] = useState<Record<string, string>>({});
  
  const [formData, setFormData] = useState<UpdateVersionRequest>({
    release_date: undefined,
    release_type: ReleaseType.FEATURE,
    eol_date: undefined,
    min_server_version: '',
    max_server_version: '',
    recommended_server_version: '',
  });

  useEffect(() => {
    if (id) {
      loadVersion();
    }
  }, [id]);

  const loadVersion = async () => {
    try {
      setLoading(true);
      setError(null);
      const data = await versionsApi.getById(id!);
      setVersion(data);

      // Check if version can be edited
      if (data.state !== VersionState.DRAFT) {
        setError('Only draft versions can be edited');
        return;
      }

      setFormData({
        release_date: data.release_date.split('T')[0],
        release_type: data.release_type,
        eol_date: data.eol_date ? data.eol_date.split('T')[0] : undefined,
        min_server_version: data.min_server_version || '',
        max_server_version: data.max_server_version || '',
        recommended_server_version: data.recommended_server_version || '',
      });
    } catch (err: any) {
      setError(err.response?.data?.error?.message || 'Failed to load version');
      console.error('Error loading version:', err);
    } finally {
      setLoading(false);
    }
  };

  const validate = (): boolean => {
    const errors: Record<string, string> = {};

    if (!formData.release_date) {
      errors.release_date = 'Release date is required';
    }

    if (!formData.release_type) {
      errors.release_type = 'Release type is required';
    }

    setValidationErrors(errors);
    return Object.keys(errors).length === 0;
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    
    if (!validate()) {
      return;
    }

    if (!version) {
      setError('Version not found');
      return;
    }

    if (version.state !== VersionState.DRAFT) {
      setError('Only draft versions can be edited');
      return;
    }

    try {
      setSaving(true);
      setError(null);
      
      // Convert date strings (YYYY-MM-DD) to ISO datetime strings
      // Date input returns YYYY-MM-DD, we need to convert to ISO datetime
      const convertDateToISO = (dateStr: string | undefined): string | undefined => {
        if (!dateStr) return undefined;
        // Create date at midnight UTC
        const date = new Date(dateStr + 'T00:00:00.000Z');
        return date.toISOString();
      };
      
      const payload: UpdateVersionRequest = {
        ...formData,
        release_date: convertDateToISO(formData.release_date),
        eol_date: convertDateToISO(formData.eol_date),
      };
      
      const updated = await versionsApi.update(version.id, payload);
      navigate(`/versions/${updated.id}`);
    } catch (err: any) {
      const errorMessage = err.response?.data?.error?.message || 'Failed to update version';
      setError(errorMessage);
    } finally {
      setSaving(false);
    }
  };

  const handleChange = (field: keyof UpdateVersionRequest, value: any) => {
    setFormData((prev) => ({ ...prev, [field]: value }));
    // Clear validation error for this field
    if (validationErrors[field]) {
      setValidationErrors((prev) => {
        const newErrors = { ...prev };
        delete newErrors[field];
        return newErrors;
      });
    }
  };

  if (loading) {
    return (
      <div className="flex justify-center items-center h-64">
        <Spinner />
      </div>
    );
  }

  if (error && !version) {
    return (
      <div className="space-y-4">
        <Button variant="ghost" onClick={() => navigate('/versions')}>
          ← Back to Versions
        </Button>
        <Alert variant="error" title="Error">
          {error}
        </Alert>
      </div>
    );
  }

  if (!version) {
    return (
      <div className="space-y-4">
        <Button variant="ghost" onClick={() => navigate('/versions')}>
          ← Back to Versions
        </Button>
        <Alert variant="info" title="Not Found">
          Version not found
        </Alert>
      </div>
    );
  }

  if (version.state !== VersionState.DRAFT) {
    return (
      <div className="space-y-4">
        <Button variant="ghost" onClick={() => navigate(`/versions/${version.id}`)}>
          ← Back to Version Details
        </Button>
        <Alert variant="warning" title="Cannot Edit">
          Only draft versions can be edited. This version is in {version.state} state.
        </Alert>
      </div>
    );
  }

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <h1 className="text-3xl font-bold text-gray-900">Edit Version</h1>
        <Button variant="ghost" onClick={() => navigate(`/versions/${version.id}`)}>
          ← Back to Version Details
        </Button>
      </div>

      {error && (
        <Alert variant="error" title="Error" onClose={() => setError(null)}>
          {error}
        </Alert>
      )}

      <Card>
        <div className="mb-4 p-4 bg-gray-50 rounded-md">
          <p className="text-sm text-gray-600">
            <strong>Version Number:</strong> {version.version_number} (cannot be changed)
          </p>
          <p className="text-sm text-gray-600">
            <strong>Product:</strong> {version.product_id} (cannot be changed)
          </p>
        </div>

        <form onSubmit={handleSubmit} className="space-y-6">
          {/* Release Type */}
          <div>
            <label className="block text-sm font-medium text-gray-700 mb-2">
              Release Type <span className="text-red-500">*</span>
            </label>
            <Select
              value={formData.release_type}
              onChange={(e) => handleChange('release_type', e.target.value as ReleaseType)}
            >
              <option value={ReleaseType.SECURITY}>Security</option>
              <option value={ReleaseType.FEATURE}>Feature</option>
              <option value={ReleaseType.MAINTENANCE}>Maintenance</option>
              <option value={ReleaseType.MAJOR}>Major</option>
            </Select>
            {validationErrors.release_type && (
              <p className="mt-1 text-sm text-red-600">{validationErrors.release_type}</p>
            )}
          </div>

          {/* Release Date */}
          <div>
            <label className="block text-sm font-medium text-gray-700 mb-2">
              Release Date <span className="text-red-500">*</span>
            </label>
            <Input
              type="date"
              value={formData.release_date || ''}
              onChange={(e) => handleChange('release_date', e.target.value || undefined)}
            />
            {validationErrors.release_date && (
              <p className="mt-1 text-sm text-red-600">{validationErrors.release_date}</p>
            )}
          </div>

          {/* EOL Date (Optional) */}
          <div>
            <label className="block text-sm font-medium text-gray-700 mb-2">
              End of Life Date (Optional)
            </label>
            <Input
              type="date"
              value={formData.eol_date ? formData.eol_date : ''}
              onChange={(e) => handleChange('eol_date', e.target.value || undefined)}
            />
          </div>

          {/* Server Version Requirements */}
          <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
            <div>
              <label className="block text-sm font-medium text-gray-700 mb-2">
                Min Server Version
              </label>
              <Input
                type="text"
                placeholder="e.g., 1.0.0"
                value={formData.min_server_version || ''}
                onChange={(e) => handleChange('min_server_version', e.target.value || undefined)}
              />
            </div>
            <div>
              <label className="block text-sm font-medium text-gray-700 mb-2">
                Max Server Version
              </label>
              <Input
                type="text"
                placeholder="e.g., 2.0.0"
                value={formData.max_server_version || ''}
                onChange={(e) => handleChange('max_server_version', e.target.value || undefined)}
              />
            </div>
            <div>
              <label className="block text-sm font-medium text-gray-700 mb-2">
                Recommended Server Version
              </label>
              <Input
                type="text"
                placeholder="e.g., 1.5.0"
                value={formData.recommended_server_version || ''}
                onChange={(e) => handleChange('recommended_server_version', e.target.value || undefined)}
              />
            </div>
          </div>

          {/* Actions */}
          <div className="flex items-center justify-end gap-4 pt-4 border-t">
            <Button
              type="button"
              variant="secondary"
              onClick={() => navigate(`/versions/${version.id}`)}
              disabled={saving}
            >
              Cancel
            </Button>
            <Button type="submit" disabled={saving}>
              {saving ? <Spinner size="sm" /> : 'Save Changes'}
            </Button>
          </div>
        </form>
      </Card>
    </div>
  );
};

