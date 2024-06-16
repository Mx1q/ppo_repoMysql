package tests

import (
	"context"
	"errors"
	"github.com/Mx1q/ppo_repoMysql/repository/mySQL"
	"github.com/Mx1q/ppo_services/domain"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"testing"
)

func Test_recipeStepRepository_Create(t *testing.T) {
	repo := mysql.NewRecipeStepRepository(testDbInstance)

	tests := []struct {
		name       string
		recipeStep *domain.RecipeStep
		wantErr    bool
		errStr     error
	}{
		{
			name: "успешное создание",
			recipeStep: &domain.RecipeStep{
				RecipeID:    uuid.UUID{3},
				Name:        "first",
				Description: "description",
			},
			wantErr: false,
		}, // успешное создание
		{
			name: "несуществующий рецепт",
			recipeStep: &domain.RecipeStep{
				RecipeID:    uuid.UUID{111},
				Name:        "first",
				Description: "description",
			},
			wantErr: true,
			errStr:  errors.New("creating recipe step: Error 1452 (23000): Cannot add or update a child row: a foreign key constraint fails (`saladRecipes`.`recipeStep`, CONSTRAINT `recipeStep_ibfk_1` FOREIGN KEY (`recipeId`) REFERENCES `recipe` (`id`) ON DELETE CASCADE)"),
		}, // несуществующий рецепт
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := repo.Create(context.Background(), tt.recipeStep)

			if tt.wantErr {
				require.Equal(t, tt.errStr.Error(), err.Error())
			} else {
				require.Nil(t, err)
			}
		})
	}
}

