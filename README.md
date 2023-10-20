Weather App Installation Guide
Backend (Golang)

1. Clone the Repository
   Clone the backend repository to your local machine using the following command:

git clone https://github.com/nihalchrs/weather-app-backend.git

2. Create Environment File
   In the backend directory, create a .env file based on the provided env_sample file. Update the environment variables as needed, including your MySQL database configuration.

3. Set up MySQL Database
   Make sure you have MySQL installed. Create a database with the name as same as you specify in the backend's .env file.

4. Install Dependencies
   In the backend directory, run the following commands to install required packages:

   go get .
   Additionally, initialize and tidy your Go modules:

   go mod init
   go mod tidy

5. Database Schema Migration
   In the backend repository, there is a shell script inside the script directory responsible for database schema migration. Execute this script to create the necessary tables in your MySQL database:

   chmod +x scripts/runMigration.sh
   ./scripts/runMigration.sh

6. Run the Application
   To run the Golang backend, use the following command in the backend directory:

   go run main.go
