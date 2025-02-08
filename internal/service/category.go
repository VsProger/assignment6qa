package service

func (s *service) GetAllCategory() ([]string, error) {

	categories, err := s.repo.GetALLCategory()
	if err != nil {
		return nil, err
	}
	return categories, nil
}
