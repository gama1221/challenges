# todo-api
 A simple CRUD Todo API written in Go, using ScyllaDB.

## Requirements
- go installed on your machine
- docker installed on your machine
- ScyllaDB using Docker

## Features

- **User Management**: Create user accounts with unique usernames and emails.
- **Todo Management**: CRUD operations for todos, including title, description, and status.
- **Pagination and Sorting**: Fetch todos with pagination, filtering, and sorting options.

## How to setup
1. **Clone the repository using:**
```
git clone https://github.com/aanamshaikh/todo-api.git
```
2. **Environments :**
```
can be found in config.yml and define only the one you require
```
3. **Database:**
```
Make sure you have setup ScyllaDB appropriately and defined the required schema
```
4. **To build the project do:**
```
go build -o todoapp main.go
```
5. **To start:**
```
./todoapp start
```


## Todo API Behavior

- A **Todo** can only be created with an associated **User**.
- Todos can be updated or deleted using their **ID**.
- During pagination, if no status is provided for filtering, the default status will be **Pending**.
- All Todos are created with the **Pending** status.
- For sorting, the default order will be **Ascending**.
- Pagination is implemented using a **cursor** based on the last ID:
  - The last ID must be provided to fetch the next page.
  - The last ID will be returned in the JSON response.

## Status

- **Pending:** The task has been created but not started yet.
- **In Progress:** The task has been started but not completed yet.
  - Can be moved to **Completed**, **Archived**, or **Cancelled**.
  - Cannot be moved back to **Pending**.
- **Completed:** The task has been completed.
  - Can be moved to **Archived**.
  - Cannot be moved to **Pending**, **In Progress**, or **Cancelled**.
- **Archived:** The task has been completed and archived for record-keeping.
  - Can be moved to **In Progress** or **Pending**.
  - Cannot be moved to **Completed** or **Cancelled**.
- **Cancelled:** The task has been cancelled and will not be completed.
  - Cannot be moved to any other status.

## Design Decisions

1. **Implement Interfaces**: Interfaces were used to ensure that the code is flexible and extensible, promoting loose coupling between components. This approach allows for easier testing and future enhancements.

2. **Database Migrations**: Implemented migrations to manage schema changes in a structured manner, ensuring smooth transitions between database versions without data loss.

3. **Configuration Management**: Utilized a configuration management system to read application settings, enhancing maintainability and enabling environment-specific configurations without code changes.

4. **Docker for ScyllaDB**: Leveraged Docker to containerize the ScyllaDB instance, simplifying the setup process and ensuring consistent environments across development and production stages. This also aids in scaling and managing dependencies effectively.


> **Testing**: To test the API, use the Postman collection located at: 
> [Todo-Api.postman_collection.json](./Todo-Api.postman_collection.json)

