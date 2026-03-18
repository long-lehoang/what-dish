# CLAUDE.md — Tối Nay Ăn Gì (Frontend)

## Project overview

A web app that helps users decide what to eat through a fun card shuffle random experience, detailed home cooking recipes, and real-time group voting. This is the Frontend repo, communicating with the Go backend via REST API + WebSocket.

**App name:** Tối Nay Ăn Gì (literally "What to eat tonight?")
**Tagline:** "Hết phân vân — lật là ăn!" (Stop hesitating — flip and eat!)

## Tech stack

- **Framework:** Next.js 14+ (App Router, Server Components)
- **Language:** TypeScript (strict mode)
- **Styling:** Tailwind CSS 3.4+
- **Animation:** Framer Motion 11+
- **State management:** Zustand
- **Real-time:** Socket.io Client
- **PWA:** next-pwa
- **Package manager:** pnpm
- **Linting:** ESLint + Prettier
- **Testing:** Vitest (unit) + React Testing Library (component) + Playwright (e2e)
- **Deployment:** Vercel

## Development workflow

Follow this loop for every feature or change:

```
┌─────────────────────────────────────────────────────────────────┐
│  1. IMPLEMENT                                                   │
│     Vibe code / build the feature                               │
│                                                                 │
│  2. QUALITY LOOP (repeat until all green)                       │
│     ┌─────────────────────────────────────────────────────┐     │
│     │  a. Add/update unit tests (Vitest)                  │     │
│     │  b. Add/update integration tests (RTL)              │     │
│     │  c. Add/update e2e tests (Playwright)               │     │
│     │  d. Run & verify all tests pass                     │     │
│     │     → pnpm test                                     │     │
│     │     → pnpm test:e2e                                 │     │
│     │  e. Review code against best practices              │     │
│     │     → SOLID, DRY, separation of concerns            │     │
│     │     → check component size (<150 lines)             │     │
│     │     → check for proper error/loading states         │     │
│     │     → check accessibility (aria, keyboard nav)      │     │
│     │  f. Run lint & CI tools                             │     │
│     │     → pnpm lint                                     │     │
│     │     → pnpm type-check                               │     │
│     │     → pnpm format:check                             │     │
│     │  g. If anything fails → fix and repeat from (a)     │     │
│     └─────────────────────────────────────────────────────┘     │
│                                                                 │
│  3. UPDATE DOCS                                                 │
│     - Update README.md if public API or setup changed           │
│     - Update this CLAUDE.md if architecture/conventions changed │
│     - Add/update diagrams (Mermaid in docs/) for any new        │
│       data flows, component trees, or state machines            │
│     - Update API contract if endpoints changed                  │
└─────────────────────────────────────────────────────────────────┘
```

### Test expectations by layer

| Layer | Tool | What to test | Coverage target |
|-------|------|-------------|-----------------|
| Unit | Vitest | Hooks, utils, stores, pure logic | 80%+ |
| Component | Vitest + RTL | Component render, user interaction, props | Key user flows |
| E2E | Playwright | Full user journeys (random → recipe, create room → vote) | Critical paths |

### CI checks (must all pass before merge)

```bash
pnpm type-check       # TypeScript strict compilation
pnpm lint             # ESLint (no warnings allowed)
pnpm format:check     # Prettier formatting
pnpm test             # Vitest unit + component tests
pnpm test:e2e         # Playwright e2e tests
pnpm build            # Production build succeeds
```

## Architecture

### Pattern: Feature-Sliced Design (adapted)

Organize code by feature domain, not by technical role. Each feature is self-contained with its own components, hooks, and logic, while shared infrastructure lives in common layers.

```
src/
├── app/                        # Routes only — thin wrappers that compose features
├── features/                   # Feature modules (the core of the app)
│   ├── random/                 # Card shuffle random feature
│   │   ├── components/         #   CardShuffle, CardReveal, ShuffleOverlay
│   │   ├── hooks/              #   useShuffle, useRerollLimit
│   │   ├── stores/             #   history-store (7-day random history)
│   │   ├── utils/              #   shuffle algorithm, animation configs
│   │   ├── types.ts            #   Feature-specific types
│   │   └── index.ts            #   Public API barrel export
│   ├── recipe/                 # Recipe detail + cook mode
│   │   ├── components/         #   RecipeHero, IngredientList, StepCard, CookMode, Timer
│   │   ├── hooks/              #   useTimer, useWakeLock, useServingScale
│   │   └── index.ts
│   ├── vote/                   # Group voting feature
│   │   ├── components/         #   RoomLobby, Tournament, SwipeVote, RankVote, ResultsScreen
│   │   ├── hooks/              #   useSocket, useVoteRoom, useCountdown
│   │   ├── stores/             #   vote-store (room state, live results)
│   │   └── index.ts
│   ├── explore/                # Browse & search dishes
│   │   ├── components/         #   DishGrid, SearchBar, FilterSheet
│   │   ├── hooks/              #   useDishSearch, useInfiniteScroll
│   │   └── index.ts
│   └── dish/                   # Shared dish domain (used across features)
│       ├── components/         #   DishCard, DishBadge, DishTags
│       ├── types.ts            #   Dish, Ingredient, Step interfaces
│       └── index.ts
├── shared/                     # Cross-cutting shared code
│   ├── ui/                     # Design system components (Button, Modal, Badge, Skeleton...)
│   ├── hooks/                  # Generic hooks (useDebounce, useMediaQuery, useLocalStorage)
│   ├── lib/                    # Infrastructure
│   │   ├── api-client.ts       #   Fetch wrapper + error handling (Adapter pattern)
│   │   ├── socket-client.ts    #   Socket.io singleton (Singleton pattern)
│   │   └── utils.ts            #   Pure utility functions
│   ├── providers/              # React context providers (Theme, Toast, etc.)
│   └── constants.ts            # App-wide constants
└── styles/
    └── globals.css
```

