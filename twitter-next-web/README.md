# Twitter Next Web

Frontend for the Twitter clone, built with Next.js App Router.

## Stack

- Next.js 16
- React 19
- TypeScript
- Tailwind CSS
- TanStack Query
- Zustand
- Zod

## API Expectations

This frontend expects the Go API (`/api/v1`) and cursor-based pagination:

- Request params: `cursor`, `size`
- Response model:
  - `items`
  - `hasNext`
  - `nextCursor`

## Local Setup

### Prerequisites

- Node.js 20+
- npm

### Install + Run

```bash
npm install
npm run dev
```

Open: `http://localhost:3000`

### Environment

Create `.env.local`:

```env
NEXT_PUBLIC_API_URL=http://localhost:8080/api/v1
NEXT_PUBLIC_GOOGLE_CLIENT_ID=your-google-client-id
```

## Quality Checks

```bash
npx tsc --noEmit
npm run lint
```

## Notes

- A `postinstall` script patches a `react-hook-form` type export issue in installed dependencies to keep type-checking stable in this project environment.
- Notification UI supports SSE updates from backend notification stream.

## Main Directories

- `src/app`: routes/layouts
- `src/components`: UI components
- `src/hooks`: data hooks and client logic
- `src/store`: Zustand stores
- `src/types`: shared TypeScript contracts
- `src/lib`: utilities/validation/helpers
