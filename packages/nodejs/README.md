# Analytics Engine

`analytics-engine-js` is a powerful Node.js package that allows you to send event data and retrieve analytical statistics using a Redis-backed service. With this package, you can log events, get detailed analytics per command, and monitor system performance.

## Table of Contents

- [Installation](#installation)
- [Usage](#usage)
  - [Creating an Instance](#creating-an-instance)
  - [Sending Events](#sending-events)
  - [Retrieving Analytics](#retrieving-analytics)
  - [Getting Statistics](#getting-statistics)
- [Data Examples](#data-examples)
- [Examples](#examples)
- [License](#license)

## Installation

You can install `analytics-engine-js` via npm:

```bash
npm install analytics-engine-js
```

## Usage

### Creating an Instance

You can create an instance of the `AnalyticsEngine` by providing your authorization token and the instance URL of your analytics service.

```typescript
import AnalyticsEngine from 'analytics-engine-js';

const analyticsClient = new AnalyticsEngine({
    authorization: 'your-authorization-token', // Replace with your actual authorization token
    instanceUrl: 'http://localhost:3000', // Your service URL
});
```

### Sending Events

You can send event data to the analytics service using the `event` method. This method accepts an object containing the event details.

```typescript
const eventData = {
    name: 'commandA', // The name of the command/event
    uniqueId: 'user1', // The ID of the user triggering the event
    createdAt: Date.now(), // The timestamp in milliseconds
};

await analyticsClient.event(eventData);
console.log('Event sent successfully!');
```

### Retrieving Analytics

You can retrieve analytics data for specific commands using the `getStatistics` method. This method allows you to specify a lookback period.

```typescript
const lookbackDays = 7; // Specify the lookback period in days
const analytics = await analyticsClient.getStatistics<typeof commands>(lookbackDays);
console.log('Analytics Data:', analytics);
```

### Getting Statistics

To monitor the overall performance of the analytics service, you can use the `getStats` method, which returns system and Redis statistics.

```typescript
const stats = await analyticsClient.getStats();
console.log('System Statistics:', stats);
```

## Data Examples

### Request Data Example

The structure of the data sent when logging an event:

```typescript
export type RequestData = {
    name: string; // The name of the event
    uniqueId?: string; // The ID of the user triggering the event
    createdAt: number; // The timestamp of the event in milliseconds
} | string;

// Example of RequestData
const eventData: RequestData = {
    name: 'commandA',
    uniqueId: 'user1',
    createdAt: Date.now(),
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

You can fetch analytics data to analyze event usage over a specified period:

```typescript
type CommandList = 'commandA' | 'commandB';
const analyticsData = await analyticsClient.getStatistics<CommandList>(7);
console.log('Analytics Data:', JSON.stringify(analyticsData, null, 2));
```

## License

This project is licensed under the MIT License - see the [LICENSE](https://github.com/Digital39999/analytics-engine/LICENSE) file for details.