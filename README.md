# Analytics Engine

Analytics Engine is a microservice designed to capture and analyze user events in real-time using Redis. It aggregates and provides insights into daily, weekly, and monthly usage, offering both global and command-specific statistics.

## Why Use Analytics Engine?

### Challenges in Gathering Real-time Analytics

When building applications that require real-time insights, several challenges often arise:

1. **Scalability**: 
   - Traditional analytics solutions can struggle with large volumes of data, leading to bottlenecks. Analytics Engine scales horizontally, efficiently handling increasing loads.

2. **Data Persistence**: 
   - Ensuring that usage data is stored and retrievable after server restarts or crashes is crucial. By leveraging Redis, Analytics Engine provides reliable persistence for your analytics data.

3. **Resource Efficiency**: 
   - Analytics systems can consume significant resources if not optimized. Analytics Engine efficiently manages and stores only the necessary data, reducing overhead.

4. **Ease of Integration**: 
   - Integrating analytics into an application should be straightforward. Analytics Engine offers a simple API, making it easy to send and retrieve event data without additional complexity.

5. **Real-time Insights**:
   - Many systems struggle with delivering real-time analytics. Analytics Engine provides immediate access to usage data, enabling you to make informed decisions quickly.

### Benefits of Analytics Engine

By using Analytics Engine, you gain:

- **Scalability**: Efficiently handles large volumes of events across distributed instances.
- **Reliability**: Persistent data storage ensures that your analytics are accurate and durable.
- **Simplicity**: A clean, easy-to-use API for tracking and aggregating event data.
- **Flexibility**: Aggregate data by the day, week, or month, and get detailed insights into each command.

# Table of Contents

- [Features](#features)
- [Requirements](#requirements)
- [Environment Variables](#environment-variables)
- [Running with Docker](#running-with-docker)
- [Running with Docker Compose](#running-with-docker-compose)
- [API Usage](#api-usage)
- [Packages](/packages/README.md)
- [Contributing](#contributing)
- [License](#license)

## Features

- Track and store user events with timestamps.
- Aggregate event data into daily, weekly, and monthly summaries.
- Retrieve global usage statistics as well as per-command insights.
- Easy deployment with Docker and Docker Compose.

## Requirements

- Docker installed on your server or local machine.
- Docker Compose (optional, but recommended).

## Environment Variables

Set the following environment variables before running the service:

- `REDIS_URL`: URL of the Redis server (e.g., `redis://localhost:6379`).
- `REDIS_KEY`: Redis key prefix for storing event data (default is `analyticsEngine`).
- `PORT`: Port on which the service will run (default is `8080`).
- `API_AUTH`: API key for securing the service.
- `MAX_AGE`: Maximum retention time for event data in days (e.g., `30`).

## Running with Docker

You can easily run Analytics Engine using Docker. Follow these steps:

<details>
<summary>1. Pull the Docker Image</summary>

```bash
docker pull ghcr.io/digital39999/analytics-engine:latest
```

</details>

<details>
<summary>2. Run the Container</summary>

Run the container with the necessary environment variables:

```bash
docker run -d \
  -e REDIS_URL="redis://your-redis-url:6379" \
  -e PORT=8080 \
  -e MAX_AGE=30 \
  -e API_AUTH="your-api-key" \
  -p 8080:8080 \
  ghcr.io/digital39999/redis-analytics-engine:latest
```

</details>

<details>
<summary>3. Access the Service</summary>

The service will be available at `http://localhost:8080`.

</details>

## Running with Docker Compose

If you prefer to use Docker Compose, follow these steps:

<details>
<summary>1. Create a `docker-compose.yml` File</summary>

Here’s an example `docker-compose.yml`:

```yaml
version: '3.8'

services:
  analytics-engine:
    image: ghcr.io/digital39999/analytics-engine:latest
    environment:
      REDIS_URL: "redis://redis:6379"
      PORT: 8080
      API_AUTH: "your-api-key"
      MAX_AGE: 365
    ports:
      - "8080:8080"
    depends_on:
      - redis

  redis:
    image: redis:latest
    ports:
      - "6379:6379"
```

</details>

<details>
<summary>2. Run the Services</summary>

To start the services, use the following command:

```bash
docker-compose up -d
```

</details>

<details>
<summary>3. Access the Service</summary>

Once the services are up, you can access Analytics Engine at `http://localhost:8080`.

</details>

## API Usage

<details>
<summary>Routes Overview</summary>

### Record an Event

You can record an event by sending a POST request to `/event`. Here's an example using `curl`:

```bash
curl -X POST http://localhost:8080/event \
-H "Content-Type: application/json" \
-H "Authorization: your-api-key" \
-d '{
  "name": "login",
  "userId": "user123",
  "createdAt": 1691913600
}'
```

- **`name`**: The event name (e.g., `login`, `purchase`).
- **`userId`**: A unique identifier for the user.
- **`createdAt`**: The timestamp of the event (in Unix time).

### Get Aggregated Analytics

To retrieve aggregated analytics, send a GET request to `/analytics`:

```bash
curl -X GET http://localhost:8080/analytics
```

- **Optional Query Parameter**:
  - `lookback`: Number of days to look back for daily counts (default is `7`).

### Get System Statistics

To retrieve system statistics, send a GET request to `/stats`:

```bash
curl -X GET http://localhost:8080/stats
```

- This will return information such as the total Redis keys, CPU and RAM usage, and system uptime.

</details>

<details>
<summary>Examples</summary>

### Example Node.js Client

Here’s how you could integrate Analytics Engine into a Node.js project:

```javascript
const apiUrl = 'http://localhost:8080/event';

async function recordEvent() {
  try {
    const response = await fetch(apiUrl, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
        'Authorization': 'your-api-key'
      },
      body: JSON.stringify({
        name: 'login',
        userId: 'user123',
        createdAt: Math.floor(Date.now() / 1000) // current Unix time
      })
    }).then(res => res.json());

    if (response.error) throw new Error(response.error);
    console.log('Event recorded successfully:', response.data);
  } catch (error) {
    console.error('Error recording event:', error.message);
  }
}

recordEvent();
```

### Example Python Client

Here’s how you could integrate Analytics Engine into a Python project using `requests`:

```python
import requests
import time

api_url = 'http://localhost:8080/event'

def record_event():
    data = {
        'name': 'login',
        'userId': 'user123',
        'createdAt': int(time.time())  # current Unix time
    }
    
    response = requests.post(api_url, json=data)
    if response.status_code == 200:
        print('Event recorded successfully:', response.json())
    else:
        print('Error recording event:', response.text)

record_event()
```

</details>

## Contributing

If you'd like to contribute to this project, feel free to open a pull request or submit an issue on the GitHub repository.

## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for more details.