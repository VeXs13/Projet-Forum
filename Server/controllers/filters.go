package controllers

import (
	models "../models"
)

func FilterPostsByTopic(Topic string, Posts []models.Post) []models.Post {
	var PostsFiltered []models.Post

	for _, Post := range Posts {
		for _, CurrTopic := range Post.Tags {
			if CurrTopic.Tag == Topic {
				PostsFiltered = append(PostsFiltered, Post)
			}
		}
	}
	return PostsFiltered
}

func FilterPostsByUser(UserSearched string, Posts []models.Post) []models.Post {
	var PostsFiltered []models.Post

	for _, Post := range Posts {
		if Post.Autor == UserSearched {
			PostsFiltered = append(PostsFiltered, Post)
		}
	}
	return PostsFiltered
}

// Tri par séléction
// ! Améliorer l'algorithme de tri (selection ==> Fusion)
func SortUserByLike(Posts []models.Post) []models.Post {
	var PostsSorted []models.Post

	for i := 0; i < len(Posts); i++ {
		for j := i; j < len(Posts); j++ {
			if Posts[i].NbrLikes < Posts[j].NbrLikes {
				Posts[i], Posts[j] = Posts[j], Posts[i]
			}
		}
	}
	return PostsSorted
}

func SortUserByDate() {}
