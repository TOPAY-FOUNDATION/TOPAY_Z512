/**
 * Jest Test Wrapper for TOPAY-Z512 Test Suite
 */

import { runAllTests } from './test';

describe('TOPAY-Z512 Comprehensive Test Suite', () => {
  test('should pass all comprehensive tests', async () => {
    // Capture console output to prevent Jest from being confused by the custom test runner output
    const originalConsoleLog = console.log;
    const logs: string[] = [];
    
    console.log = (...args: any[]) => {
      logs.push(args.join(' '));
    };

    try {
      await runAllTests();
      // If we get here, all tests passed
      expect(true).toBe(true);
    } catch (error) {
      // If any test failed, the runAllTests function will throw
      throw new Error(`Test suite failed: ${error}`);
    } finally {
      // Restore console.log
      console.log = originalConsoleLog;
      
      // Optionally log the captured output for debugging
      if (process.env['VERBOSE_TESTS']) {
        logs.forEach(log => originalConsoleLog(log));
      }
    }
  }, 60000); // 60 second timeout for comprehensive tests
});