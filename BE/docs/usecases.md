# Use Cases — WhatDish Detailed Flow

## Table of Contents

- [Use Cases — WhatDish Detailed Flow](#use-cases--whatdish-detailed-flow)
  - [Table of Contents](#table-of-contents)
  - [UC01 — Random Dish](#uc01--random-dish)
    - [Main Flow](#main-flow)
    - [Alternative Flows](#alternative-flows)
    - [Postcondition](#postcondition)
  - [UC02 — Calorie-Based Suggestion](#uc02--calorie-based-suggestion)
    - [Main Flow](#main-flow-1)
    - [Alternative Flows](#alternative-flows-1)
    - [Postcondition](#postcondition-1)
  - [UC03 — Group Meal Suggestion](#uc03--group-meal-suggestion)
    - [Main Flow](#main-flow-2)
    - [Alternative Flows](#alternative-flows-2)
    - [Postcondition](#postcondition-2)
  - [UC04 — Filter Dishes](#uc04--filter-dishes)
    - [Main Flow](#main-flow-3)
    - [Alternative Flows](#alternative-flows-3)
  - [UC05 — View Recipe Details](#uc05--view-recipe-details)
    - [Main Flow](#main-flow-4)
    - [Alternative Flows](#alternative-flows-4)
  - [UC06 — Search Dishes](#uc06--search-dishes)
    - [Main Flow](#main-flow-5)
    - [Alternative Flows](#alternative-flows-5)
  - [UC07 — Register / Login](#uc07--register--login)
    - [Main Flow — Registration](#main-flow--registration)
    - [Main Flow — Login](#main-flow--login)
    - [Main Flow — OAuth](#main-flow--oauth)
    - [Alternative Flows](#alternative-flows-6)
  - [UC08 — Setup Nutrition Profile](#uc08--setup-nutrition-profile)
    - [Main Flow](#main-flow-6)
    - [Alternative Flows](#alternative-flows-7)
  - [UC09 — Save Favorites](#uc09--save-favorites)
    - [Main Flow](#main-flow-7)
    - [View Favorites List](#view-favorites-list)
  - [UC10 — Suggestion History](#uc10--suggestion-history)
    - [Main Flow](#main-flow-8)
  - [UC11 — Manage Recipes (Admin via Notion)](#uc11--manage-recipes-admin-via-notion)
    - [Main Flow — Content Management (in Notion)](#main-flow--content-management-in-notion)
    - [Main Flow — Sync to Application (Backend)](#main-flow--sync-to-application-backend)
    - [Alternative Flows](#alternative-flows-8)
  - [UC12 — Manage Users (Admin)](#uc12--manage-users-admin)
    - [Main Flow](#main-flow-9)
  - [UC13 — Manage Categories (Admin)](#uc13--manage-categories-admin)
    - [Main Flow](#main-flow-10)

---

## UC01 — Random Dish

**Actor:** Guest / User  
**Description:** The system randomly selects a dish from the database and displays an overview.  
**Precondition:** Database contains at least 1 recipe.

### Main Flow

```
1. Actor navigates to Home / Random page
2. Actor clicks "Random Dish" button (or "What should I eat today?")
3. System randomly selects 1 recipe from DB
   └─ If Actor is logged in: exclude dishes from the last 7-day history
4. System displays the result:
   ├─ Dish name
   ├─ Thumbnail image
   ├─ Basic info (cook time, difficulty, region)
   ├─ Nutrition summary (calories, protein, carbs, fat)
   └─ "View Recipe" button → UC05
5. Actor can:
   ├─ Click "Random Again" → return to step 3
   ├─ Click "View Recipe" → UC05
   ├─ Click "Save to Favorites" → UC09 (requires login)
   └─ Click "Share" → copy link or social share
6. System saves suggestion history (if logged in) → UC10
```

### Alternative Flows

- **3a.** DB is empty or connection error → display error message with retry button
- **3b.** All dishes already in 7-day history → reset filter, random from entire DB
- **5a.** Actor is not logged in and clicks "Save to Favorites" → show login prompt popup

### Postcondition

A random dish is displayed. If User is logged in, history is updated.

---

## UC02 — Calorie-Based Suggestion

**Actor:** Guest / User  
**Description:** Suggest dishes that match a target calorie intake for a meal or entire day.  
**Precondition:** Recipes have nutrition data (calories) populated.

### Main Flow

```
1. Actor navigates to "Suggest by Calories" page
2. System displays input form:
   ├─ Meal type: Breakfast / Lunch / Dinner / Full Day
   ├─ Target calories: enter number or select preset
   │   ├─ Weight loss: 300-400 kcal/meal
   │   ├─ Maintenance: 500-700 kcal/meal
   │   └─ Weight gain: 700-900 kcal/meal
   └─ (Optional) Additional filters: dish type, region, excluded ingredients
3. Actor clicks "Suggest"
4. System processes:
   ├─ If "Full Day" selected: calculate total calories → distribute breakfast(30%), lunch(40%), dinner(30%)
   ├─ Query recipes within allowed calorie range (±15% of target)
   ├─ Apply additional filters (if any)
   ├─ Random selection from filtered result set
   └─ If User is logged in + has nutrition profile → use TDEE from profile
5. System displays results:
   ├─ 1-3 suggested dishes (depending on meal type)
   ├─ Each dish: name, image, calories, macros (protein/carbs/fat)
   ├─ Total calories for the combo
   ├─ Progress bar comparing to target
   └─ "View Recipe" button for each dish → UC05
6. Actor can:
   ├─ Click "Suggest Again" → step 4 (random different set)
   ├─ Click "Replace Dish X" → re-random only that dish within equivalent calorie range
   └─ Click "Save Menu" → UC09
```

### Alternative Flows

- **4a.** No dishes found within calorie range → expand range to ±25%, if still none → show "No results found, please adjust your target"
- **4b.** User has allergy profile → automatically exclude allergen ingredients
- **2a.** User is logged in + has profile → auto-fill target calories from TDEE

### Postcondition

A list of dishes matching the calorie target is displayed. Suggestion history saved (if logged in).

---

## UC03 — Group Meal Suggestion

**Actor:** Guest / User  
**Description:** Suggest a balanced meal combo suitable for a group of people (family, friends, couple).  
**Precondition:** Database has sufficient variety of dishes.

### Main Flow

```
1. Actor navigates to "Group Suggestion" page
2. System displays form:
   ├─ Number of people: 2 / 4 / 6 / 8 / Custom
   ├─ Group type (presets):
   │   ├─ Date night → prioritize elegant, easy-to-eat dishes
   │   ├─ Family with kids → prioritize non-spicy, kid-friendly
   │   ├─ Friends gathering → prioritize appetizers, grilled items
   │   ├─ Birthday party → prioritize large combos, include cake/dessert
   │   └─ Custom
   ├─ Budget (optional): Budget / Moderate / No limit
   └─ Meal: Breakfast / Lunch / Dinner / Party
3. Actor clicks "Suggest"
4. System processes:
   ├─ Determine number of dishes needed by group size:
   │   ├─ 2 people: 2-3 dishes
   │   ├─ 4 people: 3-4 dishes
   │   ├─ 6+ people: 4-5 dishes
   │   └─ Always include: 1 main dish, 1 soup, 1+ side dishes
   ├─ Apply balancing logic:
   │   ├─ No duplicate main ingredients across dishes
   │   ├─ Diverse cooking methods (braised, stir-fried, steamed, fried, boiled)
   │   ├─ Balanced flavor profiles (salty, sweet, sour, spicy)
   │   └─ Region consistency (avoid mixing too many regions)
   ├─ Scale servings to match group size
   └─ Calculate estimated total nutrition
5. System displays:
   ├─ Meal combo (3-5 dishes) as card grid
   ├─ Each dish: name, image, cook time, calories/person
   ├─ Combo overview: total time, total ingredients, total calories/person
   ├─ Merged shopping list (combined ingredients across all dishes)
   └─ "View Recipe" button for each dish → UC05
6. Actor can:
   ├─ Click "Suggest Again" → step 4
   ├─ Click "Replace Dish X" → re-random only that dish
   ├─ Click "Export Shopping List" → download PDF/text
   └─ Click "Save Menu" → UC09
```

### Alternative Flows

- **4a.** Fewer than 3 matching dishes in DB → relax balancing criteria, suggest with fewer constraints
- **6a.** Actor drag-and-drops to reorder dishes → System recalculates total nutrition
- **2a.** "Family with kids" selected → automatically filter out spicy and raw dishes

### Postcondition

A balanced meal combo for the group is displayed along with a merged shopping list.

---

## UC04 — Filter Dishes

**Actor:** Guest / User  
**Description:** Filter the dish list by multiple combined criteria.

### Main Flow

```
1. Actor navigates to "Explore" page or uses filter bar
2. System displays filter options:
   ├─ Dish type: Braised, Stir-fried, Steamed, Grilled, Fried, Boiled, Salad, Soup, Hotpot...
   ├─ Region: Northern, Central, Southern
   ├─ Main ingredient: Chicken, Beef, Pork, Seafood, Vegetarian, Vegetables...
   ├─ Cook time: <15 min, 15-30 min, 30-60 min, >60 min
   ├─ Difficulty: Easy, Medium, Hard
   ├─ Calorie range: slider (0-1500 kcal)
   └─ Exclude ingredients: multi-select (allergies, dislikes)
3. Actor selects one or more filters → System updates list in real-time
4. System queries DB with AND conditions across all filters
5. System displays results:
   ├─ Grid/List view toggle
   ├─ Sort by: Popular, Newest, Calories (low→high), Cook time
   ├─ Pagination or infinite scroll
   └─ Each item: image, name, calories, cook time, rating
6. Actor clicks on a dish → UC05
```

### Alternative Flows

- **4a.** No results → display "No dishes found", suggest removing some filters
- **3a.** Actor clicks "Clear Filters" → reset to default, show all dishes

---

## UC05 — View Recipe Details

**Actor:** Guest / User  
**Description:** View the full recipe for a dish, including ingredients, steps, and nutrition.

### Main Flow

```
1. Actor clicks on a dish (from random, search, filter, or direct link)
2. System queries recipe detail from Recipe Service
3. System queries nutrition data from Nutrition Service
4. System displays detail page:
   ├─ Header:
   │   ├─ Dish name
   │   ├─ Hero image
   │   ├─ Average rating + number of reviews
   │   ├─ Tags: region, dish type, difficulty
   │   └─ Actions: Favorite ♡, Share, Print
   ├─ Overview:
   │   ├─ Prep time + cook time
   │   ├─ Servings (adjustable ±)
   │   └─ Difficulty level
   ├─ Nutrition (per serving):
   │   ├─ Calories
   │   ├─ Protein (g)
   │   ├─ Carbohydrates (g)
   │   ├─ Fat (g)
   │   └─ Fiber (g)
   ├─ Ingredients:
   │   ├─ Ingredient list with measurements
   │   ├─ Auto-scale based on adjusted servings
   │   └─ "Add to Shopping List" button
   └─ Steps:
       ├─ Step 1, 2, 3... with detailed instructions
       ├─ Illustrative images (if available)
       └─ Timer (for steps with wait time)
5. System records the view (analytics)
```

### Alternative Flows

- **2a.** Recipe does not exist (old URL, deleted) → display 404, suggest random dish instead
- **4a.** Nutrition data not available → display "Coming soon" instead of hiding the section

---

## UC06 — Search Dishes

**Actor:** Guest / User  
**Description:** Search for dishes by keyword (name, ingredient).

### Main Flow

```
1. Actor types keyword into search bar
2. System displays autocomplete suggestions (debounce 300ms):
   ├─ Match by dish name: "Pho" → "Pho Bo", "Pho Ga", "Pho Cuon"
   ├─ Match by ingredient: "shrimp" → "Salt & Pepper Shrimp", "Sour Soup with Shrimp"
   └─ Display max 5 suggestions
3. Actor presses Enter or selects a suggestion
4. System sends query to Search Service (Elasticsearch):
   ├─ Full-text search on: dish name, description, ingredient names, tags
   ├─ Support Vietnamese with and without diacritical marks
   ├─ Fuzzy matching (minor typo tolerance)
   └─ Ranking: relevance score + popularity
5. System displays results (same format as UC04 step 5)
6. Actor selects a result → UC05
```

### Alternative Flows

- **4a.** No results found → "No results for '[keyword]'", suggest similar keywords
- **1a.** Actor searches by ingredients ("shrimp, tomato") → system parses into multi-ingredient search

---

## UC07 — Register / Login

**Actor:** Guest  
**Description:** Create a new account or log into the system.

### Main Flow — Registration

```
1. Actor clicks "Register"
2. System displays form:
   ├─ Email
   ├─ Password (min 8 characters, must include letters + numbers)
   ├─ Confirm password
   └─ Display name
3. Actor fills in information → clicks "Register"
4. System validates:
   ├─ Email format is valid + not already registered
   ├─ Password meets strength requirements
   └─ Confirm password matches
5. System creates user, hashes password (bcrypt)
6. System sends verification email (link expires in 24h)
7. System publishes event: UserCreated → Message Bus
8. Redirect → Nutrition Profile setup page (UC08)
```

### Main Flow — Login

```
1. Actor clicks "Login"
2. System displays form: Email + Password
3. Actor clicks "Login"
4. System validates credentials
5. System returns JWT access token (15 min) + refresh token (7 days)
6. Redirect → previous page or home
```

### Main Flow — OAuth

```
1. Actor clicks "Login with Google/Facebook"
2. System redirects to OAuth provider
3. Provider callback → System receives profile info
4. System checks if email already exists:
   ├─ Yes → link account, login
   └─ No → create new user, login, redirect to UC08
```

### Alternative Flows

- **4a (Register).** Email already exists → "This email is already registered", suggest login instead
- **4a (Login).** Invalid credentials → "Incorrect email or password" (do not specify which one is wrong)
- **4b (Login).** More than 5 failed attempts → lock account for 15 minutes
- **Forgot Password:** Click "Forgot Password" → enter email → system sends reset link (expires in 1h)

---

## UC08 — Setup Nutrition Profile

**Actor:** User  
**Description:** Enter personal information so the system can calculate TDEE and personalize suggestions.  
**Precondition:** User is logged in.

### Main Flow

```
1. User navigates to "Nutrition Profile" (or redirected from registration)
2. System displays form:
   ├─ Gender: Male / Female
   ├─ Age
   ├─ Height (cm)
   ├─ Weight (kg)
   ├─ Activity level:
   │   ├─ Sedentary (office job)
   │   ├─ Light exercise (1-3 days/week)
   │   ├─ Moderate exercise (3-5 days/week)
   │   ├─ Active (6-7 days/week)
   │   └─ Very active (2x per day)
   ├─ Goal:
   │   ├─ Lose weight
   │   ├─ Gain weight / build muscle
   │   └─ Maintain
   └─ (Optional) Allergies / Dislikes: multi-select ingredients
3. User clicks "Save"
4. System calculates TDEE using Mifflin-St Jeor formula:
   ├─ BMR (Male)  = 10 × weight(kg) + 6.25 × height(cm) − 5 × age + 5
   ├─ BMR (Female) = 10 × weight(kg) + 6.25 × height(cm) − 5 × age − 161
   ├─ TDEE = BMR × Activity Multiplier
   └─ Adjust: Lose weight (−20%), Gain weight (+15%), Maintain (±0%)
5. System saves profile + TDEE to User DB
6. System publishes event: UserProfileUpdated → Message Bus
7. System displays results:
   ├─ Calculated TDEE: "You need approximately X kcal/day"
   ├─ Suggested distribution: Breakfast (30%), Lunch (40%), Dinner (30%)
   └─ "Start Getting Suggestions" button → UC02
```

### Alternative Flows

- **3a.** Required field missing → highlight error field
- **1a.** User already has a profile → display current info, allow editing

---

## UC09 — Save Favorites

**Actor:** User  
**Precondition:** User is logged in.

### Main Flow

```
1. User clicks ♡ icon on any dish card
2. System sends POST /favorites request with recipe_id
3. System saves to Engagement DB
4. UI toggles icon ♡ → ♥ (filled)
5. System publishes event: RecipeFavorited → Message Bus
   └─ Suggestion Service listens → updates preference model
```

### View Favorites List

```
1. User navigates to "Favorites"
2. System queries Engagement Service → favorites list
3. System batch queries Recipe Service → fetch detailed info
4. Display grid of saved dishes
   ├─ Sort by: Newest, Name A-Z
   ├─ Filter by type, region
   └─ Remove from favorites: click ♥ → toggle back
```

---

## UC10 — Suggestion History

**Actor:** User  
**Precondition:** User is logged in.

### Main Flow

```
1. User navigates to "Suggestion History"
2. System queries Suggestion Service → fetch suggestion_history
3. System displays timeline:
   ├─ Grouped by date
   ├─ Each entry: timestamp, suggestion type (random/calories/group), dish list
   ├─ Filter: by suggestion type, by date range
   └─ Click on a dish → UC05
4. Statistics:
   ├─ Most frequently suggested dish
   ├─ Preferred dish types
   └─ Average daily calories (if UC02 is used frequently)
```

---

## UC11 — Manage Recipes (Admin via Notion)

**Actor:** Admin  
**Precondition:** Admin has access to the Notion workspace.

### Main Flow — Content Management (in Notion)

```
1. Admin opens the shared Notion database "WhatDish Recipes"
2. Admin can:
   ├─ Add new recipe:
   │   ├─ Create new page in Notion database
   │   ├─ Fill properties: Name, Category, Region, Difficulty, Cook Time, Servings
   │   ├─ Fill nutrition properties: Calories, Protein, Carbs, Fat
   │   ├─ Add tags via multi-select property
   │   ├─ Write ingredients as bulleted list (e.g. "200g thịt bò — thái lát mỏng")
   │   ├─ Write steps as numbered list
   │   ├─ Add images (cover + step photos)
   │   └─ Set Status property to "Published"
   ├─ Edit recipe: modify any property or content block in Notion
   ├─ Archive recipe: set Status = "Archived"
   └─ Add tips: use callout blocks
3. Changes are saved instantly in Notion (auto-save)
```

### Main Flow — Sync to Application (Backend)

```
1. Sync triggers automatically:
   ├─ On server startup
   ├─ Every N minutes (configurable via SYNC_INTERVAL_MINUTES)
   └─ Manually via POST /api/v1/admin/sync (admin auth required)
2. SyncService fetches all pages from Notion database
   ├─ Filters: Status == "Published" only
   └─ Paginates through all results
3. For each Notion page:
   ├─ Parse properties → Dish model
   ├─ Fetch child blocks → parse into Ingredients[], Steps[], Tips
   ├─ Extract nutrition data from properties
   └─ Match by external_id (Notion page ID)
4. Upsert into PostgreSQL:
   ├─ Insert new recipes
   ├─ Update changed recipes
   ├─ Soft-delete recipes removed from Notion (or set to Draft/Archived)
   └─ All within a database transaction
5. System publishes event: recipe.synced → EventBus
   └─ Cache invalidated, search index updated (tsvector trigger)
6. Sync result logged: { added: N, updated: N, deleted: N, duration: "2.3s" }
```

### Alternative Flows

- **2a.** Notion API rate limit hit → retry with exponential backoff, log warning
- **3a.** Notion block type not recognized → skip with warning log, continue sync
- **3b.** Ingredient parsing fails (can't extract amount/unit) → store raw text as ingredient name
- **4a.** Database transaction fails → rollback, log error, report partial sync

---

## UC12 — Manage Users (Admin)

**Actor:** Admin  
**Precondition:** Logged in with role = ADMIN.

### Main Flow

```
1. Admin navigates to "Manage Users"
2. System displays user list (paginated)
   ├─ Info: ID, name, email, created date, status, role
   └─ Search: by name or email
3. Admin can:
   ├─ View user details: profile, favorites count, suggestion count
   ├─ Suspend account: set status = SUSPENDED (policy violation)
   ├─ Reactivate account: set status = ACTIVE
   ├─ Change role: USER ↔ ADMIN (super admin only)
   └─ Delete account: soft delete (GDPR compliance)
4. System publishes event: UserSuspended / UserDeleted → Message Bus
   └─ Other services clean up related data
```

---

## UC13 — Manage Categories (Admin)

**Actor:** Admin  
**Precondition:** Logged in with role = ADMIN.

### Main Flow

```
1. Admin navigates to "Manage Categories"
   (Phase 1: managed via seed SQL files; Phase 2: admin UI)
2. System displays category types:
   ├─ Dish type: Braised, Stir-fried, Steamed, Grilled, Fried, Boiled...
   ├─ Region: Northern, Central, Southern
   ├─ Main ingredient: Chicken, Beef, Pork, Seafood, Vegetarian...
   ├─ Tags: Weight Loss, Easy, Quick, Party, Kid-friendly...
   └─ Units of measure: g, ml, tablespoon, cup, piece, slice...
3. Admin can:
   ├─ Add: name + slug + icon (optional) + display order
   ├─ Edit: rename, reorder
   ├─ Delete: only if no recipe is currently using it
   └─ Merge: combine 2 duplicate categories
4. System validates: no duplicate names within the same type
5. System publishes event: CategoryUpdated → EventBus
   └─ Cache invalidated

Note: Categories and tags must align with Notion database select options.
When adding a new category, also add it as a select option in Notion.
```