import React, { useEffect, useState } from "react";
import { type Project } from "../types/Project";
import { useNavigate } from "react-router";
import { useProjectsApi } from "../api/projects";
import "./Dashboard.css"
import { ProjectLink } from "./ProjectLink";

function Dashboard() {
    const [projects, setProjects] = useState<Project[]>([])
    const [newProjectFormData, setNewProjectFormData] = useState({name: "", description: ""})
    const navigate = useNavigate();
    const projectApi = useProjectsApi();

    const handleChange = (e: React.ChangeEvent<HTMLInputElement>) => {
        const { name, value } = e.target;
        setNewProjectFormData((prev) => ({
            ...prev,
            [name]: value,
        }));
    };

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

    const toggleCreateDialog = () => {
        const menu = document.getElementById('dashboard-create-dialog') as HTMLDialogElement
        menu.showModal()
    }

    const closeCreateDialog = () => {
        const menu = document.getElementById('dashboard-create-dialog') as HTMLDialogElement
        menu.close()
        setNewProjectFormData({name: "", description: ""})
    }

    return ( 
        <div className="dashboard">
            <header className="dashboard-header-bar">
                <h3>Projects</h3>
                <button onClick={toggleCreateDialog}>Create</button>
            </header>
            <div className="dashboard-container">
                <div className="dashboard-list-container">
                    {projects?.map((project) => (<ProjectLink project={project} reloadProjects={loadProjects} />))}
                </div>
            </div>

            <dialog id="dashboard-create-dialog">
                <button onClick={closeCreateDialog}>&times;</button>
                <form id="dashboard-create-form" onSubmit={addProject}>
                    <input id="dashboard-name-input"
                        value={newProjectFormData.name}
                        onChange={handleChange}
                        type="text"
                        name="name" 
                        placeholder="Enter Your Projects Name."/>
                    <input id="dashboard-description-input"
                        value={newProjectFormData.description}
                        onChange={handleChange}
                        type="text"
                        name="description"
                        placeholder="Describe Your Project." />
                    <button type="submit">Create</button>
                </form>
            </dialog>
        </div>
    )
}

export default Dashboard;