module.exports = async ({ github, context }) => {
  const fs = require('fs');

  let backendCoverage = 'N/A';
  let frontendStats = { statements: 'N/A', branches: 'N/A', functions: 'N/A', lines: 'N/A' };

  try {
    const backendOutput = fs.readFileSync('coverage/backend/coverage.txt', 'utf8');
    const totalLine = backendOutput.split('\n').find(line => line.includes('total:'));
    if (totalLine) {
      const parts = totalLine.trim().split(/\s+/);
      backendCoverage = parts[parts.length - 1];
    }
  } catch (e) {
    console.log('Error reading backend coverage:', e);
  }

  try {
    const frontendOutput = JSON.parse(fs.readFileSync('coverage/frontend/coverage-summary.json', 'utf8'));
    const total = frontendOutput.total;
    frontendStats = {
      statements: total.statements.pct + '%',
      branches: total.branches.pct + '%',
      functions: total.functions.pct + '%',
      lines: total.lines.pct + '%'
    };
  } catch (e) {
     console.log('Error reading frontend coverage:', e);
  }

  const header = '## Code Coverage Report';
  const body = `${header}

| Component | Metric | Percentage |
|---|---|---|
| **Backend** | Total | ${backendCoverage} |
| **Frontend** | Statements | ${frontendStats.statements} |
| **Frontend** | Branches | ${frontendStats.branches} |
| **Frontend** | Functions | ${frontendStats.functions} |
| **Frontend** | Lines | ${frontendStats.lines} |
`;

  const { owner, repo, number } = context.issue;

  const comments = await github.rest.issues.listComments({
    owner,
    repo,
    issue_number: number,
  });

  const botComment = comments.data.find(c => c.body.startsWith(header) && c.user.type === 'Bot');

  if (botComment) {
    await github.rest.issues.updateComment({
      owner,
      repo,
      comment_id: botComment.id,
      body: body
    });
  } else {
    await github.rest.issues.createComment({
      owner,
      repo,
      issue_number: number,
      body: body
    });
  }
};
