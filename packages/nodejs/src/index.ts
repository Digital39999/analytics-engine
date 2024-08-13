import { AnalyticsData, ConstructorOptions, RawStatsData, RequestData, ResponseType, StatsData } from './types';
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

		if (!data || data.status !== 200) throw new Error('Request failed.');
		else if ('error' in data) throw new Error(data.error);

		return data.data;
	}

	private getHeaders(): Record<string, string> {
		return {
			'Authorization': this.options.authorization,
			'Content-Type': 'application/json',
		};
	}

	public async event(type: string, data: RequestData): Promise<boolean> {
		return !!(await this.parseAxiosRequest<string>(axios({
			method: 'POST',
			url: `${this.options.instanceUrl}/event`,
			headers: this.getHeaders(),
			data: typeof data === 'string' ? {
				name: data,
				type,
			} : {
				...data,
				type,
			},
		})));
	}

	public async getStatistics<T extends string>(type: string, lookback?: number): Promise<AnalyticsData<T>> {
		return await this.parseAxiosRequest<AnalyticsData<T>>(axios({
			method: 'GET',
			url: `${this.options.instanceUrl}/analytics` + (type ? `?type=${type}` : '') + (lookback ? `&lookback=${lookback}` : ''),
			headers: this.getHeaders(),
		}));
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
			goRoutimeCount: data.go_routines,
		};
	}
}
