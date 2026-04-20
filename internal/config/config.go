package config

import (
	"time"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/trhys/Recipe-Repo-2/internal/database"
)

type ApiConfig struct {
	DB		*database.Queries
	Platform	string
	Secret		string
	JwtDuration	time.Duration
	S3client	*s3.Client
	S3bucket	string
	S3region	string
	S3cdn		string
	ImagePlaceholder string
}
