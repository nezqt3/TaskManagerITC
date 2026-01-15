package database
import "time"

type Task struct {
	ID    			int64
	Description 	string 
	Deadline 		time.Time 
	Status 			string 
	User 			string 
	Title 			string
	Author 			string
	IdProject 		string
}