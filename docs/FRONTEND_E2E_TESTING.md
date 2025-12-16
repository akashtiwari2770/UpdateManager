# Frontend E2E Testing Guide

This guide explains how to run and work with Playwright E2E tests for the Update Manager frontend.

## Prerequisites

1. **Node.js 18+** (v24.11.1 recommended)
   - If using nvm: `source ~/.nvm/nvm.sh && nvm use 24.11.1`
2. **Frontend dependencies installed**
   - Run `make frontend-install` or `cd src/frontend && npm install`
3. **Playwright browsers installed**
   - Run `cd src/frontend && npx playwright install --with-deps chromium` (for Chromium only)
   - Run `cd src/frontend && npx playwright install firefox webkit` (for all browsers)
   - Or install system dependencies first: `npx playwright install-deps` (installs required system libraries)
   - Then install browsers: `npx playwright install` (installs all browsers)

## Quick Start

### Using Makefile (Recommended)

```bash
# Run all E2E tests (headless)
make frontend-test-e2e

# Run E2E tests with interactive UI
make frontend-test-e2e-ui

# Run E2E tests in headed mode (see browser)
make frontend-test-e2e-headed
```

### Using npm directly

```bash
cd src/frontend

# Run all E2E tests (headless)
npm run test:e2e

# Run E2E tests with interactive UI
npm run test:e2e:ui

# Run E2E tests in headed mode (see browser)
npm run test:e2e:headed
```

## Test Modes Explained

### 1. Headless Mode (Default)
```bash
make frontend-test-e2e
# or
npm run test:e2e
```
- Tests run in background without opening browser windows
- Fastest execution
- Best for CI/CD and quick validation
- Results shown in terminal

### 2. UI Mode (Interactive)
```bash
make frontend-test-e2e-ui
# or
npm run test:e2e:ui
```
- Opens Playwright's interactive test runner UI
- See tests running in real-time
- Debug tests step-by-step
- Best for development and debugging

### 3. Headed Mode (Visible Browser)
```bash
make frontend-test-e2e-headed
# or
npm run test:e2e:headed
```
- Opens actual browser windows
- See what the tests are doing
- Useful for visual debugging
- Slower than headless mode

## Test Structure

Tests are located in `src/frontend/tests/e2e/`:

```
tests/
├── e2e/
│   ├── smoke.spec.ts          # Basic smoke tests
│   └── navigation.spec.ts     # Navigation tests
└── helpers/
    ├── page-objects.ts        # Page object models
    └── test-data.ts           # Test data fixtures
```

## Running Specific Tests

### Run a specific test file
```bash
cd src/frontend
npx playwright test tests/e2e/smoke.spec.ts
```

### Run tests matching a pattern
```bash
cd src/frontend
npx playwright test --grep "navigation"
```

### Run a specific test by name
```bash
cd src/frontend
npx playwright test --grep "app loads successfully"
```

## Test Configuration

The Playwright configuration is in `src/frontend/playwright.config.ts`:

- **Test Directory**: `./tests/e2e`
- **Base URL**: `http://localhost:3000`
- **Browsers**: Chromium, Firefox, WebKit (Safari)
- **Execution Mode**: Sequential (one test at a time) - reduces flakiness
- **Retries**: 1 retry locally, 2 retries in CI
- **Timeouts**: 30s global, 15s actions, 30s navigation
- **Auto-start server**: Yes (starts dev server automatically)
- **Screenshots**: On failure
- **Videos**: On failure
- **Traces**: On retry

## Viewing Test Results

### HTML Report (Default)
After running tests, an HTML report is generated:

```bash
cd src/frontend
npx playwright show-report
```

This opens an interactive HTML report showing:
- Test results
- Screenshots on failure
- Videos on failure
- Test traces

### Terminal Output
Test results are also displayed in the terminal with:
- ✅ Passed tests
- ❌ Failed tests
- ⏱️ Execution time
- 📊 Summary statistics

## Debugging Tests

