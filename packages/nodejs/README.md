# Analytics Engine

`analytics-engine-js` is a robust Node.js package designed to interact with a Redis-backed service for logging events, retrieving detailed analytics, and monitoring system performance. This package provides methods to send event data, fetch analytics for specific commands, and obtain system statistics.

## Table of Contents

- [Installation](#installation)
- [Usage](#usage)
  - [Creating an Instance](#creating-an-instance)
  - [Sending Events](#sending-events)
  - [Retrieving Analytics](#retrieving-analytics)
  - [Flushing Statistics](#flushing-statistics)
  - [Getting System Statistics](#getting-system-statistics)
- [Data Structures](#data-structures)
  - [Request Data Example](#request-data-example)
  - [Response Data Example](#response-data-example)
- [Examples](#examples)
- [License](#license)

## Installation

Install the package via npm:

```bash
npm install analytics-engine-js
```

## Usage

### Creating an Instance

Instantiate the `AnalyticsEngine` by providing your authorization token and the instance URL of your analytics service.

```typescript
import AnalyticsEngine from 'analytics-engine-js';

const analyticsClient = new AnalyticsEngine({
    authorization: 'your-authorization-token', // Replace with your actual authorization token
    instanceUrl: 'http://localhost:3000', // Your service URL
});
```

### Sending Events

Use the `event` method to send event data to the analytics service. The method accepts an object with event details.

```typescript
const eventData = {
    name: 'commandA', // The name of the command/event
    uniqueId: 'user1', // The ID of the user triggering the event
    createdAt: Date.now(), // The timestamp in milliseconds
    type: 'commands', // The type of event
};

await analyticsClient.event(eventData);
console.log('Event sent successfully!');
```

### Retrieving Analytics

Retrieve analytics data for specific commands using the `getStatistics` method. You can specify options like lookback period and filters.

```typescript
const statistics = await analyticsClient.getStatistics({
    lookback: 7, // Optional: Specify lookback period in days
    uniqueId: 'user1', // Optional: Filter by unique user ID
    type: 'commands', // Optional: Filter by event type
});
console.log('Analytics Data:', statistics);
```

### Flushing Statistics

Use the `flushStatistics` method to delete analytics data based on specified criteria.

```typescript
const flushResult = await analyticsClient.flushStatistics({
    type: 'commands', // Optional: Specify the type of data to flush
});
console.log('Flush successful:', flushResult);
```

### Getting System Statistics

Monitor the overall performance of the analytics service by using the `getStats` method, which returns system and Redis statistics.

```typescript
const stats = await analyticsClient.getStats();
console.log('System Statistics:', stats);
```

## Data Structures

### Request Data Example

The structure of the data sent when logging an event:

```typescript
export type RequestData = {
    name: string; // The name of the event
    uniqueId?: string; // The ID of the user triggering the event
    createdAt?: number; // The timestamp of the event in milliseconds
    type?: string; // The type of event
};

// Example of RequestData
const eventData: RequestData = {
    name: 'commandA',
    uniqueId: 'user1',
    createdAt: Date.now(),
    type: 'commands',
};
```

### Response Data Example

The structure of the response returned when retrieving analytics data:

```typescript
export type ResponseType<T> = {
    status: HttpStatusCode.Ok;
    data: T; // The data returned on a successful request
} | {
    status: Omit<HttpStatusCode, HttpStatusCode.Ok>; // Any other status codes
    error: string; // Error message
};

// Example ResponseType for Analytics Data
const response: ResponseType<AnalyticsData<string>> = {
    status: HttpStatusCode.Ok,
    data: {
        global: {
            daily: {
                '2024-08-01': 10,
                '2024-08-02': 15,
            },
            weekly: {
                '2024-08-01': 40,
            },
            monthly: {
                '2024-08': 150,
            },
        },
        usages: {
            commandA: {
                daily: {
                    '2024-08-01': 5,
                    '2024-08-02': 8,
                },
                weekly: {
                    '2024-08-01': 30,
                },
                monthly: {
                    '2024-08': 100,
                },
            },
            commandB: {
                daily: {
                    '2024-08-01': 3,
                    '2024-08-02': 7,
                },
                weekly: {
                    '2024-08-01': 10,
                },
                monthly: {
                    '2024-08': 50,
                },
            },
        },
    },
};
```

## Examples

### Fetching Analytics Data

Here is an example of how to fetch and print analytics data:

```typescript
type CommandList = 'commandA' | 'commandB';
const analyticsData = await analyticsClient.getStatistics<CommandList>({
    lookback: 7,
    uniqueId: 'user1',
    type: 'commands,
});
console.log('Analytics Data:', JSON.stringify(analyticsData, null, 2));
```

## License

This project is licensed under the MIT License - see the [LICENSE](https://github.com/Digital39999/analytics-engine/LICENSE) file for details.