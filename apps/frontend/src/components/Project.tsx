import { useLocation } from "react-router";
import type { Project } from "../types/Project";

interface ProjectState {
    project: Project
}

export function Project() {
    const location = useLocation();
    const state: ProjectState = location.state || {}

    return (
        <>
            <h1>{state.project.name}</h1>
            <p>{state.project.description}</p>
        </>
    )
}