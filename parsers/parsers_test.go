package parsers

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestIdentifiesEventSource(t *testing.T) {

	trueSourceHints := []string{"User-Agent:GitHub-Hookshot", "User-Agent:GitHub-Hookshot/", "User-Agent:GitHub-Hookshot/5684589df"}
	falseSourceHints := []string{"User-Agent", "GitHub-Hookshot", "GitHub"}

	for i := range falseSourceHints {
		_, err := findEventSource([]string{falseSourceHints[i]})
		assert.NotNil(t, err)
	}

	for i := range trueSourceHints {
		eventSource, err := findEventSource([]string{trueSourceHints[i]})
		assert.Nil(t, err)
		assert.Equal(t, eventSource.Name, "GitHub", "Expected event source to be identified as GitHub")
	}
}

func TestParseJsonEvent(t *testing.T) {

	data := []byte(gitHubPushData.data)
	object, eventSource, eventType, err := ParseWithHints(&data, gitHubPushData.hints)
	assert.Nil(t, err)
	assert.NotNil(t, object)
	assert.Equal(t, "GitHub", eventSource.Name)
	assert.Equal(t, "push", eventType.Type)

	// manually deserialize our test data
	var object2 interface{}
	err = json.Unmarshal([]byte(gitHubPushData.data), &object2)
	assert.Nil(t, err)
	assert.NotNil(t, object2)

	// compare the parsed vs manual objects
	assert.Equal(t, object2, object)
}

/*
	TEST DATA
*/

type TestData struct {
	content string
	hints   []string
	data    string
}

