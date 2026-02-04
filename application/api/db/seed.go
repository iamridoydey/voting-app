package db

import (
    "context"
    "time"

    "go.mongodb.org/mongo-driver/bson"

    "voting-app/api/models"
)

func SeedLanguages() error {
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    count, _ := Collection.CountDocuments(ctx, bson.M{})
    if count == 0 {
        languages := []interface{}{
            models.Language{Name: "Rust", Votes: 0, Image: "https://www.rust-lang.org/logos/rust-logo-512x512.png"},
            models.Language{Name: "Go", Votes: 0, Image: "https://raw.githubusercontent.com/golang-samples/gopher-vector/master/gopher.png"},
            models.Language{Name: "Python", Votes: 0, Image: "https://upload.wikimedia.org/wikipedia/commons/c/c3/Python-logo-notext.svg"},
            models.Language{Name: "JavaScript", Votes: 0, Image: "https://upload.wikimedia.org/wikipedia/commons/6/6a/JavaScript-logo.png"},
            models.Language{Name: "Java", Votes: 0, Image: "https://upload.wikimedia.org/wikipedia/en/3/30/Java_programming_language_logo.svg"},
            models.Language{Name: "C++", Votes: 0, Image: "https://upload.wikimedia.org/wikipedia/commons/1/18/ISO_C%2B%2B_Logo.svg"},
            models.Language{Name: "C#", Votes: 0, Image: "https://upload.wikimedia.org/wikipedia/commons/4/4f/Csharp_Logo.png"},
            models.Language{Name: "PHP", Votes: 0, Image: "https://www.php.net/images/logos/php-logo.svg"},
            models.Language{Name: "TypeScript", Votes: 0, Image: "https://raw.githubusercontent.com/remojansen/logo.ts/master/ts.png"},
        }
        
        _, err := Collection.InsertMany(ctx, languages)
        return err
    }
    return nil
}