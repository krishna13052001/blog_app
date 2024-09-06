# README

Welcome to the blog_app!

This is a simple blogging application built with Go. It allows users to create, read, update, and delete blog posts.

## Installation

1. Clone the repository:

    ```
    git clone https://github.com/krishna13052001/blog_app.git
    ```

2. Change into the project directory:

    ```
    cd blog_app
    ```

3. Install the dependencies:

    ```
    go mod download
    ```

4. Build the application:

    ```
    go build
    ```

5. Run the application:

    ```
    ./blog_app
    ```

## Usage

Once the application is running, you can access it in your web browser at `http://localhost:8080`. From there, you can create new blog posts, view existing posts, update posts, and delete posts.

## Contributing

Contributions are welcome! If you find any issues or have suggestions for improvements, please open an issue or submit a pull request.

## Routes

The blog_app has the following routes:

- `/posts` - GET: Retrieves all blog posts
- `/posts/{id}` - GET: Retrieves a specific blog post by ID
- `/posts` - POST: Creates a new blog post
- `/posts/{id}` - PUT: Updates a specific blog post by ID
- `/posts/{id}` - DELETE: Deletes a specific blog post by ID

## Frontend Contribution

We welcome contributions to the frontend of the blog_app. If you are interested in contributing, please follow these steps:

We appreciate your contributions and look forward to your pull requests!
