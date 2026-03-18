# API Design — WhatDish

## Conventions

- **Base URL:** `https://api.whatdish.app/v1`
- **Format:** JSON
- **Auth:** Bearer JWT token in header `Authorization: Bearer <token>`
- **Pagination:** `?page=1&limit=20` (default limit=20, max=100)
- **Sorting:** `?sort=created_at&order=desc`
- **Error format:**
  ```json
  {
    "error": {
      "code": "RECIPE_NOT_FOUND",
      "message": "Recipe with id xxx not found",
      "status": 404
    }
  }
  ```

### Auth Levels

| Symbol | Meaning |
|--------|---------|
| 🔓 | Public — no auth required |
| 🔐 | User — JWT token required |
| 🔑 | Admin — JWT + role ADMIN required |

---

## 1. Recipe Service (`/recipes`)

> **Note:** Recipe content is managed in Notion and synced to PostgreSQL.
> There are no POST/PUT/DELETE endpoints for recipes — use the admin sync endpoint instead.

### 🔓 `GET /recipes`

List dishes (with filters, pagination).

**Query Params:**
```
page          int     (default: 1)
limit         int     (default: 20, max: 100)
sort          string  (name, cook_time, created_at, favorite_count)
order         string  (asc, desc)
dish_type     string  (slug: braised, stir-fried, steamed, grilled...)
region        string  (slug: northern, central, southern)
main_ingredient string (slug: chicken, beef, pork, seafood, vegetarian)
meal_type     string  (slug: breakfast, lunch, dinner, snack)
difficulty    string  (EASY, MEDIUM, HARD)
max_cook_time int     (minutes)
tags          string  (comma-separated slugs)
exclude_ingredients string (comma-separated ingredient IDs — for allergies)
```

**Response 200:**
```json
{
  "data": [
    {
      "id": "uuid",
      "name": "Pho Bo Ha Noi",
      "slug": "pho-bo-ha-noi",
      "image_url": "https://...",
      "cook_time": 120,
      "difficulty": "MEDIUM",
      "servings": 4,
      "region": { "id": "uuid", "name": "Northern" },
      "dish_type": { "id": "uuid", "name": "Soup" },
      "nutrition_summary": {
        "calories": 480,
        "protein": 28
      },
      "favorite_count": 234,
      "tags": ["Traditional", "Broth"]
    }
  ],
  "pagination": {
    "page": 1,
    "limit": 20,
    "total": 156,
    "total_pages": 8
  }
}
```

### 🔓 `GET /recipes/:id`

Full recipe details.

**Response 200:**
```json
{
  "data": {
    "id": "uuid",
    "name": "Pho Bo Ha Noi",
    "slug": "pho-bo-ha-noi",
    "description": "Traditional Vietnamese beef noodle soup...",
    "image_url": "https://...",
    "prep_time": 30,
    "cook_time": 120,
    "total_time": 150,
    "servings": 4,
    "difficulty": "MEDIUM",
    "region": { "id": "uuid", "name": "Northern", "slug": "northern" },
    "dish_type": { "id": "uuid", "name": "Soup", "slug": "soup" },
    "tags": ["Traditional", "Broth"],
    "ingredients": [
      {
        "id": "uuid",
        "ingredient": { "id": "uuid", "name": "Beef bones" },
        "amount": 1.5,
        "unit": "kg",
        "note": "cleaned, blanched in boiling water",
        "is_optional": false,
        "group_name": "For the broth"
      }
    ],
    "steps": [
      {
        "step_number": 1,
        "title": "Prepare bones",
        "description": "Clean bones, blanch in boiling water for 5 minutes...",
        "image_url": "https://...",
        "duration": 10
      }
    ],
    "nutrition": {
      "calories": 480,
      "protein": 28,
      "carbs": 52,
      "fat": 16,
      "fiber": 2,
      "per_serving": true
    },
    "favorite_count": 234,
    "view_count": 1890
  }
}
```

### 🔓 `GET /recipes/random`

Random 1 dish.

**Query Params:**
```
exclude_ids    string  (comma-separated UUIDs — dishes to skip)
dish_type      string
region         string
difficulty     string
max_cook_time  int
```

**Response 200:** (same format as GET /recipes/:id)

### 🔓 `GET /recipes/search`

Full-text search (proxied to Elasticsearch).

**Query Params:**
```
q       string  (search keyword — required)
page    int
limit   int
```

**Response 200:** (same format as GET /recipes, with additional `score` field)

### 🔓 `GET /recipes/autocomplete`

Search suggestions (fast, lightweight).

**Query Params:**
```
q       string  (keyword — min 2 characters)
limit   int     (default: 5, max: 10)
```

**Response 200:**
```json
{
  "data": [
    { "id": "uuid", "name": "Pho Bo", "type": "recipe" },
    { "id": "uuid", "name": "Shrimp", "type": "ingredient" }
  ]
}
```

