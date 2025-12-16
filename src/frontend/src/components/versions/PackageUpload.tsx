import React, { useState, useRef } from 'react';
import { PackageType } from '@/types';
import { Button, Card, Input, Select, Alert } from '@/components/ui';

export interface PackageUploadProps {
  versionId: string;
  onUpload: (file: File, metadata: PackageMetadata, onProgress?: (progress: number) => void) => Promise<void>;
  onCancel: () => void;
  loading?: boolean;
}

export interface PackageMetadata {
  package_type: PackageType;
  os?: string;
  architecture?: string;
}

const MAX_FILE_SIZE = 10 * 1024 * 1024 * 1024; // 10GB
const ALLOWED_FILE_TYPES = [
  '.zip',
  '.tar',
  '.tar.gz',
  '.tgz',
  '.deb',
  '.rpm',
  '.msi',
  '.exe',
  '.dmg',
  '.pkg',
];

export const PackageUpload: React.FC<PackageUploadProps> = ({
  versionId,
  onUpload,
  onCancel,
  loading = false,
}) => {
  const [file, setFile] = useState<File | null>(null);
  const [packageType, setPackageType] = useState<PackageType>(PackageType.FULL_INSTALLER);
  const [os, setOs] = useState<string>('');
  const [architecture, setArchitecture] = useState<string>('');
  const [dragActive, setDragActive] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [uploadProgress, setUploadProgress] = useState(0);
  const fileInputRef = useRef<HTMLInputElement>(null);

  const validateFile = (fileToValidate: File): string | null => {
    if (fileToValidate.size > MAX_FILE_SIZE) {
      return `File size exceeds maximum allowed size of ${MAX_FILE_SIZE / (1024 * 1024 * 1024)}GB`;
    }

    const fileExtension = '.' + fileToValidate.name.split('.').pop()?.toLowerCase();
    if (!ALLOWED_FILE_TYPES.includes(fileExtension)) {
      return `File type not allowed. Allowed types: ${ALLOWED_FILE_TYPES.join(', ')}`;
    }

    return null;
  };

  const handleDrag = (e: React.DragEvent) => {
    e.preventDefault();
    e.stopPropagation();
    if (e.type === 'dragenter' || e.type === 'dragover') {
      setDragActive(true);
    } else if (e.type === 'dragleave') {
      setDragActive(false);
    }
  };

  const handleDrop = (e: React.DragEvent) => {
    e.preventDefault();
    e.stopPropagation();
    setDragActive(false);

    if (e.dataTransfer.files && e.dataTransfer.files[0]) {
      const droppedFile = e.dataTransfer.files[0];
      const validationError = validateFile(droppedFile);
      if (validationError) {
        setError(validationError);
        return;
      }
      setFile(droppedFile);
      setError(null);
    }
  };

  const handleFileSelect = (e: React.ChangeEvent<HTMLInputElement>) => {
    if (e.target.files && e.target.files[0]) {
      const selectedFile = e.target.files[0];
      const validationError = validateFile(selectedFile);
      if (validationError) {
        setError(validationError);
        return;
      }
      setFile(selectedFile);
      setError(null);
    }
  };

  const handleUpload = async () => {
    if (!file) {
      setError('Please select a file to upload');
      return;
    }

    try {
      setError(null);
      setUploadProgress(0);
      
      const metadata: PackageMetadata = {
        package_type: packageType,
        os: os || undefined,
        architecture: architecture || undefined,
      };

      // Upload with progress tracking
      await onUpload(file, metadata, (progress) => {
        setUploadProgress(progress);
      });
      
      setUploadProgress(100);
      
      // Reset form after successful upload
      setTimeout(() => {
        setFile(null);
        setPackageType(PackageType.FULL_INSTALLER);
        setOs('');
        setArchitecture('');
        setUploadProgress(0);
        if (fileInputRef.current) {
          fileInputRef.current.value = '';
        }
      }, 1000);
    } catch (err: any) {
      setError(err.response?.data?.error?.message || 'Failed to upload package');
      setUploadProgress(0);
    }
  };

  const formatFileSize = (bytes: number): string => {
    if (bytes === 0) return '0 Bytes';
    const k = 1024;
    const sizes = ['Bytes', 'KB', 'MB', 'GB', 'TB'];
    const i = Math.floor(Math.log(bytes) / Math.log(k));
    return Math.round(bytes / Math.pow(k, i) * 100) / 100 + ' ' + sizes[i];
  };

  return (
    <div className="space-y-6">
      {error && (
        <Alert variant="error" title="Error" onClose={() => setError(null)}>
          {error}
        </Alert>
      )}

      {/* File Upload Area */}
      <Card>
        <div
          className={`
            border-2 border-dashed rounded-lg p-8 text-center
            transition-colors
            ${dragActive ? 'border-blue-500 bg-blue-50' : 'border-gray-300'}
            ${file ? 'border-green-500 bg-green-50' : ''}
          `}
          onDragEnter={handleDrag}
          onDragLeave={handleDrag}
          onDragOver={handleDrag}
          onDrop={handleDrop}
        >
          {file ? (
            <div className="space-y-2">
              <div className="text-green-600 font-medium">{file.name}</div>
              <div className="text-sm text-gray-600">{formatFileSize(file.size)}</div>
              <Button
                variant="ghost"
                size="sm"
                onClick={() => {
                  setFile(null);
                  if (fileInputRef.current) {
                    fileInputRef.current.value = '';
                  }
                }}
              >
                Remove File
              </Button>
            </div>
          ) : (
            <div className="space-y-4">
              <div className="text-gray-500">
                <svg
                  className="mx-auto h-12 w-12 text-gray-400"
                  stroke="currentColor"
                  fill="none"
                  viewBox="0 0 48 48"
                >
                  <path
                    d="M28 8H12a4 4 0 00-4 4v20m32-12v8m0 0v8a4 4 0 01-4 4H12a4 4 0 01-4-4v-4m32-4l-3.172-3.172a4 4 0 00-5.656 0L28 28M8 32l9.172-9.172a4 4 0 015.656 0L28 28m0 0l4 4m4-24h8m-4-4v8m-12 4h.02"
                    strokeWidth={2}
                    strokeLinecap="round"
                    strokeLinejoin="round"
                  />
                </svg>
              </div>
              <div>
                <p className="text-sm text-gray-600">
                  Drag and drop a file here, or{' '}
                  <button
                    type="button"
                    className="text-blue-600 hover:text-blue-500"
                    onClick={() => fileInputRef.current?.click()}
                  >
                    browse
                  </button>
                </p>
                <p className="text-xs text-gray-500 mt-1">
                  Maximum file size: {MAX_FILE_SIZE / (1024 * 1024 * 1024)}GB
                </p>
              </div>
              <input
                ref={fileInputRef}
                type="file"
                className="hidden"
                onChange={handleFileSelect}
                accept={ALLOWED_FILE_TYPES.join(',')}
              />
            </div>
          )}
        </div>

        {/* Upload Progress */}
        {uploadProgress > 0 && (
          <div className="mt-4">
            <div className="flex items-center justify-between mb-2">
              <span className="text-sm text-gray-600">Upload Progress</span>
              <span className="text-sm text-gray-600">{uploadProgress}%</span>
            </div>
            <div className="w-full bg-gray-200 rounded-full h-2">
              <div
                className="bg-blue-600 h-2 rounded-full transition-all duration-300"
                style={{ width: `${uploadProgress}%` }}
              />
            </div>
          </div>
        )}
      </Card>

      {/* Package Metadata */}
      <Card title="Package Information">
        <div className="space-y-4">
          <Select
            label="Package Type"
            value={packageType}
            onChange={(e) => setPackageType(e.target.value as PackageType)}
            options={[
              { value: PackageType.FULL_INSTALLER, label: 'Full Installer' },
              { value: PackageType.UPDATE, label: 'Update' },
              { value: PackageType.DELTA, label: 'Delta' },
              { value: PackageType.ROLLBACK, label: 'Rollback' },
            ]}
          />

          <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
            <Input
              label="Operating System (Optional)"
              value={os}
              onChange={(e) => setOs(e.target.value)}
              placeholder="e.g., Linux, Windows, macOS"
            />
            <Input
              label="Architecture (Optional)"
              value={architecture}
              onChange={(e) => setArchitecture(e.target.value)}
              placeholder="e.g., x86_64, arm64, amd64"
            />
          </div>
        </div>
      </Card>

      {/* Actions */}
      <div className="flex items-center justify-end gap-4 pt-4 border-t">
        <Button variant="secondary" onClick={onCancel} disabled={loading}>
          Cancel
        </Button>
        <Button
          variant="primary"
          onClick={handleUpload}
          isLoading={loading}
          disabled={!file || uploadProgress > 0}
        >
          Upload Package
        </Button>
      </div>
    </div>
  );
};

