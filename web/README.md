# FATE Vault Web

Vue.js frontend application for FATE Vault.

## Setup

1. Install dependencies:
```bash
npm install
```

2. Start development server:
```bash
npm run dev
```

The application will be available at `http://localhost:3000`

## Features

- **Characters Page**: Browse characters with pagination (10 per page)
- **Character Detail/Edit**: View and edit character details with save functionality
- **Stunts Page**: Placeholder for future implementation
- **User Page**: Placeholder for future implementation

## API Integration

The web service connects to the backend API running on `http://localhost:8080`. Make sure the backend service is running before using the web application.

The API proxy is configured in `vite.config.js` to forward `/api/*` requests to the backend.

## Build

To build for production:
```bash
npm run build
```

The built files will be in the `dist` directory.


