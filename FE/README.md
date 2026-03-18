# Tối Nay Ăn Gì — Frontend

> "Hết phân vân — lật là ăn!" (Stop hesitating — flip and eat!)

A Vietnamese food decision web app that helps users pick what to eat through a fun card shuffle experience, detailed home cooking recipes, and real-time group voting.

## Tech Stack

| Category | Technology |
|----------|-----------|
| Framework | Next.js 15 (App Router, Server Components) |
| Language | TypeScript (strict mode) |
| Styling | Tailwind CSS 3.4 |
| Animation | Framer Motion 11 |
| State | Zustand 5 |
| Real-time | Socket.io Client |
| Testing | Vitest + React Testing Library + Playwright |
| Package Manager | pnpm |

## Getting Started

### Prerequisites

- Node.js 18+
- pnpm 8+

### Install & Run

```bash
# Install dependencies
pnpm install

# Start dev server (port 3000)
pnpm dev

# Production build
pnpm build && pnpm start
```

### Environment Variables

Create a `.env.local` file (see `.env.example`):

```env
NEXT_PUBLIC_API_URL=http://localhost:8080
NEXT_PUBLIC_WS_URL=ws://localhost:8080
NEXT_PUBLIC_SUPABASE_URL=
NEXT_PUBLIC_SUPABASE_ANON_KEY=
```

> **Note:** The app works fully offline with mock data when the backend is unavailable. All API calls gracefully fall back to built-in mock responses.

## Commands

```bash
pnpm dev              # Dev server with Turbopack
pnpm build            # Production build
pnpm start            # Start production server
pnpm lint             # ESLint (zero warnings)
pnpm format           # Prettier format all files
pnpm format:check     # Prettier check (CI)
pnpm type-check       # TypeScript strict check
pnpm test             # Vitest unit + component tests
pnpm test:watch       # Vitest watch mode
pnpm test:coverage    # Vitest with coverage report
pnpm test:e2e         # Playwright e2e tests
pnpm test:e2e:ui      # Playwright with UI mode
```

## Features

### Random Dish Picker
Animated card shuffle that randomly selects a dish for you. Includes category/difficulty/time filters, a visual dish pool, and a max 3 re-roll limit with playful Vietnamese messages.

### Recipe & Cook Mode
Detailed recipe pages with ingredient lists (serving scaling), step-by-step instructions, and a fullscreen Cook Mode with built-in timers and Wake Lock API support.

### Group Voting
Create a room, share the code, and vote together. Three voting strategies:
- **Tournament** — 1v1 bracket elimination
- **Swipe** — Tinder-style like/dislike
- **Ranking** — Drag-and-drop ordering

Works fully client-side with mock data when no backend is connected.

### Explore
Browse and search dishes with filters. Infinite scroll, debounced search, and responsive grid layout.

## Architecture

**Pattern:** Feature-Sliced Design — code organized by feature domain, not technical role.

```
src/
├── app/                  # Routes (thin wrappers composing features)
├── features/
│   ├── dish/             # Shared dish domain (DishCard, types)
│   ├── random/           # Card shuffle (CardShuffle, DishPool, FilterBar)
│   ├── recipe/           # Recipe detail + CookMode compound component
│   ├── vote/             # Room lobby, 3 vote strategies, results
│   └── explore/          # Search, filters, infinite scroll grid
├── shared/
│   ├── ui/               # Design system (Button, Modal, Badge, Skeleton, Navbar)
│   ├── hooks/            # useDebounce, useMediaQuery, useLocalStorage
│   ├── lib/              # api-client (Adapter), socket-client (Singleton), utils
│   ├── providers/        # ThemeProvider, ToastProvider
│   └── constants.ts
└── styles/globals.css
```

### Design Patterns

| Pattern | Usage |
|---------|-------|
| **Adapter** | `api-client.ts` — typed fetch wrapper with mock fallback |
| **Singleton** | `socket-client.ts` — shared Socket.io instance |
| **Strategy** | Vote components — Tournament / Swipe / Ranking |
| **State Machine** | Vote store — `idle → waiting → voting → finished` |
| **Compound Component** | CookMode — Root, Step, Timer, Controls |
| **Observer** | Zustand stores — selective re-rendering via slices |

### Rendering Strategy

| Route | Type | Reason |
|-------|------|--------|
| `/` | SSG | Static homepage |
| `/random` | CSR | Heavy client-side animation |
| `/dish/[slug]` | ISR (1h) | SEO for recipes |
| `/dish/[slug]/cook` | CSR | Fullscreen interactive |
| `/vote/*` | CSR | Real-time, no SEO needed |
| `/explore` | SSR | Searchable, SEO matters |

## Color Palette

| Token | Hex | Usage |
|-------|-----|-------|
| `primary` | `#FF6B35` | Main accent (orange) |
| `secondary` | `#E63946` | Red accent |
| `accent` | `#FFB703` | Yellow highlights |
| `background` | `#FFF8F0` | Light mode background (cream) |
| `dark-bg` | `#1A1A2E` | Dark mode background |
| `dark-card` | `#16213E` | Dark mode cards |

## API Endpoints

The frontend communicates with a Go backend via REST + WebSocket:

```
GET    /api/dishes              — List + filter dishes
GET    /api/dishes/:slug        — Dish detail
GET    /api/dishes/random       — Random dish (with filters)
POST   /api/rooms               — Create vote room
GET    /api/rooms/:id           — Room info
POST   /api/rooms/:id/join      — Join room
POST   /api/rooms/:id/vote      — Submit vote
GET    /api/rooms/:id/results   — Vote results
WS     /ws/rooms/:id            — Real-time vote updates
```

## Mock Data Mode

When the backend is unavailable, the app automatically falls back to mock data:
- **Dishes**: 10 pre-built Vietnamese dishes with images, ingredients, and steps
- **Vote rooms**: Full room lifecycle simulated client-side
- **Results**: Generated locally based on vote type and user selections

No configuration needed — the `api-client` catches network errors and returns mock responses transparently.

## License

Private project.
