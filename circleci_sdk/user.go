package circleci

import "context"

type UserService service

func (s *UserService) GetUserInfo(ctx context.Context) (*UserInfo, error) {
	path := "/me"
	result := &UserInfo{}
	if err := s.client.http.Get(ctx, path, result); err != nil {
		return nil, err
	}
	return result, nil
}

func (s *UserService) GetCollaborations(ctx context.Context) ([]CollaborationInfo, error) {
	path := "/me/collaborations"
	result := []CollaborationInfo{}
	if err := s.client.http.Get(ctx, path, &result); err != nil {
		return nil, err
	}
	return result, nil
}

func (s *UserService) GetUserInfoById(ctx context.Context, id string) (*UserInfo, error) {
	path := "/user/" + id
	result := &UserInfo{}
	if err := s.client.http.Get(ctx, path, result); err != nil {
		return nil, err
	}
	return result, nil
}

type CollaborationInfo struct {
	ID        string `json:"id"`
	VcsType   string `json:"vcs-type"`
	Name      string `json:"name"`
	AvatarURL string `json:"avatar_url"`
	Slug      string `json:"slug"`
}

type UserInfo struct {
	ID      string  `json:"id"`
	Login   string  `json:"login"`
	Name    string  `json:"name"`
	Message *string `json:"message,omitempty"`
}
