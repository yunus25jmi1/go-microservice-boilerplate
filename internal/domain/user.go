package domain

type User struct {
    ID   string
    Name string
}

type UserRepository interface {
    GetByID(ctx context.Context, id string) (*User, error)
    Create(ctx context.Context, user *User) error
    // Add more methods as needed
}

type UserService struct {
    repo UserRepository
}

func NewUserService(repo UserRepository) UserService {
    return UserService{repo: repo}
}

func (s UserService) GetUser(ctx context.Context, id string) (*User, error) {
    return s.repo.GetByID(ctx, id)
}

func (s UserService) CreateUser(ctx context.Context, user *User) error {
    return s.repo.Create(ctx, user)
}