func Test_recipeStepRepository_Update(t *testing.T) {
	repo := mysql.NewRecipeStepRepository(testDbInstance)
	recipeId := uuid.UUID{1}

	tests := []struct {
		name       string
		recipeStep *domain.RecipeStep
		expected   []*domain.RecipeStep
		wantErr    bool
		errStr     error
	}{
		{
			name: "Успешное обновление",
			recipeStep: &domain.RecipeStep{
				ID:          uuid.UUID{1},
				RecipeID:    uuid.UUID{2},
				Name:        "first",
				Description: "first",
				StepNum:     1,
			},
			expected: []*domain.RecipeStep{
				{
					ID:          uuid.UUID{1},
					RecipeID:    uuid.UUID{2},
					Name:        "first",
					Description: "first",
					StepNum:     1,
				},
			},
			wantErr: false,
		}, // успешное обновление
		{
			name: "Номер шага больше максимального, который существует",
			recipeStep: &domain.RecipeStep{
				ID:          uuid.UUID{1},
				RecipeID:    uuid.UUID{2},
				Name:        "first",
				Description: "first",
				StepNum:     2,
			},
			wantErr: true,
			errStr:  errors.New("updating recipe step: step num out of range"),
		}, // Номер шага больше максимального, который существует
		{
			name: "перемещение первого в конец",
			recipeStep: &domain.RecipeStep{
				ID:          uuid.UUID{2},
				RecipeID:    recipeId,
				Name:        "first",
				Description: "first",
				StepNum:     5,
			},
			expected: []*domain.RecipeStep{
				{
					ID:          uuid.UUID{3},
					Name:        "second",
					Description: "second",
					RecipeID:    recipeId,
					StepNum:     1,
				},
				{
					ID:          uuid.UUID{4},
					Name:        "third",
					Description: "third",
					RecipeID:    recipeId,
					StepNum:     2,
				},
				{
					ID:          uuid.UUID{5},
					Name:        "fourth",
					Description: "fourth",
					RecipeID:    recipeId,
					StepNum:     3,
				},
				{
					ID:          uuid.UUID{6},
					Name:        "fifth",
					Description: "fifth",
					RecipeID:    recipeId,
					StepNum:     4,
				},
				{
					ID:          uuid.UUID{2},
					RecipeID:    recipeId,
					Name:        "first",
					Description: "first",
					StepNum:     5,
				},
			},
			wantErr: false,
		}, // перемещение первого в конец
		{
			name: "перемещение последнего в начало",
			recipeStep: &domain.RecipeStep{
				ID:          uuid.UUID{2},
				RecipeID:    recipeId,
				Name:        "first",
				Description: "first",
				StepNum:     1,
			},
			expected: []*domain.RecipeStep{
				{
					ID:          uuid.UUID{2},
					RecipeID:    recipeId,
					Name:        "first",
					Description: "first",
					StepNum:     1,
				},
				{
					ID:          uuid.UUID{3},
					Name:        "second",
					Description: "second",
					RecipeID:    recipeId,
					StepNum:     2,
				},
				{
					ID:          uuid.UUID{4},
					Name:        "third",
					Description: "third",
					RecipeID:    recipeId,
					StepNum:     3,
				},
				{
					ID:          uuid.UUID{5},
					Name:        "fourth",
					Description: "fourth",
					RecipeID:    recipeId,
					StepNum:     4,
				},
				{
					ID:          uuid.UUID{6},
					Name:        "fifth",
					Description: "fifth",
					RecipeID:    recipeId,
					StepNum:     5,
				},
			},
			wantErr: false,
		}, // перемещение последнего в начало
		{
			name: "перемещение из середины вперед",
			recipeStep: &domain.RecipeStep{
				ID:          uuid.UUID{3},
				RecipeID:    recipeId,
				Name:        "second",
				Description: "second",
				StepNum:     4,
			},
			expected: []*domain.RecipeStep{
				{
					ID:          uuid.UUID{2},
					RecipeID:    recipeId,
					Name:        "first",
					Description: "first",
					StepNum:     1,
				},
				{
					ID:          uuid.UUID{4},
					Name:        "third",
					Description: "third",
					RecipeID:    recipeId,
					StepNum:     2,
				},
				{
					ID:          uuid.UUID{5},
					Name:        "fourth",
					Description: "fourth",
					RecipeID:    recipeId,
					StepNum:     3,
				},
				{
					ID:          uuid.UUID{3},
					Name:        "second",
					Description: "second",
					RecipeID:    recipeId,
					StepNum:     4,
				},
				{
					ID:          uuid.UUID{6},
					Name:        "fifth",
					Description: "fifth",
					RecipeID:    recipeId,
					StepNum:     5,
				},
			},
			wantErr: false,
		}, // перемещение из середины вперед
		{
			name: "перемещение из середины назад",
			recipeStep: &domain.RecipeStep{
				ID:          uuid.UUID{3},
				RecipeID:    recipeId,
				Name:        "second",
				Description: "second",
				StepNum:     2,
			},
			expected: []*domain.RecipeStep{
				{
					ID:          uuid.UUID{2},
					RecipeID:    recipeId,
					Name:        "first",
					Description: "first",
					StepNum:     1,
				},
				{
					ID:          uuid.UUID{3},
					Name:        "second",
					Description: "second",
					RecipeID:    recipeId,
					StepNum:     2,
				},
				{
					ID:          uuid.UUID{4},
					Name:        "third",
					Description: "third",
					RecipeID:    recipeId,
					StepNum:     3,
				},
				{
					ID:          uuid.UUID{5},
					Name:        "fourth",
					Description: "fourth",
					RecipeID:    recipeId,
					StepNum:     4,
				},
				{
					ID:          uuid.UUID{6},
					Name:        "fifth",
					Description: "fifth",
					RecipeID:    recipeId,
					StepNum:     5,
				},
			},
			wantErr: false,
		}, // перемещение из середины назад
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := repo.Update(context.Background(), tt.recipeStep)

			if tt.wantErr {
				require.Equal(t, tt.errStr.Error(), err.Error())
			} else {
				require.Nil(t, err)
				steps, err := repo.GetAllByRecipeID(context.Background(), tt.recipeStep.RecipeID)
				require.Nil(t, err)
				require.Equal(t, tt.expected, steps)
			}
		})
	}
}

func Test_recipeStepRepository_GetAllByRecipeID(t *testing.T) {
	repo := mysql.NewRecipeStepRepository(testDbInstance)

	tests := []struct {
		name     string
		recipeID uuid.UUID
		wantErr  bool
		expected []*domain.RecipeStep
		errStr   error
	}{
		{
			name:     "успешное получение",
			recipeID: uuid.UUID{4},
			wantErr:  false,
			expected: []*domain.RecipeStep{
				{
					ID:          uuid.UUID{8},
					Name:        "first",
					Description: "first",
					RecipeID:    uuid.UUID{4},
					StepNum:     1,
				},
				{
					ID:          uuid.UUID{9},
					Name:        "second",
					Description: "second",
					RecipeID:    uuid.UUID{4},
					StepNum:     2,
				},
				{
					ID:          uuid.UUID{10},
					Name:        "third",
					Description: "third",
					RecipeID:    uuid.UUID{4},
					StepNum:     3,
				},
			},
		}, // успешное получение
		{
			name:     "отсутствие id рецепта",
			recipeID: uuid.UUID{9},
			expected: []*domain.RecipeStep{},
			wantErr:  false,
		}, // отсутствие id рецепта
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			steps, err := repo.GetAllByRecipeID(context.Background(), tt.recipeID)

			if tt.wantErr {
				require.Equal(t, tt.errStr.Error(), err.Error())
			} else {
				require.Nil(t, err)
				require.Equal(t, tt.expected, steps)
			}
		})
	}
}

