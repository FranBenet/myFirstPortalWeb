
# CARS VIEWER

Cars Viewer is a project that includes a web server and a web interface working together with an API.
The web includes a gallery of cars, a search bar and a filter menu. It also allows you to compare cars, create a favourite list and review the last comparison made by the user.




## Usage

To run the application, follow these steps:
- Clone the Repository.
- Install the API. This API provides the data for the car models and needs to be installed separately. Follow the instructions in the API's README file to install it. (cars/api/readme)
- In terminal, navigate to the API directory (cars/api) and start the API using the following command: `make run`
- In another terminal(split terminal), navigate to the root directory for the project (/cars) and start the server by running: `go run ./cmd`
- Finally, access your browser and go to: http://localhost:8080 to get in the website.