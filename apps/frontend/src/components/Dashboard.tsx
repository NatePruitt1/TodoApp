import { useEffect, useState } from "react";
import { type Project } from "../types/Project";
import { useProjectsApi } from "../api/projects";
import "./Dashboard.css"
import { ProjectLink } from "./ProjectLink";
import { Header } from "./Header";

function Dashboard() {
    const [projects, setProjects] = useState<Project[]>([])
    const projectApi = useProjectsApi();

    const loadProjects = async () => {
        try {
            const projects = await projectApi.list()
            setProjects(projects.projects)
        } catch(error) {
            console.error(error)
            alert("Failed to load projects.")
        }
    }

    useEffect(() => {
        loadProjects();
    }, [])

    return ( 
        <div className="dashboard">
            <Header title="Todo - Projects" />
            <div className="dashboard-container">
                <div className="dashboard-list-container">
                    {projects?.map((project) => (<ProjectLink key={project.id} project={project} reloadProjects={loadProjects} />))}
                </div>
            </div>

            
        </div>
    )
}

export default Dashboard;