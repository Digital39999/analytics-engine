import { HttpStatusCode } from 'axios';

export type ResponseType<T> = {
	status: HttpStatusCode.Ok;
	data: T;
} | {
	status: Omit<HttpStatusCode, HttpStatusCode.Ok>;
	error: string;
}

export type RequestData = {
	name: string;
	userId: string;
	createdAt: number;
} | string;

export type ConstructorOptions = {
	authorization: string;
	instanceUrl: string;
}

export type StatsData = {
	totalKeys: number;

	cpuUsage: number;
	ramUsage: string;

	ramUsageBytes: number;

	systemUptime: string;
	goRoutimeCount: number;
}

export type RawStatsData = {
	total_redis_keys: number;

	cpu_usage: number;
	ram_usage: string;

	ram_usage_bytes: number;

	system_uptime: string;
	go_routines: number;
}

// YYYY-MM-DD
export type AnalyticsData<T extends string = string> = {
	global: {
		daily: Record<string, number>;
		weekly: Record<string, number>;
		monthly: Record<string, number>;
	};
	usages: {
		[key in T]: {
			daily: Record<string, number>;
			weekly: Record<string, number>;
			monthly: Record<string, number>;
		};
	};
};
