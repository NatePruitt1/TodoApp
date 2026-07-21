import { useNavigate } from "react-router";
import type { Project } from "../types/Project";
import { useProjectsApi } from "../api/projects";

export function ProjectLink({project, reloadProjects}: {project: Project, reloadProjects: {(): void}}) {
    const navigate = useNavigate();
    const projectApi = useProjectsApi();

    const handleClick = () => {
        navigate("/projects", {state: {project: project}})
    }

    const handleDelete = async () => {
        try{
            await projectApi.remove(project.id)
            reloadProjects()
        } catch(error) {
            console.error(error)
            alert("Failed to delete project.")
        }
    }
    
    return (
        <div className="project-list-item" >
            <p>{project.name}</p>
            <div className="project-list-item-inner">
                <p>{project.description}</p>
                <div className="project-list-item-inner">
                    <button onClick={handleClick}>Go</button>
                    <button onClick={handleDelete}>Delete</button>
                </div>
            </div>
        </div>
    )
}