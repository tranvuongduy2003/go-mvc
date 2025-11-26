import { Method } from "axios";
import {
  ApiRequestParams,
  ApiTestParams,
  makeApiRequest,
  TestCase,
  TestResult,
  validateJsonSchema,
} from "../../shared/index.js";

/**
 * Handle basic HTTP requests (GET, POST, PUT, PATCH, DELETE)
 */
export async function handleHttpRequest(
  method: Method,
  args: any
): Promise<any> {
  const params: ApiRequestParams = {
    url: args.url as string,
    method,
    headers: args.headers as Record<string, string>,
    queryParams: args.queryParams as Record<string, string>,
    body: args.body,
    timeout: args.timeout as number,
  };

  const response = await makeApiRequest(params);

  return {
    success: true,
    request: {
      method,
      url: params.url,
      headers: params.headers,
      queryParams: params.queryParams,
      body: params.body,
    },
    response: {
      status: response.status,
      statusText: response.statusText,
      headers: response.headers,
      body: response.body,
      responseTime: `${response.responseTime}ms`,
    },
  };
}

/**
 * Handle API testing with validations
 */
export async function handleApiTest(args: any): Promise<any> {
  const params: ApiTestParams = {
    url: args.url as string,
    method: args.method as Method,
    headers: args.headers as Record<string, string>,
    queryParams: args.queryParams as Record<string, string>,
    body: args.body,
    timeout: args.timeout as number,
    expectedStatus: args.expectedStatus as number,
    maxResponseTime: args.maxResponseTime as number,
    jsonSchema: args.jsonSchema,
  };

  const response = await makeApiRequest(params);
  const tests: TestCase[] = [];

  // Test 1: Status code
  if (params.expectedStatus !== undefined) {
    const passed = response.status === params.expectedStatus;
    tests.push({
      name: "Status Code",
      passed,
      message: passed
        ? `Status code is ${response.status} as expected`
        : `Expected status ${params.expectedStatus}, got ${response.status}`,
    });
  }

  // Test 2: Response time
  if (params.maxResponseTime !== undefined) {
    const passed = response.responseTime <= params.maxResponseTime;
    tests.push({
      name: "Response Time",
      passed,
      message: passed
        ? `Response time ${response.responseTime}ms is within limit`
        : `Response time ${response.responseTime}ms exceeds ${params.maxResponseTime}ms`,
    });
  }

  // Test 3: JSON Schema validation
  if (params.jsonSchema) {
    const validation = validateJsonSchema(response.body, params.jsonSchema);
    tests.push({
      name: "JSON Schema",
      passed: validation.valid,
      message: validation.valid
        ? "Response body matches expected schema"
        : `Schema validation failed: ${validation.errors.join(", ")}`,
    });
  }

  const allPassed = tests.every((t) => t.passed);

  const result: TestResult = {
    passed: allPassed,
    tests,
    response,
  };

  return {
    success: true,
    testResults: {
      summary: {
        totalTests: tests.length,
        passed: tests.filter((t) => t.passed).length,
        failed: tests.filter((t) => !t.passed).length,
        allPassed,
      },
      tests: result.tests,
    },
    request: {
      method: params.method,
      url: params.url,
      headers: params.headers,
      queryParams: params.queryParams,
      body: params.body,
    },
    response: {
      status: response.status,
      statusText: response.statusText,
      headers: response.headers,
      body: response.body,
      responseTime: `${response.responseTime}ms`,
    },
  };
}
