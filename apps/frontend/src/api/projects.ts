import { type Project } from "../types/Project";
import { useApi } from "./useApi";

export function useProjectsApi() {
    const api = useApi();
    return {
        list: () => api.get<{ projects: Project[] }>("/projects"),
        create: (body: {name: string; description?: string}) =>
            api.post<Project>("/projects", body),
        get: (projectId: string) => 
            api.get<Project>(`/projects/${projectId}`),
        update: (projectId: string, body: {name?: string, description?: string}) => 
            api.patch(`/projects/${projectId}`, body),
        remove: (projectId: string) =>
            api.delete<void>(`/projects/${projectId}`),
    }
}