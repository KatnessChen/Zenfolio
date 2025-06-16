# Transaction Tracker Frontend

A modern React application for managing trading transactions with AI-powered image processing capabilities.

## 🚀 Tech Stack

- **Framework**: React 19 with TypeScript
- **Build Tool**: Vite 6
- **Package Manager**: pnpm
- **Routing**: React Router v7
- **State Management**: Redux Toolkit + RTK Query
- **UI Components**: shadcn/ui + Tailwind CSS
- **Code Quality**: ESLint + Prettier
- **Development**: Hot reload, fast refresh, TypeScript strict mode

## 📁 Project Structure

```
src/
├── components/          # Reusable UI components
│   ├── ui/             # shadcn/ui components
│   └── Navigation.tsx  # App navigation
├── pages/              # Route-level components
│   ├── HomePage.tsx
│   ├── TransactionExtractPage.tsx
│   └── TransactionHistoryPage.tsx
├── store/              # Redux store configuration
│   └── index.ts
├── hooks/              # Custom React hooks
│   └── redux.ts        # Typed Redux hooks
├── types/              # TypeScript type definitions
│   └── index.ts
├── utils/              # Utility functions
│   └── index.ts
├── lib/                # shadcn/ui utilities
│   └── utils.ts
└── App.tsx             # Main application component
```

## 🛠️ Development Setup

### Prerequisites

- Node.js 18+
- pnpm (recommended package manager)

### Installation

1. **Install dependencies**:

   ```bash
   pnpm install
   ```

2. **Set up environment variables**:

   ```bash
   cp .env.example .env
   ```

   Edit `.env` to configure your backend API URL.

3. **Start development server**:

   ```bash
   pnpm dev
   ```

   The app will be available at `http://localhost:5173`

### Available Scripts

- `pnpm dev` - Start development server with hot reload
- `pnpm build` - Build for production
- `pnpm preview` - Preview production build locally
- `pnpm lint` - Run ESLint
- `pnpm lint:fix` - Fix ESLint issues automatically
- `pnpm format` - Format code with Prettier
- `pnpm format:check` - Check code formatting
- `pnpm type-check` - Run TypeScript type checking

## 🎨 UI Components

This project uses [shadcn/ui](https://ui.shadcn.com/) components built on top of Tailwind CSS. Components are located in `src/components/ui/` and can be imported using the `@/` alias:

```tsx
import { Button } from '@/components/ui/button'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
```

### Adding New Components

Add new shadcn/ui components:

```bash
pnpm dlx shadcn@latest add [component-name]
```

## 🗂️ State Management

The app uses Redux Toolkit for state management with typed hooks:

```tsx
import { useAppDispatch, useAppSelector } from '@/hooks/redux'

// In your component
const dispatch = useAppDispatch()
const data = useAppSelector((state) => state.someSlice.data)
```

## 🔗 API Integration

API calls should be configured in the store using RTK Query. The base URL is configured via environment variables:

```typescript
// In your API slice
import { createApi, fetchBaseQuery } from '@reduxjs/toolkit/query/react'

export const apiSlice = createApi({
  reducerPath: 'api',
  baseQuery: fetchBaseQuery({
    baseUrl: import.meta.env.VITE_API_BASE_URL,
  }),
  endpoints: (builder) => ({
    // Define your endpoints here
  }),
})
```

## 🎯 Features

- **Transaction Extraction**: Upload trading screenshots to extract transaction data using AI
- **Transaction History**: View and manage extracted transactions
- **Responsive Design**: Mobile-friendly interface
- **Dark/Light Theme**: Built-in theme support via shadcn/ui
- **Type Safety**: Full TypeScript support with strict mode

## 🔧 Configuration

### Environment Variables

| Variable            | Description          | Default                 |
| ------------------- | -------------------- | ----------------------- |
| `VITE_API_BASE_URL` | Backend API base URL | `http://localhost:8080` |
| `VITE_APP_NAME`     | Application name     | `Transaction Tracker`   |

### Path Aliases

The project uses TypeScript path mapping for clean imports:

- `@/*` maps to `./src/*`

Example:

```tsx
import { Button } from '@/components/ui/button'
import { useAppSelector } from '@/hooks/redux'
import { Transaction } from '@/types'
```

## 📦 Build & Deployment

1. **Build for production**:

   ```bash
   pnpm build
   ```

2. **Preview production build**:
   ```bash
   pnpm preview
   ```

The built files will be in the `dist/` directory, ready for deployment to any static hosting service.

## 🤝 Development Guidelines

- Use TypeScript strictly - avoid `any` types
- Follow the established folder structure
- Use typed Redux hooks (`useAppDispatch`, `useAppSelector`)
- Implement responsive design with Tailwind CSS
- Add proper error handling for API calls
- Write meaningful commit messages

## 📚 Additional Resources

- [Vite Documentation](https://vitejs.dev/)
- [React Router Documentation](https://reactrouter.com/)
- [Redux Toolkit Documentation](https://redux-toolkit.js.org/)
- [shadcn/ui Documentation](https://ui.shadcn.com/)
- [Tailwind CSS Documentation](https://tailwindcss.com/)
