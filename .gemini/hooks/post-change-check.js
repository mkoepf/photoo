#!/usr/bin/env node
const { execSync } = require('child_process');
const fs = require('fs');

// Prevent recursion if the hook itself or scripts it calls modify files
if (process.env.SKIP_CHECK === 'true') {
  process.exit(0);
}

let input;
try {
  input = JSON.parse(fs.readFileSync(0, 'utf-8'));
} catch (e) {
  process.exit(0); // If no JSON input, skip
}

// Only run for successful tool calls
if (input.tool_response && input.tool_response.error) {
  process.exit(0);
}

try {
  // Run the check script
  // We set SKIP_CHECK=true to prevent potential recursion if check.sh modifies files (e.g. go fmt)
  const output = execSync('./scripts/check.sh', { 
    env: { ...process.env, SKIP_CHECK: 'true' },
    encoding: 'utf-8' 
  });
  
  // Provide success feedback
  console.log(JSON.stringify({
    hookSpecificOutput: {
      additionalContext: `

✅ [Hook]: Code quality checks passed after modification.
${output}`
    }
  }));
} catch (err) {
  // If check.sh fails (exit code != 0), replace the tool's success with the error message
  // This forces the agent to fix the issues immediately.
  console.log(JSON.stringify({
    decision: "deny",
    reason: `❌ [Hook]: Code quality checks FAILED after modification. Please fix the following issues:

${err.stdout || err.message}`,
    systemMessage: "Quality check failed"
  }));
}
