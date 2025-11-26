import axios, { AxiosRequestConfig, AxiosResponse } from "axios";
import { ApiRequestParams, ApiResponse } from "../types/index.js";

/**
 * Make an HTTP API request
 * @param params - Request parameters
 * @returns API response with timing information
 */
export async function makeApiRequest(
  params: ApiRequestParams
): Promise<ApiResponse> {
  const startTime = Date.now();

  try {
    const config: AxiosRequestConfig = {
      method: params.method,
      url: params.url,
      headers: params.headers || {},
      params: params.queryParams || {},
      data: params.body,
      timeout: params.timeout || 30000,
      validateStatus: () => true, // Don't throw on any status
    };

    const response: AxiosResponse = await axios(config);
    const endTime = Date.now();

    return {
      status: response.status,
      statusText: response.statusText,
      headers: response.headers as Record<string, string>,
      body: response.data,
      responseTime: endTime - startTime,
    };
  } catch (error: any) {
    const endTime = Date.now();
    throw new Error(
      `Request failed: ${error.message || "Unknown error"} (Response time: ${
        endTime - startTime
      }ms)`
    );
  }
}