### 🔑 `POST /recipes`

Create new recipe (Admin).

**Request Body:**
```json
{
  "name": "Bun Bo Hue",
  "description": "...",
  "image_url": "https://...",
  "prep_time": 20,
  "cook_time": 90,
  "servings": 4,
  "difficulty": "MEDIUM",
  "dish_type_id": "uuid",
  "region_id": "uuid",
  "main_ingredient_id": "uuid",
  "meal_type_id": "uuid",
  "tag_ids": ["uuid", "uuid"],
  "ingredients": [
    {
      "ingredient_id": "uuid",
      "amount": 500,
      "unit": "g",
      "note": "sliced",
      "is_optional": false,
      "group_name": "Meat"
    }
  ],
  "steps": [
    {
      "step_number": 1,
      "title": "Preparation",
      "description": "...",
      "duration": 15
    }
  ],
  "nutrition": {
    "calories": 520,
    "protein": 32,
    "carbs": 48,
    "fat": 20,
    "fiber": 3
  },
  "status": "PUBLISHED"
}
```

**Response 201:** (recipe object)

### 🔑 `PUT /recipes/:id`

Update recipe (Admin). Same body as POST, supports partial update.

### 🔑 `DELETE /recipes/:id`

Soft delete recipe (Admin). **Response 204.**

### 🔓 `GET /categories`

**Query Params:** `type` (DISH_TYPE, REGION, MAIN_INGREDIENT, MEAL_TYPE)

**Response 200:**
```json
{
  "data": [
    { "id": "uuid", "name": "Braised", "slug": "braised", "type": "DISH_TYPE", "icon_url": "..." }
  ]
}
```

---

## 2. Suggestion Service (`/suggestions`)

### 🔓 `POST /suggestions/random`

Smart random (excludes history if authenticated).

**Request Body:**
```json
{
  "filters": {
    "dish_type": "stir-fried",
    "region": "southern",
    "max_cook_time": 30,
    "exclude_ingredients": ["uuid-shellfish"]
  }
}
```

**Response 200:**
```json
{
  "data": {
    "session_id": "uuid",
    "recipe": { ... },
    "suggestion_type": "RANDOM"
  }
}
```

### 🔓 `POST /suggestions/by-calories`

Calorie-based suggestion.

**Request Body:**
```json
{
  "target_calories": 500,
  "meal_type": "lunch",
  "tolerance_pct": 15,
  "filters": {
    "region": "northern",
    "exclude_ingredients": []
  }
}
```

**Response 200:**
```json
{
  "data": {
    "session_id": "uuid",
    "meal_type": "lunch",
    "target_calories": 500,
    "recipes": [
      {
        "recipe": { ... },
        "calories": 480,
        "role": "main"
      }
    ],
    "total_calories": 480,
    "calorie_diff": -20,
    "suggestion_type": "BY_CALORIES"
  }
}
```

### 🔓 `POST /suggestions/by-group`

Group meal suggestion.

**Request Body:**
```json
{
  "group_size": 4,
  "group_type": "family_with_kids",
  "meal_type": "dinner",
  "budget": "moderate",
  "filters": {
    "exclude_ingredients": ["uuid-spicy-chili"]
  }
}
```

**Response 200:**
```json
{
  "data": {
    "session_id": "uuid",
    "group_size": 4,
    "group_type": "family_with_kids",
    "dishes": [
      { "recipe": { ... }, "role": "main", "scaled_servings": 4 },
      { "recipe": { ... }, "role": "soup", "scaled_servings": 4 },
      { "recipe": { ... }, "role": "side", "scaled_servings": 4 }
    ],
    "shopping_list": [
      { "ingredient": "Chicken", "total_amount": 800, "unit": "g" },
      { "ingredient": "Tomato", "total_amount": 4, "unit": "pieces" }
    ],
    "total_nutrition": {
      "calories_per_person": 650,
      "protein_per_person": 38
    },
    "estimated_cook_time": 45,
    "suggestion_type": "BY_GROUP"
  }
}
```

### 🔐 `GET /suggestions/history`

User's suggestion history.

**Query Params:**
```
page            int
limit           int
session_type    string  (RANDOM, BY_CALORIES, BY_GROUP)
from_date       string  (ISO 8601)
to_date         string  (ISO 8601)
```

**Response 200:**
```json
{
  "data": [
    {
      "session_id": "uuid",
      "session_type": "BY_GROUP",
      "input_params": { "group_size": 4 },
      "recipes": [
        { "id": "uuid", "name": "Thit Kho Tau", "image_url": "..." }
      ],
      "total_calories": 650,
      "created_at": "2025-03-18T12:00:00Z"
    }
  ],
  "pagination": { ... },
  "stats": {
    "total_suggestions": 45,
    "most_suggested_recipe": { "id": "uuid", "name": "Canh Chua Ca", "count": 5 },
    "avg_daily_calories": 1650
  }
}
```

