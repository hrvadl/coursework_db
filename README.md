# My coursework project. This was intended to be a financial application

Tech stack:

- Go
- Chi router
- Gorm
- MySQL
- HTMX
- TailwindCSS

## How to run?

### In Docker

1. Create .prod.env file in the root of the project
2. Fill the .prod.env file with needed variables
3. Make sure you have Docker & Docker compose installed
4. Run `make up`

### Partially Locally

1. Create .env file in the root of the project
2. Fill the .env file with needed variables
3. Make sure you have Docker & Docker compose installed
4. Run `make dev` to run the DB in the Docker
5. Run `make run` to run the application locally

### Fully locally

1. Install MySQL locally
2. Create needed DB, user etc.
3. Fill the .env file with needed variables
4. Run `make run` to run the application

## Screenshots

<img width="1205" alt="image" src="https://github.com/hrvadl/coursework_db/assets/93580374/982729df-c271-43b1-bd1e-95e06ddbc7bd">
<img width="1210" alt="image" src="https://github.com/hrvadl/coursework_db/assets/93580374/7f9770de-6563-466c-8c6d-f6f932740f5e">
<img width="1202" alt="image" src="https://github.com/hrvadl/coursework_db/assets/93580374/c9331ef5-5b8e-43c7-bb5c-14a31f4f3f61">
