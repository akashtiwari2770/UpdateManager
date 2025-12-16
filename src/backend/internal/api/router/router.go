package router

import (
	"net/http"
	"strings"

	"updatemanager/internal/api/handlers"
	"updatemanager/internal/api/middleware"
	"updatemanager/internal/service"
)

// Router handles HTTP routing
type Router struct {
	mux *http.ServeMux
}

// NewRouter creates a new router
func NewRouter(services *service.ServiceFactory) *Router {
	mux := http.NewServeMux()

	// Initialize handlers
	productHandler := handlers.NewProductHandler(services.ProductService)
	versionHandler := handlers.NewVersionHandler(services.VersionService, services.PendingUpdatesService)
	compatibilityHandler := handlers.NewCompatibilityHandler(services.CompatibilityService)
	notificationHandler := handlers.NewNotificationHandler(services.NotificationService)
	upgradePathHandler := handlers.NewUpgradePathHandler(services.UpgradePathService)
	updateDetectionHandler := handlers.NewUpdateDetectionHandler(services.UpdateDetectionService)
	updateRolloutHandler := handlers.NewUpdateRolloutHandler(services.UpdateRolloutService)
	auditLogHandler := handlers.NewAuditLogHandler(services.AuditLogService)
	customerHandler := handlers.NewCustomerHandler(services.CustomerService)
	tenantHandler := handlers.NewTenantHandler(services.TenantService)
	deploymentHandler := handlers.NewDeploymentHandler(services.DeploymentService)
	pendingUpdatesHandler := handlers.NewPendingUpdatesHandler(services.PendingUpdatesService)
	subscriptionHandler := handlers.NewSubscriptionHandler(services.SubscriptionService)
	licenseHandler := handlers.NewLicenseHandler(services.LicenseService)
	licenseAllocationHandler := handlers.NewLicenseAllocationHandler(services.LicenseAllocationService)

	// API v1 routes
	apiV1 := "/api/v1"

	// Health check
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"status":"healthy"}`))
		}
	})

	// Product routes - must be registered before more specific routes
	// GET/POST /api/v1/products
	mux.HandleFunc(apiV1+"/products", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			productHandler.ListProducts(w, r)
		case http.MethodPost:
			productHandler.CreateProduct(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	// GET /api/v1/products/active
	mux.HandleFunc(apiV1+"/products/active", productHandler.GetActiveProducts)

	// GET /api/v1/products/by-product-id/:product_id
	mux.HandleFunc(apiV1+"/products/by-product-id/", productHandler.GetProductByProductID)

	// Complex product routes - handle nested paths
	mux.HandleFunc(apiV1+"/products/", func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path

		// Upgrade Path routes: /api/v1/products/:product_id/upgrade-paths
		if strings.Contains(path, "/upgrade-paths") {
			if strings.HasSuffix(path, "/block") {
				upgradePathHandler.BlockUpgradePath(w, r)
			} else if strings.Count(path, "/") >= 6 {
				upgradePathHandler.GetUpgradePath(w, r)
			} else if r.Method == http.MethodPost {
				upgradePathHandler.CreateUpgradePath(w, r)
			} else {
				http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			}
			return
		}

		// Compatibility routes: /api/v1/products/:product_id/versions/:version_number/compatibility
		if strings.Contains(path, "/versions/") && strings.Contains(path, "/compatibility") {
			if r.Method == http.MethodPost {
				compatibilityHandler.ValidateCompatibility(w, r)
			} else if r.Method == http.MethodGet {
				compatibilityHandler.GetCompatibility(w, r)
			} else {
				http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			}
			return
		}

		// Version routes under products: /api/v1/products/:product_id/versions
		if strings.Contains(path, "/versions") && !strings.Contains(path, "/compatibility") {
			if r.Method == http.MethodGet {
				versionHandler.GetVersionsByProduct(w, r)
			} else if r.Method == http.MethodPost {
				versionHandler.CreateVersion(w, r)
			} else {
				http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			}
			return
		}

		// Product CRUD: /api/v1/products/:id
		switch r.Method {
		case http.MethodGet:
			productHandler.GetProduct(w, r)
		case http.MethodPut:
			productHandler.UpdateProduct(w, r)
		case http.MethodDelete:
			productHandler.DeleteProduct(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	// Version routes
	// GET /api/v1/versions (list all versions)
	// GET/PUT /api/v1/versions/:id
	// POST /api/v1/versions/:id/submit
	// POST /api/v1/versions/:id/approve
	// POST /api/v1/versions/:id/release
	mux.HandleFunc(apiV1+"/versions", func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path
		basePath := apiV1 + "/versions"
		
		// Handle list all versions (exact match, no trailing slash, no ID)
		if path == basePath {
			if r.Method == http.MethodGet {
				versionHandler.ListVersions(w, r)
			} else {
				http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			}
			return
		}

		// Handle version actions with suffixes
		if strings.HasSuffix(path, "/submit") {
			versionHandler.SubmitForReview(w, r)
		} else if strings.HasSuffix(path, "/approve") {
			versionHandler.ApproveVersion(w, r)
		} else if strings.HasSuffix(path, "/release") {
			versionHandler.ReleaseVersion(w, r)
		} else if strings.HasSuffix(path, "/packages") {
			// Handle GET/POST /api/v1/versions/:id/packages
			if r.Method == http.MethodGet {
				versionHandler.ListPackages(w, r)
			} else if r.Method == http.MethodPost {
				versionHandler.UploadPackage(w, r)
			} else {
				http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			}
		} else if strings.Contains(path, "/packages/") && strings.HasSuffix(path, "/download") {
			// Handle GET /api/v1/versions/:id/packages/:package_id/download
			versionHandler.DownloadPackage(w, r)
		} else {
			// Handle GET/PUT /api/v1/versions/:id
			switch r.Method {
			case http.MethodGet:
				versionHandler.GetVersion(w, r)
			case http.MethodPut:
				versionHandler.UpdateVersion(w, r)
			default:
				http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			}
		}
	})
	
	// Also handle /api/v1/versions/ for compatibility (with trailing slash)
	mux.HandleFunc(apiV1+"/versions/", func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path
		basePath := apiV1 + "/versions/"
		
		// Handle list all versions (exact match with trailing slash, no ID)
		// Check if path is exactly /api/v1/versions/ or has no ID after the slash
		if path == basePath {
			if r.Method == http.MethodGet {
				versionHandler.ListVersions(w, r)
			} else {
				http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			}
			return
		}

		// Handle version actions with suffixes
		if strings.HasSuffix(path, "/submit") {
			versionHandler.SubmitForReview(w, r)
		} else if strings.HasSuffix(path, "/approve") {
			versionHandler.ApproveVersion(w, r)
		} else if strings.HasSuffix(path, "/release") {
			versionHandler.ReleaseVersion(w, r)
		} else if strings.HasSuffix(path, "/packages") {
			// Handle GET/POST /api/v1/versions/:id/packages
			if r.Method == http.MethodGet {
				versionHandler.ListPackages(w, r)
			} else if r.Method == http.MethodPost {
				versionHandler.UploadPackage(w, r)
			} else {
				http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			}
		} else if strings.Contains(path, "/packages/") && strings.HasSuffix(path, "/download") {
			// Handle GET /api/v1/versions/:id/packages/:package_id/download
			versionHandler.DownloadPackage(w, r)
		} else {
			// Handle GET/PUT /api/v1/versions/:id
			switch r.Method {
			case http.MethodGet:
				versionHandler.GetVersion(w, r)
			case http.MethodPut:
				versionHandler.UpdateVersion(w, r)
			default:
				http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			}
		}
	})

	// Compatibility routes
	// GET /api/v1/compatibility
	mux.HandleFunc(apiV1+"/compatibility", compatibilityHandler.ListCompatibility)

	// Notification routes
	// POST /api/v1/notifications
	// GET /api/v1/notifications?recipient_id=xxx
	mux.HandleFunc(apiV1+"/notifications", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			notificationHandler.GetNotifications(w, r)
		case http.MethodPost:
			notificationHandler.CreateNotification(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	// GET /api/v1/notifications/unread-count
	mux.HandleFunc(apiV1+"/notifications/unread-count", notificationHandler.GetUnreadCount)

	// POST /api/v1/notifications/mark-all-read
	mux.HandleFunc(apiV1+"/notifications/mark-all-read", notificationHandler.MarkAllAsRead)

	// Upgrade Path routes are handled in the /api/v1/products/ handler above
	// POST /api/v1/products/:product_id/upgrade-paths
	// GET /api/v1/products/:product_id/upgrade-paths/:from_version/:to_version
	// POST /api/v1/products/:product_id/upgrade-paths/:from_version/:to_version/block

	// Update Detection routes
	// GET/POST /api/v1/update-detections
	mux.HandleFunc(apiV1+"/update-detections", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			updateDetectionHandler.ListDetections(w, r)
		case http.MethodPost:
			updateDetectionHandler.DetectUpdate(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})
	// PUT /api/v1/update-detections/:id/available-version
	mux.HandleFunc(apiV1+"/update-detections/", updateDetectionHandler.UpdateAvailableVersion)

	// Update Rollout routes
	// GET/POST /api/v1/update-rollouts
	mux.HandleFunc(apiV1+"/update-rollouts", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			updateRolloutHandler.ListRollouts(w, r)
		case http.MethodPost:
			updateRolloutHandler.InitiateRollout(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})
	// GET/PUT /api/v1/update-rollouts/:id
	// PUT /api/v1/update-rollouts/:id/status
	// PUT /api/v1/update-rollouts/:id/progress
	mux.HandleFunc(apiV1+"/update-rollouts/", func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path
		if strings.HasSuffix(path, "/status") {
			updateRolloutHandler.UpdateRolloutStatus(w, r)
		} else if strings.HasSuffix(path, "/progress") {
			updateRolloutHandler.UpdateRolloutProgress(w, r)
		} else {
			// Handle GET/PUT /api/v1/update-rollouts/:id
			switch r.Method {
			case http.MethodGet:
				updateRolloutHandler.GetRollout(w, r)
			default:
				http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			}
		}
	})

	// Audit Log routes
	// GET /api/v1/audit-logs
	mux.HandleFunc(apiV1+"/audit-logs", auditLogHandler.GetAuditLogs)

	// Customer Management routes
	// GET/POST /api/v1/customers
	mux.HandleFunc(apiV1+"/customers", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			customerHandler.ListCustomers(w, r)
		case http.MethodPost:
			customerHandler.CreateCustomer(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	// Complex customer routes - handle nested paths
	mux.HandleFunc(apiV1+"/customers/", func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path

		// Tenant routes: /api/v1/customers/:customer_id/tenants
		if strings.Contains(path, "/tenants") {
			parts := strings.Split(strings.TrimPrefix(path, apiV1+"/customers/"), "/")
			if len(parts) >= 2 && parts[1] == "tenants" {
				// Check if there's a tenant ID (parts[2])
				if len(parts) >= 3 {
					// Deployment routes: /api/v1/customers/:customer_id/tenants/:tenant_id/deployments
					if len(parts) >= 4 && parts[3] == "deployments" {
						if len(parts) >= 5 {
							// Pending updates route: /api/v1/customers/:customer_id/tenants/:tenant_id/deployments/pending-updates
							if parts[4] == "pending-updates" {
								pendingUpdatesHandler.GetTenantPendingUpdates(w, r)
								return
							}

							// Updates route: /api/v1/customers/:customer_id/tenants/:tenant_id/deployments/:deployment_id/updates
							if len(parts) >= 6 && parts[5] == "updates" {
								pendingUpdatesHandler.GetDeploymentPendingUpdates(w, r)
								return
							}

							// GET/PUT/DELETE /api/v1/customers/:customer_id/tenants/:tenant_id/deployments/:deployment_id
							switch r.Method {
							case http.MethodGet:
								deploymentHandler.GetDeployment(w, r)
							case http.MethodPut:
								deploymentHandler.UpdateDeployment(w, r)
							case http.MethodDelete:
								deploymentHandler.DeleteDeployment(w, r)
							default:
								http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
							}
							return
						}

						// GET/POST /api/v1/customers/:customer_id/tenants/:tenant_id/deployments
						switch r.Method {
						case http.MethodGet:
							deploymentHandler.ListDeployments(w, r)
						case http.MethodPost:
							deploymentHandler.CreateDeployment(w, r)
						default:
							http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
						}
						return
					}

					// Statistics route: /api/v1/customers/:customer_id/tenants/:tenant_id/statistics
					if len(parts) >= 4 && parts[3] == "statistics" {
						tenantHandler.GetTenantStatistics(w, r)
						return
					}
					
					func (r *Router) Handler() http.Handler {
						handler := http.Handler(r.mux)
						handler = middleware.RecoveryMiddleware(handler)
						handler = middleware.LoggingMiddleware(handler)
						handler = middleware.CORSMiddleware(handler)
						return handler
					}


					// GET/PUT/DELETE /api/v1/customers/:customer_id/tenants/:tenant_id
					switch r.Method {
					case http.MethodGet:
						tenantHandler.GetTenant(w, r)
					case http.MethodPut:
						tenantHandler.UpdateTenant(w, r)
					case http.MethodDelete:
						tenantHandler.DeleteTenant(w, r)
					default:
						http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
					}
					return
				}

				// GET/POST /api/v1/customers/:customer_id/tenants
				switch r.Method {
				case http.MethodGet:
					customerHandler.GetCustomerTenants(w, r)
				case http.MethodPost:
					tenantHandler.CreateTenant(w, r)
				default:
					http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
				}
				return
			}
		}

		// Pending updates route: /api/v1/customers/:customer_id/deployments/pending-updates
		if strings.Contains(path, "/deployments/pending-updates") {
			pendingUpdatesHandler.GetCustomerPendingUpdates(w, r)
			return
		}

		// Subscription routes: /api/v1/customers/:customer_id/subscriptions
		if strings.Contains(path, "/subscriptions") {
			parts := strings.Split(strings.TrimPrefix(path, apiV1+"/customers/"), "/")
			if len(parts) >= 2 && parts[1] == "subscriptions" {
				// Check if there's a subscription ID (parts[2])
				if len(parts) >= 3 {
					// License routes: /api/v1/customers/:customer_id/subscriptions/:subscription_id/licenses
					if len(parts) >= 4 && parts[3] == "licenses" {
						if len(parts) >= 5 {
							// Allocation routes: /api/v1/customers/:customer_id/subscriptions/:subscription_id/licenses/:license_id/allocations
							if len(parts) >= 6 && parts[5] == "allocations" {
								if len(parts) >= 7 {
									// Release route: /api/v1/customers/:customer_id/subscriptions/:subscription_id/licenses/:license_id/allocations/:allocation_id/release
									if strings.HasSuffix(path, "/release") {
										licenseAllocationHandler.ReleaseAllocation(w, r)
										return
									}
									// GET /api/v1/customers/:customer_id/subscriptions/:subscription_id/licenses/:license_id/allocations/:allocation_id
									// (Currently no specific handler for single allocation, handled by GetAllocations)
								}
								// GET /api/v1/customers/:customer_id/subscriptions/:subscription_id/licenses/:license_id/allocations
								licenseAllocationHandler.GetAllocations(w, r)
								return
							}
							// Allocate route: /api/v1/customers/:customer_id/subscriptions/:subscription_id/licenses/:license_id/allocate
							if strings.HasSuffix(path, "/allocate") {
								licenseAllocationHandler.AllocateLicense(w, r)
								return
							}
							// Utilization route: /api/v1/customers/:customer_id/subscriptions/:subscription_id/licenses/:license_id/utilization
							if strings.HasSuffix(path, "/utilization") {
								licenseAllocationHandler.GetLicenseUtilization(w, r)
								return
							}
							// Statistics route: /api/v1/customers/:customer_id/subscriptions/:subscription_id/licenses/:license_id/statistics
							if strings.HasSuffix(path, "/statistics") {
								licenseHandler.GetLicenseStatistics(w, r)
								return
							}
							// Renew route: /api/v1/customers/:customer_id/subscriptions/:subscription_id/licenses/:license_id/renew
							if strings.HasSuffix(path, "/renew") {
								licenseHandler.RenewLicense(w, r)
								return
							}
							// GET/PUT/DELETE /api/v1/customers/:customer_id/subscriptions/:subscription_id/licenses/:license_id
							switch r.Method {
							case http.MethodGet:
								licenseHandler.GetLicense(w, r)
							case http.MethodPut:
								licenseHandler.UpdateLicense(w, r)
							case http.MethodDelete:
								licenseHandler.RevokeLicense(w, r)
							default:
								http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
							}
							return
						}
						// GET/POST /api/v1/customers/:customer_id/subscriptions/:subscription_id/licenses
						switch r.Method {
						case http.MethodGet:
							licenseHandler.ListLicenses(w, r)
						case http.MethodPost:
							licenseHandler.AssignLicense(w, r)
						default:
							http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
						}
						return
					}
					// Statistics route: /api/v1/customers/:customer_id/subscriptions/:subscription_id/statistics
					if strings.HasSuffix(path, "/statistics") {
						subscriptionHandler.GetSubscriptionStatistics(w, r)
						return
					}
					// Renew route: /api/v1/customers/:customer_id/subscriptions/:subscription_id/renew
					if strings.HasSuffix(path, "/renew") {
						subscriptionHandler.RenewSubscription(w, r)
						return
					}
					// GET/PUT/DELETE /api/v1/customers/:customer_id/subscriptions/:subscription_id
					switch r.Method {
					case http.MethodGet:
						subscriptionHandler.GetSubscription(w, r)
					case http.MethodPut:
						subscriptionHandler.UpdateSubscription(w, r)
					case http.MethodDelete:
						subscriptionHandler.DeleteSubscription(w, r)
					default:
						http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
					}
					return
				}
				// GET/POST /api/v1/customers/:customer_id/subscriptions
				switch r.Method {
				case http.MethodGet:
					subscriptionHandler.ListSubscriptions(w, r)
				case http.MethodPost:
					subscriptionHandler.CreateSubscription(w, r)
				default:
					http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
				}
				return
			}
		}

		// License routes for tenants: /api/v1/customers/:customer_id/tenants/:tenant_id/licenses
		// Note: This is handled within the /tenants section above, but we check here for license-specific routes
		if strings.Contains(path, "/tenants") && strings.Contains(path, "/licenses") {
			parts := strings.Split(strings.TrimPrefix(path, apiV1+"/customers/"), "/")
			if len(parts) >= 5 && parts[1] == "tenants" && parts[4] == "licenses" {
				licenseAllocationHandler.GetAllocationsByTenant(w, r)
				return
			}
			// Check for deployments/licenses route
			if len(parts) >= 7 && parts[1] == "tenants" && parts[4] == "deployments" && parts[6] == "licenses" {
				licenseAllocationHandler.GetAllocationsByDeployment(w, r)
				return
			}
		}

		// Statistics route: /api/v1/customers/:id/statistics
		if strings.HasSuffix(path, "/statistics") {
			customerHandler.GetCustomerStatistics(w, r)
			return
		}

		// GET/PUT/DELETE /api/v1/customers/:id
		switch r.Method {
		case http.MethodGet:
			customerHandler.GetCustomer(w, r)
		case http.MethodPut:
			customerHandler.UpdateCustomer(w, r)
		case http.MethodDelete:
			customerHandler.DeleteCustomer(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	// Pending Updates routes (admin view)
	// GET /api/v1/updates/pending
	mux.HandleFunc(apiV1+"/updates/pending", pendingUpdatesHandler.GetAllPendingUpdates)

	return &Router{mux: mux}
}

// Handler returns the HTTP handler with middleware
func (r *Router) Handler() http.Handler {
	handler := http.Handler(r.mux)
	handler = middleware.RecoveryMiddleware(handler)
	handler = middleware.LoggingMiddleware(handler)
	handler = middleware.CORSMiddleware(handler)
	return handler
}
