package google

type LocationHistory struct {
	Filepath string
}

func (p *LocationHistory) String() string {
	return p.Filepath
}

func (p *LocationHistory) Parse() (*models.Data, error) {
	var data models.GoogleData

	file, err := Load(p.Filepath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	err = decoder.Decode(&data)
	if err != nil {
		return nil, fmt.Errorf("failed to decode json file '%s': %w", p.Filepath, err)
	}

	return conversions.GoogleDataToData(&data)
}