### Design patterns used

| Pattern | Where | Why |
|---------|-------|-----|
| **Adapter** | `shared/lib/api-client.ts` | Wraps fetch into a typed API client. Backend URL, auth headers, error mapping in one place. Easy to swap HTTP lib later. |
| **Singleton** | `shared/lib/socket-client.ts` | Single Socket.io instance shared across vote components. Lazy init on first use, auto-cleanup on disconnect. |
| **Observer** | Zustand stores | Components subscribe to state slices. Store changes trigger re-renders only in subscribed components. |
| **Strategy** | `features/vote/components/` | Three vote modes (Tournament, Swipe, Ranking) implement the same interface but different UX logic. The room page picks the strategy based on `voteType`. |
| **State Machine** | `features/vote/stores/vote-store.ts` | Vote room lifecycle: `idle → waiting → voting → finished`. Transitions enforced in the store, preventing invalid states. |
| **Compound Component** | `features/recipe/components/CookMode` | CookMode.Root, CookMode.Step, CookMode.Timer compose together. Parent manages state, children render UI slices. |
| **Container/Presenter** | All features | Page-level components (containers) fetch data and manage state. Inner components (presenters) are pure UI receiving props. |
| **Facade** | Feature `index.ts` barrels | Each feature exposes a clean public API via barrel exports. Other features import only from the barrel, never from internal files. |

### Data flow

```
[Next.js Server Component]
        │
        │ fetch on server (dishes list, dish detail)
        ▼
[API Client (Adapter)]  ─────────►  [Go Backend REST API]
        │
        │ passes data as props
        ▼
[Client Component]
        │
        ├── local state (useState)        → UI interactions
        ├── Zustand store (Observer)       → cross-component state
        └── Socket.io (Singleton)          → real-time vote updates
                │
                └──── WebSocket ──────────►  [Go Backend WS]
```

### Rendering strategy

| Route | Render type | Why |
|-------|------------|-----|
| `/` | SSG (static) | Homepage doesn't change often |
| `/random` | CSR (client) | Heavy animation, all client-side |
| `/dish/[slug]` | SSR → ISR (revalidate 1h) | SEO for recipes, cached |
| `/dish/[slug]/cook` | CSR (client) | Fullscreen interactive, no SEO needed |
| `/vote/*` | CSR (client) | Real-time, no SEO needed |
| `/explore` | SSR | Searchable, filterable, SEO matters |

## Project structure (full)

```
src/
├── app/                        # Next.js App Router (routes only)
│   ├── layout.tsx
│   ├── page.tsx                # → imports from features/random or shared
│   ├── random/page.tsx         # → imports from features/random
│   ├── dish/[slug]/page.tsx     # → imports from features/recipe
│   ├── dish/[slug]/cook/page.tsx
│   ├── vote/page.tsx           # → imports from features/vote
│   ├── vote/create/page.tsx
│   ├── vote/[roomId]/page.tsx
│   ├── explore/page.tsx        # → imports from features/explore
│   └── account/                # Phase 2
├── features/                   # See architecture section above
├── shared/                     # See architecture section above
├── styles/globals.css
└── __tests__/                  # E2e tests (Playwright)
    ├── random.spec.ts
    ├── recipe.spec.ts
    └── vote.spec.ts
```

## API communication

Backend runs at `NEXT_PUBLIC_API_URL` (default: `http://localhost:8080`). Dish content is sourced from a **Notion database** (CMS), synced and cached by the Go backend. The frontend never talks to Notion directly — all data comes through the backend REST API.

