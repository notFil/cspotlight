export function logError(error: Error, context?: string) {
  console.error('Error:', error);
  if (context) {
    console.error('Context:', context);
  }
}