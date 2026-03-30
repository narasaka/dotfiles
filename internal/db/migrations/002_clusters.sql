-- 002_clusters.sql

CREATE TABLE IF NOT EXISTS clusters (
    id TEXT PRIMARY KEY DEFAULT (lower(hex(randomblob(16)))),
    name TEXT NOT NULL,
    display_name TEXT NOT NULL DEFAULT '',
    provider TEXT NOT NULL,                     -- 'gke', 'eks', 'aks', etc.
    provider_cluster_id TEXT NOT NULL DEFAULT '', -- provider-specific cluster ID (e.g. 'us-central1-a/my-cluster')
    project_id TEXT NOT NULL DEFAULT '',         -- cloud project/account ID
    location TEXT NOT NULL DEFAULT '',           -- region/zone
    status TEXT NOT NULL DEFAULT 'pending',      -- pending, connecting, connected, error, disconnected
    k8s_version TEXT NOT NULL DEFAULT '',
    endpoint TEXT NOT NULL DEFAULT '',
    node_count INTEGER NOT NULL DEFAULT 0,
    credentials TEXT NOT NULL DEFAULT '{}',      -- encrypted JSON credentials
    metadata TEXT NOT NULL DEFAULT '{}',         -- additional provider-specific metadata
    last_synced_at DATETIME,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(provider, provider_cluster_id)
);

CREATE TABLE IF NOT EXISTS cluster_node_pools (
    id TEXT PRIMARY KEY DEFAULT (lower(hex(randomblob(16)))),
    cluster_id TEXT NOT NULL,
    name TEXT NOT NULL,
    provider_pool_id TEXT NOT NULL DEFAULT '',
    machine_type TEXT NOT NULL DEFAULT '',
    disk_size_gb INTEGER NOT NULL DEFAULT 100,
    disk_type TEXT NOT NULL DEFAULT 'pd-standard',
    node_count INTEGER NOT NULL DEFAULT 0,
    min_nodes INTEGER NOT NULL DEFAULT 0,
    max_nodes INTEGER NOT NULL DEFAULT 0,
    autoscaling BOOLEAN NOT NULL DEFAULT 0,
    preemptible BOOLEAN NOT NULL DEFAULT 0,
    spot_instance BOOLEAN NOT NULL DEFAULT 0,
    status TEXT NOT NULL DEFAULT 'pending',     -- pending, provisioning, running, error, deleting
    k8s_version TEXT NOT NULL DEFAULT '',
    labels TEXT NOT NULL DEFAULT '{}',
    taints TEXT NOT NULL DEFAULT '[]',
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (cluster_id) REFERENCES clusters(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS cluster_events (
    id TEXT PRIMARY KEY DEFAULT (lower(hex(randomblob(16)))),
    cluster_id TEXT NOT NULL,
    event_type TEXT NOT NULL,                   -- 'node_pool_created', 'node_pool_scaled', 'sync', 'error', etc.
    message TEXT NOT NULL DEFAULT '',
    metadata TEXT NOT NULL DEFAULT '{}',
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (cluster_id) REFERENCES clusters(id) ON DELETE CASCADE
);

CREATE INDEX IF NOT EXISTS idx_clusters_provider ON clusters(provider);
CREATE INDEX IF NOT EXISTS idx_clusters_status ON clusters(status);
CREATE INDEX IF NOT EXISTS idx_cluster_node_pools_cluster_id ON cluster_node_pools(cluster_id);
CREATE INDEX IF NOT EXISTS idx_cluster_events_cluster_id ON cluster_events(cluster_id);
CREATE INDEX IF NOT EXISTS idx_cluster_events_created_at ON cluster_events(created_at);
