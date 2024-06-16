package tests

import (
	"context"
	"github.com/Mx1q/ppo_repoMysql/repository/mySQL"
	"github.com/Mx1q/ppo_services/domain"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"testing"
)

func Test_saladRepository_GetAll(t *testing.T) {
	repo := mysql.NewSaladRepository(testDbInstance)

	tests := []struct {
		name       string
		filter     *domain.RecipeFilter
		beforeTest func() (*domain.RecipeFilter, []*domain.Salad)
		page       int
		wantErr    bool
		errStr     error
	}{
		{
			name: "успешное получение, фильтр по типам и ингредиентам",
			page: 1,
			beforeTest: func() (*domain.RecipeFilter, []*domain.Salad) {
				filter := new(domain.RecipeFilter)

				filter.AvailableIngredients = make([]uuid.UUID, 2)
				filter.AvailableIngredients[0], _ = uuid.Parse("f1fc4bfc-799c-4471-a971-1bb00f7dd30a")
				filter.AvailableIngredients[1], _ = uuid.Parse("01000000-0000-0000-0000-000000000000")
				filter.Status = domain.PublishedSaladStatus
				filter.SaladTypes = make([]uuid.UUID, 1)
				filter.SaladTypes[0], _ = uuid.Parse("7e17866b-2b97-4d2b-b399-42ceeebd5480")

				saladId, _ := uuid.Parse("fbabc2aa-cd4a-42b0-b68d-d3cf67fba06f")
				saladId2, _ := uuid.Parse("01000000-0000-0000-0000-000000000000")

				expected := make([]*domain.Salad, 2)
				expected[0] = &domain.Salad{
					ID:          saladId,
					AuthorID:    uuid.Nil,
					Name:        "цезарь",
					Description: "",
				}
				expected[1] = &domain.Salad{
					ID:          saladId2,
					AuthorID:    uuid.Nil,
					Name:        "овощной",
					Description: "",
				}
				return filter, expected
			},
			wantErr: false,
		}, // успешное получение.
		{
			name: "успешное получение, фильтр только по ингредиентам",
			page: 1,
			beforeTest: func() (*domain.RecipeFilter, []*domain.Salad) {
				filter := new(domain.RecipeFilter)

				filter.AvailableIngredients = make([]uuid.UUID, 2)
				filter.AvailableIngredients[0], _ = uuid.Parse("f1fc4bfc-799c-4471-a971-1bb00f7dd30a")
				filter.AvailableIngredients[1], _ = uuid.Parse("02000000-0000-0000-0000-000000000000")
				filter.Status = domain.PublishedSaladStatus

				saladId, _ := uuid.Parse("fbabc2aa-cd4a-42b0-b68d-d3cf67fba06f")
				saladId2, _ := uuid.Parse("03000000-0000-0000-0000-000000000000")

				expected := make([]*domain.Salad, 2)
				expected[0] = &domain.Salad{
					ID:          saladId,
					AuthorID:    uuid.Nil,
					Name:        "цезарь",
					Description: "",
				}
				expected[1] = &domain.Salad{
					ID:          saladId2,
					AuthorID:    uuid.Nil,
					Name:        "сельдь под шубой",
					Description: "",
				}
				return filter, expected
			},
			wantErr: false,
		}, // успешное получение. фильтр только по ингредиентам
		{
			name: "успешное получение, фильтр только по типам",
			page: 1,
			beforeTest: func() (*domain.RecipeFilter, []*domain.Salad) {
				filter := new(domain.RecipeFilter)

				filter.Status = domain.PublishedSaladStatus
				filter.SaladTypes = make([]uuid.UUID, 1)
				filter.SaladTypes[0], _ = uuid.Parse("01000000-0000-0000-0000-000000000000")

				saladId, _ := uuid.Parse("02000000-0000-0000-0000-000000000000")
				saladId2, _ := uuid.Parse("01000000-0000-0000-0000-000000000000")

				expected := make([]*domain.Salad, 2)
				expected[1] = &domain.Salad{
					ID:          saladId,
					AuthorID:    uuid.Nil,
					Name:        "сезонный",
					Description: "",
				}
				expected[0] = &domain.Salad{
					ID:          saladId2,
					AuthorID:    uuid.Nil,
					Name:        "овощной",
					Description: "",
				}
				return filter, expected
			},
			wantErr: false,
		}, // успешное получение. фильтр только по типам
		{
			name: "успешное получение, пустой фильтр",
			page: 1,
			beforeTest: func() (*domain.RecipeFilter, []*domain.Salad) {
				filter := new(domain.RecipeFilter)

				filter.Status = domain.PublishedSaladStatus
				saladId, _ := uuid.Parse("fbabc2aa-cd4a-42b0-b68d-d3cf67fba06f")
				saladId2, _ := uuid.Parse("01000000-0000-0000-0000-000000000000")
				saladId3, _ := uuid.Parse("02000000-0000-0000-0000-000000000000")
				saladId4, _ := uuid.Parse("03000000-0000-0000-0000-000000000000")
				saladId5, _ := uuid.Parse("04000000-0000-0000-0000-000000000000")

				expected := make([]*domain.Salad, 5)
				expected[0] = &domain.Salad{
					ID:          saladId,
					AuthorID:    uuid.Nil,
					Name:        "цезарь",
					Description: "",
				}
				expected[1] = &domain.Salad{
					ID:          saladId2,
					AuthorID:    uuid.Nil,
					Name:        "овощной",
					Description: "",
				}
				expected[2] = &domain.Salad{
					ID:          saladId4,
					AuthorID:    uuid.Nil,
					Name:        "сельдь под шубой",
					Description: "",
				}
				expected[3] = &domain.Salad{
					ID:          saladId3,
					AuthorID:    uuid.Nil,
					Name:        "сезонный",
					Description: "",
				}
				expected[4] = &domain.Salad{
					ID:          saladId5,
					AuthorID:    uuid.Nil,
					Name:        "греческий",
					Description: "",
				}
				return filter, expected
			},
			wantErr: false,
		}, // успешное получение. пустой фильтр
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			filter := new(domain.RecipeFilter)
			expected := make([]*domain.Salad, 1)
			if tt.beforeTest != nil {
				filter, expected = tt.beforeTest()
			}
			salads, _, err := repo.GetAll(context.Background(), filter, tt.page)

			if tt.wantErr {
				require.Equal(t, tt.errStr.Error(), err.Error())
			} else {
				require.Nil(t, err)
				require.Equal(t, expected, salads)
			}
		})
	}
}
