import { Method } from "axios";

/**
 * API Request Parameters
 */
export interface ApiRequestParams {
  url: string;
  method: Method;
  headers?: Record<string, string>;
  queryParams?: Record<string, string>;
  body?: any;
  timeout?: number;
}

/**
 * API Test Parameters - extends ApiRequestParams with test-specific fields
 */
export interface ApiTestParams extends ApiRequestParams {
  expectedStatus?: number;
  maxResponseTime?: number;
  jsonSchema?: any;
}

/**
 * API Response structure
 */
export interface ApiResponse {
  status: number;
  statusText: string;
  headers: Record<string, string>;
  body: any;
  responseTime: number;
}

/**
 * Test Result structure
 */
export interface TestResult {
  passed: boolean;
  tests: TestCase[];
  response: ApiResponse;
}

/**
 * Individual Test Case
 */
export interface TestCase {
  name: string;
  passed: boolean;
  message: string;
}
