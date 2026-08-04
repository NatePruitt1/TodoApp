import { useNavigate } from "react-router";
import type { Project } from "../types/Project";
import { useProjectsApi } from "../api/projects";
import trashUrl from "../../public/material-symbols--delete-outline.svg"
import arrowUrl from "../../public/material-symbols--arrow-menu-open.svg"

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
            <h3>{project.name}</h3>
            <div className="project-list-item-inner">
                <p>{project.description}</p>
                <div className="project-list-item-inner">
                    <button className="small-icon-button" onClick={handleDelete}><img src={trashUrl} /></button>
                    <button id="project-list-arrow" className="small-icon-button" onClick={handleClick}><img src={arrowUrl} /></button>
                </div>
            </div>
        </div>
    )
}