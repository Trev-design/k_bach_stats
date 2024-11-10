package db

import (
	"user_manager/graph/model"
)

type XPData struct {
	ID   string
	Name string
}

func getIDsFromExistingExperiences(experiences []*model.ExperienceCredentials) []string {
	ids := make([]string, len(experiences))
	for _, experience := range experiences {
		ids = append(ids, experience.ExperienceID)
	}

	return ids
}

func getaRtingsFromNewExperiences(experiences []*model.NewExperienceCredentials) []int {
	ratings := make([]int, len(experiences))
	for _, experience := range experiences {
		ratings = append(ratings, experience.Rating)
	}

	return ratings
}

func getRatingsFromExistingExperiences(experiences []*model.ExperienceCredentials) []int {
	ratings := make([]int, len(experiences))
	for _, experience := range experiences {
		ratings = append(ratings, experience.Rating)
	}

	return ratings
}

func getProfileIDs(experiences []*model.NewExperienceCredentials) []string {
	ids := make([]string, len(experiences))

	for _, experience := range experiences {
		ids = append(ids, experience.ProfileID)
	}

	return ids
}

func getProfileIDsFromExisting(experiences []*model.ExperienceCredentials) []string {
	ids := make([]string, len(experiences))

	for _, experience := range experiences {
		ids = append(ids, experience.ProfileID)
	}

	return ids
}

func getExperienceNames(experiences []*model.NewExperienceCredentials) []string {
	names := make([]string, len(experiences))

	for _, experience := range experiences {
		names = append(names, experience.Experience)
	}

	return names
}

func getExperienceIDs(experiences []*model.ExperienceCredentials) []string {
	ids := make([]string, len(experiences))

	for _, experience := range experiences {
		ids = append(ids, experience.ExperienceID)
	}

	return ids
}

func getExperiencesFromExisting(existing []*model.ExperienceCredentials, ratingIDs []string, names []string) []*model.Experience {
	results := make([]*model.Experience, len(existing))

	for index, experience := range existing {
		results = append(
			results,
			&model.Experience{
				ID:         ratingIDs[index],
				Experience: names[index],
				Rating:     experience.Rating,
			},
		)
	}

	return results
}

func getExperienceProfileJoinData(experiencIDs []string, profileIDs []string) []any {
	joinData := make([]any, len(experiencIDs))

	for index, experienceID := range experiencIDs {
		joinData = append(
			joinData,
			experienceID,
			profileIDs[index],
		)
	}

	return joinData
}
