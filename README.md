
# ZincSearch Backend - Email Indexing and Search

This repository contains a backend implementation for indexing and searching emails using [ZincSearch](https://github.com/zinclabs/zinc). The project is designed to process large volumes of email data efficiently, index the data into ZincSearch, and provide a RESTful API for searching indexed emails.

## Features

- **Concurrent Email Processing**: Leverages goroutines and worker pools to process email files efficiently.
- **ZincSearch Integration**: Indexes email data into ZincSearch for fast and scalable search capabilities.
- **RESTful API**: Built with Go and the `chi` router for handling search requests.
- **CORS Support**: Configured to support cross-origin requests.
- **Performance Profiling**: Includes CPU profiling for optimizing processing performance.

## Prerequisites

- Go 1.18+
- ZincSearch instance running locally
- Enron Email dataset ([enron_mail_20110402](http://www.cs.cmu.edu/~enron/enron_mail_20110402.tgz))

## Getting Started

### 1. Clone the repository

```bash
git clone https://github.com/jd-oh/EnronMails_ZincSearchVueGo.git
cd EnronMails_ZincSearchVueGo/Backend
```

### 2. Initialize the Go module

```bash
go mod init zincsearch_vue_go

```

### 3. Install dependencies

```bash
go get github.com/go-chi/chi/v5
```

### 4. Configure ZincSearch

Ensure a ZincSearch instance is running locally at `http://localhost:4080`. The default credentials are:

- **Username**: `admin`
- **Password**: `admin123`

Refer to the [ZincSearch documentation](https://docs.zincsearch.com/) for installation instructions.

### 5. Set up your email dataset

Place the email dataset (e.g., the Enron dataset) in the project directory. Updates the `folderPath` variable in `Indexer.go` to point to the location of the previously extracted dataset.

### 6. Run the Indexer

```bash
go run Indexer.go
```

This will process and index the emails into ZincSearch.

### 7. Run the API Server

```bash
go run main.go
```

The server will start on `http://localhost:8080`.

## API Endpoints

### Search Emails

**Endpoint:** `POST /api/search`

**Request Body:**

```json
{
  "term": "search term",
  "from": 0,
  "size": 10,
  "field": "body"
}
```

- `term`: Search term (required).
- `from`: Starting index for pagination (default: 0).
- `size`: Number of results to return (default: 10).
- `field`: Field to search (default: `body`).

**Example Response:**

```json
[
  {
    "message_id": "<unique-id@example.com>",
    "date": "2025-01-01",
    "from": "sender@example.com",
    "to": "recipient@example.com",
    "subject": "Subject Line",
    "body": "Email content...",
    "folder": "inbox"
  }
]
```

## Performance Profiling

The `Indexer.go` script includes CPU profiling to help optimize performance. The profile is saved to `cpu_profile.prof` and can be analyzed using Go’s `pprof` tool:

```bash
go tool pprof cpu_profile.prof
```
or 
```bash
go tool pprof -http=:8081 cpu_profile.prof
```

# Vue & Tailwind Frontend - Email Search
This repository contains the frontend application built with Vue.js for the Email Indexer project. The application allows users to interact with the indexed email data through a user-friendly interface.

## Features

- **Search Emails**: Perform keyword-based searches on indexed emails.
- **Advanced Filters**: Narrow down results using fields like "from", "to", "subject", etc.
- **Pagination Support**: Browse large sets of results with ease.
- **Responsive Design**: Optimized for both desktop and mobile devices.

## Prerequisites

Before running the frontend, ensure you have the following installed on your machine:

- **Node.js** (version 16 or later recommended)

## Installation

```bash
cd EnronMails_ZincSearchVueGo/Frontend/email-search-vue
```

Install dependences 
```bash
npm install
```

Running the Development Server
```bash
npm run dev
```

Open your browser and navigate to: http://localhost:5173/


## Try the Deployed Application

You can try the application directly without setting it up locally. It's deployed at the following address:

**[http://3.149.23.100/](http://3.149.23.100/)**

Experience the full functionality of the platform and provide feedback!



