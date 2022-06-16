package main

import "forumynov.com/db"

type Post struct {
	Id       int
	Username string
	Date     string
}

type Comment struct {
	Id       int
	Username string
	Date     string
	PostID   int
}

func Posts() (posts []Post, err error) {
	rows, err := db.DB.Query("SELECT id, username, date FROM post ORDER BY date DESC")
	if err != nil {
		return
	}
	for rows.Next() {
		th := Post{}
		if err = rows.Scan(&th.Id, &th.Username, &th.Date); err != nil {

			return
		}
		posts = append(posts, th)
	}
	rows.Close()
	return
}

func Comments() (comments []Comment, err error) {
	rows, err := db.DB.Query("SELECT id, username, date, postID FROM comment ORDER BY date DESC")
	if err != nil {
		return
	}
	for rows.Next() {
		th := Comment{}
		if err = rows.Scan(&th.Id, &th.Username, &th.Date, &th.PostID); err != nil {

			return
		}
		comments = append(comments, th)
	}
	rows.Close()
	return
}
