# Groupie Tracker

This project is a **web application** built in **Go (Golang)** to provide users with a search feature for tracking information about musical artists. The application is structured using MVC principles and integrates various assets and APIs to offer dynamic and interactive features.

---

## Table of Contents
- [Project Structure](#project-structure)
- [Key Features](#key-features)
- [Setup and Installation](#setup-and-installation)
- [How to Set the Server Port](#how-to-set-the-server-port)
- [API Integration](#api-integration)
- [Contributors](#contributors)

---

## Project Structure

The application is organized as follows:

- **`main.go`**: Entry point of the application.
- **`go.mod`**: Manages Go module dependencies.
- **`.gitignore`**: Specifies files to be ignored by Git.
- **`readme.md`**: Project documentation.

### Directories:
- **`handlers/`**: Contains route logic for different pages:
  - `artist.go`
  - `helpers.go`
  - `home.go`
  - `search.go`
  
- **`models/`**: Defines data structures like:
  - `artist.go`
  - `artistFull.go`
  - `date.go`
  - `location.go`
  - `relation.go`
  
- **`services/`**: Contains API logic:
  - `api.go`
  
- **`static/`**: Stores static assets like:
  - **CSS (`css/`)**:
    - `styleartist.css`
    - `styleFilters.css`
    - `styles.css`
    - `pallete.md`: Color palette or additional documentation.s
  - **Images (`img/`)**:
    - `1.png`
    - `2.png`
    - `3.png`
    - `favicon.png`
  - **JavaScript (`js/`)**:
    - `filters.js`
    - `geolocation.js`
    - `locfix.js`
    - `search.js`
  
- **`templates/`**: HTML templates for rendering:
  - `artist.html`
  - `error.html`
  - `filterModal.html`
  - `home.html`
  - `welcome.html`


---

## Key Features

1. **Search Bar**:
   - Provides users with a real-time search feature to look up artists.
2. **API Integration**:
   - Fetches artist data, locations, and events dynamically.
3. **Filters**:
   - Use filters to limit your search.
4. **Geolocation**:
   - See on a map the locations where the artists have performed.

---

## Setup and Installation

Follow these steps to set up the project:

### Prerequisites
- [Go](https://golang.org/dl/) (version 1.16 or higher)
- Basic knowledge of terminal commands.

### Clone the repository:
```bash
git clone https://github.com/your-repo/groupie-tracker-search-bar.git
cd groupie-tracker-search-bar
```

### Install dependencies:
```bash
go mod tidy
```

### Run the application:
```bash
go run main.go || go run .
```

### Open your browser and navigate to:
```link
http://localhost:8080
```

---

## How to Set the Server Port

By default, the application runs on port `8080`. To specify a different port, follow one of these options:

### 1. Set the PORT Environment Variable
You can specify the desired port by setting the `PORT` environment variable before running the application:

```bash
export PORT=9090
go run main.go
```

This will start the server on port `9090`.

### 2. Use Command-Line Arguments
Alternatively, you can modify the `main.go` file to accept a `-port` flag and set the port dynamically. Example:

```bash
go run main.go -port=9090
```

(Additional implementation would be required in the code to support this flag.)

If no port is specified via an environment variable, the server defaults to port `8080`.

---


## API Integration

The application integrates with an external API to retrieve:

- Artist information
- Event dates and locations
- Relations between artists

APIs are managed in the `services/api.go` file.

---

## Contributors

- Cemvalot  
- Pkalliag  
- Gpatoula

---

