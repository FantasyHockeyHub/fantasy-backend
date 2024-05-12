package service

import (
	"context"
	"github.com/Frozen-Fantasy/fantasy-backend.git/config"
	"github.com/Frozen-Fantasy/fantasy-backend.git/pkg/models/players"
	"github.com/Frozen-Fantasy/fantasy-backend.git/pkg/models/store"
	"github.com/Frozen-Fantasy/fantasy-backend.git/pkg/models/tournaments"
	"github.com/Frozen-Fantasy/fantasy-backend.git/pkg/models/user"
	"github.com/Frozen-Fantasy/fantasy-backend.git/pkg/storage"
	"github.com/google/uuid"
)

//go:generate mockgen -source=service.go -destination=mocks/mock.go

type User interface {
	SignUp(input user.SignUpInput) error
	SignIn(input user.SignInInput) (user.Tokens, error)
	RefreshTokens(refreshTokenID string) (user.Tokens, error)
	CreateSession(userID uuid.UUID) (user.Tokens, error)
	Logout(refreshTokenID string) error
	SendVerificationCode(email string) error
	CheckEmailVerification(email string, inputCode int) error
	CheckEmailExists(email string) (bool, error)
	CheckNicknameExists(nickname string) (bool, error)
	ChangePassword(inp user.ChangePasswordModel) error
	ForgotPassword(email string) error
	ResetPassword(inp user.ResetPasswordInput) error
	GetUserInfo(userID uuid.UUID) (user.UserInfoModel, error)
	CheckUserDataExists(inp user.UserExistsDataInput) error
	DeleteProfile(userID uuid.UUID) error
	GetCoinTransactions(profileID uuid.UUID) ([]user.CoinTransactionsModel, error)
}

type TokenManager interface {
	CreateJWT(userID string) (int64, string, error)
	ParseJWT(accessToken string) (string, error)
	CreateRefreshToken() (string, error)
}

type Teams interface {
	CreateTeamsNHL(context.Context, []tournaments.Standing) error
	CreateTeamsKHL(ctx context.Context, teams []tournaments.TeamKHL) error
	GetMatchesDay(ctx context.Context, league tournaments.League) ([]tournaments.Matches, error)
}

type Tournaments interface {
	GetTournaments(context.Context, tournaments.League) ([]tournaments.Tournament, error)
	GetMatchesByTournamentsId(context.Context, tournaments.ID) ([]tournaments.GetTournamentsTotalInfo, error)
	GetRosterByTournamentID(userID uuid.UUID, tournamentID int) (players.TournamentRosterResponse, error)
	CreateTournamentTeam(inp tournaments.TournamentTeamModel) error
	CheckUserTeam(tournamentInfo tournaments.Tournament, userTeam []int) error
	GetTeamCost(team []int) (float32, error)
	GetTournamentTeam(userID uuid.UUID, tournamentID int) (players.UserTeamResponse, error)
	EditTournamentTeam(inp tournaments.TournamentTeamModel) error
}

type Store interface {
	GetAllProducts() ([]store.Product, error)
	BuyProduct(buy store.BuyProductModel) error
}

type Players interface {
	CreatePlayers(playersData []players.Player) error
	GetPlayers(playersFilter players.PlayersFilter) ([]players.PlayerResponse, error)
	GetPlayerCards(filter players.PlayerCardsFilter) ([]players.PlayerCardResponse, error)
	CardUnpacking(id int, userID uuid.UUID) error
}

type Services struct {
	User
	TokenManager
	Teams
	Tournaments
	Store
	Players
}

type Deps struct {
	Cfg      config.ServiceConfiguration
	Storage  *storage.PostgresStorage
	RStorage *storage.RedisStorage
	Jwt      *Manager
}

func NewServices(deps Deps) *Services {
	userService := NewUserService(deps.Storage, deps.RStorage, deps.Jwt, deps.Cfg)
	playersService := NewPlayersService(deps.Storage)
	tournamentsService := NewTournamentsService(deps.Storage, playersService)
	storeService := NewStoreService(deps.Storage)
	teamsService := NewTeamsService(deps.Storage)
	return &Services{
		User:         userService,
		TokenManager: deps.Jwt,
		Teams:        teamsService,
		Tournaments:  tournamentsService,
		Store:        storeService,
		Players:      playersService,
	}
}
