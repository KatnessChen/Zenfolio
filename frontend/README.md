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

### Development with Docker

All development and testing should be done using Docker containers. See the main project README and Makefile for setup and available commands.

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

## 🔧 Configuration

### Environment Variables

See the example configuration in [`frontend/.env.example`](./.env.example).

### Path Aliases

The project uses TypeScript path mapping for clean imports:

- `@/*` maps to `./src/*`

Example:

```tsx
import { Button } from '@/components/ui/button'
import { useAppSelector } from '@/hooks/redux'
import { Transaction } from '@/types'
```

## � Storybook

Storybook is used for developing and testing UI components in isolation. It provides an interactive development environment for building and showcasing components.

### Running Storybook

To start the Storybook development server:

```bash
docker-compose exec frontend pnpm storybook
```

Storybook will be available at `http://localhost:6006`

## �📦 Build & Deployment

WIP

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
