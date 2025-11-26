import { Tool } from "@modelcontextprotocol/sdk/types.js";

/**
 * API Agent Tool Definitions
 */
export const apiAgentTools: Tool[] = [
  {
    name: "api_get",
    description: "Send a GET request to a REST API endpoint",
    inputSchema: {
      type: "object",
      properties: {
        url: {
          type: "string",
          description: "The full URL to send the GET request to",
        },
        headers: {
          type: "object",
          description: "Optional HTTP headers as key-value pairs",
        },
        queryParams: {
          type: "object",
          description: "Optional query parameters as key-value pairs",
        },
        timeout: {
          type: "number",
          description: "Request timeout in milliseconds (default: 30000)",
        },
      },
      required: ["url"],
    },
  },
  {
    name: "api_post",
    description: "Send a POST request to a REST API endpoint",
    inputSchema: {
      type: "object",
      properties: {
        url: {
          type: "string",
          description: "The full URL to send the POST request to",
        },
        headers: {
          type: "object",
          description: "Optional HTTP headers as key-value pairs",
        },
        queryParams: {
          type: "object",
          description: "Optional query parameters as key-value pairs",
        },
        body: {
          type: ["object", "string", "array"],
          description: "Request body (JSON object, array, or string)",
        },
        timeout: {
          type: "number",
          description: "Request timeout in milliseconds (default: 30000)",
        },
      },
      required: ["url"],
    },
  },
  {
    name: "api_put",
    description: "Send a PUT request to a REST API endpoint",
    inputSchema: {
      type: "object",
      properties: {
        url: {
          type: "string",
          description: "The full URL to send the PUT request to",
        },
        headers: {
          type: "object",
          description: "Optional HTTP headers as key-value pairs",
        },
        queryParams: {
          type: "object",
          description: "Optional query parameters as key-value pairs",
        },
        body: {
          type: ["object", "string", "array"],
          description: "Request body (JSON object, array, or string)",
        },
        timeout: {
          type: "number",
          description: "Request timeout in milliseconds (default: 30000)",
        },
      },
      required: ["url"],
    },
  },
  {
    name: "api_patch",
    description: "Send a PATCH request to a REST API endpoint",
    inputSchema: {
      type: "object",
      properties: {
        url: {
          type: "string",
          description: "The full URL to send the PATCH request to",
        },
        headers: {
          type: "object",
          description: "Optional HTTP headers as key-value pairs",
        },
        queryParams: {
          type: "object",
          description: "Optional query parameters as key-value pairs",
        },
        body: {
          type: ["object", "string", "array"],
          description: "Request body (JSON object, array, or string)",
        },
        timeout: {
          type: "number",
          description: "Request timeout in milliseconds (default: 30000)",
        },
      },
      required: ["url"],
    },
  },
  {
    name: "api_delete",
    description: "Send a DELETE request to a REST API endpoint",
    inputSchema: {
      type: "object",
      properties: {
        url: {
          type: "string",
          description: "The full URL to send the DELETE request to",
        },
        headers: {
          type: "object",
          description: "Optional HTTP headers as key-value pairs",
        },
        queryParams: {
          type: "object",
          description: "Optional query parameters as key-value pairs",
        },
        body: {
          type: ["object", "string", "array"],
          description: "Optional request body (JSON object, array, or string)",
        },
        timeout: {
          type: "number",
          description: "Request timeout in milliseconds (default: 30000)",
        },
      },
      required: ["url"],
    },
  },
  {
    name: "api_test",
    description:
      "Send an API request and run automated tests on the response (status code, response time, schema validation)",
    inputSchema: {
      type: "object",
      properties: {
        url: {
          type: "string",
          description: "The full URL to send the request to",
        },
        method: {
          type: "string",
          enum: ["GET", "POST", "PUT", "PATCH", "DELETE"],
          description: "HTTP method to use",
        },
        headers: {
          type: "object",
          description: "Optional HTTP headers as key-value pairs",
        },
        queryParams: {
          type: "object",
          description: "Optional query parameters as key-value pairs",
        },
        body: {
          type: ["object", "string", "array"],
          description: "Optional request body (JSON object, array, or string)",
        },
        expectedStatus: {
          type: "number",
          description: "Expected HTTP status code (e.g., 200, 201, 404)",
        },
        maxResponseTime: {
          type: "number",
          description: "Maximum acceptable response time in milliseconds",
        },
        jsonSchema: {
          type: "object",
          description:
            "JSON schema to validate the response body against (simplified schema with type, properties, required)",
        },
        timeout: {
          type: "number",
          description: "Request timeout in milliseconds (default: 30000)",
        },
      },
      required: ["url", "method"],
    },
  },
];
