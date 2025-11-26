/**
 * Validate JSON data against a simplified schema
 * @param data - Data to validate
 * @param schema - JSON schema
 * @returns Validation result with errors
 */
export function validateJsonSchema(
  data: any,
  schema: any
): { valid: boolean; errors: string[] } {
  const errors: string[] = [];

  if (!schema || typeof schema !== "object") {
    return { valid: true, errors: [] };
  }

  // Simple schema validation (type checking)
  if (schema.type) {
    const actualType = Array.isArray(data) ? "array" : typeof data;
    if (actualType !== schema.type) {
      errors.push(`Expected type ${schema.type}, got ${actualType}`);
      return { valid: false, errors };
    }
  }

  // Validate object properties
  if (schema.type === "object" && schema.properties) {
    if (typeof data !== "object" || data === null) {
      errors.push("Expected an object");
      return { valid: false, errors };
    }

    for (const [key, propSchema] of Object.entries(schema.properties)) {
      if (schema.required && schema.required.includes(key) && !(key in data)) {
        errors.push(`Missing required property: ${key}`);
      } else if (key in data) {
        const result = validateJsonSchema(data[key], propSchema);
        errors.push(...result.errors.map((e) => `${key}.${e}`));
      }
    }
  }

  // Validate array items
  if (schema.type === "array" && schema.items && Array.isArray(data)) {
    data.forEach((item, index) => {
      const result = validateJsonSchema(item, schema.items);
      errors.push(...result.errors.map((e) => `[${index}].${e}`));
    });
  }

  return { valid: errors.length === 0, errors };
}