func Test_recipeStepRepository_GetById(t *testing.T) {
	repo := mysql.NewRecipeStepRepository(testDbInstance)

	tests := []struct {
		name     string
		id       uuid.UUID
		wantErr  bool
		expected *domain.RecipeStep
		errStr   error
	}{
		{
			name:    "успешное получение",
			id:      uuid.UUID{8},
			wantErr: false,
			expected: &domain.RecipeStep{
				ID:          uuid.UUID{8},
				Name:        "first",
				Description: "first",
				RecipeID:    uuid.UUID{4},
				StepNum:     1,
			},
		}, // успешное получение
		{
			name:    "отсутствие id",
			id:      uuid.Nil,
			wantErr: true,
			errStr:  errors.New("getting recipe step by id: record not found"),
		}, // отсутствие id
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			step, err := repo.GetById(context.Background(), tt.id)

			if tt.wantErr {
				require.Equal(t, tt.errStr.Error(), err.Error())
			} else {
				require.Nil(t, err)
				require.Equal(t, tt.expected, step)
			}
		})
	}
}

func Test_recipeStepRepository_DeleteAllByRecipeID(t *testing.T) {
	repo := mysql.NewRecipeStepRepository(testDbInstance)

	tests := []struct {
		name     string
		recipeId uuid.UUID
		wantErr  bool
		errStr   error
	}{
		{
			name:     "успешное удаление",
			recipeId: uuid.UUID{2},
			wantErr:  false,
		}, // успешное удаление
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := repo.DeleteAllByRecipeID(context.Background(), tt.recipeId)

			if tt.wantErr {
				require.Equal(t, tt.errStr.Error(), err.Error())
			} else {
				require.Nil(t, err)
			}
		})
	}
}

func Test_recipeStepRepository_DeleteById(t *testing.T) {
	repo := mysql.NewRecipeStepRepository(testDbInstance)
	recipeId := uuid.UUID{1}

	tests := []struct {
		name     string
		id       uuid.UUID
		recipeId uuid.UUID // to check changing nums of other recipe steps
		expected []*domain.RecipeStep
		wantErr  bool
		errStr   error
	}{
		{
			name:     "успешное удаление",
			id:       uuid.UUID{2},
			recipeId: recipeId,
			expected: []*domain.RecipeStep{
				{
					ID:          uuid.UUID{3},
					Name:        "second",
					Description: "second",
					RecipeID:    recipeId,
					StepNum:     1,
				},
				{
					ID:          uuid.UUID{4},
					Name:        "third",
					Description: "third",
					RecipeID:    recipeId,
					StepNum:     2,
				},
				{
					ID:          uuid.UUID{5},
					Name:        "fourth",
					Description: "fourth",
					RecipeID:    recipeId,
					StepNum:     3,
				},
				{
					ID:          uuid.UUID{6},
					Name:        "fifth",
					Description: "fifth",
					RecipeID:    recipeId,
					StepNum:     4,
				},
			},
			wantErr: false,
		}, // успешное удаление перового шага
		{
			name:     "успешное удаление шага из середины рецепта",
			id:       uuid.UUID{4},
			recipeId: recipeId,
			expected: []*domain.RecipeStep{
				{
					ID:          uuid.UUID{3},
					Name:        "second",
					Description: "second",
					RecipeID:    recipeId,
					StepNum:     1,
				},
				{
					ID:          uuid.UUID{5},
					Name:        "fourth",
					Description: "fourth",
					RecipeID:    recipeId,
					StepNum:     2,
				},
				{
					ID:          uuid.UUID{6},
					Name:        "fifth",
					Description: "fifth",
					RecipeID:    recipeId,
					StepNum:     3,
				},
			},
			wantErr: false,
		}, // успешное удаление шага из середины рецепта
		{
			name:     "успешное удаление последнего шага",
			id:       uuid.UUID{6},
			recipeId: recipeId,
			expected: []*domain.RecipeStep{
				{
					ID:          uuid.UUID{3},
					Name:        "second",
					Description: "second",
					RecipeID:    recipeId,
					StepNum:     1,
				},
				{
					ID:          uuid.UUID{5},
					Name:        "fourth",
					Description: "fourth",
					RecipeID:    recipeId,
					StepNum:     2,
				},
			},
			wantErr: false,
		}, // успешное удаление последнего шага
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := repo.DeleteById(context.Background(), tt.id)

			if tt.wantErr {
				require.Equal(t, tt.errStr.Error(), err.Error())
			} else {
				require.Nil(t, err)
				steps, err := repo.GetAllByRecipeID(context.Background(), tt.recipeId)
				require.Nil(t, err)
				require.Equal(t, tt.expected, steps)
			}
		})
	}
}
