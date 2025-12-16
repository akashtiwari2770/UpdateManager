import React, { useState } from 'react';
import { AuditLog } from '@/types';
import { Button } from '@/components/ui';

interface AuditLogDetailsProps {
  auditLog: AuditLog;
}

export const AuditLogDetails: React.FC<AuditLogDetailsProps> = ({ auditLog }) => {
  const [isExpanded, setIsExpanded] = useState(false);
  const [copied, setCopied] = useState(false);

  const handleCopy = async () => {
    try {
      await navigator.clipboard.writeText(JSON.stringify(auditLog.details, null, 2));
      setCopied(true);
      setTimeout(() => setCopied(false), 2000);
    } catch (error) {
      console.error('Failed to copy:', error);
    }
  };

  const formatJson = (obj: Record<string, any>): string => {
    return JSON.stringify(obj, null, 2);
  };

  return (
    <div className="border-t border-gray-200 pt-4 mt-4">
      <div className="flex items-center justify-between mb-2">
        <button
          onClick={() => setIsExpanded(!isExpanded)}
          className="text-sm font-medium text-blue-600 hover:text-blue-700 flex items-center gap-2"
        >
          {isExpanded ? (
            <>
              <svg className="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M5 15l7-7 7 7" />
              </svg>
              Hide Details
            </>
          ) : (
            <>
              <svg className="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M19 9l-7 7-7-7" />
              </svg>
              Show Details
            </>
          )}
        </button>
        {isExpanded && (
          <Button variant="ghost" size="sm" onClick={handleCopy}>
            {copied ? (
              <>
                <svg className="w-4 h-4 mr-1" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M5 13l4 4L19 7" />
                </svg>
                Copied!
              </>
            ) : (
              <>
                <svg className="w-4 h-4 mr-1" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M8 16H6a2 2 0 01-2-2V6a2 2 0 012-2h8a2 2 0 012 2v2m-6 12h8a2 2 0 002-2v-8a2 2 0 00-2-2h-8a2 2 0 00-2 2v8a2 2 0 002 2z" />
                </svg>
                Copy JSON
              </>
            )}
          </Button>
        )}
      </div>

      {isExpanded && (
        <div className="bg-gray-50 rounded-lg p-4 mt-2">
          <div className="space-y-4">
            <div>
              <h4 className="text-sm font-semibold text-gray-700 mb-2">Full Details</h4>
              <pre className="bg-white p-3 rounded border border-gray-200 text-xs overflow-x-auto">
                {formatJson(auditLog.details)}
              </pre>
            </div>
            <div className="grid grid-cols-2 gap-4 text-sm">
              <div>
                <span className="font-medium text-gray-700">IP Address:</span>
                <span className="ml-2 text-gray-600">{auditLog.ip_address}</span>
              </div>
              <div>
                <span className="font-medium text-gray-700">User Agent:</span>
                <span className="ml-2 text-gray-600 truncate block">{auditLog.user_agent}</span>
              </div>
            </div>
          </div>
        </div>
      )}
    </div>
  );
};

