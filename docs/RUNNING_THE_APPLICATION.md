# Running Update Manager - Backend & Frontend

## Quick Start Guide

### Prerequisites

1. **MongoDB** - Must be running (Docker or local)
2. **Node.js 18+** - For frontend (v24.11.1 recommended)
   - If using nvm: `nvm use 24.11.1` (or `nvm use node` for latest)
   - The Makefile will automatically use nvm if available
3. **Go 1.21+** - For backend

### Step-by-Step Instructions

#### 1. Start MongoDB

**Option A: Using Docker (Recommended)**
```bash
make db-start
```

**Option B: Local MongoDB**
```bash
# Ensure MongoDB is running on localhost:27017
```

#### 2. Setup Database (First time only)
```bash
make db-setup
make db-indexes
```

#### 3. Start Backend Server

**Terminal 1:**
```bash
# From project root
make run
```

Or manually:
```bash
cd src/backend
go run ./cmd/server
```

Backend will run on: **http://localhost:8080**

#### 4. Start Frontend Development Server

**Terminal 2:**
```bash
# From project root - Makefile will handle Node.js version
make frontend-dev
```

Or manually:
```bash
# Switch to correct Node.js version (if using nvm)
source ~/.nvm/nvm.sh && nvm use 24.11.1

# Then start the dev server
cd src/frontend
npm run dev
```

Frontend will run on: **http://localhost:3000**

### Access the Application

- **Frontend UI**: http://localhost:3000
- **Backend API**: http://localhost:8080
- **API Health Check**: http://localhost:8080/health
- **API Documentation**: See `docs/api-specification.md`

### Running Frontend E2E Tests

See detailed guide: `docs/FRONTEND_E2E_TESTING.md`

**Quick commands:**
```bash
# Run all E2E tests (headless)
make frontend-test-e2e

# Run with interactive UI
make frontend-test-e2e-ui

# Run in headed mode (see browser)
make frontend-test-e2e-headed
```

## Running Both Services Together

### Option 1: Two Separate Terminals (Recommended for Development)

**Terminal 1 - Backend:**
```bash
make run
```

**Terminal 2 - Frontend:**
```bash
cd src/frontend && npm run dev
```

### Option 2: Background Processes

**Start Backend in Background:**
```bash
make run &
```

**Start Frontend in Background:**
```bash
cd src/frontend && npm run dev &
```

### Option 3: Using tmux or screen

```bash
# Create a new tmux session
tmux new-session -d -s updatemanager

# Split window
tmux split-window -h

# Run backend in left pane
tmux send-keys -t updatemanager:0.0 "make run" C-m

# Run frontend in right pane
tmux send-keys -t updatemanager:0.1 "cd src/frontend && npm run dev" C-m

# Attach to session
tmux attach-session -t updatemanager
```

## Verification

### Check Backend is Running
```bash
curl http://localhost:8080/health
# Should return: {"status":"healthy"}
```

### Check Frontend is Running
```bash
curl http://localhost:3000
# Should return HTML content
```

Or open in browser: http://localhost:3000

## Troubleshooting

### Backend won't start
- Check MongoDB is running: `make db-status`
- Check port 8080 is not in use: `lsof -i :8080`
- Check backend logs for errors

### Frontend won't start
- **Check Node.js version**: `node --version` (should be 18+)
  - If using nvm: `source ~/.nvm/nvm.sh && nvm use 24.11.1`
  - The Makefile (`make frontend-dev`) will automatically switch to the correct version
- Check port 3000 is not in use: `lsof -i :3000`
- Install dependencies: `cd src/frontend && npm install`
- If you see PostCSS/Tailwind errors, ensure you're using Node.js 18+

### Database connection issues
- Start MongoDB: `make db-start`
- Check MongoDB logs: `make db-logs`
- Verify connection string in backend `.env` file

## Stopping Services

### Stop Backend
Press `Ctrl+C` in the backend terminal

### Stop Frontend
Press `Ctrl+C` in the frontend terminal

### Stop MongoDB
```bash
make db-stop
```

## Development Workflow

1. **Start MongoDB** (if not already running)
   ```bash
   make db-start
   ```

2. **Start Backend** (Terminal 1)
   ```bash
   make run
   ```

3. **Start Frontend** (Terminal 2)
   ```bash
   cd src/frontend && npm run dev
   ```

4. **Open Browser**
   - Navigate to: http://localhost:3000
   - You should see the Update Manager application

5. **Make Changes**
   - Frontend: Changes hot-reload automatically
   - Backend: Restart required for changes (`Ctrl+C` then `make run`)

## What You'll See (Phase 1)

After starting both services and opening http://localhost:3000, you should see:

1. **Header** with:
   - Update Manager logo
   - Search bar (placeholder)
   - Notification bell (placeholder)
   - User menu (placeholder)

2. **Sidebar** with navigation:
   - Dashboard
   - Products
   - Versions
   - Updates
   - Notifications
   - Audit Logs

3. **Main Content Area** showing:
   - Dashboard page (placeholder content)
   - Navigation between pages works
   - Active route highlighting

4. **All pages are placeholders** showing:
   - "Dashboard" - "Dashboard content will be implemented in Phase 9"
   - "Products" - "Product management will be implemented in Phase 2"
   - "Versions" - "Version management will be implemented in Phase 3"
   - etc.

## Next Steps

- Phase 1 is complete - foundation is ready
- Phase 2 will implement Product Management features
- All API endpoints are ready and connected
- Frontend can make API calls (though pages are placeholders)