### 1. Using UI Mode (Easiest)
```bash
make frontend-test-e2e-ui
```
- Click on a test to see it run step-by-step
- Pause at any point
- Inspect page state
- Step through actions

### 2. Using Headed Mode
```bash
make frontend-test-e2e-headed
```
- Watch the browser as tests execute
- See what's happening visually
- Identify UI issues

### 3. Using Playwright Inspector
```bash
cd src/frontend
PWDEBUG=1 npx playwright test
```
- Opens Playwright Inspector
- Step through test execution
- Inspect selectors and page state

### 4. Using VS Code Extension
Install the "Playwright Test for VSCode" extension:
- Run tests from VS Code
- Set breakpoints
- Debug directly in editor

## Writing New Tests

### Test File Structure
```typescript
import { test, expect } from '@playwright/test';

test.describe('Feature Name', () => {
  test.beforeEach(async ({ page }) => {
    await page.goto('/');
  });

  test('should do something', async ({ page }) => {
    // Test code here
    await expect(page.locator('h1')).toContainText('Expected Text');
  });
});
```

### Using Page Objects
```typescript
import { DashboardPage } from '../helpers/page-objects';

test('example', async ({ page }) => {
  const dashboard = new DashboardPage(page);
  await dashboard.goto('/');
  await dashboard.navigateToProducts();
});
```

## Common Test Scenarios

### Testing Navigation
```typescript
test('navigate to products page', async ({ page }) => {
  await page.goto('/');
  await page.click('text=Products');
  await expect(page).toHaveURL(/.*products/);
});
```

### Testing Forms
```typescript
test('fill and submit form', async ({ page }) => {
  await page.fill('input[name="name"]', 'Test Product');
  await page.selectOption('select[name="type"]', 'server');
  await page.click('button[type="submit"]');
  await expect(page.locator('.success')).toBeVisible();
});
```

### Testing API Integration
```typescript
test('loads data from API', async ({ page }) => {
  await page.goto('/products');
  await expect(page.locator('.product-item')).toHaveCount(5);
});
```

## Troubleshooting

### Tests fail because frontend isn't running
- Playwright should auto-start the dev server
- If not, manually start: `make frontend-dev`
- Check that port 3000 is available

### Tests timeout
- Increase timeout in `playwright.config.ts`
- Check if backend is running (some tests need it)
- Verify MongoDB is running

### Browser not found or missing dependencies
```bash
cd src/frontend

# Install system dependencies (required for Firefox/WebKit)
npx playwright install-deps

# Install browsers
npx playwright install chromium  # For Chromium only
npx playwright install firefox webkit  # For Firefox and WebKit
npx playwright install  # For all browsers
```

### Node.js version issues
```bash
source ~/.nvm/nvm.sh && nvm use 24.11.1
```

### Clear test artifacts
```bash
cd src/frontend
rm -rf test-results playwright-report
```

## CI/CD Integration

For CI/CD, tests run in headless mode:

```bash
cd src/frontend
npm run test:e2e
```

The configuration automatically:
- Runs in headless mode
- Retries failed tests (2 retries in CI)
- Generates HTML reports
- Captures screenshots/videos on failure

## Best Practices

1. **Use Page Objects** - Keep tests maintainable
2. **Use data-testid** - Prefer stable selectors
3. **Wait for elements** - Use `waitFor()` when needed
4. **Isolate tests** - Each test should be independent
5. **Clean up** - Reset state between tests
6. **Use fixtures** - Share common setup code

## Test Coverage (Phase 1)

Currently implemented tests:
- ✅ Smoke tests (app loads, basic navigation)
- ✅ Navigation tests (route changes, active highlighting)
- ✅ Sidebar tests (collapse/expand)

Upcoming in future phases:
- Component interaction tests
- Form validation tests
- API integration tests
- Accessibility tests
- Cross-browser tests

## Additional Resources

- [Playwright Documentation](https://playwright.dev/)
- [Playwright Best Practices](https://playwright.dev/docs/best-practices)
- [Test Generator](https://playwright.dev/docs/codegen) - Record tests automatically

