package event

const (
	RecipeSynced        = "recipe.synced"
	RecipeDeleted       = "recipe.deleted"
	UserProfileUpdated  = "user.profile_updated"
	UserDeleted         = "user.deleted"
	EngagementFavorited = "engagement.favorited"
	EngagementViewed    = "engagement.viewed"
)

type Event struct {
	Type    string
	Payload any
}

type RecipeSyncedPayload struct {
	Added   int
	Updated int
	Deleted int
}

type RecipeDeletedPayload struct {
	RecipeID string
}
