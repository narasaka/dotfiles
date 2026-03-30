-- 001_init.sql

CREATE TABLE IF NOT EXISTS users (
    id TEXT PRIMARY KEY DEFAULT (lower(hex(randomblob(16)))),
    email TEXT UNIQUE NOT NULL,
    password_hash TEXT NOT NULL,
    name TEXT NOT NULL DEFAULT '',
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS apps (
    id TEXT PRIMARY KEY DEFAULT (lower(hex(randomblob(16)))),
    name TEXT UNIQUE NOT NULL,
    display_name TEXT NOT NULL DEFAULT '',
    git_url TEXT NOT NULL,
    git_branch TEXT NOT NULL DEFAULT 'main',
    git_subpath TEXT NOT NULL DEFAULT '',
    dockerfile_path TEXT NOT NULL DEFAULT 'Dockerfile',
    registry_image TEXT NOT NULL DEFAULT '',
    namespace TEXT NOT NULL DEFAULT 'default',
    replicas INTEGER NOT NULL DEFAULT 1,
    port INTEGER NOT NULL DEFAULT 8080,
    env_vars TEXT NOT NULL DEFAULT '{}',
    auto_deploy BOOLEAN NOT NULL DEFAULT 1,
    webhook_secret TEXT NOT NULL DEFAULT '',
    ingress_host TEXT NOT NULL DEFAULT '',
    ingress_tls BOOLEAN NOT NULL DEFAULT 0,
    status TEXT NOT NULL DEFAULT 'inactive',
    current_build_id TEXT,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (current_build_id) REFERENCES builds(id)
);

CREATE TABLE IF NOT EXISTS builds (
    id TEXT PRIMARY KEY DEFAULT (lower(hex(randomblob(16)))),
    app_id TEXT NOT NULL,
    commit_sha TEXT NOT NULL DEFAULT '',
    commit_message TEXT NOT NULL DEFAULT '',
    commit_author TEXT NOT NULL DEFAULT '',
    image_tag TEXT NOT NULL,
    status TEXT NOT NULL DEFAULT 'pending',
    kaniko_job_name TEXT NOT NULL DEFAULT '',
    logs TEXT NOT NULL DEFAULT '',
    started_at DATETIME,
    finished_at DATETIME,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (app_id) REFERENCES apps(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS deployments (
    id TEXT PRIMARY KEY DEFAULT (lower(hex(randomblob(16)))),
    app_id TEXT NOT NULL,
    build_id TEXT NOT NULL,
    k8s_deployment_name TEXT NOT NULL,
    replicas_desired INTEGER NOT NULL DEFAULT 1,
    replicas_ready INTEGER NOT NULL DEFAULT 0,
    status TEXT NOT NULL DEFAULT 'pending',
    rolled_back_to TEXT,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (app_id) REFERENCES apps(id) ON DELETE CASCADE,
    FOREIGN KEY (build_id) REFERENCES builds(id)
);

CREATE TABLE IF NOT EXISTS settings (
    key TEXT PRIMARY KEY,
    value TEXT NOT NULL DEFAULT '',
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

INSERT OR IGNORE INTO settings (key, value) VALUES
    ('registry_url', ''),
    ('registry_username', ''),
    ('registry_password', ''),
    ('default_namespace', 'kubedeck-apps'),
    ('default_domain', ''),
    ('kaniko_image', 'gcr.io/kaniko-project/executor:latest'),
    ('session_secret', lower(hex(randomblob(32))));

CREATE INDEX IF NOT EXISTS idx_builds_app_id ON builds(app_id);
CREATE INDEX IF NOT EXISTS idx_builds_status ON builds(status);
CREATE INDEX IF NOT EXISTS idx_deployments_app_id ON deployments(app_id);
CREATE INDEX IF NOT EXISTS idx_apps_status ON apps(status);
