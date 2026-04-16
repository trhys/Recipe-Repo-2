package main

import (
	"encoding/json"
	"fmt"
	"io"
	"mime"
	"net/http"
	"os"
	"time"

	"github.com/google/uuid"
        "github.com/aws/aws-sdk-go-v2/service/s3"
        "github.com/trhys/Recipe-Repo-2/internal/database"
        "github.com/trhys/Recipe-Repo-2/internal/auth"
)

type ingredient struct {
        ID              uuid.UUID `json:"id"`
        Name            string `json:"name"`
	ImageKey	string `json:"image_key"`
        CreatedAt       time.Time `json:"created_at"`
        UpdatedAt       time.Time `json:"updated_at"`
}

func (cfg *apiConfig) handlerCreateIngredient(w http.ResponseWriter, r *http.Request) {
	// Authorization
        token, err := auth.GetBearerToken(r.Header)
        if err != nil {
                respondFail(w, 401, "Failed to retrieve bearer token", err)
                return
        }

        subject, err := auth.ValidateJWT(token, cfg.secret)
        if err != nil {
                respondFail(w, 401, "Couldn't validate token", err)
                return
        }

	admin, err := cfg.db.CheckAdmin(r.Context(), subject)
	if err != nil {
		respondFail(w, 500, "Something went wrong", err)
		return
	} else if admin == false {
		respondFail(w, 403, "Unauthorized access", fmt.Errorf("Must be administrator"))
		return
	}

	// User is admin --- proceed
        
	var req struct{
		Name string `json:"name"`
	}

	// Get request payload 
        jsonString := r.FormValue("payload")

        // Unmarshal JSON
        if err := json.Unmarshal([]byte(jsonString), &req); err != nil {
                respondFail(w, 500, "Failed to unmarshal payload", err)
                return
        }

	file, fileHeader, err := r.FormFile("image")
        key := uuid.New().String()
        if err == nil {
                defer file.Close()

                mediaType, _, err := mime.ParseMediaType(fileHeader.Header.Get("Content-Type"))
                if err != nil {
                        respondFail(w, 400, "Couldn't parse media type", err)
                        return
                }

                if mediaType != "image/jpeg" && mediaType != "image/png" {
                        respondFail(w, 400, "Invalid media type", fmt.Errorf("Must be jpg or png. Got: %s", mediaType))
                        return
                }

                tmp, err := os.CreateTemp("", "image_upload")
                if err != nil {
                        respondFail(w, 500, "Something went wrong", err)
                        return
                }
                defer os.Remove(tmp.Name())
                defer tmp.Close()

                _, fail := io.Copy(tmp, file)
                if fail != nil {
                        respondFail(w, 500, "Couldn't save image", err)
                        return
                }

                tmp.Seek(0, io.SeekStart)

                // Upload to s3
                if _, err := cfg.s3client.PutObject(r.Context(), &s3.PutObjectInput{
                        Bucket: &cfg.s3bucket,
                        Key: &key,
                        Body: tmp,
                        ContentType: &mediaType,
                }); err != nil {
                        respondFail(w, 500, "Failed to upload to s3 bucket", err)
                        return
                }

        } else if err != nil {
                if err != http.ErrMissingFile {
                        respondFail(w, 400, "Something went wrong during upload", err)
                        return
                }
        }

	query := database.CreateIngredientParams{
		Name: req.Name,
		ImageKey: key,
	}

	ing, err := cfg.db.CreateIngredient(r.Context(), query)
	if err != nil {
		respondFail(w, 500, "Database error during write", err)
		return
	}

	res := ingredient{
		ID: ing.ID,
		Name: ing.Name,
		ImageKey: ing.ImageKey,
		CreatedAt: ing.CreatedAt,
		UpdatedAt: ing.UpdatedAt,
	}

	respondJSON(w, 200, res)
}
