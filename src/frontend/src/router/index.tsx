import React from 'react';
import { BrowserRouter, Routes, Route } from 'react-router-dom';
import { MainLayout } from '@/components/layout';
import {
  Dashboard,
  Products,
  Versions,
  Updates,
  Notifications,
  AuditLogs,
  Compatibility,
  RolloutDetails,
  InitiateRollout,
  Customers,
  Licenses,
} from '@/pages';
import { CreateNotificationForm } from '@/components/notifications';
import {
  ProductDetails,
  CreateProductForm,
  EditProductForm,
} from '@/components/products';
import {
  VersionDetails,
  CreateVersionForm,
  EditVersionForm,
} from '@/components/versions';
import { CustomerForm } from '@/components/customers';
import { TenantForm } from '@/components/tenants';
import { DeploymentForm, DeploymentDetails } from '@/components/deployments';
import { SubscriptionForm, SubscriptionDetails } from '@/components/subscriptions';
import { LicenseForm, LicenseDetails } from '@/components/licenses';
import { AllocateLicenseForm } from '@/components/license-allocations';
import { CustomerDetails } from '@/pages/CustomerDetails';
import { TenantDetails } from '@/pages/TenantDetails';

export const AppRouter: React.FC = () => {
  return (
    <BrowserRouter>
      <MainLayout>
        <Routes>
          <Route path="/" element={<Dashboard />} />
          <Route path="/products" element={<Products />} />
          <Route path="/products/new" element={<CreateProductForm />} />
          <Route path="/products/:id" element={<ProductDetails />} />
          <Route path="/products/:id/edit" element={<EditProductForm />} />
          <Route path="/versions" element={<Versions />} />
          <Route path="/versions/new" element={<CreateVersionForm />} />
          <Route path="/versions/:id" element={<VersionDetails />} />
          <Route path="/versions/:id/edit" element={<EditVersionForm />} />
          <Route path="/products/:productId/versions/new" element={<CreateVersionForm />} />
          <Route path="/updates" element={<Updates />} />
          <Route path="/updates/rollout/new" element={<InitiateRollout />} />
          <Route path="/updates/rollouts/:id" element={<RolloutDetails />} />
          <Route path="/notifications" element={<Notifications />} />
          <Route path="/notifications/new" element={<CreateNotificationForm />} />
          <Route path="/audit-logs" element={<AuditLogs />} />
          <Route path="/compatibility" element={<Compatibility />} />
          <Route path="/customers" element={<Customers />} />
          <Route path="/customers/new" element={<CustomerForm />} />
          <Route path="/customers/:id" element={<CustomerDetails />} />
          <Route path="/customers/:id/edit" element={<CustomerForm />} />
          <Route path="/licenses" element={<Licenses />} />
          <Route path="/customers/:customerId/tenants/new" element={<TenantForm />} />
          <Route path="/customers/:customerId/tenants/:tenantId" element={<TenantDetails />} />
          <Route path="/customers/:customerId/tenants/:tenantId/edit" element={<TenantForm />} />
          <Route
            path="/customers/:customerId/tenants/:tenantId/deployments/new"
            element={<DeploymentForm />}
          />
          <Route
            path="/customers/:customerId/tenants/:tenantId/deployments/:deploymentId"
            element={<DeploymentDetails />}
          />
          <Route
            path="/customers/:customerId/tenants/:tenantId/deployments/:deploymentId/edit"
            element={<DeploymentForm />}
          />
          {/* Subscription Routes */}
          <Route
            path="/customers/:customerId/subscriptions/new"
            element={<SubscriptionForm />}
          />
          <Route
            path="/customers/:customerId/subscriptions/:subscriptionId"
            element={<SubscriptionDetails />}
          />
          <Route
            path="/customers/:customerId/subscriptions/:subscriptionId/edit"
            element={<SubscriptionForm />}
          />
          {/* License Routes */}
          <Route
            path="/customers/:customerId/subscriptions/:subscriptionId/licenses/new"
            element={<LicenseForm />}
          />
          <Route
            path="/customers/:customerId/subscriptions/:subscriptionId/licenses/:licenseId"
            element={<LicenseDetails />}
          />
          <Route
            path="/customers/:customerId/subscriptions/:subscriptionId/licenses/:licenseId/edit"
            element={<LicenseForm />}
          />
          <Route
            path="/customers/:customerId/subscriptions/:subscriptionId/licenses/:licenseId/allocate"
            element={<AllocateLicenseForm />}
          />
        </Routes>
      </MainLayout>
    </BrowserRouter>
  );
};

