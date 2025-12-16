import React from 'react';
import { CompatibilityMatrixList } from '@/components/compatibility';

export const Compatibility: React.FC = () => {
  return (
    <div className="space-y-6">
      <div>
        <h1 className="text-2xl font-bold text-gray-900">Compatibility Matrices</h1>
        <p className="mt-1 text-sm text-gray-500">
          View and manage compatibility matrices for all product versions.
        </p>
      </div>
      <CompatibilityMatrixList />
    </div>
  );
};

