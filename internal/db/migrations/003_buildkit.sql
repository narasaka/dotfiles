-- 003_buildkit.sql
-- Migrate from Kaniko to BuildKit

ALTER TABLE builds RENAME COLUMN kaniko_job_name TO build_job_name;

UPDATE settings SET key = 'buildkit_addr', value = 'tcp://kubeploy-buildkitd:1234'
    WHERE key = 'kaniko_image';
