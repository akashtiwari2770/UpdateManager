import React from 'react';
import { useNavigate } from 'react-router-dom';
import { Version, VersionState } from '@/types';
import { Card, Badge, Spinner } from '@/components/ui';

interface RecentUpdatesProps {
  updates: Version[];
  loading?: boolean;
}

export const RecentUpdates: React.FC<RecentUpdatesProps> = ({ updates, loading }) => {
  const navigate = useNavigate();

  const getStateBadgeColor = (state: VersionState) => {
    switch (state) {
      case VersionState.RELEASED:
        return 'bg-green-100 text-green-800';
      case VersionState.APPROVED:
        return 'bg-blue-100 text-blue-800';
      case VersionState.PENDING_REVIEW:
        return 'bg-yellow-100 text-yellow-800';
      case VersionState.DRAFT:
        return 'bg-gray-100 text-gray-800';
      default:
        return 'bg-gray-100 text-gray-800';
    }
  };

  const formatDate = (dateString: string) => {
    return new Date(dateString).toLocaleDateString('en-US', {
      year: 'numeric',
      month: 'short',
      day: 'numeric',
    });
  };

  return (
    <Card title="Recent Updates">
      {loading ? (
        <div className="flex items-center justify-center py-8">
          <Spinner />
        </div>
      ) : updates.length === 0 ? (
        <div className="text-center py-8 text-gray-500">
          <p>No recent updates</p>
        </div>
      ) : (
        <>
          <div className="overflow-x-auto">
            <table className="min-w-full divide-y divide-gray-200">
              <thead className="bg-gray-50">
                <tr>
                  <th className="px-4 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                    Product
                  </th>
                  <th className="px-4 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                    Version
                  </th>
                  <th className="px-4 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                    Release Date
                  </th>
                  <th className="px-4 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                    Status
                  </th>
                </tr>
              </thead>
              <tbody className="bg-white divide-y divide-gray-200">
                {updates.slice(0, 10).map((version) => (
                  <tr
                    key={version.id}
                    onClick={() => navigate(`/versions/${version.id}`)}
                    className="hover:bg-gray-50 cursor-pointer"
                  >
                    <td className="px-4 py-3 whitespace-nowrap text-sm font-medium text-gray-900">
                      {version.product_id}
                    </td>
                    <td className="px-4 py-3 whitespace-nowrap text-sm text-gray-600">
                      {version.version_number}
                    </td>
                    <td className="px-4 py-3 whitespace-nowrap text-sm text-gray-500">
                      {formatDate(version.release_date)}
                    </td>
                    <td className="px-4 py-3 whitespace-nowrap">
                      <Badge className={getStateBadgeColor(version.state)}>
                        {version.state.replace('_', ' ')}
                      </Badge>
                    </td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
          {updates.length >= 10 && (
            <div className="mt-4 text-center">
              <button
                onClick={() => navigate('/versions')}
                className="text-sm text-blue-600 hover:text-blue-700 font-medium"
              >
                View All Versions →
              </button>
            </div>
          )}
        </>
      )}
    </Card>
  );
};

