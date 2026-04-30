package viewmodel

import (
	"fmt"
	"time"

	"github.com/google/uuid"
        "github.com/trhys/Recipe-Repo-2/internal/database"
)

type User struct{
	ID              uuid.UUID 	`json:"id"`
        Name            string 		`json:"name"`
}

type PrivateUserViewModel struct{
	User
        Email           string 		`json:"email"`
        CreatedAt       time.Time	`json:"created_at"`
	Recipes 	[]Recipe 	`json:"recipes"`
}

type PublicUserViewModel struct{
	User
	Recipes 	[]Recipe 	`json:"recipes"`
}

type SessionViewModel struct{
	User
	Email	string	`json:"email"`
	JWT	string	`json:"token"`
	RT	string	`json:"refresh_token"`
}

func (builder *VMFactory) GeneratePrivateUser(user database.GetUserRow, recipes []database.Recipe) PrivateUserViewModel {
	model := PrivateUserViewModel{
		User: User{
			ID:	user.ID,
			Name:	user.Name,
		},
		Email:		user.Email,
		CreatedAt:	user.CreatedAt,
		Recipes:	make([]Recipe, 0, len(recipes)),
	}

	for _, r := range recipes {
		model.Recipes = append(model.Recipes, Recipe{
			ID:		r.ID,
			Title:		r.Title,
			CreatedAt: 	r.CreatedAt,
			UpdatedAt: 	&r.UpdatedAt,
			ImageURL:	fmt.Sprintf("%s/%s", builder.S3cdn, r.ImageKey),
		})
	}

	return model
}

func (builder *VMFactory) GeneratePublicUser(user database.GetUserRow, recipes []database.Recipe) PublicUserViewModel {
	model := PublicUserViewModel{
		User: User{
			ID:	user.ID,
			Name:	user.Name,
		},
		Recipes:	make([]Recipe, 0, len(recipes)),
	}

	for _, r := range recipes {
		model.Recipes = append(model.Recipes, Recipe{
			ID:		r.ID,
			Title:		r.Title,
			CreatedAt: 	r.CreatedAt,
			UpdatedAt: 	&r.UpdatedAt,
			ImageURL:	fmt.Sprintf("%s/%s", builder.S3cdn, r.ImageKey),
		})
	}

	return model
}

func GenerateSession(user database.GetUserHashRow, token, refreshToken string) SessionViewModel {
	return SessionViewModel{
		User: User{
			ID:	user.ID,
			Name:	user.Name,
		},
		Email:	user.Email,
		JWT:	token,
		RT:	refreshToken,
	}
}
