package models

type Artist struct {
    UserID    string `json:"user_id"`
    Username  string `json:"username"`
    Email     string `json:"email"`
    Role      string `json:"role"`
    Followers int    `json:"followers"`
}

type ArtistWithSongs struct {
    Artist Artist `json:"artist"`
    Songs  []Song `json:"songs"`
} 