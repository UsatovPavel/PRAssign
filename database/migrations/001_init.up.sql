CREATE TABLE teams (
    team_name VARCHAR(40) PRIMARY KEY
);

CREATE TABLE users (
    user_id VARCHAR(40) PRIMARY KEY,
    username VARCHAR(40) NOT NULL,
    team_name VARCHAR(40) NOT NULL REFERENCES teams(team_name) ON DELETE CASCADE,
    is_active BOOLEAN NOT NULL DEFAULT TRUE
);

CREATE INDEX idx_users_team_active ON users(team_name, is_active);

CREATE TABLE pull_requests (
    pull_request_id VARCHAR(40) PRIMARY KEY,
    pull_request_name VARCHAR(40) NOT NULL,
    author_id VARCHAR(40) NOT NULL REFERENCES users(user_id),
    status VARCHAR(40) NOT NULL CHECK (status IN ('OPEN','MERGED')) DEFAULT 'OPEN',
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    merged_at TIMESTAMPTZ NULL
);

CREATE TABLE pull_request_reviewers (
    pull_request_id VARCHAR(40) REFERENCES pull_requests(pull_request_id) ON DELETE CASCADE,
    user_id VARCHAR(40) REFERENCES users(user_id),
    assigned_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    PRIMARY KEY (pull_request_id, user_id)
);

CREATE INDEX idx_pr_reviewers_pr ON pull_request_reviewers(pull_request_id);
CREATE INDEX idx_pr_reviewers_user ON pull_request_reviewers(user_id);
