import { useNavigate } from "react-router";
import type { Project } from "../types/Project";
import { useProjectsApi } from "../api/projects";
import trashUrl from "../../public/icons8-trash-32.png"
import arrowUrl from "../../public/icons8-arrow-30.png"

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
                    <button className="small-icon-button" onClick={handleDelete}><img className="medium-icon" src={trashUrl} /></button>
                    <button id="project-list-arrow" className="small-icon-button" onClick={handleClick}><img className="medium-icon" src={arrowUrl} /></button>
                </div>
            </div>
        </div>
    )
}