```
GET    /api/dishes              — List + filter
GET    /api/dishes/:slug        — Dish detail + ingredients + steps
GET    /api/dishes/random       — Random 1 dish (with filters + exclusion)
POST   /api/rooms               — Create vote room
GET    /api/rooms/:id           — Room info
POST   /api/rooms/:id/join      — Join room
POST   /api/rooms/:id/vote      — Submit vote
GET    /api/rooms/:id/results   — Vote results
WS     /ws/rooms/:id            — Real-time vote updates
```

## Key features & UX specs

### Card shuffle animation sequence
1. **Shuffle** (0.8s) — Mini cards scatter and shuffle rapidly
2. **Converge** (0.5s) — Cards gather to center, stacking up
3. **Select** (0.6s) — One card separates, scales up, 3D flip to reveal front
4. **Reveal** (0.5s) — Show dish image + name, glow effect + confetti particles
5. **Settle** — Selected card rests at center, remaining cards fade to background

Important: Animation MUST be smooth on mobile. Use only `transform` and `opacity` (GPU accelerated). Test on iPhone SE viewport (375px).

### Random rules
- Max 3 re-rolls → show message: "Thôi ăn cái này đi, đừng kén nữa 😄" (Just eat this, stop being picky!)
- Save 7-day history in Zustand (persist with localStorage) to avoid repeating dishes

### Cook Mode
- Fullscreen, large font, high contrast
- One step per screen, swipe or tap next
- Built-in timer via `useTimer` hook
- Wake Lock API (`useWakeLock` hook) to keep screen on while cooking

## Coding conventions

### TypeScript
- Strict mode enabled, never use `any` — use `unknown` and narrow types
- Interface for object shapes, Type for unions/intersections
- Props types defined inline above component, only extract to separate file if shared

### Components
- Functional components only, no class components
- One component per file, filename = PascalCase (e.g., `CardShuffle.tsx`)
- Co-locate component + tests in the same folder when complex
- Default export for page components, named export for shared components
- Keep components under 150 lines — extract sub-components if exceeding

### Naming
- Components: PascalCase (`DishCard`, `FilterBar`)
- Hooks: camelCase prefixed with `use` (`useShuffle`, `useWakeLock`)
- Utils/lib: camelCase (`formatCurrency`, `createApiClient`)
- Types: PascalCase (`Dish`, `VoteRoom`, `Ingredient`)
- Files: PascalCase for components, kebab-case for utils/lib
- CSS: Tailwind utility classes; avoid custom CSS unless needed for animation keyframes

### Tailwind
- Use consistent design tokens from the palette
- Mobile-first: write mobile classes first, use `md:` `lg:` for desktop
- Avoid `@apply` in CSS files — write utilities directly in JSX
- Dark mode: use `dark:` prefix, class strategy

## Color palette (Tailwind config)

```
primary:    #FF6B35  (orange — main accent)
secondary:  #E63946  (red)
accent:     #FFB703  (yellow)
background: #FFF8F0  (cream — light mode)
dark-bg:    #1A1A2E  (warm dark — dark mode)
dark-card:  #16213E  (card dark mode)
```

## Fonts

- **Headings:** Be Vietnam Pro (weight 500, 700)
- **Body:** Inter (weight 400, 500)
- Load via `next/font/google` in root layout

## Environment variables

```
NEXT_PUBLIC_API_URL=http://localhost:8080
NEXT_PUBLIC_WS_URL=ws://localhost:8080
NEXT_PUBLIC_SUPABASE_URL=
NEXT_PUBLIC_SUPABASE_ANON_KEY=
```

## Commands

```bash
pnpm dev              # Dev server (port 3000)
pnpm build            # Production build
pnpm start            # Start production server
pnpm lint             # ESLint (zero warnings)
pnpm format           # Prettier format all files
pnpm format:check     # Prettier check (CI)
pnpm type-check       # TypeScript strict check (tsc --noEmit)
pnpm test             # Vitest unit + component tests
pnpm test:watch       # Vitest watch mode
pnpm test:coverage    # Vitest with coverage report
pnpm test:e2e         # Playwright e2e tests
pnpm test:e2e:ui      # Playwright with UI mode
```

## Important notes

- All UI text is in **Vietnamese**. All code, comments, variable names, and commit messages are in **English**.
- Dish content is managed in **Notion** — editors add/edit dishes there, backend syncs to PostgreSQL cache. Frontend only consumes via API.
- Dish images may come from Notion (S3-signed URLs) or Supabase Storage; use `next/image` with appropriate loader config and handle both sources.
- SEO: every page must export `metadata` (title, description, og:image).
- Responsive: mobile-first. Breakpoints: sm(640) md(768) lg(1024). Prioritize testing on 375px mobile viewport.
- Accessibility: semantic HTML, aria-labels on interactive elements, keyboard navigation for card shuffle.
- Diagrams are maintained in `docs/` folder using Mermaid format. Update them whenever architecture or data flows change.