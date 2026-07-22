import React, { useEffect, useState } from "react";
import { type Project } from "../types/Project";
import { useNavigate } from "react-router";
import { useProjectsApi } from "../api/projects";
import "./Dashboard.css"
import { ProjectLink } from "./ProjectLink";

import personUrl from "../../public/material-symbols--person.svg"
import addUrl from "../../public/material-symbols--add.svg"
import homeUrl from "../../public/material-symbols--home-outline.svg"
import { Header } from "./Header";

function Dashboard() {
    const [projects, setProjects] = useState<Project[]>([])
    const [newProjectFormData, setNewProjectFormData] = useState({name: "", description: ""})
    const navigate = useNavigate();
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

    const addProject = async (event: React.SubmitEvent<HTMLFormElement>) => {
        event.preventDefault();
        try {
            const project = await projectApi.create({name: newProjectFormData.name, description: newProjectFormData.description})
            setNewProjectFormData({name: "", description: ""})
            navigate("/projects", {state: {project: project}})
        } catch (e) {
            console.error(e)
            alert("Failed to create project")
        }
    }

   

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