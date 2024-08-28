import { AnalyticsData, ConstructorOptions, FlushOptions, RawStatsData, RequestData, ResponseType, StatisticOptions, StatsData } from './types';
import axios, { AxiosError, AxiosResponse } from 'axios';

export * from './types';

export default class AnalyticsEngine {
	private options: ConstructorOptions;

	constructor ({
		authorization,
		instanceUrl,
	}: ConstructorOptions) {
		this.options = {
			authorization,
			instanceUrl,
		};

		this.checkInstanceUrl(instanceUrl);
	}

	private async checkInstanceUrl(url: string): Promise<void> {
		if (!url) throw new Error('Instance URL is required.');

		const instanceUrl = new URL(url);
		if (!instanceUrl.hostname) throw new Error('Invalid instance URL.');

		const response = await axios<string>({ url, method: 'GET' }).catch((err: AxiosError<string>) => err.response);
		if (!response || response.status !== 200) throw new Error('Invalid instance URL.');
		else if ('error' in response && typeof response.error === 'string') throw new Error(response.error);

		return;
	}

	private async parseAxiosRequest<T>(response: Promise<AxiosResponse<ResponseType<T>>>): Promise<T> {
		const data = await response.then((res) => res.data).catch((err: AxiosError<ResponseType<T>>) => err.response?.data);

		if (!data) throw new Error('Request failed.');
		else if ('error' in data) throw new Error(data.error);
		else if (!data.data) throw new Error('Invalid response data.');

		return data.data;
	}

	public qp<T extends Record<string, unknown>>(url: string, params?: T) {
		if (!params) return url;

		const query = new URLSearchParams();

		for (const [key, value] of Object.entries(params)) {
			if (value === undefined) continue;
			query.append(key, String(value));
		}

		return url + '?' + query.toString();
	}

	private getHeaders(): Record<string, string> {
		return {
			'Authorization': this.options.authorization,
			'Content-Type': 'application/json',
		};
	}

	public async event(data: RequestData): Promise<boolean> {
		return !!(await this.parseAxiosRequest<string>(axios({
			method: 'POST',
			url: this.qp(`${this.options.instanceUrl}/event`),
			headers: this.getHeaders(),
			data: {
				...data,
				createdAt: data.createdAt || Date.now(),
			},
		})));
	}

	public async getStatistics<T extends string>(options?: StatisticOptions): Promise<AnalyticsData<T>> {
		return await this.parseAxiosRequest<AnalyticsData<T>>(axios({
			method: 'GET',
			url: this.qp(`${this.options.instanceUrl}/analytics`, options),
			headers: this.getHeaders(),
		}));
	}

	public async flushStatistics(options?: FlushOptions): Promise<boolean> {
		return !!(await this.parseAxiosRequest<string>(axios({
			method: 'DELETE',
			url: this.qp(`${this.options.instanceUrl}/analytics`, options),
			headers: this.getHeaders(),
		})));
	}

	public async getStats(): Promise<StatsData> {
		const data = await this.parseAxiosRequest<RawStatsData>(axios({
			method: 'GET',
			url: `${this.options.instanceUrl}/stats`,
			headers: this.getHeaders(),
		}));

		return {
			totalKeys: data.total_redis_keys,
			cpuUsage: data.cpu_usage,
			ramUsage: data.ram_usage,
			ramUsageBytes: data.ram_usage_bytes,
			systemUptime: data.system_uptime,
			systemUptimeSeconds: data.system_uptime_seconds,
			goRoutimeCount: data.go_routines,
		};
	}
}
