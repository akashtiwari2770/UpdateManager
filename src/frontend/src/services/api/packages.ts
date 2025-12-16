import apiClient from './client';
import { PackageInfo, PackageType } from '@/types';

export interface UploadPackageRequest {
  file: File;
  package_type: PackageType;
  os?: string;
  architecture?: string;
}

export const packagesApi = {
  upload: async (
    versionId: string,
    data: UploadPackageRequest,
    onProgress?: (progress: number) => void
  ): Promise<PackageInfo> => {
    const formData = new FormData();
    formData.append('file', data.file);
    formData.append('package_type', data.package_type);
    if (data.os) {
      formData.append('os', data.os);
    }
    if (data.architecture) {
      formData.append('architecture', data.architecture);
    }

    // Get axios instance to set headers properly for FormData
    const axiosInstance = apiClient.getInstance();
    
    const response = await axiosInstance.post<PackageInfo>(
      `/versions/${versionId}/packages`,
      formData,
      {
        headers: {
          'Content-Type': 'multipart/form-data',
        },
        onUploadProgress: (progressEvent) => {
          if (progressEvent.total && onProgress) {
            const percentCompleted = Math.round(
              (progressEvent.loaded * 100) / progressEvent.total
            );
            onProgress(percentCompleted);
          }
        },
      }
    );

    // Handle response - backend might return {success: true, data: {...}}
    if (response.data && (response.data as any).package) {
      return (response.data as any).package;
    }
    if (response.data && (response.data as any).data) {
      return (response.data as any).data;
    }
    return response.data;
  },

  list: async (versionId: string): Promise<PackageInfo[]> => {
    const response = await apiClient.get<any>(`/versions/${versionId}/packages`);
    
    // Handle various response formats
    if (response.data && response.data.packages) {
      return response.data.packages;
    }
    if (response.data && response.data.data) {
      return Array.isArray(response.data.data) ? response.data.data : [];
    }
    if (Array.isArray(response.data)) {
      return response.data;
    }
    return [];
  },

  download: async (versionId: string, packageId: string): Promise<Blob> => {
    const response = await apiClient.get(
      `/versions/${versionId}/packages/${packageId}/download`,
      {
        responseType: 'blob',
      }
    );
    return response.data;
  },

  delete: async (versionId: string, packageId: string): Promise<void> => {
    await apiClient.delete(`/versions/${versionId}/packages/${packageId}`);
  },
};