var gitHubPushData = TestData{
	content: "baxterthehacker committed to public-repo - https://github.com/baxterthehacker/public-repo/commit/4d2ab4e76d0d405d17d1a0f2b8a6071394e3ab40",
	hints:   []string{"Content-Type:application/json", "User-Agent:GitHub-Hookshot/458f8", "X-Github-Event:push"},
	data: `{
  "ref": "refs/heads/gh-pages",
  "after": "4d2ab4e76d0d405d17d1a0f2b8a6071394e3ab40",
  "before": "993b46bdfc03ae59434816829162829e67c4d490",
  "created": false,
  "deleted": false,
  "forced": false,
  "compare": "https://github.com/baxterthehacker/public-repo/compare/993b46bdfc03...4d2ab4e76d0d",
  "commits": [
    {
      "id": "4d2ab4e76d0d405d17d1a0f2b8a6071394e3ab40",
      "distinct": true,
      "message": "Trigger pages build",
      "timestamp": "2014-07-25T12:37:40-04:00",
      "url": "https://github.com/baxterthehacker/public-repo/commit/4d2ab4e76d0d405d17d1a0f2b8a6071394e3ab40",
      "author": {
        "name": "Kyle Daigle",
        "email": "kyle.daigle@github.com",
        "username": "kdaigle"
      },
      "committer": {
        "name": "Kyle Daigle",
        "email": "kyle.daigle@github.com",
        "username": "kdaigle"
      },
      "added": [

      ],
      "removed": [

      ],
      "modified": [
        "index.html"
      ]
    }
  ],
  "head_commit": {
    "id": "4d2ab4e76d0d405d17d1a0f2b8a6071394e3ab40",
    "distinct": true,
    "message": "Trigger pages build",
    "timestamp": "2014-07-25T12:37:40-04:00",
    "url": "https://github.com/baxterthehacker/public-repo/commit/4d2ab4e76d0d405d17d1a0f2b8a6071394e3ab40",
    "author": {
      "name": "Kyle Daigle",
      "email": "kyle.daigle@github.com",
      "username": "kdaigle"
    },
    "committer": {
      "name": "Kyle Daigle",
      "email": "kyle.daigle@github.com",
      "username": "kdaigle"
    },
    "added": [

    ],
    "removed": [

    ],
    "modified": [
      "index.html"
    ]
  },
  "repository": {
    "id": 20000106,
    "name": "public-repo",
    "full_name": "baxterthehacker/public-repo",
    "owner": {
      "name": "baxterthehacker",
      "email": "baxterthehacker@users.noreply.github.com"
    },
    "private": false,
    "html_url": "https://github.com/baxterthehacker/public-repo",
    "description": "",
    "fork": false,
    "url": "https://github.com/baxterthehacker/public-repo",
    "forks_url": "https://api.github.com/repos/baxterthehacker/public-repo/forks",
    "keys_url": "https://api.github.com/repos/baxterthehacker/public-repo/keys{/key_id}",
    "collaborators_url": "https://api.github.com/repos/baxterthehacker/public-repo/collaborators{/collaborator}",
    "teams_url": "https://api.github.com/repos/baxterthehacker/public-repo/teams",
    "hooks_url": "https://api.github.com/repos/baxterthehacker/public-repo/hooks",
    "issue_events_url": "https://api.github.com/repos/baxterthehacker/public-repo/issues/events{/number}",
    "events_url": "https://api.github.com/repos/baxterthehacker/public-repo/events",
    "assignees_url": "https://api.github.com/repos/baxterthehacker/public-repo/assignees{/user}",
    "branches_url": "https://api.github.com/repos/baxterthehacker/public-repo/branches{/branch}",
    "tags_url": "https://api.github.com/repos/baxterthehacker/public-repo/tags",
    "blobs_url": "https://api.github.com/repos/baxterthehacker/public-repo/git/blobs{/sha}",
    "git_tags_url": "https://api.github.com/repos/baxterthehacker/public-repo/git/tags{/sha}",
    "git_refs_url": "https://api.github.com/repos/baxterthehacker/public-repo/git/refs{/sha}",
    "trees_url": "https://api.github.com/repos/baxterthehacker/public-repo/git/trees{/sha}",
    "statuses_url": "https://api.github.com/repos/baxterthehacker/public-repo/statuses/{sha}",
    "languages_url": "https://api.github.com/repos/baxterthehacker/public-repo/languages",
    "stargazers_url": "https://api.github.com/repos/baxterthehacker/public-repo/stargazers",
    "contributors_url": "https://api.github.com/repos/baxterthehacker/public-repo/contributors",
    "subscribers_url": "https://api.github.com/repos/baxterthehacker/public-repo/subscribers",
    "subscription_url": "https://api.github.com/repos/baxterthehacker/public-repo/subscription",
    "commits_url": "https://api.github.com/repos/baxterthehacker/public-repo/commits{/sha}",
    "git_commits_url": "https://api.github.com/repos/baxterthehacker/public-repo/git/commits{/sha}",
    "comments_url": "https://api.github.com/repos/baxterthehacker/public-repo/comments{/number}",
    "issue_comment_url": "https://api.github.com/repos/baxterthehacker/public-repo/issues/comments/{number}",
    "contents_url": "https://api.github.com/repos/baxterthehacker/public-repo/contents/{+path}",
    "compare_url": "https://api.github.com/repos/baxterthehacker/public-repo/compare/{base}...{head}",
    "merges_url": "https://api.github.com/repos/baxterthehacker/public-repo/merges",
    "archive_url": "https://api.github.com/repos/baxterthehacker/public-repo/{archive_format}{/ref}",
    "downloads_url": "https://api.github.com/repos/baxterthehacker/public-repo/downloads",
    "issues_url": "https://api.github.com/repos/baxterthehacker/public-repo/issues{/number}",
    "pulls_url": "https://api.github.com/repos/baxterthehacker/public-repo/pulls{/number}",
    "milestones_url": "https://api.github.com/repos/baxterthehacker/public-repo/milestones{/number}",
    "notifications_url": "https://api.github.com/repos/baxterthehacker/public-repo/notifications{?since,all,participating}",
    "labels_url": "https://api.github.com/repos/baxterthehacker/public-repo/labels{/name}",
    "releases_url": "https://api.github.com/repos/baxterthehacker/public-repo/releases{/id}",
    "created_at": 1400625583,
    "updated_at": "2014-07-01T17:21:25Z",
    "pushed_at": 1406306262,
    "git_url": "git://github.com/baxterthehacker/public-repo.git",
    "ssh_url": "git@github.com:baxterthehacker/public-repo.git",
    "clone_url": "https://github.com/baxterthehacker/public-repo.git",
    "svn_url": "https://github.com/baxterthehacker/public-repo",
    "homepage": null,
    "size": 612,
    "stargazers_count": 0,
    "watchers_count": 0,
    "language": null,
    "has_issues": true,
    "has_downloads": true,
    "has_wiki": true,
    "forks_count": 0,
    "mirror_url": null,
    "open_issues_count": 25,
    "forks": 0,
    "open_issues": 25,
    "watchers": 0,
    "default_branch": "master",
    "stargazers": 0,
    "master_branch": "master"
  },
  "pusher": {
    "name": "baxterthehacker",
    "email": "baxterthehacker@users.noreply.github.com"
  }}`,
}