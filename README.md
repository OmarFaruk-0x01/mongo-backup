# Mongo DB Backup Script

This is a GoLang script for backing up MongoDB databases and pushing the backup files to AWS S3. It provides a command-line interface for specifying the MongoDB URI, database name, AWS S3 bucket details, and other options.

## Prerequisites

Before using this script, make sure you have the following prerequisites:

- GoLang installed on your system
- MongoDB installed and running
- An AWS S3 bucket created with the necessary access credentials

## Installation

1. Clone the repository or download the script file.
2. Build the script using the Go compiler:
```sh
go build -o mongo-backup
```
## Usage

The script supports various command-line flags to customize the backup process. Here are the available flags:

- `-uri`: MongoDB URI (required) - The connection string for the MongoDB server.
- `-db`: MongoDB database name (required) - The name of the database to be backed up.
- `-bucket`: AWS S3 bucket name (required) - The name of the S3 bucket where the backup will be stored.
- `-region`: AWS S3 bucket region (required) - The region where the S3 bucket is located.
- `-aws-key`: AWS access key (required) - The access key for your AWS account.
- `-aws-sec`: AWS secret key (required) - The secret key for your AWS account.

To run the script, execute the following command:
```sh
./mongo-backup -uri <mongodb_uri> -db <database_name> -bucket <s3_bucket_name> -region <s3_bucket_region> -aws-key <aws_access_key> -aws-sec <aws_secret_key>
```

Replace `<mongodb_uri>`, `<database_name>`, `<s3_bucket_name>`, `<s3_bucket_region>`, `<aws_access_key>`, and `<aws_secret_key>` with the appropriate values.

## Functionality

1. **Dumping the Database**: The script uses the `mongodump` command to create a compressed backup of the specified MongoDB database. The backup file is saved in the `.bkp/backups` directory with a filename in the format `<database_name>-<timestamp>.gz`.

2. **Pushing to AWS S3**: After the backup file is created, it is uploaded to the specified AWS S3 bucket using the AWS SDK. The backup file is stored in the `database-backups/<database_name>` directory within the S3 bucket.

3. **Validation and Error Handling**: The script validates the provided inputs to ensure that all required parameters are provided. It checks for the existence of the backup file and handles errors during the backup and upload processes.

## Configuration

By default, the script uses the user's home directory to store the backup files in the `.bkp/backups` directory. You can modify the `root_dir` and `backup_dir` variables in the `Application` struct to change the backup directory location.

## License

This script is provided under the [MIT License](LICENSE).

## Acknowledgments

This script utilizes the following third-party libraries:

- [github.com/aws/aws-sdk-go](https://github.com/aws/aws-sdk-go) - AWS SDK for GoLang

## Disclaimer

This script is provided as-is without any warranty. Use it at your own risk.
