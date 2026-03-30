import { create } from 'zustand'

export interface App {
  id: string
  name: string
  display_name: string
  git_url: string
  git_branch: string
  git_subpath: string
  dockerfile_path: string
  registry_image: string
  namespace: string
  replicas: number
  port: number
  env_vars: string
  auto_deploy: boolean
  webhook_secret: string
  ingress_host: string
  ingress_tls: boolean
  status: string
  current_build_id: string | null
  created_at: string
  updated_at: string
}

export interface Build {
  id: string
  app_id: string
  commit_sha: string
  commit_message: string
  commit_author: string
  image_tag: string
  status: string
  kaniko_job_name: string
  logs: string
  started_at: string | null
  finished_at: string | null
  created_at: string
}

export interface Deployment {
  id: string
  app_id: string
  build_id: string
  k8s_deployment_name: string
  replicas_desired: number
  replicas_ready: number
  status: string
  rolled_back_to: string | null
  created_at: string
}

interface AppState {
  apps: App[]
  setApps: (apps: App[]) => void
  updateApp: (app: App) => void
  removeApp: (id: string) => void
}

export const useAppStore = create<AppState>((set) => ({
  apps: [],
  setApps: (apps) => set({ apps }),
  updateApp: (app) =>
    set((state) => ({
      apps: state.apps.map((a) => (a.id === app.id ? app : a)),
    })),
  removeApp: (id) =>
    set((state) => ({
      apps: state.apps.filter((a) => a.id !== id),
    })),
}))
