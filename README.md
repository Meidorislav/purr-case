# PurrCase

![PurrCase](frontend/assets/logo.svg)

Cat-themed loot case marketplace. Browse cases, open them for randomized rewards, and manage your inventory – all powered by Xsolla payments.

## Preview

![[catalog.png]]

![[cart.png]]

![[inventory.png]]

![[opencase.png]]
## Team

| Name             | Role                             |
| ---------------- | -------------------------------- |
| Vladislav Bakin  | Leader, Items, Open case, Deploy |
| Alina Selivanova | Login, Graphic Designer          |
| Elena Yakina     | Frontend                         |
| Matvey Kedrov    | Payments, Inventory              |

## Stack

**Backend:** Go, Chi, PostgreSQL, golang-migrate
**Frontend:** React 18, TypeScript, Redux Toolkit, Vite
**Infrastructure:** Docker, Docker Compose, nginx

## Getting Started

```bash
cp .env.example .env
# Fill in Xsolla credentials in .env
docker-compose up --build
```

| Service  | URL                   |
| -------- | --------------------- |
| Frontend | http://localhost:5173 |
| Backend  | http://localhost:8080 |
| Database | localhost:5432        |

### Production

```bash
docker-compose -f docker-compose.prod.yaml up --build
```

Serves on ports 80/443 via nginx with SSL.

## Environment Variables

| Variable                    | Description                          |
| --------------------------- | ------------------------------------ |
| `PORT`                      | Backend port (default: `8080`)       |
| `merchant_id`               | Xsolla merchant ID                   |
| `XSOLLA_PROJECT_ID`         | Xsolla project ID                    |
| `XSOLLA_API_KEY`            | Xsolla server API key                |
| `XSOLLA_WEBHOOK_SECRET_KEY` | Webhook signature secret             |
| `XSOLLA_SANDBOX`            | Enable sandbox mode (`true`/`false`) |
| `XSOLLA_RETURN_URL`         | Redirect URL after payment           |
| `VITE_XSOLLA_CLIENT_ID`     | Xsolla client ID (frontend)          |
| `VITE_XSOLLA_LOGIN_ID`      | Xsolla login ID (frontend)           |
| `VITE_XSOLLA_RETURN_URL`    | Redirect URL after login (frontend)  |
| `DB_USER`                   | PostgreSQL user                      |
| `DB_PASSWORD`               | PostgreSQL password                  |
| `DB_NAME`                   | PostgreSQL database name             |
| `DB_PORT`                   | PostgreSQL port (default: `5432`)    |

## API

All authenticated endpoints require a Bearer JWT token.

| Method | Path | Auth | Description |
|--------|------|------|-------------|
| `GET` | `/health` | — | Health check |
| `GET` | `/me` | required | Current user info |
| `GET` | `/items` | optional | Item catalog |
| `GET` | `/items/sku/{sku}` | optional | Item by SKU |
| `GET` | `/items/virtual_items` | optional | Virtual currency items |
| `GET` | `/inventory` | required | User inventory |
| `GET` | `/inventory/{sku}` | required | Quantity of specific item |
| `POST` | `/inventory/consume` | required | Consume an item |
| `POST` | `/inventory/unpack` | required | Unpack a case/bundle |
| `POST` | `/cases/{sku}/open` | required | Open a case, receive reward |
| `POST` | `/payments/checkout` | required | Create checkout session |
| `POST` | `/payments/webhook` | — | Xsolla webhook receiver |

### Open a case

```
POST /cases/{sku}/open
→ { caseSKU, wonItem }
```

### Checkout

```
POST /payments/checkout
Body: { items: [{ sku: string, quantity: number }] }
→ { orderId, status, itemsCount, checkoutUrl }
```

## Project Structure

```
purr-case/
├── backend/
│   ├── cmd/purr-case/       # Entry point
│   ├── internal/
│   │   ├── db/              # Database connection
│   │   ├── dto/             # Data Transfer Objects
│   │   ├── httpapi/         # HTTP handlers & routing
│   │   ├── migrations/      # SQL migrations
│   │   └── service/         # Business logic
│   └── data/items.json      # Static item catalog
├── frontend/
│   └── src/
│       ├── pages/           # Route pages
│       ├── widgets/         # Reusable components
│       ├── shared/          # Hooks, store, UI primitives
│       └── app/             # Redux store config
├── docker-compose.yaml
├── docker-compose.prod.yaml
└── .env.example
```

## Frontend Scripts

```bash
npm run dev      # Development server
npm run build    # Production build
npm run preview  # Preview production build
npm run lint     # Run ESLint
```
