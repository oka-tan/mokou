//Package config wraps Mokou's configuration
package config

import (
	"log"
	"os"

	"github.com/BurntSushi/toml"
)

//Config is the Mokou configuration struct
type Config struct {
	PostgresConfig   PostgresConfig `toml:"postgres"`
	MediaConfig      *S3Config      `toml:"media"`
	ThumbnailsConfig *S3Config      `toml:"thumbnails"`
	OekakiConfig     *S3Config      `toml:"oekaki"`
	AsagiConfig      *AsagiConfig   `toml:"asagi"`
	BadgerConfig     *BadgerConfig  `toml:"badger"`
}

//PostgresConfig is the configuration struct for
//connection to a Koiwai db instance.
type PostgresConfig struct {
	ConnectionString string `toml:"connection_string"`
	BatchSize        int    `toml:"batch_size"`
}

//S3Config is the configuration struct for
//connection to S3 storage.
type S3Config struct {
	S3Endpoint        string `toml:"s3_endpoint"`
	S3AccessKeyID     string `toml:"s3_access_key_id"`
	S3SecretAccessKey string `toml:"s3_secret_access_key"`
	S3UseSSL          bool   `toml:"s3_use_ssl"`
	S3Region          string `toml:"s3_region"`
	S3BucketName      string `toml:"s3_bucket_name"`
}

//AsagiConfig is the configuration struct for
//connection to an AsagiConfig db instance.
type AsagiConfig struct {
	ConnectionString string             `toml:"connection_string"`
	ImagesFolder     *string            `toml:"images_folder"`
	Boards           []AsagiBoardConfig `toml:"boards"`
}

//BadgerConfig configures the import from badger's dumps
type BadgerConfig struct {
	JsonFolder string              `toml:"json_folder"`
	Boards     []BadgerBoardConfig `toml:"boards"`
}

//AsagiBoardConfig is the configuration struct
//an Asagi board to be imported
type AsagiBoardConfig struct {
	Name             string `toml:"name"`
	ImportImages     bool   `toml:"import_images"`
	ImportThumbnails bool   `toml:"import_thumbnails"`
	EnableCode       bool   `toml:"enable_code"`
	EnableSpoiler    bool   `toml:"enable_spoiler"`
	EnableFortune    bool   `toml:"enable_fortune"`
	EnableExif       bool   `toml:"enable_exif"`
	EnableOekaki     bool   `toml:"enable_oekaki"`
}

//BadgerBoardConfig is the configuration struct for
//a badger board to be imported
type BadgerBoardConfig struct {
	Name string `toml:"name"`
}

//LoadConfig reads config.json and unmarshals it into a Config struct.
func LoadConfig() Config {
	configFile := os.Getenv("MOKOU_CONFIG")

	if configFile == "" {
		configFile = "./config.json"
	}

	f, err := os.Open(configFile)

	if err != nil {
		log.Fatalf("Error opening configuration file: %s", err)
	}

	var conf Config

	if _, err := toml.NewDecoder(f).Decode(&conf); err != nil {
		log.Fatalf("Error unmarshalling configuration file: %s", err)
	}

	return conf
}
