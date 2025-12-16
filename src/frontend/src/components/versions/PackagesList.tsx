import React, { useState } from 'react';
import { PackageInfo, PackageType } from '@/types';
import { Button, Card, Select, Badge } from '@/components/ui';

export interface PackagesListProps {
  packages: PackageInfo[];
  onDownload: (packageId: string, fileName: string) => Promise<void>;
  onUpload?: () => void;
  canUpload?: boolean;
}

export const PackagesList: React.FC<PackagesListProps> = ({
  packages,
  onDownload,
  onUpload,
  canUpload = false,
}) => {
  const [osFilter, setOsFilter] = useState<string>('');
  const [archFilter, setArchFilter] = useState<string>('');

  const formatFileSize = (bytes: number): string => {
    if (bytes === 0) return '0 Bytes';
    const k = 1024;
    const sizes = ['Bytes', 'KB', 'MB', 'GB', 'TB'];
    const i = Math.floor(Math.log(bytes) / Math.log(k));
    return Math.round(bytes / Math.pow(k, i) * 100) / 100 + ' ' + sizes[i];
  };

  const getPackageTypeLabel = (type: PackageType): string => {
    return type.replace('_', ' ').replace(/\b\w/g, (l) => l.toUpperCase());
  };

  const getPackageTypeColor = (type: PackageType): string => {
    switch (type) {
      case PackageType.FULL_INSTALLER:
        return 'blue';
      case PackageType.UPDATE:
        return 'green';
      case PackageType.DELTA:
        return 'purple';
      case PackageType.ROLLBACK:
        return 'red';
      default:
        return 'gray';
    }
  };

  const filteredPackages = packages.filter((pkg) => {
    if (osFilter && pkg.os !== osFilter) return false;
    if (archFilter && pkg.architecture !== archFilter) return false;
    return true;
  });

  const uniqueOS = Array.from(new Set(packages.map((p) => p.os).filter(Boolean)));
  const uniqueArch = Array.from(new Set(packages.map((p) => p.architecture).filter(Boolean)));

  if (packages.length === 0) {
    return (
      <div className="text-center py-12">
        <p className="text-gray-500 mb-4">No packages available for this version.</p>
        {canUpload && onUpload && (
          <Button variant="primary" onClick={onUpload}>
            Upload Package
          </Button>
        )}
      </div>
    );
  }

  return (
    <div className="space-y-4">
      {/* Header and Filters */}
      <div className="flex items-center justify-between">
        <div>
          <h3 className="text-lg font-semibold text-gray-900">
            Packages ({filteredPackages.length})
          </h3>
          <p className="text-sm text-gray-500 mt-1">
            Total: {packages.length} package{packages.length !== 1 ? 's' : ''}
          </p>
        </div>
        {canUpload && onUpload && (
          <Button variant="primary" onClick={onUpload}>
            Upload Package
          </Button>
        )}
      </div>

      {/* Filters */}
      {(uniqueOS.length > 0 || uniqueArch.length > 0) && (
        <Card>
          <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
            {uniqueOS.length > 0 && (
              <Select
                label="Filter by OS"
                value={osFilter}
                onChange={(e) => setOsFilter(e.target.value)}
                options={[
                  { value: '', label: 'All Operating Systems' },
                  ...uniqueOS.map((os) => ({ value: os!, label: os! })),
                ]}
              />
            )}
            {uniqueArch.length > 0 && (
              <Select
                label="Filter by Architecture"
                value={archFilter}
                onChange={(e) => setArchFilter(e.target.value)}
                options={[
                  { value: '', label: 'All Architectures' },
                  ...uniqueArch.map((arch) => ({ value: arch!, label: arch! })),
                ]}
              />
            )}
          </div>
        </Card>
      )}

      {/* Packages Table */}
      <Card>
        <div className="overflow-x-auto">
          <table className="min-w-full divide-y divide-gray-200">
            <thead className="bg-gray-50">
              <tr>
                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                  Package Type
                </th>
                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                  File Name
                </th>
                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                  Size
                </th>
                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                  OS
                </th>
                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                  Architecture
                </th>
                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                  Checksum
                </th>
                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                  Uploaded
                </th>
                <th className="px-6 py-3 text-right text-xs font-medium text-gray-500 uppercase tracking-wider">
                  Actions
                </th>
              </tr>
            </thead>
            <tbody className="bg-white divide-y divide-gray-200">
              {filteredPackages.map((pkg) => (
                <tr key={pkg.id} className="hover:bg-gray-50">
                  <td className="px-6 py-4 whitespace-nowrap">
                    <Badge color={getPackageTypeColor(pkg.package_type)}>
                      {getPackageTypeLabel(pkg.package_type)}
                    </Badge>
                  </td>
                  <td className="px-6 py-4 whitespace-nowrap">
                    <div className="text-sm font-medium text-gray-900">{pkg.file_name}</div>
                  </td>
                  <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-500">
                    {formatFileSize(pkg.file_size)}
                  </td>
                  <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-500">
                    {pkg.os || <span className="text-gray-400">—</span>}
                  </td>
                  <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-500">
                    {pkg.architecture || <span className="text-gray-400">—</span>}
                  </td>
                  <td className="px-6 py-4 whitespace-nowrap">
                    <div className="text-xs font-mono text-gray-500 max-w-xs truncate">
                      {pkg.checksum_sha256}
                    </div>
                  </td>
                  <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-500">
                    {pkg.uploaded_at
                      ? new Date(pkg.uploaded_at).toLocaleDateString()
                      : <span className="text-gray-400">—</span>}
                  </td>
                  <td className="px-6 py-4 whitespace-nowrap text-right text-sm font-medium">
                    <Button
                      variant="ghost"
                      size="sm"
                      onClick={() => onDownload(pkg.id, pkg.file_name)}
                    >
                      Download
                    </Button>
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>

        {filteredPackages.length === 0 && (
          <div className="text-center py-8 text-gray-500">
            No packages match the selected filters.
          </div>
        )}
      </Card>
    </div>
  );
};

