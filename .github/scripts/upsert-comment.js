// Shared helper for finding and updating (or creating) a bot comment on a PR.
// Usage from actions/github-script:
//   const upsert = require('./.github/scripts/upsert-comment.js');
//   await upsert({ github, context, marker: 'Test Coverage Report', body: '...' });
module.exports = async ({ github, context, marker, body }) => {
  const { data: comments } = await github.rest.issues.listComments({
    owner: context.repo.owner,
    repo: context.repo.repo,
    issue_number: context.issue.number,
  });

  const botComment = comments.find(comment =>
    comment.user.type === 'Bot' && comment.body.includes(marker)
  );

  if (botComment) {
    await github.rest.issues.updateComment({
      owner: context.repo.owner,
      repo: context.repo.repo,
      comment_id: botComment.id,
      body,
    });
  } else {
    await github.rest.issues.createComment({
      owner: context.repo.owner,
      repo: context.repo.repo,
      issue_number: context.issue.number,
      body,
    });
  }
};
