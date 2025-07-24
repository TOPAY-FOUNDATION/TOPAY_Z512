/**
 * Individual Test Runner for TOPAY-Z512
 * Runs each test category separately for detailed verification
 */

/// <reference types="node" />

// Import all test functions from test.ts
import {
  testConstants,
  testUtilities,
  testHashOperations,
  testKeyPairOperations,
  testKEMOperations,
  testFragmentation,
  testPerformance,
  testIntegration,
  TestRunner
} from './test';

async function runIndividualTests(): Promise<void> {
  console.log('🧪 TOPAY-Z512 Individual Test Runner');
  console.log('='.repeat(50));
  console.log('Running each test category individually...\n');

  const testCategories = [
    { name: 'Constants', fn: testConstants },
    { name: 'Utility Functions', fn: testUtilities },
    { name: 'Hash Operations', fn: testHashOperations },
    { name: 'Key Pair Operations', fn: testKeyPairOperations },
    { name: 'KEM Operations', fn: testKEMOperations },
    { name: 'Fragmentation', fn: testFragmentation },
    { name: 'Performance', fn: testPerformance },
    { name: 'Integration', fn: testIntegration }
  ];

  let overallPassed = 0;
  let overallFailed = 0;
  let overallDuration = 0;

  for (let i = 0; i < testCategories.length; i++) {
    const category = testCategories[i]!;

    console.log(`\n${'='.repeat(60)}`);
    console.log(`📋 Test Category ${i + 1}/${testCategories.length}: ${category.name}`);
    console.log(`${'='.repeat(60)}`);

    const runner = new TestRunner();

    try {
      await category.fn(runner);
      const summary = runner.getSummary();
      overallPassed += summary.totalPassed;
      overallFailed += summary.totalFailed;
      overallDuration += summary.totalDuration;

      console.log(`\n✅ ${category.name} completed successfully!`);
    } catch (error) {
      console.log(`\n❌ ${category.name} failed: ${error}`);
      overallFailed++;
    }

    // Small delay between categories
    await new Promise(resolve => setTimeout(resolve, 100));
  }

  // Final summary
  console.log(`\n${'='.repeat(60)}`);
  console.log('🎯 OVERALL TEST SUMMARY');
  console.log(`${'='.repeat(60)}`);
  console.log(`Total Tests: ${overallPassed + overallFailed}`);
  console.log(`Passed: ${overallPassed}`);
  console.log(`Failed: ${overallFailed}`);
  console.log(
    `Success Rate: ${((overallPassed / (overallPassed + overallFailed)) * 100).toFixed(1)}%`
  );
  console.log(`Total Duration: ${overallDuration}ms`);

  if (overallFailed > 0) {
    console.log('\n❌ Some test categories failed!');
    process.exit(1);
  } else {
    console.log('\n✅ All test categories passed!');
  }
}

// Run if this file is executed directly
if (require.main === module) {
  runIndividualTests().catch(console.error);
}

export { runIndividualTests };
