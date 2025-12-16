# Update Manager Frontend

React frontend application for the Update Manager product release management system.

## Tech Stack

- **React 19** - UI library
- **TypeScript** - Type safety
- **Vite** - Build tool and dev server
- **React Router** - Routing
- **Zustand** - State management
- **Axios** - HTTP client
- **Tailwind CSS** - Styling
- **Playwright** - E2E testing

## Getting Started

### Prerequisites

- Node.js 18+ and npm
- Backend server running on `http://localhost:8080`

### Installation

```bash
# Install dependencies
npm install

# Install Playwright browsers (first time only)
npx playwright install --with-deps chromium
```

### Development

```bash
# Start dev server (runs on http://localhost:3000)
npm run dev

# Run type checking
npm run type-check

# Run linter
npm run lint

# Format code
npm run format
```

### Testing

```bash
# Run E2E tests
npm run test:e2e

# Run E2E tests with UI
npm run test:e2e:ui

# Run E2E tests in headed mode
npm run test:e2e:headed
```

### Building

```bash
# Build for production
npm run build

# Preview production build
npm run preview
```

## Project Structure

```
src/
├── components/          # React components
│   ├── ui/             # Base UI components (Button, Input, Card, etc.)
│   └── layout/         # Layout components (Header, Sidebar, MainLayout)
├── pages/              # Page components
├── services/           # API services
│   └── api/           # API client and endpoint services
├── store/              # State management (Zustand stores)
├── types/              # TypeScript type definitions
├── hooks/              # Custom React hooks
├── router/             # Routing configuration
├── App.tsx            # Main App component
└── main.tsx           # Entry point
```

## Environment Variables

Create a `.env` file in the frontend directory:

```env
VITE_API_BASE_URL=http://localhost:8080
VITE_API_VERSION=v1
VITE_APP_NAME=Update Manager
VITE_APP_VERSION=1.0.0
```

## Development Phases

This frontend is being built iteratively following the phases outlined in `docs/FRONTEND_DEVELOPMENT_PHASES.md`.

### Phase 1: Foundation & Project Setup ✅
- [x] React project setup with Vite
- [x] TypeScript configuration
- [x] Base UI components
- [x] Layout components (Header, Sidebar)
- [x] API integration layer
- [x] Routing setup
- [x] State management (Zustand)
- [x] Playwright test setup

### Next: Phase 2 - Product Management

## API Integration

The API client is configured in `src/services/api/client.ts` and automatically:
- Adds authentication tokens from localStorage
- Adds `X-User-ID` header for audit logging
- Handles common HTTP errors
- Provides typed responses

## Testing

Playwright tests are located in `tests/e2e/`. The test suite includes:
- Smoke tests
- Navigation tests
- Component interaction tests
- API integration tests

Run tests with:
```bash
npm run test:e2e
```

## License

Proprietary - Accops Technologies