---

## 3. User Service (`/auth`, `/users`)

### 🔓 `POST /auth/register`

```json
{
  "email": "user@example.com",
  "password": "securePass123",
  "name": "Nguyen Van A"
}
```

**Response 201:**
```json
{
  "data": {
    "user": { "id": "uuid", "email": "...", "name": "..." },
    "message": "Please check your email to verify your account"
  }
}
```

### 🔓 `POST /auth/login`

```json
{
  "email": "user@example.com",
  "password": "securePass123"
}
```

**Response 200:**
```json
{
  "data": {
    "access_token": "eyJ...",
    "refresh_token": "eyJ...",
    "expires_in": 900,
    "user": { "id": "uuid", "email": "...", "name": "...", "role": "USER" }
  }
}
```

### 🔓 `POST /auth/refresh`

```json
{ "refresh_token": "eyJ..." }
```

### 🔓 `POST /auth/forgot-password`

```json
{ "email": "user@example.com" }
```

### 🔓 `POST /auth/reset-password`

```json
{ "token": "reset-token-from-email", "new_password": "newSecure123" }
```

### 🔐 `GET /users/me`

Get current user info.

### 🔐 `PUT /users/me/profile`

Update nutrition profile.

```json
{
  "gender": "MALE",
  "age": 28,
  "height_cm": 175,
  "weight_kg": 70,
  "activity_level": "MODERATE",
  "goal": "MAINTAIN",
  "allergies": [
    { "ingredient_id": "uuid", "allergy_type": "ALLERGY" }
  ]
}
```

**Response 200:**
```json
{
  "data": {
    "profile": { ... },
    "calculated": {
      "bmr": 1695,
      "tdee": 2627,
      "daily_target": 2627,
      "meal_targets": {
        "breakfast": 788,
        "lunch": 1051,
        "dinner": 788
      }
    }
  }
}
```

### 🔑 `GET /users`

List users (Admin). Supports search, filter, pagination.

### 🔑 `PUT /users/:id/status`

Suspend/reactivate account (Admin).

```json
{ "status": "SUSPENDED", "reason": "Policy violation" }
```

---

## 4. Nutrition Service (`/nutrition`)

### 🔓 `GET /nutrition/recipes/:recipe_id`

Nutrition for 1 recipe.

### 🔓 `GET /nutrition/recipes?ids=uuid1,uuid2,uuid3`

Batch query nutrition data.

### 🔓 `GET /nutrition/goals`

List preset nutrition goals.

### 🔓 `POST /nutrition/calculate-tdee`

Calculate TDEE without creating an account.

```json
{
  "gender": "FEMALE",
  "age": 25,
  "height_cm": 162,
  "weight_kg": 55,
  "activity_level": "LIGHT",
  "goal": "LOSE_WEIGHT"
}
```

**Response 200:**
```json
{
  "data": {
    "bmr": 1305,
    "tdee": 1794,
    "daily_target": 1435,
    "meal_breakdown": {
      "breakfast": { "calories": 431, "label": "30%" },
      "lunch": { "calories": 574, "label": "40%" },
      "dinner": { "calories": 431, "label": "30%" }
    }
  }
}
```

### 🔑 `POST /nutrition/recipes/:recipe_id`

Add/update nutrition data (Admin).

---

## 5. Engagement Service (`/favorites`, `/views`)

### 🔐 `POST /favorites`

```json
{ "recipe_id": "uuid" }
```

### 🔐 `DELETE /favorites/:recipe_id`

### 🔐 `GET /favorites`

**Query Params:** `page`, `limit`, `sort` (created_at, recipe_name)

### 🔐 `GET /favorites/check?recipe_ids=uuid1,uuid2`

Batch check if recipes are favorited (for UI icon display).

**Response 200:**
```json
{
  "data": {
    "uuid1": true,
    "uuid2": false
  }
}
```

### 🔓 `POST /views`

Record a view (no auth required, uses session_id).

```json
{
  "recipe_id": "uuid",
  "source": "RANDOM"
}
```

---

## 6. Rate Limiting

| Endpoint Pattern | Guest | User | Admin |
|---|---|---|---|
| `GET /recipes/*` | 60/min | 120/min | unlimited |
| `POST /suggestions/*` | 20/min | 60/min | unlimited |
| `POST /auth/login` | 10/min | — | — |
| `POST /auth/register` | 5/min | — | — |
| `GET /search/*` | 30/min | 60/min | unlimited |
| `POST /favorites` | — | 30/min | unlimited |

---

## 7. Versioning Strategy

API version is included in the URL path: `/v1/`, `/v2/`.

When a breaking change is needed: create a new version, keep the old version running for at least 6 months, and send a `Deprecation` header to notify clients.

```
Deprecation: true
Sunset: Sat, 01 Jan 2027 00:00:00 GMT
Link: <https://api.whatdish.app/v2/recipes>; rel="successor-version"
```