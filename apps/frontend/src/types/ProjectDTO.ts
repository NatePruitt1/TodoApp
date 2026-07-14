import type { Project } from "./Project";

export interface AllProjectsDTO {
    success: string;
    data: {
        projects: Project[]
    }
}